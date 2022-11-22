package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "iqama",
			Description: "Get Today's Iqama",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"iqama": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: func() string {
						resp, err := http.Get("https://iqama.ccil-kbw.com/iqamatimes.php")
						if err != nil {
							log.Fatalln(err)
						}
						defer resp.Body.Close()

						body, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							log.Fatalln(err)
						}
						return iqamaPrettify(body)
					}(),
				},
			})
		},
	}
)
