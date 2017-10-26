package mtproto

import (
	"reflect"
	"log"
	"fmt"
)

const (
	//UPDATE_TYPE_NEW_MESSAGE             string = "NewMessage"
	//UPDATE_TYPE_CHANNEL_NEW_MESSAGE     string = "ChannelNewMessage"
	//UPDATE_TYPE_READ_CHANNEL_INBOX      string = "ReadChannelInbox"
	//UPDATE_TYPE_CHANNEL_TOO_LONG        string = "ChannelTooLong"
	//UPDATE_TYPE_READ_HISTORY_INBOX      string = "ReadHistoryInbox"
	//UPDATE_TYPE_USER_PHOTO              string = "UserPhoto"
	//UPDATE_TYPE_USER_TYPING             string = "UserTyping"
	//UPDATE_TYPE_CHAT_PARTICIPANT_ADD    string = "ChatParticipantAdd"
	//UPDATE_TYPE_CHAT_PARTICIPANT_ADMIN  string = "ChatParticipantAdmin"
	//UPDATE_TYPE_CHAT_PARTICIPANT_DELETE string = "ChatParticipantDelete"
	//UPDATE_TYPE_CHAT_USER_TYPING        string = "ChatUserTyping"
)

const (
	UPDATE_DIFFERENCE_EMPTY    = "EMPTY"
	UPDATE_DIFFERENCE_SLICE    = "SLICE"
	UPDATE_DIFFERENCE_TOO_LONG = "TOO_LONG"
)

type Update struct {
	Type      string
	UserID    int32
	InviterID int32
	ChatID    int32
	Pts       int32
	PtsCount  int32
	Message   *Message
	Version   int32
	Date      int32
	ChannelID int32
	MaxID     int32
	Flags     int32
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
	OtherUpdates      []Update
	Chats             []Chat
	Channels          []Channel
	Users             []User
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
	OtherUpdates []Update
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
func NewUpdate(input TL) *Update {
	update := new(Update)
	update.Type = reflect.TypeOf(input).String()
	switch u := input.(type) {
	case TL_updateNewMessage:
		update.Pts = u.Pts
		update.PtsCount = u.Pts_count
		update.Message = NewMessage(u.Message)
	case TL_updateNewChannelMessage:
		update.Message = NewMessage(u.Message)
		update.PtsCount = u.Pts_count
		update.Pts = u.Pts
	case TL_updateReadChannelInbox:
		update.ChannelID = u.Channel_id
		update.MaxID = u.Max_id
	case TL_updateChannelTooLong:
		update.Pts = u.Pts
		update.ChannelID = u.Channel_id
		update.Flags = u.Flags
	case TL_updateReadHistoryInbox:
		update.Pts = u.Pts
		update.PtsCount = u.Pts_count
		update.MaxID = u.Max_id
	case TL_updateUserPhoto:
		update.UserID = u.User_id
		update.Date = u.Date
		// Save NewUserProfilePhoto(u.Photo)
	case TL_updateContactLink:
		update.UserID = u.User_id
	case TL_updateEditChannelMessage:
		update.Pts = u.Pts
		update.PtsCount = u.Pts_count
		update.Message = NewMessage(u.Message)
	case TL_updateEditMessage:
		update.Pts = u.Pts
		update.PtsCount = u.Pts_count
		update.Message = NewMessage(u.Message)
	default:
		update.Type = reflect.TypeOf(u).String()
	}
	return update
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
			Pts_total_limit: 100,
			Qts:             qts,
			Date:            date,
		},
		resp,
	}
	x := <-resp
	updateDifference := new(UpdateDifference)
	switch  u := x.(type) {
	case TL_updates_differenceEmpty:
		updateDifference.Type = UPDATE_DIFFERENCE_EMPTY
		updateDifference.IsSlice = false
		updateDifference.IntermediateState.Date = u.Date
		updateDifference.IntermediateState.Seq = u.Seq
		return updateDifference
	case TL_updates_difference:
		updateDifference.IsSlice = false
		updateDifference.IntermediateState = *NewUpdateState(u.State)
		for _, m := range u.New_messages {
			updateDifference.NewMessages = append(updateDifference.NewMessages, *NewMessage(m))
		}
		for _, ch := range u.Chats {
			switch ch.(type) {
			case TL_chatFull, TL_chat, TL_chatForbidden, TL_chatEmpty:
				updateDifference.Chats = append(updateDifference.Chats, *NewChat(ch))
			case TL_channel, TL_channelForbidden, TL_channelFull:
				updateDifference.Channels = append(updateDifference.Channels, *NewChannel(ch))
			}

		}
		for _, user := range u.Users {
			updateDifference.Users = append(updateDifference.Users, *NewUser(user))
		}
		for _, update := range u.Other_updates {
			updateDifference.OtherUpdates = append(updateDifference.OtherUpdates, *NewUpdate(update))
		}
		return updateDifference
	case TL_updates_differenceSlice:
		updateDifference.Type = UPDATE_DIFFERENCE_SLICE
		updateDifference.IsSlice = true
		updateDifference.IntermediateState = *NewUpdateState(u.Intermediate_state)
		for _, m := range u.New_messages {
			updateDifference.NewMessages = append(updateDifference.NewMessages, *NewMessage(m))
		}
		for _, ch := range u.Chats {
			switch ch.(type) {
			case TL_chatFull, TL_chat, TL_chatForbidden, TL_chatEmpty:
				updateDifference.Chats = append(updateDifference.Chats, *NewChat(ch))
			case TL_channel, TL_channelForbidden, TL_channelFull:
				updateDifference.Channels = append(updateDifference.Channels, *NewChannel(ch))
			}

		}
		for _, user := range u.Users {
			updateDifference.Users = append(updateDifference.Users, *NewUser(user))
		}
		for _, update := range u.Other_updates {
			updateDifference.OtherUpdates = append(updateDifference.OtherUpdates, *NewUpdate(update))
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
		updateDifference.OtherUpdates = []Update{}
		for _, m := range u.New_messages {
			updateDifference.NewMessages = append(updateDifference.NewMessages, *NewMessage(m))
		}
		for _, u := range u.Other_updates {
			updateDifference.OtherUpdates = append(updateDifference.OtherUpdates, *NewUpdate(u))
		}
	case TL_updates_channelDifferenceTooLong:
		updateDifference.TooLong = true
		updateDifference.Pts = u.Pts
		updateDifference.Flags = u.Flags
		updateDifference.Timeout = u.Timeout
		updateDifference.NewMessages = []Message{}
		updateDifference.OtherUpdates = []Update{}
		for _, m := range u.Messages {
			updateDifference.NewMessages = append(updateDifference.NewMessages, *NewMessage(m))
		}
	case TL_rpc_error:
		log.Println("Update_GetChannelDiffrence::", u.error_code, u.error_message)
	}
	return updateDifference
}
