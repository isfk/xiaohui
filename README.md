# xiaohui

1. 先修改版本号
2. `go run main.go`
3. push

> xiaohui.service

```sh
[Unit]
Description=xiaohui
After=network.target

[Service]
ExecStart=/data/xiaohui/xiaohui cron
ExecStop=/bin/pkill xiaohui
RestartSec=3
Restart=always
User=root
Group=root

[Install]
WantedBy=multi-user.target
```

```sh
yt-dlp -f 'ba' -x --audio-format mp3 --playlist-start 1 --playlist-end 10 https://www.youtube.com/@eastsinglecom/videos -o '%(upload_date)s_%(title)s.mp3'
```

```sh
git tag -a v0.0.7 -m 'v0.0.7' -f; git push origin tag v0.0.7 -f
```
