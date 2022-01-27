package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	// setup logging configuration
	setupLogging()

	cfg := &Configuration{}

	app := &cli.App{
		Name:  "amalgam",
		Usage: "Create macOS Universal binaries from Github releases",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "owner",
				Required:    true,
				Usage:       "Github repo owner username ",
				Destination: &cfg.Owner,
			},
			&cli.StringFlag{
				Name:        "repo",
				Required:    true,
				Usage:       "Github repository name",
				Destination: &cfg.Repository,
			},
			&cli.StringFlag{
				Name:        "tag",
				Required:    false,
				Usage:       "Github repository name",
				Value:       "latest",
				Destination: &cfg.Tag,
			},
			&cli.StringFlag{
				Name:        "amd64",
				Required:    true,
				Usage:       "Substring for the amd64 binary",
				Destination: &cfg.Amd64Substring,
			},
			&cli.StringFlag{
				Name:        "arm64",
				Required:    true,
				Usage:       "Substring for the arm64 binary",
				Destination: &cfg.Arm64Substring,
			},
			&cli.BoolFlag{
				Name:        "compressed",
				Required:    false,
				Usage:       "Do the releases use compressed archives",
				Destination: &cfg.Compressed,
			},
			&cli.BoolFlag{
				Name:        "overwrite",
				Required:    false,
				Usage:       "Delete pre-existing universal asset?",
				Destination: &cfg.Overwrite,
			},
			&cli.StringFlag{
				Name:        "identifier",
				Required:    false,
				Usage:       "Identifier for universal binary",
				Value:       "all",
				Destination: &cfg.UniversalIdentifer,
			},
		},
		Action: func(context *cli.Context) error {
			log.Debug("Num flags: ", context.NumFlags())
			if context.NumFlags() < 4 {
				cli.ShowAppHelpAndExit(context, 1)
			}

			token := os.Getenv("GITHUB_TOKEN")
			if token == "" {
				log.Fatal("GITHUB_TOKEN environment variable is empty.")
				os.Exit(2)
			}

			cfg.GithubToken = token

			err := CreateUniveralBinary(cfg)
			if err != nil {
				log.Fatal("Exiting. Encountered exception: ", err)
				os.Exit(1)
			}

			log.Info("Complete.")
			return nil
		},
	}

	// run command line application and exit on error
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Encountered error: %s", err.Error())
		log.Fatal("Exiting...")
		os.Exit(1)
	}

}
