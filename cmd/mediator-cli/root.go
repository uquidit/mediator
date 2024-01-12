package main

import (
	"fmt"
	"os"
	"uqtu/mediator/apiclient"
	"uqtu/mediator/clicommands"
	"uqtu/mediator/clicommands/securechangeapi"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	URL                string
	InsecureSkipVerify bool = false

	rootCmd = &cobra.Command{
		Use:   "mediator",
		Short: "A low-level CLI to manage Mediator back-end",
		Long: `This CLI provides commands to manage scripts used by Mediator back-end. It includes:
* List registered scripts
* Register new scripts
* Un-register useless scripts
* Refresh script checksum
* Test scripts
`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if URL == "" {
				return fmt.Errorf("provided Back-End URL is empty")
			}
			apiclient.InitHelpers(URL, InsecureSkipVerify)
			return nil
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&InsecureSkipVerify, "sslskipverify", "", false, "Skip SSL certificate verification (insecure)")
	rootCmd.PersistentFlags().StringVarP(&URL, "url", "u", "", "Back-end URL (required)")
	rootCmd.MarkPersistentFlagRequired("url")

	rootCmd.AddCommand(clicommands.MediatorSettingsCmd)
	rootCmd.AddCommand(securechangeapi.MediatorSecurechangeAPICmd)
	rootCmd.AddCommand(clicommands.ScriptCmd)
}
