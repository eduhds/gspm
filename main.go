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

//go:generate go-winres make --in windows/winres.json

const appname = "gspm"
const version = "1.0.1"
const description = "Thanks for using gspm, the Git Services Package Manager.\n"
const asciiArt = "\n ,adPPYb,d8  ,adPPYba,  8b,dPPYba,   88,dPYba,,adPYba,  \n" +
	"a8\"    `Y88  I8[    \"\"  88P'    \"8a  88P'   \"88\"    \"8a \n" +
	"8b       88   `\"Y8ba,   88       d8  88      88      88 \n" +
	"\"8a,   ,d88  aa    ]8I  88b,   ,a8\"  88      88      88 \n" +
	" `\"YbbdP\"Y8  `\"YbbdP\"'  88`YbbdP\"'   88      88      88 \n" +
	" aa,    ,88             88                              \n" +
	"  \"Y8bbdP\"              88                              \n"

var (
	service            string = "github"
	customConfigDir    string = ""
	customShellCommand string = ""
	downloadPrefix            = util.GetDownloadsDir()
)

func (CliArgs) Version() string {
	return fmt.Sprintf("%s v%s", appname, version)
}

func (CliArgs) Description() string {
	return description
}

func (CliArgs) Epilogue() string {
	return "For more information visit https://github.com/eduhds/gspm"
}

func main() {
	var args CliArgs
	arg.MustParse(&args)

	if args.Service != "" {
		validService := false
		for _, s := range gitservice.SupportedServices {
			if s == args.Service {
				service = args.Service
				validService = true
				break
			}
		}
		if !validService {
			tui.ShowError("Invalid service: " + args.Service)
			os.Exit(1)
		}
	}

	if service == "github" && args.GitHubToken != "" {
		gitservice.GHToken = args.GitHubToken
	}

	if service == "gitlab" && args.GitLabToken != "" {
		gitservice.GLToken = args.GitLabToken
	}

	if service == "bitbucket" && args.BitbucketToken != "" {
		gitservice.BBToken = args.BitbucketToken
	}

	if args.ConfigDir != "" {
		customConfigDir = args.ConfigDir
		tui.ShowInfo("Using custom config directory: " + customConfigDir)
	}

	if args.ShellCommand != "" {
		customShellCommand = strings.TrimSpace(args.ShellCommand)
		if len(strings.Split(customShellCommand, " ")) != 2 {
			tui.ShowError("Invalid shell command: " + customShellCommand)
			os.Exit(1)
		}
		tui.ShowInfo("Using custom shell command: " + customShellCommand)
	}

	config := ResolveConfig()

	if args.Command == "" {
		// Interactive mode
		util.ClearScreen()

		quit := false

		for !quit {
			tui.TextInfo(asciiArt)
			tui.ShowInfo(fmt.Sprintf("v%s", version))
			tui.ShowLine()

			option := tui.ShowOptions("What command do you want to use?", []string{"add", "remove", "install", "edit", "info", "list", "<cancel>"})

			if option == "<cancel>" {
				quit = true
			} else {
				repo := ""

				if option != "list" && option != "install" {
					var knownPackages []string
					for _, item := range PlatformPackages(config) {
						if MatchService(item) {
							knownPackages = append(knownPackages, item.Name)
						}
					}

					if len(knownPackages) > 0 {
						knownPackages = append(knownPackages, "<none>")
						repo = tui.ShowOptions("Select a package", knownPackages)
					} else {
						repo = "<none>"
					}

					if repo == "<none>" {
						repo = tui.ShowTextInput(fmt.Sprintf("What repository do you want \"%s\"? (Format: username/repository)", option), false, "")
					}
				}

				program := os.Args[0] // strings.TrimSpace(strings.Join(os.Args, " "))
				args := []string{option, repo}

				if len(os.Args) > 1 {
					args = append(args, os.Args[1:]...)
				}

				cmd := exec.Command(program, args...)
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
			} else {
				util.ClearScreen()
			}
		}

		return
	}

	// Non-interactive mode
	tui.ShowInfo(fmt.Sprintf("%s v%s", appname, version))
	tui.ShowLine()

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

			var targetPackage GSPackage
			targetPackage.Name = value

			if !IsValidPackageName(targetPackage.Name) {
				tui.ShowError("Invalid package name: " + targetPackage.Name)
				continue
			}

			for {
				gsp, err := ResolvePackage(targetPackage.Name, config, mustExist)

				if err != nil {
					tui.ShowError(err.Error())

					if mustExist && len(PlatformPackages(config)) > 0 {
						var knownPackages []string
						for _, item := range PlatformPackages(config) {
							if MatchService(item) {
								knownPackages = append(knownPackages, item.Name)
							}
						}

						knownPackages = append(knownPackages, "<cancel>")
						targetPackage.Name = tui.ShowOptions("Select a package", knownPackages)

						if targetPackage.Name == "<cancel>" {
							targetPackage = GSPackage{}
							break
						}
						continue
					} else {
						targetPackage = GSPackage{}
						break
					}
				} else {
					targetPackage = gsp
					break
				}
			}

			if targetPackage.Name == "" {
				continue
			}

			withScript := len(args.Scripts) > 0 && args.Scripts[index] != ""
			supportScript := args.Command == "add" || args.Command == "update" || args.Command == "edit" || args.Command == "remove"

			if supportScript {
				if withScript {
					targetPackage.Script = args.Scripts[index]
				}
			}

			if args.Command == "add" || args.Command == "update" {
				config = CommandAdd(config, targetPackage)
			} else if args.Command == "remove" {
				config = CommandRemove(config, targetPackage, withScript)
			} else if args.Command == "edit" {
				config = CommandEdit(config, targetPackage, withScript)
			} else if args.Command == "info" {
				CommandInfo(targetPackage)
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
