package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"os"
	"github.com/plasmakatt/bdobot/chatbots/discord"
)

type Auth struct {
	Token string
}

func readAuth(authFile string) Auth {
	raw, err := ioutil.ReadFile(authFile)
	if err != nil {
		fmt.Println("An error occurred while reading auth file:", err)
		os.Exit(1)
	}
	var a Auth
	json.Unmarshal(raw, &a)
	return a
} 

func main() {
	auth := readAuth("./auth.json")
	bot := discord.DiscordBot{Token: auth.Token}
	bot.Run()
}