# Translation Key Usage Tracker

Scan a TypeScript React codebase for translation key usages and report the usage count per key as JSON or CSV.

This CLI:
- Loads translation keys from a YAML file
- Walks your codebase (rooted at `XXL_FES_PATH`) and scans all `.tsx` files
- Counts raw substring occurrences of each key
- Can output all keys or only unused keys (count == 0)
- Writes results to JSON or CSV


## Requirements

- Go installed and available on your PATH
- A YAML file containing your translation keys
- Environment variable `XXL_FES_PATH` set to the root directory to scan


## Installation

Using the Makefile in this repository:

- Download dependencies (optional if you already have a populated go.mod):
  
      make deps

- Build the binary (outputs `./translation-key-usage-tracker`):
  
      make build

- Install system-wide (may prompt for sudo):
  
      make install

- Install for current user only (no sudo):
  
      make install-user

Make sure the install location is on your PATH. The binary name is `translation-key-usage-tracker`.


## YAML input format

The tool expects a YAML file with a top-level `keys` mapping. The map keys are the translation keys to track; the values are not used by the scanner and may be any string.

    keys:
      checkout.shipping.title: "Shipping"
      cart.empty: "Your cart is empty"
      account.profile.save: "Save"

By default, the tool reads `translations.yaml` from the current directory. You can override the input path with `-input`.


## Usage

1) Set the scanning root:

    export XXL_FES_PATH=/absolute/path/to/your/frontend/repo

2) Run the tracker:

- JSON output (default):

      translation-key-usage-tracker \
        -input translations.yaml \
        -output output.json \
        -format json

- CSV output:

      translation-key-usage-tracker \
        -input translations.yaml \
        -output output.csv \
        -format csv

- Only include unused keys (count == 0):

      translation-key-usage-tracker \
        -input translations.yaml \
        -output unused.json \
        -format json \
        -unused-only

- Verbose mode:

      translation-key-usage-tracker \
        -input translations.yaml \
        -output output.json \
        -format json \
        -verbose


### CLI flags

- `-input string`      path to the YAML file containing translation keys (default: `translations.yaml`)
- `-unused-only`       include only unused translation keys in the output
- `-output string`     path to the output file (default: `output.json`)
- `-format string`     output format (`json` or `csv`; default: `json`)
- `-verbose`           print configuration and progress information

Environment:
- `XXL_FES_PATH`       root directory to scan for `.tsx` files


## Output examples

- JSON:

      {
        "checkout.shipping.title": 3,
        "cart.empty": 0,
        "account.profile.save": 7
      }

- CSV:

      key,count
      checkout.shipping.title,3
      cart.empty,0
      account.profile.save,7


## Makefile targets

- `make help`          Show target descriptions
- `make deps`          Download and tidy Go module dependencies
- `make build`         Build the binary (`./translation-key-usage-tracker`)
- `make install`       Install to `/usr/local/bin` (uses sudo if needed)
- `make install-user`  Install to `~/bin` (ensure it’s on your PATH)
- `make uninstall`     Remove installed binaries from system and user locations
- `make clean`         Remove build artifacts and sample outputs
- `make test`          Run Go tests
- `make fmt`           `go fmt ./...`
- `make vet`           `go vet ./...`
- `make run`           Build and run with example arguments
  - This checks `XXL_FES_PATH` and that `translations.yaml` exists
  - It executes:
    
        ./translation-key-usage-tracker -input translations.yaml -output output.json -format json


## How it works

- The `XXL_FES_PATH` environment variable points to the root directory to scan.
- The tool loads translation keys from a YAML file with a top-level `keys` mapping.
- It recursively scans for `.tsx` files and counts raw substring occurrences of each key using `strings.Count`.
- Results are aggregated across files and written via the chosen writer:
  - JSON writer: writes a map of `key -> count`
  - CSV writer: writes rows of `key,count`

Note: Matching is literal substring counting. It does not parse AST or interpret frameworks’ i18n APIs.


## Notes and limitations

- File types: Only `.tsx` files are scanned. If keys appear in `.ts`, `.js`, `.jsx`, etc., they will not be counted.
- Matching strategy: Raw substring counting is fast but:
  - May over-count when keys appear in comments or unrelated string literals.
  - May miss usages if your code constructs keys dynamically (e.g., string concatenation, template building).
- Output path: The tool writes to the `-output` path exactly as provided. Ensure your path/extension matches the `-format` you choose.
- Errors: Individual file read errors are logged and scanning continues. YAML/FS/IO errors will stop the run with a non-zero exit.


## Troubleshooting

- “XXL_FES_PATH environment variable not set”
  - Set it before running:
    
        export XXL_FES_PATH=/path/to/project

- No keys found or all counts are zero
  - Ensure your YAML has a top-level `keys` mapping and that the keys in the YAML match the substrings present in `.tsx` files.

- CSV vs JSON outputs
  - Set `-format` to `csv` or `json` and ensure your `-output` file extension matches your expectations.

- Permission denied when installing
  - Use `make install-user` or rerun `make install` with appropriate privileges.

- Long runtime on large repos
  - Scanning is concurrent per file, but counting is still substring-based across all keys. Consider narrowing `XXL_FES_PATH` or keys list if needed.


## Development

- Format: `make fmt`
- Lint: `make vet`
- Test: `make test`
- Dependencies: `make deps`

Contributions to improve matching accuracy (e.g., AST parsing, i18n-aware heuristics, or configurable file globs) are welcome.