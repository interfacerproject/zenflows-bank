// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2023 Dyne.org foundation <foundation@dyne.org>.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package storage

import (
	"encoding/csv"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Token struct {
	Id              string
	EthereumAddress string
	Idea            int64
	Strengths       int64
}

type TokensFile struct {
	Tokens   []*Token
	FileName string
}

func (m TokensFile) IterateRows() <-chan []interface{} {
	c := make(chan []interface{})
	go func() {
		for _, v := range m.Tokens {
			c <- []interface{}{
				v.Id,
				v.EthereumAddress,
				v.Idea,
				v.Strengths,
			}
		}
		close(c)
	}()
	return c
}

func (ts TokensFile) Header() []interface{} {
	return []interface{}{
		"ID",
		"EthereumAddress",
		"Idea",
		"Strengths",
	}
}

func (ts TokensFile) ExportCSV() {
	ExportCSV(ts.FileName, ts)
}

func (ts TokensFile) ExportXLS() {
	ExportXLS(ts.FileName, ts, "Balances")
}

func (ts TokensFile) Export() {
	ext := filepath.Ext(ts.FileName)
	switch ext {
	case ".csv":
		ts.ExportCSV()
	case ".xlsx":
		ts.ExportXLS()
	default:
		panic("Unsupported extension " + ext)
	}
}

func (ts *TokensFile) ImportCSV() {
	f, err := os.Open(ts.FileName)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)

	// Remove header
	_, err = csvReader.Read()
	if err == io.EOF {
		return
	}
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		idea, err := strconv.ParseInt(rec[2], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		strengths, err := strconv.ParseInt(rec[3], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		ts.Tokens = append(ts.Tokens,
			&Token{
				Id:              rec[0],
				EthereumAddress: rec[1],
				Idea:            idea,
				Strengths:       strengths,
			},
		)
	}
}

func (ts *TokensFile) ImportXLSX() {
	f, err := excelize.OpenFile(ts.FileName)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	readCell := func(i, j int) string {
		cell, err := excelize.CoordinatesToCellName(i, j)
		if err != nil {
			panic(err)
		}
		val, err := f.GetCellValue("Balances", cell)
		if err != nil {
			panic(err)
		}
		return val
	}

	for i := 2; ; i = i + 1 {
		id := readCell(1, i)
		if id == "" {
			break
		}
		ethAddr := readCell(2, i)
		idea, err := strconv.ParseInt(readCell(3, i), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		strengths, err := strconv.ParseInt(readCell(4, i), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		ts.Tokens = append(ts.Tokens,
			&Token{
				Id:              id,
				EthereumAddress: ethAddr,
				Idea:            idea,
				Strengths:       strengths,
			},
		)
	}
}

func (ts *TokensFile) Import() {
	ext := filepath.Ext(ts.FileName)
	switch ext {
	case ".csv":
		ts.ImportCSV()
	case ".xlsx":
		ts.ImportXLSX()
	default:
		panic("Unsupported extension " + ext)
	}
}
