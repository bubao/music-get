package qq

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/winterssy/easylog"
	"github.com/xiaomLee/music-get/utils"
	"io/ioutil"
	"math"
	"strings"
)

type SingerListResponse struct {
	Code int            `json:"code"`
	Data SingerListData `json:"data"`
}

type SingerListData struct {
	Area       int        `json:"area"`
	Sex        int        `json:"sex"`
	SingerList SingerList `json:"singerlist"`
	Total      int        `json:"total"`
}

type SingerItem struct {
	Country    string `json:"country"`
	SingerId   int    `json:"singer_id"`
	SingerMid  string `json:"singer_mid"`
	SingerName string `json:"singer_name"`
	SingerPic  string `json:"singer_pic"`
}
type SingerList []SingerItem

type SingerListParam struct {
	Area    int `json:"area"`
	Sex     int `json:"sex"`
	Genre   int `json:"genre"`
	Index   int `json:"index"`
	Sin     int `json:"sin"`
	CurPage int `json:"cur_page"`
}

var (
	singerListApi = "https://u.y.qq.com/cgi-bin/musicu.fcg?-=getUCGI9326033526469224&g_tk=5381&loginUin=0&hostUin=0&format=json&inCharset=utf8&outCharset=utf-8&notice=0&platform=yqq.json&needNewCode=0&data="
	albumListApi  = "https://c.y.qq.com/v8/fcg-bin/fcg_v8_singer_album.fcg?g_tk=5381&loginUin=0&hostUin=0&format=jsonp&inCharset=utf8&outCharset=utf-8&notice=0&platform=yqq&needNewCode=0&order=time&begin=0&num=300&exstatus=1&singermid="
)

func GetSingerListData(param SingerListParam) (*SingerListData, error) {
	params := map[string]interface{}{
		"comm": map[string]interface{}{"ct": 24, "cv": 0},
		"singerList": map[string]interface{}{
			"module": "Music.SingerListServer",
			"method": "get_singer_list",
			"param":  param,
		},
	}
	paramBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	//fmt.Println(singerListApi + string(paramBytes))
	resp, err := request(singerListApi+string(paramBytes), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var resData struct {
		Code int                `json:"code"`
		List SingerListResponse `json:"singerList"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&resData); err != nil {
		return nil, errors.New("singer list encode err")
	}
	if resData.Code != 0 {
		return nil, errors.New("singer list encode err")
	}
	return &resData.List.Data, nil
}

func FetchAllSinger(char string) (SingerList, error) {
	index := utils.CharToIndex(char)
	param := SingerListParam{
		Area:    -100,
		Genre:   -100,
		Sex:     -100,
		Index:   index,
		Sin:     0,
		CurPage: 1,
	}
	data, err := GetSingerListData(param)
	if err != nil {
		return nil, err
	}
	//80条每页
	list := data.SingerList
	page := 2
	maxPage := int(math.Ceil(float64(data.Total) / 80))
	for page < maxPage {
		param.CurPage = page
		param.Sin = (page - 1) * 80
		data, err := GetSingerListData(param)
		if err != nil {
			return nil, err
		}
		list = append(list, data.SingerList...)
		page++
	}

	return list, nil
}

func FetchAllSingerAndStore(char string, path string) error {
	list, err := FetchAllSinger(char)
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(list, "", "	")
	if err != nil {
		easylog.Error(err)
	}
	return ioutil.WriteFile(path, bytes, 0666)
}

func GetSingerMidByName(path string, name string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	var list SingerList
	if err := json.Unmarshal(data, &list); err != nil {
		return ""
	}
	for _, singer := range list {
		if strings.Contains(singer.SingerName, name) {
			return singer.SingerMid
		}
	}
	return ""
}

type AlbumItem struct {
	AlbumId    string `json:"albumID"`
	AlbumMid   string `json:"albumMID"`
	AlbumName  string `json:"albumName"`
	Albumtype  string `json:"albumtype"`
	Company    string `json:"company"`
	Desc       string `json:"desc"`
	Lan        string `json:"lan"`
	SingerID   string `json:"singerId"`
	SingerMID  string `json:"singerMid"`
	SingerName string `json:"singerName"`
}

type AlbumList []AlbumItem

func GetAlbumList(singerMid string) (AlbumList, error) {
	resp, err := request(albumListApi+singerMid, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var resData struct {
		Code int `json:"code"`
		Data struct {
			List  AlbumList `json:"list"`
			Total int       `json:"total"`
		} `json:"data"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&resData); err != nil {
		return nil, err
	}
	if resData.Code != 0 {
		return nil, err
	}
	//fmt.Printf("res: %v \n", resData.Data.Total)
	return resData.Data.List, nil
}

func GenerateAlbumUrl(singerMid string) []string {
	albumList, err := GetAlbumList(singerMid)
	if err != nil {
		easylog.Error(err)
		return nil
	}
	var list []string
	format := "https://y.qq.com/n/yqq/album/%s.html"
	for _, album := range albumList {
		url := fmt.Sprintf(format, album.AlbumMid)
		list = append(list, url)
	}
	return list
}
