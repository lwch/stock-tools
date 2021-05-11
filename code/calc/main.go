package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"

	"github.com/lwch/runtime"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ./calc xx.csv")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	runtime.Assert(err)
	defer f.Close()
	r := csv.NewReader(f)
	// skip header
	_, err = r.Read()
	runtime.Assert(err)
	var opens, closes []float64
	for {
		data, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		open, err := strconv.ParseFloat(data[1], 64)
		runtime.Assert(err)
		close, err := strconv.ParseFloat(data[2], 64)
		runtime.Assert(err)
		opens = append(opens, open)
		closes = append(closes, close)
	}
	sort.Float64s(opens)
	sort.Float64s(closes)
	show(opens, "open")
	show(closes, "close")
}

func show(arr []float64, prefix string) {
	fmt.Printf(`%s:
  avg=%.02f, mean=%.02f
  min=%.02f, max=%.02f
  P10=%.02f, P70=%.02f, P90=%.02f
`,
		prefix, avg(arr), arr[len(arr)/2],
		arr[0], arr[len(arr)-1],
		percent(arr, 10), percent(arr, 70), percent(arr, 90))
}

func sum(arr []float64) float64 {
	var ret float64
	for _, n := range arr {
		ret += n
	}
	return ret
}

func avg(arr []float64) float64 {
	return sum(arr) / float64(len(arr))
}

func percent(arr []float64, n int) float64 {
	return arr[len(arr)*n/100]
}
