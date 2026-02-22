package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/n0xum/structify/internal/application"
	"github.com/n0xum/structify/internal/application/command"
	"github.com/n0xum/structify/internal/application/query"
	"github.com/n0xum/structify/internal/generator"
)

type Command struct {
	FS            *flag.FlagSet
	ToSQL         bool
	ToRepo        bool
	ModelFile     string
	InterfaceFile string
	OutputFile    string
	ShowVersion   bool
}

func NewCommand() *Command {
	cmd := &Command{
		FS: flag.NewFlagSet("structify", flag.ContinueOnError),
	}

	cmd.FS.BoolVar(&cmd.ToSQL, "to-sql", false, "Generate PostgreSQL CREATE TABLE statements")
	cmd.FS.BoolVar(&cmd.ToSQL, "to-schema", false, "Generate PostgreSQL CREATE TABLE statements (alias)")
	cmd.FS.BoolVar(&cmd.ToRepo, "to-repo", false, "Generate repository implementation from interface")
	cmd.FS.StringVar(&cmd.ModelFile, "model", "", "Model Go file with struct definitions (for --to-repo)")
	cmd.FS.StringVar(&cmd.InterfaceFile, "interface", "", "Go file containing the repository interface (for --to-repo)")
	cmd.FS.StringVar(&cmd.OutputFile, "o", "", "Output file")
	cmd.FS.StringVar(&cmd.OutputFile, "output", "", "Output file")
	cmd.FS.BoolVar(&cmd.ShowVersion, "version", false, "Show version")
	cmd.FS.BoolVar(&cmd.ShowVersion, "v", false, "Show version (shorthand)")

	return cmd
}

func (c *Command) Parse(args []string) error {
	return c.FS.Parse(args[1:])
}

func (c *Command) Validate() error {
	if c.ToRepo {
		if c.ModelFile == "" {
			return fmt.Errorf("--to-repo requires --model")
		}
		if c.InterfaceFile == "" {
			return fmt.Errorf("--to-repo requires --interface")
		}
		return nil
	}
	if !c.ToSQL {
		fmt.Fprintln(os.Stderr, "No output flag specified. Use one of:")
		fmt.Fprintln(os.Stderr, "  --to-sql       Generate PostgreSQL schema")
		fmt.Fprintln(os.Stderr, "  --to-repo      Generate repository implementation")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Parsing and validating structs only (no output generated)")
		return fmt.Errorf("no output flag specified")
	}
	return nil
}

type App struct {
	cmd           *Command
	version       string
	queryHandler  *query.Handler
	cmdHandler    *command.Handler
	parserWrapper *application.ParserWrapper
}

func New(version string) *App {
	cmd := NewCommand()

	parserWrapper := application.NewParserWrapper()
	queryHandler := query.NewHandler(parserWrapper)

	compositeGen := generator.NewCompositeGenerator()
	cmdHandler := command.NewHandler(compositeGen)

	return &App{
		cmd:           cmd,
		version:       version,
		queryHandler:  queryHandler,
		cmdHandler:    cmdHandler,
		parserWrapper: parserWrapper,
	}
}

func (a *App) Run(args []string) error {
	if len(args) < 2 {
		a.printUsage()
		return fmt.Errorf("no arguments specified")
	}

	if err := a.cmd.Parse(args); err != nil {
		return err
	}

	if a.cmd.ShowVersion {
		fmt.Println("structify version", a.version)
		return nil
	}

	if err := a.cmd.Validate(); err != nil {
		return err
	}

	ctx := context.Background()

	if a.cmd.ToRepo {
		return a.runRepoGeneration(ctx)
	}

	// Standard struct-based flow (--to-sql, or parse-only)
	inputFiles := a.cmd.FS.Args()

	parseQuery := &query.ParseQuery{Files: inputFiles}
	parseResult, err := a.queryHandler.Parse(ctx, parseQuery)
	if err != nil {
		return err
	}

	var output string
	if a.cmd.ToSQL {
		if parseResult.Count == 0 {
			return fmt.Errorf("no structs found")
		}
		cmd := &command.GenerateSchemaCommand{Entities: parseResult.EntityList}
		output, err = a.cmdHandler.GenerateSchema(ctx, cmd)
		if err != nil {
			return err
		}
	} else {
		fmt.Fprintf(os.Stderr, "Found %d struct(s):\n", parseResult.Count)
		for _, ent := range parseResult.EntityList {
			fmt.Fprintf(os.Stderr, "  %s\n", ent.Name)
		}
		return nil
	}

	return a.writeOutput(output, a.cmd.OutputFile)
}

