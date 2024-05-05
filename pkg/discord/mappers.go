package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ccil-kbw/robot/pkg/iqama/v2"
	"github.com/ccil-kbw/robot/pkg/masjid_info"
	"go.uber.org/zap"
)

func iqamaDiscordInteraction(logger *zap.Logger, resp v2.IqamaDailyTimes, respMasjidInfo masjid_info.MasjidInfo) *discordgo.InteractionResponseData {
	logger.Info("iqamaDiscordInteraction", zap.Time("Date", resp.Date))
	return &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				URL:         respMasjidInfo.Website,
				Type:        discordgo.EmbedTypeRich,
				Title:       "Iqama Times for " + respMasjidInfo.Name + " on " + resp.Date.Format("2006-01-02"),
				Description: "Contributed By " + respMasjidInfo.ContributedBy,
				Color:       0x05993e,
				Image: &discordgo.MessageEmbedImage{
					// This is a test URL
					URL:    "https://cdn.discordapp.com/attachments/1159517401809952828/1232860492180099162/fares____bedouins_in_desert_discussing_night_cold_snow_far_away_4e47cde4-6cfa-423a-b966-9921d252e668.png?ex=6638d60e&is=6637848e&hm=99fe7744a0c41e7abe4b858eedd07c6a0d26eab652c9ce0fbef262704d41d620&",
					Width:  300,
					Height: 300,
				},
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
