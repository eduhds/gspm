package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/eduhds/gspm/internal/gitservice"
	"github.com/eduhds/gspm/internal/tui"
	"github.com/eduhds/gspm/internal/util"
)

const appname = "gspm"
const version = "0.1.4"
const description = "Thanks for using gspm, the Git Services Package Manager.\n"
const asciiArt = "\n ,adPPYb,d8  ,adPPYba,  8b,dPPYba,   88,dPYba,,adPYba,  \n" +
	"a8\"    `Y88  I8[    \"\"  88P'    \"8a  88P'   \"88\"    \"8a \n" +
	"8b       88   `\"Y8ba,   88       d8  88      88      88 \n" +
	"\"8a,   ,d88  aa    ]8I  88b,   ,a8\"  88      88      88 \n" +
	" `\"YbbdP\"Y8  `\"YbbdP\"'  88`YbbdP\"'   88      88      88 \n" +
	" aa,    ,88             88                              \n" +
	"  \"Y8bbdP\"              88                              \n"

var customConfigDir string = ""

func (args) Version() string {
	return fmt.Sprintf("%s v%s", appname, version)
}

func (args) Description() string {
	return description
}

func (args) Epilogue() string {
	return "For more information visit https://github.com/eduhds/gspm"
}

var downloadPrefix = util.GetDownloadsDir()

func main() {
	var args args
	arg.MustParse(&args)

	if args.GitHubToken != "" {
		gitservice.GHToken = args.GitHubToken
	}

	if args.ConfigDir != "" {
		customConfigDir = args.ConfigDir
		tui.ShowInfo("Using custom config directory: " + customConfigDir)
	}

	if args.Command == "" {
		// Interactive mode
		tui.TextInfo(asciiArt)
		tui.ShowInfo(fmt.Sprintf("v%s", version))
		tui.ShowLine()

		quit := false

		for !quit {
			option := tui.ShowOptions("What command do you want to use?", []string{"add", "remove", "update", "install", "edit", "info", "list", "<cancel>"})

			if option == "<cancel>" {
				quit = true
			} else {
				repo := ""

				if option != "list" && option != "install" {
					repo = tui.ShowTextInput(fmt.Sprintf("What repository do you want \"%s\"? (Format: username/repository)", option), false, "")
				}

				program := strings.TrimSpace(strings.Join(os.Args, " "))

				cmd := exec.Command(program, option, repo)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				err := cmd.Run()

				if err != nil {
					tui.ShowError(err.Error())
				}

				tui.ShowLine()
				quit = !tui.ShowConfirm("Continue to another command?")
			}

			if quit {
				tui.TextSuccess(description)
				tui.ShowLine()
			}
		}

		return
	}

	// Non-interactive mode
	tui.ShowInfo(fmt.Sprintf("%s v%s", appname, version))
	tui.ShowLine()

	config := ResolveConfig()

	switch args.Command {
	case "list":
		CommandList(config)
	case "install":
		if len(args.Repos) > 0 || len(args.Scripts) > 0 {
			tui.ShowWarning("Ignoring args for install command.")
		}

		CommandInstall(config)
	case "add", "update", "remove", "edit", "info":
		if len(args.Repos) == 0 {
			tui.ShowError("No packages provided")
			return
		}

		mustExist := args.Command != "add"

		for index, value := range args.Repos {
			if index > 0 {
				tui.ShowLine()
			}

			gsp, err := ResolvePackage(value, config, mustExist)

			if err != nil {
				tui.ShowError(err.Error())
				continue
			}

			withScript := len(args.Scripts) > 0 && args.Scripts[index] != ""

			if args.Command == "add" || args.Command == "update" || args.Command == "edit" {
				if withScript {
					gsp.Script = args.Scripts[index]
				}
			}

			if args.Command == "add" {
				config = CommandAdd(config, gsp)
			} else if args.Command == "update" {
				config = CommandUpdate(config, gsp)
			} else if args.Command == "remove" {
				config = CommandRemove(config, gsp)
			} else if args.Command == "edit" {
				config = CommandEdit(config, gsp, withScript)
			} else if args.Command == "info" {
				CommandInfo(gsp)
			}
		}

		if args.Command != "info" {
			WriteConfig(config)
		}
	default:
		tui.ShowError("Unknown command: " + args.Command)
		os.Exit(1)
	}
}
