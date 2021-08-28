package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

type Config struct {
	Platform *string
	Language *string
}

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
	})

	// Setup
	cobra.CheckErr(rootCmd.Execute())
}

var rootCmd = &cobra.Command{
	Use:   "tldr",
	Short: "command line client for tldr",
	Long:  "command line client for tldr written in go",
	Run:   main_call,
}

func main_call(cmd *cobra.Command, args []string) {
	// Print version
	v, err := cmd.Flags().GetBool("version")
	if err != nil {
		log.Fatal(err.Error())
	}

	// Build config
	platform, err := cmd.Flags().GetString("platform")
	if err != nil {
		log.Fatal(err.Error())
	}

	language, err := cmd.Flags().GetString("language")
	if err != nil {
		log.Fatal(err.Error())
	}
	// Update cache
	update, err := cmd.Flags().GetBool("update_cache")
	if err != nil {
		log.Fatal(err.Error())
	}

	hide_custom, err := cmd.Flags().GetBool("hide_custom")
	if err != nil {
		log.Fatal(err.Error())
	}

	hide_official, err := cmd.Flags().GetBool("hide_official")
	if err != nil {
		log.Fatal(err.Error())
	}

	hide_cache_age, err := cmd.Flags().GetBool("hide_cache_age")
	if err != nil {
		log.Fatal(err.Error())
	}

	hide_cache_age_warning, err := cmd.Flags().GetBool("hide_cache_age_warning")
	if err != nil {
		log.Fatal(err.Error())
	}

	show_cache_details, err := cmd.Flags().GetBool("show_cache_details")
	if err != nil {
		log.Fatal(err.Error())
	}

	show_build_info, err := cmd.Flags().GetBool("show_build_info")
	if err != nil {
		log.Fatal(err.Error())
	}

	cfg := &Config{
		Platform: &platform,
		Language: &language,
	}

	pageDir := getCachePath() + "/tldr-main/pages/"
	pageDirExists := isFileExists(pageDir)

	if !pageDirExists || update {

		if err = updateCache(); err != nil {
			fmt.Println("Cache faild updated ，err:", err.Error())
			return
		}
		fmt.Println("Cache successfully updated")
	}

	available, modifiedtime := getCacheAge()
	if !available {
		fmt.Printf("Cache is not available - please check %s", getCachePath())
		return
	}

	if show_build_info {
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
		log.Fatal(err.Error())
	}

	if list {
		tldrs := getCachedCommandList(cfg)
		for i, tldr := range tldrs {
			fmt.Printf("%s\t\t", tldr)
			if i%5 == 0 {
				fmt.Printf("\n")
			}
		}
		fmt.Printf("\n")

		fmt.Printf("---------------------\n")
		fmt.Printf("%q\n", getLocalCommandList())
		return
	}

	if len(args) != 1 {
		cmd.Help()
		return
	}

	if !hide_cache_age {

		age := time.Since(modifiedtime)

		fmt.Println("Last modified time : ", modifiedtime.Local())
		fmt.Println("Cache age          : ", humanizeDuration(age))

		if !hide_cache_age_warning {

			if age.Hours() > 7*24 {
				fmt.Println("Local cache is older than 7 days - run command with -u option to update")
			}
		}
	}

	// Get command name
	command := args[0]

	found := false
	// check custom location first
	if !hide_custom {
		pageLoc, page, err := checkCustom(command)
		if err != nil {
			log.Fatal(err)
		}
		if len(page) > 0 {
			if show_cache_details {
				fmt.Println(GrayString("[custom - " + pageLoc + " ]"))
			} else {
				fmt.Println(GrayString("[custom]"))
			}

			fmt.Println(output(page))
			found = true
		}
	}

	// Get page from local folder
	if !hide_official {
		pageLoc, page, err := checkLocalCache(cfg, command)
		if err != nil {
			log.Fatal(err)
		}
		if len(page) > 0 {
			if show_cache_details {
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
	fmt.Println("GitCommit  : ", GitCommit)
	fmt.Println("GitVersion : ", GitVersion)
	fmt.Println("Version    : ", Version)
	fmt.Println("BuildDate  : ", BuildDate)
	fmt.Println("BuiltOnIp  : ", BuiltOnIp)
	fmt.Println("BuiltOnOs  : ", BuiltOnOs)
	fmt.Println("GoVersion  : ", GoVersion)
	fmt.Println("OsArch     : ", OsArch)
}
