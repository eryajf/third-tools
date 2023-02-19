package main

import (
	"fmt"

	"github.com/mozillazg/go-pinyin"
)

func main() {
	fmt.Println(ConvertName("你好"))
}

func ConvertName(src string) (dst string) {
	args := pinyin.NewArgs()
	args.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}

	for _, singleResult := range pinyin.Pinyin(src, args) {
		for _, result := range singleResult {
			dst = dst + result
		}
	}
	return
}
