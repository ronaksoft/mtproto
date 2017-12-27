package mtproto

import (
	"fmt"
	"log"
	"github.com/pkg/errors"
)


func (m *MTProto) Auth_SendCode(phonenumber string) (string, error) {
	var authSentCode TL_auth_sentCode
	flag := true
	for flag {
		resp := make(chan TL, 1)
		m.queueSend <- packetToSend{TL_auth_sendCode{
			Flags:          1,
			Current_number: TL_boolTrue{},
			Phone_number:   phonenumber,
			Api_id:         int32(m.appId),
			Api_hash:       m.appHash,
		}, resp}
		x := <-resp
		switch x.(type) {
		case TL_auth_sentCode:
			authSentCode = x.(TL_auth_sentCode)
			flag = false
		case TL_rpc_error:
			x := x.(TL_rpc_error)
			if x.error_code != 303 {
				return "", fmt.Errorf("RPC error: %v", x)
			}
			var newDc int32
			n, _ := fmt.Sscanf(x.error_message, "PHONE_MIGRATE_%d", &newDc)
			if n != 1 {
				n, _ := fmt.Sscanf(x.error_message, "NETWORK_MIGRATE_%d", &newDc)
				if n != 1 {
					return "", fmt.Errorf("RPC error_string: %s", x.error_message)
				}
			}

			newDcAddr, ok := m.dclist[newDc]
			if !ok {
				return "", fmt.Errorf("Wrong DC index: %d", newDc)
			}
			err := m.reconnect(newDcAddr)
			fmt.Println("Reconnected")
			if err != nil {
				return "", err
			}
		default:
			return "", fmt.Errorf("Got: %T", x)
		}

	}

	if authSentCode.Flags&1 == 0 {
		return "", errors.New("Cannot sign up yet")
	}

	return authSentCode.Phone_code_hash, nil
}

func (m *MTProto) Auth_SignIn(phonenumber string, hash, code string) (TL_auth_authorization, error) {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_auth_signIn{phonenumber, hash, code},
		resp,
	}
	x := <-resp
	auth, ok := x.(TL_auth_authorization)
	if !ok {
		return TL_auth_authorization{}, fmt.Errorf("RPC: %#v", x)
	}
	userSelf := auth.User.(TL_user)
	fmt.Printf("Signed in: id %d name <%s %s>\n", userSelf.Id, userSelf.First_name, userSelf.Last_name)
	return auth, nil
}

func (m *MTProto) Auth_CheckPhone(phonenumber string) bool {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_auth_checkPhone{
			"989121228718",
		},
		resp,
	}
	x := <-resp
	if v, ok := x.(TL_auth_checkedPhone); ok {
		if toBool(v) {
			return true
		}
	}
	return false
}

func (m *MTProto) Users_GetFullSelf() (User, error) {
	return m.users_getFullUsers(TL_inputUserSelf{})
}

func (m *MTProto) users_getFullUsers(id TL) (User, error) {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_users_getFullUser{
			Id: id,
		},
		resp,
	}
	x := <-resp
	user, ok := x.(TL_userFull)
	if !ok {
		log.Println(fmt.Sprintf("RPC: %#v", x))
		return User{}, fmt.Errorf("RPC: %#v", x)
	}
	newUser := NewUser(user.User)
	return *newUser, nil
}
