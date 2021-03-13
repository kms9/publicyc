package util

import (
	"strconv"
	"strings"
)
// 该函数比较两个版本号是否相等，是否大于或小于的关系
// 返回值：0表示v1与v2相等；1表示v1大于v2；-1表示v1小于v2
func VersionCompare(v1, v2 string) (ret int) {
	// 替换一些常见的版本符号
	replaceMap := map[string]string{
		"V": "",
		"v": "",
		"-": ".",
	}
	for k, v := range replaceMap {
		if strings.Contains(v1, k) {
			v1 = strings.Replace(v1, k, v, -1)
		}
		if strings.Contains(v2, k) {
			v2 = strings.Replace(v2, k, v, -1)
		}
	}
	v1s := strings.Split(v1, ".")
	v2s := strings.Split(v2, ".")
	loopMax := len(v2s)
	if len(v1s) > len(v2s) {
		loopMax = len(v1s)
	}
	for i := 0; i < loopMax; i++ {
		var x, y string
		if len(v1s) > i {
			x = v1s[i]
		}
		if len(v2s) > i {
			y = v2s[i]
		}
		xi,_ := strconv.Atoi(x)
		yi,_ := strconv.Atoi(y)
		if xi > yi {
			ret = 1
		} else if xi < yi {
			ret = -1
		}
		if ret != 0 {
			break
		}
	}
	return
}
