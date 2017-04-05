package main

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type GeneratedType struct {
	Name      string
	LowerName string
}

func getRenderedPath(suffix, inputPath string) (string, error) {
	if !strings.HasSuffix(inputPath, ".go") {
		return "", fmt.Errorf("Input path %s doesn't have .go extension", inputPath)
	}
	trimmed := strings.TrimSuffix(inputPath, ".go")
	dir, file := filepath.Split(trimmed)
	return filepath.Join(dir, fmt.Sprintf("%s_%s.go", file, suffix)), nil
}

type generateTemplateData struct {
	Package string
	Types   []GeneratedType
}

func render(suffix string, w io.Writer, packageName string, types []GeneratedType) error {

	switch suffix {
	case "ponzi":
		return ponziTmpl.Execute(w, generateTemplateData{packageName, types})
	}
	return errors.New("Unknown template")
}
