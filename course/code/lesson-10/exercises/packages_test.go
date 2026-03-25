package exercises

import (
	"testing"
)

func TestIsInitialized(t *testing.T) {
	if !IsInitialized() {
		t.Error("Package should be initialized")
	}
}

func TestGetInitCount(t *testing.T) {
	count := GetInitCount()
	if count < 1 {
		t.Errorf("Expected initCount >= 1, got %d", count)
	}
}

func TestUseAliases(t *testing.T) {
	result := UseAliases()
	if result != "hello" {
		t.Errorf("Expected 'hello', got '%s'", result)
	}
}

func TestUseExternal(t *testing.T) {
	result := UseExternal()
	if result == "" {
		t.Error("Expected non-empty string")
	}
}

func TestPackageDocumentation(t *testing.T) {
	pd := NewPackageDocumentation("mypackage", "A test package", "v1.0.0")
	
	if pd.Name != "mypackage" {
		t.Errorf("Expected name 'mypackage', got '%s'", pd.Name)
	}
	
	if pd.Version != "v1.0.0" {
		t.Errorf("Expected version 'v1.0.0', got '%s'", pd.Version)
	}
	
	expected := "# mypackage\n\nA test package\n\nVersion: v1.0.0"
	if pd.String() != expected {
		t.Errorf("Expected:\n%s\ngot:\n%s", expected, pd.String())
	}
}

func TestBlankImport(t *testing.T) {
	result := BlankImport()
	if result == "" {
		t.Error("Expected non-empty string")
	}
}

func TestVendorDirectory(t *testing.T) {
	result := VendorDirectory()
	if result == "" {
		t.Error("Expected non-empty string")
	}
}

func TestModuleCommands(t *testing.T) {
	result := ModuleCommands()
	if result == "" {
		t.Error("Expected non-empty string")
	}
}
