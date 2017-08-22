package mtproto

import (
	"fmt"
	"reflect"
	"log"
	"github.com/pkg/errors"
)

func (m *MTProto) Auth_SendCode(phonenumber string) (string, error) {
	var authSentCode TL_auth_sentCode
	flag := true
	for flag {
		resp := make(chan TL, 1)
		m.queueSend <- packetToSend{TL_auth_sendCode{
			flags:          1,
			current_number: TL_boolTrue{},
			phone_number:   phonenumber,
			api_id:         appId,
			api_hash:       appHash,
		}, resp}
		x := <-resp
		switch x.(type) {
		case TL_auth_sentCode:
			authSentCode = x.(TL_auth_sentCode)
			flag = false
		case TL_rpc_error:
			x := x.(TL_rpc_error)
			if x.error_code != 303 {
				return "", fmt.Errorf("RPC error_code: %d", x.error_code)
			}
			var newDc int32
			n, _ := fmt.Sscanf(x.error_message, "PHONE_MIGRATE_%d", &newDc)
			if n != 1 {
				n, _ := fmt.Sscanf(x.error_message, "NETWORK_MIGRATE_%d", &newDc)
				if n != 1 {
					return "", fmt.Errorf("RPC error_string: %s", x.error_message)
				}
			}

			newDcAddr, ok := m.dclist[newDc]
			if !ok {
				return "", fmt.Errorf("Wrong DC index: %d", newDc)
			}
			err := m.reconnect(newDcAddr)
			fmt.Println("Reconnected")
			if err != nil {
				return "", err
			}
		default:
			return "", fmt.Errorf("Got: %T", x)
		}

	}

	if authSentCode.flags&1 == 0 {
		return "", errors.New("Cannot sign up yet")
	}

	return authSentCode.phone_code_hash, nil
}

func (m *MTProto) Auth_SignIn(phonenumber string, hash, code string) error {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_auth_signIn{phonenumber, hash, code},
		resp,
	}
	x := <-resp
	auth, ok := x.(TL_auth_authorization)
	if !ok {
		return fmt.Errorf("RPC: %#v", x)
	}
	userSelf := auth.user.(TL_user)
	fmt.Printf("Signed in: id %d name <%s %s>\n", userSelf.id, userSelf.first_name, userSelf.last_name)
	return nil
}

func (m *MTProto) Auth_CheckPhone(phonenumber string) bool {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_auth_checkPhone{
			"989121228718",
		},
		resp,
	}
	x := <-resp
	if v, ok := x.(TL_auth_checkedPhone); ok {
		if toBool(v) {
			return true
		}
	}
	return false
}

func (m *MTProto) Contacts_GetContacts(hash string) ([]Contact, []User) {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{TL_contacts_getContacts{hash}, resp}
	x := <-resp
	list, ok := x.(TL_contacts_contacts)
	if !ok {
		log.Println(fmt.Sprintf("RPC: %#v", x))
		return []Contact{}, []User{}
	}
	TContacts := make([]Contact, 0, len(list.contacts))
	TUsers := make([]User, 0, len(list.users))
	for _, v := range list.contacts {
		TContacts = append(
			TContacts,
			*NewContact(v),
		)
	}
	for _, v := range list.users {
		switch u := v.(type) {
		case TL_user, TL_userEmpty:
			TUsers = append(TUsers, *NewUser(u))
		case TL_userProfilePhoto:
			TUsers[len(TUsers)-1].Photo = NewUserProfilePhoto(u)
		case TL_userStatusRecently, TL_userStatusOffline, TL_userStatusOnline, TL_userStatusLastWeek, TL_userStatusLastMonth:
			TUsers[len(TUsers)-1].Status = NewUserStatus(u)
		}
	}
	return TContacts, TUsers
}

