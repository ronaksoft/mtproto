package mtproto

import (
	"log"
	"reflect"
)

func NewInputPeerUser(userID int32, accessHash int64) TL {
	return TL_inputPeerUser{
		User_id:     userID,
		Access_hash: accessHash,
	}
}
func NewInputPeerChat(chatID int32) TL {
	return TL_inputPeerChat{
		Chat_id: chatID,
	}
}
func NewInputPeerChannel(channelID int32, accessHash int64) TL {
	return TL_inputPeerChannel{
		Channel_id:  channelID,
		Access_hash: accessHash,
	}
}
func NewPeerChat(chatID int32) TL {
	return TL_peerChat{
		Chat_id: chatID,
	}
}
func NewPeerChannel(channelID int32) TL {
	return TL_peerChannel{
		Channel_id: channelID,
	}
}

func NewPeer(in TL) (p *Peer) {
	p = new(Peer)
	switch x := in.(type) {
	case TL_peerChannel:
		p.Type = PEER_TYPE_CHANNEL
		p.ID = x.Channel_id
	case TL_peerChat:
		p.Type = PEER_TYPE_CHAT
		p.ID = x.Chat_id
	case TL_peerUser:
		p.Type = PEER_TYPE_USER
		p.ID = x.User_id
	}
	return p
}
func NewFileLocation(in TL) (fl *FileLocation) {
	fl = new(FileLocation)
	switch x := in.(type) {
	case TL_fileLocationUnavailable:
		return nil
	case TL_fileLocation:
		fl.DC = x.Dc_id
		fl.LocalID = x.Local_id
		fl.Secret = x.Secret
		fl.VolumeID = x.Volume_id
	}
	return
}
func NewPhoto(in TL) (photo *Photo) {
	photo = new(Photo)
	switch x := in.(type) {
	case TL_photo:
		photo.flags = x.Flags
		photo.ID = x.Id
		photo.AccessHash = x.Access_hash
		photo.Date = x.Date
		photo.Sizes = make([]*PhotoSize, 0, len(x.Sizes))
		for _, v := range x.Sizes {
			photo.Sizes = append(photo.Sizes, NewPhotoSize(v))
		}
	default:
		return nil
	}
	return
}
func NewPhotoSize(in TL) (ps *PhotoSize) {
	ps = new(PhotoSize)
	switch x := in.(type) {
	case TL_photoSizeEmpty:
		return nil
	case TL_photoSize:
		ps.Type = x._Type
		ps.Size = x.Size
		ps.Width = x.W
		ps.Height = x.H
		ps.Location = NewFileLocation(x.Location)
	}
	return
}
func NewDocument(in TL) (d *Document) {
	d = new(Document)
	switch x := in.(type) {
	case TL_document:
		d.ID = x.Id
		d.AccessHash = x.Access_hash
		d.Mimetype = x.Mime_type
		d.Date = x.Date
		d.DcID = x.Dc_id
		d.Size = x.Size
		d.Thumb = NewPhotoSize(x.Thumb)
		//TODO:: Document Attribute
	default:
		log.Println("NewDocument::", reflect.TypeOf(x).String())
	}
	return d
}
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
func NewContact(in TL) (contact *Contact) {
	contact = new(Contact)
	switch c := in.(type) {
	case TL_contact:
		contact.UserID = c.User_id
		contact.Mutual = toBool(c.Mutual)
	default:
		log.Println("GetContact::Error::Invalid Type")
		return nil
	}
	return
}

