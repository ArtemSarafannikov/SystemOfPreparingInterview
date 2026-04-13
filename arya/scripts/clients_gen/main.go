//go:build ignore

// clients_gen.go
package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

var (
	configPath = flag.String("config", "config/clients.yaml", "config file path")
)

// Config .
type Config struct {
	Pkg     string `yaml:"package"`
	Clients []*ClientConfig
	Mock    *Mock
}

// ClientConfig .
type ClientConfig struct {
	Service  string
	Pkg      string `yaml:"package"`
	Client   string
	Field    string
	PkgAlias string
	Timeout  string
}

type Mock struct {
	PbFilePath  string `yaml:"pb_file_path"`
	MocksOutput string `yaml:"mocks_output"`
	Pkg         string `yaml:"package"`
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatalf("no command specified")
	}

	config := readConfig(*configPath)

	switch args[0] {
	case "container":
		genClients(*config)
	case "mocks":
		genMocks(*config)
	default:
		log.Fatalf("unknown command: %s", args[0])
	}
}

func readConfig(path string) *Config {
	yamlFile, err := os.ReadFile(path) // nolint: gosec
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	config := &Config{}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Fatalf("error parsing config file: %v", err)
	}
	return config
}

func genClients(config Config) {
	if len(config.Clients) == 0 {
		return
	}

	for _, client := range config.Clients {
		if client.Field == "" {
			client.Field = UnCapitalize(client.Client)
		}
		client.PkgAlias = fmt.Sprintf("%sPkg", client.Field)
		if client.Timeout == "" {
			client.Timeout = "3s"
		}
	}

	templateFile, err := os.ReadFile("scripts/clients_gen/clients.template")
	if err != nil {
		log.Fatalf("cannot read template: %v", err)
	}
	funcManp := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"Title":   Capitalize,
	}

	t, err := template.New("").Funcs(funcManp).Parse(string(templateFile))
	if err != nil {
		log.Fatalf("cannot parse template: %v", err)
	}

	err = t.Execute(os.Stdout, config)
	if err != nil {
		log.Fatalf("cannot execute template: %v", err)
	}
}

// Capitalize .
func Capitalize(s string) string {
	return cases.Title(language.English, cases.Compact).String(s[:1] + s[1:])
}

// UnCapitalize .
func UnCapitalize(s string) string {
	return cases.Lower(language.English, cases.Compact).String(s[:1] + s[1:])
}

func genMocks(config Config) {
	genDir := config.Mock.PbFilePath
	mocksDir := config.Mock.MocksOutput

	fmt.Println("🚀 Starting mock generation...")

	if err := os.MkdirAll(mocksDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create mocks directory: %v\n", err)
		os.Exit(1)
	}

	successCount := 0
	errorCount := 0

	err := filepath.WalkDir(genDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("⚠️  Error accessing path %s: %v\n", path, err)
			return nil
		}

		if d.IsDir() || !strings.HasSuffix(d.Name(), "_grpc.pb.go") {
			return nil
		}

		baseName := strings.TrimSuffix(d.Name(), "_grpc.pb.go")
		mockFile := filepath.Join(mocksDir, baseName+"_mock.go")

		cmd := exec.Command("mockgen",
			"-source", path,
			"-destination", mockFile,
			"-package", config.Mock.Pkg,
		)

		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Printf("❌ Failed to generate mock for %s: %v\n", filepath.Base(path), err)
			if len(output) > 0 {
				fmt.Printf("   Output: %s\n", string(output))
			}
			errorCount++
		} else {
			if err := addRequiredImports(mockFile); err != nil {
				fmt.Printf("⚠️  Failed to add imports to %s: %v\n", filepath.Base(mockFile), err)
				errorCount++
			} else {
				fmt.Printf("✅ Successfully generated: %s\n", filepath.Base(mockFile))
				successCount++
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("❌ Error walking through directory: %v\n", err)
		os.Exit(1)
	}

	// Выводим итоги
	fmt.Printf("\n📊 Generation results:\n")
	fmt.Printf("   ✅ Successful: %d\n", successCount)
	fmt.Printf("   ❌ Failed: %d\n", errorCount)
	fmt.Printf("   📁 Total processed: %d\n", successCount+errorCount)

	if errorCount > 0 {
		fmt.Println("\n⚠️  Some mocks failed to generate. Check the errors above.")
		os.Exit(1)
	}

	if successCount == 0 {
		fmt.Println("\nℹ️  No _grpc.pb.go files found in ./gen directory")
	} else {
		fmt.Println("\n🎉 All mocks generated successfully!")
	}
}

func addRequiredImports(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, content, 0644)
}
