package cli

import (
	"apitraning/pkg/integrations/amoCRM"
	"apitraning/pkg/repository"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func Workerfunc(repo *repository.Repository) {
	app := &cli.App{
		Name:  "prod",
		Usage: "work",
		Commands: []*cli.Command{
			{
				Name:  "start",
				Usage: "Start the webhook worker",
				Action: func(c *cli.Context) error {
					amoCRM.WebhookWorker(repo)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
