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

type Peer struct {
	Type string
	ID   int32
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
	ID         int64
	AccessHash int64
	Date       int32
	Mimetype   string
	Size       int32
	Thumb      *PhotoSize
	DcID       int32
	Version    int32

	attributes []TL // DocumentAttribute
}

func (d *Document) GetInputFileLocation() TL_inputDocumentFileLocation {
	return TL_inputDocumentFileLocation{
		Id:          d.ID,
		Access_hash: d.AccessHash,
	}
}
func (p *PhotoSize) GetInputFileLocation() TL_inputFileLocation {
	return TL_inputFileLocation{
		p.Location.VolumeID,
		p.Location.LocalID,
		p.Location.Secret,
	}
}
func (f *FileLocation) GetInputFileLocation() TL_inputFileLocation {
	return TL_inputFileLocation{
		f.VolumeID,
		f.LocalID,
		f.Secret,
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
