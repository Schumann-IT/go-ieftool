package cmd

import (
	"com.schumann-it.go-ieftool/internal"
	"github.com/spf13/cobra"
)

var createKeySets = &cobra.Command{
	Use:   "create-key-sets",
	Short: "Create Key Sets",
	Long:  `Create Key Sets for Policies.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		e := internal.MustNewEnvironmentsFromFlags(cmd.Flags())

		log.Infof("creating keysets for enviornment(s): %s", e.String())
		err := e.CreateKeySets()
		if err != nil {
			log.Fatalf("errors occurred during key set creation process: \n%s", err.Error())
		}

		return nil
	},
}

func init() {
	globalFlags(createKeySets)
	rootCmd.AddCommand(createKeySets)
}
