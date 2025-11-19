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

	iqamaClient iqamav2.Iqama

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, obs *rec.Recorder){
		"iqama": func(s *discordgo.Session, i *discordgo.InteractionCreate, obs *rec.Recorder) {
			resp, err := iqamaClient.GetTodayTimes()
			if err != nil {
				log.Printf("Error getting today's iqama times: %v", err)
				if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Error retrieving prayer times: %v", err),
					},
				}); err != nil {
					log.Printf("Error responding to interaction: %v", err)
				}
				return
			}
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: mappers.IqamaTimesToDiscordInteractionResponseData(*resp),
			}); err != nil {
				log.Printf("Error responding to iqama command: %v", err)
			}
		},
		"rec-schedule": func(s *discordgo.Session, i *discordgo.InteractionCreate, obs *rec.Recorder) {
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Type:        discordgo.EmbedTypeRich,
							Title:       "Recording Schedule",
							Description: "OBS Recording Schedule",
							Color:       0x05294e,
							Fields: func() []*discordgo.MessageEmbedField {
								resp := []*discordgo.MessageEmbedField{}
								// TODO: Implement actual schedule fetching
								return resp
							}(),
						},
					},
				},
			}); err != nil {
				log.Printf("Error responding to rec-schedule command: %v", err)
			}
		},
		"rec-start": func(s *discordgo.Session, i *discordgo.InteractionCreate, obs *rec.Recorder) {
			log.Println("discord command called: /rec-start")
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Type:        discordgo.EmbedTypeRich,
							Title:       "Recording Start",
							Description: "OBS Recording Start Operation",
							Color:       0x05294e,
							Fields: func() []*discordgo.MessageEmbedField {
								if obs == nil {
									return []*discordgo.MessageEmbedField{
										{Name: "Error", Value: "OBS Client not initialized"},
									}
								}
								now := time.Now()
								// TODO: Implement actual recording start logic
								return []*discordgo.MessageEmbedField{
									{
										Name:  "OBS Recording Started",
										Value: fmt.Sprintf("Recording scheduled at %d:%02d", now.Hour(), now.Minute()),
									},
								}
							}(),
						},
					},
				},
			}); err != nil {
				log.Printf("Error responding to rec-start command: %v", err)
			}
		},
		"rec-status": func(s *discordgo.Session, i *discordgo.InteractionCreate, obs *rec.Recorder) {
			log.Println("discord command called: /rec-status")
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Type:        discordgo.EmbedTypeRich,
							Title:       "Recording Status",
							Description: "OBS Recording Status",
							Color:       0x05294e,
							Fields: func() []*discordgo.MessageEmbedField {
								if obs == nil {
									return []*discordgo.MessageEmbedField{
										{Name: "Error", Value: "OBS Client not initialized"},
									}
								}
								isRecording, err := obs.IsRecording()
								if err != nil {
									log.Printf("error calling obs.IsRecording endpoint: %v", err)
									return []*discordgo.MessageEmbedField{
										{
											Name:  "Record Status",
											Value: fmt.Sprintf("Could not access OBS: %v", err),
										},
									}
								}

								log.Println("successfully called obs.IsRecording()")
								return []*discordgo.MessageEmbedField{
									{
										Name:  "Record Status",
										Value: fmt.Sprintf("Recording: %v", isRecording),
									},
								}
							}(),
						},
					},
				},
			}); err != nil {
				log.Printf("Error responding to rec-status command: %v", err)
			}
		},
	}
)

// Run starts the Discord bot and blocks until shutdown (Ctrl+C).
// It initializes the bot with the provided guild ID and token, registers slash commands,
// and sets up handlers for iqama times and OBS recording control.
// The obs parameter can be nil if recording features are disabled.
// The notifyChan is used to send prayer time notifications to a Discord channel.
// Set removeCommands to true to clean up registered commands on shutdown.
// Returns an error if bot initialization, command registration, or session setup fails.
func Run(guildID, botToken *string, removeCommands *bool, obs *rec.Recorder, notifyChan chan string) error {
	var err error
	var s *discordgo.Session
	fmt.Println("Starting Discord Bot")

	// Initialize iqama client
	iqamaClient, err = iqamav2.NewIqamaCSV(iqamav2.GetDefaultIqamaCSVPath())
	if err != nil {
		return fmt.Errorf("failed to initialize iqama client: %w", err)
	}

	s, err = discordgo.New("Bot " + *botToken)
	if err != nil {
		return fmt.Errorf("invalid bot parameters: %w", err)
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
		return fmt.Errorf("cannot open the session: %w", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *guildID, v)
		if err != nil {
			return fmt.Errorf("cannot create '%v' command: %w", v.Name, err)
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
				log.Printf("Warning: Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
	return nil
}
