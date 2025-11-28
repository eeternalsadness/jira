# Copilot Instructions

## Project Overview

A Go CLI tool for interacting with Jira APIs. Enables developers to manage Jira issues, projects, and configurations from the command line.

## Technologies & Libraries

- **Go 1.23.4**
- **Cobra** - CLI framework (following cobra.dev recommendations)
- **Viper** - Configuration management with YAML
- **Standard library** - HTTP client for Jira REST API calls

## Project Structure (Cobra Best Practices)

```
├── main.go                 # Entry point - minimal, calls cmd.Execute()
├── cmd/
│   └── root.go             # Root command and global flags, link child commands here
├── internal/               # Private application code
│   ├── cli/                # all subcommands go here
│   │   ├── configure/      # code for the 'configure' subcommand
│   │   │   └── command.go  # main file for the 'configure' subcommad, exports NewCommand() which returns the 'configure' subcommand
│   │   └── get/            # code for the 'get' subcommand
│   │       └── command.go  # main file for the 'get' subcommad, exports NewCommand() which returns the 'get' subcommand
│   └── util/               # common utils for the commands
└── pkg/                    # Public API (if needed for external use)
    └── jira/               # all jira API integration goes here
```

## Go/Cobra Coding Standards

### General Guideline

- Make sure each command has the following attributes: `Use`, `Short`, `Long`, `Aliases` (if necessary), `Example`
- Use `RunE` for commands that can return errors
  - Use `cmd.SilenceUsage = true` to suppress printing usage on runtime errors
- Use `PersistentPreRunE` and `PreRunE` to handle configuration logic before command logic is executed
- Use `Persistent*` functions and attributes if child commands should inherit such functions or attributes
- Keep the `*RunE` and `*PersistentPreRunE` functions short and readable. If they're more than 10 lines long, create a new function and call that function inside `*RunE` or `*PersistentPreRunE`

### Naming Conventions

- CamelCase for exported functions/types
- Use descriptive names: `CreateIssueOptions` not `Options`
- Command variables: `createIssueCmd`, `getRootCmd`
- Flag names: kebab-case (`--project-id`)

### Order of Precedence for Configuration Values

Each command should parse configuration values in the following order (from highest priority to lowest):

1. Command-line flag (e.g. `--port 3000`)
2. Environment variable (e.g. `APP_PORT=4000`)
3. Configuration file value (e.g. `port: 5000` in `config.yaml`)
4. Some sensible default (e.g. `8000`)

### Command Examples

```go
// Proper command structure
var createIssueCmd = &cobra.Command{
  Use:   "issue",
  Short: "Create a new Jira issue",
  Long:  `Create a new Jira issue with interactive prompts...`,
  Aliases: []string{"iss"},
  Example: `# create a new Jira issue with the default project ID and issue type
jira create issue`,
  PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
    // handle configuration logic here
    return handleConfig(cmd)
  },
  RunE: func(cmd *cobra.Command, args []string) error {
    // handle command logic here
    // use viper.Get* functions to get config values instead of getting them from flags
  },
}
```

