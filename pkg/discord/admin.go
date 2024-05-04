package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/ccil-kbw/robot/pkg/rec"
	"time"
)

var (
	adminCommands = []*discordgo.ApplicationCommand{
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

	adminCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, obs *rec.Recorder){
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
