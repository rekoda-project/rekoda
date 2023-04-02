package recorder

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/grafov/m3u8"
	"github.com/hako/durafmt"
	lru "github.com/hashicorp/golang-lru"
	log "github.com/sirupsen/logrus"
	"github.com/rekoda-project/rekoda/internal/config"
	"github.com/wmw64/twitchpl"
)

var (
	USER_AGENT      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:86.0) Gecko/20100101 Firefox/86.0"
	backoffSchedule = []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second, 10 * time.Second}
)

type Recorder struct {
	Segments Segment
	Online   []string
	Client   *http.Client
}

type Segment struct {
	URI           string
	totalDuration time.Duration
}

// New is used to create Recorder object and initializing http.Client
func New() *Recorder {
	return &Recorder{
		Client: &http.Client{Timeout: 5 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second},
		},
	}
}

// Start is main infinite loop function used to start recording channels and checking streams to come out live
func Start() {
	ctxLog := log.WithField("general", "REC")
	ctxLog.Info("Recorder starting")

	c := config.InitConfig() // Initialize config
	r := New()               // Initialize recorder

	if i := c.GetRecAmountChannels(); i < 1 {
		ctxLog.Error("There's nothing to record, try to add and enable channels for recording. Type 'rekoda channel' for more info.")
		return
	}

	// Capture <Ctrl>+<C>
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		cleanup()
		os.Exit(1)
	}()

	// Main cycle where all the magic happens ✨
	for {
		ctxLog.Debugf("Channels being recorded right now: %v", r.Online)
		for _, u := range c.Channels {
			defer recoverFromPanic()
			ctxLog.Tracef("Checking %v", u.User)
			if u.Enabled && !r.IsOnline(u.User) {
				cLog := ctxLog.WithField("channel", u.User)
				cLog.Info("Trying to get m3u8 live playlist")
				pl, err := twitchpl.Get(context.Background(), u.User, true)
				if err != nil {
					cLog.Info("Channel is offline or banned.")
					continue
				}
				cLog.Info("🤩 Went online! ")
				switch u.Quality {
				case "best":
					cLog.Info("Opening stream: best quality")
					url := pl.Best().AsURL()
					cLog.Debugf("URL: %v", url)
					r.Rec(cLog, c, u, url)
				case "worst":
					cLog.Info("Opening stream: worst quality")
					url := pl.Worst().AsURL()
					cLog.Debugf("URL: %v", url)
					r.Rec(cLog, c, u, url)
				case "audio_only":
					cLog.Info("Opening stream: audio_only")
					url := pl.Audio().AsURL()
					cLog.Debugf("URL: %v", url)
					r.Rec(cLog, c, u, url)

				}
			}
			time.Sleep(1 * time.Second)
		}
		ctxLog.Infof("Sleep for 1 minute")
		time.Sleep(1 * time.Minute)
	}

}

// Rec is used to create channel's foldera and to compose stream file name
// then executes 2 goroutines with 1 shared channel to send data from one to another
func (r *Recorder) Rec(log *log.Entry, c *config.Config, channel config.Channels, hlsURL string) {
	// Get local time
	now, err := TimeIn(time.Now(), "Local")
	if err != nil {
		log.Error(err)
	}
	log.Tracef("Local time: %v", now)

	// Define file name
	fname := fmt.Sprintf("%v_%v.ts", channel.User, now.Format("2006-01-02_15-04-05"))
	fLog := log.WithField("file", fname)

	// Create channel directory
	sep := string(os.PathSeparator)
	log.Trace("Trying to create channel directory")
	channelDir := c.StreamsDir + sep + channel.User
	if err := os.MkdirAll(channelDir, 0777); err != nil {
		fLog.Error(err)
		return
	}
	filepath := channelDir + sep + fname
	fLog.Debugf("Stream directory: '%s'", channelDir)
	fLog.Infof("Writing stream to file: %v", filepath)

	dlc := make(chan *Segment, 1024)
	go r.GetPlaylist(fLog, channel, hlsURL, dlc)
	go r.DownloadSegment(fLog, filepath, dlc)
}

