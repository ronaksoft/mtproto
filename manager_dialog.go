package mtproto


type Dialog struct {
	Type           string
	Pts            int32
	PeerID         int32
	PeerAccessHash int64
	TopMessageID   int32
	TopMessage     *Message
	Chat           *Chat
	User           *User
	UnreadCount    int32
	NotifySettings interface{}
}

// NewDialog returns a pointer to Dialog struct
// input :		TL_dialog
func NewDialog(input TL) (d *Dialog) {
	d = new(Dialog)
	if dialog, ok := input.(TL_dialog); ok {
		switch pt := dialog.Peer.(type) {
		case TL_peerChat:
			d.Type = DIALOG_TYPE_CHAT
			d.PeerID = pt.Chat_id
		case TL_peerUser:
			d.Type = DIALOG_TYPE_USER
			d.PeerID = pt.User_id
		case TL_peerChannel:
			d.Type = DIALOG_TYPE_CHANNEL
			d.PeerID = pt.Channel_id
		default:
			return nil
		}
		d.Pts = dialog.Pts
		d.TopMessageID = dialog.Top_message
		d.UnreadCount = dialog.Unread_count

		return d
	}
	return nil

}

// GetInputPeer returns either of the struct below:
//	1. TL_inputPeerChat
//	2. TL_inputPeerChannel
//	3. TL_inputPeerUser
func (d *Dialog) GetInputPeer() TL {
	switch d.Type {
	case DIALOG_TYPE_CHAT:
		return TL_inputPeerChat{
			Chat_id: d.PeerID,
		}
	case DIALOG_TYPE_CHANNEL:
		return TL_inputPeerChannel{
			Channel_id:  d.PeerID,
			Access_hash: d.PeerAccessHash,
		}
	case DIALOG_TYPE_USER:
		return TL_inputPeerUser{
			User_id: d.PeerID,
		}
	default:
		return nil
	}
}


func (m *MTProto) Messages_GetDialogs(offsetID, offsetDate, limit int32, offsetInputPeer TL) ([]Dialog, int) {
	resp := make(chan TL, 1)
	for {
		m.queueSend <- packetToSend{
			TL_messages_getDialogs{
				Offset_id:   offsetID,
				Offset_date: offsetDate,
				Limit:       limit,
				Offset_peer: offsetInputPeer,
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
			for _, v := range d.Messages {
				m := NewMessage(v)
				mMessages[m.ID] = m
			}
			for _, v := range d.Chats {
				c := NewChat(v)
				mChats[c.ID] = c
			}
			for _, v := range d.Users {
				u := NewUser(v)
				mUsers[u.ID] = u
			}
			for _, v := range d.Dialogs {
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
			return dialogs, int(d.Count)
		case TL_messages_dialogs:
			for _, v := range d.Messages {
				m := NewMessage(v)
				mMessages[m.ID] = m
			}
			for _, v := range d.Chats {
				c := NewChat(v)
				mChats[c.ID] = c
			}
			for _, v := range d.Dialogs {
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
			return dialogs, len(d.Chats)
		default:
			return []Dialog{}, 0
		}
	}

}