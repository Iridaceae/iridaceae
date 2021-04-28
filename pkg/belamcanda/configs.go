package belamcanda

import (
	"fmt"
	"strings"

	"github.com/TensRoses/iris/internal/configparser"
)

// TODO: add options to check for valid configuration name.
var (
	ClientID      = configparser.Register("iris.clientid", "ClientID of the bot", nil)
	ClientSecrets = configparser.Register("iris.clientsecret", "ClientSecret of the bot", nil)
	BotToken      = configparser.Register("iris.authtoken", "authentication token of the bot", nil)

	Loaded = false
)

// LoadGeneralConfig loads required configs
func LoadGeneralConfig() error {
	if Loaded {
		return nil
	}

	Loaded = true
	configparser.AddSource(&configparser.EnvSource{})
	configparser.Load()

	required := []*configparser.Options{
		ClientID,
		ClientSecrets,
		BotToken,
	}

	for _, v := range required {
		if v.LoadedValue == nil {
			env := strings.ToUpper(strings.ReplaceAll(v.Name, ".", "_"))
			return fmt.Errorf("didn't contain required config options: %q (%s as envars)", v.Name, env)
		}
	}
	return nil
}
