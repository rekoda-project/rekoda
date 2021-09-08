# 📼 rekōdā - Automatic Twitch Recorder
 > Rekōdā, レコダ - means <i>Recorder</i> in Japanese

Your best new friend for your Kappa needs

# 🤩 It's awesome!  Here's why
- 💪 **Powerful**. Record **multiple** live streams at once. It's like ```streamlink``` on steroids
- 🤏 **Small**. Single **~10 MB** binary file solution
- 0️⃣  **Independent**. Zero run-time dependency. Does not depends on ```ffmpeg``` or ```streamlink```
- 🧠 **Smart**. Avoids making lot of files. If stream goes down, it waits and tries to reconnect to continue writing into the same file.
- ⚙ **Efficient**. Low CPU and RAM usage. ~20 MB average RAM usage
- 🚀 **Crossplatform**. Supports ```Windows```, ```Linux```, ```macOS```, ```Rasberry Pi (ARM)``` and other OS & platforms

# 📦 Installation
You need ```go``` installed. Get it here ➡ [golang.org](https://golang.org/)
```console
go install github.com/wmw9/rekoda@latest
```

# 🔬 Basic usage 
```console
wmw@ubuntu:~$ rekoda channel add sodapoppin rwxrob
[2021 Sep 8, Wednesday][12:56:50 MSK] [INFO] [INIT] Log level: info (default)
[2021 Sep 8, Wednesday][12:56:50 MSK] [INFO] [INIT] Using default config file path: /home/wmw/rekoda/rekoda.toml
[2021 Sep 8, Wednesday][12:56:50 MSK] [INFO] [CLI] Channel added 'sodapoppin' in config file
[2021 Sep 8, Wednesday][12:56:50 MSK] [INFO] [CLI] Channel added 'rwxrob' in config file
wmw@ubuntu:~$ rekoda rec
[2021 Sep 8, Wednesday][12:57:02 MSK] [INFO] [REC] Recorder starting
[2021 Sep 8, Wednesday][12:57:02 MSK] [INFO] [INIT] Log level: info (default)
[2021 Sep 8, Wednesday][12:57:02 MSK] [INFO] [INIT] Using default config file path: /home/wmw/rekoda/rekoda.toml
[2021 Sep 8, Wednesday][12:57:04 MSK] [INFO] [REC] [sodapoppin] Trying to get m3u8 live playlist
[2021 Sep 8, Wednesday][12:57:04 MSK] [INFO] [REC] [sodapoppin] 🤩 Went online!
[2021 Sep 8, Wednesday][12:57:04 MSK] [INFO] [REC] [sodapoppin] Opening stream: best quality
[2021 Sep 8, Wednesday][12:57:04 MSK] [INFO] [REC] [sodapoppin] [sodapoppin_2021-09-08_12-57-04.ts] Writing stream to file: /home/wmw/rekoda/streams/sodapoppin/sodapoppin_2021-09-08_12-57-04.ts
[2021 Sep 8, Wednesday][12:57:05 MSK] [INFO] [REC] [SEG] [DOWNLOAD] [sodapoppin] [sodapoppin_2021-09-08_12-57-04.ts] Written 879 kB (2 seconds)
[2021 Sep 8, Wednesday][12:57:05 MSK] [INFO] [REC] [SEG] [DOWNLOAD] [sodapoppin] [sodapoppin_2021-09-08_12-57-04.ts] Written 2.3 MB (4 seconds)
[2021 Sep 8, Wednesday][12:57:05 MSK] [INFO] [REC] [rwxrob] Trying to get m3u8 live playlist
[2021 Sep 8, Wednesday][12:57:06 MSK] [INFO] [REC] [rwxrob] 🤩 Went online!
[2021 Sep 8, Wednesday][12:57:06 MSK] [INFO] [REC] [rwxrob] Opening stream: best quality
[2021 Sep 8, Wednesday][12:57:06 MSK] [INFO] [REC] [rwxrob] [rwxrob_2021-09-08_12-57-06.ts] Writing stream to file: /home/wmw/rekoda/streams/rwxrob/rwxrob_2021-09-08_12-57-06.ts
[2021 Sep 8, Wednesday][12:57:06 MSK] [INFO] [REC] [SEG] [DOWNLOAD] [rwxrob] [rwxrob_2021-09-08_12-57-06.ts] Written 830 kB (2 seconds)
[2021 Sep 8, Wednesday][12:57:06 MSK] [INFO] [REC] [SEG] [DOWNLOAD] [sodapoppin] [sodapoppin_2021-09-08_12-57-04.ts] Written 3.7 MB (6 seconds)
[2021 Sep 8, Wednesday][12:57:06 MSK] [INFO] [REC] [SEG] [DOWNLOAD] [rwxrob] [rwxrob_2021-09-08_12-57-06.ts] Written 2.1 MB (4 seconds)
^C[2021 Sep 8, Wednesday][12:57:07 MSK] [INFO] [CLI] Interrupted! SIGTERM signal. <Ctrl>+<C> pressed. Graceful shutdown...
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
