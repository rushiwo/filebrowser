package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configCatCmd)
}

var configCatCmd = &cobra.Command{
	Use:   "cat",
	Short: "Prints the configuration",
	Long:  `Prints the configuration.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		st, s := getStorageSettings()
		auther, err := st.Auth.Get(s.AuthMethod)
		checkErr(err)
		printSettings(s, auther)
	},
}
