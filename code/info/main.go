package main

import (
	"fmt"
	"os"
	"tools/code/utils"

	"github.com/lwch/runtime"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ./info xx.csv")
		os.Exit(1)
	}

	data, err := utils.LoadCSV(os.Args[1])
	runtime.Assert(err)

	var opens, closes []float64
	var upDays, downDays int
	for _, row := range data.Rows {
		opens = append(opens, row.Open)
		closes = append(closes, row.Close)
		if row.Close-row.Open > 0 {
			upDays++
		} else if row.Close-row.Open < 0 {
			downDays++
		}
	}
	show(opens, "开盘价")
	show(closes, "收盘价")
	fmt.Printf("总天数: %d, 上涨天数: %d(%.02f%%), 下跌天数: %d(%.02f%%)\n", len(opens),
		upDays, float64(upDays)*100./float64(len(opens)), downDays, float64(downDays)*100./float64(len(opens)))
}

func show(arr []float64, prefix string) {
	last := arr[len(arr)-1]
	fmt.Printf(`%s:
  均值=%.02f, 中位数=%.02f, 最后=%.02f
  最小值=%.02f, 最大值=%.02f, 标准差=%.02f
  P10=%.02f, P70=%.02f, P90=%.02f
`,
		prefix, utils.Avg(arr), utils.Percent(50)(arr), last,
		utils.Min(arr), utils.Max(arr), utils.Stdev(arr),
		utils.Percent(10)(arr), utils.Percent(70)(arr), utils.Percent(90)(arr))
}
