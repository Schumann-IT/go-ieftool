package cmd

import (
	"com.schumann-it.go-ieftool/internal"
	"github.com/spf13/cobra"
)

var remove = &cobra.Command{
	Use:   "remove",
	Short: "Delete remote b2c policies.",
	Long:  `Delete remote b2c policies from B2C identity experience framework.`,
	Run: func(cmd *cobra.Command, args []string) {
		e := internal.MustNewEnvironmentsFromFlags(cmd.Flags())

		log.Infof("deleting policies for enviornment(s): %s", e.String())

		err := e.DeleteRemotePolicies()
		if err != nil {
			log.Fatalf("failed to delete policies %s", err.Error())
		}
	},
}

func init() {
	globalFlags(remove)
	rootCmd.AddCommand(remove)
}
