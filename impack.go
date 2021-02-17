package impack

import (
	"bytes"
	"context"
	"go/token"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/gopackages"
	"golang.org/x/tools/go/packages"
)

// Linter performs a simple memory packing on go a package
type Linter struct {
	sizes types.Sizes
}

// NewLinter creates a new ImpackLinter by setting up the sizes and alignments
func NewLinter(compiler, arch string) *Linter {
	return &Linter{
		sizes: types.SizesFor(compiler, arch),
	}
}

// Lint finds all structs in a package and lints them
func (linter *Linter) Lint(ctx context.Context, projectRoot string) error {
	fset := token.NewFileSet()
	cfg := &packages.Config{
		Fset:    fset,
		Context: ctx,
		Dir:     projectRoot,
		Mode:    packages.LoadAllSyntax,
		Tests:   false,
	}
	pkgs, err := decorator.Load(cfg, "")
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		// avoid bound check
		_ = pkg.GoFiles[len(pkg.Syntax)-1]
		for i, dstFile := range pkg.Syntax {
			fileChanged := false
			for _, decl := range dstFile.Decls {
				if err = ctx.Err(); err != nil {
					return err
				}
				genDecl, ok := decl.(*dst.GenDecl)
				if ok {
					for _, s := range genDecl.Specs {
						t, ok := s.(*dst.TypeSpec)
						if ok {
							st, ok := t.Type.(*dst.StructType)
							if ok {
								if tn, ok := pkg.Types.Scope().Lookup(t.Name.Name).(*types.TypeName); ok {
									if typesStruct, ok := tn.Type().Underlying().(*types.Struct); ok {
										fileChanged = true
										err = linter.lintStruct(st, typesStruct)
										if err != nil {
											return err
										}
									}
								}
							}
						}
					}
				}
			}
			if fileChanged {
				if err = writeFile(pkg.GoFiles[i], dstFile); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (linter *Linter) lintStruct(st *dst.StructType, typesStruct *types.Struct) error {
	nf := typesStruct.NumFields()
	fieldSizes := make(map[string]int64, nf)
	for i := 0; i < nf; i++ {
		f := typesStruct.Field(i)
		fieldSizes[f.Name()] = linter.typeSize(f.Type())
	}
	sort.SliceStable(st.Fields.List, func(i, j int) bool {
		field1 := st.Fields.List[i]
		field2 := st.Fields.List[j]
		name1 := field1.Names[0].Name
		name2 := field2.Names[0].Name
		size1 := fieldSizes[name1]
		size2 := fieldSizes[name2]
		if size1 == size2 {
			// if sizes are equal then order alphabetically
			var name1Lower = strings.ToLower(name1)
			var name2Lower = strings.ToLower(name2)
			if name1Lower == name2Lower {
				return name1 < name2
			}
			return name1Lower < name2Lower
		}
		// lower sizes on top
		return size1 < size2

	})
	return nil
}

func (linter *Linter) typeSize(t types.Type) int64 {
	switch tp := t.(type) {
	case *types.Array:
		n := tp.Len()
		if n <= 0 {
			return 0
		}
		size := linter.typeSize(tp.Elem())
		alignment := linter.sizes.Alignof(tp.Elem())
		return align(size, alignment)*(n-1) + size
	}
	return linter.sizes.Sizeof(t)
}

func align(x int64, a int64) int64 {
	y := x + a - 1
	return y - y%a
}

func writeFile(path string, dstFile *dst.File) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	// we use a buffer in case the decorator fails. If it does, the file's content won't be destroyed
	var b []byte
	buf := bytes.NewBuffer(b)
	dir := filepath.Dir(path)
	r := decorator.NewRestorerWithImports(path, gopackages.New(dir))
	if err = r.Fprint(buf, dstFile); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, fileInfo.Mode())
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, buf)
	return err
}