// DownloadSegment is mainly used as a goroutine which accepts new .ts chunks to be downloaded from GetPlaylist() function and then merges them into local file.
// Also updates and report total duration and bytes of current stream
func (r *Recorder) DownloadSegment(log *log.Entry, filepath string, dlc chan *Segment) {
	defer recoverFromPanic()
	ctxLog := log.WithField("status", "DOWNLOAD").WithField("func", "SEG")
	var totalBytes uint64 = 0
	var bytes int64
	var req *http.Request

	out, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		ctxLog.Error(err)
	}
	defer out.Close()

	for v := range dlc {
		req, err = http.NewRequest("GET", v.URI, nil)
		if err != nil {
			ctxLog.Info(err)
		}
		res, err := r.doRequestWithRetries(req)
		if err != nil {
			ctxLog.Error(err)
			continue
		}

		if res.StatusCode != 200 {
			ctxLog.Errorf("Received HTTP %v for %v\n", res.StatusCode, v.URI)
			res.Body.Close()
			continue
		}
		bytes, err = io.Copy(out, res.Body)
		if err != nil {
			ctxLog.Error(err)
		}
		res.Body.Close()

		totalBytes += uint64(bytes)
		duration := durafmt.Parse(v.totalDuration).LimitFirstN(1).String()
		// speed := uint64(v.totalDuration * time.Millisecond)

		ctxLog.Infof("Written %v (%v)", humanize.Bytes(totalBytes), duration)
	}
}

// GetPlaylist is an infinite loop function used to download m3u8 live playlist file from server
// then it parses and drops old chunks which are already downloaded.
// If new chunks are present they are being sent to DownloadSegment() function via channel to be downloaded.
// New chunks are marked as old after being sent by adding their unique filename in cache.
// When m3u8 live playlist link is expired (usually 24 hours) it tries to refresh it by generating a new one
func (r *Recorder) GetPlaylist(log *log.Entry, channel config.Channels, urlStr string, dlc chan *Segment) {
	r.AddOnline(channel.User)
	defer r.RemoveOnline(channel.User)
	defer recoverFromPanic()

	ctxLog := log.WithField("status", "DOWNLOAD").WithField("func", "GET")

	var recDuration time.Duration = 0
	var req *http.Request

	cache, _ := lru.New(1024)

	playlistUrl, err := url.Parse(urlStr)
	if err != nil {
		ctxLog.Error(err)
	}
	for {
		req, err = http.NewRequest("GET", urlStr, nil)
		if err != nil {
			ctxLog.Error(err)
		}
		ctxLog.Debugf("URL: %v", urlStr[len(urlStr)-10:]) // Change it later

		res, err := r.doRequestWithRetries(req)
		if err != nil {
			ctxLog.Error(err)
			urlStr, err = r.RefreshPlaylist(ctxLog, channel)
			if err != nil {
				ctxLog.Errorf("Failed to refresh m3u8 playlist: '%v'", err)
			}
			time.Sleep(1 * time.Second)
			continue
		}

		playlist, listType, err := m3u8.DecodeFrom(res.Body, true)
		if err != nil {
			ctxLog.Error(err)
			urlStr, err = r.RefreshPlaylist(ctxLog, channel)
			if err != nil {
				ctxLog.Errorf("Failed to refresh m3u8 playlist: '%v'", err)
			}
			time.Sleep(1 * time.Second)
			continue
		}
		res.Body.Close()

		if listType == m3u8.MEDIA {
			mpl := playlist.(*m3u8.MediaPlaylist)
			for _, v := range mpl.Segments {
				if v != nil {
					var msURI string

					if strings.HasPrefix(v.URI, "http") {
						msURI, err = url.QueryUnescape(v.URI)
						if err != nil {
							ctxLog.Fatal(err)
						}
					} else {
						msUrl, err := playlistUrl.Parse(v.URI)
						if err != nil {
							ctxLog.Errorf("Failed to parse m3u8: '%v'", err)
							urlStr, err = r.RefreshPlaylist(ctxLog, channel)
							if err != nil {
								ctxLog.Errorf("Failed to refresh m3u8 playlist: '%v'", err)
							}
							time.Sleep(1 * time.Second)
							continue
						}
						msURI, err = url.QueryUnescape(msUrl.String())
						if err != nil {
							ctxLog.Fatal(err)
						}
					}
					_, hit := cache.Get(msURI)
					if !hit {
						cache.Add(msURI, nil)
						recDuration += time.Duration(int64(v.Duration * 1000000000))
						dlc <- &Segment{msURI, recDuration}
					}
				}
			}
			if mpl.Closed {
				ctxLog.Info("Stream ended. Waiting 10m for stream to come online again before closing file") // Often streamers restart their translation for various reasons
				urlStr, err = r.WaitForRestart(ctxLog, channel)
				if err != nil {
					ctxLog.Infof("Closing file. Reason: '%v'", err)
					close(dlc)
					return
				}
				ctxLog.Info("🚀 Went online again!") // Often streamers restart their translation for various reasons
				continue
			} else {
				time.Sleep(time.Duration(int64(mpl.TargetDuration * 1000000000)))
			}
		} else {
			ml := playlist.(*m3u8.MasterPlaylist)
			ctxLog.Fatalf("Not a valid media playlist: '%v'", ml)
		}
	}
}

