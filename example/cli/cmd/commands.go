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

var (
    _fPhone          *string
    _fPeerType       *string
    _fPeerID         *int32
    _fMaxID          *int32
    _fMinID          *int32
    _fPeerAccessHash *int64
    _fLimit          *int32
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
        tableState.SetHeader([]string{"Type", "Date", "Pts", "Qts", "Seq", "Unread Counts"})
        tableState.Append([]string{
            "Update State",
            fmt.Sprintf("%d", updateState.Date),
            fmt.Sprintf("%d", updateState.Pts),
            fmt.Sprintf("%d", updateState.Qts),
            fmt.Sprintf("%d", updateState.Seq),
            fmt.Sprintf("%d", updateState.UnreadCounts),
        })
        tableState.Append([]string{
            updateDifference.Type,
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

        idx := 0
        tableMessages := tablewriter.NewWriter(os.Stdout)
        tableUpdates := tablewriter.NewWriter(os.Stdout)
        messageRows := make([][]string, 0)
        updateRows := make([][]string, 0)
        for {
            if len(updateDifference.NewMessages) > 0 {
                tableMessages.SetHeader([]string{"Index", "Message ID", "Time", "From", "To", "Body"})
                tableMessages.SetColMinWidth(4, 50)
                tableMessages.SetCaption(true, "Table 2. :: New Messages")
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
                    messageRows = append(messageRows, tableRow)
                    //tableMessages.Append(tableRow)

                }

            }

            if len(updateDifference.OtherUpdates) > 0 {

                tableUpdates.SetHeader([]string{"Index", "Update Type", "Date", "UserID", "ChannelID", "ChatID", "MessageID", "Pts", "Pts Count"})
                tableUpdates.SetCaption(true, "Table 3. :: Other Updates")
                for _, u := range updateDifference.OtherUpdates {
                    idx++
                    tableRow := []string{
                        fmt.Sprintf("%d", idx),
                        u.GetType(u),
                    }
                    if v, ok := u.GetInt32(u, "Date"); ok {
                        tableRow = append(tableRow, time.Unix(int64(v), 0).Format("2006-01-02 15:04:05"))
                    } else {
                        tableRow = append(tableRow, "No Time")
                    }

                    if userID, ok := u.GetInt32(u, "UserID"); ok {
                        tableRow = append(tableRow, fmt.Sprintf("%d", userID))
                    } else {
                        tableRow = append(tableRow, "-")
                    }
                    if channelID, ok := u.GetInt32(u, "ChannelID"); ok {
                        tableRow = append(tableRow, fmt.Sprintf("%d", channelID))
                    } else {
                        tableRow = append(tableRow, "-")
                    }
                    if chatID, ok := u.GetInt32(u, "ChatID"); ok {
                        tableRow = append(tableRow, fmt.Sprintf("%d", chatID))
                    } else {
                        tableRow = append(tableRow, "-")
                    }
                    if messageID, ok := u.GetInt32(u, "MessageID"); ok {
                        tableRow = append(tableRow, fmt.Sprintf("%d", messageID))
                    } else {
                        tableRow = append(tableRow, "-")
                    }
                    if pts, ok := u.GetInt32(u, "Pts"); ok {
                        tableRow = append(tableRow, fmt.Sprintf("%d", pts))
                    } else {
                        tableRow = append(tableRow, "-")
                    }
                    if ptsCount, ok := u.GetInt32(u, "UserID"); ok {
                        tableRow = append(tableRow, fmt.Sprintf("%d", ptsCount))
                    } else {
                        tableRow = append(tableRow, "-")
                    }

                    updateRows = append(updateRows, tableRow)
                }
            }

            if !updateDifference.IsSlice {
                break
            }
            updateDifference = _MT.Updates_GetDifference(
                updateDifference.IntermediateState.Pts,
                updateDifference.IntermediateState.Qts,
                updateDifference.IntermediateState.Date,
            )
        }

        tableMessages.AppendBulk(messageRows)
        tableMessages.Render()
        fmt.Println()
        fmt.Println()

        tableUpdates.AppendBulk(updateRows)
        tableUpdates.Render()
        fmt.Println()
        fmt.Println()
    },
}

