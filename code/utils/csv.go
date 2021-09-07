package utils

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/lwch/runtime"
)

type Row struct {
	Date   time.Time
	Open   float64
	Close  float64
	Low    float64
	High   float64
	Volumn float64
	Turn   float64
	Delta  float64
	DeltaP float64
}

type Data struct {
	Rows []Row
}

func LoadCSV(dir string) (*Data, error) {
	f, err := os.Open(dir)
	runtime.Assert(err)
	defer f.Close()
	reader := csv.NewReader(f)
	hdr, err := reader.Read()
	if err != nil {
		return nil, err
	}
	date := -1
	open := -1
	close := -1
	low := -1
	high := -1
	volumn := -1
	turn := -1
	delta := -1
	deltaP := -1
	for i, name := range hdr {
		switch name {
		case "date":
			date = i
		case "open":
			open = i
		case "close":
			close = i
		case "low":
			low = i
		case "high":
			high = i
		case "volumn":
			volumn = i
		case "turn":
			turn = i
		case "delta":
			delta = i
		case "deltap":
			deltaP = i
		}
	}
	readFloat := func(row []string, idx int, dst *float64) error {
		if idx == -1 {
			return nil
		}
		if len(row) <= idx {
			return nil
		}
		*dst, err = strconv.ParseFloat(row[idx], 32)
		return err
	}
	ret := new(Data)
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return ret, nil
			}
			return nil, err
		}
		var next Row
		if date != -1 && len(row) > date {
			next.Date, err = time.ParseInLocation("20060102", row[date], time.Local)
			runtime.Assert(err)
		}
		err = readFloat(row, open, &next.Open)
		if err != nil {
			return nil, err
		}
		err = readFloat(row, close, &next.Close)
		if err != nil {
			return nil, err
		}
		err = readFloat(row, low, &next.Low)
		if err != nil {
			return nil, err
		}
		err = readFloat(row, high, &next.High)
		if err != nil {
			return nil, err
		}
		err = readFloat(row, volumn, &next.Volumn)
		if err != nil {
			return nil, err
		}
		err = readFloat(row, turn, &next.Turn)
		if err != nil {
			return nil, err
		}
		err = readFloat(row, delta, &next.Delta)
		if err != nil {
			return nil, err
		}
		err = readFloat(row, deltaP, &next.DeltaP)
		if err != nil {
			return nil, err
		}
		ret.Rows = append(ret.Rows, next)
	}
}
