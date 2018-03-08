package mtproto

type Dialog struct {
    Type           string
    Pts            int32
    PeerID         int32
    PeerAccessHash int64
    TopMessageID   int32
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

func (m *MTProto) Messages_GetDialogs(
    offsetID, offsetDate, limit int32, offsetInputPeer TL,
) ([]Dialog, map[int32]User, map[int32]Chat, map[int32]Channel, map[int32]Message, int) {
    resp := make(chan TL, 1)

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
    mMessages := make(map[int32]Message)
    mChats := make(map[int32]Chat)
    mChannels := make(map[int32]Channel)
    mUsers := make(map[int32]User)
    var dialogs []Dialog
    switch input := x.(type) {
    case TL_messages_dialogsSlice:
        for _, v := range input.Messages {
            m := NewMessage(v)
            if m != nil {
                mMessages[m.ID] = *m
            }

        }
        for _, v := range input.Chats {
            switch v.(type) {
            case TL_chatEmpty, TL_chat, TL_chatFull, TL_chatForbidden:
                c := NewChat(v)
                mChats[c.ID] = *c
            case TL_channel, TL_channelFull, TL_channelForbidden:
                c := NewChannel(v)
                mChannels[c.ID] = *c
            }
        }
        for _, v := range input.Users {
            u := NewUser(v)
            mUsers[u.ID] = *u
        }
        for _, v := range input.Dialogs {
            d := NewDialog(v)
            switch d.Type {
            case DIALOG_TYPE_USER:
                d.PeerAccessHash = mUsers[d.PeerID].AccessHash
            case DIALOG_TYPE_CHAT:
                d.PeerAccessHash = mChats[d.PeerID].AccessHash
            case DIALOG_TYPE_CHANNEL:
                d.PeerAccessHash = mChannels[d.PeerID].AccessHash
            }
            dialogs = append(dialogs, *d)
        }
        return dialogs, mUsers, mChats, mChannels, mMessages, int(input.Count)
    case TL_messages_dialogs:
        for _, v := range input.Messages {
            m := NewMessage(v)
            if m != nil {
                mMessages[m.ID] = *m
            }

        }
        for _, v := range input.Chats {
            switch v.(type) {
            case TL_chatEmpty, TL_chat, TL_chatFull, TL_chatForbidden:
                c := NewChat(v)
                mChats[c.ID] = *c
            case TL_channel, TL_channelFull, TL_channelForbidden:
                c := NewChannel(v)
                mChannels[c.ID] = *c
            }
        }
        for _, v := range input.Users {
            u := NewUser(v)
            mUsers[u.ID] = *u
        }
        for _, v := range input.Dialogs {
            d := NewDialog(v)
            switch d.Type {
            case DIALOG_TYPE_USER:
                d.PeerAccessHash = mUsers[d.PeerID].AccessHash
            case DIALOG_TYPE_CHAT:
                d.PeerAccessHash = mUsers[d.PeerID].AccessHash
            case DIALOG_TYPE_CHANNEL:
                d.PeerAccessHash = mChannels[d.PeerID].AccessHash
            }
            dialogs = append(dialogs, *d)
        }
        return dialogs, mUsers, mChats, mChannels, mMessages, len(input.Chats)
    default:
        return nil, nil, nil, nil, nil, 0
    }

}
