package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/op/go-logging.v1"
)

var log = logging.MustGetLogger("cmd")

var rootCmd = &cobra.Command{
	Use:   "ieftool",
	Short: "Tooling for Azure B2C Identity Experience Framework",
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		lvl := logging.INFO
		isDebug, _ := cmd.Flags().GetBool("debug")
		if isDebug {
			lvl = logging.DEBUG
		}
		logging.SetLevel(lvl, "")
		logging.SetFormatter(logging.MustStringFormatter(
			`%{color}%{level}(%{module})%{color:reset} %{message}`,
		))
		log = logging.MustGetLogger("cmd")
	},
}

var completion = &cobra.Command{
	Use:                "completion [bash|zsh|fish|powershell]",
	Short:              "Generate completion script",
	Long:               `To load completions.`,
	DisableFlagParsing: true,
	ValidArgs:          []string{"bash", "zsh", "fish", "powershell"},
	Args:               cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completion)
}

func globalFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("debug", false, "Enable debug mode")
	cmd.Flags().StringP("config", "c", "./config.yaml", "Path to the ieftool configuration file")
	cmd.Flags().StringP("environment", "e", "", "Environment to deploy (default: all environments)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
