package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/lonelycode/montag-ai/server/api/db"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "server",
				Value:    "https://montag.example.com",
				Usage:    "server to use for montag resources",
				EnvVars:  []string{"MONTAG_SERVER"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "key",
				Value:    "YOURAPIKEY",
				Usage:    "API key to validate against the server",
				EnvVars:  []string{"MONTAG_KEY"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "storage",
				Value:    "./montag.sqlite",
				Usage:    "database to use for value storage",
				EnvVars:  []string{"MONTAG_DB"},
				Required: true,
			},
		},
		Name:  "montag-cli",
		Usage: "run montag scripts locally, with resources provided by a montag server",
		Action: func(cCtx *cli.Context) error {
			return runsScript(cCtx.Args().Get(0))
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getDB() *gorm.DB {
	var d *gorm.DB
	binaryPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	d = db.GetSqlite(filepath.Dir(binaryPath))

	return d
}

func runsScript(script string) error {
	if script == "" {
		return fmt.Errorf("a tengo script filename is required")
	}

	// create a db
	d := getDB()

	return nil
}
