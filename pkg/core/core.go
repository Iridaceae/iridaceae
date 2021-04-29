// Package core contains bot functionality for iris. Currently following Pomodoro technique.
package core

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/TensRoses/iris/internal/datastore"
	"github.com/TensRoses/iris/internal/irislog"
)

// GetBotToken will handles authToken.
func GetBotToken() string {
	token := BotToken.GetString()
	if !strings.HasSuffix(token, "Bot ") {
		token = "Bot " + token
	}

	return token
}

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
	discord       *discordgo.Session
	logger        *irislog.IrisLogger
	cmdHandlers   map[string]botCommand
	poms          UserPomodoroMap
	// record metrics here
	// metrics metrics.Recorder
}

// NewIris creates a new instance of Iris that can deploy over Heroku.
func NewIris() *Iris {
	// setup new logLevel
	logger := irislog.NewLogger(logLevel, "botevents")

	err := LoadConfig()
	if err != nil {
		logger.Error(err)
	}

	ir := &Iris{
		logger: logger,
		poms:   NewUserPomodoroMap(),
	}

	ir.registerCmdHandlers()
	ir.helpMessage = ir.buildHelpMessage()
	ir.inviteMessage = fmt.Sprintf("Click here: <"+baseAuthURLTemplate+"> to invite me to the server", ClientID.GetString())
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
		helpBuffer.WriteString(fmt.Sprintf("   Example: `%s%s %s`\n", CmdPrefix.GetString(), cmdStr, cmd.exampleParams))
	}

	helpBuffer.WriteString("\n" + ir.inviteMessage)

	return helpBuffer.String()
}

// Start will start the bot, blocking til completed.
func (ir *Iris) Start() error {
	var err error
	ir.discord, err = discordgo.New(GetBotToken())
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

	ir.logger.Info("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	return ir.discord.Close()
}

// onReady should prepare metrics collector and setup web interface for configuration (features).
func (ir *Iris) onReady(s *discordgo.Session, event *discordgo.Ready) {
	numGuilds := int64(len(s.State.Guilds))
	ir.logger.Info(fmt.Sprintf("Iris connected and ready - userName: %s#%s numGuilds: %d", event.User.Username, event.User.Discriminator, numGuilds))
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

	cmdPrefixLen := len(CmdPrefix.GetString())

	// dispatch the command iff we have our prefix, (case-insensitive) otherwise throws an errors
	if len(msg) > cmdPrefixLen && strings.EqualFold(CmdPrefix.GetString(), msg[0:cmdPrefixLen]) {
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
				_, err := s.ChannelMessageSend(m.ChannelID, "Command error/not supported - dm **@aarnphm**")
				if err != nil {
					ir.logger.Warn(err.Error())
				}
			}
		}
	}
}

// onPomEnded handles when Pom ends. It should add new users to current mongoDB if users hasn't existed in the database, else updates the minutes studied.
// NOTE: this should be refactored into multiple functions.
func (ir *Iris) onPomEnded(notif NotifyInfo, completed bool) {
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
		err = datastore.FetchUser(notif.User.DiscordID)
		if err != nil {
			// create new users entry
			hash, err = datastore.NewUser(notif.User.DiscordID, notif.User.DiscordTag, notif.User.GUIDID, pomDuration.String())
			ir.logger.Info("inserted %s to mongoDB. Hash: %s", notif.User.DiscordID, hash)
			if err != nil {
				ir.logger.Warn(err.Error())
			}
		} else {
			// users already in database, just updates timing
			err = datastore.UpdateUser(notif.User.DiscordID, int(pomDuration.Minutes()))
			if err != nil {
				ir.logger.Warn(err.Error())
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
			Color:       msgColor,
			Description: notifDesc,
		}

		data := &discordgo.MessageSend{
			Content: message,
			Embed:   embed,
		}

		_, err := ir.discord.ChannelMessageSendComplex(notif.User.ChannelID, data)
		if err != nil {
			ir.logger.Warn(err.Error())
		}
	} else {
		_, err := ir.discord.ChannelMessageSend(notif.User.ChannelID, fmt.Sprintf("%s, pom canceled!", user.Mention()))
		if err != nil {
			ir.logger.Warn(err.Error())
		}
	}

	// ir.metrics.RecordRunningPoms(int64(ir.poms.Count()))
}

func (ir *Iris) onCmdStartPom(s *discordgo.Session, m *discordgo.MessageCreate, ex string) {
	// number regex to ensure no other entities are allowed
	rgx := regexp.MustCompile(`^([0-9])+`)

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		ir.logger.Error(err)
	}

	// make sure that this converts to time instead of any other funky usecase
	// ex here are time period for pom sessions
	if rgx.MatchString(ex) {
		ex = strings.ReplaceAll(ex, "`", "")
		ex = strings.TrimSpace(ex)
	} else {
		ir.logger.Warn(fmt.Sprintf("unknown time format. Accepts numbers only. got %s instead", ex))
	}

	if ex != "" {
		newDuration, _ := strconv.Atoi(ex)
		pomDuration = time.Minute * time.Duration(newDuration)
	} else {
		pomDuration = defaultPomDuration
	}

	notif := NotifyInfo{
		TitleID: ex,
		User: &datastore.User{
			DiscordID:  m.Author.ID,
			DiscordTag: m.Author.Discriminator,
			GUIDID:     channel.GuildID,
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
			Color:       msgColor,
			Description: notifDesc,
		}

		data := &discordgo.MessageSend{
			Content: content,
			Embed:   embed,
		}
		_, err := s.ChannelMessageSendComplex(m.ChannelID, data)
		if err != nil {
			ir.logger.Warn(err.Error())
		}
	} else {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("A pomodoro is already running for %s", m.Author.Mention()))
		if err != nil {
			ir.logger.Warn(err.Error())
		}
	}
}

func (ir *Iris) onCmdCancelPom(s *discordgo.Session, m *discordgo.MessageCreate, ex string) {
	if exists := ir.poms.RemoveIfExists(m.Author.ID); !exists {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No pom is currently running for %s", m.Author.Mention()))
		if err != nil {
			ir.logger.Warn(err.Error())
		}
	}
	// if this removal is success then call onPomEnded
}

func (ir *Iris) onCmdHelp(s *discordgo.Session, m *discordgo.MessageCreate, ex string) {
	_, err := s.ChannelMessageSend(m.ChannelID, ir.helpMessage)
	if err != nil {
		ir.logger.Warn(err.Error())
	}
}

func (ir *Iris) onCmdInvite(s *discordgo.Session, m *discordgo.MessageCreate, ex string) {
	_, err := s.ChannelMessageSend(m.ChannelID, ir.inviteMessage)
	if err != nil {
		ir.logger.Warn(err.Error())
	}
}
