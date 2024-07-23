package main

import (
	"image/color"
	"io"
	"os"
	"strings"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/op"
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
	topLabelColor     = color.NRGBAModel.Convert(color.RGBA{R: 0, G: 0, B: 255, A: 255}).(color.NRGBA)
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
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	widgets := []layout.Widget{

		func(gtx C) D {
			label := material.H3(th, appName)
			label.Color = topLabelColor
			label.Alignment = text.Middle
			return label.Layout(gtx)
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
				default:
					translateSource = "MyMemory"
				}
				gtx.Execute(op.InvalidateCmd{})
			}
			return layout.Center.Layout(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(material.RadioButton(th, radioButtonsGroup, "mm", "MyMemory").Layout),
					layout.Rigid(material.RadioButton(th, radioButtonsGroup, "gt", "Google Translate").Layout),
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

			return material.Button(th, &translateBtn, "Translate").Layout(gtx)

		},
		func(gtx C) D {
			result := material.Body1(th, log)
			if isError {
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
				return material.Button(th, &CopyButton, "Copy").Layout(gtx)
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
