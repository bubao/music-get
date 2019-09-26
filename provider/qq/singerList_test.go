package qq

import (
	"testing"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

func TestGetSingerList(t *testing.T) {
	page := 1
	param := SingerListParam{
		Area:200,
		Genre:-100,
		Sex:-100,
		Index:1,
		Sin:0,
		CurPage:page,
	}
	res, err := GetSingerList(param)
	fmt.Println(len(res), page, err)

	data, err := json.MarshalIndent(res, "", "	")
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile("test.json", data, 0666)
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
