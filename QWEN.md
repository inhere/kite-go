# Qwen Code Context for Kite Project

## Project Information
- **Project Name**: Kite
- **Project Type**: Single CLI application 
- **Description**: Kite is a personal developer tool command application that provides a collection of utilities for developers

## Technical Context
- **Language/Version**: Go 1.23
- **Primary Dependencies**: github.com/gookit/config, github.com/gookit/rux, github.com/gookit/gcli, github.com/gookit/ini
- **Storage**: Files (reading project files like README.md, YAML configs)
- **Testing**: Go testing package
- **Target Platform**: Cross-platform (Linux, macOS, Windows)
- **Architecture**: CLI-first interface with text in/out protocol

## Project Structure
- `/cmd/kite`: Main CLI application

## Constitution Compliance
- All features must serve developer productivity and workflow enhancement
- All functionality must be accessible through a command-line interface
- All code must follow TDD practices
- Focus on testing areas involving CLI command integration and file system operations
- Implementation must follow Go best practices and be cross-platform compatible
