# xiaohui

1. 先修改版本号
2. `go run main.go`
3. push

```
git tag -a v0.0.7 -m 'v0.0.7' -f; git push origin tag v0.0.7 -f
```

> xiaohui.service

```sh
[Unit]
Description=wutong
After=network.target

[Service]
ExecStart=/data/xiaohui/xiaohui cron
ExecStop=/bin/pkill xiaohui
RestartSec=3
Restart=always

[Install]
WantedBy=multi-user.target
```
