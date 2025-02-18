// SPDX-FileCopyrightText: Copyright 2025 The OSPS Authors
// SPDX-License-Identifier: Apache-2.0

package baseline

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ossf/security-baseline/pkg/types"
	"gopkg.in/yaml.v3"
)

const LexiconFilename = "lexicon.yaml"

// Loader is an object that reads the baseline data
type Loader struct {
	DataPath string
}

func NewLoader() *Loader {
	return &Loader{}
}

// Load reads the baseline data and returns the representation types
func (l *Loader) Load() (*types.Baseline, error) {
	b := &types.Baseline{
		Categories: make(map[string]types.Category, len(types.Categories)),
	}

	// Load the lexicon:
	lexicon, err := l.loadLexicon()
	if err != nil {
		return nil, fmt.Errorf("error reading lexicon: %w", err)
	}
	b.Lexicon = lexicon

	for _, catCode := range types.Categories {
		cat, err := l.loadCategory(catCode)
		if err != nil {
			return nil, fmt.Errorf("loading category %q: %w", catCode, err)
		}
		b.Categories[catCode] = *cat
	}

	// return b, b.validate()
	return b, nil
}

// loadLexicon
func (l *Loader) loadLexicon() ([]types.LexiconEntry, error) {
	file, err := os.Open(filepath.Join(l.DataPath, LexiconFilename))
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var lexicon []types.LexiconEntry

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	if err := decoder.Decode(&lexicon); err != nil {
		return nil, fmt.Errorf("error decoding YAML: %v", err)
	}
	return lexicon, nil
}

// loadCategory loads a category definition from its YAML source
func (l *Loader) loadCategory(catCode string) (*types.Category, error) {
	file, err := os.Open(filepath.Join(l.DataPath, fmt.Sprintf("OSPS-%s.yaml", catCode)))
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var category = &types.Category{}

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	if err := decoder.Decode(category); err != nil {
		return nil, fmt.Errorf("error decoding %s YAML: %w", catCode, err)
	}
	return category, nil
}
