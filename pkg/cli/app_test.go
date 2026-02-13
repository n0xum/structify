package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	app := New("1.0.0")
	if app == nil {
		t.Fatal("New() returned nil")
	}
}

func TestAppRunNoArgs(t *testing.T) {
	app := New("1.0.0")
	err := app.Run([]string{"structify"})
	if err == nil {
		t.Error("Run() with no args should return error")
	}
}

func TestAppRunVersion(t *testing.T) {
	app := New("1.2.3")
	err := app.Run([]string{"structify", "--version"})
	if err != nil {
		t.Errorf("Run() --version error = %v", err)
	}
}

func TestAppRunNoFlag(t *testing.T) {
	app := New("1.0.0")
	err := app.Run([]string{"structify", "somefile.go"})
	if err == nil {
		t.Error("Run() without output flag should return error")
	}
}

func TestAppRunToSQLNoStructs(t *testing.T) {
	// Create a temp file with no exported structs
	tmp, err := os.CreateTemp("", "empty_*.go")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	tmp.WriteString("package test\n")
	tmp.Close()

	app := New("1.0.0")
	err = app.Run([]string{"structify", "--to-sql", tmp.Name()})
	if err == nil {
		t.Error("Run() --to-sql with no structs should return error")
	}
}

func TestAppRunToDBCodeNoStructs(t *testing.T) {
	tmp, err := os.CreateTemp("", "empty_*.go")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	tmp.WriteString("package test\n")
	tmp.Close()

	app := New("1.0.0")
	err = app.Run([]string{"structify", "--to-db-sql", tmp.Name()})
	if err == nil {
		t.Error("Run() --to-db-sql with no structs should return error")
	}
}

func TestAppRunFromJSONNoFile(t *testing.T) {
	app := New("1.0.0")
	err := app.Run([]string{"structify", "--from-json", "somefile.go"})
	if err == nil {
		t.Error("Run() --from-json without --json-file should return error")
	}
}

func TestAppRunFromJSON(t *testing.T) {
	app := New("1.0.0")
	err := app.Run([]string{"structify", "--from-json", "--json-file", "test.json", "somefile.go"})
	// convertJSON not yet implemented, should return error
	if err == nil {
		t.Error("Run() --from-json should return error (not implemented)")
	}
}

func TestAppRunToSQLWithOutput(t *testing.T) {
	// Use the real fixture file
	fixture := "../../test/fixtures/user.go"
	if _, err := os.Stat(fixture); os.IsNotExist(err) {
		t.Skip("fixture not found")
	}

	tmp := filepath.Join(t.TempDir(), "out.sql")
	app := New("1.0.0")
	err := app.Run([]string{"structify", "--to-sql", "--output", tmp, fixture})
	if err != nil {
		t.Errorf("Run() --to-sql error = %v", err)
	}
	if _, err := os.Stat(tmp); os.IsNotExist(err) {
		t.Error("output file not created")
	}
}

func TestAppRunToDBCodeWithOutput(t *testing.T) {
	fixture := "../../test/fixtures/user.go"
	if _, err := os.Stat(fixture); os.IsNotExist(err) {
		t.Skip("fixture not found")
	}

	tmp := filepath.Join(t.TempDir(), "out.go")
	app := New("1.0.0")
	err := app.Run([]string{"structify", "--to-db-sql", "--output", tmp, fixture})
	if err != nil {
		t.Errorf("Run() --to-db-sql error = %v", err)
	}
}

func TestWriteOutputStdout(t *testing.T) {
	app := New("1.0.0")
	app.cmd.OutputFile = ""
	err := app.writeOutput("test output")
	if err != nil {
		t.Errorf("writeOutput() error = %v", err)
	}
}

func TestWriteOutputFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.txt")
	app := New("1.0.0")
	app.cmd.OutputFile = tmp
	err := app.writeOutput("hello world")
	if err != nil {
		t.Errorf("writeOutput() error = %v", err)
	}
	data, _ := os.ReadFile(tmp)
	if string(data) != "hello world" {
		t.Errorf("writeOutput() wrote %q, want %q", data, "hello world")
	}
}

func TestPrintUsage(t *testing.T) {
	app := New("1.0.0")
	// Should not panic
	app.printUsage()
}
