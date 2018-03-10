package cmd

import (
    "github.com/spf13/cobra"
    "fmt"
    "strings"
    "time"
    "strconv"
    "github.com/ronaksoft/mtproto"
    "github.com/olekukonko/tablewriter"
    "os"
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
        minutes, _ := strconv.Atoi(cmd.Flag("minutes").Value.String())


        updateState := _MT.Updates_GetState()
        updateDifference := _MT.Updates_GetDifference(
            updateState.Pts-int32(numberOfUpdates),
            0,
            int32(time.Now().Add(- time.Duration(minutes) * time.Minute).Unix()),
        )

        tableState := tablewriter.NewWriter(os.Stdout)
        tableState.SetHeader([]string{"Date", "Pts", "Qts", "Seq", "Unread Counts"})
        tableState.Append([]string{
            fmt.Sprintf("%d", updateState.Date),
            fmt.Sprintf("%d", updateState.Pts),
            fmt.Sprintf("%d", updateState.Qts),
            fmt.Sprintf("%d", updateState.Seq),
            fmt.Sprintf("%d", updateState.UnreadCounts),

        })
        tableState.Append([]string{
            fmt.Sprintf("%d", updateDifference.IntermediateState.Date),
            fmt.Sprintf("%d", updateDifference.IntermediateState.Pts),
            fmt.Sprintf("%d", updateDifference.IntermediateState.Qts),
            fmt.Sprintf("%d", updateDifference.IntermediateState.Seq),
            fmt.Sprintf("%d", updateDifference.IntermediateState.UnreadCounts),
        })
        tableState.SetCaption(true, "Table 1. :: Update States")
        tableState.Render()
        fmt.Println()
        fmt.Println()

        tableMessages := tablewriter.NewWriter(os.Stdout)
        tableMessages.SetHeader([]string{"Index", "Message ID", "Time", "From", "To", "Body"})
        tableMessages.SetColMinWidth(4, 50)
        tableMessages.SetCaption(true, "Table 2. :: New Messages")

        idx := 0
        for _, m := range updateDifference.NewMessages {
            idx++
            tableRow := []string{
                fmt.Sprintf("%d", idx),
                fmt.Sprintf("%d", m.ID),
                time.Unix(int64(m.Date), 0).Format("2006-01-02 15:04:05"),
                fmt.Sprintf("%s %s", updateDifference.Users[m.From].FirstName, updateDifference.Users[m.From].LastName),
            }
            switch m.To.Type {
            case mtproto.PEER_TYPE_USER:
                tableRow = append(
                    tableRow,
                    fmt.Sprintf("%s %s",
                        updateDifference.Users[m.To.ID].FirstName,
                        updateDifference.Users[m.To.ID].LastName,
                    ),
                )
            case mtproto.PEER_TYPE_CHAT:
                tableRow = append(
                    tableRow,
                    fmt.Sprintf("%s(@%s) %s",
                        updateDifference.Chats[m.To.ID].Title,
                        updateDifference.Chats[m.To.ID].Username,
                        m.To.Type,
                    ),
                )
            case mtproto.PEER_TYPE_CHANNEL:
                tableRow = append(
                    tableRow,
                    fmt.Sprintf("%s(@%s) %s",
                        updateDifference.Channels[m.To.ID].Title,
                        updateDifference.Channels[m.To.ID].Username,
                        m.To.Type,
                    ),
                )
            }
            if len(m.Body) > 20 {
                tableRow = append(tableRow, m.Body[:20])
            } else {
                tableRow = append(tableRow, m.Body)
            }
            tableMessages.Append(tableRow)

        }
        tableMessages.Render()
        fmt.Println()
        fmt.Println()

        tableUpdates := tablewriter.NewWriter(os.Stdout)
        tableUpdates.SetHeader([]string{"Index", "Update Type", "Date", "UserID", "ChannelID", "ChatID", "MessageID", "Pts", "Pts Count"})
        idx = 0
        for _, u := range updateDifference.OtherUpdates {
            idx++
            tableRow := []string{
                fmt.Sprintf("%d", idx),
                u.Type,
            }
            if u.Date > 0 {
                tableRow = append(tableRow, time.Unix(int64(u.Date), 0).Format("2006-01-02 15:04:05"))
            } else {
                tableRow = append(tableRow, "No Time")
            }
            tableRow = append(tableRow,
                fmt.Sprintf("%d", u.UserID),
                fmt.Sprintf("%d", u.ChannelID),
                fmt.Sprintf("%d", u.ChatID),
                fmt.Sprintf("%d", u.MessageID),
                fmt.Sprintf("%d", u.Pts),
                fmt.Sprintf("%d", u.PtsCount),
            )
            tableUpdates.Append(tableRow)
        }

        tableUpdates.SetCaption(true, "Table 3. :: Other Updates")
        tableUpdates.Render()
        fmt.Println()
        fmt.Println()
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
    GetUpdatesCmd.Flags().Int("minutes", 10, "")
    GetDialogsCmd.Flags().String("peerType", "", "")
}
