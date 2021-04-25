package config

import "github.com/spf13/viper"

// Configs holds struct for the bot configuration data.
type Configs struct {
	// for future references
	CmdPrefix string `yaml:"cmdPrefix"`
}

// Secrets is the Bot's per-user data, some of which is secret.
type Secrets struct {
	AuthToken string `yaml:"authToken"`
	ClientID  string `yaml:"clientID"`
}

// LoadConfigFile will setup config from given path, will use viper to process and returns error if one occurred.
func LoadConfigFile(path string) (*Configs, error) {
	viper.SetConfigName("default")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	cfg := &Configs{
		CmdPrefix: viper.GetString("cmdPrefix"),
	}
	return cfg, err
}

// LoadSecretsFile will setup config from given path, will use viper to process and returns error if one occurred.
func LoadSecretsFile(path string) (*Secrets, error) {
	viper.SetConfigName("discord")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	// https://stackoverflow.com/a/47185439/8643197
	// err := viper.MergeInConfig()
	err := viper.ReadInConfig()
	sec := &Secrets{
		AuthToken: viper.GetString("authToken"),
		ClientID:  viper.GetString("clientID"),
	}
	return sec, err
}
