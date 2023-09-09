package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	platform string
	language string
)

func main() {

	// Default flags
	rootCmd.Flags().SortFlags = false
	rootCmd.Flags().BoolP("help", "h", false, "show this help message and exit")
	rootCmd.Flags().BoolP("version", "v", false, "show program's version number and exit")
	rootCmd.Flags().BoolP("update_cache", "u", false, "Update the local cache of pages and exit")
	rootCmd.Flags().StringP("platform", "p", "linux", "Override the operating system [linux, osx, sunos, windows, common]")
	rootCmd.Flags().StringP("language", "L", "en", "Override the default language")
	rootCmd.Flags().BoolP("list", "l", false, "List all available commands for operating system matching regex- use regex")
	rootCmd.Flags().BoolP("hide_custom", "c", false, "Hide Custom TLDR results from "+getCustomPath())
	rootCmd.Flags().BoolP("hide_official", "o", false, "Hide Official TLDR result from "+getOfficalPath())
	rootCmd.Flags().BoolP("hide_cache_age", "a", false, "Hide Official TLDR cache age")
	rootCmd.Flags().BoolP("hide_cache_age_warning", "w", false, "Hide cachee age warning when older than 7 days")
	rootCmd.Flags().BoolP("show_cache_details", "d", false, "Show cache location details")
	rootCmd.Flags().BoolP("show_build_info", "b", false, "Show build information")
	rootCmd.Flags().BoolP("no-color", "n", false, "Disable color output")
	rootCmd.Flags().BoolP("edit", "e", false, "Edit custom tldr document")

	rootCmd.SetUsageTemplate(`Usage:{{if .Runnable}}
	{{.UseLine}} <tldr to lookup> {{end}}{{if .HasAvailableSubCommands}}
	{{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}
  
  Aliases:
	{{.NameAndAliases}}{{end}}{{if .HasExample}}
  
  Examples:
  {{.Example}}{{end}}{{if .HasAvailableSubCommands}}
  
  Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
	{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}
  
  Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}
  
  Global Flags:
  {{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}
  
  Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
	{{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
  
  Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
  `)

	origHelpFunc := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		origHelpFunc(cmd, args)
		fmt.Println("")
		displayBuildInfo()
		displayCacheInfo()
	})

	// Setup
	cobra.CheckErr(rootCmd.Execute())
}

var rootCmd = &cobra.Command{
	Use:   "tldr",
	Short: "command line client for tldr",
	Long:  "command line client for tldr written in go",
	Run:   mainCall,
}

func mainCall(cmd *cobra.Command, args []string) {

	// Disable color if cmd line set
	flagNoColor, err := cmd.Flags().GetBool("no-color")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if flagNoColor {
		color.NoColor = true // disables colorized output
	}

	// Print version
	v, err := cmd.Flags().GetBool("version")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Build config
	platform, err = cmd.Flags().GetString("platform")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	language, err = cmd.Flags().GetString("language")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Update cache
	update, err := cmd.Flags().GetBool("update_cache")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	hideCustom, err := cmd.Flags().GetBool("hide_custom")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	hideOfficial, err := cmd.Flags().GetBool("hide_official")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	hideCacheAge, err := cmd.Flags().GetBool("hide_cache_age")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	hideCacheAgeWarning, err := cmd.Flags().GetBool("hide_cache_age_warning")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	showCacheDetails, err := cmd.Flags().GetBool("show_cache_details")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	showBuildInfo, err := cmd.Flags().GetBool("show_build_info")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	pageDir := getCachePath() + "/tldr-main/pages/"
	pageDirExists := isFileExists(pageDir)

	if !pageDirExists || update {

		if err = updateCache(); err != nil {
			ErrorOutput.Printf("Cache faild updated ，err: %s\n", err.Error())
			return
		}
		fmt.Println("Cache successfully updated")
	}

	available, modifiedtime := getCacheAge()
	if !available {
		ErrorOutput.Printf("Cache is not available - please check %s\n", getCachePath())
		return
	}

	edit, err := cmd.Flags().GetBool("edit")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if showBuildInfo {
		displayBuildInfo()
		return
	}

	if v {
		fmt.Printf("%s - TLDR cache: %s\n", getVersion(), modifiedtime.Local().String())
		return
	}

	// Print list of available commands
	list, err := cmd.Flags().GetBool("list")
	if err != nil {
		fmt.Println(err.Error())
		return

	}

	if list {
		showAvailableTldrs()
		return
	}

	if len(args) != 1 {
		cmd.Help()
		return
	}

	if !hideCacheAge {

		age := time.Since(modifiedtime)

		fmt.Println("Last modified time : ", modifiedtime.Local())
		fmt.Println("Cache age          : ", humanizeDuration(age))

		if !hideCacheAgeWarning {

			if age.Hours() > 7*24 {
				ErrorOutput.Printf("Local cache is older than 7 days - run command with -u option to update\n")
			}
		}
	}

	// Get command name
	command := args[0]

	if edit {
		filename := buildCustomPath(command)
		fmt.Printf("Edit custom doc: %s - %s\n", command, filename)
		editFile(filename)
		return
	}

	found := false
	// check custom location first
	if !hideCustom {
		pageLoc, page, err := checkCustom(command)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if len(page) > 0 {
			if showCacheDetails {
				fmt.Println(GrayString("[custom - " + pageLoc + " ]"))
			} else {
				fmt.Println(GrayString("[custom]"))
			}

			fmt.Println(output(page))
			found = true
		}
	}

	// Get page from local folder
	if !hideOfficial {
		pageLoc, page, err := checkLocalCache(command)
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(page) > 0 {
			if showCacheDetails {
				fmt.Println(GrayString("[local - " + pageLoc + " ]"))
			} else {
				fmt.Println(GrayString("[local]"))
			}
			fmt.Println(output(page))
			found = true
		}
	}

	if !found {
		fmt.Printf("`%s` documentation is not available. Consider contributing Pull Request to https://github.com/tldr-pages/tldr or create local tldr file: %s\n", command, buildCustomPath(command))
	}

}

func displayBuildInfo() {
	fmt.Println("BuildDate    : ", BuildDate)
	fmt.Println("LatestCommit : ", LatestCommit)
	fmt.Println("Version      : ", Version)
	fmt.Println("GitRepo      : ", GitRepo)
	fmt.Println("GitBranch    : ", GitBranch)
}

func displayCacheInfo() {
	available, modifiedtime := getCacheAge()
	if !available {
		ErrorOutput.Printf("Cache is not available - please check %s\n", getCachePath())
		return
	}
	age := time.Since(modifiedtime)
	ErrorOutput.Printf("Last modified time : %s\n", modifiedtime.Local())
	ErrorOutput.Printf("Cache age          : %s\n", humanizeDuration(age))

}

func showAvailableTldrs() {
	tldrs := getCachedCommandList()
	for i, tldr := range tldrs {
		fmt.Printf("%s\t\t", tldr)
		if i%5 == 0 {
			fmt.Printf("\n")
		}
	}
	fmt.Printf("\n")

	fmt.Printf("---------------------\n")
	fmt.Printf("%q\n", getLocalCommandList())

}

func editFile(filename string) {

	var editor string

	switch runtime.GOOS {
	case "linux", "darwin":
		editor = os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}
		break

	case "windows":
		editor = "notepad"
	default:
		editor = "vi"
	}

	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
