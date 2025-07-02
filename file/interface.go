package file

import (
	"context"
	"fmt"
	"github.com/ribeirosaimon/testgen/domain"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type file struct {
	filePath    string
	PackageName string
	Interface   *domain.Interface
	Methods     []domain.Method
}

func New(ctx context.Context, path string) *file {
	myFile := &file{
		filePath:    path,
		Methods:     make([]domain.Method, 0),
		PackageName: "meupackage",
	}
	myFile.createInterfaceMethods(ctx)
	return myFile
}

func (f *file) createInterfaceMethods(ctx context.Context) error {
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, f.filePath, nil, parser.AllErrors)
	f.PackageName = strings.ToLower(node.Name.Name)
	if err != nil {
		fmt.Println("Erro ao analisar o arquivo:", err)
		return err
	}

	fmt.Printf("Interfaces encontradas em %s:\n", f.filePath)
	resp := make(map[string]int)

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		funcDecl, funcOk := decl.(*ast.FuncDecl)
		if funcOk {
			i := f.normalizeFunc(funcDecl, fs)
			resp[funcDecl.Name.Name] = i
		}
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		f.findInterfaceAndMethods(ctx, genDecl, fs)
	}
	for i := range f.Methods {
		if v, ok := resp[f.Methods[i].MethodName]; ok {
			for count := range v {
				f.Methods[i].CountIfs = append(f.Methods[i].CountIfs, fmt.Sprintf("%d", count+1))
			}
		}
	}
	return nil
}

func (f *file) normalizeFunc(funcDecl *ast.FuncDecl, fs *token.FileSet) int {
	if funcDecl.Recv == nil {
		return 0
	}
	if recvField := funcDecl.Recv.List[0]; recvField != nil {
		return countIfsInFunc(funcDecl)
	}
	return 0
}

func (f *file) findInterfaceAndMethods(ctx context.Context, genDecl *ast.GenDecl, fs *token.FileSet) {
	myInterface := domain.Interface{}
	for _, spec := range genDecl.Specs {
		typeSpec := spec.(*ast.TypeSpec)
		if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
			fmt.Printf("- %s (na linha %d)\n", typeSpec.Name.Name, fs.Position(typeSpec.Pos()).Line)
			myInterface.InterfaceName = typeSpec.Name.Name

			// Listar métodos da interface
			for _, method := range interfaceType.Methods.List {
				if len(method.Names) == 0 {
					continue // pode ser um mét\odo embutido, ignoramos por simplicidade
				}
				myMethods := domain.Method{
					MethodName: method.Names[0].Name,
				}
				funcType, ok := method.Type.(*ast.FuncType)
				if !ok {
					continue
				}

				var params []string
				if funcType.Params != nil {
					for _, param := range funcType.Params.List {
						typeStr := fmt.Sprintf("%s", exprToString(param.Type))
						for range param.Names {
							params = append(params, typeStr)
						}
						// Caso o parâmetro não tenha nome, ainda adicionamos o tipo
						if len(param.Names) == 0 {
							params = append(params, typeStr)
						}
					}
				}

				// 4. Resultados
				var results []string
				if funcType.Results != nil {
					for _, result := range funcType.Results.List {
						typeStr := fmt.Sprintf("%s", exprToString(result.Type))
						for range result.Names {
							results = append(results, typeStr)
						}
						if len(result.Names) == 0 {
							results = append(results, typeStr)
						}
					}
				}

				myMethods.Params = params
				myMethods.Results = results
				f.Methods = append(f.Methods, myMethods)
			}
		} else {
			return
		}
	}
	f.Interface = &myInterface
}

func exprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return exprToString(e.X) + "." + e.Sel.Name
	case *ast.StarExpr:
		return "*" + exprToString(e.X)
	case *ast.ArrayType:
		return "[]" + exprToString(e.Elt)
	case *ast.MapType:
		return "map[" + exprToString(e.Key) + "]" + exprToString(e.Value)
	case *ast.FuncType:
		return "func" // simplificado
	default:
		return fmt.Sprintf("%T", expr)
	}
}

func countIfsInFunc(funcDecl *ast.FuncDecl) int {
	count := 0
	if funcDecl.Body == nil {
		return 0
	}
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		if _, ok := n.(*ast.IfStmt); ok {
			count++
		}
		return true
	})
	return count
}
