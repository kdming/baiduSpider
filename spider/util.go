package spider

import (
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 工具类

// 字符串编码
func EncodeStr(str string) string {
	encodeStr, _, _ := transform.String(simplifiedchinese.GBK.NewEncoder(), str)
	return encodeStr
}
