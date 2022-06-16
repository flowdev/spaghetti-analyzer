package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/flowdev/spaghetti-analyzer/analdata"
	"github.com/flowdev/spaghetti-analyzer/deps"
	"github.com/flowdev/spaghetti-analyzer/doc"
	"github.com/flowdev/spaghetti-analyzer/parse"
	"github.com/flowdev/spaghetti-analyzer/stat"
	"github.com/flowdev/spaghetti-analyzer/tree"
	"github.com/flowdev/spaghetti-analyzer/x/dirs"
	"github.com/flowdev/spaghetti-analyzer/x/pkgs"
	"github.com/flowdev/spaghetti-cutter/config"
	"github.com/flowdev/spaghetti-cutter/data"
)

func main() {
	rc := analyze(os.Args[1:])
	if rc != 0 {
		os.Exit(rc)
	}
}

func analyze(args []string) int {
	const (
		usageShort     = " (shorthand)"
		defaultRoot    = "."
		usageRoot      = "root directory of the project (NO CRAWLING is done!)"
		defaultDoc     = "*"
		usageDoc       = "write '" + doc.FileName + "' for packages (separated by ','; '' for none)"
		defaultNoLinks = false
		usageNoLinks   = "don't use links in '" + doc.FileName + "' files"
		defaultStats   = false
		usageStats     = "write '" + stat.FileName + "' for project"
		defaultDirTree = false
		usageDirTree   = "write a directory tree (starting at 'root') to: " + tree.File
	)
	var startDir string
	var docPkgs string
	var noLinks bool
	var doStats bool
	var dirTree bool
	fs := flag.NewFlagSet("spaghetti-analyzer", flag.ExitOnError)
	fs.StringVar(&startDir, "root", defaultRoot, usageRoot)
	fs.StringVar(&startDir, "r", defaultRoot, usageRoot+usageShort)
	fs.StringVar(&docPkgs, "doc", defaultDoc, usageDoc)
	fs.StringVar(&docPkgs, "d", defaultDoc, usageDoc+usageShort)
	fs.BoolVar(&noLinks, "nolinks", defaultNoLinks, usageNoLinks)
	fs.BoolVar(&noLinks, "l", defaultNoLinks, usageNoLinks+usageShort)
	fs.BoolVar(&doStats, "stats", defaultStats, usageStats)
	fs.BoolVar(&doStats, "s", defaultStats, usageStats+usageShort)
	fs.BoolVar(&dirTree, "dirtree", defaultDirTree, usageDirTree)
	fs.BoolVar(&dirTree, "t", defaultDirTree, usageDirTree+usageShort)
	err := fs.Parse(args)
	if err != nil {
		log.Printf("FATAL - %v", err)
		return 2
	}

	root, err := dirs.ValidateRoot(startDir, config.File)
	if err != nil {
		log.Printf("FATAL - %v", err)
		return 3
	}
	cfgFile := filepath.Join(root, config.File)
	var cfg config.Config
	if root != "" {
		cfgBytes, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			log.Printf("FATAL - unable to read configuration file %q: %v", cfgFile, err)
			return 4
		}
		cfg, err = config.Parse(cfgBytes, cfgFile)
		if err != nil {
			log.Printf("FATAL - %v", err)
			return 5
		}

		log.Printf("INFO - configuration 'god': %s", cfg.God)
		log.Printf("INFO - configuration 'tool': %s", cfg.Tool)
		log.Printf("INFO - configuration 'db': %s", cfg.DB)
	} else {
		root, err = filepath.Abs(startDir)
		if err != nil {
			log.Printf("FATAL - %v", err)
			return 3
		}
		cfg.God, err = data.NewSimplePatternList([]string{"main"}, "god")
		if err != nil {
			log.Printf("FATAL - %v", err)
			return 5
		}
	}

	log.Printf("INFO - dependency tables for package(s): %s", docPkgs)
	log.Printf("INFO - no links in '"+doc.FileName+"' files: %t", noLinks)
	log.Printf("INFO - write dependency statistics: %t", doStats)
	log.Printf("INFO - write directory tree: %t", dirTree)

	packs, err := parse.DirTree(root)
	if err != nil {
		log.Printf("FATAL - %v", err)
		return 6
	}

	depMap := make(analdata.DependencyMap, 256)
	rootPkg := parse.RootPkg(packs)
	log.Printf("INFO - root package: %s", rootPkg)
	pkgInfos := pkgs.UniquePackages(packs)
	for _, pkgInfo := range pkgInfos {
		deps.Fill(pkgInfo.Pkg, rootPkg, cfg, &depMap)
	}

	if doStats {
		writeStatistics(root, depMap)
	} else {
		log.Print("INFO - No dependency statistics wanted.")
	}

	if docPkgs != "" {
		writeDepTables(docPkgs, root, rootPkg, noLinks, depMap)
	} else {
		log.Print("INFO - No dependency table wanted.")
	}

	if dirTree {
		err := writeDirTree(root, path.Base(rootPkg), packs)
		if err != nil {
			log.Printf("FATAL - %v", err)
			return 7
		}
	} else {
		log.Print("INFO - No directory tree wanted.")
	}

	return 0
}

