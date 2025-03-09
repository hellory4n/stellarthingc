// ui :)
package ui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/hellory4n/stellarthing/core"
	"github.com/hellory4n/stellarthing/platform"
)

type ButtonStyle int
const (
	ButtonStylePrimary ButtonStyle = iota
	ButtonStyleSecondary
	ButtonStyleDanger
)

// makes a button :)
func Button(pos core.Vec2, size core.Vec2, text string, style ButtonStyle, onPressed func()) {	
	// center text :)
	var textSize rl.Vector2 = rl.MeasureTextEx(boldFont, text, float32(DefaultFontSize), 1)
	var textPos core.Vec2 = core.NewVec2(
		(size.X / 2) - (float64(textSize.X) / 2),
		(size.Y / 2) - (float64(textSize.Y) / 2),
	)
	var textColor core.Color = core.ColorWhite

	// base texture
	switch style {
		case ButtonStylePrimary: DrawPrimaryButton(pos, size)
		case ButtonStyleSecondary: DrawSecondaryButton(pos, size)
		case ButtonStyleDanger: DrawDangerButton(pos, size)
		default: panic("MATE\n")
	}

	// is it focused?
	buttonRect := core.NewRect(pos.X, pos.Y, size.X, size.Y)
	if buttonRect.HasPoint(platform.MousePosition()) {
		DrawHoverOutline(pos, size)

		// and is it pressed?
		if platform.IsMouseButtonJustReleased(platform.MouseButtonLeft) {
			onPressed()
		}
		// different mouse thingy so its not highlighted for 1 frame
		if platform.IsMouseButtonHeld(platform.MouseButtonLeft) {
			textColor = core.ColorBlack
			DrawPressedButton(pos, size)
		}
	}

	// finally draw text :)
	DrawBoldText(text, textPos.Add(pos), DefaultFontSize, textColor)
}