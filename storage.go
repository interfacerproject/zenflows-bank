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

package main

import (
	//"encoding/json"
	"errors"
	"github.com/tarantool/go-tarantool"
	"log"
	"strconv"
	"time"
)

type TTStorage struct {
	db *tarantool.Connection
}

const MAX_RETRY int = 10

func (storage *TTStorage) Init(host, user, pass string) error {
	var err error
	for done, retry := false, 0; !done; retry++ {
		storage.db, err = tarantool.Connect(host, tarantool.Opts{
			User: user,
			Pass: pass,
		})
		done = retry == MAX_RETRY || err == nil
		if !done {
			log.Println("Could not connect to tarantool, retrying...")
			time.Sleep(3 * time.Second)
		} else {
			log.Println("Connected to tarantool")
		}
	}
	return err
}

func (storage *TTStorage) Balances() (map[string]*Token, error) {
	const limit uint32 = 100
	var offset uint32 = 0

	balances := make(map[string]*Token)

	for {
		resp, err := storage.db.Select("TXS", "primary", offset, limit, tarantool.IterEq, []interface{}{})
		if err != nil {
			return nil, err
		}
		if resp.Error != "" {
			return nil, errors.New(resp.Error)
		}
		if len(resp.Data) == 0 {
			break
		}
		for i := 0; i < len(resp.Data); i = i + 1 {
			id := resp.Data[i].([]interface{})[1].(string)

			n, err := strconv.ParseInt(resp.Data[i].([]interface{})[4].(string), 10, 64)
			if err != nil {
				return nil, err
			}
			if balances[id] == nil {
				balances[id] = &Token{}
			}
			switch resp.Data[i].([]interface{})[2].(string) {
			case "idea":
				balances[id].Idea = balances[id].Idea + n
			case "strength":
				balances[id].Strength = balances[id].Strength + n
			}
		}
		offset = offset + limit
	}
	return balances, nil
}
