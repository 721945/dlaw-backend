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

// rootCmd represents the base command when called without any subcommands
func WrapCommand(opt fx.Option) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "dlaw-backend",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			//isMigrate := cmd.Flag("migrate").Value.String() == "true"

			logger := libs.NewLogger()
			ctx := context.Background()
			app := fx.New(
				opt,
				fx.Invoke(RunInit),
			)
			//log.Fatal("ERROR HERE")

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
	err := godotenv.Load(".env.asb")

	if err != nil {
		log.Println("Error loading .env file")
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dlaw-backend.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("migrate", "t", false, "Help message for migrate")
}
