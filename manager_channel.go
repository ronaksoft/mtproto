package mtproto

import (
	"log"
	"reflect"
	"fmt"
)

const (
	CHANNEL_DATA_EMPTY   int = 0x00
	CHANNEL_DATA_REGULAR int = 0x01
	CHANNEL_DATA_FULL    int = 0x02
)

type Channel struct {
	_State            int
	Flags             ChannelFlags
	ID                int32
	AccessHash        int64
	Title             string
	About             string
	Username          string
	Photo             TL // ChatPhoto
	Date              int32
	Version           int32
	PinnedMessageID   int32
	RestrictionReason string
	AdminRights       ChannelAdminRights  // flags_14?ChannelAdminRights
	BannedRights      ChannelBannedRights // flags_15?ChannelBannedRights
	Counters          ChannelCounters
}
type ChannelCounters struct {
	Admins       int32
	Kicked       int32
	Banned       int32
	Unread       int32
	Participants int32
}

func (ch *Channel) GetPeer() TL {
	return TL_peerChannel{
		Channel_id: ch.ID,
	}
}
func (ch *Channel) GetInputPeer() TL {
	return TL_inputPeerChannel{
		Channel_id:  ch.ID,
		Access_hash: ch.AccessHash,
	}
}

type ChannelFlags struct {
	Creator         bool // flags_0?true
	Left            bool // flags_2?true
	Editor          bool // flags_3?true
	Broadcast       bool // flags_5?true
	Verified        bool // flags_7?true
	Megagroup       bool // flags_8?true
	Restricted      bool // flags_9?true
	Democracy       bool // flags_10?true
	Signatures      bool // flags_11?true
	Min             bool // flags_12?true
	AdminRightsSet  bool //flags_14
	BannedRightsSet bool //flags_15
}

func (f *ChannelFlags) loadFlags(flags int32) {
	if flags&1<<0 != 0 {
		f.Creator = true
	}
	if flags&1<<2 != 0 {
		f.Left = true
	}
	if flags&1<<3 != 0 {
		f.Editor = true
	}
	if flags&1<<5 != 0 {
		f.Broadcast = true
	}
	if flags&1<<7 != 0 {
		f.Verified = true
	}
	if flags&1<<8 != 0 {
		f.Megagroup = true
	}
	if flags&1<<9 != 0 {
		f.Restricted = true
	}
	if flags&1<<10 != 0 {
		f.Democracy = true
	}
	if flags&1<<11 != 0 {
		f.Signatures = true
	}
	if flags&1<<12 != 0 {
		f.Min = true
	}
	if flags&1<<14 != 0 {
		f.AdminRightsSet = true
	}
	if flags&1<<15 != 0 {
		f.BannedRightsSet = true
	}
}

type ChannelAdminRights struct {
	ChangeInfo     bool // flags_0?true
	PostMessages   bool // flags_1?true
	EditMessages   bool // flags_2?true
	DeleteMessages bool // flags_3?true
	BanUsers       bool // flags_4?true
	InviteUsers    bool // flags_5?true
	InviteLink     bool // flags_6?true
	PinMessages    bool // flags_7?true
	AddAdmins      bool // flags_9?true
}

func (f *ChannelAdminRights) loadFlags(flags int32) {
	if flags&1<<0 != 0 {
		f.ChangeInfo = true
	}
	if flags&1<<1 != 0 {
		f.PostMessages = true
	}
	if flags&1<<2 != 0 {
		f.EditMessages = true
	}
	if flags&1<<3 != 0 {
		f.DeleteMessages = true
	}
	if flags&1<<4 != 0 {
		f.BanUsers = true
	}
	if flags&1<<5 != 0 {
		f.InviteUsers = true
	}
	if flags&1<<6 != 0 {
		f.InviteLink = true
	}
	if flags&1<<7 != 0 {
		f.PinMessages = true
	}
	if flags&1<<9 != 0 {
		f.AddAdmins = true
	}
}

type ChannelBannedRights struct {
	UntilDate    int32
	ViewMessages bool // flags_0?true
	SendMessages bool // flags_1?true
	SendMedia    bool // flags_2?true
	SendStickers bool // flags_3?true
	SendGifs     bool // flags_4?true
	SendGames    bool // flags_5?true
	SendInline   bool // flags_6?true
	EmbedLinks   bool // flags_7?true
}

