// Copyright 2016 - 2026 The excelize Authors. All rights reserved.
// Dedup demo fixture: Type 4 semantic duplicate.

package excelize

import "math"

// newDemoNumericFormulaArg intentionally performs the same work as newNumberFormulaArg
// with different naming so the semantic duplicate is easy to inspect in the demo PR.
func newDemoNumericFormulaArg(value float64) formulaArg {
	if math.IsNaN(value) {
		return newErrorFormulaArg(formulaErrorNUM, formulaErrorNUM)
	}
	return formulaArg{Type: ArgNumber, Number: value}
}
