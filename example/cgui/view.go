package main

import "github.com/jroimartin/gocui"

type View interface {
    Name () string
    Render(g *gocui.Gui) error
    KeyBindings (g *gocui.Gui) error
}