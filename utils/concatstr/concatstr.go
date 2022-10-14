package concatstr

import "bytes"

// ConcatString 拼接字符串
func ConcatString(str ...string) string {
	var buf bytes.Buffer
	for _, s := range str {
		buf.WriteString(s)
	}
	return buf.String()
}
