package main

import (
    "os"
    "github.com/ronaksoft/mtproto/example/cli/cmd"
    "strings"
    "github.com/c-bata/go-prompt"
    "fmt"
)

func main() {
    fmt.Println("Press Ctrl+D to exit")
    p := prompt.New(
        cmdExecutor,
        cmdCompleter,
        prompt.OptionTitle("Telegram Client"),
        prompt.OptionPrefix(">>> "),
        prompt.OptionInputTextColor(prompt.Blue),
        prompt.OptionSuggestionBGColor(prompt.Black),
        prompt.OptionSuggestionTextColor(prompt.White),
        prompt.OptionSelectedSuggestionTextColor(prompt.DarkGreen),
    )
    p.Run()

}

func cmdCompleter(d prompt.Document) []prompt.Suggest {
    var args []string
    w := d.GetWordBeforeCursor()

    // Suggests flags
    for _, arg := range strings.Split(d.TextBeforeCursor(), " ") {
        if !strings.HasPrefix(arg, "-") {
            args = append(args, arg)
        }
    }

    // Suggest Commands
    if strings.HasPrefix(w, "--") {
        return getCommandFlags(args, w)
    }

    return getCommands(args, w)
}

func cmdExecutor(in string) {
    args := strings.Fields(in)
    os.Args = append([]string{}, args...)
    if err := cmd.RootCmd.Execute(); err != nil {
        fmt.Println(err.Error())
    }
}



