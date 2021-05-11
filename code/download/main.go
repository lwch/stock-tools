package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lwch/runtime"
)

// https://blog.csdn.net/leijia_xing/article/details/81139005

const urlAddr = "http://q.stock.sohu.com/hisHq"

func main() {
	code := flag.String("code", "cn_000001", "股票代号cn_编号")
	begin := flag.String("begin", "-100", "开始时间，负数表示向前N天，或yyyymmdd")
	flag.Parse()

	var beginTime time.Time
	endTime := time.Now()

	if strings.HasPrefix(*begin, "-") {
		n, err := strconv.ParseInt(*begin, 10, 64)
		runtime.Assert(err)
		beginTime = endTime.Add(time.Duration(n) * 24 * time.Hour)
	} else {
		var err error
		beginTime, err = time.ParseInLocation("20060102", *begin, time.Local)
		runtime.Assert(err)
	}

	args := make(url.Values)
	args.Set("code", *code)
	args.Set("start", beginTime.Format("20060102"))
	args.Set("end", endTime.Format("20060102"))
	args.Set("stat", "1")
	args.Set("order", "A")
	args.Set("period", "d")

	resp, err := http.Get(urlAddr + "?" + args.Encode())
	runtime.Assert(err)
	defer resp.Body.Close()

	type dt struct {
		HQ [][10]string `json:"hq"`
	}
	var ret []dt
	runtime.Assert(json.NewDecoder(resp.Body).Decode(&ret))
	if len(ret) == 0 {
		panic("no data")
	}
	data := ret[0]
	f, err := os.Create(fmt.Sprintf("%s_%s_%s.csv", *code,
		beginTime.Format("20060102"), endTime.Format("20060102")))
	runtime.Assert(err)
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	runtime.Assert(w.Write([]string{"date", "open", "close", "volume", "turn"}))
	for _, date := range data.HQ {
		t, err := time.ParseInLocation("2006-01-02", date[0], time.Local)
		runtime.Assert(err)
		open := date[1]   // 开盘
		close := date[2]  // 收盘
		volume := date[7] // 成交量
		turn := date[8]   // 成交额
		runtime.Assert(w.Write([]string{
			t.Format("20060102"),
			open,
			close,
			volume,
			turn,
		}))
	}
}
