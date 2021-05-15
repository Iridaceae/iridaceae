package helpers

import (
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/joho/godotenv"

	"github.com/Iridaceae/iridaceae/pkg/configmanager"
)

var (
	ConcertinaClientID, _      = configmanager.RegisterOption("concertina.clientid", "ClientID of test bot", nil)
	ConcertinaClientSecrets, _ = configmanager.RegisterOption("concertina.clientsecret", "ClientSecret of test bot", nil)
	ConcertinaBotToken, _      = configmanager.RegisterOption("concertina.authtoken", "authentication token of test bot", nil)
	IridaceaeClientID, _       = configmanager.RegisterOption("iris.clientid", "IridaceaeClientID of the bot", nil)
	IridaceaeClientSecrets, _  = configmanager.RegisterOption("iris.clientsecret", "ClientSecret of the bot", nil)
	IridaceaeBotToken, _       = configmanager.RegisterOption("iris.authtoken", "authentication token of the bot", nil)

	CmdPrefix, _ = configmanager.RegisterOption("iris.cmdprefix", "prefix for iris", "ir!")

	Loaded     = false
	CI         = true
	configroot = "config"
)

// GetBotToken will handles authToken.
func GetBotToken(token configmanager.Options) string {
	tokenStr := token.ToString()
	if !strings.HasSuffix(tokenStr, "Bot ") {
		tokenStr = "Bot " + tokenStr
	}
	return tokenStr
}

// GetConfigRoot returns our configparser dir.
func GetConfigRoot() string {
	rootDir, _ := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	trimmed := strings.ReplaceAll(string(rootDir), "\n", "")
	return path.Join(trimmed, configroot)
}

// LoadConfig will load given client id, secrets, and token for setting bot.
func LoadConfig(cid, cs, token configmanager.Options) error {
	if Loaded {
		return nil
	}

	Loaded = true
	configmanager.RegisterSource(&configmanager.EnvSource{})
	configmanager.LoadOptions()

	required := []configmanager.Options{cid, cs, token}

	for _, v := range required {
		if v.GetValue() == nil {
			env := strings.ToUpper(strings.ReplaceAll(v.GetName(), ".", "_"))
			return fmt.Errorf("didn't contain required configparser options: %q (%s as envars)", v.GetName(), env)
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
	return godotenv.Load(path.Join(GetConfigRoot(), "defaults.env"))
}
