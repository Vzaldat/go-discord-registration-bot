package main

import (
	"fmt"

	"github.com/Vzaldat/registration-bot/controller"
	"github.com/bwmarrin/discordgo"
)

func main() {

	dg, err := discordgo.New("")
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	dg.AddHandler(controller.MessageCreate)
	dg.AddHandler(controller.ReactionAdd)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
	}

	fmt.Println("Bot is up and running. Press ctrl + c to stop")
	<-make(chan struct{})
}
