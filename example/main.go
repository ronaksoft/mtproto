package main

import (
    "github.com/ronaksoft/mtproto"
    "log"
    "fmt"
    "strings"
    "time"
)

var (
    _MT *mtproto.MTProto
)

func main() {
    if v, err := mtproto.NewMTProto("./auth_key", "", 0); err != nil {
        log.Println(err.Error())
        _MT = v
        if err := _MT.Connect(); err != nil {
            log.Println("Connect:", err.Error())
        }
        if phoneCodeHash, err := _MT.Auth_SendCode("989121228718"); err != nil {
            log.Println("SendCode:", err.Error())
        } else {
            var phoneCode string
            fmt.Print("Enter Code:")
            fmt.Scanln(&phoneCode)
            phoneCode = strings.TrimSpace(phoneCode)
            fmt.Println("Code:", phoneCode)
            _MT.Auth_SignIn("989121228718", phoneCodeHash, phoneCode)
        }
    } else {
        _MT = v
        if err := _MT.Connect(); err != nil {
            log.Println("Connect:", err.Error())
        }
    }

    // Get Update State
    //updateState := _MT.Updates_GetState()
    //fmt.Println("Update State Time:", time.Unix(int64(updateState.Date), 0).Format("2006-01-02 15:04:05"))
    //pretty.Println(updateState)
    //updateDifference := _MT.Updates_GetDifference(updateState.Pts, updateState.Qts, updateState.Date)
    //
    //fmt.Println("New Messages:")
    //for _, m := range updateDifference.NewMessages {
    //    fmt.Println(m)
    //}
    //fmt.Println("New Updates:")
    //for _, u := range updateDifference.OtherUpdates {
    //    fmt.Println(u)
    //}

    fmt.Println("Loading Messages")
    i := 200000
    msgIDs := make([]int32, 20)
    for {
        for j := 0; j < 20; j++ {
            msgIDs[j] = int32(i + j)
        }
        messages, users, chats := _MT.Messages_GetMessages(msgIDs)

        for _, m := range messages {
            fmt.Println("==================================")
            fmt.Println(m.ID)
            fmt.Println(fmt.Sprintf("From:  %s %s",users[m.From].FirstName, users[m.From].LastName))
            switch m.To.Type {
            case mtproto.PEER_TYPE_USER:
                fmt.Println(fmt.Sprintf("To:    %s %s", users[m.To.ID].FirstName, users[m.To.ID].LastName))
            case mtproto.PEER_TYPE_CHAT:
                fmt.Println(fmt.Sprintf("To:    %s %s", chats[m.To.ID].Title, chats[m.To.ID].Username))
            }
            fmt.Println(fmt.Sprintf("%s %s",  time.Unix(int64(m.Date), 0).Format("2006-01-02"), m.Body))
            fmt.Println("Out:", m.Flags.Out, "Unread:", m.Flags.MediaUnread)
        }
        fmt.Scanln()
        i += 20
    }


}
