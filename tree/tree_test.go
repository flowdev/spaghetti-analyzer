package tree_test

import (
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
		Cmds: map[string]func(*testscript.TestScript, bool, []string){
			"dirTree": callGenerate,
		},
		// TestWork: true,
	})
}

func callGenerate(ts *testscript.TestScript, _ bool, args []string) {
	workDir := ts.Getenv("WORK")
	treeFile := workDir + "/dirtree.actual"
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
	if len(args) > 0 {
		name = args[0]
	}

	output, err := tree.Generate(workDir, name, packs)
	if err != nil {
		ts.Fatalf("ERROR: Unable to generate directory tree: %v", err)
	}

	err = os.WriteFile(treeFile, []byte(output+"\n"), 0666)
	if err != nil {
		ts.Fatalf("ERROR: Unable to write file '%s': %v", treeFile, err)
	}
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
