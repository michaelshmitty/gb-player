package ui

import (
	"gb-player/gb"
)

type View struct {
	controller *Controller
	console    *gb.Console
}

func NewView(controller *Controller, console *gb.Console) *View {
	return &View{controller, console}
}
