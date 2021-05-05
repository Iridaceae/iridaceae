package pkg

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/Iridaceae/iridaceae/internal/configparser"
)

// TODO: add options to check for valid configuration name.
var (
	ConcertinaClientID, _      = configparser.Register("concertina.clientid", "ClientID of test bot", nil)
	ConcertinaClientSecrets, _ = configparser.Register("concertina.clientsecret", "ClientSecret of test bot", nil)
	ConcertinaBotToken, _      = configparser.Register("concertina.authtoken", "authentication token of test bot", nil)
	IridaceaeClientID, _       = configparser.Register("iris.clientid", "IridaceaeClientID of the bot", nil)
	IridaceaeClientSecrets, _  = configparser.Register("iris.clientsecret", "ClientSecret of the bot", nil)
	IridaceaeBotToken, _       = configparser.Register("iris.authtoken", "authentication token of the bot", nil)
	CmdPrefix, _               = configparser.Register("iris.cmdprefix", "prefix for iris", "/i ")
	Loaded                     = false

	// VERSION is defined via git.
	VERSION     = "unknown"
	IrisSession *discordgo.Session
	IrisUser    *discordgo.User
)

const (
	BaseAuthURLTemplate string = "https://discord.com/api/oauth2/authorize?client_id=%s&scope=bot"
)

// GetBotToken will handles authToken.
func GetBotToken(token *configparser.Options) string {
	tokenStr := token.GetString()
	if !strings.HasSuffix(tokenStr, "Bot ") {
		tokenStr = "Bot " + tokenStr
	}
	return tokenStr
}

// LoadConfig will load given clientid, secrets, and token for setting bot.
func LoadConfig(clientid, clientsecret, token *configparser.Options) error {
	if Loaded {
		return nil
	}

	Loaded = true
	configparser.AddSource(&configparser.EnvSource{})
	configparser.Load()

	required := []*configparser.Options{
		clientid,
		clientsecret,
		token,
	}

	for _, v := range required {
		if v.LoadedValue == nil {
			env := strings.ToUpper(strings.ReplaceAll(v.Name, ".", "_"))
			return fmt.Errorf("didn't contain required config options: %q (%s as envars)", v.Name, env)
		}
	}
	return nil
}
