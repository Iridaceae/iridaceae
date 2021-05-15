package helpers

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"

	"github.com/Iridaceae/iridaceae/pkg/configmanager"
)

var (
	ConcertinaClientID, _      = configmanager.Register("concertina.clientid", "ClientID of test bot", nil)
	ConcertinaClientSecrets, _ = configmanager.Register("concertina.clientsecret", "ClientSecret of test bot", nil)
	ConcertinaBotToken, _      = configmanager.Register("concertina.authtoken", "authentication token of test bot", nil)
	IridaceaeClientID, _       = configmanager.Register("iris.clientid", "IridaceaeClientID of the bot", nil)
	IridaceaeClientSecrets, _  = configmanager.Register("iris.clientsecret", "ClientSecret of the bot", nil)
	IridaceaeBotToken, _       = configmanager.Register("iris.authtoken", "authentication token of the bot", nil)
	CmdPrefix, _               = configmanager.Register("iris.cmdprefix", "prefix for iris", "-ir")
	Loaded                     = false
	CI                         = true
)

const (
	BaseAuthURLTemplate string = "https://discord.com/api/oauth2/authorize?client_id=%s&scope=bot"
)

// GetBotToken will handles authToken.
func GetBotToken(token *configmanager.Options) string {
	tokenStr := token.GetString()
	if !strings.HasSuffix(tokenStr, "Bot ") {
		tokenStr = "Bot " + tokenStr
	}
	return tokenStr
}

// GetConfigRoot returns our config dir.
func GetConfigRoot() string {
	rootDir, _ := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	trimmed := strings.ReplaceAll(string(rootDir), "\n", "")
	return strings.Join([]string{trimmed, "config"}, "/")
}

// LoadConfig will load given client id, secrets, and token for setting bot.
func LoadConfig(cid, cs, token *configmanager.Options) error {
	if Loaded {
		return nil
	}

	Loaded = true
	configmanager.AddSource(&configmanager.EnvSource{})
	configmanager.Load()

	required := []*configmanager.Options{cid, cs, token}

	for _, v := range required {
		if v.LoadedValue == nil {
			env := strings.ToUpper(strings.ReplaceAll(v.Name, ".", "_"))
			return fmt.Errorf("didn't contain required config options: %q (%s as envars)", v.Name, env)
		}
	}
	return nil
}

// LoadGlobalEnv ensures our envars are loaded.
func LoadGlobalEnv() error {
	// We assume that everything is run with CI, thus the usecase of this is when not running with CI.
	CI = false
	if CI {
		return nil
	}
	return godotenv.Load(strings.Join([]string{GetConfigRoot(), "defaults.env"}, "/"))
}
