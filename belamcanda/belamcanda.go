// Package belamcanda contains bot functionality for iris. Currently following Pomodoro technique.
package belamcanda

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/TensRoses/iris/internal/config"
	"github.com/TensRoses/iris/internal/db"
	"github.com/TensRoses/iris/internal/log"
)

const (
	logLevel            int           = 2 // refers to internal/log/log.go for level definition
	defaultPomDuration  time.Duration = 25 * time.Minute
	discordBotPrefix    string        = "Bot "
	baseAuthURLTemplate string        = "https://discord.com/api/oauth2/authorize?client_id=%s&scope=bot"
)

// pomDuration defines default sessions (should always be 25 mins).
var pomDuration time.Duration

type cmdHandler func(s *discordgo.Session, m *discordgo.MessageCreate, ex string)

type botCommand struct {
	handler       cmdHandler
	desc          string
	exampleParams string
}

// Iris defines the structure for the bot's functionality.
type Iris struct {
	helpMessage   string
	inviteMessage string
	Config        config.Configs
	secrets       config.Secrets
	discord       *discordgo.Session
	logger        log.Logging
	cmdHandlers   map[string]botCommand
	poms          db.UserPomodoroMap
	// record metrics here
	// metrics metrics.Recorder
}

// NewIris creates a new instance of Iris that can deploy over Heroku.
func NewIris(config config.Configs, secrets config.Secrets, logger log.Logging) *Iris {
	// setup new logLevel
	logger.SetLoggingLevel(logLevel)

	ir := &Iris{
		Config:  config,
		secrets: secrets,
		logger:  logger,
		poms:    db.NewUserPomodoroMap(),
	}

	ir.registerCmdHandlers()
	ir.helpMessage = ir.buildHelpMessage()
	ir.inviteMessage = fmt.Sprintf("Click here: <"+baseAuthURLTemplate+"> to invite me to the server", ir.secrets.ClientID)
	return ir
}

func (ir *Iris) registerCmdHandlers() {
	ir.cmdHandlers = map[string]botCommand{
		"help":   {handler: ir.onCmdHelp, desc: "Show this help message", exampleParams: ""},
		"pom":    {handler: ir.onCmdStartPom, desc: "Start a pom work cycle. You can optionally specify the period of time (default: 25 mins)", exampleParams: "50"},
		"stop":   {handler: ir.onCmdCancelPom, desc: "cancle current pom cycle", exampleParams: ""},
		"invite": {handler: ir.onCmdInvite, desc: "Create an invite link you can use to have the bot join the server", exampleParams: ""},
		// "simp":   {handler: ir.onCmdSimp, desc: "notify another friend with the good stuff", exampleParams: ""},
	}
}

func (ir *Iris) buildHelpMessage() string {
	helpBuffer := bytes.Buffer{}
	helpBuffer.WriteString("Made by **@aarnphm**\n")

	// just use map iteration order
	for cmdStr, cmd := range ir.cmdHandlers {
		helpBuffer.WriteString(fmt.Sprintf("\n•  **%s**  -  %s\n", cmdStr, cmd.desc))
		helpBuffer.WriteString(fmt.Sprintf("   Example: `%s%s %s`\n", ir.Config.CmdPrefix, cmdStr, cmd.exampleParams))
	}

	helpBuffer.WriteString("\n" + ir.inviteMessage)

	return helpBuffer.String()
}

