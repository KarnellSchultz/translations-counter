# Translation Key Usage Tracker

A Go-based CLI tool that analyzes your React/TypeScript codebase to track the usage of translation keys. This tool helps identify which translation keys are being used in your project and, more importantly, which ones are not being used at all.

## Features

- 🔍 **Scans TypeScript/TSX files** for translation key usage
- 📊 **Generates usage reports** in JSON or CSV format
- 🚀 **Concurrent processing** for fast analysis of large codebases
- 🎯 **Unused key detection** to help clean up translation files
- 📝 **Multiple output formats** for easy integration with other tools

## Installation

### Prerequisites

- Go 1.19 or higher (for building from source)
- Make (optional, for using the Makefile)

### Using Make (Recommended)

The project includes a Makefile with convenient installation options:

#### System-wide Installation
```bash
# Install to /usr/local/bin (may require sudo)
make install
```

#### User Installation
```bash
# Install to ~/bin (no sudo required)
make install-user

# Add ~/bin to your PATH if not already there:
echo 'export PATH=$PATH:$HOME/bin' >> ~/.bashrc  # or ~/.zshrc
source ~/.bashrc  # or source ~/.zshrc
```

### Manual Build

```bash
# Clone the repository
git clone <repository-url>
cd translations-key-usage-tracker

# Build the binary
go build -o translation-key-usage-tracker main.go

# Move to a directory in your PATH
mv translation-key-usage-tracker /usr/local/bin/
```

### From Source

```bash
go install github.com/yourusername/translations-key-usage-tracker@latest
```

## Usage

### Basic Usage

```bash
translation-key-usage-tracker -input translations.yaml -root-path /path/to/your/project
```

### Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-input` | Path to the YAML file containing translation keys | `translations.yaml` |
| `-output` | Path to the output file | `output.json` |
| `-format` | Output format (`json` or `csv`) | `json` |
| `-root-path` | Root directory of the project to scan | `.` (current directory) |
| `-unused-only` | Include only unused translation keys in output | `false` |
| `-verbose` | Print configuration and progress information | `false` |

### Examples

#### Find all unused translation keys in JSON format
```bash
translation-key-usage-tracker \
  -input translations.yaml \
  -root-path /path/to/react-project \
  -unused-only \
  -output unused-keys.json
```

#### Generate a CSV report of all translation key usage
```bash
translation-key-usage-tracker \
  -input translations.yaml \
  -root-path /path/to/react-project \
  -format csv \
  -output translation-usage.csv
```

#### Verbose mode for debugging
```bash
translation-key-usage-tracker \
  -input translations.yaml \
  -root-path ./src \
  -verbose
```

## Input Format

The tool expects a YAML file with translation keys in the following format:

```yaml
keys:
  home.title: "Welcome to our application"
  home.description: "This is the home page"
  nav.about: "About Us"
  nav.contact: "Contact"
  footer.copyright: "© 2024 Your Company"
```

## Output Formats

### JSON Output

When using `-format json`, the output will be a JSON object with keys and their usage counts:

```json
{
  "home.title": 5,
  "home.description": 2,
  "nav.about": 3,
  "nav.contact": 3,
  "footer.copyright": 0
}
```

### CSV Output

When using `-format csv`, the output will be a CSV file with headers:

```csv
key,count
home.title,5
home.description,2
nav.about,3
nav.contact,3
footer.copyright,0
```

### Unused Keys Only

When using the `-unused-only` flag, the output will only include keys with a count of 0:

```json
{
  "footer.copyright": 0,
  "deprecated.old_feature": 0
}
```

## Development

### Building from Source

```bash
# Clone the repository
git clone <repository-url>
cd translations-key-usage-tracker

# Download dependencies
make deps

# Build the project
make build

# Run tests
make test
```

### Available Make Commands

```bash
make help          # Show all available commands
make build         # Build the binary
make install       # Install to /usr/local/bin (requires sudo)
make install-user  # Install to ~/bin (no sudo needed)
make uninstall     # Remove from system
make clean         # Remove built binaries and output files
make test          # Run tests
make run           # Build and run with example arguments
make fmt           # Format code
make vet           # Run go vet for linting
make deps          # Download and tidy dependencies
```

### Running in Development

For development, you can use the `make run` command, but you'll need to set up the environment:

```bash
# Set the project path (if using make run)
export XXL_FES_PATH=/path/to/your/react/project

# Ensure you have a translations.yaml file in the current directory
# Then run
make run
```

## How It Works

1. **Parse Translation File**: The tool reads the YAML file containing all translation keys
2. **Scan Project Files**: Recursively walks through the specified directory, finding all `.tsx` files
3. **Concurrent Processing**: Uses goroutines to process multiple files simultaneously for better performance
4. **Count Occurrences**: For each translation key, counts how many times it appears in each file
5. **Aggregate Results**: Combines counts from all files
6. **Generate Output**: Writes results in the specified format (JSON or CSV)

## Performance Considerations

- The tool uses concurrent processing with goroutines to handle large codebases efficiently
- File reading and processing are done in parallel
- Results are aggregated using channels for thread-safe operations

## Use Cases

- **Clean up unused translations**: Identify and remove translation keys that are no longer used
- **Translation audit**: Get a complete overview of translation key usage
- **CI/CD integration**: Automate checks for unused translations in your build pipeline
- **Refactoring assistance**: Understand which translations are heavily used before refactoring
- **Documentation**: Generate reports of translation usage for documentation purposes

## Troubleshooting

### Common Issues

1. **"Go is not installed" error**
   - Install Go from [https://go.dev/doc/install](https://go.dev/doc/install)
   - On macOS: `brew install go`

2. **"translations.yaml not found" error**
   - Ensure the YAML file exists in the specified location
   - Use the `-input` flag to specify a different path

3. **No files being scanned**
   - Verify the `-root-path` points to the correct directory
   - Check that your TypeScript files have the `.tsx` extension

4. **Permission denied during installation**
   - Use `make install-user` for user-level installation
   - Or use `sudo make install` for system-wide installation

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[Specify your license here]

## Author

[Your name/organization]