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
	sshKey, err := os.ReadFile(d.conf.Config.XiaoHui.SSHPath + "/id_rsa")
	if err != nil {
		d.log.Sugar().Errorf("ReadFile err %v", err)
	}

	auth, err := ssh.NewPublicKeys("git", []byte(sshKey), "")
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
