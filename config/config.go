package config

import (
	"bytes"
	"context"
	"fmt"
	"magmar/util"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3Config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	if *env == "local" {
		v.AddConfigPath("./config/")
		v.AddConfigPath("../config/")

		err := v.ReadInConfig()
		if err != nil {
			fmt.Printf("Error when reading config: %v\n", err)
			os.Exit(1)
		}
	} else {
		buf, err := getConfig(context.TODO(), *env)
		if err != nil {
			fmt.Printf("Error when reading config from s3: %v\n", err)
			os.Exit(1)
		}

		err = v.ReadConfig(buf)
		if err != nil {
			fmt.Printf("Error when reading config from s3: %v\n", err)
			os.Exit(1)
		}
	}

	v.AutomaticEnv()

	return &ViperConfig{v}
}

func getRootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

func getConfig(ctx context.Context, env string) (*bytes.Buffer, error) {
	cfg, err := s3Config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println("getConfig error", err)
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	resp, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(util.ConfigBucketName),
		Key:    aws.String(fmt.Sprintf("%s.yml", env)),
	})

	if err != nil {
		fmt.Println("getConfig error", err)
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		fmt.Println("getConfig error", err)
		return nil, err
	}

	return buf, nil
}
