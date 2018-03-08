package main

import (
    "github.com/ronaksoft/mtproto"
    "log"
    "time"
    "fmt"
)

var (
    _MT *mtproto.MTProto
)

func main() {
    if v, err := mtproto.NewMTProto("../auth_key", "", 0); err != nil {
        log.Println(err.Error())
    } else {
        _MT = v
        if err := _MT.Connect(); err != nil {
            log.Println("Connect:", err.Error())
        }
    }

    dialogs, users, chats, channels, messages, dialogsCount := _MT.Messages_GetDialogs(0, int32(time.Now().Unix()), 100, mtproto.TL_inputPeerSelf{})
    fmt.Println("Total Dialogs (fetched/all):", len(dialogs), dialogsCount)
    for _, d := range dialogs {
        fmt.Println("===============================")
        fmt.Println("Dialog Type:", d.Type, "Top MessageID:", d.TopMessageID)
        switch d.Type {
        case mtproto.DIALOG_TYPE_USER:
            fmt.Println("From:", d.User.FirstName, d.User.LastName, "Username:", d.User.Username, "(", d.User.ID, ")")
        case mtproto.DIALOG_TYPE_CHAT:
            fmt.Println("From:", d.User.FirstName, d.User.LastName, "Username:", d.User.Username, "(", d.User.ID, ")")
            fmt.Println("Chat Title:", d.Chat.Title, d.Chat.Username, d.Chat.ID, d.Chat.AccessHash)
        case mtproto.DIALOG_TYPE_CHANNEL:
            fmt.Println("From:", d.User.FirstName, d.User.LastName, "Username:", d.User.Username, "(", d.User.ID, ")")
        }

    }

}
