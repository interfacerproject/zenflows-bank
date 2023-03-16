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
)

type Fabcoin struct {
	EthereumAddress string
	Idea            int64
	Strengths       int64
	Fabcoin       int64
	TxId       string
}

type FabcoinFile struct {
	Fabcoins   []*Fabcoin
	FileName string
}

func (m FabcoinFile) IterateRows() <-chan []interface{} {
	c := make(chan []interface{})
	go func() {
		for _, v := range m.Fabcoins {
			c <- []interface{}{
				v.EthereumAddress,
				v.Idea,
				v.Strengths,
				v.Fabcoin,
				v.TxId,
			}
		}
		close(c)
	}()
	return c
}

func (ts FabcoinFile) Header() []interface{} {
	return []interface{}{
		"EthereumAddress",
		"Idea",
		"Strengths",
		"Fabcoin",
		"TxId",
	}
}

func (ts FabcoinFile) ExportCSV() {
	ExportCSV(ts.FileName, ts)
}

func (ts FabcoinFile) ExportXLS() {
	ExportXLS(ts.FileName, ts, "Balances")
}

func (ts FabcoinFile) Export() {
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

