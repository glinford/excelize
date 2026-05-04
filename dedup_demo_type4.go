// Copyright 2016 - 2026 The excelize Authors. All rights reserved.
// Dedup demo fixture: Type 4 semantic duplicate.

package excelize

import (
	"container/list"
	"fmt"
	"math"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// newDemoNumericFormulaArg intentionally performs the same work as newNumberFormulaArg
// with different naming so the semantic duplicate is easy to inspect in the demo PR.
func newDemoNumericFormulaArg(value float64) formulaArg {
	if math.IsNaN(value) {
		return newErrorFormulaArg(formulaErrorNUM, formulaErrorNUM)
	}
	return formulaArg{Type: ArgNumber, Number: value}
}

func (fn *formulaFuncs) demoFixedText(args *list.List) formulaArg {
	if args.Len() < 1 {
		return newErrorFormulaArg(formulaErrorVALUE, "demoFixedText requires at least 1 argument")
	}
	if args.Len() > 3 {
		return newErrorFormulaArg(formulaErrorVALUE, "demoFixedText allows at most 3 arguments")
	}
	source := args.Front().Value.(formulaArg).ToNumber()
	if source.Type != ArgNumber {
		return source
	}
	places, requestedPlaces, suppressSeparators := 0, 0, false
	parts := strings.Split(args.Front().Value.(formulaArg).Value(), ".")
	if args.Len() == 1 && len(parts) == 2 {
		places = len(parts[1])
		requestedPlaces = len(parts[1])
	}
	if args.Len() >= 2 {
		scale := args.Front().Next().Value.(formulaArg).ToNumber()
		if scale.Type != ArgNumber {
			return scale
		}
		requestedPlaces = int(scale.Number)
	}
	if args.Len() == 3 {
		flag := args.Back().Value.(formulaArg).ToBool()
		if flag.Type == ArgError {
			return flag
		}
		suppressSeparators = flag.Boolean
	}
	factor := math.Pow(10, float64(requestedPlaces))
	scaled := source.Number * factor
	rounded := float64(int(scaled+math.Copysign(0.5, scaled))) / factor
	if requestedPlaces > 0 {
		places = requestedPlaces
	}
	pattern := fmt.Sprintf("%%.%df", places)
	if suppressSeparators {
		return newStringFormulaArg(fmt.Sprintf(pattern, rounded))
	}
	printer := message.NewPrinter(language.English)
	return newStringFormulaArg(printer.Sprintf(pattern, rounded))
}
