package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Magmar ...
var Magmar *ViperConfig

// ViperConfig ...
type ViperConfig struct {
	*viper.Viper
}

func init() {
	Magmar = initViperConfig()
}

func initViperConfig() *ViperConfig {
	v := viper.New()

	var env *string
	if value := os.Getenv("env"); value != "" {
		env = &value
	} else {
		env = pflag.String("env", "local", "help message for environment")
	}

	pflag.Parse()
	v.BindPFlags(pflag.CommandLine)

	v.SetConfigName(*env)

	v.SetConfigType("yml")
	v.AddConfigPath("./config/")
	v.AddConfigPath("../config/")
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		fmt.Printf("Error when reading config: %v\n", err)
		os.Exit(1)
	}

	if v.GetString("env") == "local" {
		v.Set("absPath", getRootDir())
	}

	prvTokenKey, err := ioutil.ReadFile("key/token_key")
	if err != nil {
		fmt.Printf("Error when reading prvTokenKey: %v\n", err)
		os.Exit(1)
	}
	v.Set("prvTokenKey", prvTokenKey)

	pubTokenKey, err := ioutil.ReadFile("key/token_key.pub")
	if err != nil {
		fmt.Printf("Error when reading pubTokenKey: %v\n", err)
		os.Exit(1)
	}
	v.Set("pubTokenKey", pubTokenKey)

	return &ViperConfig{v}
}

func getRootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