// Start will start the bot, blocking til completed.
func (ir *Iris) Start() error {
	if ir.secrets.AuthToken == "" {
		return errors.New("no authToken found")
	}

	var err error
	ir.discord, err = discordgo.New(discordBotPrefix + ir.secrets.AuthToken)
	if err != nil {
		return err
	}

	// onReady will prepare our metrics, which will get from prometheus
	ir.discord.AddHandler(ir.onReady)
	ir.discord.AddHandler(ir.onMessageReceived)

	err = ir.discord.Open()
	if err != nil {
		return err
	}

	ir.logger.Infof("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	return ir.discord.Close()
}

// onReady should prepare metrics collector and setup web interface for configuration (features).
func (ir *Iris) onReady(s *discordgo.Session, event *discordgo.Ready) {
	numGuilds := int64(len(s.State.Guilds))
	ir.logger.Infof("Iris connected and ready - userName: %s#%s numGuilds: %d", event.User.Username, event.User.Discriminator, numGuilds)
	// should include metrics collection down here
}

// onMessageReceived will be called everytime a new message is created on any channel that the bot is listenning to.
// It will dispatch know commands to command handlers, passing along necessary info.
func (ir *Iris) onMessageReceived(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore message created by the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	msg := m.Content

	cmdPrefixLen := len(ir.Config.CmdPrefix)

	// dispatch the command iff we have our prefix, (case-insensitive) otherwise throws an errors
	if len(msg) > cmdPrefixLen && strings.EqualFold(ir.Config.CmdPrefix, msg[0:cmdPrefixLen]) {
		afterPrefix := msg[cmdPrefixLen:]
		cmd := strings.SplitN(afterPrefix, " ", 2)

		if f, ok := ir.cmdHandlers[strings.ToLower(cmd[0])]; ok {
			rest := ""
			if len(cmd) > 1 {
				rest = cmd[1]
			}

			if f.handler != nil {
				f.handler(s, m, rest)
			} else {
				s.ChannelMessageSend(m.ChannelID, "Command error/not supported - dm **@aarnphm**")
			}
		}
	}
}

// onPomEnded handles when Pom ends. It should add new users to current mongoDB if users hasn't existed in the database, else updates the minutes studied.
// this should be refactored into multiple functions (future ref)
func (ir *Iris) onPomEnded(notif db.NotifyInfo, completed bool) {
	var (
		err        error
		hash       string
		toMention  []string
		notifTitle string
		notifDesc  string
	)
	user, er := ir.discord.User(notif.User.DiscordID)
	if er == nil {
		toMention = append(toMention, user.Mention())
	}

	if completed {
		// update users' progress to databse
		err = db.FetchUser(notif.User.DiscordID)
		if err != nil {
			// create new users entry
			hash, err = db.NewUser(notif.User.DiscordID, notif.User.DiscordTag, notif.User.GuidID, pomDuration.String())
			ir.logger.Infof("inserted %s to mongoDB. Hash: %s", notif.User.DiscordID, hash)
			if err != nil {
				ir.logger.Warnf(err.Error())
			}
		} else {
			// users already in database, just updates timing
			err = db.UpdateUser(notif.User.DiscordID, int(pomDuration.Minutes()))
			if err != nil {
				ir.logger.Warnf(err.Error())
			}
		}
		// notif title
		notifTitle = "Pomodoro"

		notifDesc = ":timer: Work cycle complete. :timer:\n :blush: Time to take a break! :blush:"

		message := ""

		if len(toMention) > 0 {
			mentions := strings.Join(toMention, " ")
			message = fmt.Sprintf("%s\n%s", message, mentions)
		}

		embed := &discordgo.MessageEmbed{
			Type:        "rich",
			Title:       notifTitle,
			Color:       130,
			Description: notifDesc,
		}

		data := &discordgo.MessageSend{
			Content: message,
			Embed:   embed,
		}

		ir.discord.ChannelMessageSendComplex(notif.User.ChannelID, data)
	} else {
		ir.discord.ChannelMessageSend(notif.User.ChannelID, fmt.Sprintf("%s, pom canceled!", user.Mention()))
	}

	// ir.metrics.RecordRunningPoms(int64(ir.poms.Count()))
}

func (ir *Iris) onCmdStartPom(s *discordgo.Session, m *discordgo.MessageCreate, ex string) {
	// number regex to ensure no other entities are allowed
	rgx := regexp.MustCompile(`^([0-9])+`)

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		ir.logger.Fatal(err)
	}

	// make sure that this converts to time instead of any other funky usecase
	// ex here are time period for pom sessions
	if rgx.MatchString(ex) {
		ex = strings.ReplaceAll(ex, "`", "")
		ex = strings.TrimSpace(ex)
	} else {
		ir.logger.Warnf(fmt.Sprintf("unknown time format. Accepts numbers only. got %s instead", ex))
	}

	if ex != "" {
		newDuration, _ := strconv.Atoi(ex)
		pomDuration = time.Minute * time.Duration(newDuration)
	} else {
		pomDuration = defaultPomDuration
	}

	notif := db.NotifyInfo{
		TitleID: ex,
		User: &db.User{
			DiscordID:  m.Author.ID,
			DiscordTag: m.Author.Discriminator,
			GuidID:     channel.GuildID,
			ChannelID:  m.ChannelID,
		},
	}

	if ir.poms.CreateIfEmpty(pomDuration, ir.onPomEnded, notif) {
		// notif title
		var (
			notifTitle string
			notifDesc  string
		)
		notifTitle = "Pomodoro"

		notifDesc = fmt.Sprintf(":peach: Your timer is set to **%d minutes** :peach:\n :blush: Happy working :blush:", int(pomDuration.Minutes()))

		content := fmt.Sprintf("%s\n", m.Author.Mention())

		embed := &discordgo.MessageEmbed{
			Type:        "rich",
			Title:       notifTitle,
			Color:       130,
			Description: notifDesc,
		}

		data := &discordgo.MessageSend{
			Content: content,
			Embed:   embed,
		}
		s.ChannelMessageSendComplex(m.ChannelID, data)
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("A pomodoro is already running for %s", m.Author.Mention()))
	}
}

func (ir *Iris) onCmdCancelPom(s *discordgo.Session, m *discordgo.MessageCreate, ex string) {
	if exists := ir.poms.RemoveIfExists(m.Author.ID); !exists {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No pom is currently running for %s", m.Author.Mention()))
	}
	// if this removal is success then call onPomEnded
}

func (ir *Iris) onCmdHelp(s *discordgo.Session, m *discordgo.MessageCreate, ex string) {
	s.ChannelMessageSend(m.ChannelID, ir.helpMessage)
}

func (ir *Iris) onCmdInvite(s *discordgo.Session, m *discordgo.MessageCreate, ex string) {
	s.ChannelMessageSend(m.ChannelID, ir.inviteMessage)
}
