package cmd

import (
	"com.schumann-it.go-ieftool/internal"
	"github.com/spf13/cobra"
)

var deleteKeySets = &cobra.Command{
	Use:   "delete-key-sets",
	Short: "Delete Key Sets",
	Long:  `Delete Key Sets for Policies.`,
	Run: func(cmd *cobra.Command, args []string) {
		e := internal.MustNewEnvironmentsFromFlags(cmd.Flags())

		log.Infof("deleting keysets for enviornment(s): %s", e.String())
		err := e.DeleteKeySets()
		if err != nil {
			log.Fatalf("errors occurred during key set deletion process: \n%s", err.Error())
		}
	},
}

func init() {
	globalFlags(deleteKeySets)
	rootCmd.AddCommand(deleteKeySets)
}
