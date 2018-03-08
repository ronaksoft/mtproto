package get_updates

import (
    "github.com/ronaksoft/mtproto"
    "log"
    "fmt"
    "strings"
    "time"
    "github.com/kr/pretty"
)

var (
    _MT *mtproto.MTProto
)

func main() {
    appId := int64(48841)
    appHash := "3151c01673d412c18c055f089128be50"
    if v, err := mtproto.NewMTProto(appId, appHash, "./auth_key", "", 0); err != nil {
        log.Println(err.Error())
    } else {
        _MT = v
        if err := _MT.Connect(); err != nil {
            log.Println("Connect:", err.Error())
        }
    }

    updateState := _MT.Updates_GetState()
    fmt.Println("Update State Time:", time.Unix(int64(updateState.Date), 0).Format("2006-01-02 15:04:05"))
    pretty.Println(updateState)

    updateDifference := _MT.Updates_GetDifference(
        updateState.Pts-100,
        0,
        int32(time.Date(2017, time.Month(03), 07, 0, 0, 0, 0, time.Local).Unix()),
    )
    PrintUpdateDifference(updateDifference)

}

func LoadMessages(startPoint int) {
    fmt.Println("Loading Messages")

    msgIDs := make([]int32, 20)
    for {
        for j := 0; j < 20; j++ {
            msgIDs[j] = int32(startPoint + j)
        }
        messages, users, chats := _MT.Messages_GetMessages(msgIDs)

        for _, m := range messages {
            fmt.Println("==================================")
            fmt.Println(m.ID)
            fmt.Println(fmt.Sprintf("From:  %s %s", users[m.From].FirstName, users[m.From].LastName))
            switch m.To.Type {
            case mtproto.PEER_TYPE_USER:
                fmt.Println(fmt.Sprintf("To:    %s %s", users[m.To.ID].FirstName, users[m.To.ID].LastName))
            case mtproto.PEER_TYPE_CHAT:
                fmt.Println(fmt.Sprintf("To:    %s %s", chats[m.To.ID].Title, chats[m.To.ID].Username))
            }
            fmt.Println(fmt.Sprintf("%s %s", time.Unix(int64(m.Date), 0).Format("2006-01-02"), m.Body))
            fmt.Println("Out:", m.Flags.Out, "Unread:", m.Flags.MediaUnread)
        }
        time.Sleep(1)
        startPoint += 20
    }

}

func PrintUpdateDifference(updateDifference *mtproto.UpdateDifference) {
    fmt.Println("Total:", updateDifference.Total)
    fmt.Println("Intermediate State:", updateDifference.IntermediateState)

    fmt.Println("New Messages:")
    for _, m := range updateDifference.NewMessages {
        fmt.Println("-------------------------------")
        fmt.Println("MessageID:", m.ID, m.MediaType)
        fmt.Println("Time:", time.Unix(int64(m.Date), 0).Format("2006-01-02"))
        fmt.Println("From:", updateDifference.Users[m.From].FirstName, updateDifference.Users[m.From].LastName)
        switch m.To.Type {
        case mtproto.PEER_TYPE_CHAT:
            fmt.Println("To:", updateDifference.Chats[m.To.ID].Title, updateDifference.Chats[m.To.ID].Username, m.To.Type)
        case mtproto.PEER_TYPE_CHANNEL:
            fmt.Println("To:", updateDifference.Channels[m.To.ID].Title, updateDifference.Channels[m.To.ID].Username, m.To.Type)
        }
        fmt.Println(m.Body)
        time.Sleep(1 * time.Second)
    }

    fmt.Println("New Updates:")
    for _, u := range updateDifference.OtherUpdates {
        fmt.Println("=================================")
        fmt.Println("", u.Type, u.Date)
        fmt.Println(u.UserID, u.ChatID, u.ChannelID)
        time.Sleep(1 * time.Second)
    }
}

