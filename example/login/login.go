package main

import (
    "strings"
    "github.com/ronaksoft/mtproto"
    "log"
    "fmt"
)

var (
    _MT *mtproto.MTProto
)

func main() {
    appId := int64(48841)
    appHash := "3151c01673d412c18c055f089128be50"
    if v, err := mtproto.NewMTProto(appId, appHash,"../auth_key", "", 0); err != nil {
        log.Println(err.Error())
        return
    } else {
        _MT = v
        if err := _MT.Connect(); err != nil {
            log.Println("Connect:", err.Error())
        }
    }
    if phoneCodeHash, err := _MT.Auth_SendCode("989121228718"); err != nil {
        log.Println("SendCode:", err.Error())
    } else {
        var phoneCode string
        fmt.Print("Enter Code:")
        fmt.Scanln(&phoneCode)
        phoneCode = strings.TrimSpace(phoneCode)
        fmt.Println("Code:", phoneCode)
        _MT.Auth_SignIn("989121228718", phoneCodeHash, phoneCode)
    }

}

