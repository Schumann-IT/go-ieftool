package cmd

import (
	"com.schumann-it.go-ieftool/internal"
	"github.com/spf13/cobra"
)

var list = &cobra.Command{
	Use:   "list [path to policies]",
	Short: "List remote b2c policies.",
	Long:  `List remote b2c policies from B2C identity experience framework.`,
	Run: func(cmd *cobra.Command, args []string) {
		e := internal.MustNewEnvironmentsFromFlags(cmd.Flags())

		log.Infof("fetching list of policies for enviornment(s): %s", e.String())

		ps, err := e.FetchRemotePolicies()
		if err != nil {
			log.Fatalf("failed to fetch policies %v", err)
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
