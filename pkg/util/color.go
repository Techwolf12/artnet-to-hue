package util

import "math"

func ConvertDMXToHue(data []byte) (float64, float64, int) {
	if len(data) < 3 {
		return 0, 0, 0 // Not enough data
	}

	// Extract RGB values from DMX data
	r := float64(data[0]) / 255.0
	g := float64(data[1]) / 255.0
	b := float64(data[2]) / 255.0

	// Convert RGB to XYZ
	x, y, z := rgbToXyz(r, g, b)

	// Convert XYZ to xy
	cx := x / (x + y + z)
	cy := y / (x + y + z)

	// Calculate brightness
	brightness := int(y * 254)

	return cx, cy, brightness
}

// Convert RGB to XYZ
func rgbToXyz(r, g, b float64) (float64, float64, float64) {
	r = pivotRgb(r)
	g = pivotRgb(g)
	b = pivotRgb(b)

	x := r*0.4124564 + g*0.3575761 + b*0.1804375
	y := r*0.2126729 + g*0.7151522 + b*0.0721750
	z := r*0.0193339 + g*0.1191920 + b*0.9503041

	return x, y, z
}

// Pivot RGB values
func pivotRgb(value float64) float64 {
	if value > 0.04045 {
		return math.Pow((value+0.055)/1.055, 2.4)
	}
	return value / 12.92
}
