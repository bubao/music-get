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
	ConfigPath string
	Conf       = &Config{}
	Singer     string
	Url        string
)

type Config struct {
	Cookies                      []*http.Cookie `json:"cookies,omitempty"`
	Workspace                    string         `json:"-"`
	DownloadDir                  string         `json:"-"`
	DownloadSubDir               string         `json:"-"`
	DownloadOverwrite            bool           `json:"-"`
	ConcurrentDownloadTasksCount int            `json:"-"`
}

var (
	downloadOverwrite            bool
	concurrentDownloadTasksCount int
)

func init() {
	flag.BoolVar(&downloadOverwrite, "f", false, "overwrite already downloaded music")
	flag.IntVar(&concurrentDownloadTasksCount, "n", 5, "concurrent download tasks count, max 16")
	flag.StringVar(&Singer, "singer", "", "get music by singer name, otherwise get music by url")
	flag.StringVar(&Url, "url", "", "get music by url, otherwise get music by singer name")
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

	ConfigPath = filepath.Join(pwd, "config")
	if err = utils.BuildPathIfNotExist(ConfigPath); err != nil {
		return err
	}

	if err = utils.BuildPathIfNotExist(filepath.Join(ConfigPath, "qq")); err != nil {
		return err
	}
	if err = utils.BuildPathIfNotExist(filepath.Join(ConfigPath, "netease")); err != nil {
		return err
	}

	downloadDir := filepath.Join(pwd, "downloads")
	if err = utils.BuildPathIfNotExist(downloadDir); err != nil {
		return err
	}

	if err := load(filepath.Join(ConfigPath, "main.json")); err != nil {
		easylog.Warn("may you run this first time")
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
	path := filepath.Join(ConfigPath, "main.json")
	return ioutil.WriteFile(path, data, 0644)
}
