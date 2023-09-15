package jobs

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/isfk/xiaohui/config"
	"go.uber.org/zap"
)

type AudioJob struct {
	log  *zap.Logger
	conf *config.Conf
}

func NewAudioJob(log *zap.Logger, conf *config.Conf) *AudioJob {
	return &AudioJob{
		log:  log,
		conf: conf,
	}
}

func (d *AudioJob) Run() {
	tag := fmt.Sprintf("v%d", time.Now().Unix())
	d.log.Info("audio cron start", zap.String("time", time.Now().Format("15:04:05")))
	// 下载
	// yt-dlp -f 'ba' -x --audio-format mp3 https://www.youtube.com/@eastsinglecom/videos -o '%(upload_date)s_%(title)s.mp3'
	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "-x", "--audio-format", "mp3", "https://www.youtube.com/@eastsinglecom/videos", "-o", d.conf.Config.XiaoHui.AudioPath+"/%(upload_date)s_%(title)s.mp3")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		d.log.Error(err.Error())
		return
	}
	cmd.Stderr = cmd.Stdout

	if err != nil {
		d.log.Error(err.Error())
		return
	}

	err = cmd.Start()

	if err != nil {
		d.log.Error(err.Error())
		return
	}

	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		d.log.Info(string(tmp))
		if err != nil {
			break
		}
	}

	if err = cmd.Wait(); err != nil {
		d.log.Error(err.Error())
		return
	}

	// 生成 list.js
	type audio struct {
		Artist string `json:"artist"`
		Name   string `json:"name"`
		Url    string `json:"url"`
	}

	d.log.Info("audio", zap.String("time", time.Now().Format("15:04:05")))
	dir, _ := os.ReadDir(d.conf.Config.XiaoHui.AudioPath)
	files := []*audio{}

	for _, d := range dir {
		files = append(files, &audio{
			Artist: d.Name()[0:8],
			Name:   d.Name()[9:len(d.Name())],
			Url:    "https://cdn.jsdelivr.net/gh/isfk/xiaohui@" + tag + "/mp3/" + url.QueryEscape(d.Name()),
		})
	}

	vd, _ := json.Marshal(files)
	_ = os.WriteFile(d.conf.Config.XiaoHui.AudioPath+"/list.js", vd, 0o666)

	// git
	gitCmd1 := exec.Command("/usr/bin/git ", "add", "./")
	_ = gitCmd1.Run().Error()

	gitCmd2 := exec.Command("/usr/bin/git ", "commit", "-am", "update")
	_ = gitCmd2.Run().Error()

	gitCmd3 := exec.Command("/usr/bin/git ", "push")
	_ = gitCmd3.Run().Error()

	gitCmd4 := exec.Command("/usr/bin/git", "tag", "-a", tag, "-m '"+tag+"'")
	_ = gitCmd4.Run().Error()

	gitCmd5 := exec.Command("/usr/bin/git", "push", "origin", "tag", tag, "-f")
	_ = gitCmd5.Run().Error()

	d.log.Info("audio cron done.")
}

func (*AudioJob) String() string {
	return "audio"
}
