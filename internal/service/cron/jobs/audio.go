package jobs

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
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
	// yt-dlp -f 'ba' -x --audio-format mp3 --playlist-start 1 --playlist-end 10 https://www.youtube.com/@eastsinglecom/videos -o '%(upload_date)s_%(title)s.mp3'
	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "-x", "--audio-format", "mp3", "--playlist-start", "1", "--playlist-end", "10", "https://www.youtube.com/@eastsinglecom/videos", "-o", d.conf.Config.XiaoHui.AudioPath+"/%(upload_date)s_%(title)s.mp3")

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

	for _, f := range dir {
		if len(f.Name()) < 8 {
			d.log.Error("错除了", zap.String("name", f.Name()))
			continue
		}
		files = append(files, &audio{
			Artist: f.Name()[0:8],
			Name:   f.Name()[9:len(f.Name())],
			Url:    "https://cdn.jsdelivr.net/gh/isfk/xiaohui@" + tag + "/mp3/" + url.QueryEscape(f.Name()),
		})
	}

	vd, _ := json.Marshal(files)
	_ = os.WriteFile(d.conf.Config.XiaoHui.DocsPath+"/list.js", vd, 0o666)

	// git
	auth, err := ssh.NewPublicKeys("git", []byte(`-----BEGIN OPENSSH PRIVATE KEY-----
	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
	-----END OPENSSH PRIVATE KEY-----`), "")
	if err != nil {
		d.log.Sugar().Errorf("NewPublicKeys err %v", err)
		return
	}

	r, err := git.PlainOpen(d.conf.Config.XiaoHui.Path)
	if err != nil {
		d.log.Sugar().Errorf("PlainOpen err %v", err)
		return
	}

	w, err := r.Worktree()
	if err != nil {
		d.log.Sugar().Errorf("Worktree err %v", err)
		return
	}

	err = w.AddGlob("./")
	if err != nil {
		d.log.Sugar().Errorf("AddGlob err %v", err)
		return
	}

	_, err = w.Commit(tag, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "isfk",
			Email: "isfk@live.cn",
			When:  time.Now(),
		},
	})
	if err != nil {
		d.log.Sugar().Errorf("Commit err %v", err)
		return
	}

	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		Auth:       auth,
	})
	if err != nil {
		d.log.Sugar().Errorf("Push err %v", err)
		return
	}

	h, err := r.Head()
	if err != nil {
		d.log.Sugar().Errorf("Head err %v", err)
		return
	}

	_, err = r.CreateTag(tag, h.Hash(), &git.CreateTagOptions{})
	if err != nil {
		d.log.Sugar().Errorf("CreateTag err %v", err)
		return
	}

	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		RefSpecs:   []gitConfig.RefSpec{gitConfig.RefSpec("refs/tags/*:refs/tags/*")},
		Auth:       auth,
	})

	if err != nil {
		d.log.Sugar().Errorf("Push Tag err %v", err)
		return
	}

	d.log.Info("audio cron done.")
}

func (*AudioJob) String() string {
	return "audio"
}
