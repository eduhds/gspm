package main

type GSPackage struct {
	Name        string
	Tag         string
	AssetUrl    string
	Script      string
	Platform    string
	LastModfied string
}

type GSConfig struct {
	GSPM struct {
		Version string
	}
	Packages []GSPackage
}

type args struct {
	ConfigDir string   `arg:"env:GSPM_CONFIG_DIR"`
	Command   string   `arg:"positional" help:"Command to run. Must be add, remove, update, install, edit, info or list."`
	Repos     []string `arg:"positional" help:"Repos from Git Services (GitHub supported only for now). Format: username/repository"`
	Scripts   []string `arg:"-s,--script,separate" help:"Script to run after download a asset. Use {{ASSET}} to reference the asset path."`
}
