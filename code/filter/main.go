package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"tools/code/utils"

	"github.com/lwch/runtime"
)

func main() {
	begin := flag.String("begin", "-30", "开始时间，负数表示向前N天，或yyyymmdd")
	column := flag.String("column", "close", "筛选字段，支持open,close,low,high")
	f := flag.String("func", "avg", "聚合函数，支持max,min,sum,avg,stdev,p<n>")
	gt := flag.Float64("gt", 0, "大于该值")
	lt := flag.Float64("lt", 0, "小于该值")
	out := flag.String("out", "", "输出目录")
	flag.Parse()

	if *gt < math.SmallestNonzeroFloat32 && *lt < math.SmallestNonzeroFloat32 {
		fmt.Println("缺少gt或lt参数")
		os.Exit(1)
	}
	var pfunc utils.PFunc
	switch {
	case *f == "max":
		pfunc = utils.Max
	case *f == "min":
		pfunc = utils.Min
	case *f == "sum":
		pfunc = utils.Sum
	case *f == "avg":
		pfunc = utils.Avg
	case *f == "stdev":
		pfunc = utils.Stdev
	case len(*f) > 0 && (*f)[0] == 'p':
		p, err := strconv.ParseInt((*f)[1:], 10, 32)
		runtime.Assert(err)
		pfunc = utils.PercentFloat32(int(p))
	}
	if len(*out) == 0 {
		fmt.Println("缺少out参数")
		os.Exit(1)
	}

	beginTime := utils.BeginToTime(*begin)

	if flag.NArg() < 1 {
		fmt.Println("usage: ./filter <input-dir>")
		os.Exit(1)
	}

	files, err := filepath.Glob(path.Join(flag.Arg(0), "*.csv"))
	runtime.Assert(err)

	runtime.Assert(os.MkdirAll(*out, 0755))
	for _, file := range files {
		data, err := utils.LoadCSV(file)
		if err != nil {
			continue
		}
		var rows []utils.Row
		for _, row := range data.Rows {
			if row.Date.After(beginTime) {
				rows = append(rows, row)
			}
		}
		if len(rows) == 0 {
			continue
		}
		values := make([]float64, len(rows))
		for i, row := range rows {
			switch *column {
			case "open":
				values[i] = row.Open
			case "close":
				values[i] = row.Close
			case "high":
				values[i] = row.High
			case "low":
				values[i] = row.Low
			}
		}
		n := pfunc(values)
		want := false
		if *gt > math.SmallestNonzeroFloat32 {
			if n > *gt {
				want = true
			}
		} else {
			if n < *lt {
				want = true
			}
		}
		if !want {
			continue
		}
		fi, err := os.Stat(file)
		runtime.Assert(err)
		for fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			file, err = os.Readlink(file)
			runtime.Assert(err)
			fi, err = os.Stat(file)
			runtime.Assert(err)
		}
		file, err = filepath.Abs(file)
		runtime.Assert(err)
		os.Symlink(file, path.Join(*out, path.Base(file)))
		fmt.Printf("%s added\n", path.Base(file))
	}
}
