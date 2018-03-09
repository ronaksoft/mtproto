package cmd

import (
    "github.com/spf13/cobra"
    "fmt"
    "strings"
    "github.com/kr/pretty"
    "time"
    "strconv"
    "github.com/ronaksoft/mtproto"
)

var LoginCmd = &cobra.Command{
    Use: "login",
    Run: func(cmd *cobra.Command, args []string) {
        phone := cmd.Flag("phone").Value.String()
        if phoneCodeHash, err := _MT.Auth_SendCode(phone); err != nil {
            fmt.Println("SendCode:", err.Error())
        } else {
            var phoneCode string
            fmt.Print("Enter Code:")
            fmt.Scanln(&phoneCode)
            phoneCode = strings.TrimSpace(phoneCode)
            _MT.Auth_SignIn(phone, phoneCodeHash, phoneCode)
        }
    },
}

var GetUpdatesCmd = &cobra.Command{
    Use: "getUpdates",
    Run: func(cmd *cobra.Command, args []string) {
        numberOfUpdates, _ := strconv.Atoi(cmd.Flag("numberOfUpdates").Value.String())
        updateState := _MT.Updates_GetState()
        fmt.Println("Update State Time:", time.Unix(int64(updateState.Date), 0).Format("2006-01-02 15:04:05"))
        pretty.Println(updateState)

        updateDifference := _MT.Updates_GetDifference(
            updateState.Pts-int32(numberOfUpdates),
            0,
            int32(time.Date(2017, time.Month(03), 07, 0, 0, 0, 0, time.Local).Unix()),
        )

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

    },
}

var GetDialogsCmd = &cobra.Command{
    Use: "getDialogs",
    Run: func(cmd *cobra.Command, args []string) {
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

    },
}

func init() {
    RootCmd.AddCommand(LoginCmd, GetUpdatesCmd, GetDialogsCmd)
    LoginCmd.Flags().String("phone", "989121228718", "")
    GetUpdatesCmd.Flags().Int("numberOfUpdates", 10, "")
    GetDialogsCmd.Flags().String("peerType", "", "")
}
