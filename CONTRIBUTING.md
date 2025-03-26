# Contributing

Contributions are always welcome, no matter how large or small!

We want this community to be friendly and respectful to each other. Please follow it in all your interactions with the project. Before contributing, please read the [code of conduct](./CODE_OF_CONDUCT.md).

## Development workflow

This project is a simple GO project without template.
TODO: more info

To get started with the project, run `go install` in the root directory to install the required dependencies for each package:

```sh
go install
```

You can use various commands from the root directory to work with the project.

To run interactive mode:

```sh
go run .
```

To run a command:

```sh
go run . <command> <user>/<repository>
```

### Commit message convention

We follow the [conventional commits specification](https://www.conventionalcommits.org/en) for our commit messages:

-   `fix`: bug fixes, e.g. fix crash due to deprecated method.
-   `feat`: new features, e.g. add new method to the module.
-   `refactor`: code refactor, e.g. migrate from class components to hooks.
-   `docs`: changes into documentation, e.g. add usage example for the module..
-   `test`: adding or updating tests, e.g. add integration tests using detox.
-   `chore`: tooling changes, e.g. change CI config.

TODO: Our pre-commit hooks verify that your commit message matches this format when committing.

### Linting and tests

TODO: Our pre-commit hooks verify that the linter and tests pass when committing.

### Sending a pull request

> **Working on your first pull request?** You can learn how from this _free_ series: [How to Contribute to an Open Source Project on GitHub](https://app.egghead.io/playlists/how-to-contribute-to-an-open-source-project-on-github).

When you're sending a pull request:

-   Prefer small pull requests focused on one change.
-   Verify that linters and tests are passing.
-   Review the documentation to make sure it looks good.
-   Follow the pull request template when opening a pull request.
-   For pull requests that change the API or implementation, discuss with maintainers first by opening an issue.

This contributing guide was inspired by [create-react-native-library](https://github.com/callstack/react-native-builder-bob).
