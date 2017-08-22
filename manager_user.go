package mtproto

type UserStatus struct {
	Status    string
	Online    bool
	Timestamp int32
}
type UserProfilePhoto struct {
	ID         int64
	PhotoSmall FileLocation
	PhotoLarge FileLocation
}
type User struct {
	flags                int32
	ID                   int32
	Username             string
	FirstName            string
	LastName             string
	Phone                string
	Photo                *UserProfilePhoto
	Status               *UserStatus
	Inactive             bool
	Mutual               bool
	Verified             bool
	Restricted           bool
	AccessHash           int64
	BotInfoVersion       int32
	BotInlinePlaceHolser string
	RestrictionReason    string
}
func (user *User) IsSelf() bool {
	if user.flags&1<<10 != 0 {
		return true
	}
	return false
}
func (user *User) IsContact() bool {
	if user.flags&1<<11 != 0 {
		return true
	}
	return false
}
func (user *User) IsMutualContact() bool {
	if user.flags&1<<12 != 0 {
		return true
	}
	return false
}
func (user *User) IsDeleted() bool {
	if user.flags&1<<13 != 0 {
		return true
	}
	return false
}
func (user *User) IsBot() bool {
	if user.flags&1<<14 != 0 {
		return true
	}
	return false
}
func (user *User) GetInputPeer() TL {
	if user.IsSelf() {
		return TL_inputPeerSelf{}
	} else {
		return TL_inputPeerUser{}
	}
}
func (user *User) GetPeer() TL {
	return TL_peerUser{
		user_id: user.ID,
	}
}
func NewUserStatus(userStatus TL) (s *UserStatus) {
	s = new(UserStatus)
	switch status := userStatus.(type) {
	case TL_userStatusEmpty:
		return nil
	case TL_userStatusOnline:
		s.Status = USER_STATUS_ONLINE
		s.Online = true
		s.Timestamp = status.expires
	case TL_userStatusOffline:
		s.Status = USER_STATUS_OFFLINE
		s.Online = false
		s.Timestamp = status.was_online
	case TL_userStatusRecently:
		s.Status = USER_STATUS_RECENTLY
		s.Online = false
	case TL_userStatusLastWeek:
		s.Status = USER_STATUS_LAST_WEEK
	case TL_userStatusLastMonth:
		s.Status = USER_STATUS_LAST_MONTH
	}
	return
}
func NewUserProfilePhoto(userProfilePhoto TL) (u *UserProfilePhoto) {
	u = new(UserProfilePhoto)
	switch pp := userProfilePhoto.(type) {
	case TL_userProfilePhotoEmpty:
		return nil
	case TL_userProfilePhoto:
		u.ID = pp.photo_id
		switch big := pp.photo_big.(type) {
		case TL_fileLocationUnavailable:
		case TL_fileLocation:
			u.PhotoLarge.DC = big.dc_id
			u.PhotoLarge.LocalID = big.local_id
			u.PhotoLarge.Secret = big.secret
			u.PhotoLarge.VolumeID = big.volume_id
		}
		switch small := pp.photo_small.(type) {
		case TL_fileLocationUnavailable:
		case TL_fileLocation:
			u.PhotoSmall.DC = small.dc_id
			u.PhotoSmall.LocalID = small.local_id
			u.PhotoLarge.Secret = small.secret
			u.PhotoSmall.VolumeID = small.volume_id
		}
	}
	return
}
func NewUser(in TL) (user *User) {
	user = new(User)
	switch u := in.(type) {
	case TL_userEmpty:
		user.ID = u.id
	case TL_user:
		user.ID = u.id
		user.Username = u.username
		user.FirstName = u.first_name
		user.LastName = u.last_name
		user.AccessHash = u.access_hash
		user.BotInfoVersion = u.bot_info_version
		user.BotInlinePlaceHolser = u.bot_inline_placeholder
		user.RestrictionReason = u.restriction_reason
		user.Phone = u.phone
		if u.flags&1<<5 != 0 {
			user.Photo = NewUserProfilePhoto(u.photo)
		}
		if u.flags&1<<6 != 0 {
			user.Status = NewUserStatus(u.status)
		}

	default:
		//fmt.Println(reflect.TypeOf(u).String())
		return nil
	}
	return
}