# jetfind

Jetfind is a configurable file search tool that combines parallel file system scanning with an intuitive terminal interface. It features configurable filtering algorithms, gitignore-like skip file system, and customizable themes.

## Features

- **Parallel Scanning**: Multi-threaded file system traversal
- **Interactive TUI**: Keyboard-driven interface built with Bubble Tea
- **Fuzzy Filtering**: Multiple filtering algorithms including fuzzy matching with Jaro-Winkler
- **Ignore File Support**: Respects `.findignore` files with glob pattern matching
- **Configurable Themes**: Customizable colors and styling
- **Command Integration**: Execute commands on selected files with `--post-cmd` flag
- **Cross-Platform**: Works on macOS, Linux, and Windows

## Installation

### Prerequisites

Go 1.19 or later is required to build jetfind from source.

### From Source

```bash
git clone https://github.com/yourusername/jetfind.git
cd jetfind
make build
```

The binary will be created at `./bin/jetfind`.

### Global Installation

To install jetfind globally and make it available from anywhere:

```bash
make install
```

This installs the binary to `$GOPATH/bin` (typically `~/go/bin`).

**Important**: Make sure `$GOPATH/bin` is in your PATH. Add this to your shell configuration file (`~/.zshrc`, `~/.bashrc`, or `~/.bash_profile`):

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Then reload your shell configuration:
```bash
source ~/.zshrc  # or your appropriate shell config file
```

Now you can run `jetfind` from anywhere in your terminal.

## Usage

### Basic Usage

```bash
# Launch interactive file finder
jetfind

# type filtering query
# press enter to select a file

# Execute command on selected file
jetfind --post-cmd vim
jetfind --cd-to --post-cmd code
jetfind --post-cmd cat
```

### Configuration

Jetfind uses a YAML configuration file located at the standard config directory for your operating system:

- **Linux**: `~/.config/jetfind/config.yaml`
- **macOS**: `~/Library/Application Support/jetfind/config.yaml`
- **Windows**: `%APPDATA%\jetfind\config.yaml`

The configuration file is automatically created with default values on first run if it doesn't exist.

#### Configuration Options

```yaml
filter:
  type: "fuzzy"           # Filter type: fuzzy
  algorithm: "jarowinkler" # Algorithm: jarowinkler
  threshold: 0.9          # Similarity threshold (0.0-1.0)

findignore:
  enable: false           # Enable .findignore file support
  hidden_ignore: false    # Ignore hidden files/directories

tui:
  highlighted_file:
    foreground: "#FFFFFF"
  query_box:
    text_foreground: "#F9FAFB"
    text_background: "#374151" 
    border_foreground: "#6B7280"
```

**Filter Configuration:**
- `type`: Filtering method (`fuzzy`, `contains`, `null` are currently supported)
- `algorithm`: Fuzzy matching algorithm (`jarowinkler`, `ngram`)
- `threshold`: Minimum similarity score (0.0-1.0, higher = more strict)

**Findignore Configuration:**
- `enable`: Whether to use `.findignore` files
- `hidden_ignore`: Automatically ignore hidden files and directories

**TUI Configuration:**
- `highlighted_file`: Colors for selected file in the list
- `query_box`: Styling for the search input box

### Ignore Files

Create a `.findignore` file in the configuration directory (where the config.yml is placed) to exclude files and directories:

```
# Comments start with #
*.log
*.tmp
node_modules/
.git/
build/
dist/

# Negate with ! to include files that would otherwise be ignored
!important.log
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [XDG](https://github.com/adrg/xdg) - Cross-platform config directories
- [YAML v3](https://gopkg.in/yaml.v3) - Configuration parsing

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

