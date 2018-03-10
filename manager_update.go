package mtproto

import (
    "reflect"
    "log"
    "fmt"
    "github.com/fatih/structs"
)

const (
    UPDATE_TYPE_NEW_MESSAGE                    = "NewMessage"
    UPDATE_TYPE_CHANNEL_NEW_MESSAGE            = "ChannelNewMessage"
    UPDATE_TYPE_READ_CHANNEL_INBOX             = "ReadChannelInbox"
    UPDATE_TYPE_READ_CHANNEL_OUTBOX            = "ReadChannelOutbox"
    UPDATE_TYPE_CHANNEL_TOO_LONG               = "ChannelTooLong"
    UPDATE_TYPE_READ_HISTORY_INBOX             = "ReadHistoryInbox"
    UPDATE_TYPE_READ_HISTORY_OUTBOX            = "ReadHistoryOutbox"
    UPDATE_TYPE_USER_PHOTO                     = "UserPhoto"
    UPDATE_TYPE_EDIT_MESSAGE                   = "EditMessage"
    UPDATE_TYPE_EDIT_CHANNEL_MESSAGE           = "EditChannelMessage"
    UPDATE_TYPE_CONTACT_LINK                   = "ContactLink"
    UPDATE_TYPE_DRAFT_MESSAGE                  = "DraftMessage"
    UPDATE_TYPE_SAVED_GIFS                     = "SavedGIFs"
    UPDATE_TYPE_MESSAGE_ID                     = "MessageID"
    UPDATE_TYPE_DELETE_MESSAGES                = "DeleteMessages"
    UPDATE_TYPE_CONTACT_REGISTERED             = "ContactRegistered"
    UPDATE_TYPE_USER_BLOCKED                   = "UserBlocked"
    UPDATE_TYPE_CHANNEL_READ_MESSAGES_CONTENTS = "ChannelReadMessagesContents"
    UPDATE_TYPE_USER_TYPING                    = "UserTyping"
    UPDATE_TYPE_CHAT_PARTICIPANT_ADD           = "ChatParticipantAdd"
    UPDATE_TYPE_CHAT_PARTICIPANT_ADMIN         = "ChatParticipantAdmin"
    UPDATE_TYPE_CHAT_PARTICIPANT_DELETE        = "ChatParticipantDelete"
    UPDATE_TYPE_CHAT_USER_TYPING               = "ChatUserTyping"
)

const (
    UPDATE_DIFFERENCE_EMPTY    = "EMPTY"
    UPDATE_DIFFERENCE_COMPLETE = "COMPLETE"
    UPDATE_DIFFERENCE_SLICE    = "SLICE"
    UPDATE_DIFFERENCE_TOO_LONG = "TOO_LONG"
)

type IUpdate interface {
    GetType(IUpdate) string
    GetInt32(IUpdate, string) (int32, bool)
    GetString(IUpdate, string) (string, bool)
    GetMap(IUpdate) map[string]interface{}
}

type UpdateCore struct{}

func (u UpdateCore) GetInt32(i IUpdate, keyName string) (int32, bool) {
    if f, ok := structs.New(i).FieldOk(keyName); !ok {
        return 0, false
    } else {
        return f.Value().(int32), true
    }
}
func (u UpdateCore) GetString(i IUpdate, keyName string) (string, bool) {
    if f, ok := structs.New(i).FieldOk(keyName); !ok {
        return "", false
    } else {
        return f.Value().(string), true
    }
}
func (u UpdateCore) GetType(i IUpdate) string {
    return reflect.TypeOf(i).String()
}
func (u UpdateCore) GetMap(i IUpdate) map[string]interface{} {
    return structs.Map(i)
}

