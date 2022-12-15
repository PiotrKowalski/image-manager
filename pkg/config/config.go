package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

var (
	loaded = false

	ErrIncorrectInt    = errors.New("incorrect int given")
	ErrIncorrectString = errors.New("incorrect string given")
)

func init() {
	load()
}

// ReadEnvInt
func ReadEnvInt(key string) (int, error) {

	ok := viper.IsSet(key)
	if !ok {
		return 0, nil
	}

	value := viper.GetInt(key)
	if value == 0 {
		return 0, fmt.Errorf("key is 0: [%w]", ErrIncorrectInt)
	}

	return value, nil
}

// ReadEnvInt
func ReadEnvString(key string) (string, error) {

	ok := viper.IsSet(key)
	if !ok {
		return "", nil
	}

	value := viper.GetString(key)
	if value == "" {
		return "", fmt.Errorf("key is empty: [%w]", ErrIncorrectString)
	}

	return value, nil
}

func load() {
	viper.AutomaticEnv()
	loaded = true
}
