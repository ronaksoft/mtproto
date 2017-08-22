package mtproto

import (
	"log"
	"reflect"
)

const (
	PEER_TYPE_USER    = "USER"
	PEER_TYPE_CHAT    = "CHAT"
	PEER_TYPE_CHANNEL = "CHANNEL"
)
const (
	MESSAGE_TYPE_EMPTY     = "EMPTY"
	MESSAGE_TYPE_NORMAL    = "NORMAL"
	MESSAGE_TYPE_SERVICE   = "SERVICE"
	MESSAGE_TYPE_FORWARDED = "FORWARDED"
)
const (
	MESSAGE_ACTION_CHAT_CREATED         = "CHAT_CREATED"
	MESSAGE_ACTION_CHAT_EDIT_TITLE      = "CHAT_EDIT_TITLE"
	MESSAGE_ACTION_CHAT_EDIT_PHOTO      = "CHAT_EDIT_PHOTO"
	MESSAGE_ACTION_CHAT_DELETE_PHOTO    = "CHAT_DELETE_PHOTO"
	MESSAGE_ACTION_CHAT_ADD_USER        = "CHAT_ADD_USER"
	MESSAGE_ACTION_CHAT_DELETE_USER     = "CHAT_DELETE_USER"
	MESSAGE_ACTION_CHAT_JOINED_BY_LINK  = "CHAT_JOINED_BY_LINK"
	MESSAGE_ACTION_CHAT_MIGRATE_TO      = "CHAT_MIGRATE_TO"
	MESSAGE_ACTION_CHANNEL_CREATED      = "CHANNEL_CREATED"
	MESSAGE_ACTION_CHANNEL_MIGRATE_FROM = "CHANNEL_MIGRATE"
	MESSAGE_ACTION_GAME_SCORE           = "GAME_SCORE"
	MESSAGE_ACTION_HISTORY_CLEAN        = "HISTORY_CLEAN"
)
const (
	MESSAGE_MEDIA_TYPE_EMPTY    = "EMPTY"
	MESSAGE_MEDIA_TYPE_PHOTO    = "PHOTO"
	MESSAGE_MEDIA_TYPE_VIDEO    = "VIDEO"
	MESSAGE_MEDIA_TYPE_GEO      = "GEO"
	MESSAGE_MEDIA_TYPE_CONTACT  = "CONTACT"
	MESSAGE_MEDIA_TYPE_DOCUMENT = "DOCUMENT"
	MESSAGE_MEDIA_TYPE_AUDIO    = "AUDIO"
)
const (
	USER_STATUS_OFFLINE    = "OFFLINE"
	USER_STATUS_ONLINE     = "ONLINE"
	USER_STATUS_RECENTLY   = "RECENTLY"
	USER_STATUS_LAST_WEEK  = "LAST_WEEK"
	USER_STATUS_LAST_MONTH = "LAST_MONTH"
)
const (
	DIALOG_TYPE_CHAT    = "CHAT"
	DIALOG_TYPE_USER    = "USER"
	DIALOG_TYPE_CHANNEL = "CHANNEL"
)

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
			chat_id: d.PeerID,
		}
	case DIALOG_TYPE_CHANNEL:
		return TL_inputPeerChannel{
			channel_id:  d.PeerID,
			access_hash: d.PeerAccessHash,
		}
	case DIALOG_TYPE_USER:
		return TL_inputPeerUser{
			user_id: d.PeerID,
		}
	default:
		return nil
	}
}
type Peer struct {
	Type 	string
	ID 		int32
}
type GeoPoint struct {
	Longtitude float32
	Latitude   float32
}
type FileLocation struct {
	DC       int32
	VolumeID int64
	LocalID  int32
	Secret   int64
}
type Photo struct {
	flags      int32
	ID         int64
	AccessHash int64
	UserID     int32
	Date       int32
	Caption    string
	Geo        *GeoPoint
	Sizes      []*PhotoSize
}
type PhotoSize struct {
	Type     string
	Location *FileLocation
	Width    int32
	Height   int32
	Size     int32
}
type Contact struct {
	UserID int32
	Mutual bool
}
type Document struct {
	ID          int64
	AccessHash int64
	Date        int32
	Mimetype   string
	Size        int32
	Thumb       *PhotoSize
	DcID       int32
	Version     int32

	attributes  []TL // DocumentAttribute
}
func (d *Document) GetInputFileLocation() TL_inputDocumentFileLocation {
	return TL_inputDocumentFileLocation{
		id: d.ID,
		access_hash: d.AccessHash,
	}
}
func NewPeer(in TL) (p *Peer) {
	p = new(Peer)
	switch x:= in.(type) {
	case TL_peerChannel:
		p.Type = PEER_TYPE_CHANNEL
		p.ID = x.channel_id
	case TL_peerChat:
		p.Type = PEER_TYPE_CHAT
		p.ID = x.chat_id
	case TL_peerUser:
		p.Type = PEER_TYPE_USER
		p.ID = x.user_id
	}
	return p
}
func NewFileLocation(in TL) (fl *FileLocation) {
	fl = new(FileLocation)
	switch x := in.(type) {
	case TL_fileLocationUnavailable:
		return nil
	case TL_fileLocation:
		fl.DC = x.dc_id
		fl.LocalID = x.local_id
		fl.Secret = x.secret
		fl.VolumeID = x.volume_id
	}
	return
}
func NewPhoto(in TL) (photo *Photo) {
	photo = new(Photo)
	switch x := in.(type) {
	case TL_photo:
		photo.flags = x.flags
		photo.ID = x.id
		photo.AccessHash = x.access_hash
		photo.Date = x.date
		photo.Sizes = make([]*PhotoSize, 0, len(x.sizes))
		for _, v := range x.sizes {
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
		ps.Type = x._type
		ps.Size = x.size
		ps.Width = x.w
		ps.Height = x.h
		ps.Location = NewFileLocation(x.location)
	}
	return
}
func NewDocument(in TL) (d *Document) {
	d = new(Document)
	switch x := in.(type) {
	case TL_document:
		d.ID = x.id
		d.AccessHash  = x.access_hash
		d.Mimetype = x.mime_type
		d.Date = x.date
		d.DcID = x.dc_id
		d.Size = x.size
		d.Thumb = NewPhotoSize(x.thumb)
		//TODO:: Document Attribute
	default:
		log.Println("NewDocument::", reflect.TypeOf(x).String())
	}
	return d
}
func NewDialog(input TL) (d *Dialog) {
	d = new(Dialog)
	if dialog, ok := input.(TL_dialog); ok {
		switch pt := dialog.peer.(type) {
		case TL_peerChat:
			d.Type = DIALOG_TYPE_CHAT
			d.PeerID = pt.chat_id
		case TL_peerUser:
			d.Type = DIALOG_TYPE_USER
			d.PeerID = pt.user_id
		case TL_peerChannel:
			d.Type = DIALOG_TYPE_CHANNEL
			d.PeerID = pt.channel_id
		default:
			return nil
		}
		d.Pts = dialog.pts
		d.TopMessageID = dialog.top_message
		d.UnreadCount = dialog.unread_count

		return d
	}
	return nil

}
func NewContact(in TL) (contact *Contact) {
	contact = new(Contact)
	switch c := in.(type) {
	case TL_contact:
		contact.UserID = c.user_id
		contact.Mutual = toBool(c.mutual)
	default:
		log.Println("GetContact::Error::Invalid Type")
		return nil
	}
	return
}

