package cmd

import (
	"context"
	"github.com/721945/dlaw-backend/bootstrap"
	"github.com/721945/dlaw-backend/libs"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func WrapCommand(opt fx.Option) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "dlaw-backend",
		Short: "A brief description of your application",
		Long:  `A longer description that spans multiple lines and likely contains`,
		Run: func(cmd *cobra.Command, args []string) {

			logger := libs.NewLogger()
			ctx := context.Background()

			app := fx.New(
				opt,
				fx.Invoke(RunInit),
			)

			err := app.Start(ctx)

			defer app.Stop(ctx)

			if err != nil {
				logger.Fatal(err)
			}

		},
	}

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	rootCmd := WrapCommand(bootstrap.CommonModules)

	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}

}

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println(err)
		log.Println("Error loading .env file")
	}

}
