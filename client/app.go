// Example of using Mtproto for Telegram bot
package main

import (
"fmt"
"os"
"path/filepath"

"github.com/xelaj/go-dry"

"github.com/k0kubun/pp"
"github.com/riversgo007/EvaBot/mtproto/telegram"

utils "github.com/riversgo007/EvaBot/mtproto/examples/example_utils"
)

const (
	appID   = 4806267
	appHash = "d90c5efaf477553e7e2b7a168cacbd02"
	// If you don't know how to get your token, read this guide: https://core.telegram.org/bots
	botToken = "1706227378:AAE1i8R1CgpG_WJd7wr6y_bned_ggt0NJ7o"
	// bot username IS NOT required to perform auth, it's using to visualize example
	botUsername = "superamazingbot"
)

func main() {
	// helper variables
	appStorage := utils.PrepareAppStorageForExamples()
	sessionFile := filepath.Join(appStorage, "session.json")
	publicKeys := filepath.Join(appStorage, "tg_public_keys.pem")

	client, err := telegram.NewClient(telegram.ClientConfig{
		// where to store session configuration. must be set
		SessionFile: sessionFile,
		// host address of mtproto server. Actually, it can'be mtproxy, not only official
		ServerHost: "149.154.167.50:443",
		// public keys file is patrh to file with public keys, which you must get from https://my.telelgram.org
		PublicKeysFile:  publicKeys,
		AppID:           appID,   // app id, could be find at https://my.telegram.org
		AppHash:         appHash, // app hash, could be find at https://my.telegram.org
		InitWarnChannel: true,    // if we want to get errors, otherwise, client.Warnings will be set nil
	})
	utils.ReadWarningsToStdErr(client.Warnings)
	dry.PanicIfErr(err)

	// Trying to auth as bot with our bot token
	_, err = client.AuthImportBotAuthorization(
		1, // flags, it's reserved, must be set (don't mind how does it works, we don't know too)
		appID,
		appHash,
		botToken,
	)
	if err != nil {
		fmt.Println("ImportBotAuthorization error:", err.Error())
		os.Exit(1)
	}

	// Request info about username of our bot, this is not efficient way, we just want to
	// test, did auth succeed or not
	//client.GetChannelInfoByInviteLink()
	uname, err := client.ContactsResolveUsername(botUsername)
	if err != nil {
		fmt.Println("ResolveUsername error:", err.Error())
		os.Exit(1)
	}

	chatsCount := len(uname.Chats)
	// No chats for requested username of our bot
	if chatsCount > 0 {
		fmt.Println("Chats number:", chatsCount)
		os.Exit(1)
	}

	usersCount := len(uname.Users)
	// Users vector must contain single item with information about our bot
	if usersCount != 1 {
		fmt.Println("Users number:", usersCount)
		os.Exit(1)
	}

	user := uname.Users[0].(*telegram.UserObj)

	pp.Println(user)
}
