package mtproto

import (
	"reflect"
	"log"
	"fmt"
)

const (
	UPDATE_TYPE_NEW_MESSAGE             string = "NewMessage"
	UPDATE_TYPE_CHANNEL_NEW_MESSAGE     string = "ChannelNewMessage"
	UPDATE_TYPE_USER_TYPING             string = "UserTyping"
	UPDATE_TYPE_CHAT_PARTICIPANT_ADD    string = "ChatParticipantAdd"
	UPDATE_TYPE_CHAT_PARTICIPANT_ADMIN  string = "ChatParticipantAdmin"
	UPDATE_TYPE_CHAT_PARTICIPANT_DELETE string = "ChatParticipantDelete"
	UPDATE_TYPE_CHAT_USER_TYPING        string = "ChatUserTyping"
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
}
type UpdateState struct {
	Qts          int32
	Pts          int32
	Date         int32
	Seq          int32
	UnreadCounts int32
}
type UpdateDifference struct {
	IsSlice           bool
	Total             int32
	NewMessages       []Message
	OtherUpdates      []Update
	Chats             []Chat
	Users             []User
	IntermediateState UpdateState
}
type ChannelUpdateDifference struct {
	Empty       bool
	TooLong     bool
	Flags       int32
	Final       bool
	Pts         int32
	Timeout     int32
	NewMessages []Message
}

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

// input :
//	1. TL_updateNewMessage
//	2. TL_updateNewChannelMessage
func NewUpdate(input TL) *Update {
	update := new(Update)
	switch u := input.(type) {
	case TL_updateNewMessage:
		update.Type = UPDATE_TYPE_NEW_MESSAGE
		update.Pts = u.Pts
		update.PtsCount = u.Pts_count
		update.Message = NewMessage(u.Message)
	case TL_updateNewChannelMessage:
		update.Type = UPDATE_TYPE_CHANNEL_NEW_MESSAGE
		update.Message = NewMessage(u.Message)
		update.PtsCount = u.Pts_count
		update.Pts = u.Pts
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
		updateDifference.IntermediateState.Date = 0
		return updateDifference
	case TL_updates_difference:
		updateDifference.IsSlice = false
		updateDifference.IntermediateState = *NewUpdateState(u.State)
		for _, m := range u.New_messages {
			updateDifference.NewMessages = append(updateDifference.NewMessages, *NewMessage(m))
		}
		for _, ch := range u.Chats {
			updateDifference.Chats = append(updateDifference.Chats, *NewChat(ch))
		}
		for _, user := range u.Users {
			updateDifference.Users = append(updateDifference.Users, *NewUser(user))
		}
		for _, update := range u.Other_updates {
			updateDifference.OtherUpdates = append(updateDifference.OtherUpdates, *NewUpdate(update))
		}
		return updateDifference
	case TL_updates_differenceSlice:
		updateDifference.IsSlice = true
		updateDifference.IntermediateState = *NewUpdateState(u.Intermediate_state)
		for _, m := range u.New_messages {
			updateDifference.NewMessages = append(updateDifference.NewMessages, *NewMessage(m))
		}
		for _, ch := range u.Chats {
			updateDifference.Chats = append(updateDifference.Chats, *NewChat(ch))
		}
		for _, user := range u.Users {
			updateDifference.Users = append(updateDifference.Users, *NewUser(user))
		}
		for _, update := range u.Other_updates {
			updateDifference.OtherUpdates = append(updateDifference.OtherUpdates, *NewUpdate(update))
		}

		return updateDifference
	case TL_updates_differenceTooLong:
		updateDifference.IntermediateState.Pts = u.Pts
		return updateDifference
	default:
		log.Println(fmt.Sprintf("RPC: %#v", x))
		return updateDifference
	}
}

func (m *MTProto) Updates_GetChannelDifference(inputChannel TL) *ChannelUpdateDifference {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_updates_getChannelDifference{
			Channel: inputChannel,
			Filter:  TL_channelMessagesFilterEmpty{},
		},
		resp,
	}
	x := <-resp
	updateDifference := new(ChannelUpdateDifference)
	switch u := x.(type) {
	case TL_updates_channelDifferenceEmpty:
		updateDifference.Empty = true
	case TL_updates_channelDifference:
		updateDifference.Pts = u.Pts
		updateDifference.Flags = u.Flags
		updateDifference.NewMessages = []Message{}
		for _, m := range u.New_messages {
			updateDifference.NewMessages = append(updateDifference.NewMessages, *NewMessage(m))
		}
	case TL_updates_channelDifferenceTooLong:
		updateDifference.TooLong = true
	}
	return updateDifference
}
