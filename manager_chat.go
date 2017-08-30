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
	case CHAT_TYPE_CHANNEL, CHAT_TYPE_CHANNEL_FORBIDDEN:
		return TL_peerChannel{
			Channel_id: ch.ID,
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
	case CHAT_TYPE_CHANNEL, CHAT_TYPE_CHANNEL_FORBIDDEN:
		return TL_inputPeerChannel{
			Channel_id:  ch.ID,
			Access_hash: ch.AccessHash,
		}
	default:
		return nil
	}
}

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

// input:
//	1. TL_chatEmpty
//	2. TL_chatForbidden
//	3. TL_chat
//	4. TL_chatFull:
//	5. TL_channelFull:
//	6. TL_channelForbidden:
//	7. TL_channel
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
	case TL_channelFull:
	case TL_channelForbidden:
		chat.flags = ch.Flags
		chat.Type = CHAT_TYPE_CHANNEL_FORBIDDEN
		chat.ID = ch.Id
		chat.Title = ch.Title
		chat.AccessHash = ch.Access_hash
	case TL_channel:
		chat.flags = ch.Flags
		chat.Type = CHAT_TYPE_CHANNEL
		chat.ID = ch.Id
		chat.Username = ch.Username
		chat.Title = ch.Title
		chat.Date = ch.Date
		chat.Photo = NewChatProfilePhoto(ch.Photo)
		chat.Version = ch.Version
		chat.AccessHash = ch.Access_hash
	default:
		fmt.Println(reflect.TypeOf(ch).String())
		return nil
	}
	return chat
}

func (m *MTProto) Channels_GetParticipants(channel TL, offset, limit int32) []User {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_channels_getParticipants{
			Channel: channel,
			Filter:  TL_channelParticipantsRecent{},
			Offset:  offset,
			Limit:   limit,
		},
		resp,
	}
	x := <-resp
	users := make([]User, 0)
	switch input := x.(type) {
	case TL_channels_channelParticipants:
		for _, u := range input.Users {
			users = append(users, *NewUser(u))
		}
	case TL_rpc_error:
		fmt.Println(input.error_code, input.error_message)
	default:
		fmt.Println(reflect.TypeOf(input).String())
	}
	return users
}

func (m *MTProto) Channels_GetChannels(in []TL) []Chat {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_channels_getChannels{
			Id: in,
		},
		resp,
	}
	x := <-resp
	chats := make([]Chat, 0, len(in))
	switch input := x.(type) {
	case TL_messages_chats:
		for _, ch := range input.Chats {
			chats = append(chats, *NewChat(ch))
		}
		return chats
	case TL_rpc_error:
		fmt.Println(input.error_code, input.error_message)
		return chats
	default:
		fmt.Println(reflect.TypeOf(input).String())
		return chats
	}
}

func (m *MTProto) Channels_GetMessages(channel TL, ids []int32) []Message {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_channels_getMessages{
			Channel: channel,
			Id:      ids,
		},
		resp,
	}
	x := <-resp
	messages := make([]Message, 0, len(ids))
	switch input := x.(type) {
	case TL_messages_messages:
		for _, m := range input.Messages {
			messages = append(messages, *NewMessage(m))
		}
		return messages
	case TL_messages_messagesSlice:
		for _, m := range input.Messages {
			messages = append(messages, *NewMessage(m))
		}
		return messages
	case TL_messages_channelMessages:
		for _, m := range input.Messages {
			messages = append(messages, *NewMessage(m))
		}
		return messages
	case TL_rpc_error:
		fmt.Println(input.error_code, input.error_message)
		return messages
	default:
		fmt.Println(reflect.TypeOf(input).String())
		return messages
	}

}

func (m *MTProto) Messages_GetChats(chatIDs []int32) []Chat {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_messages_getChats{
			Id: chatIDs,
		},
		resp,
	}
	x := <-resp
	chats := make([]Chat, 0, len(chatIDs))
	switch input := x.(type) {
	case TL_messages_chats:
		for _, ch := range input.Chats {
			chats = append(chats, *NewChat(ch))
		}
		return chats
	case TL_rpc_error:
		fmt.Println(input.error_code, input.error_message)
		return chats
	default:
		fmt.Println(reflect.TypeOf(input).String())
		return chats
	}
}

func (m *MTProto) Messages_GetFullChat(chatID int32) *Chat {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_messages_getFullChat{
			Chat_id: chatID,
		},
		resp,
	}
	x := <-resp
	chat := new(Chat)
	switch input := x.(type) {
	case TL_messages_chatFull:
		chat = NewChat(input)
	default:

	}
	return chat
}

func (m *MTProto) Messages_GetHistory(inputPeer TL, limit int32) ([]Message, int32) {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_messages_getHistory{
			Peer:  inputPeer,
			Limit: limit,
		},
		resp,
	}
	x := <-resp
	messages := make([]Message, 0, 20)
	switch input := x.(type) {
	case TL_messages_messages:
		for _, msg := range input.Messages {
			messages = append(messages, *NewMessage(msg))
		}
		return messages, int32(len(messages))
	case TL_messages_messagesSlice:
		for _, msg := range input.Messages {
			messages = append(messages, *NewMessage(msg))
		}
		return messages, input.Count
	case TL_messages_channelMessages:
		for _, msg := range input.Messages {
			messages = append(messages, *NewMessage(msg))
		}
		return messages, input.Count
	case TL_rpc_error:
		fmt.Println(input.error_message, input.error_code)
		return messages, 0
	default:
		fmt.Println(reflect.TypeOf(input).String())
		return messages, 0
	}

}

