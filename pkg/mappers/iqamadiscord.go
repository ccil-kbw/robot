package mappers

import (
	"github.com/bwmarrin/discordgo"
	iqamav1 "github.com/ccil-kbw/robot/pkg/iqama/v1"
)

func IqamaTimesToDiscordInteractionResponseData(resp iqamav1.Resp) *discordgo.InteractionResponseData {
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
							Value: resp.Fajr.Iqama,
						},
						{
							Name:  "Dhuhr",
							Value: resp.Dhuhr.Iqama,
						},
						{
							Name:  "Asr",
							Value: resp.Asr.Iqama,
						},
						{
							Name:  "Maghrib",
							Value: resp.Maghrib.Iqama,
						},
						{
							Name:  "Isha",
							Value: resp.Isha.Iqama,
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
