package cmd

import (
	"com.schumann-it.go-ieftool/internal"
	"github.com/spf13/cobra"
)

var build = &cobra.Command{
	Use:   "build",
	Short: "Build",
	Long:  `Build source policies and replacing template variables for given environments.`,
	Run: func(cmd *cobra.Command, args []string) {
		e := internal.MustNewEnvironmentsFromFlags(cmd.Flags())
		sd := internal.MustAbsPathFromFlag(cmd.Flags(), "source")
		dd := internal.MustAbsPathFromFlag(cmd.Flags(), "destination")

		log.Infof("building enviornment(s): %s", e.String())
		err := e.Build(sd, dd)
		if err != nil {
			log.Fatalf("errors occurred during build process: \n%s", err.Error())
		}
	},
}

func init() {
	globalFlags(build)
	build.Flags().StringP("source", "s", "./src", "Source directory")
	build.Flags().StringP("destination", "d", "./build", "Destination directory")
	rootCmd.AddCommand(build)
}
