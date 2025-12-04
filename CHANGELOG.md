# Changelog

All notable changes to this project will be documented in this file.

## [v0.2.2] - 2025-12-04

### Added
- Headers added to selection lists for improved clarity

### Changed
- **BREAKING**: Configuration file format changed from CamelCase to kebab-case keys (e.g., `Domain` → `domain`, `ProjectIDs` → `project-ids`, `IssueTypeIDs` → `issue-type-ids`, `DefaultProjectID` → `default-project-id`, `DefaultIssueTypeID` → `default-issue-type-id`)

## [v0.2.1] - 2025-12-03

### Fixed
- Fixed off-by-one error in issue transition selection causing incorrect transition to be applied

## [v0.2.0] - 2025-12-03

### Added
- Interactive configuration selection in `jira configure` command

### Changed
- **BREAKING**: Command structure changed from `jira get issue` to `jira issue get`
- **BREAKING**: Command structure changed from `jira create issue` to `jira issue create`
- **BREAKING**: Command structure changed from `jira transition issue` to `jira issue transition`
- **BREAKING**: All separate configure commands (`jira configure projects`, `jira configure issueTypes`) consolidated into single `jira configure` command with interactive selection
- Improved user selection interface with clearer numbering (starting from 1)

## [v0.1.7] - 2025-11-28

### Added
- Support for multiple projects and issue types configuration
- `jira configure projects` command to configure available project IDs  
- `jira configure issueTypes` command to configure available issue types
- Project ID selection when creating issues with `-p/--project-id` flag
- Issue type selection when creating issues with `-t/--issue-type-id` flag
- Default project ID and issue type configuration support
- Issues can now be created without specifying project/issue type if defaults are configured

## [v0.1.6]

### Added
- Automatic version check
