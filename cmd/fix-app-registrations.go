package cmd

import (
	"com.schumann-it.go-ieftool/internal"
	"github.com/spf13/cobra"
)

var fix = &cobra.Command{
	Use:   "fix-app-registrations",
	Short: "Fix App Registrations",
	Long:  `Fix App Registration manifest.`,
	Run: func(cmd *cobra.Command, args []string) {
		e := internal.MustNewEnvironmentsFromFlags(cmd.Flags())

		log.Infof("fixing app registrations for enviornment(s): %s", e.String())
		err := e.FixAppRegistrations()
		if err != nil {
			log.Fatalf("errors occurred during fix app registrations process: \n%s", err.Error())
		}
	},
}

func init() {
	globalFlags(fix)
	rootCmd.AddCommand(fix)
}
