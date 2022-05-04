/*
Package cmd
Copyright © 2022 Pavel Sidlo <github.com/dyamon-cz>
*/
package cmd

import (
	"fmt"
	"github.com/dyamon-cz/gograb/internal"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"strings"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gograb [module]",
	Short: "Easily grab Go modules.",
	Long: `
 ██████╗  ██████╗      ██████╗ ██████╗  █████╗ ██████╗ 
██╔════╝ ██╔═══██╗    ██╔════╝ ██╔══██╗██╔══██╗██╔══██╗
██║  ███╗██║   ██║    ██║  ███╗██████╔╝███████║██████╔╝
██║   ██║██║   ██║    ██║   ██║██╔══██╗██╔══██║██╔══██╗
╚██████╔╝╚██████╔╝    ╚██████╔╝██║  ██║██║  ██║██████╔╝
 ╚═════╝  ╚═════╝      ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝ 

Go Grab CLI allows you to download Go modules using key
words instead of remembering url.

Examples:
$ gograb gin
$ gograb protobuf
$ gograb x context`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			//cmd.HelpFunc()(cmd, args) // print help by default
			fmt.Println(cmd.Long)
			fmt.Println()
			return
		}

		search := strings.Join(args, " ") // allow multi-word search

		fmt.Printf("Searching: %s\n", search)

		modules := internal.SearchModules(search)

		if len(modules) > 0 {
			selectUi(modules)
		} else {
			fmt.Println("No results")
		}
	},
}

var templates = &promptui.SelectTemplates{
	Label:    "{{ . }}:",
	Active:   "\U0001F4E6 {{ .Name | cyan }} {{ \"(\" | red }}{{ .Path | red }}{{ \")\" | red }}", // [box] module (package/url)
	Inactive: "  {{ .Name | cyan }} {{ \"(\" | red }}{{ .Path | red }}{{ \")\" | red }}",          // module (package/url)
	Selected: "\U0001F4E6 {{ .Name | red | cyan }}",                                               // [box] module
	Details: `
{{ "Version:" | faint }}	{{ .Version }}
{{ "Published:" | faint }}	{{ .Published }}
{{ "Imports:" | faint }}	{{ .Imports }}
{{ "Licence:" | faint }}	{{ .Licence }}

{{ .Description }}`,
}

func selectUi(modules []internal.Module) {

	searcher := func(input string, index int) bool {
		module := modules[index]
		name := strings.Replace(strings.ToLower(module.Name), " ", "", -1)
		path := strings.Replace(strings.ToLower(module.Path), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input) || strings.Contains(path, input)
	}

	prompt := promptui.Select{
		Label:     "Modules",
		Items:     modules,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	// go get package

	cmd := exec.Command("go", "get", modules[i].Path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
