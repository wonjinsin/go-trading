package config

import (
	"bytes"
	"context"
	"fmt"
	"magmar/util"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	s3Config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ViperConfig ...
type ViperConfig struct {
	*viper.Viper
}

// InitViperConfig ...
func InitViperConfig() *ViperConfig {
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
			fmt.Printf("Error when get config from s3: %v\n", err)
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
	fmt.Println("getConfig start", "env", env)
	cfg, err := s3Config.LoadDefaultConfig(ctx,
		s3Config.WithRetryer(func() aws.Retryer {
			return retry.NewStandard(func(o *retry.StandardOptions) {
				o.MaxAttempts = 5
			})
		}),
		s3Config.WithHTTPClient(
			&http.Client{
				Timeout: 30 * time.Second, // HTTP 타임아웃 설정
			}))
	if err != nil {
		fmt.Println("getConfig connect error", err)
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	resp, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(util.ConfigBucketName),
		Key:    aws.String(fmt.Sprintf("%s.yml", env)),
	})

	if err != nil {
		fmt.Println("getConfig getObject error", err)
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		fmt.Println("readFrom buffer error", err)
		return nil, err
	}

	return buf, nil
}
