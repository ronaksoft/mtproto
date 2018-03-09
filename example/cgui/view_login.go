package main

import (
    "github.com/jroimartin/gocui"
    "fmt"
)

const (
    VIEW_LOGIN = "login"
)

type LoginView struct {
    consoleView *gocui.View
}

func NewLoginView() *LoginView {
    v := new(LoginView)
    return v

}
func (v *LoginView) Name() string {
    return v.consoleView.Name()
}

func (v *LoginView) Render(g *gocui.Gui) error {
    maxX, maxY := g.Size()
    if cv, err := g.SetView(VIEW_LOGIN, (maxX-30)/2, (maxY-10)/2, (maxX+30)/2, (maxY+10)/2); err != nil {

        if err != gocui.ErrUnknownView {
            return err
        }
        cv.Title = "Login"
        cv.SelFgColor = gocui.ColorGreen
        cv.Highlight = true

        fmt.Fprintln(cv, "Phone Number:")

        v.consoleView = cv
    }
    return nil
}

func (v *LoginView) KeyBindings (g *gocui.Gui) error  {
    return nil
}