package pkg

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"

	configparser "github.com/Iridaceae/iridaceae/pkg/configmanager"
)

var (
	ConcertinaClientID, _      = configparser.Register("concertina.clientid", "ClientID of test bot", nil)
	ConcertinaClientSecrets, _ = configparser.Register("concertina.clientsecret", "ClientSecret of test bot", nil)
	ConcertinaBotToken, _      = configparser.Register("concertina.authtoken", "authentication token of test bot", nil)
	IridaceaeClientID, _       = configparser.Register("iris.clientid", "IridaceaeClientID of the bot", nil)
	IridaceaeClientSecrets, _  = configparser.Register("iris.clientsecret", "ClientSecret of the bot", nil)
	IridaceaeBotToken, _       = configparser.Register("iris.authtoken", "authentication token of the bot", nil)
	CmdPrefix, _               = configparser.Register("iris.cmdprefix", "prefix for iris", "-ir ")
	Loaded                     = false
	CI                         = true
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

// GetRootDir will returns root dir of iridaceae.
func GetRootDir() string {
	rootDir, _ := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	return strings.ReplaceAll(string(rootDir), "\n", "")
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

func LoadGlobalEnv() error {
	// We assume that everything is run with CI, thus the usecase of this is when not running with CI.
	CI = false
	if CI {
		return nil
	}
	return godotenv.Load(strings.Join([]string{GetRootDir(), "defaults.env"}, "/"))
}
