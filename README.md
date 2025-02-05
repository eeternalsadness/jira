# jira

A CLI program that performs common Jira tasks. Made with the Cobra & Viper framework.

## Requirements

To use the program, make sure Go is installed and added to `PATH`.

Make sure to create a Jira API token for the program to use.

## Installation

To install the program, run the following command:

```shell
go install github.com/eeternalsadness/jira@latest
```

## Configure

After installation, configure the program by running:

```shell
jira configure
```

The configuration file is stored under `$HOME/.config/jira/config.yaml` by default.