func (a *App) runRepoGeneration(ctx context.Context) error {
	// 1. Parse model files → entities
	parseQuery := &query.ParseQuery{Files: []string{a.cmd.ModelFile}}
	parseResult, err := a.queryHandler.Parse(ctx, parseQuery)
	if err != nil {
		return fmt.Errorf("parse model: %w", err)
	}
	if parseResult.Count == 0 {
		return fmt.Errorf("no structs found in model file %s", a.cmd.ModelFile)
	}

	// Use first entity as the target entity
	ent := parseResult.EntityList[0]

	// 2. Parse interface file → interfaces (bound to entity)
	repos, err := a.parserWrapper.ParseInterfaces(ctx, []string{a.cmd.InterfaceFile}, ent)
	if err != nil {
		return fmt.Errorf("parse interface: %w", err)
	}
	if len(repos) == 0 {
		return fmt.Errorf("no interfaces found in %s", a.cmd.InterfaceFile)
	}

	repo := repos[0]

	// 3. Determine package name
	pkgName := parseResult.Package
	if pkgName == "" {
		pkgName = "repository"
	}

	// 4. Generate
	cmd := &command.GenerateRepoCommand{
		Entity:      ent,
		Interface:   repo,
		PackageName: pkgName,
	}
	output, err := a.cmdHandler.GenerateRepository(ctx, cmd)
	if err != nil {
		return err
	}

	// 5. Determine output file (default to interface file directory with .gen.go suffix)
	outputFile := a.cmd.OutputFile
	if outputFile == "" {
		outputFile = a.getDefaultOutputFile(a.cmd.InterfaceFile)
	}

	return a.writeOutput(output, outputFile)
}

func (a *App) writeOutput(output string, outputFile string) error {
	if outputFile != "" {
		return os.WriteFile(outputFile, []byte(output), 0600)
	}
	_, err := fmt.Fprint(os.Stdout, output)
	return err
}

// getDefaultOutputFile generates a default output filename based on the interface file path.
// For example: "./repo/user_repository.go" → "./repo/user_repository.gen.go"
func (a *App) getDefaultOutputFile(interfaceFile string) string {
	// Get the directory and base filename
	dir := ""
	baseName := interfaceFile

	if lastSlash := strings.LastIndex(interfaceFile, "/"); lastSlash != -1 {
		dir = interfaceFile[:lastSlash+1]
		baseName = interfaceFile[lastSlash+1:]
	} else if lastSlash := strings.LastIndex(interfaceFile, "\\"); lastSlash != -1 {
		dir = interfaceFile[:lastSlash+1]
		baseName = interfaceFile[lastSlash+1:]
	}

	// Remove .go extension and add .gen.go
	if strings.HasSuffix(baseName, ".go") {
		baseName = strings.TrimSuffix(baseName, ".go") + ".gen.go"
	} else {
		baseName = baseName + ".gen.go"
	}

	return dir + baseName
}

func (a *App) printUsage() {
	fmt.Fprintln(os.Stderr, "structify - Go struct to PostgreSQL schema converter")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  structify [flags] <input-files...>")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Flags:")
	fmt.Fprintln(os.Stderr, "  --to-sql, --to-schema")
	fmt.Fprintln(os.Stderr, "        Generate PostgreSQL CREATE TABLE statements")
	fmt.Fprintln(os.Stderr, "  --to-repo --model <file> --interface <file>")
	fmt.Fprintln(os.Stderr, "        Generate repository implementation from interface")
	fmt.Fprintln(os.Stderr, "  --output, -o <file>")
	fmt.Fprintln(os.Stderr, "        Output file (default: stdout)")
	fmt.Fprintln(os.Stderr, "  --version, -v")
	fmt.Fprintln(os.Stderr, "        Show version")
	fmt.Fprintln(os.Stderr, "  --help")
	fmt.Fprintln(os.Stderr, "        Show this help")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintln(os.Stderr, "  structify --to-sql ./models/user.go")
	fmt.Fprintln(os.Stderr, "  structify --to-repo --model ./models/user.go --interface ./repo/user_repo.go")
	fmt.Fprintln(os.Stderr, "  structify --to-repo --model ./models/user.go --interface ./repo/user_repo.go -o ./repo/user_repo.gen.go")
	fmt.Fprintln(os.Stderr, "")
}