type UpdateNewMessage struct {
    UpdateCore
    Pts      int32
    PtsCount int32
    Message  Message
}
type UpdateNewChannelMessage struct {
    UpdateCore
    Pts      int32
    PtsCount int32
    Message  Message
}
type UpdateReadChannelInbox struct {
    UpdateCore
    ChannelID int32
    MaxID     int32
}
type UpdateReadChannelOutbox struct {
    UpdateCore
    ChannelID int32
    MaxID     int32
}
type UpdateChannelTooLong struct {
    UpdateCore
    ChannelID int32
    Pts       int32
    Flags     int32
}
type UpdateReadHistoryInbox struct {
    UpdateCore
    Pts      int32
    PtsCount int32
    MaxID    int32
    Peer     Peer
}
type UpdateReadHistoryOutbox struct {
    UpdateCore
    Pts      int32
    PtsCount int32
    MaxID    int32
    Peer     Peer
}
type UpdateContactLink struct {
    UpdateCore
    UserID      int32
    MyLink      string
    ForeignLink string
}
type UpdateContactRegistered struct {
    UpdateCore
    UserID int32
    Date   int32
}
type UpdateUserPhoto struct {
    UpdateCore
    UserID       int32
    Date         int32
    ProfilePhoto UserProfilePhoto
    Previous     bool
}
type UpdateEditChannelMessage struct {
    UpdateCore
    Pts      int32
    PtsCount int32
    Message  Message
}
type UpdateEditMessage struct {
    UpdateCore
    Pts      int32
    PtsCount int32
    Message  Message
}
type UpdatedSaveGIFs struct {
    UpdateCore
}
type UpdateMessageID struct {
    UpdateCore
    MessageID int32
    RandomID  int64
}
type UpdateDeleteMessages struct {
    UpdateCore
    Pts        int32
    PtsCount   int32
    MessageIDs []int32
}
type UpdateDraftMessage struct {
    UpdateCore
    Peer  Peer
    Draft DraftMessage
}
type UpdateUserBlocked struct {
    UpdateCore
    UserID  int32
    Blocked bool
}
type UpdateChannelReadMessagesContents struct {
    UpdateCore
    MessageIDs []int32
    ChannelID  int32
}
type UpdateUnknown struct {
    UpdateCore
    Type string
}


type UpdateState struct {
    Qts          int32
    Pts          int32
    Date         int32
    Seq          int32
    UnreadCounts int32
}
type UpdateDifference struct {
    Type              string
    IsSlice           bool
    Total             int32
    NewMessages       []Message
    OtherUpdates      []IUpdate
    Chats             map[int32]Chat
    Channels          map[int32]Channel
    Users             map[int32]User
    IntermediateState UpdateState
    Seq               int32
}
type ChannelUpdateDifference struct {
    Empty        bool
    TooLong      bool
    Flags        int32
    Final        bool
    Pts          int32
    Timeout      int32
    NewMessages  []Message
    OtherUpdates []IUpdate
}

// NewUpdateState
// input :
//	1. TL_updates_state
func NewUpdateState(input TL) *UpdateState {
    us := new(UpdateState)
    switch in := input.(type) {
    case TL_updates_state:
        us.Qts = in.Qts
        us.Pts = in.Pts
        us.Seq = in.Seq
        us.Date = in.Date
        us.UnreadCounts = in.Unread_count
    }
    return us
}

