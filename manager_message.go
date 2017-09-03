package mtproto

import (
	"fmt"
	"reflect"
)

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
	Entities      []MessageEntity
	Views         int32
	Media         interface{}
}
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
type MessageEntity struct {
	Type     string
	Offset   int32
	Length   int32
	Url      string
	language string
	UserID   int32
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

// input
//	1. TL_message
//	2. TL_messageService
func NewMessage(input TL) (m *Message) {
	m = new(Message)
	switch x := input.(type) {
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
		m.Entities = make([]MessageEntity, 0, 0)
		for _, e := range x.Entities {
			m.Entities = append(m.Entities, *NewMessageEntity(e))
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
		fmt.Println("NewMessage::UnSupported Input Format", reflect.TypeOf(x).String())
	}
	return
}

// input:
//	1. TL_messageActionEmpty
//	2. TL_messageActionChannelCreate
//	3. TL_messageActionChannelMigrateFrom
//	4. TL_messageActionChatCreate
//	5. TL_messageActionChatAddUser
//	6. TL_messageActionChatDeleteUser
//	7. TL_messageActionChatDeleteUser
//	8. TL_messageActionChatEditPhoto
//	9. TL_messageActionChatEditTitle
//	10. TL_messageActionChatJoinedByLink
//	11.	TL_messageActionChatMigrateTo
//	12.	TL_messageActionGameScore
//	13. TL_messageActionHistoryClear
//	14. TL_messageActionPinMessage
//	15. TL_messageActionPhoneCall
func NewMessageAction(input TL) (m *MessageAction) {
	m = new(MessageAction)
	switch x := input.(type) {
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
	case TL_messageActionPhoneCall:
		m.Type = MESSAGE_ACTION_PHONE_CALL
	default:
		fmt.Println("NewMessageAction::UnSupported Input Format", reflect.TypeOf(x).String())
	}
	return
}

func NewMessageEntity(input TL) (e *MessageEntity) {
	e = new(MessageEntity)
	switch x := input.(type) {
	case TL_messageEntityBold:
		e.Type = MESSAGE_ENTITY_BOLD
		e.Offset, e.Length = x.Offset, x.Length
	case TL_messageEntityEmail:
		e.Type = MESSAGE_ENTITY_EMAIL
		e.Offset, e.Length = x.Offset, x.Length
	case TL_messageEntityBotCommand:
		e.Type = MESSAGE_ENTITY_BOT_COMMAND
		e.Offset, e.Length = x.Offset, x.Length
	case TL_messageEntityHashtag:
		e.Type = MESSAGE_ENTITY_HASHTAG
		e.Offset, e.Length = x.Offset, x.Length
	case TL_messageEntityCode:
		e.Type = MESSAGE_ENTITY_CODE
		e.Offset, e.Length = x.Offset, x.Length
	case TL_messageEntityItalic:
		e.Type = MESSAGE_ENTITY_ITALIC
		e.Offset, e.Length = x.Offset, x.Length
	case TL_messageEntityMention:
		e.Type = MESSAGE_ENTITY_MENTION
		e.Offset, e.Length = x.Offset, x.Length
	case TL_messageEntityUrl:
		e.Type = MESSAGE_ENTITY_URL
		e.Offset, e.Length = x.Offset, x.Length
	case TL_messageEntityTextUrl:
		e.Type = MESSAGE_ENTITY_TEXT_URL
		e.Offset, e.Length = x.Offset, x.Length
		e.Url = x.Url
	case TL_messageEntityPre:
		e.Type = MESSAGE_ENTITY_PRE
		e.Offset, e.Length = x.Offset, x.Length
		e.language = x.Language
	case TL_messageEntityMentionName:
		e.Type = MESSAGE_ENTITY_MENTION_NAME
		e.Offset, e.Length = x.Offset, x.Length
		e.UserID = x.User_id
	default:
		fmt.Println("NewMessageEntity::UnSupported Input Format", reflect.TypeOf(x).String())
	}
	return e
}

func NewMessageForwardHeader(input TL) (fwd *MessageForwardHeader) {
	fwd = new(MessageForwardHeader)
	fwdHeader := input.(TL_messageFwdHeader)
	fwd.Date = fwdHeader.Date
	fwd.From = fwdHeader.From_id
	fwd.ChannelID = fwdHeader.Channel_id
	fwd.ChannelPost = fwdHeader.Channel_post
	return
}

// input:
//	1. TL_messageMediaPhoto
//	2. TL_messageMediaContact
//	3. TL_messageMediaDocument
func NewMessageMedia(input TL) (interface{}) {
	switch x := input.(type) {
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
		fmt.Println("NewMessageMedia::UnSupported Input Format", reflect.TypeOf(x).String())
	}
	return nil
}
