// Package discord
//
// This package needs some love, basically Discord's robot
package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ccil-kbw/robot/mappers"
	"github.com/ccil-kbw/robot/rec"

	"github.com/bwmarrin/discordgo"
	iqamav2 "github.com/ccil-kbw/robot/iqama/v2"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "iqama",
			Description: "Get Today's Iqama",
		},
		{
			Name:        "rec-start",
			Description: "Start Recording on the Main Camera (e.g unscheduled speech @ccil-kbw)",
		},
		{
			Name:        "rec-status",
			Description: "See OBS and Recording Status on the Main Camera",
		},
		{
			Name:        "rec-schedule",
			Description: "See Recording Schedule",
		},
	}

	iqamaClient = iqamav2.NewIqamaCSV("iqama_2025.csv")

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, obs *rec.Recorder){
		"iqama": func(s *discordgo.Session, i *discordgo.InteractionCreate, obs *rec.Recorder) {
			resp, _ := iqamaClient.GetTodayTimes()
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: mappers.IqamaTimesToDiscordInteractionResponseData(*resp),
			})
		},
		"rec-schedule": func(s *discordgo.Session, i *discordgo.InteractionCreate, obs *rec.Recorder) {
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Type:        discordgo.EmbedTypeRich,
							Title:       "Scheduling Start",
							Description: "OBS Scheduling Start Operation",
							Color:       0x05294e,
							Fields: func() []*discordgo.MessageEmbedField {
								resp := []*discordgo.MessageEmbedField{}

								return resp
							}(),
						},
					},
				},
			},
			)
		},
		"rec-start": func(s *discordgo.Session, i *discordgo.InteractionCreate, obs *rec.Recorder) {
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Type:        discordgo.EmbedTypeRich,
							Title:       "Scheduling Start",
							Description: "OBS Scheduling Start Operation",
							Color:       0x05294e,
							Fields: func() []*discordgo.MessageEmbedField {
								fmt.Println("discord command called: /obs-start")
								if obs == nil {
									return []*discordgo.MessageEmbedField{
										{Name: "Error", Value: "OBS Client not initialized"},
									}
								}
								now := time.Now()

								return []*discordgo.MessageEmbedField{
									{
										Name:  "OBS Recording Started",
										Value: fmt.Sprintf("Will be scheduling until %d:%d", now.Hour(), now.Minute()),
									},
								}
							}(),
						},
					},
				},
			},
			)
		},
		"rec-status": func(s *discordgo.Session, i *discordgo.InteractionCreate, obs *rec.Recorder) {
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Type:        discordgo.EmbedTypeRich,
							Title:       "Scheduling Status",
							Description: "OBS Scheduling Status Operation",
							Color:       0x05294e,
							Fields: func() []*discordgo.MessageEmbedField {
								fmt.Println("discord command called: /obs-status")
								if obs == nil {
									return []*discordgo.MessageEmbedField{
										{Name: "Error", Value: "OBS Client not initialized"},
									}
								}
								isRecording, err := obs.IsRecording()
								if err != nil {
									fmt.Printf("error calling obs.IsRecording endpoint. %v", err)
									return []*discordgo.MessageEmbedField{
										{
											Name:  "Record Status",
											Value: "Could not access OBS",
										},
									}
								}

								fmt.Println("successfully called obs.IsRecording()")
								return []*discordgo.MessageEmbedField{
									{
										Name:  "Record Status",
										Value: fmt.Sprintf("recording: %v", isRecording),
									},
								}
							}(),
						},
					},
				},
			},
			)
		},
	}
)

// Run the Discord bot. NOTE: Function can be split
func Run(guildID, botToken *string, removeCommands *bool, obs *rec.Recorder, notifyChan chan string) {
	var err error
	var s *discordgo.Session
	fmt.Println("Starting Discord Bot")
	s, err = discordgo.New("Bot " + *botToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i, obs)
		}
	})

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *guildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer func(s *discordgo.Session) {
		_ = s.Close()
	}(s)

	var channelID string
	{
		channelID = os.Getenv("MDROID_DISCORD_CHANNEL_ID")
		if channelID == "" {
			channelID = "1161464871792160789"
		}
	}
	go func() {
		previousMessageID := ""
		for {
			m, err := s.ChannelMessageSend(channelID, <-notifyChan)
			if err != nil {
				log.Println("Couldn't send Message")
			}
			if previousMessageID != "" {
				if err := s.ChannelMessageDelete(channelID, previousMessageID); err != nil {
					log.Println("Couldn't delete the Previous Message, please clean up the channel manually.")
				}
			}
			previousMessageID = m.ID
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	fmt.Println("Press Ctrl+C to exit")
	<-stop

	if *removeCommands {
		log.Println("Removing commands...")

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *guildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
