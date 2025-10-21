package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/stretchr/testify/assert/yaml"
)

type TraslationsYaml struct {
	Keys map[string]string `yaml:"keys"`
}

type Result struct {
	Keys map[string]int
}

func parseYamlFile(fileName string, t *TraslationsYaml) error {
	yamlBytes, err := os.ReadFile(fileName)
	err = yaml.Unmarshal(yamlBytes, &t)
	return err
}

func getFilePaths() ([]string, error) {
	XXL_FES_PATH := os.Getenv("XXL_FES_PATH")
	eligiblePaths := make([]string, 0)
	err := filepath.WalkDir(XXL_FES_PATH, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			if strings.HasSuffix(path, ".tsx") {
				eligiblePaths = append(eligiblePaths, path)
			}
		}
		return nil
	})
	return eligiblePaths, err
}

func processFile(filebytes []byte, keys *ResultsMap) {
	fileContent := string(filebytes)
	for key := range *keys {
		count := strings.Count(fileContent, key)
		(*keys)[key] += count
	}
}

type TargetPaths struct {
	Paths []string
}

type Status string

const (
	StatusFound    Status = "found"
	StatusNotFound Status = "not-found"
)

type ResultsMap = map[string]int

type FileWriter interface {
	Write(*ResultsMap, string) error
}

type CsvWriter struct{}

func (c CsvWriter) Write(resultsMap *ResultsMap, dest string) error {
	csvFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)

	writer.Write([]string{"key", "count"})

	for key, value := range *resultsMap {
		writer.Write([]string{key, fmt.Sprintf("%d", value)})
	}
	writer.Flush()
	return writer.Error()
}

type JsonWriter struct{}

func (j JsonWriter) Write(resultMap *ResultsMap, dest string) error {
	jsonFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	bytes, err := json.Marshal(resultMap)
	if err != nil {
		return err
	}
	_, err = jsonFile.Write(bytes)
	return err
}

const translationsFileValie = "translations.yaml"

type App struct {
	translationFile string
	unusedOnly      bool
	output          string
	outputFormat    string
}

func NewApp(translationFile string, unusedOnly bool, output string, outputFormat string) *App {
	return &App{
		translationFile: translationFile,
		unusedOnly:      unusedOnly,
		output:          output,
		outputFormat:    outputFormat,
	}
}

func (a *App) Run() error {
	translations := TraslationsYaml{}
	err := parseYamlFile(a.translationFile, &translations)
	if err != nil {
		return fmt.Errorf("failed to parse translations file: %w", err)
	}

	targetPaths, err := getFilePaths()
	if err != nil {
		return fmt.Errorf("failed to get file paths: %w", err)
	}

	resultsMap := a.processFiles(targetPaths, translations.Keys)

	valueToWrite := resultsMap
	if a.unusedOnly {
		fmt.Println("Only unused keys")
		unusedKeys := make(ResultsMap, 0)
		for key, count := range resultsMap {
			if count == 0 {
				unusedKeys[key] = 0
			}
		}
		valueToWrite = unusedKeys
	}

	// Write results
	writer, err := a.getWriter()
	if err != nil {
		return err
	}

	if err := writer.Write(&valueToWrite, a.output); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

func (a *App) processFiles(targetPaths []string, keys map[string]string) ResultsMap {
	var mu sync.Mutex
	var wg sync.WaitGroup

	keysLen := len(keys)
	wg.Add(len(targetPaths))
	resultsMap := make(ResultsMap, keysLen)

	for _, path := range targetPaths {
		go func(item string) {
			defer wg.Done()

			fileBytes, err := os.ReadFile(item)
			if err != nil {
				log.Printf("Error reading file %s: %v", item, err)
				return
			}

			localResults := make(ResultsMap, keysLen)
			for key := range keys {
				localResults[key] = 0
			}
			processFile(fileBytes, &localResults)

			mu.Lock()
			for key, count := range localResults {
				resultsMap[key] += count
			}
			mu.Unlock()
		}(path)
	}
	wg.Wait()
	return resultsMap
}

func (a *App) getWriter() (FileWriter, error) {
	switch a.outputFormat {
	case "csv":
		return CsvWriter{}, nil
	case "json":
		return JsonWriter{}, nil
	default:
		return nil, fmt.Errorf("unknown format: %s", a.outputFormat)
	}
}

func main() {
	translationsFile := flag.String("translations", translationsFileValie, "path to the YAML file containing translation keys")
	unusedOnly := flag.Bool("unused-only", false, "include only unused translation keys in the output")
	defaultOutput := "output.json"
	output := flag.String("output", defaultOutput, "path to the output file (extension will match format)")
	format := flag.String("format", "json", "output format (json or csv); default: json")
	if format != nil && *format == "csv" {
		defaultOutput = "output.csv"
	}
	flag.Parse()

	app := NewApp(*translationsFile, *unusedOnly, *output, *format)
	if err := app.Run(); err != nil {
		log.Fatalf("error: %s", err)
	}
}
