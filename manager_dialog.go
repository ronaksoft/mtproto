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

