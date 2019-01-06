package cmd

import (
	"fmt"

	"github.com/filebrowser/filebrowser/v2/rules"
	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/storage"
	"github.com/filebrowser/filebrowser/v2/users"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rulesCmd)
	rulesCmd.PersistentFlags().StringP("username", "u", "", "username of user to which the rules apply")
	rulesCmd.PersistentFlags().UintP("id", "i", 0, "id of user to which the rules apply")
}

var rulesCmd = &cobra.Command{
	Use:     "rules",
	Version: rootCmd.Version,
	Short:   "Rules management utility",
	Long: `On each subcommand you'll have available at least two flags:
"username" and "id". You must either set only one of them
or none. If you set one of them, the command will apply to
an user, otherwise it will be applied to the global set or
rules.`,
	Args: cobra.NoArgs,
}

func runRules(cmd *cobra.Command, users func(*users.User, *storage.Storage), global func(*settings.Settings, *storage.Storage)) {
	st, s := getStorageSettings()

	id := getUserIdentifier(cmd)
	if id != nil {
		user, err := st.Users.Get("", id)
		checkErr(err)

		if users != nil {
			users(user, st)
		}

		printRules(user.Rules, id)
		return
	}

	if global != nil {
		global(s, st)
	}

	printRules(s.Rules, id)
}

func getUserIdentifier(cmd *cobra.Command) interface{} {
	id := mustGetUint(cmd, "id")
	username := mustGetString(cmd, "username")

	if id != 0 {
		return id
	} else if username != "" {
		return username
	}

	return nil
}

func printRules(rules []rules.Rule, id interface{}) {
	if id == nil {
		fmt.Printf("Global Rules:\n\n")
	} else {
		fmt.Printf("Rules for user %v:\n\n", id)
	}

	for id, rule := range rules {
		fmt.Printf("(%d) ", id)
		if rule.Regex {
			if rule.Allow {
				fmt.Printf("Allow Regex: \t%s\n", rule.Regexp.Raw)
			} else {
				fmt.Printf("Disallow Regex: \t%s\n", rule.Regexp.Raw)
			}
		} else {
			if rule.Allow {
				fmt.Printf("Allow Path: \t%s\n", rule.Path)
			} else {
				fmt.Printf("Disallow Path: \t%s\n", rule.Path)
			}
		}
	}
}
