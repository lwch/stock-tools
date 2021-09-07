package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"tools/code/utils"
	"unsafe"

	"github.com/Fontinalis/fonet"
	"github.com/lwch/runtime"
)

const review = 5
const srcCount = 8
const columnCount = srcCount + 2*review
const batch = 1000

func main() {
	lr := flag.Float64("lr", 0.01, "学习率")
	round := flag.Int("round", 500, "迭代次数")
	out := flag.String("out", "", "输出目录")
	flag.Parse()

	if len(*out) == 0 {
		fmt.Println("缺少out参数")
		os.Exit(1)
	}

	if flag.NArg() < 1 {
		fmt.Println("usage: ./dnn <input-dir>")
		os.Exit(1)
	}

	from := flag.Arg(0)
	var err error
	from, err = filepath.Abs(from)
	runtime.Assert(err)
	files, err := filepath.Glob(path.Join(from, "*.csv"))
	runtime.Assert(err)

	for _, file := range files {
		fmt.Printf("开始训练：%s\n", file)
		data, err := utils.LoadCSV(file)
		if err != nil {
			continue
		}
		// 特征矩阵：
		// 开盘价，收盘价，最低价，最高价，成交量，成交额，涨幅，涨幅百分比，
		// 前1天收盘价，前1天涨幅，...前5天收盘价，前5天涨幅
		samples := normalize(data)
		if len(samples) < 100 {
			continue
		}
		labels := make([]float64, len(samples)-review-1)
		for i := review; i < len(samples)-1; i++ {
			// 填充前N天的数据
			for j := 0; j < review; j++ {
				samples[i][srcCount+j*2] = samples[i-j-1][1]
				samples[i][srcCount+j*2+1] = samples[i-j-1][6]
			}
			// delta大于0为正样本，否则为负样本
			if data.Rows[i+1].Delta > 0 {
				labels[i-review] = 1
			}
		}
		samples = samples[review:]
		samples = samples[:len(samples)-1]
		rand.Shuffle(len(labels), func(i, j int) {
			samples[i], samples[j] = samples[j], samples[i]
			labels[i], labels[j] = labels[j], labels[i]
		})
		// 训练集和测试集比例7:3
		n := len(samples) * 7 / 10
		// trainSamples := samples[:n]
		// trainLabels := labels[:n]
		testSamples := samples[n+1:]
		testLabels := labels[n+1:]

		net, err := fonet.NewNetwork([]int{columnCount, 8, 4, 2, 1}, fonet.LeakyReLU)
		runtime.Assert(err)
		trainData := make([][][]float64, len(samples))
		for i := 0; i < len(samples); i++ {
			trainData[i] = [][]float64{
				samples[i][:],
				{labels[i]},
			}
		}
		for i := 0; i < *round; i += batch {
			net.Train(trainData, batch, *lr, false)
			acc := test(net, samples, labels)
			fmt.Printf("第%d轮训练，准确率：%.2f%%\n", i+batch, acc*100)
		}
		left := *round % batch
		if left == 0 {
			continue
		}
		net.Train(trainData, batch, *lr, false)
		acc := test(net, testSamples, testLabels)
		fmt.Printf("第%d轮训练，准确率：%.2f%%\n", *round, acc*100)
	}
}

func normalize(data *utils.Data) [][columnCount]float64 {
	var min [columnCount]float64
	var max [columnCount]float64
	var offsets []uintptr
	get := func(row utils.Row, offset uintptr) float64 {
		ptr := unsafe.Pointer(&row)
		return *(*float64)(unsafe.Pointer(uintptr(ptr) + offset))
	}
	div := func(a, b float64) float64 {
		if b == 0 {
			return 0
		}
		return a / b
	}
	row := data.Rows[0]
	offsets = append(offsets, unsafe.Offsetof(row.Open))
	offsets = append(offsets, unsafe.Offsetof(row.Close))
	offsets = append(offsets, unsafe.Offsetof(row.Low))
	offsets = append(offsets, unsafe.Offsetof(row.High))
	offsets = append(offsets, unsafe.Offsetof(row.Volumn))
	offsets = append(offsets, unsafe.Offsetof(row.Turn))
	offsets = append(offsets, unsafe.Offsetof(row.Delta))
	offsets = append(offsets, unsafe.Offsetof(row.DeltaP))
	for i := 0; i < len(offsets); i++ {
		n := get(row, offsets[i])
		min[i] = n
		max[i] = n
	}
	for _, row := range data.Rows {
		for i := 0; i < len(offsets); i++ {
			n := get(row, offsets[i])
			if n < min[i] {
				min[i] = n
			}
			if n > max[i] {
				max[i] = n
			}
		}
	}
	list := make([][columnCount]float64, len(data.Rows))
	for idx, row := range data.Rows {
		var raw [columnCount]float64
		for i := 0; i < len(offsets); i++ {
			n := get(row, offsets[i])
			raw[i] = div(n-min[i], max[i]-min[i])
		}
		list[idx] = raw
	}
	return list
}

func test(net *fonet.Network, samples [][columnCount]float64, labels []float64) float64 {
	var total float64
	for i := 0; i < len(samples); i++ {
		predict := net.Predict(samples[i][:])
		diff := labels[i] - predict[0]
		if diff < 0 {
			diff = -diff
		}
		total += diff
	}
	return 1 - total/float64(len(samples))
}
