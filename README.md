# gspm

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Windows](https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)
![Linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)
![macOS](https://img.shields.io/badge/mac%20os-000000?style=for-the-badge&logo=macos&logoColor=F0F0F0)

Git Services Package Manager (GitHub supported only for now).

Support installing from releases with custom script.

## Install

### MacOS/Linux (recommended)

```sh
curl -sL https://dub.sh/gspm | bash
```

### Friendly installers

<p>
  <a href="https://github.com/eduhds/gspm/releases/download/v0.1.0/gspm-windows-amd64-setup.exe">
    <img src="assets/BadgeWindows.png" alt="Windows" width="120" />
  </a>
  <a href="https://github.com/eduhds/gspm/releases/download/v0.1.0/gspm-linux-amd64.AppImage">
    <img src="assets/BadgeLinux.png" alt="Linux" width="120" />
  </a>
  <a href="https://github.com/eduhds/gspm/releases/download/v0.1.0/gspm-darwin-amd64.dmg">
    <img src="assets/BadgeMacOS.png" alt="macOS" width="120" />
  </a>
</p>

### Installing manually

#### First time

Download manually from [releases](https://github.com/eduhds/gspm/releases).

#### Already have `gspm` installed

Use `gspm` to update itself.

```
# MacOS/Linux
gspm eduhds/gspm -s 'sudo tar -C /usr/local/bin -xzf {{ASSET}} && rm {{ASSET}}'
```

## Usage

```sh
# Add
gspm add username/repository
gspm add username/repository@tag
gspm add username/repository@latest

# Remove
gspm remove username/repository

# Install (from ~/.config/gspm.json)
gspm install

# List
gspm list

# Edit
gspm edit username/repository
```

## Development

```sh
# Run
go run . <command> <arguments>

# Build
task build:release
```

## Credits

- [alexflint/go-arg](https://github.com/alexflint/go-arg)
- [imroc/req](https://github.com/imroc/req)
- [pterm/pterm](https://github.com/pterm/pterm)
- [go-task/task](https://github.com/go-task/task)