// NewUpdate
// input :
//	1. TL_updateNewMessage
//	2. TL_updateNewChannelMessage
func NewUpdate(input TL) IUpdate {
    switch u := input.(type) {
    case TL_updateNewMessage:
        return UpdateNewMessage{
            Pts:      u.Pts,
            PtsCount: u.Pts_count,
            Message:  *NewMessage(u.Message),
        }
    case TL_updateNewChannelMessage:
        return UpdateNewChannelMessage{
            Pts:      u.Pts,
            PtsCount: u.Pts_count,
            Message:  *NewMessage(u.Message),
        }
    case TL_updateReadChannelInbox:
        return UpdateReadChannelInbox{
            ChannelID: u.Channel_id,
            MaxID:     u.Max_id,
        }
    case TL_updateReadChannelOutbox:
        return UpdateReadChannelOutbox{
            ChannelID: u.Channel_id,
            MaxID:     u.Max_id,
        }
    case TL_updateChannelTooLong:
        return UpdateChannelTooLong{
            Pts:       u.Pts,
            ChannelID: u.Channel_id,
            Flags:     u.Flags,
        }
    case TL_updateReadHistoryInbox:
        return UpdateReadHistoryInbox{
            Pts:      u.Pts,
            PtsCount: u.Pts_count,
            MaxID:    u.Max_id,
            Peer:     *NewPeer(u.Peer),
        }
    case TL_updateReadHistoryOutbox:
        return UpdateReadHistoryOutbox{
            Pts:      u.Pts,
            PtsCount: u.Pts_count,
            MaxID:    u.Max_id,
            Peer:     *NewPeer(u.Peer),
        }
    case TL_updateUserPhoto:
        var previous bool
        switch u.Previous.(type) {
        case TL_boolFalse:
            previous = false
        case TL_boolTrue:
            previous = true
        }
        return UpdateUserPhoto{
            UserID:       u.User_id,
            Date:         u.Date,
            ProfilePhoto: *NewUserProfilePhoto(u.Photo),
            Previous:     previous,
        }
    case TL_updateContactLink:
        return UpdateContactLink{
            UserID:      u.User_id,
            MyLink:      reflect.TypeOf(u.My_link).String(),
            ForeignLink: reflect.TypeOf(u.Foreign_link).String(),
        }
    case TL_updateEditChannelMessage:
        return UpdateEditChannelMessage{
            Pts:      u.Pts,
            PtsCount: u.Pts_count,
            Message:  *NewMessage(u.Message),
        }
    case TL_updateEditMessage:
        return UpdateEditMessage{
            Pts:      u.Pts,
            PtsCount: u.Pts_count,
            Message:  *NewMessage(u.Message),
        }
    case TL_updateSavedGifs:
        return UpdatedSaveGIFs{}
    case TL_updateDraftMessage:
        return UpdateDraftMessage{
            Peer:  *NewPeer(u.Peer),
            Draft: *NewDraftMessage(u.Draft),
        }
    case TL_updateMessageID:
        return UpdateMessageID{
            MessageID: u.Id,
            RandomID:  u.Random_id,
        }
    case TL_updateDeleteMessages:
        return UpdateDeleteMessages{
            Pts:        u.Pts,
            PtsCount:   u.Pts_count,
            MessageIDs: u.Messages,
        }
    case TL_updateContactRegistered:
        return UpdateContactRegistered{
            Date:   u.Date,
            UserID: u.User_id,
        }

    case TL_updateUserBlocked:
        var blocked bool
        switch u.Blocked.(type) {
        case TL_boolTrue:
            blocked = true
        }
        return UpdateUserBlocked{
            UserID:  u.User_id,
            Blocked: blocked,
        }
    case TL_updateChannelReadMessagesContents:
        return UpdateChannelReadMessagesContents{
            ChannelID:  u.Channel_id,
            MessageIDs: u.Messages,
        }
    default:
        return UpdateUnknown{
            Type: reflect.TypeOf(u).String(),
        }
    }
    return nil
}

func (m *MTProto) Updates_GetState() *UpdateState {
    resp := make(chan TL, 1)
    m.queueSend <- packetToSend{
        TL_updates_getState{},
        resp,
    }
    x := <-resp
    switch x.(type) {
    case TL_updates_state:
        return NewUpdateState(x)
    default:
        log.Println(fmt.Sprintf("RPC: %#v", x))
        return nil
    }
}

