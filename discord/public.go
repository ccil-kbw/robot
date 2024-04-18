package discord

import (
	"github.com/bwmarrin/discordgo"
	iqamav2 "github.com/ccil-kbw/robot/iqama/v2"
	"go.uber.org/zap"
)

var (
	publicCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "iqamatestv1",
			Description: "Get Today's Iqama",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "Masjid",
					Description: "Masjid to get Iqama for",
					Type:        discordgo.ApplicationCommandOptionString,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "@ccil_kbw, Laval QC",
							Value: "@ccil_kbw",
						},
					},
					Required: true,
				},
			},
		},
	}

	publicCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, logger *zap.Logger){
		"iqamatestv1": func(s *discordgo.Session, i *discordgo.InteractionCreate, logger *zap.Logger) {
			selectedOption := i.ApplicationCommandData().Options[0].StringValue()

			iqamaClient := iqamav2.NewIqamaCSV(selectedOption)
			resp, _ := iqamaClient.GetTodayTimes()
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: iqamaDiscordInteraction(logger, *resp),
			})
		},
	}
)