var GetDialogsCmd = &cobra.Command{
    Use: "getDialogs",
    Run: func(cmd *cobra.Command, args []string) {
        dialogs, users, chats, channels, messages, dialogsCount := _MT.Messages_GetDialogs(0, int32(time.Now().Unix()), *_fLimit, mtproto.TL_inputPeerSelf{})
        fmt.Println("Total Dialogs (fetched/all):", len(dialogs), dialogsCount)

        _ = chats
        _ = channels

        tableDialogs := tablewriter.NewWriter(os.Stdout)
        tableDialogs.SetHeader([]string{"Index", "Peer Type", "Peer ID", "AccessHash", "Date", "Last Sender ID", "Last Sender", "UnRead Count"})
        tableDialogs.SetCaption(true, "Table 1. :: Dialogs")

        idx := 0
        for _, d := range dialogs {
            if len(*_fPeerType) > 0 && *_fPeerType != d.Type {
                continue
            }
            idx++
            userID := messages[d.TopMessageID].From
            tableDialogs.Append([]string{
                fmt.Sprintf("%d", idx),
                d.Type,
                fmt.Sprintf("%d", d.PeerID),
                fmt.Sprintf("%d", d.PeerAccessHash),
                time.Unix(int64(messages[d.TopMessageID].Date), 0).Format("2006-01-02 15:04:05"),
                fmt.Sprintf("%d", userID),
                fmt.Sprintf("%s %s", users[messages[d.TopMessageID].From].FirstName, users[messages[d.TopMessageID].From].LastName),
                fmt.Sprintf("%d", d.UnreadCount),
            })
        }
        tableDialogs.Render()
    },
}

var GetHistoryCmd = &cobra.Command{
    Use: "getHistory",
    Run: func(cmd *cobra.Command, args []string) {
        var inputPeer mtproto.TL
        switch *_fPeerType {
        case mtproto.PEER_TYPE_USER:
            inputPeer = mtproto.NewUserInputPeer(*_fPeerID, *_fPeerAccessHash)
        case mtproto.PEER_TYPE_CHAT:
            inputPeer = mtproto.NewChatInputPeer(*_fPeerID)
        case mtproto.PEER_TYPE_CHANNEL:
        default:
            return
        }
        messages, _ := _MT.Messages_GetHistory(inputPeer, *_fLimit, *_fMinID, *_fMaxID)
        tableMessages := tablewriter.NewWriter(os.Stdout)
        tableMessages.SetHeader([]string{"Index", "Message ID", "From ID", "Date", "Flags", "Body"})
        tableMessages.SetCaption(true, "Table 1. Messages")
        idx := 0
        for _, msg := range messages {
            idx++
            tableMessages.Append([]string{
                fmt.Sprintf("%d", idx),
                fmt.Sprintf("%d", msg.ID),
                fmt.Sprintf("%d", msg.From),
                time.Unix(int64(msg.Date), 0).Format("2006-01-02 15:04:05"),
                fmt.Sprintf("Out(%t) MediaUnread(%t) Post(%t)", msg.Flags.Out, msg.Flags.MediaUnread, msg.Flags.Post),
                fmt.Sprintf("%s", msg.Body),
            })
        }
        tableMessages.Render()
    },
}

var ReadHistoryCmd = &cobra.Command{
    Use: "readHistory",
    Run: func(cmd *cobra.Command, args []string) {
        var inputPeer mtproto.TL
        switch *_fPeerType {
        case mtproto.PEER_TYPE_USER:
            inputPeer = mtproto.NewUserInputPeer(*_fPeerID, *_fPeerAccessHash)
        case mtproto.PEER_TYPE_CHAT:
            inputPeer = mtproto.NewChatInputPeer(*_fPeerID)
        case mtproto.PEER_TYPE_CHANNEL:
        default:
            return
        }
        _MT.Messages_ReadHistory(inputPeer, *_fMaxID)
    },
}

var SendMessageCmd = &cobra.Command{
    Use: "send",
    Run: func(cmd *cobra.Command, args []string) {
        var peer mtproto.TL
        switch *_fPeerType {
        case mtproto.PEER_TYPE_USER:
            peer = mtproto.NewUserInputPeer(*_fPeerID, *_fPeerAccessHash)
        case mtproto.PEER_TYPE_CHAT:
            peer = mtproto.NewChatInputPeer(*_fPeerID)
        }
        _MT.Messages_SendMessage("Test", peer, 0)
    },
}

func init() {
    RootCmd.AddCommand(
        LoginCmd, GetUpdatesCmd, GetDialogsCmd, GetHistoryCmd, ReadHistoryCmd,
        SendMessageCmd,
    )
    GetUpdatesCmd.Flags().Int("numberOfUpdates", 10, "")
    GetUpdatesCmd.Flags().Int("minutes", 10, "")

    _fPhone = LoginCmd.Flags().String("phone", "989121228718", "")
    _fPeerType = RootCmd.PersistentFlags().String("peerType", "", "")
    _fPeerID = RootCmd.PersistentFlags().Int32("peerID", 0, "")
    _fPeerAccessHash = RootCmd.PersistentFlags().Int64("peerAccessHash", 0, "")
    _fLimit = RootCmd.PersistentFlags().Int32("limit", 10, "")
    _fMaxID = RootCmd.PersistentFlags().Int32("maxID", 0, "")
    _fMinID = RootCmd.PersistentFlags().Int32("minID", 0, "")

}
