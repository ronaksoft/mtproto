package main

import (
    "github.com/jroimartin/gocui"
    "log"
    "sync"
)

const (
    VIEW_ALL = ""
)

type App struct {
    sync.Mutex
    g *gocui.Gui

    // Views
    viewOrder     []string
    views         map[string]View
    curr_view     View
    curr_view_idx int
}

func NewApp() *App {
    app := new(App)
    if cui, err := gocui.NewGui(gocui.Output256); err != nil {
        log.Panic("NewApp::", err.Error())
    } else {
        app.g = cui
        //app.g.Cursor = true
    }

    if err := app.setKeyBindings(); err != nil {
        log.Panic("NewApp::KeyBindings::", err.Error())
    }

    return app
}

func (a *App) setKeyBindings() error {
    if err := a.g.SetKeybinding(VIEW_ALL, gocui.KeyCtrlC, gocui.ModNone, a.quit); err != nil {
        return err
    }

    if err := a.g.SetKeybinding(VIEW_ALL, gocui.KeyArrowUp, gocui.ModNone, a.cursorUp); err != nil {
        return err
    }

    if err := a.g.SetKeybinding(VIEW_ALL, gocui.KeyArrowDown, gocui.ModNone, a.cursorDown); err != nil {
        return err
    }

    if err := a.g.SetKeybinding(VIEW_ALL, gocui.KeyTab, gocui.ModNone, a.switchToNextView); err != nil {
        return err
    }
    return nil
}

func (a *App) managerFunction() {
    for _, v := range a.views {
        v.Render(a.g)
        //cv, _ := a.g.View(v.Name())
    }
}

func (a *App) cursorDown(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        _, vy := v.Size()
        cx, cy := v.Cursor()
        cy = (cy + 1) % vy
        if err := v.SetCursor(cx, cy); err != nil {
            ox, oy := v.Origin()
            if err := v.SetOrigin(ox, oy+1); err != nil {
                return err
            }
        }

    }
    return nil
}

func (a *App) cursorUp(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        ox, oy := v.Origin()
        cx, cy := v.Cursor()
        _, vy := v.Size()
        cy = cy - 1
        if cy < 0 {
            cy = vy - 1
        }
        if err := v.SetCursor(cx, cy); err != nil && oy > 0 {
            if err := v.SetOrigin(ox, oy-1); err != nil {
                return err
            }
        }
    }
    return nil
}

func (a *App) switchToNextView(g *gocui.Gui, v *gocui.View) error {
    var err error
    a.curr_view_idx = (a.curr_view_idx + 1) % len(a.views)
    a.curr_view = a.views[a.curr_view_idx]
    if _, err = a.g.SetCurrentView(a.curr_view.Name()); err != nil {
        return err
    }
    return nil
}

func (a *App) quit(g *gocui.Gui, v *gocui.View) error {
    return gocui.ErrQuit
}

func (a *App) AddView(views ...View) {
    for _, v := range views {
        a.viewOrder = append(a.viewOrder, v.Name())
        a.views[v.Name()] = v
        v.KeyBindings(a.g)
    }
}

func (a *App) SetMainView(v View) {


}

func (a *App) SetCurrentView(v View) {
    a.g.SetCurrentView(v.Name())
}

func (a *App) Run() {
    defer a.g.Close()

    a.managerFunction()

    if err := a.g.MainLoop(); err != nil && err != gocui.ErrQuit {
        log.Panicln(err)
    }
}
