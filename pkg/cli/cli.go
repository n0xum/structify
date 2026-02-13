package cli

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/n0xum/structify/internal/application"
	"github.com/n0xum/structify/internal/application/command"
	"github.com/n0xum/structify/internal/application/query"
	"github.com/n0xum/structify/internal/generator"
)

type Command struct {
	FS          *flag.FlagSet
	ToSQL       bool
	ToDBCode    bool
	FromJSON    bool
	JSONFile    string
	OutputFile  string
	ShowVersion bool
}

func NewCommand() *Command {
	cmd := &Command{
		FS: flag.NewFlagSet("structify", flag.ContinueOnError),
	}

	cmd.FS.BoolVar(&cmd.ToSQL, "to-sql", false, "Generate PostgreSQL CREATE TABLE statements")
	cmd.FS.BoolVar(&cmd.ToSQL, "to-schema", false, "Generate PostgreSQL CREATE TABLE statements (alias)")
	cmd.FS.BoolVar(&cmd.ToDBCode, "to-db-sql", false, "Generate database/sql CRUD code")
	cmd.FS.BoolVar(&cmd.ToDBCode, "to-dbcode", false, "Generate database/sql CRUD code (alias)")
	cmd.FS.BoolVar(&cmd.FromJSON, "from-json", false, "Convert JSON to Go struct")
	cmd.FS.StringVar(&cmd.JSONFile, "json-file", "", "JSON input file")
	cmd.FS.StringVar(&cmd.JSONFile, "f", "", "JSON input file (shorthand)")
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
	if c.FromJSON && c.JSONFile == "" {
		return fmt.Errorf("--from-json requires --json-file")
	}
	if !c.ToSQL && !c.ToDBCode && !c.FromJSON {
		fmt.Fprintln(os.Stderr, "No output flag specified. Use one of:")
		fmt.Fprintln(os.Stderr, "  --to-sql      Generate PostgreSQL schema")
		fmt.Fprintln(os.Stderr, "  --to-db-sql    Generate database/sql code")
		fmt.Fprintln(os.Stderr, "  --from-json    Convert JSON to Go struct")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Parsing and validating structs only (no output generated)")
		return fmt.Errorf("no output flag specified")
	}
	return nil
}

type App struct {
	cmd          *Command
	version      string
	queryHandler *query.Handler
	cmdHandler   *command.Handler
}

func New(version string) *App {
	cmd := NewCommand()

	parserWrapper := application.NewParserWrapper()
	queryHandler := query.NewHandler(parserWrapper)

	compositeGen := generator.NewCompositeGenerator()
	cmdHandler := command.NewHandler(compositeGen)

	return &App{
		cmd:          cmd,
		version:      version,
		queryHandler: queryHandler,
		cmdHandler:   cmdHandler,
	}
}

func (a *App) Run(args []string) error {
	if len(args) < 2 {
		a.printUsage()
		return fmt.Errorf("no input files specified")
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

	inputFiles := a.cmd.FS.Args()

	if a.cmd.FromJSON {
		return a.convertJSON()
	}

	ctx := context.Background()

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
	} else if a.cmd.ToDBCode {
		if parseResult.Count == 0 {
			return fmt.Errorf("no structs found")
		}
		pkgName := parseResult.Package
		if pkgName == "" {
			pkgName = "models"
		}
		cmd := &command.GenerateSchemaCommand{PackageName: pkgName, Entities: parseResult.EntityList}
		output, err = a.cmdHandler.GenerateCode(ctx, cmd)
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

	return a.writeOutput(output)
}

func (a *App) convertJSON() error {
	return fmt.Errorf("JSON conversion not yet implemented")
}

func (a *App) writeOutput(output string) error {
	if a.cmd.OutputFile != "" {
		return os.WriteFile(a.cmd.OutputFile, []byte(output), 0600)
	}
	_, err := fmt.Fprint(os.Stdout, output)
	return err
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
	fmt.Fprintln(os.Stderr, "  --to-db-sql, --to-dbcode")
	fmt.Fprintln(os.Stderr, "        Generate database/sql CRUD code")
	fmt.Fprintln(os.Stderr, "  --from-json")
	fmt.Fprintln(os.Stderr, "        Convert JSON to Go struct")
	fmt.Fprintln(os.Stderr, "  --json-file, -f <file>")
	fmt.Fprintln(os.Stderr, "        JSON input file for --from-json")
	fmt.Fprintln(os.Stderr, "  --output, -o <file>")
	fmt.Fprintln(os.Stderr, "        Output file (default: stdout)")
	fmt.Fprintln(os.Stderr, "  --version, -v")
	fmt.Fprintln(os.Stderr, "        Show version")
	fmt.Fprintln(os.Stderr, "  --help")
	fmt.Fprintln(os.Stderr, "        Show this help")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintln(os.Stderr, "  structify --to-sql ./models/user.go")
	fmt.Fprintln(os.Stderr, "  structify --to-db-sql ./models/*.go -o user_repo.go")
	fmt.Fprintln(os.Stderr, "")
}
