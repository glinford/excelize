// Copyright 2016 - 2026 The excelize Authors. All rights reserved. Use of
// this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package excelize

import (
	"strconv"
	"strings"
)

// parseRangeAxis splits a worksheet axis into a column token and row number.
func parseRangeAxis(axis string) (string, int, error) {
	isColumnRune := func(r rune) bool {
		return ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z') || (r == 36)
	}
	if strings.IndexFunc(axis, isColumnRune) == 0 {
		idx := strings.LastIndexFunc(axis, isColumnRune)
		if idx >= 0 && idx < len(axis)-1 {
			column, rowText := strings.ReplaceAll(axis[:idx+1], "$", ""), axis[idx+1:]
			if row, err := strconv.Atoi(rowText); err == nil && row > 0 {
				return column, row, nil
			}
		}
	}
	return "", -1, newInvalidCellNameError(axis)
}
