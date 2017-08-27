package mtproto

import "reflect"

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
	IsSlice 			bool
	Total             int32
	NewMessages       []Message
	OtherUpdates      []Update
	Chats             []Chat
	Users             []User
	IntermediateState UpdateState
}

func NewUpdateState(input TL) *UpdateState {
	us := new(UpdateState)
	switch in := input.(type) {
	case TL_updates_state:
		us.Qts = in.qts
		us.Pts = in.pts
		us.Seq = in.seq
		us.Date = in.date
		us.UnreadCounts = in.unread_count
	}
	return us
}

func NewUpdate(input TL) *Update {
	update := new(Update)
	switch u := input.(type) {
	case TL_updateNewMessage:
		update.Type = UPDATE_TYPE_NEW_MESSAGE
		update.Pts = u.pts
		update.PtsCount = u.pts_count
		update.Message = NewMessage(u.message)
	case TL_updateNewChannelMessage:
		update.Type = UPDATE_TYPE_CHANNEL_NEW_MESSAGE
		update.Message = NewMessage(u.message)
		update.PtsCount = u.pts_count
		update.Pts = u.pts
	default:
		update.Type = reflect.TypeOf(u).String()
	}
	return update
}
