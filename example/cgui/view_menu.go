package main

import (
    "github.com/jroimartin/gocui"
    "fmt"
)

const (
    VIEW_MENU = "menu"
)

type Item struct {
    Name string
    Action    func() error
}
type MenuView struct {
    consoleView *gocui.View
    items       []Item
}

func NewMenuView() *MenuView {
    v := new(MenuView)
    v.items = make([]Item, 0)
    return v
}

func (v *MenuView) AddItem(name string, act func() error) {
    v.items = append(v.items, Item{
        Name: name,
        Action: act,
    })
}

func (v *MenuView) Name() string {
    return v.consoleView.Name()
}

func (v *MenuView) Render(g *gocui.Gui) error {
    if cv, err := g.SetView(VIEW_MENU, 1, 1, 25, 2 + len(v.items)); err != nil {

        if err != gocui.ErrUnknownView {
            return err
        }
        cv.Title = "Menu"
        cv.SelFgColor = gocui.ColorGreen | gocui.AttrUnderline
        cv.Highlight = true
        cv.Frame = true
        idx := 1
        for _, item := range v.items {
            fmt.Fprintln(cv, fmt.Sprintf("%02d %s",idx, item.Name))
            idx++
        }
        v.consoleView = cv
    }
    return nil
}


func (v *MenuView) runAction(g *gocui.Gui, cv *gocui.View) error {
    _, cy := cv.Cursor()
    line, _ := cv.Line(cy)
    for _, item := range v.items {
        if item.Name == line[3:] {
            item.Action()
        }
    }
    return nil

}
func (v *MenuView) KeyBindings (g *gocui.Gui) error {
    if err := g.SetKeybinding(VIEW_MENU, gocui.KeyEnter, gocui.ModNone, v.runAction); err != nil {
        return err
    }
    return nil
}