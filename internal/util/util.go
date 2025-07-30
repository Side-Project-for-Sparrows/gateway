package util

import "strings"

// 서비스 이름 추출
func ExtractServiceKey(path string) string {
	return strings.Split(strings.Trim(path, "/"), "/")[0] // /board/list → "board"
}
