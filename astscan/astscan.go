// astscan 测试

package astscan

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"

	"github.com/fatih/astrewrite"
)

// max depth for scanning dir, you can change it before running Dir
var (
	MaxDepth = 7
)

var (
	_callback Callback
	_check    Checker
)

func getPos(fset *token.FileSet, pos token.Pos) string {
	if fset == nil {
		return ""
	}
	return fset.Position(pos).String()
}

func scan(n ast.Node, fset *token.FileSet, pkg string) (ast.Node, bool) {
	if n == nil {
		return nil, true
	}

	switch v := n.(type) {
	case *ast.BasicLit:
		if v.Kind == token.STRING {
			if _check(v.Value) {
				_callback(Item{pkg, getPos(fset, v.ValuePos), TypeString, v.Value, n})
			}
		}
	case *ast.Comment:
		if _check(v.Text) {
			_callback(Item{pkg, getPos(fset, v.Slash), TypeComment, v.Text, n})
		}
	}
	return n, true
}

// File parses the file and checks the String and Comment on ast.
func File(file string, check Checker, callback Callback) error {
	if err := checkParams(check, callback); err != nil {
		return err
	}

	fset := token.NewFileSet()
	fd, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	astrewrite.Walk(fd, func(n ast.Node) (ast.Node, bool) {
		return scan(n, fset, "")
	})
	return nil
}

// Dir scans directory with deep iteration firstly, and parses every
// directory to check the String and Comment on ast. Dir's depth limits
// by MaxDepth.
func Dir(dir string, check Checker, callback Callback) error {
	if err := checkParams(check, callback); err != nil {
		return err
	}

	return scanDir(dir, 0)
}

func scanDir(dir string, depth int) error {
	if MaxDepth > 0 && depth > MaxDepth {
		return errors.New("astscan: Dir reaches max depth")
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for pkg, n := range pkgs {
		astrewrite.Walk(n, func(n ast.Node) (ast.Node, bool) {
			return scan(n, fset, pkg)
		})
	}

	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, fd := range dirs {
		if !fd.IsDir() {
			continue
		}

		if err := scanDir(filepath.Join(dir, fd.Name()), depth+1); err != nil {
			return err
		}
	}
	return nil
}

func checkParams(check Checker, callback Callback) error {
	if check == nil {
		return errors.New("check cannot be nil")
	}
	if callback == nil {
		return errors.New("callback cannot be nil")
	}

	_check, _callback = check, callback
	return nil
}
