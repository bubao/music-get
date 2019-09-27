package utils

import (
	"fmt"
	"testing"
)

func TestIsHan(t *testing.T) {
	str := "中国 123"
	runes := []rune(str)

	fmt.Println(runes[0], string(runes[0]), IsHan(runes[0]))
	if IsHan(runes[0]) {
		arr := ToPinyin(string(runes[0]))
		fmt.Println(arr, string(arr[0][0][0]))
	}
}

func TestCharToIndex(t *testing.T) {
	fmt.Println(CharToIndex("b"))
}
