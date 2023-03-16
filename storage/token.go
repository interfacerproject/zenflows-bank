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
	"path/filepath"
	"encoding/csv"
	"strconv"
	"os"
	"github.com/xuri/excelize/v2"
)

type Token struct {
	Id string
	EthereumAddress string
	Idea            int64
	Strengths       int64
}

type Tokens []*Token

func (ts Tokens) ExportCSV(fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	wr := csv.NewWriter(f)
	wr.Write([]string{
		"ID",
		"EthereumAddress",
		"Idea",
		"Strengths",
	})
	for _, v := range ts {
		wr.Write([]string{
			v.Id,
			v.EthereumAddress,
			strconv.FormatInt(v.Idea, 10),
			strconv.FormatInt(v.Strengths, 10),
		})
	}
	wr.Flush()
}

func (ts Tokens) ExportXLS(fileName string) {
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
	f.SetSheetRow("Balances", "A1", &[]interface{}{
		"ID",
		"EthereumAddress",
		"Idea",
		"Strengths",
	})
    for idx, row := range ts {
        cell, err := excelize.CoordinatesToCellName(1, idx+2)
        if err != nil {
			panic(err)
        }
        f.SetSheetRow("Balances", cell, &[]interface{}{
			row.Id,
			row.EthereumAddress,
			row.Idea,
			row.Strengths,
		})
    }
	f.SetActiveSheet(index)
	if err := f.SaveAs(fileName); err != nil {
        panic(err)
    }
}

func (ts Tokens) Export(fileName string) {
	ext := filepath.Ext(fileName)
	switch ext {
	case ".csv":
		ts.ExportCSV(fileName)
	case ".xlsx":
		ts.ExportXLS(fileName)
	default:
		panic("Unsupported extension " + ext)
	}
}
