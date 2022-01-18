// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := cli.NewApp()

	// Plugin Information
	app.Name = "secret-vault"
	app.HelpName = "secret-vault"
	app.Usage = "Vela Vault secret plugin for sourcing secrets into pipelines"
	app.Copyright = "Copyright (c) 2022 Target Brands, Inc. All rights reserved."
	app.Authors = []*cli.Author{
		{
			Name:  "Vela Admins",
			Email: "vela@target.com",
		},
	}

	// Plugin Metadata
	app.Compiled = time.Now()
	app.Action = run

	// Plugin Flags
	app.Flags = flags()

	// Plugin Start
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
func run(c *cli.Context) error {
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
		"code": "https://github.com/go-vela/secret-vault",

		// TODO think about a place for secret plugin docs
		// "docs":     "",
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
			RawItems: c.String("items"),
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
