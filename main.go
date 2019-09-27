package main

import (
	"flag"
	"fmt"
	"github.com/winterssy/easylog"
	"github.com/xiaomLee/music-get/conf"
	"github.com/xiaomLee/music-get/handler"
	"github.com/xiaomLee/music-get/provider/qq"
	"github.com/xiaomLee/music-get/utils"
	"path/filepath"
	"strings"
)

func main() {
	flag.Parse()

	if conf.Singer == "" && conf.Url == "" {
		easylog.Fatal("Missing required args. " +
			"get music by singer name pls use: music-get -singer \"Michael Jackson\" . " +
			"get music by url pls use: music-get -url http://y.qq.com/album/xxxx.html")
	}

	if err := conf.Init(); err != nil {
		easylog.Fatal(err)
	}

	if conf.Singer != "" {
		easylog.Info(fmt.Sprintf("start get music by singer %s", conf.Singer))
		urlList := getAlbumUrlList(conf.Singer)
		for _, url := range urlList {
			getMusicByUrl(url)
		}
	}

	if conf.Url != "" {
		easylog.Info(fmt.Sprintf("start get music by url %s", conf.Url))
		getMusicByUrl(conf.Url)
	}

	if err := conf.Conf.Save(); err != nil {
		easylog.Errorf("Save config failed: %s", err.Error())
	}
}

func getAlbumUrlList(name string) []string {
	//获取姓名首字母,若是中文转成拼音
	indexChar := utils.GetIndexChar(name)
	if indexChar == "" {
		return nil
	}
	indexChar = strings.ToLower(indexChar)
	singerConf := filepath.Join(conf.ConfigPath, fmt.Sprintf("qq/%s-singer.json", indexChar))
	if exist, _ := utils.ExistsPath(singerConf); !exist {
		if err := qq.FetchAllSingerAndStore(indexChar, singerConf); err != nil {
			easylog.Error(err)
			return nil
		}
	}
	singerMid := qq.GetSingerMidByName(singerConf, name)
	if singerMid == "" {
		easylog.Error(fmt.Sprintf("not fund singer %s", name))
		return nil
	}
	return qq.GenerateAlbumUrl(singerMid)
}

func getMusicByUrl(url string) {
	req, err := handler.Parse(url)
	if err != nil {
		easylog.Fatal(err)
	}

	if req.RequireLogin() {
		easylog.Info("Unauthorized, please login")
		if err = req.Login(); err != nil {
			easylog.Fatalf("Login failed: %s", err.Error())
		}
		easylog.Info("Login successful")
	}

	if err = req.Do(); err != nil {
		easylog.Fatal(err)
	}

	mp3List, err := req.Prepare()
	if err != nil {
		easylog.Fatal(err)
	}

	if len(mp3List) == 0 {
		return
	}

	n := conf.Conf.ConcurrentDownloadTasksCount
	switch {
	case n > 1:
		handler.ConcurrentDownload(mp3List, n)
	default:
		handler.SingleDownload(mp3List)
	}
}
