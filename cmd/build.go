package cmd

import (
	"com.schumann-it.go-ieftool/internal"
	"com.schumann-it.go-ieftool/pkg/b2c"
	"github.com/spf13/cobra"
)

var build = &cobra.Command{
	Use:   "build",
	Short: "Build",
	Long:  `Build source policy and replacing template variables for given environments.`,
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

var buildNew = &cobra.Command{
	Use:   "buildNew",
	Short: "Build",
	Long:  `Build source policy and replacing template variables for given environments.`,
	Run: func(cmd *cobra.Command, args []string) {
		cf, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalf("could not parse flag 'config': \n%s", err.Error())
		}

		en, err := cmd.Flags().GetString("environment")
		if err != nil {
			log.Fatalf("could not parse flag 'environment': \n%s", err.Error())
		}
		if en == "" {
			log.Fatalf("must provide flag 'environment'")
		}

		sd := internal.MustAbsPathFromFlag(cmd.Flags(), "source")
		dd := internal.MustAbsPathFromFlag(cmd.Flags(), "destination")

		a, err := b2c.NewApi(cf, sd, dd)
		if err != nil {
			log.Fatalf("could build environment %s: %s", en, err.Error())
		}
		err = a.BuildPolicies(en)
		if err != nil {
			log.Fatalf("could build environment %s: %s", en, err.Error())
		}
	},
}

func init() {
	globalFlags(build)
	build.Flags().StringP("source", "s", "./src", "File directory")
	build.Flags().StringP("destination", "d", "./build", "Destination directory")
	rootCmd.AddCommand(build)

	globalFlags(buildNew)
	buildNew.Flags().StringP("source", "s", "./src", "File directory")
	buildNew.Flags().StringP("destination", "d", "./build", "Destination directory")
	rootCmd.AddCommand(buildNew)
}
