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

