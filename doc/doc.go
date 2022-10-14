// Package doc documents of the Go package structure as a dependency table.
package doc

import (
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/flowdev/spaghetti-analyzer/analdata"
	"github.com/flowdev/spaghetti-analyzer/x/table"
	"github.com/flowdev/spaghetti-cutter/data"
)

const (
	// FileName is the name of the documentation file (package_dependencies.md)
	FileName = "package_dependencies.md"
	// Title is the mark down title of the package dependencies
	Title = `# Dependency Table For: `
	// Legend is the explanation for the dependency table
	Legend = `
### Legend

* Rows - Importing packages
* Columns - Imported packages


#### Meaning Of Row And Row Header Formatting

* **Bold** - God package (can use all packages)
` + "* `Code` - Database package (can only use tool and other database packages)" + `
* _Italic_ - Tool package (foundational, no dependencies)
* No formatting - Standard package (can only use tool and database packages)


#### Meaning Of Letters In Table Columns

* G - God package (can use all packages)
* D - Database package (can only use tool and other database packages)
* T - Tool package (foundational, no dependencies)
* S - Standard package (can only use tool and database packages)
`
)

// WriteDocs generates documentation for the packages 'dtPkgs' and writes it to
// files.
// If dtPkgs contains more than 1 package, the others will be used as link targets
// instead of reporting all the details in one table.
func WriteDocs(
	dtPkgs []string,
	depMap analdata.DependencyMap,
	rootPkg, root string,
) {
	pl, err := data.NewSimplePatternList(dtPkgs, "--doc")
	if err != nil {
		log.Printf("ERROR - Unable to create dependency tables: %v", err)
		return
	}

	docFiles := docFilesForPkgs(dtPkgs)
	for i := range pl {
		writeDoc(i, pl, docFiles, depMap, rootPkg, root)
	}
}
func docFilesForPkgs(pkgs []string) []string {
	files := make([]string, len(pkgs))
	for i, p := range pkgs {
		docFile := filepath.Join(analdata.PkgForPattern(p), FileName)
		files[i] = docFile
	}
	return files
}

func writeDoc(
	idx int,
	links data.PatternList,
	docFiles []string,
	depMap analdata.DependencyMap,
	rootPkg, root string,
) {
	doc := generateTable(idx, links, docFiles, depMap, rootPkg)
	if doc == "" {
		return
	}
	docFile := docFiles[idx]
	log.Printf("INFO - Write dependency table to file: %s", docFile)
	docFile = filepath.Join(root, docFile)
	err := ioutil.WriteFile(docFile, []byte(doc), 0644)
	if err != nil {
		log.Printf("ERROR - Unable to write dependency table to file %s: %v", docFile, err)
	}
}

// generateTable generates the dependency matrix for the idx package from links.
func generateTable(
	idx int,
	links data.PatternList,
	docFiles []string,
	depMap analdata.DependencyMap,
	rootPkg string,
) string {
	startPkg := filepath.ToSlash(filepath.Dir(docFiles[idx]))
	pattern := links[idx].Pattern
	depMap = analdata.FilterDepMap(depMap, idx, links)
	if len(depMap) == 0 {
		log.Printf("INFO - Won't write depenency table for package %q because it has no dependencies.", pattern)
		return ""
	}
	allRows := make([]string, 0, len(depMap))
	allCols := make([]string, 0, len(depMap)*2)
	allColsMap := make(map[string]data.PkgType, len(depMap)*2)

	for pkg, pkgImps := range depMap {
		allRows = append(allRows, pkg)
		for impName, impType := range pkgImps.Imports {
			if _, ok := allColsMap[impName]; !ok {
				allColsMap[impName] = impType
				allCols = append(allCols, impName)
			}
		}
	}

	sort.Strings(allRows)
	sort.Strings(allCols)

	sb := &strings.Builder{}
	sb.WriteString(Title + path.Join(rootPkg, pattern) + "\n\n")

	tableData := make([][]string, 0, len(allRows)+1)
	headerRow := make([]string, 0, len(allCols)+1)
	tableAlign := make([]table.Align, 0, len(allCols)+1)

	// (column) header line: | | C o l 1 - G | C o l 2 | ... | C o l N - T |
	headerRow = append(headerRow, "")
	for _, col := range allCols {
		sb2 := &strings.Builder{}
		colIdx := data.DocMatchStringIndex(col, links)
		if colIdx >= 0 && colIdx != idx {
			sb2.WriteRune('[')
		}
		for _, r := range col {
			sb2.WriteRune(r)
			sb2.WriteRune(' ')
		}
		sb2.WriteString("- ")
		letter := data.TypeLetter(allColsMap[col])
		sb2.WriteRune(letter)
		if colIdx >= 0 && colIdx != idx {
			sb2.WriteString("](")
			sb2.WriteString(RelPath(startPkg, filepath.ToSlash(docFiles[colIdx])))
			sb2.WriteString(")")
		}
		headerRow = append(headerRow, sb2.String())
	}
	tableData = append(tableData, headerRow)

	// separator line: |:---|:--:|:-:| ... |:---:|
	tableAlign = append(tableAlign, table.AlignLeft)
	for range allCols {
		tableAlign = append(tableAlign, table.AlignCenter)
	}

	// normal rows: | **Row1** | **G** | | ... | **T** |
	for _, row := range allRows {
		dataRow := make([]string, 0, len(allCols)+1)
		sb2 := &strings.Builder{}
		pkgImps := depMap[row]

		format := data.TypeFormat(pkgImps.PkgType)
		sb2.WriteString(format)
		sb2.WriteString(row)
		sb2.WriteString(format)
		dataRow = append(dataRow, sb2.String())

		for _, col := range allCols {
			sb2.Reset()
			if impType, ok := pkgImps.Imports[col]; ok {
				sb2.WriteString(format)
				sb2.WriteRune(data.TypeLetter(impType))
				sb2.WriteString(format)
			}
			dataRow = append(dataRow, sb2.String())
		}
		tableData = append(tableData, dataRow)
	}

	sb.WriteString(table.Generate(tableData, tableAlign))
	sb.WriteString(Legend)
	return sb.String()
}

// RelPath calculates the relative path from 'basepath' to 'targetpath'.
// This is very similar to filepath.Rel() but not OS specific but it is working
// by purely lexical processing like the path package.
func RelPath(basepath, targetpath string) string {
	base := splitPath(path.Clean(basepath))
	targ := splitPath(path.Clean(targetpath))

	n := len(base)
	m := len(targ)
	i := 0
	for i < n && i < m && base[i] == targ[i] {
		i++
	}

	ret := ""
	for j := i; j < n; j++ { // go backward from base
		ret = path.Join(ret, "..")
	}
	for j := i; j < m; j++ { // go forward to target
		ret = path.Join(ret, targ[j])
	}

	return ret
}
func splitPath(p string) []string {
	ret := make([]string, 0, 64)
	for p != "" {
		base, last := path.Split(p)
		ret = append(ret, last)
		p = removeTrailingSlash(base)
	}
	return reverse(ret)
}
func removeTrailingSlash(p string) string {
	if !strings.HasSuffix(p, "/") {
		return p
	}
	return p[:len(p)-1]
}
func reverse(ss []string) []string {
	n := len(ss)
	ts := make([]string, n)
	n--
	for i, s := range ss {
		ts[n-i] = s
	}
	return ts
}
