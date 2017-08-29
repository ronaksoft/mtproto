package mtproto

import (
	"log"
	"reflect"
)


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

