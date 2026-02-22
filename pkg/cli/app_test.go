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

func TestAppRunToSQLWithOutput(t *testing.T) {
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

func TestAppRunToRepoMissingModel(t *testing.T) {
	app := New("1.0.0")
	err := app.Run([]string{"structify", "--to-repo", "--interface", "repo.go"})
	if err == nil {
		t.Error("Run() --to-repo without --model should return error")
	}
}

func TestAppRunToRepoMissingInterface(t *testing.T) {
	app := New("1.0.0")
	err := app.Run([]string{"structify", "--to-repo", "--model", "model.go"})
	if err == nil {
		t.Error("Run() --to-repo without --interface should return error")
	}
}

func TestAppRunToRepoSuccess(t *testing.T) {
	modelFixture := "../../test/fixtures/user.go"
	ifaceFixture := "../../test/fixtures/user_repository.go"
	if _, err := os.Stat(modelFixture); os.IsNotExist(err) {
		t.Skip("model fixture not found")
	}
	if _, err := os.Stat(ifaceFixture); os.IsNotExist(err) {
		t.Skip("interface fixture not found")
	}

	tmp := filepath.Join(t.TempDir(), "out.gen.go")
	app := New("1.0.0")
	err := app.Run([]string{"structify", "--to-repo", "--model", modelFixture, "--interface", ifaceFixture, "-o", tmp})
	if err != nil {
		t.Errorf("Run() --to-repo error = %v", err)
	}
	if _, err := os.Stat(tmp); os.IsNotExist(err) {
		t.Error("output file not created")
	}
}

func TestAppRunToRepoNoInterfaces(t *testing.T) {
	modelFixture := "../../test/fixtures/user.go"
	if _, err := os.Stat(modelFixture); os.IsNotExist(err) {
		t.Skip("fixture not found")
	}

	// Use model file as interface file (no interfaces in it)
	app := New("1.0.0")
	err := app.Run([]string{"structify", "--to-repo", "--model", modelFixture, "--interface", modelFixture})
	if err == nil {
		t.Error("Run() --to-repo with no interfaces should return error")
	}
}

func TestAppRunToRepoModelParseError(t *testing.T) {
	app := New("1.0.0")
	err := app.Run([]string{"structify", "--to-repo", "--model", "nonexistent.go", "--interface", "nonexistent.go"})
	if err == nil {
		t.Error("Run() --to-repo with nonexistent model should return error")
	}
}

func TestWriteOutputStdout(t *testing.T) {
	app := New("1.0.0")
	app.cmd.OutputFile = ""
	err := app.writeOutput("test output", "")
	if err != nil {
		t.Errorf("writeOutput() error = %v", err)
	}
}

func TestWriteOutputFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.txt")
	app := New("1.0.0")
	app.cmd.OutputFile = tmp
	err := app.writeOutput("hello world", tmp)
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
	app.printUsage()
}

func TestAppRunToRepoInterfaceParseError(t *testing.T) {
	modelFixture := "../../test/fixtures/user.go"
	if _, err := os.Stat(modelFixture); os.IsNotExist(err) {
		t.Skip("model fixture not found")
	}

	app := New("1.0.0")
	err := app.Run([]string{"structify", "--to-repo", "--model", modelFixture, "--interface", "nonexistent.go"})
	if err == nil {
		t.Error("Run() --to-repo with nonexistent interface should return error")
	}
}

func TestGetDefaultOutputFile(t *testing.T) {
	tests := []struct {
		name          string
		interfaceFile string
		expected      string
	}{
		{
			name:          "simple path",
			interfaceFile: "repo/user_repository.go",
			expected:      "repo/user_repository.gen.go",
		},
		{
			name:          "relative path with ./",
			interfaceFile: "./repo/user_repository.go",
			expected:      "./repo/user_repository.gen.go",
		},
		{
			name:          "absolute path",
			interfaceFile: "/abs/path/repo/user_repository.go",
			expected:      "/abs/path/repo/user_repository.gen.go",
		},
		{
			name:          "path with directory traversal",
			interfaceFile: "../repo/user_repository.go",
			expected:      "../repo/user_repository.gen.go",
		},
		{
			name:          "Windows-style path",
			interfaceFile: "repo\\user_repository.go",
			expected:      "repo\\user_repository.gen.go",
		},
	}

	app := New("1.0.0")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := app.getDefaultOutputFile(tt.interfaceFile)
			if got != tt.expected {
				t.Errorf("getDefaultOutputFile() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAppRunInvalidFlag(t *testing.T) {
	app := New("1.0.0")
	err := app.Run([]string{"structify", "--invalid-flag"})
	if err == nil {
		t.Error("Run() with invalid flag should return error")
	}
}
