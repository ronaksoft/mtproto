package mtproto

import (
	"log"
	"reflect"
)

// Peer
type Peer struct {
	Type string
	ID   int32
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

type GeoPoint struct {
	Longtitude float32
	Latitude   float32
}


// Photo
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

// PhotoSize
type PhotoSize struct {
	Type     string
	Location *FileLocation
	Width    int32
	Height   int32
	Size     int32
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
func (p *PhotoSize) GetInputFileLocation() TL_inputFileLocation {
	return TL_inputFileLocation{
		p.Location.VolumeID,
		p.Location.LocalID,
		p.Location.Secret,
	}
}


// Document
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
func (d *Document) GetInputFileLocation() TL_inputDocumentFileLocation {
	return TL_inputDocumentFileLocation{
		Id:          d.ID,
		Access_hash: d.AccessHash,
	}
}

// FileLocation
type FileLocation struct {
	DC       int32
	VolumeID int64
	LocalID  int32
	Secret   int64
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
func (f *FileLocation) GetInputFileLocation() TL_inputFileLocation {
	return TL_inputFileLocation{
		f.VolumeID,
		f.LocalID,
		f.Secret,
	}
}




