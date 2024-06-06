# gspm

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

Git Services Package Manager

Support installing from releases with custom script.

> Work in progress!

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
task build
```

## Credits

- [alexflint/go-arg](https://github.com/alexflint/go-arg)
- [imroc/req](https://github.com/imroc/req)
- [pterm/pterm](https://github.com/pterm/pterm)
- [go-task/task](https://github.com/go-task/task)
