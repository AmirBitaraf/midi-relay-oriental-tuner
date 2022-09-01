package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func runUI() {
	a := app.New()
	w := a.NewWindow("MIDI Relay Tuner")

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))
	w.SetFixedSize(true)
	w.Resize(fyne.Size{Height: 400, Width: 600})
	w.SetIcon(resourceLogoPng)
	w.ShowAndRun()
}
