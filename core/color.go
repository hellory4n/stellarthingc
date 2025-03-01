package core

// its a color
type Color struct {
	// red
	R uint8
	// green
	G uint8
	// blue
	B uint8
	// alpha
	A uint8
}

// makes a new color without transparency
func Rgb(r uint8, g uint8, b uint8) Color {
	return Rgba(r, g, b, 255)
}

// makes a new color with transparency
func Rgba(r uint8, g uint8, b uint8, a uint8) Color {
	return Color{R: r, G: g, B: b, A: a}
}

// the whitest white that has ever whited
var ColorWhite = Rgba(255, 255, 255, 255)

// the blackest black that has ever blacked
var ColorBlack = Rgba(0, 0, 0, 255)

// transparent.
var ColorTransparent = Rgba(0, 0, 0, 0)