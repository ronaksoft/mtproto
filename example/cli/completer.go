package main

import (
    "github.com/c-bata/go-prompt"
)

var (
    _COMMANDS []Command
)

type CommandFlag struct {
    Name        string
    Description string
}

type Command struct {
    Name        string
    Description string
    Flags       []CommandFlag
    SubCommands []Command
}

func init() {
    _COMMANDS = []Command{
        {
            "tg",
            "Telegram Client",
            []CommandFlag{
                {"--account", "Choose a name for this account (default: unknown)"},
            },
            []Command{
                {
                    "login",
                    "Login to telegram using provided phone",
                    []CommandFlag{
                        {"--phone", "The phone number you are trying to login"},
                    },
                    []Command{},
                },
                {
                    "getUpdates",
                    "Get the last 'numberOfUpdates' updates",
                    []CommandFlag{
                        {"--numberOfUpdates", "The number of last updates"},
                        {"--minutes", "Get updates from n minutes ago"},
                    },
                    []Command{},
                },
                {
                    "getDialogs",
                    "Get dialogs of the logged in user",
                    []CommandFlag{
                        {"--peerType", "Filter Peer Type"},
                    },
                    []Command{},
                },
            },
        },
    }
}

func getCommandFlags(args []string, w string) []prompt.Suggest {
    if len(args) < 1 {
        return []prompt.Suggest{}
    }
    suggests := make([]prompt.Suggest, 0)
    cmdList := _COMMANDS
    for _, arg := range args {
        found := false
        for _, cmd := range cmdList {
            if cmd.Name == arg {
                cmdList = cmd.SubCommands
                for _, flag := range cmd.Flags {
                    suggests = append(suggests, prompt.Suggest{
                        flag.Name, flag.Description,
                    })
                }
                found = true
                break
            }
        }
        if !found {
            break
        }
    }
    return prompt.FilterContains(suggests, w, true)

}

func getCommands(args []string, w string) []prompt.Suggest {
    suggests := make([]prompt.Suggest, 0)
    cmdList := _COMMANDS
    for _, arg := range args {
        found := false
        for _, cmd := range cmdList {
            if cmd.Name == arg {
                cmdList = cmd.SubCommands
                found = true
                break
            }
        }
        if !found {
            break
        }
    }
    for _, cmd := range cmdList {
        suggests = append(suggests, prompt.Suggest{
            cmd.Name, cmd.Description,
        })
    }

    return prompt.FilterContains(suggests, w, true)
}
