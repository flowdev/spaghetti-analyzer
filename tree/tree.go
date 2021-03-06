// Package tree creates a directory tree with information about Go packages.
package tree

import (
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/flowdev/spaghetti-analyzer/x/dirs"
	"github.com/flowdev/spaghetti-analyzer/x/pkgs"
)

const (
	File         = "dirtree.txt"
	newLine      = "\n"
	separator    = " -\t"
	emptyItem    = "    "
	middleItem   = "├── "
	continueItem = "│   "
	lastItem     = "└── "
)

var (
	reSpaces   = regexp.MustCompile(`[\s]+`)
	reFullStop = regexp.MustCompile(`([.:?!])(?s:\s.*)?$`)
)

func Generate(root, name string, packs []*pkgs.Package) (string, error) {
	sb := &strings.Builder{}

	err := generateTree(root, name, sb, "", "", packs)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func generateTree(root, name string, sb *strings.Builder, prefix string, pkg string, packs []*pkgs.Package) error {
	sb.WriteString(name)
	sb.WriteString(separator)
	sb.WriteString(docForPkg(pkg, name, packs))
	sb.WriteString(newLine)

	files, err := os.ReadDir(root) // files is already sorted :)
	if err != nil {
		log.Printf("ERROR - Unable to read the directory %q: %v", root, err)
		return err
	}

	// reduce all files to only the items we want to include
	var items []fs.DirEntry
	for _, file := range files {
		if file.IsDir() && dirs.IncludeFile(file.Name()) {
			items = append(items, file)
		}
	}

	lastI := len(items) - 1
	p := ""
	for i, item := range items {
		if pkg == "" {
			p = item.Name()
		} else {
			p = path.Join(pkg, item.Name())
		}

		pref := prefix
		if i == lastI {
			sb.WriteString(prefix + lastItem)
			pref += emptyItem
		} else {
			sb.WriteString(prefix + middleItem)
			pref += continueItem
		}
		err = generateTree(filepath.Join(root, item.Name()), item.Name(), sb, pref, p, packs)
		if err != nil {
			return err
		}
	}
	return nil
}

func docForPkg(pkg, name string, packs []*pkgs.Package) string {
	if pkg == "" {
		pkg = name
	}
	for _, p := range packs {
		if strings.HasSuffix(p.PkgPath, pkg) {
			for _, f := range p.Syntax {
				if f.Doc != nil {
					return firstSentenceOf(f.Doc.Text())
				}
			}
			return "Package " + name + " ..."
		}
	}
	return ""
}

func firstSentenceOf(text string) string {
	text = reSpaces.ReplaceAllString(text, " ")    // replace multiple spaces, tabs and new lines with a single space
	return reFullStop.ReplaceAllString(text, "$1") // cut after first '. ', ': ', '? ' or '! '
}
