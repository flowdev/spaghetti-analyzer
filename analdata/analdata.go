package analdata

import (
	"sort"
	"strings"

	"github.com/flowdev/spaghetti-cutter/data"
)

// PkgImports contains the package type and the imported internal packages with their types.
type PkgImports struct {
	PkgType data.PkgType
	Imports map[string]data.PkgType
}

// DependencyMap is mapping importing package to imported packages.
// importingPackageName -> (importedPackageNames -> PkgType)
// An imported package name could be added multiple times to the same importing
// package name due to test packages.
type DependencyMap map[string]PkgImports

// SortedPkgNames returns the sorted keys (package names) of the dependency map.
func (dm DependencyMap) SortedPkgNames() []string {
	names := make([]string, 0, len(dm))
	for pkg := range dm {
		names = append(names, pkg)
	}
	sort.Strings(names)
	return names
}

// FilterDepMap filters allMap to contain only packages matching idx and its transitive
// dependencies.  Entries matching other indices in links are filtered, too.
func FilterDepMap(allMap DependencyMap, idx int, links data.PatternList) DependencyMap {
	if idx < 0 || len(links) == 0 {
		return allMap
	}

	fltrMap := make(DependencyMap, len(allMap))
	for pkg := range allMap {
		if i := data.DocMatchStringIndex(pkg, links); i >= 0 && i == idx {
			copyDepsRecursive(allMap, pkg, fltrMap, links, idx)
		}
	}
	return fltrMap
}
func copyDepsRecursive(
	allMap DependencyMap,
	startPkg string,
	fltrMap DependencyMap,
	links data.PatternList,
	idx int,
) {
	if i := data.DocMatchStringIndex(startPkg, links); i >= 0 && i != idx {
		return
	}
	imps, ok := allMap[startPkg]
	if !ok {
		return
	}
	fltrMap[startPkg] = imps
	for pkg := range imps.Imports {
		copyDepsRecursive(allMap, pkg, fltrMap, links, idx)
	}
}

// PkgForPattern returns the (parent) package of the given package pattern.
// If pkg doesn't contain any wildcard '*' the whole string is returned.
// Otherwise everything up to the last '/' before the wildcard or
// the empty string if there is no '/' before it.
func PkgForPattern(pkg string) string {
	i := strings.IndexRune(pkg, '*')
	if i < 0 {
		return pkg
	}
	i = strings.LastIndex(pkg[:i], "/")
	if i > 0 {
		return pkg[:i]
	}
	return ""
}
