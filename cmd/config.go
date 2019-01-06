package cmd

import (
	"encoding/json"
	nerrors "errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/filebrowser/filebrowser/v2/auth"
	"github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	v "github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management utility",
	Long:  `Configuration management utility.`,
	Args:  cobra.NoArgs,
}

func addConfigFlags(f *pflag.FlagSet) {

	addUserFlags(f)

	vaddP(f, "signup", "s", false, "allow users to signup")
	vadd(f, "shell", "", "shell command to which other commands should be appended")

	vadd(f, "auth.method", string(auth.MethodJSONAuth), "authentication type")
	vadd(f, "auth.header", "", "HTTP header for auth.method=proxy")

	vadd(f, "recaptcha.host", "https://www.google.com", "use another host for ReCAPTCHA. recaptcha.net might be useful in China")
	vadd(f, "recaptcha.key", "", "ReCaptcha site key")
	vadd(f, "recaptcha.secret", "", "ReCaptcha secret")

	vadd(f, "branding.name", "", "replace 'File Browser' by this name")
	vadd(f, "branding.files", "", "path to directory with images and custom styles")
	vadd(f, "branding.disableExternal", false, "disable external links such as GitHub links")

	if err := v.BindPFlags(f); err != nil {
		panic(err)
	}
}

func getAuthentication(cmd *cobra.Command) (settings.AuthMethod, auth.Auther) {
	method := settings.AuthMethod(mustGetString(cmd, "auth.method"))

	var auther auth.Auther
	if method == auth.MethodProxyAuth {
		header := mustGetString(cmd, "auth.header")
		if header == "" {
			panic(nerrors.New("you must set the flag 'auth.header' for method 'proxy'"))
		}
		auther = &auth.ProxyAuth{Header: header}
	}

	if method == auth.MethodNoAuth {
		auther = &auth.NoAuth{}
	}

	if method == auth.MethodJSONAuth {
		jsonAuth := &auth.JSONAuth{}

		host := mustGetString(cmd, "recaptcha.host")
		key := mustGetString(cmd, "recaptcha.key")
		secret := mustGetString(cmd, "recaptcha.secret")

		if key != "" && secret != "" {
			jsonAuth.ReCaptcha = &auth.ReCaptcha{
				Host:   host,
				Key:    key,
				Secret: secret,
			}
		}

		auther = jsonAuth
	}

	if auther == nil {
		panic(errors.ErrInvalidAuthMethod)
	}

	return method, auther
}

func printSettings(s *settings.Settings, auther auth.Auther) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "Sign up:\t%t\n", s.Signup)
	fmt.Fprintf(w, "Auth method:\t%s\n", s.AuthMethod)
	fmt.Fprintf(w, "Shell:\t%s\t\n", strings.Join(s.Shell, " "))
	fmt.Fprintln(w, "\nBranding:")
	fmt.Fprintf(w, "\tName:\t%s\n", s.Branding.Name)
	fmt.Fprintf(w, "\tFiles override:\t%s\n", s.Branding.Files)
	fmt.Fprintf(w, "\tDisable external links:\t%t\n", s.Branding.DisableExternal)
	fmt.Fprintln(w, "\nDefaults:")
	fmt.Fprintf(w, "\tScope:\t%s\n", s.Defaults.Scope)
	fmt.Fprintf(w, "\tLocale:\t%s\n", s.Defaults.Locale)
	fmt.Fprintf(w, "\tView mode:\t%s\n", s.Defaults.ViewMode)
	fmt.Fprintf(w, "\tCommands:\t%s\n", strings.Join(s.Defaults.Commands, " "))
	fmt.Fprintf(w, "\tSorting:\n")
	fmt.Fprintf(w, "\t\tBy:\t%s\n", s.Defaults.Sorting.By)
	fmt.Fprintf(w, "\t\tAsc:\t%t\n", s.Defaults.Sorting.Asc)
	fmt.Fprintf(w, "\tPermissions:\n")
	fmt.Fprintf(w, "\t\tAdmin:\t%t\n", s.Defaults.Perm.Admin)
	fmt.Fprintf(w, "\t\tExecute:\t%t\n", s.Defaults.Perm.Execute)
	fmt.Fprintf(w, "\t\tCreate:\t%t\n", s.Defaults.Perm.Create)
	fmt.Fprintf(w, "\t\tRename:\t%t\n", s.Defaults.Perm.Rename)
	fmt.Fprintf(w, "\t\tModify:\t%t\n", s.Defaults.Perm.Modify)
	fmt.Fprintf(w, "\t\tDelete:\t%t\n", s.Defaults.Perm.Delete)
	fmt.Fprintf(w, "\t\tShare:\t%t\n", s.Defaults.Perm.Share)
	fmt.Fprintf(w, "\t\tDownload:\t%t\n", s.Defaults.Perm.Download)
	w.Flush()

	b, err := json.MarshalIndent(auther, "", "  ")
	checkErr(err)
	fmt.Printf("\nAuther configuration (raw):\n\n%s\n\n", string(b))
}
