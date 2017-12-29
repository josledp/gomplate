package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

// input - models an input file...
type input struct {
	name     string
	target   io.Writer
	contents string
}

func (t *input) toGoTemplate(g *Gomplate) (*template.Template, error) {
	tmpl := template.New(t.name)
	tmpl.Option("missingkey=error")
	tmpl.Funcs(g.funcMap)
	tmpl.Delims(g.leftDelim, g.rightDelim)
	return tmpl.Parse(t.contents)
}

// gatherTemplates - gather and prepare input template(s) and output file(s) for rendering
func gatherTemplates(o *GomplateOpts) (ins []*input, err error) {
	// the arg-provided input string gets a special name
	if o.input != "" {
		o.inputFiles = []string{"<arg>"}
		ins = []*input{{
			name:     "<arg>",
			contents: o.input,
		}}
	}

	// input dirs presume output dirs are set too
	if o.inputDir != "" {
		excludes, err := executeCombinedGlob(o.excludeGlob)
		if err != nil {
			return nil, err
		}

		o.inputFiles, o.outputFiles, err = walkDir(o.inputDir, o.outputDir, excludes)
		if err != nil {
			return nil, err
		}
	}

	ins = make([]*input, len(o.inputFiles))

	if len(o.outputFiles) == 0 {
		o.outputFiles = []string{"-"}
	}

	for i, filename := range o.inputFiles {
		if ins[i] == nil {
			ins[i] = &input{}
		}
		if ins[i].name == "" {
			ins[i].name = filename
		}
		if ins[i].contents == "" {
			contents, err := readInput(filename)
			if err != nil {
				return nil, err
			}
			ins[i].contents = contents
		}
		if ins[i].target == nil {
			target, err := openOutFile(o.outputFiles[i])
			if err != nil {
				return nil, err
			}
			addCleanupHook(func() {
				// nolint: errcheck
				target.Close()
			})
			ins[i].target = target
		}
	}

	return ins, nil
}

func walkDir(dir, outDir string, excludes []string) ([]string, []string, error) {
	dir = filepath.Clean(dir)
	outDir = filepath.Clean(outDir)

	si, err := os.Stat(dir)
	if err != nil {
		return nil, nil, err
	}

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}

	if err = os.MkdirAll(outDir, si.Mode()); err != nil {
		return nil, nil, err
	}

	inFiles := []string{}
	outFiles := []string{}
	for _, entry := range entries {
		nextInPath := filepath.Join(dir, entry.Name())
		nextOutPath := filepath.Join(outDir, entry.Name())

		if inList(excludes, nextInPath) {
			continue
		}

		if entry.IsDir() {
			i, o, err := walkDir(nextInPath, nextOutPath, excludes)
			if err != nil {
				return nil, nil, err
			}
			inFiles = append(inFiles, i...)
			outFiles = append(outFiles, o...)
		} else {
			inFiles = append(inFiles, nextInPath)
			outFiles = append(outFiles, nextOutPath)
		}
	}
	return inFiles, outFiles, nil
}

func inList(list []string, entry string) bool {
	for _, file := range list {
		if file == entry {
			return true
		}
	}

	return false
}

func openOutFile(filename string) (out *os.File, err error) {
	if filename == "-" {
		return os.Stdout, nil
	}
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
}

func readInput(filename string) (string, error) {
	var err error
	var inFile *os.File
	if filename == "-" {
		inFile = os.Stdin
	} else {
		inFile, err = os.Open(filename)
		if err != nil {
			return "", fmt.Errorf("failed to open %s\n%v", filename, err)
		}
		// nolint: errcheck
		defer inFile.Close()
	}
	bytes, err := ioutil.ReadAll(inFile)
	if err != nil {
		err = fmt.Errorf("read failed for %s\n%v", filename, err)
		return "", err
	}
	return string(bytes), nil
}

// takes an array of glob strings and executes it as a whole,
// returning a merged list of globbed files
func executeCombinedGlob(globArray []string) ([]string, error) {
	var combinedExcludes []string
	for _, glob := range globArray {
		excludeList, err := filepath.Glob(glob)
		if err != nil {
			return nil, err
		}

		combinedExcludes = append(combinedExcludes, excludeList...)
	}

	return combinedExcludes, nil
}
