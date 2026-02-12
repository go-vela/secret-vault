// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"log"
	"net/mail"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := &cli.Command{
		Name:      "secret-vault",
		Usage:     "Vela Vault secret plugin for sourcing secrets into pipelines",
		Copyright: "Copyright 2020 Target Brands, Inc. All rights reserved.",
		Authors: []any{
			&mail.Address{
				Name:    "Vela Admins",
				Address: "vela@target.com",
			},
		},
		Action: run,
		Flags:  flags(),
	}

	// Plugin Start
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
func run(_ context.Context, c *cli.Command) error {
	// set the log level for the plugin
	switch c.String("log.level") {
	case "t", "trace", "Trace", "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
	case "d", "debug", "Debug", "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "w", "warn", "Warn", "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "e", "error", "Error", "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "f", "fatal", "Fatal", "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "p", "panic", "Panic", "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	case "i", "info", "Info", "INFO":
		fallthrough
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.WithFields(logrus.Fields{
		"code":     "https://github.com/go-vela/secret-vault",
		"docs":     "https://go-vela.github.io/docs/plugins/registry/secret/vault/",
		"registry": "https://hub.docker.com/r/target/secret-vela",
	}).Info("Vela Secret Vault Plugin")

	// setup plugin
	p := Plugin{
		Config: &Config{
			Addr:       c.String("config.addr"),
			AuthMethod: c.String("config.auth-method"),
			Password:   c.String("config.password"),
			Token:      c.String("config.token"),
			Username:   c.String("config.username"),
		},
		Read: &Read{
			RawItems:    c.String("items"),
			OutputsPath: c.String("vela.masked-outputs"),
		},
	}

	// validate the plugin configuration
	err := p.Validate()
	if err != nil {
		return err
	}

	// execute the plugin
	return p.Exec()
}
