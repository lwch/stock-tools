package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/lwch/runtime"
)

// BeginToTime 解析begin参数并转换为time.Time
func BeginToTime(begin string) time.Time {
	now := time.Now()
	if strings.HasPrefix(begin, "-") {
		n, err := strconv.ParseInt(begin, 10, 64)
		runtime.Assert(err)
		return now.Add(time.Duration(n) * 24 * time.Hour)
	}
	t, err := time.ParseInLocation("20060102", begin, time.Local)
	runtime.Assert(err)
	return t
}