func (m *MTProto) Channels_GetParticipants(channel TL, offset, limit int32) []User {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_channels_getParticipants{
			channel: channel,
			filter:  TL_channelParticipantsRecent{},
			offset:  offset,
			limit:   limit,
		},
		resp,
	}
	x := <-resp
	users := make([]User, 0)
	switch input := x.(type) {
	case TL_channels_channelParticipants:
		for _, u := range input.users {
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
			id: in,
		},
		resp,
	}
	x := <-resp
	chats := make([]Chat, 0, len(in))
	switch input := x.(type) {
	case TL_messages_chats:
		for _, ch := range input.chats {
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
			channel: channel,
			id:      ids,
		},
		resp,
	}
	x := <-resp
	messages := make([]Message, 0, len(ids))
	switch input := x.(type) {
	case TL_messages_messages:
		for _, m := range input.messages {
			messages = append(messages, *NewMessage(m))
		}
		return messages
	case TL_messages_messagesSlice:
		for _, m := range input.messages {
			messages = append(messages, *NewMessage(m))
		}
		return messages
	case TL_messages_channelMessages:
		for _, m := range input.messages {
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

func (m *MTProto) Messages_GetDialogs(offsetID, offsetDate, limit int32, offsetInputPeer TL) ([]Dialog, int) {
	resp := make(chan TL, 1)
	for {
		m.queueSend <- packetToSend{
			TL_messages_getDialogs{
				offset_id:   offsetID,
				offset_date: offsetDate,
				limit:       limit,
				offset_peer: offsetInputPeer,
			},
			resp,
		}
		x := <-resp
		mMessages := make(map[int32]*Message)
		mChats := make(map[int32]*Chat)
		mUsers := make(map[int32]*User)
		var dialogs []Dialog
		switch d := x.(type) {
		case TL_messages_dialogsSlice:
			for _, v := range d.messages {
				m := NewMessage(v)
				mMessages[m.ID] = m
			}
			for _, v := range d.chats {
				c := NewChat(v)
				mChats[c.ID] = c
			}
			for _, v := range d.users {
				u := NewUser(v)
				mUsers[u.ID] = u
			}
			for _, v := range d.dialogs {
				d := NewDialog(v)
				d.TopMessage = mMessages[d.TopMessageID]
				switch d.Type {
				case DIALOG_TYPE_USER:
					d.PeerAccessHash = mUsers[d.PeerID].AccessHash
					d.User = mUsers[d.PeerID]
				case DIALOG_TYPE_CHAT:
					d.Chat = mChats[d.PeerID]
				case DIALOG_TYPE_CHANNEL:
					d.PeerAccessHash = mChats[d.PeerID].AccessHash
					d.Chat = mChats[d.PeerID]
				}
				dialogs = append(dialogs, *d)
			}
			return dialogs, int(d.count)
		case TL_messages_dialogs:
			for _, v := range d.messages {
				m := NewMessage(v)
				mMessages[m.ID] = m
			}
			for _, v := range d.chats {
				c := NewChat(v)
				mChats[c.ID] = c
			}
			for _, v := range d.dialogs {
				d := NewDialog(v)
				d.TopMessage = mMessages[d.TopMessageID]
				switch d.Type {
				case DIALOG_TYPE_USER:
					d.PeerAccessHash = mUsers[d.PeerID].AccessHash
					d.User = mUsers[d.PeerID]
				case DIALOG_TYPE_CHAT:
					d.Chat = mChats[d.PeerID]
				case DIALOG_TYPE_CHANNEL:
					d.PeerAccessHash = mChats[d.PeerID].AccessHash
					d.Chat = mChats[d.PeerID]
				}
				dialogs = append(dialogs, *d)
			}
			return dialogs, len(d.chats)
		default:
			return []Dialog{}, 0
		}
	}

}

func (m *MTProto) Messages_GetChats(chatIDs []int32) []Chat {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_messages_getChats{
			id: chatIDs,
		},
		resp,
	}
	x := <-resp
	chats := make([]Chat, 0, len(chatIDs))
	switch input := x.(type) {
	case TL_messages_chats:
		for _, ch := range input.chats {
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

func (m *MTProto) Messages_GetFullChat (chatID int32) *Chat {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_messages_getFullChat{
			chat_id: chatID,
		},
		resp,
	}
	x := <-resp
	chat := new(Chat)
	switch input := x.(type) {
	case TL_messages_chatFull:
		chat =  NewChat(input)
	default:

	}
	return chat
}

func (m *MTProto) Messages_GetHistory(inputPeer TL, limit int32) ([]Message, int32) {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_messages_getHistory{
			peer:  inputPeer,
			limit: limit,
		},
		resp,
	}
	x := <-resp
	messages := make([]Message, 0, 20)
	switch input := x.(type) {
	case TL_messages_messages:
		for _, msg := range input.messages {
			messages = append(messages, *NewMessage(msg))
		}
		return messages, int32(len(messages))
	case TL_messages_messagesSlice:
		for _, msg := range input.messages {
			messages = append(messages, *NewMessage(msg))
		}
		return messages, input.count
	case TL_messages_channelMessages:
		for _, msg := range input.messages {
			messages = append(messages, *NewMessage(msg))
		}
		return messages, input.count
	case TL_rpc_error:
		fmt.Println(input.error_message, input.error_code)
		return messages, 0
	default:
		fmt.Println(reflect.TypeOf(input).String())
		return messages, 0
	}

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
			pts:  pts,
			qts:  qts,
			date: date,
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
		updateDifference.IntermediateState = *NewUpdateState(u.state)
		for _ , m := range u.new_messages {
			updateDifference.NewMessages = append(updateDifference.NewMessages, *NewMessage(m))
		}
		for _, ch := range u.chats {
			updateDifference.Chats = append(updateDifference.Chats, *NewChat(ch))
		}
		for _, user := range u.users {
			updateDifference.Users = append(updateDifference.Users, *NewUser(user))
		}
		for _, update := range u.other_updates {
			updateDifference.OtherUpdates = append(updateDifference.OtherUpdates, *NewUpdate(update))
		}
		return updateDifference
	case TL_updates_differenceSlice:
		updateDifference.IsSlice = true
		updateDifference.IntermediateState = *NewUpdateState(u.intermediate_state)
		for _ , m := range u.new_messages {
			updateDifference.NewMessages = append(updateDifference.NewMessages, *NewMessage(m))
		}
		for _, ch := range u.chats {
			updateDifference.Chats = append(updateDifference.Chats, *NewChat(ch))
		}
		for _, user := range u.users {
			updateDifference.Users = append(updateDifference.Users, *NewUser(user))
		}
		for _, update := range u.other_updates {
			updateDifference.OtherUpdates = append(updateDifference.OtherUpdates, *NewUpdate(update))
		}

		return updateDifference
	default:
		log.Println(fmt.Sprintf("RPC: %#v", x))
		return updateDifference
	}
}

func (m *MTProto) Upload_GetFile(in TL, offset, limit int32) []byte {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_upload_getFile{
			offset: offset,
			limit: limit,
			location: in,
		},
		resp,
	}
	x := <-resp
	switch f := x.(type) {
	case TL_upload_file:
		return f.bytes
	default:
		log.Println(reflect.TypeOf(f).String(), f)
	}
	return []byte{}
}
