package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ccil-kbw/robot/iqama/v2"
	"go.uber.org/zap"
)

func iqamaDiscordInteraction(logger *zap.Logger, resp v2.IqamaDailyTimes) *discordgo.InteractionResponseData {
	logger.Info("iqamaDiscordInteraction", zap.Time("Date", resp.Date))
	return &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				URL:         "https://ccil-kbw.com/iqama",
				Type:        discordgo.EmbedTypeRich,
				Title:       "Iqama Time",
				Description: "Iqama pulled from https://ccil-kbw.com/iqama",
				Color:       0x05993e,
				Fields: func() []*discordgo.MessageEmbedField {
					return []*discordgo.MessageEmbedField{
						{
							Name:  "Fajr",
							Value: v2.FormatTime(resp.Fajr.Iqama),
						},
						{
							Name:  "Dhuhr",
							Value: v2.FormatTime(resp.Dhuhr.Iqama),
						},
						{
							Name:  "Asr",
							Value: v2.FormatTime(resp.Asr.Iqama),
						},
						{
							Name:  "Maghrib",
							Value: v2.FormatTime(resp.Maghrib.Iqama),
						},
						{
							Name:  "Isha",
							Value: v2.FormatTime(resp.Isha.Iqama),
						},
						{
							Name:  "Friday Prayer 1",
							Value: "Traduction en Francais: 12:00, Arabe: 12:15",
						},
						{
							Name:  "Friday Prayer 2",
							Value: "Traduction en Anglais: 13:00, Arabe: 13:15",
						},
					}
				}(),
			},
		},
	}
}