func writeStatistics(root string, depMap analdata.DependencyMap) {
	log.Printf("INFO - Writing package statistics to file: %s", stat.FileName)
	statMD := stat.Generate(depMap)
	if statMD == "" {
		return
	}
	statFile := filepath.Join(root, stat.FileName)
	err := ioutil.WriteFile(statFile, []byte(statMD), 0644)
	if err != nil {
		log.Printf("ERROR - Unable to write package statistics to file %s: %v", statFile, err)
	}
}

func writeDirTree(root, name string, packs []*pkgs.Package) error {
	treeFile := filepath.Join(root, tree.File)
	log.Printf("INFO - Writing directory tree to file: %s", treeFile)
	tr, err := tree.Generate(root, name, packs)
	if err != nil {
		log.Print("ERROR - Unable to generate directory tree")
		return err
	}
	err = ioutil.WriteFile(treeFile, []byte(tr), 0644)
	if err != nil {
		log.Printf("ERROR - Unable to write directory tree to file %s: %v", treeFile, err)
		return err
	}
	return nil
}

func writeDepTables(docPkgs, root, rootPkg string, noLinks bool, depMap analdata.DependencyMap) {
	log.Print("INFO - Writing dependency tables:")
	dtPkgs := findDepTablesAsSlice(docPkgs, root, rootPkg, "documentation")

	linkDocPkgs := map[string]struct{}{}
	if !noLinks {
		linkDocPkgs = dirs.FindDepTables(doc.FileName, doc.Title, dtPkgs, root, rootPkg)
	}
	doc.WriteDocs(dtPkgs, depMap, linkDocPkgs, rootPkg, root)
}

func findDepTablesAsSlice(pkgNames, root, rootPkg, pkgType string) []string {
	var packs []string
	if pkgNames == "*" { // find all existing files
		pkgMap := dirs.FindDepTables(doc.FileName, doc.Title, nil, root, rootPkg)
		packs = make([]string, 0, len(pkgMap))
		for p := range pkgMap {
			packs = append(packs, p)
		}
	} else { // write explicitly given docs
		packs = splitPackageNames(pkgNames, pkgType)
	}
	return packs
}

func splitPackageNames(docPkgs, pkgType string) []string {
	splitPkgs := strings.Split(docPkgs, ",")
	retPkgs := make([]string, 0, len(splitPkgs))
	for i, splitPkg := range splitPkgs {
		pkg := strings.TrimSpace(splitPkg)
		if pkg == "" {
			log.Printf("INFO - Not writing %s for %d-th package because the name is empty.", pkgType, i+1)
			continue
		}
		retPkgs = append(retPkgs, pkg)
	}
	return retPkgs
}
