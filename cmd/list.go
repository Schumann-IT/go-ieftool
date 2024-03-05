package cmd

import (
	"com.schumann-it.go-ieftool/internal"
	"github.com/spf13/cobra"
)

var list = &cobra.Command{
	Use:   "list [path to policy]",
	Short: "List remote b2c policy.",
	Long:  `List remote b2c policy from B2C identity experience framework.`,
	Run: func(cmd *cobra.Command, args []string) {
		e := internal.MustNewEnvironmentsFromFlags(cmd.Flags())

		log.Infof("fetching list of policy for enviornment(s): %s", e.String())

		ps, err := e.FetchRemotePolicies()
		if err != nil {
			log.Fatalf("failed to fetch policy %v", err)
		}

		for n, l := range ps {
			log.Error(n)
			for _, i := range l {
				log.Error(i)
			}
		}
	},
}

func init() {
	globalFlags(list)
	rootCmd.AddCommand(list)
}
