package tree_test

import (
	"fmt"
	"go/ast"
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"

	"github.com/flowdev/spaghetti-analyzer/tree"
	"github.com/flowdev/spaghetti-analyzer/x/pkgs"
)

func TestGenerate(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
	})
}

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"dirTree": callGenerate,
	}))
}

func callGenerate() int {
	packs := []*pkgs.Package{
		genPkg("bla/cmd/exe1", "Package exe1 is all about exe.cution first. Or third."),
		genPkg("bla/cmd/exe2", "Package exe2 is all about exe!cution second! And..."),
		genPkg("bla/pkg/db", "Package db contains the data:base code: CODE"),
		genPkg("bla/pkg/db/model", "Package model holds models? Or anything at all?"),
		genPkg("bla/pkg/db/store", ""),
		genPkg("bla/pkg/domain1", "Package domain1 is ..."),
		genPkg("bla/pkg/domain2", "Package domain2 is."),
	}

	name := "project-root"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	output, err := tree.Generate(".", name, packs)
	if err != nil {
		fmt.Printf("ERROR: Unable to generate directory tree: %v", err)
		return 1
	}

	fmt.Println(output)
	return 0
}
func genPkg(pkgPath, comment string) *pkgs.Package {
	pkg := &pkgs.Package{PkgPath: pkgPath}

	if comment != "" {
		pkg.Syntax = []*ast.File{
			{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{{Text: comment}},
				},
			},
		}
	}
	return pkg
}