func (m *MTProto) Updates_GetDifference(pts, qts, date int32) *UpdateDifference {
    resp := make(chan TL, 1)
    m.queueSend <- packetToSend{
        TL_updates_getDifference{
            Flags:           1,
            Pts:             pts,
            Pts_total_limit: 200,
            Qts:             qts,
            Date:            date,
        },
        resp,
    }
    x := <-resp
    updateDifference := new(UpdateDifference)
    updateDifference.Users = make(map[int32]User)
    updateDifference.Chats = make(map[int32]Chat)
    updateDifference.Channels = make(map[int32]Channel)
    switch  u := x.(type) {
    case TL_updates_differenceEmpty:
        updateDifference.Type = UPDATE_DIFFERENCE_EMPTY
        updateDifference.IsSlice = false
        updateDifference.IntermediateState.Date = u.Date
        updateDifference.IntermediateState.Seq = u.Seq
        return updateDifference
    case TL_updates_difference:
        updateDifference.Type = UPDATE_DIFFERENCE_COMPLETE
        updateDifference.IsSlice = false
        updateDifference.IntermediateState = *NewUpdateState(u.State)
        for _, m := range u.New_messages {
            msg := NewMessage(m)
            if msg != nil {
                updateDifference.NewMessages = append(updateDifference.NewMessages, *msg)
            }

        }
        for _, ch := range u.Chats {
            switch ch.(type) {
            case TL_chatFull, TL_chat, TL_chatForbidden, TL_chatEmpty:
                newChat := NewChat(ch)
                if newChat != nil {
                    updateDifference.Chats[newChat.ID] = *newChat
                }
            case TL_channel, TL_channelForbidden, TL_channelFull:
                newChannel := NewChannel(ch)
                if newChannel != nil {
                    updateDifference.Channels[newChannel.ID] = *newChannel
                }
            }

        }
        for _, u := range u.Users {
            newUser := NewUser(u)
            if newUser != nil {
                updateDifference.Users[newUser.ID] = *newUser
            }
        }
        for _, update := range u.Other_updates {
            updateDifference.OtherUpdates = append(updateDifference.OtherUpdates, NewUpdate(update))
        }
        return updateDifference
    case TL_updates_differenceSlice:
        updateDifference.Type = UPDATE_DIFFERENCE_SLICE
        updateDifference.IsSlice = true
        updateDifference.IntermediateState = *NewUpdateState(u.Intermediate_state)
        for _, m := range u.New_messages {
            msg := NewMessage(m)
            if msg != nil {
                updateDifference.NewMessages = append(updateDifference.NewMessages, *msg)
            }

        }
        for _, ch := range u.Chats {
            switch ch.(type) {
            case TL_chatFull, TL_chat, TL_chatForbidden, TL_chatEmpty:
                newChat := NewChat(ch)
                if newChat != nil {
                    updateDifference.Chats[newChat.ID] = *newChat
                }
            case TL_channel, TL_channelForbidden, TL_channelFull:
                newChannel := NewChannel(ch)
                if newChannel != nil {
                    updateDifference.Channels[newChannel.ID] = *newChannel
                }
            }

        }
        for _, u := range u.Users {
            newUser := NewUser(u)
            if newUser != nil {
                updateDifference.Users[newUser.ID] = *newUser
            }
        }
        for _, update := range u.Other_updates {
            updateDifference.OtherUpdates = append(updateDifference.OtherUpdates, NewUpdate(update))
        }

        return updateDifference
    case TL_updates_differenceTooLong:
        updateDifference.Type = UPDATE_DIFFERENCE_TOO_LONG
        updateDifference.IntermediateState.Pts = u.Pts
        return updateDifference
    default:
        log.Println(fmt.Sprintf("RPC: %#v", x))
        return updateDifference
    }
}

func (m *MTProto) Updates_GetChannelDifference(inputChannel TL, pts, limit int32) *ChannelUpdateDifference {
    resp := make(chan TL, 1)
    m.queueSend <- packetToSend{
        TL_updates_getChannelDifference{
            Channel: inputChannel,
            Filter:  TL_channelMessagesFilterEmpty{},
            Pts:     pts,
            Limit:   limit,
        },
        resp,
    }
    x := <-resp
    updateDifference := new(ChannelUpdateDifference)
    switch u := x.(type) {
    case TL_updates_channelDifferenceEmpty:
        updateDifference.Empty = true
        updateDifference.Pts = u.Pts
        updateDifference.Flags = u.Flags
        updateDifference.Timeout = u.Timeout

    case TL_updates_channelDifference:
        updateDifference.Pts = u.Pts
        updateDifference.Flags = u.Flags
        updateDifference.Timeout = u.Timeout
        updateDifference.NewMessages = []Message{}
        updateDifference.OtherUpdates = []IUpdate{}
        for _, m := range u.New_messages {
            msg := NewMessage(m)
            if msg != nil {
                updateDifference.NewMessages = append(updateDifference.NewMessages, *msg)
            }

        }
        for _, u := range u.Other_updates {
            updateDifference.OtherUpdates = append(updateDifference.OtherUpdates, NewUpdate(u))
        }
    case TL_updates_channelDifferenceTooLong:
        updateDifference.TooLong = true
        updateDifference.Pts = u.Pts
        updateDifference.Flags = u.Flags
        updateDifference.Timeout = u.Timeout
        updateDifference.NewMessages = []Message{}
        updateDifference.OtherUpdates = []IUpdate{}
        for _, m := range u.Messages {
            msg := NewMessage(m)
            if msg != nil {
                updateDifference.NewMessages = append(updateDifference.NewMessages, *msg)
            }

        }

    case TL_rpc_error:
        log.Println("Update_GetChannelDiffrence::", u.error_code, u.error_message)
    }
    return updateDifference
}
