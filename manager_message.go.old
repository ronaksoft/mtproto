package mtproto

import (
	"fmt"
	"reflect"
	"log"
)

type MessageAction struct {
	Type      string
	Title     string
	ChatID    int32
	ChannelID int32
	GameID    int64
	GameScore int32
	UserID    int32
	UserIDs   []int32
	Photo     *Photo
}
type Message struct {
	flags         int32
	Type          string
	ID            int32
	From          int32
	To            *Peer
	Date          int32
	Body          string
	MediaType     string
	Action        *MessageAction
	ForwardHeader *MessageForwardHeader
	Views         int32
	Media         interface{}
}
type MessageForwardHeader struct {
	From        int32
	Date        int32
	ChannelID   int32
	ChannelPost int32
}
type MessageMedia interface{}
type MessageMediaPhoto struct {
	Caption string
	Photo   Photo
}
type MessageMediaContact struct {
	Firstname string
	Lastname  string
	UserID    int32
	Phone     string
}
type MessageMediaDocument struct {
	Caption  string
	Document Document
}

func NewMessageMedia(in TL) (interface{}) {
	switch x := in.(type) {
	case TL_messageMediaPhoto:
		mm := new(MessageMediaPhoto)
		mm.Caption = x.caption
		mm.Photo = *NewPhoto(x.photo)
		return mm
	case TL_messageMediaContact:
		mm := new(MessageMediaContact)
		mm.UserID = x.user_id
		mm.Firstname = x.first_name
		mm.Lastname = x.last_name
		mm.Phone = x.phone_number
		return mm
	case TL_messageMediaDocument:
		mm := new(MessageMediaDocument)
		mm.Caption = x.caption
		mm.Document = *NewDocument(x.document)
		return mm
	default:
		log.Println("NewMessageMedia::", reflect.TypeOf(x).String())
	}
	return nil
}

func NewMessageForwardHeader(in TL) (fwd *MessageForwardHeader) {
	fwd = new(MessageForwardHeader)
	fwdHeader := in.(TL_messageFwdHeader)
	fwd.Date = fwdHeader.date
	fwd.From = fwdHeader.from_id
	fwd.ChannelID = fwdHeader.channel_id
	fwd.ChannelPost = fwdHeader.channel_post
	return
}
func NewMessageAction(in TL) (m *MessageAction) {
	m = new(MessageAction)
	switch x := in.(type) {
	case TL_messageActionEmpty:
	case TL_messageActionChannelCreate:
		m.Type = MESSAGE_ACTION_CHANNEL_CREATED
		m.Title = x.title
	case TL_messageActionChannelMigrateFrom:
		m.Type = MESSAGE_ACTION_CHANNEL_MIGRATE_FROM
		m.Title = x.title
		m.ChatID = x.chat_id
	case TL_messageActionChatCreate:
		m.Type = MESSAGE_ACTION_CHAT_CREATED
		m.Title = x.title
		m.UserIDs = x.users
	case TL_messageActionChatAddUser:
		m.Type = MESSAGE_ACTION_CHAT_ADD_USER
		m.UserIDs = x.users
	case TL_messageActionChatDeleteUser:
		m.Type = MESSAGE_ACTION_CHAT_DELETE_USER
		m.UserID = x.user_id
	case TL_messageActionChatDeletePhoto:
		m.Type = MESSAGE_ACTION_CHAT_DELETE_PHOTO
	case TL_messageActionChatEditPhoto:
		m.Type = MESSAGE_ACTION_CHAT_EDIT_PHOTO
		m.Photo = NewPhoto(x.photo)
	case TL_messageActionChatEditTitle:
		m.Type = MESSAGE_ACTION_CHAT_EDIT_TITLE
		m.Title = x.title
	case TL_messageActionChatJoinedByLink:
		m.Type = MESSAGE_ACTION_CHAT_JOINED_BY_LINK
		m.UserID = x.inviter_id
	case TL_messageActionChatMigrateTo:
		m.Type = MESSAGE_ACTION_CHAT_MIGRATE_TO
		m.ChannelID = x.channel_id
	case TL_messageActionGameScore:
		m.Type = MESSAGE_ACTION_GAME_SCORE
		m.GameID = x.game_id
		m.GameScore = x.score
	case TL_messageActionHistoryClear:
		m.Type = MESSAGE_ACTION_HISTORY_CLEAN
	case TL_messageActionPinMessage:
	default:
		return nil
	}
	return
}
func NewMessage(in TL) (m *Message) {
	m = new(Message)
	switch x := in.(type) {
	case TL_message:
		m.flags = x.flags
		m.Type = MESSAGE_TYPE_NORMAL
		m.ID = x.id
		m.Date = x.date
		m.From = x.from_id
		m.Body = x.message
		m.To = NewPeer(x.to_id)
		m.Views = x.views
		if x.media != nil {
			m.Media = NewMessageMedia(x.media)
		}
		if x.fwd_from != nil {
			m.ForwardHeader = NewMessageForwardHeader(x.fwd_from)
		}
	case TL_messageService:
		m.flags = x.flags
		m.Type = MESSAGE_TYPE_SERVICE
		m.ID = x.id
		m.Date = x.date
		m.From = x.from_id
		m.To = NewPeer(x.to_id)
		m.Action = NewMessageAction(x.action)
		m.ForwardHeader = new(MessageForwardHeader)
	default:
		fmt.Println("GER", reflect.TypeOf(x).String())
	}
	return
}