// RefreshPlaylist is used to update direct link to m3u8 live playlist
func (r *Recorder) RefreshPlaylist(log *log.Entry, channel config.Channels) (string, error) {
	ctxLog := log.WithField("func", "REFRESH")
	var url string

	ctxLog.Info("Trying to refresh m3u8 playlist")
	pl, err := twitchpl.Get(context.Background(), channel.User, true)
	if err != nil {
		return "", err
	}
	switch channel.Quality {
	case "best":
		url = pl.Best().AsURL()
	case "worst":
		url = pl.Worst().AsURL()
	case "audio_only":
		url = pl.Audio().AsURL()
	}
	ctxLog.Debugf("m3u8 URL updated: %v", url)
	return url, nil
}

// WaitForRestart is used to prevent making multiple stream files by writing stream to the same file in case of channel coming online again in a few minutes
func (r *Recorder) WaitForRestart(log *log.Entry, channel config.Channels) (string, error) {
	ctxLog := log.WithField("func", "WAIT")
	var url string
	var pl *twitchpl.PlaylistManager
	var err error

	for i := 1; i <= 20; i++ {
		ctxLog.Infof("Sleep for 30 seconds. Try %v/20", i)
		time.Sleep(30 * time.Second)
		ctxLog.Info("Checking if channel went online again (restart)")
		pl, err = twitchpl.Get(context.Background(), channel.User, true)
		if err != nil {
			ctxLog.Info("Channel is still offline")
			continue
		}
		switch channel.Quality {
		case "best":
			url = pl.Best().AsURL()
		case "worst":
			url = pl.Worst().AsURL()
		case "audio_only":
			url = pl.Audio().AsURL()
		}
		return url, nil
	}
	return "", errors.New("channel was offline for 10 minutes")
}

// TimeIn is used to get local time
func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

// doRequestWithRetries makes GET request, if failed it retries 3 more times with backoff timer
func (r *Recorder) doRequestWithRetries(req *http.Request) (*http.Response, error) {
	defer recoverFromPanic()
	ctxLog := log.WithField("general", "GET")
	var err error
	var res *http.Response

	// req.Close = true
	// req.Header.Set("Connection", "close") // prevent 'too many open files' error
	req.Header.Set("User-Agent", USER_AGENT)

	for _, backoff := range backoffSchedule {
		res, err = r.Client.Do(req)
		if err == nil {
			break
		}
		ctxLog.Errorf("Request error: '%v' Retrying in %v", err, backoff)
		time.Sleep(backoff)

	}
	if err != nil {
		return res, err
	}
	return res, err
}

func cleanup() {
	log.WithField("general", "CLI").Info("Interrupted! SIGTERM signal. <Ctrl>+<C> pressed. Graceful shutdown...")
}

// IsOnline is used to check Online struct if specified channel is being recorder right now
func (r *Recorder) IsOnline(channel string) bool {
	for _, v := range r.Online {
		if v == channel {
			return true
		}
	}
	return false
}

// AddOnline marks channel name being recorder right now by adding it to Channel struct
func (r *Recorder) AddOnline(u string) {
	r.Online = append(r.Online, u)
}

// RemoveOnline removes channel name from Online struct, usually invokes when stream ends.
func (r *Recorder) RemoveOnline(u string) {
	for i, v := range r.Online {
		if v == u {
			r.Online = append(r.Online[:i], r.Online[i+1:]...)
		}
	}
}

func recoverFromPanic() {
	if r := recover(); r != nil {
		log.Errorf("Recovering from panic: '%v'", r)
	}
}
