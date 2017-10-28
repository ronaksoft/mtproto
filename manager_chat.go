package mtproto

import (
	"fmt"
	"reflect"
)

const (
	CHAT_TYPE_EMPTY             = "EMPTY"
	CHAT_TYPE_CHAT              = "CHAT"
	CHAT_TYPE_CHAT_FORBIDDEN    = "CHAT_FORBIDDEN"
)

type ChatProfilePhoto struct {
	PhotoSmall FileLocation
	PhotoBig   FileLocation
}
type Chat struct {
	flags        int32
	Type         string
	ID           int32
	Username     string
	Title        string
	Photo        *ChatProfilePhoto
	Participants int32
	Members      []ChatMember
	Date         int32
	Left         bool
	Version      int32
	AccessHash   int64
	Address      string
	Venue        string
	CheckedIn    bool
}
type ChatMember struct {
	UserID    int32
	InviterID int32
	Date      int32
}
type ChannelParticipantFilter struct{}

func (ch *Chat) GetPeer() TL {
	switch ch.Type {
	case CHAT_TYPE_CHAT, CHAT_TYPE_CHAT_FORBIDDEN:
		return TL_peerChat{
			Chat_id: ch.ID,
		}
	default:
		return nil
	}
}
func (ch *Chat) GetInputPeer() TL {
	switch ch.Type {
	case CHAT_TYPE_CHAT, CHAT_TYPE_CHAT_FORBIDDEN:
		return TL_inputPeerChat{
			Chat_id: ch.ID,
		}
	default:
		return nil
	}
}

// NewChatProfilePhoto
// input :
//	1. TL_chatPhotoEmpty
//	2. TL_chatPhoto
func NewChatProfilePhoto(input TL) (photo *ChatProfilePhoto) {
	photo = new(ChatProfilePhoto)
	switch p := input.(type) {
	case TL_chatPhotoEmpty:
		return nil
	case TL_chatPhoto:
		switch big := p.Photo_big.(type) {
		case TL_fileLocationUnavailable:
		case TL_fileLocation:
			photo.PhotoBig.DC = big.Dc_id
			photo.PhotoBig.LocalID = big.Local_id
			photo.PhotoBig.Secret = big.Secret
			photo.PhotoBig.VolumeID = big.Volume_id
		}
		switch small := p.Photo_small.(type) {
		case TL_fileLocationUnavailable:
		case TL_fileLocation:
			photo.PhotoSmall.DC = small.Dc_id
			photo.PhotoSmall.LocalID = small.Local_id
			photo.PhotoSmall.Secret = small.Secret
			photo.PhotoSmall.VolumeID = small.Volume_id
		}
	}
	return photo
}

// NewChat
// input:
//	1. TL_chatEmpty
//	2. TL_chatForbidden
//	3. TL_chat
//	4. TL_chatFull:
func NewChat(input TL) (chat *Chat) {
	chat = new(Chat)
	chat.Members = []ChatMember{}
	switch ch := input.(type) {
	case TL_chatEmpty:
		chat.Type = CHAT_TYPE_EMPTY
		chat.ID = ch.Id
	case TL_chatForbidden:
		chat.Type = CHAT_TYPE_CHAT_FORBIDDEN
		chat.ID = ch.Id
		chat.Title = ch.Title
	case TL_chat:
		chat.flags = ch.Flags
		chat.Type = CHAT_TYPE_CHAT
		chat.ID = ch.Id
		chat.Title = ch.Title
		chat.Date = ch.Date
		chat.Photo = NewChatProfilePhoto(ch.Photo)
		chat.Version = ch.Version
		chat.Participants = ch.Participants_count
	case TL_chatFull:
		chat.ID = ch.Id
		participants := ch.Participants.(TL_chatParticipants)
		chat.Version = participants.Version
		for _, tl := range ch.Participants.(TL_chatParticipants).Participants {
			m := tl.(TL_chatParticipant)
			chat.Members = append(chat.Members, ChatMember{m.User_id, m.Inviter_id, m.Date})
		}
	default:
		fmt.Println(reflect.TypeOf(ch).String())
		return nil
	}
	return chat
}






