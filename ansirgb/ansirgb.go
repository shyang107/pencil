// ref: https://github.com/stroborobo/ansirgb

package ansirgb

import (
	"fmt"
	"image/color"
)

// Color represents an ANSI color
type Color struct {
	color.Color
	Code int
}

// String returns a small description of the color
// containing ANSI code and a hex representation of
// the RGB value, e.g. " 26: 0000 4B4B D2D2"
func (c *Color) String() string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("%3d: \033[38;5;%dm%04X %04X %04X\033[0m", c.Code, c.Code, r, g, b)
}

// Fg returns the escape code for foreground color
func (c *Color) Fg() string {
	// What the heck are you trying to do here, dearest user?
	//if c.IsTransparent() {}
	return fmt.Sprintf("\033[38;5;%dm", c.Code)
}

// Bg returns the escape code for background color
func (c *Color) Bg() string {
	if c.IsTransparent() {
		return fmt.Sprintf("\033[49m")
	}
	return fmt.Sprintf("\033[48;5;%dm", c.Code)
}

// IsTransparent returns if the color is a transparent one
func (c *Color) IsTransparent() bool {
	return c.Code == -1
}

// Equals :
func (c *Color) Equals(rhs *Color) bool {
	return c.Code == rhs.Code
}

// CreatePalette creates a palette contains all 256 ANSI colors
func CreatePalette() color.Palette {
	// Palette contains all 256 ANSI colors
	var palette = make(color.Palette, 0, 255)

	var r uint8
	var g uint8
	var b uint8

	// ignore 8 Bit colors, they are very dependent on the user's terminal
	// while other colors usually aren't.

	// 6 * 6 * 6 = 216 colors
	var start uint8 = 95
	var step uint8 = 40
	r = 0
	g = 0
	b = 0
	for i := 16; i <= 231; i++ {
		palette = append(palette, &Color{&color.RGBA{r, g, b, 255}, i})
		b = next(b, start, step)
		if i > 16 {
			if (i-15)%6 == 0 {
				g = next(g, start, step)
			}
			if (i-15)%(6*6) == 0 {
				r = next(r, start, step)
			}
		}
	}

	// grayscale
	r = 8
	g = 8
	b = 8
	for i := 232; i <= 255; i++ {
		palette = append(palette, &Color{&color.RGBA{r, g, b, 255}, i})
		r += 10
		g += 10
		b += 10
	}

	// transparent
	palette = append(palette, &Color{&color.RGBA{0, 0, 0, 0}, -1})

	return palette
}

func next(n, start, step uint8) uint8 {
	if n == 0 {
		return start
	} else if n == 255 {
		return 0
	}
	return n + step
}

// Convert returns the ANSI color closest to c.
func Convert(c color.Color) *Color {
	return Palette.Convert(c).(*Color)
}

// Index returns the index of the ANSI palette color closest to c
func Index(c color.Color) int {
	return Convert(c).Code
}

// ANSIRGBModel can convert any Color to one from its own color model. The conversion may be lossy.
var ANSIRGBModel = color.ModelFunc(ansirgbaModel)

func ansirgbaModel(c color.Color) color.Color {
	// if _, ok := c.(color.RGBA); ok {
	// 	cl := Convert(c)
	// 	return cl
	// }
	cl := Convert(c)
	return cl
}