func (f *ChannelBannedRights) loadFlags(flags int32) {
	if flags&1<<0 != 0 {
		f.ViewMessages = true
	}
	if flags&1<<1 != 0 {
		f.SendMessages = true
	}
	if flags&1<<2 != 0 {
		f.SendMedia = true
	}
	if flags&1<<3 != 0 {
		f.SendStickers = true
	}
	if flags&1<<4 != 0 {
		f.SendGifs = true
	}
	if flags&1<<5 != 0 {
		f.SendGames = true
	}
	if flags&1<<6 != 0 {
		f.SendInline = true
	}
	if flags&1<<7 != 0 {
		f.EmbedLinks = true
	}
}

// input:
//	1. TL_channelFull:
//	2. TL_channelForbidden:
//	3. TL_channel
func NewChannel(input TL) *Channel {
	channel := new(Channel)
	switch ch := input.(type) {
	case TL_channelFull:
		channel._State = CHANNEL_DATA_FULL
		channel.ID = ch.Id
		channel.About = ch.About
		channel.PinnedMessageID = ch.Pinned_msg_id
		channel.Counters.Admins = ch.Admins_count
		channel.Counters.Banned = ch.Banned_count
		channel.Counters.Kicked = ch.Kicked_count
		channel.Counters.Unread = ch.Unread_count
		channel.Counters.Participants = ch.Participants_count
		channel.Flags.loadFlags(ch.Flags)
	case TL_channelForbidden:
		channel._State = CHANNEL_DATA_EMPTY
	case TL_channel:
		channel._State = CHANNEL_DATA_REGULAR
		channel.ID = ch.Id
		channel.Title = ch.Title
		channel.AccessHash = ch.Access_hash
		channel.Username = ch.Username
		channel.Date = ch.Date
		channel.RestrictionReason = ch.Restriction_reason
		channel.Flags.loadFlags(ch.Flags)
		if channel.Flags.AdminRightsSet && ch.Admin_rights != nil {
			channel.AdminRights.loadFlags(ch.Admin_rights.(TL_channelAdminRights).Flags)
		}
		if channel.Flags.BannedRightsSet && ch.Banned_rights != nil {
			channel.BannedRights.UntilDate = ch.Banned_rights.(TL_channelBannedRights).Until_date
			channel.BannedRights.loadFlags(ch.Banned_rights.(TL_channelBannedRights).Flags)
		}
	default:
		log.Println("NewChannel::ERROR::", reflect.TypeOf(ch))
		return nil
	}
	return channel
}

func (m *MTProto) Channels_GetParticipants(channel TL, offset, limit int32) []User {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_channels_getParticipants{
			Channel: channel,
			Filter:  TL_channelParticipantsRecent{},
			Offset:  offset,
			Limit:   limit,
		},
		resp,
	}
	x := <-resp
	users := make([]User, 0)
	switch input := x.(type) {
	case TL_channels_channelParticipants:
		for _, u := range input.Users {
			users = append(users, *NewUser(u))
		}
	case TL_rpc_error:
		fmt.Println(input.error_code, input.error_message)
	default:
		fmt.Println(reflect.TypeOf(input).String())
	}
	return users
}

func (m *MTProto) Channels_GetChannels(in []TL) []Channel {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_channels_getChannels{
			Id: in,
		},
		resp,
	}
	x := <-resp
	channels := make([]Channel, 0, len(in))
	switch input := x.(type) {
	case TL_messages_chats:
		for _, ch := range input.Chats {
			channels = append(channels, *NewChannel(ch))
		}
		return channels
	case TL_rpc_error:
		fmt.Println(input.error_code, input.error_message)
		return channels
	default:
		fmt.Println(reflect.TypeOf(input).String())
		return channels
	}
}

func (m *MTProto) Channels_GetFullChannel(channelID int32, accessHash int64) *Channel {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_channels_getFullChannel{
			Channel: TL_inputChannel{
				Channel_id:  channelID,
				Access_hash: accessHash,
			},
		},
		resp,
	}
	x := <-resp
	channel := new(Channel)
	switch input := x.(type) {
	case TL_messages_chatFull:
		channel = NewChannel(input.Chats[0])
	default:
		return nil
	}
	return channel
}
