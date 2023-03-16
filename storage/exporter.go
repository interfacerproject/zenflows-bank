package storage

import (
	"encoding/csv"
	"github.com/xuri/excelize/v2"
	"os"
	"strconv"
)

type Exportable interface {
	Header() []interface{}
	IterateRows() <-chan []interface{}
}

func ExportXLS(fileName string, data Exportable, sheetName string) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	index, err := f.NewSheet("Balances")
	if err != nil {
		panic(err)
	}
	header := data.Header()
	f.SetSheetRow("Balances", "A1", &header)
	idx := 0
	for row := range data.IterateRows() {
		cell, err := excelize.CoordinatesToCellName(1, idx+2)
		if err != nil {
			panic(err)
		}
		f.SetSheetRow(sheetName, cell, &row)
		idx = idx + 1
	}
	f.SetActiveSheet(index)
	if err := f.SaveAs(fileName); err != nil {
		panic(err)
	}
}

func ExportCSV(fileName string, data Exportable) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// TODO: write []interface{} -> []string function
	wr := csv.NewWriter(f)
	header := data.Header()
	line := make([]string, len(header))
	for k, w := range header {
		if i, ok := w.(int64); ok {
			line[k] = strconv.FormatInt(i, 10)
		} else {
			line[k] = w.(string)
		}
	}
	for v := range data.IterateRows() {
		line := make([]string, len(v))
		for k, w := range v {
			if i, ok := w.(int64); ok {
				line[k] = strconv.FormatInt(i, 10)
			} else {
				line[k] = w.(string)
			}
		}
		wr.Write(line)
	}
	wr.Flush()
}
