package dirs

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ValidateRoot validates the root of a project.
// It looks for the configuration file '.spaghetti-cutter.hjson' in the 'root'
// directory and returns it's absolute path or the empty string if not found.
func ValidateRoot(root, cfgFile string) (string, error) {
	if root == "" {
		root = "."
	}

	absDir, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("unable to find absolute directory (for %q): %w", root, err)
	}

	p := filepath.Join(absDir, cfgFile)
	if _, err = os.Stat(p); err == nil {
		return absDir, nil
	}

	return "", nil
}

// FindDepTables is finding packages containing a dependency table on disk
// starting at 'root' and adding them to those given in 'startPkgs'.
func FindDepTables(file, title string, startPkgs []string, root, rootPkg string) map[string]struct{} {
	val := struct{}{}
	// prefill doc packages from startPkgs
	retPkgs := make(map[string]struct{}, 128)
	for _, p := range startPkgs {
		retPkgs[p] = val
	}

	// walk the file system to find more 'file's
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() { // we are only interested in directories
			return nil
		}
		if err != nil {
			log.Printf("WARN - Unable to list directory %q: %v", path, err)
			return filepath.SkipDir
		}

		// no valid package starts with '.' and we don't want to search in testdata
		if strings.HasPrefix(info.Name(), ".") || info.Name() == "testdata" {
			return filepath.SkipDir
		}

		depFile := filepath.Join(path, file)
		if _, err := os.Lstat(depFile); err == nil {
			pkg, err := filepath.Rel(root, path)
			if err != nil {
				log.Printf("WARN - Unable to compute package for %q: %v", path, err)
				return nil // sub-directories might work
			}
			pattern, err := readPatternFromFile(depFile, title, rootPkg)
			if err != nil {
				log.Printf("WARN - Problem reading pattern from file %q: %v", depFile, err)
				err = nil
			}
			if pattern == "" {
				pattern = strings.ReplaceAll(pkg, "\\", "/") // packages like URLs have always '/'s
			}
			if pattern == "." {
				retPkgs["/"] = val
			} else {
				retPkgs[pattern] = val
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("ERROR - Unable to walk the path %q: %v", root, err)
	}
	return retPkgs
}

func readPatternFromFile(depFile, prefix, rootPkg string) (string, error) {
	lines, err := readFirstLines(depFile, 15)
	prefix = strings.ToLower(prefix)
	for _, l := range lines {
		if strings.HasPrefix(strings.ToLower(l), prefix) {
			pattern := l[len(prefix):]
			pattern = strings.TrimSpace(pattern)
			if strings.HasPrefix(pattern, rootPkg) {
				pattern = pattern[len(rootPkg):]
				if pattern != "" && pattern[0] == '/' {
					pattern = pattern[1:]
				}
				if pattern == "" {
					pattern = "/"
				}
			}
			return pattern, err
		}
	}
	return "", err
}

func readFirstLines(fileName string, n int) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0, n)
	scanner := bufio.NewScanner(file)
	for i := 0; i < n && scanner.Scan(); i++ {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return lines, err
	}
	return lines, nil
}
