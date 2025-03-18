package main

type GSPackage struct {
	Name        string
	Tag         string
	AssetUrl    string
	Script      string
	Platform    string
	LastModfied string
	Service     string
}

type GSConfig struct {
	GSPM struct {
		Version string
	}
	Packages []GSPackage
}

type CliArgs struct {
	ConfigDir      string   `arg:"env:GSPM_CONFIG_DIR"`
	ShellCommand   string   `arg:"env:GSPM_SHELL_COMMAND"`
	GitHubToken    string   `arg:"env:GSPM_GITHUB_TOKEN"`
	GitLabToken    string   `arg:"env:GSPM_GITLAB_TOKEN"`
	BitbucketToken string   `arg:"env:GSPM_BITBUCKET_TOKEN"`
	Command        string   `arg:"positional" help:"Command to run. Must be add, remove, install, edit, info or list."`
	Repos          []string `arg:"positional" help:"Repos from Git Services (GitHub supported only for now). Format: username/repository"`
	Scripts        []string `arg:"-s,--script,separate" help:"Script to run after download a asset. Use {{ASSET}} to reference the asset path."`
	Service        string   `arg:"--service" default:"github" help:"Git Service (github, gitlab, bitbucket)"`
}
