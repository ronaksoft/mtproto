package mtproto

import (
	"fmt"
	"reflect"
)

const (
	CHAT_TYPE_EMPTY             = "EMPTY"
	CHAT_TYPE_CHAT              = "CHAT"
	CHAT_TYPE_CHAT_FORBIDDEN    = "CHAT_FORBIDDEN"
	CHAT_TYPE_CHANNEL           = "CHANNEL"
	CHAT_TYPE_CHANNEL_FORBIDDEN = "CHANNEL_FORBIDDEN"
)

type ChatProfilePhoto struct {
	PhotoSmall FileLocation
	PhotoBig   FileLocation
}
type Chat struct {
	flags        int32
	ID           int32
	Username     string
	Type         string
	Title        string
	Photo        *ChatProfilePhoto
	Participants int32
	Members 	  []ChatMember
	Date         int32
	Left         bool
	Version      int32
	AccessHash   int64
	Address      string
	Venue        string
	CheckedIn    bool
}
type ChatMember struct {
	UserID 		int32
	InviterID	int32
	Date 		int32
}
type ChannelParticipantFilter struct {}

func (ch *Chat) GetPeer() TL {
	switch ch.Type {
	case CHAT_TYPE_CHAT, CHAT_TYPE_CHAT_FORBIDDEN:
		return TL_peerChat{
			chat_id: ch.ID,
		}
	case CHAT_TYPE_CHANNEL, CHAT_TYPE_CHANNEL_FORBIDDEN:
		return TL_peerChannel{
			channel_id: ch.ID,

		}
	default:
		return nil
	}
}
func (ch *Chat) GetInputPeer() TL {
	switch ch.Type {
	case CHAT_TYPE_CHAT, CHAT_TYPE_CHAT_FORBIDDEN:
		return TL_inputPeerChat{
			chat_id: ch.ID,
		}
	case CHAT_TYPE_CHANNEL, CHAT_TYPE_CHANNEL_FORBIDDEN:
		return TL_inputPeerChannel{
			channel_id:  ch.ID,
			access_hash: ch.AccessHash,
		}
	default:
		return nil
	}
}

func NewChatProfilePhoto(input TL) (photo *ChatProfilePhoto) {
	photo = new(ChatProfilePhoto)
	switch p := input.(type) {
	case TL_chatPhotoEmpty:
		return nil
	case TL_chatPhoto:
		switch big := p.photo_big.(type) {
		case TL_fileLocationUnavailable:
		case TL_fileLocation:
			photo.PhotoBig.DC = big.dc_id
			photo.PhotoBig.LocalID = big.local_id
			photo.PhotoBig.Secret = big.secret
			photo.PhotoBig.VolumeID = big.volume_id
		}
		switch small := p.photo_small.(type) {
		case TL_fileLocationUnavailable:
		case TL_fileLocation:
			photo.PhotoSmall.DC = small.dc_id
			photo.PhotoSmall.LocalID = small.local_id
			photo.PhotoSmall.Secret = small.secret
			photo.PhotoSmall.VolumeID = small.volume_id
		}
	}
	return photo
}
func NewChat(input TL) (chat *Chat) {
	chat = new(Chat)
	chat.Members = []ChatMember{}
	switch ch := input.(type) {
	case TL_chatEmpty:
		chat.Type = CHAT_TYPE_EMPTY
		chat.ID = ch.id
	case TL_chatForbidden:
		chat.Type = CHAT_TYPE_CHAT_FORBIDDEN
		chat.ID = ch.id
		chat.Title = ch.title
	case TL_chat:
		chat.flags = ch.flags
		chat.Type = CHAT_TYPE_CHAT
		chat.ID = ch.id
		chat.Title = ch.title
		chat.Date = ch.date
		chat.Photo = NewChatProfilePhoto(ch.photo)
		chat.Version = ch.version
		chat.Participants = ch.participants_count
	case TL_chatFull:
		chat.ID = ch.id
		participants := ch.participants.(TL_chatParticipants)
		chat.Version = participants.version
		for _, tl := range ch.participants.(TL_chatParticipants).participants {
			m := tl.(TL_chatParticipant)
			chat.Members = append(chat.Members, ChatMember{m.user_id, m.inviter_id, m.date})
		}
	case TL_channelFull:

	case TL_channelForbidden:
		chat.flags = ch.flags
		chat.Type = CHAT_TYPE_CHANNEL_FORBIDDEN
		chat.ID = ch.id
		chat.Title = ch.title
		chat.AccessHash = ch.access_hash
	case TL_channel:
		chat.flags = ch.flags
		chat.Type = CHAT_TYPE_CHANNEL
		chat.ID = ch.id
		chat.Username = ch.username
		chat.Title = ch.title
		chat.Date = ch.date
		chat.Photo = NewChatProfilePhoto(ch.photo)
		chat.Version = ch.version
		chat.AccessHash = ch.access_hash
	default:
		fmt.Println(reflect.TypeOf(ch).String())
		return nil
	}
	return chat

}
func NewInputPeerUser(userID int32, accessHash int64) TL {
	return TL_inputPeerUser{
		user_id: userID,
		access_hash: accessHash,
	}
}
func NewInputPeerChat(chatID int32) TL {
	return TL_inputPeerChat{
		chat_id: chatID,
	}
}
func NewInputPeerChannel(channelID int32, accessHash int64) TL {
	return TL_inputPeerChannel{
		channel_id:  channelID,
		access_hash: accessHash,
	}
}
func NewPeerChat(chatID int32) TL {
	return TL_peerChat{
		chat_id: chatID,
	}
}
func NewPeerChannel(channelID int32) TL {
	return TL_peerChannel{
		channel_id: channelID,
	}
}
