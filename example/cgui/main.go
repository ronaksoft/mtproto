package main

func main() {
    app := NewApp()
    menuView := NewMenuView()
    loginView := NewLoginView()

    menuView.AddItem("Login", func() error {
        loginView := NewLoginView()
        loginView.Render(app.g)
        app.SetCurrentView(loginView)
        return nil
    })

    menuView.AddItem("Get Updates", func() error {
        return nil
    })

    menuView.AddItem("Get Dialogs", func() error {
        return nil
    })

    app.AddView(menuView, loginView)
    app.Run()

}
