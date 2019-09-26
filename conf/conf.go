package conf

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/winterssy/easylog"
	"github.com/xiaomLee/music-get/utils"
)

const (
	MaxConcurrentDownloadTasksCount = 16
	DefaultDownloadBr               = 128
)

var (
	confPath string
	Conf     = &Config{}
)

type Config struct {
	Cookies                      []*http.Cookie `json:"cookies,omitempty"`
	Workspace                    string         `json:"-"`
	DownloadDir                  string         `json:"-"`
	DownloadOverwrite            bool           `json:"-"`
	ConcurrentDownloadTasksCount int            `json:"-"`
	RunMode string            `json:"-"`
}

var (
	downloadOverwrite            bool
	concurrentDownloadTasksCount int
	runMode string
)

func init() {
	flag.BoolVar(&downloadOverwrite, "f", false, "overwrite already downloaded music")
	flag.IntVar(&concurrentDownloadTasksCount, "n", 5, "concurrent download tasks count, max 16")
	flag.StringVar(&runMode, "m", "singer", "decide run mode; default singer, get music by singer name; otherwise url,get music by url")
}

func Init() error {
	if concurrentDownloadTasksCount < 1 || concurrentDownloadTasksCount > MaxConcurrentDownloadTasksCount {
		easylog.Warn("Invalid n parameter, use default value")
		concurrentDownloadTasksCount = 1
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	confPath = filepath.Join(pwd, "music-get.json")
	if err = load(confPath); err != nil {
		easylog.Warn("Load config file failed, you may run for the first time")
	}

	downloadDir := filepath.Join(pwd, "downloads")
	if err = utils.BuildPathIfNotExist(downloadDir); err != nil {
		return err
	}

	Conf.Workspace = pwd
	Conf.DownloadDir = downloadDir
	Conf.DownloadOverwrite = downloadOverwrite
	Conf.ConcurrentDownloadTasksCount = concurrentDownloadTasksCount
	return nil
}

func load(confPath string) error {
	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &Conf)
}

func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(confPath, data, 0644)
}
