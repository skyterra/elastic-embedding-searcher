package cmd

import (
	"github.com/skyterra/elastic-embedding-searcher/api"
	"github.com/skyterra/elastic-embedding-searcher/elastic"
	runner "github.com/skyterra/elastic-embedding-searcher/runner"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "elastic-embedding-searcher",
	Short: "An advanced semantic search engine.",

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// skip execution if the command is to show help
		if cmd.Name() == "help" || cmd.Name() == "h" {
			return nil
		}

		addr, err := cmd.Flags().GetString("elastic-address")
		if err != nil {
			return err
		}

		modelPath, err := cmd.Flags().GetString("model-path")
		if err != nil {
			return err
		}

		// start ModelX process.
		if err = runner.StartModelX(3, "modelx/server.py", modelPath); err != nil {
			return err
		}

		// init elastic search without username and password.
		if err = elastic.Init(addr, "", ""); err != nil {
			return err
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		// start service.
		port, err := cmd.Flags().GetInt16("port")
		if err != nil {
			return err
		}

		api.Start(port)
		return nil
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if err := api.Cleanup(); err != nil {
			log.Printf("failed to exec api cleanup. err:%s\n", err.Error())
		}

		if err := elastic.Cleanup(); err != nil {
			log.Printf("failed to exec elastic cleanup. err:%s\n", err.Error())
		}

		err := runner.StopModelX()
		if err != nil {
			log.Printf("failed to stop modelX. err:%s\n", err.Error())
		}

	},
}

func init() {
	rootCmd.PersistentFlags().StringP("elastic-address", "e", "", "Set elasticsearch address.")
	rootCmd.PersistentFlags().StringP("model-path", "m", "", "No set to load model from hugging-face website. set local director to load model from disk.")
	rootCmd.PersistentFlags().Int16P("port", "p", 8081, "Set grpc port. default port is 8081")
}

// Execute run the root command.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}
