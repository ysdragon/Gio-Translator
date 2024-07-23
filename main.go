package main

import (
	"image"
	"image/color"
	"io"
	"os"
	"strings"
	"unicode"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"

	"Gio-Translator/sources"
)

type (
	// alias for Context.
	C = layout.Context
	// alias for Dimensions.
	D = layout.Dimensions
)

var (
	appName = "Gio Translator"
	list    = &widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}
	th                = material.NewTheme()
	topLabelColor     = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	globalColor       = color.NRGBA{R: 135, G: 39, B: 36, A: 255}
	errorColor        = color.NRGBAModel.Convert(color.RGBA{R: 255, G: 0, B: 0, A: 255}).(color.NRGBA)
	ops               op.Ops
	translateField    component.TextField
	fromField         component.TextField
	toField           component.TextField
	translateBtn      widget.Clickable
	radioButtonsGroup = new(widget.Enum)
	CopyButton        widget.Clickable

	log             string
	showcopyBtn     bool
	isError         bool
	translateSource string
	err             error
)

func main() {
	go func() {
		// Create a new window
		w := new(app.Window)
		// Set window options
		w.Option(
			// Set the window title
			app.Title(appName),
			// Set the window Size (width, height)
			app.Size(unit.Dp(385), unit.Dp(600)),
		)
		for {
			switch event := w.Event().(type) {
			case app.DestroyEvent:
				os.Exit(0)
			case app.FrameEvent:
				frame(app.NewContext(&ops, event), w)
				event.Frame(&ops)
			}
		}
	}()

	// Start the app
	app.Main()
}

// frame lays out the entire frame and returns the reusltant ops buffer.
func frame(gtx C, w *app.Window) D {
	// Theme options
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	th.Palette.Bg = color.NRGBA{R: 32, G: 34, B: 36, A: 255}
	th.Palette.Fg = color.NRGBA{R: 215, G: 218, B: 222, A: 255}
	th.Palette.ContrastBg = globalColor
	th.Palette.ContrastFg = color.NRGBA{R: 255, G: 255, B: 255, A: 255}

	// Set the background color
	macro := op.Record(gtx.Ops)
	rect := image.Rectangle{
		Max: image.Point{
			X: gtx.Constraints.Max.X,
			Y: gtx.Constraints.Max.Y,
		},
	}

	paint.FillShape(gtx.Ops, th.Palette.Bg, clip.Rect(rect).Op())
	background := macro.Stop()

	background.Add(gtx.Ops)

	widgets := []layout.Widget{
		func(gtx C) D {
			toplabel := material.H3(th, appName)
			toplabel.Color = topLabelColor
			toplabel.Alignment = text.Middle
			return toplabel.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return translateField.Layout(gtx, th, "Text")
			})

		},

		func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return fromField.Layout(gtx, th, "From (e.g. en)")
			})
		},
		func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return toField.Layout(gtx, th, "To (e.g. ar)")
			})
		},
		func(gtx C) D {
			if radioButtonsGroup.Update(gtx) {
				switch radioButtonsGroup.Value {
				case "mm":
					translateSource = "MyMemory"
				case "gt":
					translateSource = "GoogleTranslate"
				case "ai":
					translateSource = "AI"
				default:
					translateSource = "MyMemory"
				}
				gtx.Execute(op.InvalidateCmd{})
			}
			return layout.Center.Layout(gtx, func(gtx C) D {
				mmRButton := material.RadioButton(th, radioButtonsGroup, "mm", "MyMemory")
				mmRButton.IconColor = globalColor
				gtRButton := material.RadioButton(th, radioButtonsGroup, "gt", "Google Translate")
				gtRButton.IconColor = globalColor
				aiRButton := material.RadioButton(th, radioButtonsGroup, "ai", "AI")
				aiRButton.IconColor = globalColor
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(mmRButton.Layout),
					layout.Rigid(gtRButton.Layout),
					layout.Rigid(aiRButton.Layout),
				)
			})
		},
		func(gtx C) D {
			// Check if the translate button has been clicked
			if translateBtn.Clicked(gtx) {
				go func() {
					// Get text value from translateField
					text := translateField.Text()
					// Get text value from fromField
					from := fromField.Text()
					// Get text value from toField
					to := toField.Text()

					// Check
					if text == "" {
						log = "Text field cannot be empty!"
						isError = true
						showcopyBtn = false
						return
					} else if from == "" {
						log = "From language field cannot be empty!"
						isError = true
						showcopyBtn = false
						return
					} else if to == "" {
						log = "To language field cannot be empty!"
						isError = true
						showcopyBtn = false
						return
					}

					switch translateSource {
					case "MyMemory":
						log, err = sources.MyMemory(text, from, to)
						isError = false
						showcopyBtn = true
						if err != nil {
							log = "Error"
							isError = true
							showcopyBtn = false
							return
						}
					case "GoogleTranslate":
						log, err = sources.GTrasnlate(text, from, to)
						isError = false
						showcopyBtn = true
						if err != nil {
							log = "Error"
							isError = true
							showcopyBtn = false
							return
						}
					case "AI":
						log, err = sources.DuckDuckGoAiTranslate(text, from, to)
						isError = false
						showcopyBtn = true
						if err != nil {
							log = "Error"
							isError = true
							showcopyBtn = false
							return
						}
					default:
						log, err = sources.MyMemory(text, from, to)
						isError = false
						showcopyBtn = true
						if err != nil {
							log = "Error"
							isError = true
							showcopyBtn = false
							return
						}

					}

					w.Invalidate()
				}()
			}

			button := material.Button(th, &translateBtn, "Translate")
			button.Background = globalColor
			return button.Layout(gtx)

		},
		func(gtx C) D {
			result := material.Body1(th, log)

			if isArabic(log) {
				result.Alignment = text.End
			} else if isError {
				result.Color = errorColor
			}

			return result.Layout(gtx)
		},
		func(gtx C) D {
			// Check if showcopyBtn is true
			if showcopyBtn {
				// Check if the copy button has been clicked
				if CopyButton.Clicked(gtx) {
					gtx.Execute(clipboard.WriteCmd{
						Data: io.NopCloser(strings.NewReader(log)),
					})
				}
				// Now return the copy button
				copyButton := material.Button(th, &CopyButton, "Copy")
				copyButton.Background = globalColor
				return copyButton.Layout(gtx)
			}
			return D{}
		},
	}

	return material.List(th, list).Layout(gtx, len(widgets), func(gtx C, i int) D {
		maxWidth := gtx.Dp(unit.Dp(385))
		gtx.Constraints.Max.X = maxWidth
		return layout.UniformInset(unit.Dp(10)).Layout(gtx, widgets[i])
	})

}

// Function to detect if the given string is arabic.
func isArabic(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Arabic, r) {
			return true
		}
	}
	return false
}
