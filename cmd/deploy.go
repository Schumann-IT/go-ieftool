package cmd

import (
	"com.schumann-it.go-ieftool/internal"
	"github.com/spf13/cobra"
)

var deploy = &cobra.Command{
	Use:   "deploy [path to policy]",
	Short: "Deploy b2c policy.",
	Long:  `Deploy b2c policy to B2C identity experience framework.`,
	Run: func(cmd *cobra.Command, args []string) {
		e := internal.MustNewEnvironmentsFromFlags(cmd.Flags())
		bd := internal.MustAbsPathFromFlag(cmd.Flags(), "build-dir")

		log.Infof("deploying enviornment(s): %s", e.String())
		err := e.Deploy(bd)
		if err != nil {
			log.Fatalf("failed to deploy policy %v", err)
		}
	},
}

func init() {
	globalFlags(deploy)
	deploy.Flags().StringP("build-dir", "b", "./build", "Build directory")
	rootCmd.AddCommand(deploy)
}
