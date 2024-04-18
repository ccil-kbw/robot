// Package discord
//
// This package needs some love, basically Discord's robot
package discord

import (
	"github.com/ccil-kbw/robot/rec"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type Bot interface {
	StartBot()
}

type bot struct {
	adminGuildID string
	session      *discordgo.Session
	logger       *zap.Logger
	isPublic     bool
}

// NewDiscordBot creates a new Discord Bot
func NewDiscordBot(logger *zap.Logger, adminGuildID, botToken string, isPublic bool) Bot {
	session, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	return &bot{
		adminGuildID: adminGuildID,
		session:      session,
		isPublic:     isPublic,
		logger:       logger,
	}
}

// StartBot starts the bot
func (d *bot) StartBot() {
	d.run(true, nil, nil)
}

// run starts the bot and listens for incoming commands
func (d *bot) run(removeCommands bool, obs *rec.Recorder, notifyChan chan string) {
	var err error

	d.addInteractionCreateHandlers()
	d.addReadyHandlers()

	err = d.session.Open()
	if err != nil {
		d.logger.Fatal("Cannot open the session", zap.Error(err))
	}

	d.logger.Info("Starting Discord Bot",
		zap.Bool("isPublic", d.isPublic),
		zap.String("InvitationLink", "https://discord.com/api/oauth2/authorize?client_id="+d.session.State.User.ID+"&scope=bot%20applications.commands"),
	)
	d.addApplicationCommands()

	defer func(s *discordgo.Session) {
		d.logger.Info("Closing Discord Bot")
		_ = s.Close()
	}(d.session)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	d.logger.Info("Press Ctrl+C to exit")
	defer d.logger.Info("Gracefully shutting down.")

	defer func() {
		if removeCommands {
			d.deleteApplicationCommands()
		}
	}()

	<-stop
	time.Sleep(10 * time.Second) // Graceful shutdown, giving some time for the bot to respond to the commands

}

func (d *bot) addInteractionCreateHandlers() {
	d.logger.Info("Adding Interaction Create Handlers")
	defer d.logger.Info("Interaction Create Handlers Added")
	d.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := publicCommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i, d.logger.With(
				zap.String("Command", i.ApplicationCommandData().Name),
				zap.String("User", i.Member.User.Username),
				zap.String("Guild", i.GuildID),
				zap.String("Channel", i.ChannelID),
				zap.String("InteractionID", i.ID),
			))
		}
	})
}

func (d *bot) addReadyHandlers() {
	d.logger.Info("Adding Ready Handlers")
	defer d.logger.Info("Ready Handlers Added")
	d.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		guildNames := make([]string, len(r.Guilds))
		guildIDs := make([]string, len(r.Guilds))
		for i, v := range r.Guilds {
			d.logger.Info("Registering Guild", zap.String("Name", v.Name), zap.String("ID", v.ID), zap.Int("MemberCount", v.MemberCount))
			guildNames[i] = v.Name
			guildIDs[i] = v.ID
		}

		d.logger.Info(
			"Logged in",
			zap.String("User", s.State.User.Username),
			zap.String("Discriminator", s.State.User.Discriminator),
			zap.Strings("GuildNames", guildNames),
			zap.Strings("GuildIDs", guildIDs),
		)
	})
}

func (d *bot) addApplicationCommands() {
	d.logger.Info("Adding Application Commands")
	defer d.logger.Info("Application Commands Added")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(publicCommands))
	for i, v := range publicCommands {
		d.logger.Info("Registering Command", zap.String("Name", v.Name), zap.String("Description", v.Description))
		cmd, err := d.session.ApplicationCommandCreate(d.session.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
}

func (d *bot) deleteApplicationCommands() {
	commands, err := d.session.ApplicationCommands(d.session.State.User.ID, "")
	if err != nil {
		d.logger.Error("Cannot fetch application commands", zap.Error(err))
	}

	for _, v := range commands {
		d.logger.Info("Deleting Command", zap.String("Command", v.Name), zap.String("GuildID", v.GuildID))
		err := d.session.ApplicationCommandDelete(d.session.State.User.ID, "", v.ID)
		if err != nil {
			d.logger.Panic("Cannot delete command", zap.String("Command", v.Name), zap.Error(err), zap.String("GuildID", v.GuildID))
		}
	}

}
