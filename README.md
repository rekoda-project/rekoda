# 📼 rekōdā - Automatic Twitch Recorder
 > Rekōdā, レコダ - means <i>Recorder</i> in Japanese

Is a CLI tool inspired by [```streamlink```](https://github.com/streamlink/streamlink) but with more capabilities. Mainly used to run in the background.
Just _run'n'forget_. It endlessly watches multiple twitch.tv channels you added to come online and starts recording/downloading them simultaneously into local MPEGTS ```.ts``` files which can be played via Video Player

Your best new friend for your Kappa needs

# 🤩 It's awesome!  Here's why
- 💪 **Powerful**. Record **multiple** live streams simultaneously. It's like ```streamlink``` on steroids
- 🤏 **Small**. Single **~10 MB** binary file solution
- 0️⃣  **Independent**. Zero run-time dependency. Does not depends on ```ffmpeg``` or ```streamlink```
- 🧠 **Smart**. Avoids making lot of files. If stream goes down, it waits and tries to reconnect to continue writing into the same file.
- ⚙ **Efficient**. Low CPU and RAM usage. ~20 MB average RAM usage
- 🚀 **Crossplatform**. Supports ```Windows```, ```Linux```, ```macOS```, ```Raspberry Pi (arm)``` and other OS & platforms

# 📦 Installation
You need ```go``` installed. Get it here ➡ [golang.org](https://golang.org/)
```console
go install github.com/wmw9/rekoda@latest
```

# 🔬 Basic usage 
```console
wmw@ubuntu:~$ rekoda channel add sodapoppin rwxrob
[INFO] [INIT] Log level: info (default)
[INFO] [INIT] Using default config file path: /home/wmw/rekoda/rekoda.toml
[INFO] [CLI] Channel added 'sodapoppin' in config file
[INFO] [CLI] Channel added 'rwxrob' in config file
wmw@ubuntu:~$ rekoda rec
[INFO] [REC] Recorder starting
[INFO] [INIT] Log level: info (default)
[INFO] [INIT] Using default config file path: /home/wmw/rekoda/rekoda.toml
[INFO] [REC] [sodapoppin] Trying to get m3u8 live playlist
[INFO] [REC] [sodapoppin] 🤩 Went online!
[INFO] [REC] [sodapoppin] Opening stream: best quality
[INFO] [REC] [sodapoppin] [sodapoppin_2021-09-08_12-57-04.ts] Writing stream to file: /home/wmw/rekoda/streams/sodapoppin/sodapoppin_2021-09-08_12-57-04.ts
[INFO] [REC] [SEG] [DOWNLOAD] [sodapoppin] [sodapoppin_2021-09-08_12-57-04.ts] Written 879 kB (2 seconds)
[INFO] [REC] [SEG] [DOWNLOAD] [sodapoppin] [sodapoppin_2021-09-08_12-57-04.ts] Written 2.3 MB (4 seconds)
[INFO] [REC] [rwxrob] Trying to get m3u8 live playlist
[INFO] [REC] [rwxrob] 🤩 Went online!
[INFO] [REC] [rwxrob] Opening stream: best quality
[INFO] [REC] [rwxrob] [rwxrob_2021-09-08_12-57-06.ts] Writing stream to file: /home/wmw/rekoda/streams/rwxrob/rwxrob_2021-09-08_12-57-06.ts
[INFO] [REC] [SEG] [DOWNLOAD] [rwxrob] [rwxrob_2021-09-08_12-57-06.ts] Written 830 kB (2 seconds)
[INFO] [REC] [SEG] [DOWNLOAD] [sodapoppin] [sodapoppin_2021-09-08_12-57-04.ts] Written 3.7 MB (6 seconds)
[INFO] [REC] [SEG] [DOWNLOAD] [rwxrob] [rwxrob_2021-09-08_12-57-06.ts] Written 2.1 MB (4 seconds)
^C [INFO] [CLI] Interrupted! SIGTERM signal. <Ctrl>+<C> pressed. Graceful shutdown...
wmw@ubuntu:~$
```

# 🤝 Contributing
Contributions, issues and feature requests are welcome! 👍 <br>
Feel free to check [open issues](https://github.com/wmw9/rekoda/issues).

## 🌟 Show your support 
Give a ⭐️ if this project helped you!

# 📝 ToDo
- [x] Record multiple streams simultaneously
- [ ] Daemon mode with  ```systemd``` support
- [ ] Record chat history
- [ ] Progress Bars
- [ ] Download VoDs (past broadcasts) and clips capabilities. ```rekoda download``` command
- [ ] Discord, Telegram reports
- [ ] Packages for apt, dnf, pacman, choco, brew etc package managers
- [ ] Dockerfile one-liner

# 🧠 What I Learned
- Goroutines and channels
- Advanced http.Client usage with graceful degradation
- Golang project layout best practices
- Making CLI applications using Cobra package
- Importance of log levels: trace, debug, info
- Writing tests
- TOML file usage

# 📑 License 
(c) 2021 Ivan Smyshlyaev. [MIT License](https://tldrlegal.com/license/mit-license)
