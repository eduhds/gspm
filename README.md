# gspm

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Windows](https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)
![Linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)
![macOS](https://img.shields.io/badge/mac%20os-000000?style=for-the-badge&logo=macos&logoColor=F0F0F0)

Git Services Package Manager (GitHub supported only for now).

Support installing from releases with custom script.

<p align="center">
  <img src="./screencast.gif" style="display: block; margin: auto;" />
</p>

## ✨ Features

-   Get asset from repository releases.
-   Run custom script to install, uncompress, move, etc.
-   Save repository/package info for future installs.
-   Interactive mode.
-   Multi-platform: Windows, Linux, MacOS.
-   Public and Private repositories

## ⬇️ Install

<p>
  <a href="https://github.com/eduhds/gspm/releases/download/v0.2.2/gspm_Windows_x86_64.zip"><img src="assets/BadgeWindows.png" alt="Windows" width="120" /></a>
  <a href="https://github.com/eduhds/gspm/releases/download/v0.2.2/gspm_Linux_x86_64.tar.gz"><img src="assets/BadgeLinux.png" alt="Linux" width="120" /></a>
  <a href="https://github.com/eduhds/gspm/releases/download/v0.2.2/gspm_Darwin_x86_64.tar.gz"><img src="assets/BadgeMacOS.png" alt="macOS" width="120" /></a>
</p>

### MacOS/Linux (recommended)

```sh
curl -sL https://dub.sh/gspm | bash
```

### Installing manually

Download manually from [releases](https://github.com/eduhds/gspm/releases).

### Use `gspm` to update itself:

-   MacOS/Linux

```sh
gspm add eduhds/gspm -s 'sudo tar -C /usr/local/bin -xzf {{ASSET}} gspm && rm {{ASSET}}'
```

## 📖 Usage

### CLI

```sh
Usage: gspm [--configdir CONFIGDIR] [--githubtoken GITHUBTOKEN] [--script SCRIPT] [COMMAND [REPOS [REPOS ...]]]

Positional arguments:
  COMMAND                Command to run. Must be add, remove, update, install, edit, info or list.
  REPOS                  Repos from Git Services (GitHub supported only for now). Format: username/repository

Options:
  --configdir CONFIGDIR [env: GSPM_CONFIG_DIR]
  --githubtoken GITHUBTOKEN [env: GSPM_GITHUB_TOKEN]
  --script SCRIPT, -s SCRIPT
                         Script to run after download a asset. Use {{ASSET}} to reference the asset path.
  --help, -h             display this help and exit
  --version              display version and exit
```

### Examples

```sh
# Add
gspm add username/repository
gspm add username/repository@tag
gspm add username/repository@latest

# Info, Edit
gspm info username/repository
gspm edit username/repository

# Using inline Script
gspm <add|edit> username/repository -s 'your script here'

# Remove (only from ~/.config/gspm.json)
gspm remove username/repository
# Remove (from ~/.config/gspm.json and running script to remove from system)
gspm remove username/repository -s 'your script here'

# List
gspm list

# Install (from ~/.config/gspm.json)
gspm install

# Custom config dir path
GSPM_CONFIG_DIR=/path/to/custom/dir gspm <command> ...

# GitHub private repositories
GSPM_GITHUB_TOKEN='your token here' gspm add username/repository
```

## 🛠️ Development

```sh
# Run
go run . <command> <arguments>

# Build
task build:release
```

## 🤝 Support

[![BuyMeACoffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-ffdd00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=black)](https://www.buymeacoffee.com/eduhds)
[![Ko-Fi](https://img.shields.io/badge/Ko--fi-F16061?style=for-the-badge&logo=ko-fi&logoColor=white)](https://ko-fi.com/eduhds)

## 📜 License

[GPL-3.0 license](./LICENSE.txt)

## 🫂 Contributing

Contributions are welcome, see [Contributions Guide](./CONTRIBUTING.md) and [Code of Conduct](./CODE_OF_CONDUCT.md).

## 🙏 Credits/Thanks

-   [alexflint/go-arg](https://github.com/alexflint/go-arg)
-   [imroc/req](https://github.com/imroc/req)
-   [pterm/pterm](https://github.com/pterm/pterm)
-   [go-task/task](https://github.com/go-task/task)
