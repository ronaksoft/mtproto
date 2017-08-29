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
		mm.Caption = x.Caption
		mm.Photo = *NewPhoto(x.Photo)
		return mm
	case TL_messageMediaContact:
		mm := new(MessageMediaContact)
		mm.UserID = x.User_id
		mm.Firstname = x.First_name
		mm.Lastname = x.Last_name
		mm.Phone = x.Phone_number
		return mm
	case TL_messageMediaDocument:
		mm := new(MessageMediaDocument)
		mm.Caption = x.Caption
		mm.Document = *NewDocument(x.Document)
		return mm
	default:
		log.Println("NewMessageMedia::", reflect.TypeOf(x).String())
	}
	return nil
}

func NewMessageForwardHeader(in TL) (fwd *MessageForwardHeader) {
	fwd = new(MessageForwardHeader)
	fwdHeader := in.(TL_messageFwdHeader)
	fwd.Date = fwdHeader.Date
	fwd.From = fwdHeader.From_id
	fwd.ChannelID = fwdHeader.Channel_id
	fwd.ChannelPost = fwdHeader.Channel_post
	return
}
func NewMessageAction(in TL) (m *MessageAction) {
	m = new(MessageAction)
	switch x := in.(type) {
	case TL_messageActionEmpty:
	case TL_messageActionChannelCreate:
		m.Type = MESSAGE_ACTION_CHANNEL_CREATED
		m.Title = x.Title
	case TL_messageActionChannelMigrateFrom:
		m.Type = MESSAGE_ACTION_CHANNEL_MIGRATE_FROM
		m.Title = x.Title
		m.ChatID = x.Chat_id
	case TL_messageActionChatCreate:
		m.Type = MESSAGE_ACTION_CHAT_CREATED
		m.Title = x.Title
		m.UserIDs = x.Users
	case TL_messageActionChatAddUser:
		m.Type = MESSAGE_ACTION_CHAT_ADD_USER
		m.UserIDs = x.Users
	case TL_messageActionChatDeleteUser:
		m.Type = MESSAGE_ACTION_CHAT_DELETE_USER
		m.UserID = x.User_id
	case TL_messageActionChatDeletePhoto:
		m.Type = MESSAGE_ACTION_CHAT_DELETE_PHOTO
	case TL_messageActionChatEditPhoto:
		m.Type = MESSAGE_ACTION_CHAT_EDIT_PHOTO
		m.Photo = NewPhoto(x.Photo)
	case TL_messageActionChatEditTitle:
		m.Type = MESSAGE_ACTION_CHAT_EDIT_TITLE
		m.Title = x.Title
	case TL_messageActionChatJoinedByLink:
		m.Type = MESSAGE_ACTION_CHAT_JOINED_BY_LINK
		m.UserID = x.Inviter_id
	case TL_messageActionChatMigrateTo:
		m.Type = MESSAGE_ACTION_CHAT_MIGRATE_TO
		m.ChannelID = x.Channel_id
	case TL_messageActionGameScore:
		m.Type = MESSAGE_ACTION_GAME_SCORE
		m.GameID = x.Game_id
		m.GameScore = x.Score
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
		m.flags = x.Flags
		m.Type = MESSAGE_TYPE_NORMAL
		m.ID = x.Id
		m.Date = x.Date
		m.From = x.From_id
		m.Body = x.Message
		m.To = NewPeer(x.To_id)
		m.Views = x.Views
		if x.Media != nil {
			m.Media = NewMessageMedia(x.Media)
		}
		if x.Fwd_from != nil {
			m.ForwardHeader = NewMessageForwardHeader(x.Fwd_from)
		}
	case TL_messageService:
		m.flags = x.Flags
		m.Type = MESSAGE_TYPE_SERVICE
		m.ID = x.Id
		m.Date = x.Date
		m.From = x.From_id
		m.To = NewPeer(x.To_id)
		m.Action = NewMessageAction(x.Action)
		m.ForwardHeader = new(MessageForwardHeader)
	default:
		fmt.Println("GER", reflect.TypeOf(x).String())
	}
	return
}
