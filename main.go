package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

type audio struct {
	Artist string `json:"artist"`
	Name   string `json:"name"`
	Url    string `json:"url"`
}

func main() {
	dir, _ := os.ReadDir("./mp3")
	files := []*audio{}

	for _, d := range dir {
		fmt.Println(d.Name()[9:(len(d.Name()) - 4)])
		files = append(files, &audio{
			Artist: d.Name()[0:8],
			Name:   d.Name()[9:len(d.Name())],
			Url:    "https://cdn.jsdelivr.net/gh/isfk/xiaohui/" + url.QueryEscape(d.Name()),
		})
	}

	vd, _ := json.Marshal(files)
	fmt.Println(string(vd))
	_ = os.WriteFile("./docs/list.js", vd, 666)
	fmt.Println("ok")
}
