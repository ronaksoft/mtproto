package cmd

import (
    "github.com/spf13/cobra"
    "github.com/ronaksoft/mtproto"
    "fmt"
)

var (
    _MT *mtproto.MTProto
)

var RootCmd = &cobra.Command{
    Use: "tg",
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        appId := int64(48841)
        appHash := "3151c01673d412c18c055f089128be50"
        account := cmd.Flag("account").Value.String()
        if v, err := mtproto.NewMTProto(
            appId,
            appHash,
            fmt.Sprintf("../%s_auth_key", account), "", 0); err != nil {
            fmt.Println(err.Error())
            return
        } else {
            _MT = v
            if err := _MT.Connect(); err != nil {
                fmt.Println("Connect:", err.Error())
            }
        }
    },
}

func init() {
    RootCmd.PersistentFlags().String("account", "unknown", "")
}
