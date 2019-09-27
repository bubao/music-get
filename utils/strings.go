package utils

import (
	"github.com/mozillazg/go-pinyin"
	"strings"
	"unicode"
)

//字母转26位字母表中的位置; 大写
func CharToIndex(char string) int {
	if len(char) > 1 {
		return -100
	}
	char = strings.ToUpper(char)
	runes := []rune(char)
	return int(runes[0]) - 64
}

func GetIndexChar(s string) string {
	indexChar := ""
	runes := []rune(s)
	if IsHan(runes[0]) {
		arr := ToPinyin(string(runes[0]))
		indexChar = string(arr[0][0][0])
	} else {
		indexChar = string(runes[0])
	}
	return indexChar
}

func IsHan(r rune) bool {
	return unicode.Is(unicode.Han, r)
}

func ToPinyin(s string) [][]string {
	a := pinyin.NewArgs()
	return pinyin.Pinyin(s, a)
}
