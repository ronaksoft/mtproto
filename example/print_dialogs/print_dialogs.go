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
    appId := int64(48841)
    appHash := "3151c01673d412c18c055f089128be50"
    if v, err := mtproto.NewMTProto(appId, appHash, "../auth_key", "", 0); err != nil {
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
        fmt.Println("Dialog Type:", d.Type)
        fmt.Println("Top MessageID:", d.TopMessageID)
        userID := messages[d.TopMessageID].From
        fmt.Println("From:", users[userID].FirstName, users[userID].LastName, "@", users[userID].Username, "(", userID, ")")
        switch d.Type {
        case mtproto.DIALOG_TYPE_USER:
            fmt.Println("Peer Info:", users[d.PeerID].FirstName, users[d.PeerID].LastName)
        case mtproto.DIALOG_TYPE_CHAT:
            fmt.Println("Peer Info:", chats[d.PeerID].Title, chats[d.PeerID].Username, d.PeerID, d.PeerAccessHash)
        case mtproto.DIALOG_TYPE_CHANNEL:
            fmt.Println("Peer Info:", channels[d.PeerID].Title, channels[d.PeerID].Username)
        }

    }

}
