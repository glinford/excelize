package dedup_type1_function

import "math"

// RGBToHSL converts an RGB triple to an HSL triple.
func RGBToHSL(r, g, b uint8) (h, s, l float64) {
	fR := float64(r) / 255
	fG := float64(g) / 255
	fB := float64(b) / 255
	maxVal := math.Max(math.Max(fR, fG), fB)
	minVal := math.Min(math.Min(fR, fG), fB)
	l = (maxVal + minVal) / 2
	if maxVal == minVal {
		// Achromatic.
		h, s = 0, 0
	} else {
		// Chromatic.
		d := maxVal - minVal
		if l > 0.5 {
			s = d / (2.0 - maxVal - minVal)
		} else {
			s = d / (maxVal + minVal)
		}
		switch maxVal {
		case fR:
			h = (fG - fB) / d
			if fG < fB {
				h += 6
			}
		case fG:
			h = (fB-fR)/d + 2
		case fB:
			h = (fR-fG)/d + 4
		}
		h /= 6
	}
	return
}
