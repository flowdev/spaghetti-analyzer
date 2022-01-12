package deps

import (
	"strings"

	"github.com/flowdev/spaghetti-analyzer/analdata"
	"github.com/flowdev/spaghetti-analyzer/x/pkgs"
	"github.com/flowdev/spaghetti-cutter/config"
	"github.com/flowdev/spaghetti-cutter/data"
)

// Fill fills the dependency map of the given package.
func Fill(pkg *pkgs.Package, rootPkg string, cfg config.Config, depMap *analdata.DependencyMap) {
	if pkgs.IsTestPackage(pkg) {
		return
	}

	relPkg, strictRelPkg := pkgs.RelativePackageName(pkg, rootPkg)
	unqPkg := pkgs.UniquePackageName(relPkg, strictRelPkg)
	pkgImps := importsOf(pkg, relPkg, strictRelPkg, rootPkg, cfg)

	if _, fullmatch := isPackageInList(cfg.God, nil, relPkg, strictRelPkg); fullmatch {
		pkgImps.PkgType = data.TypeGod
	}
	if _, fullmatch := isPackageInList(cfg.DB, nil, relPkg, strictRelPkg); fullmatch {
		pkgImps.PkgType = data.TypeDB
	}
	if _, fullmatch := isPackageInList(cfg.Tool, nil, relPkg, strictRelPkg); fullmatch {
		pkgImps.PkgType = data.TypeTool
	}

	if len(pkgImps.Imports) > 0 {
		(*depMap)[unqPkg] = pkgImps
	}
}

func importsOf(
	pkg *pkgs.Package,
	relPkg, strictRelPkg, rootPkg string,
	cfg config.Config,
) analdata.PkgImports {
	imps := analdata.PkgImports{}

	for _, p := range pkg.Imports {
		if !strings.HasPrefix(p.PkgPath, rootPkg) {
			continue
		}
		relImp, strictRelImp := pkgs.RelativePackageName(p, rootPkg)

		imps.Imports = saveDep(imps.Imports, relImp, strictRelImp, cfg)
	}
	return imps
}

func isPackageInList(pl data.PatternList, dollars []string, pkg, strictPkg string) (atAll, full bool) {
	if strictPkg != "" {
		if atAll, full := pl.MatchString(strictPkg, dollars); atAll {
			return true, full
		}
	}
	return pl.MatchString(pkg, dollars)
}

func saveDep(im map[string]data.PkgType, relImp, strictRelImp string, cfg config.Config) map[string]data.PkgType {
	if len(im) == 0 {
		im = make(map[string]data.PkgType, 32)
	}
	unqImp := pkgs.UniquePackageName(relImp, strictRelImp)

	if _, full := isPackageInList(cfg.Tool, nil, relImp, strictRelImp); full {
		im[unqImp] = data.TypeTool
	} else if _, full := isPackageInList(cfg.DB, nil, relImp, strictRelImp); full {
		im[unqImp] = data.TypeDB
	} else if _, full := isPackageInList(cfg.God, nil, relImp, strictRelImp); full {
		im[unqImp] = data.TypeGod
	} else {
		im[unqImp] = data.TypeStandard
	}
	return im
}
