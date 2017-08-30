package mtproto

import (
	"log"
	"reflect"
)

func (m *MTProto) Upload_GetFile(in TL, offset, limit int32) []byte {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_upload_getFile{
			Offset:   offset,
			Limit:    limit,
			Location: in,
		},
		resp,
	}
	x := <-resp
	switch f := x.(type) {
	case TL_upload_file:
		return f.Bytes
	case TL_upload_fileCdnRedirect:
	case TL_rpc_error:
		if f.error_code == 303 {
			// Migrate Code
		}
	default:
		log.Println(reflect.TypeOf(f).String(), f)
	}
	return []byte{}
}

func (m *MTProto) Upload_GetCdnFile(fileToken []byte, offset, limit int32) []byte {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_upload_getCdnFile{
			fileToken,
			offset,
			limit,
		},
		resp,
	}
	x := <-resp
	switch f := x.(type) {
	case TL_upload_cdnFileReuploadNeeded:
		m.queueSend <- packetToSend{
			TL_upload_reuploadCdnFile{
				Request_token: f.Request_token,
				File_token:    fileToken,
			},
			resp,
		}
		z := <-resp
		switch reflect.TypeOf(z).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(z)
			for i := 0; i < s.Len(); i++ {
				//hash := s.Interface().(TL_cdnFileHash)
				//TODO:: what to do now ?!!
			}
		}
	case TL_upload_cdnFile:
		return f.Bytes
	}
	return []byte{}
}
