package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/stretchr/testify/assert/yaml"
	"golang.org/x/tools/go/analysis/passes/waitgroup"
)

//
//	"account.email.change.error.duplicated": "Provided email belongs to another account"
//	"account.faq": "FAQ"

type T struct {
	Keys map[string]string `yaml:"keys"`
}

type Result struct {
	Keys map[string]string
}

func processFile(paths *[]string, key string) int {
	var count int
	for _, f := range *paths {
		file, err := os.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}
		count += strings.Count(string(file), key)
	}
	return count
}

func parseYamlFile(fileName string, t *T) error {
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

type TargetPaths struct {
	Paths []string
}

const translationsFileValie = "translations.yaml"

// const testFile = "test.yaml"

func main() {
	file := flag.String("file", translationsFileValie, "A path to a yaml file")
	translations := T{}
	err := parseYamlFile(*file, &translations)
	if err != nil {
		log.Fatalf("erorr: %s", err)
	}

	targetPaths, err := getFilePaths()
	if err != nil {
		log.Fatalf("erorr: %s", err)
	}

	for translationKey := range translations.Keys {
		go func() {
			var wg sync.WaitGroup
			count := processFile(&targetPaths, translationKey)
			if count == 0 {
				fmt.Printf("No keys found for: %s\n", translationKey)
			}
		}()
		// if count > 0 {
		// 	fmt.Printf("Found %d of %s\n", count, translationKey)
		// }
	}

	// XXL_FES_PATH := os.Getenv("XXL_FES_PATH")
	// for _, value := range translations.Keys {
	// 	fmt.Printf("📍-%s\n", value)
	// 	err = filepath.WalkDir(XXL_FES_PATH, func(path string, d os.DirEntry, err error) error {
	// 		if err != nil {
	// 			return err
	// 		}
	// 		if !d.IsDir() {
	// 			// log.Println(path)
	// 			if strings.HasSuffix(path, ".tsx") {
	// 				go processFile(path, value)
	// 			}
	// 		}
	// 		return nil
	// 	})
	// }

	if err != nil {
		log.Fatal(err)
	}
}
