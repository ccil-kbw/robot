package discord

import (
	"github.com/bwmarrin/discordgo"
	iqamav2 "github.com/ccil-kbw/robot/iqama/v2"
	"go.uber.org/zap"
	"os"
	"strings"
)

var (
	publicCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "iqamatestv1",
			Description: "Get Today's Iqama",
			Options:     getCitySubcommands(),
		},
	}

	cityNamesMap = make(map[string]string)

	publicCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, logger *zap.Logger){
		"iqamatestv1": func(s *discordgo.Session, i *discordgo.InteractionCreate, logger *zap.Logger) {
			cityName := cityNamesMap[i.ApplicationCommandData().Options[0].Name]
			masjidName := i.ApplicationCommandData().Options[0].Options[0].StringValue()
			choice := cityName + "/" + masjidName
			iqamaClient := iqamav2.NewIqamaCSV(choice)
			resp, _ := iqamaClient.GetTodayTimes()
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: iqamaDiscordInteraction(logger, *resp),
			})
		},
	}
)

func getCitySubcommands() []*discordgo.ApplicationCommandOption {
	cities, err := os.ReadDir("assets/masjids_data")
	if err != nil {
		panic(err)
	}

	var subcommands []*discordgo.ApplicationCommandOption
	for _, city := range cities {
		if city.IsDir() {
			originalCityName := city.Name()
			lowerCaseCityName := strings.ToLower(strings.ReplaceAll(originalCityName, " ", "_"))
			cityNamesMap[lowerCaseCityName] = originalCityName
			subcommands = append(subcommands, &discordgo.ApplicationCommandOption{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        lowerCaseCityName,
				Description: "Choose a masjid in " + originalCityName,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "masjid",
						Description: "Masjid to get Iqama for",
						Choices:     getMasjidChoices("assets/masjids_data/" + originalCityName),
						Required:    true,
					},
				},
			})
		}
	}

	return subcommands
}

func getMasjidChoices(cityPath string) []*discordgo.ApplicationCommandOptionChoice {
	masjids, err := os.ReadDir(cityPath)
	if err != nil {
		panic(err)
	}

	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, masjid := range masjids {
		if masjid.IsDir() {
			masjidName := masjid.Name()
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  masjidName,
				Value: masjidName,
			})
		}
	}

	return choices
}
