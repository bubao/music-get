package qq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestGetSingerList(t *testing.T) {
	page := 1
	param := SingerListParam{
		Area:    -100,
		Genre:   -100,
		Sex:     -100,
		Index:   1,
		Sin:     0,
		CurPage: page,
	}
	res, err := GetSingerListData(param)
	fmt.Println(len(res.SingerList), page, err)

	data, err := json.MarshalIndent(res.SingerList, "", "	")
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile("test.json", data, 0666)
}

func TestFetchAllSinger(t *testing.T) {
	char := "m"
	data, err := FetchAllSinger(char)
	bytes, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile(strings.ToLower(char)+"-singer.json", bytes, 0666)
}

func TestGetSingerMidByName(t *testing.T) {
	name := "Michael Jackson"
	fmt.Println(GetSingerMidByName("m-singer.json", name))
}

func TestGetAlbumList(t *testing.T) {
	res, err := GetAlbumList("0025NhlN2yWrP4")
	fmt.Println(len(res), err)

	data, err := json.MarshalIndent(res, "", "	")
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile("test-album.json", data, 0666)
}
