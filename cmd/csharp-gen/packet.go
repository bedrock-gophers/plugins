package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/constant"
	"go/parser"
	"go/token"
	"go/types"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

type packetFieldKind uint8

const (
	packetFieldValue packetFieldKind = iota
	packetFieldBool
	packetFieldSigned
	packetFieldUnsigned
	packetFieldFloat
	packetFieldString
	packetFieldBytes
	packetFieldVec2
	packetFieldVec3
	packetFieldUUID
)

type packetFieldSpec struct {
	Name, Type string
	Kind       packetFieldKind
	Index      int
}

type packetTypeSpec struct {
	Name   string
	ID     uint32
	Fields []packetFieldSpec
}

func inspectPackets(directory string) ([]packetTypeSpec, error) {
	files := token.NewFileSet()
	packages, err := parser.ParseDir(files, directory, func(info os.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	pkg, ok := packages["packet"]
	if !ok {
		return nil, fmt.Errorf("gophertunnel packet package not found")
	}
	structures := map[string]*ast.StructType{}
	for _, file := range pkg.Files {
		for _, declaration := range file.Decls {
			generated, ok := declaration.(*ast.GenDecl)
			if !ok || generated.Tok != token.TYPE {
				continue
			}
			for _, raw := range generated.Specs {
				typeSpec, ok := raw.(*ast.TypeSpec)
				if !ok || !typeSpec.Name.IsExported() {
					continue
				}
				structure, ok := typeSpec.Type.(*ast.StructType)
				if ok {
					structures[typeSpec.Name.Name] = structure
				}
			}
		}
	}
	ids, err := gophertunnelPacketIDs(files, pkg)
	if err != nil {
		return nil, err
	}
	result := make([]packetTypeSpec, 0, len(ids))
	for name, id := range ids {
		structure := structures[name]
		if structure == nil {
			return nil, fmt.Errorf("gophertunnel packet.%s struct not found", name)
		}
		definition := packetTypeSpec{Name: name, ID: id}
		fieldIndex := 0
		for _, field := range structure.Fields.List {
			names := field.Names
			if len(names) == 0 {
				embedded := embeddedPacketFieldName(field.Type)
				if embedded == "" {
					return nil, fmt.Errorf("packet.%s has unsupported embedded field %s", name, formatGoExpression(field.Type))
				}
				names = []*ast.Ident{ast.NewIdent(embedded)}
			}
			for _, name := range names {
				if name.IsExported() {
					typeName, kind := packetCSharpField(field.Type)
					definition.Fields = append(definition.Fields, packetFieldSpec{
						Name: name.Name, Type: typeName, Kind: kind, Index: fieldIndex,
					})
				}
				fieldIndex++
			}
		}
		result = append(result, definition)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].ID != result[j].ID {
			return result[i].ID < result[j].ID
		}
		return result[i].Name < result[j].Name
	})
	return result, nil
}

func gophertunnelPacketIDs(files *token.FileSet, pkg *ast.Package) (map[string]uint32, error) {
	var idFile *ast.File
	for path, file := range pkg.Files {
		if filepath.Base(path) == "id.go" {
			idFile = file
			break
		}
	}
	if idFile == nil {
		return nil, fmt.Errorf("gophertunnel packet/id.go not found")
	}
	info := &types.Info{Defs: map[*ast.Ident]types.Object{}}
	if _, err := (&types.Config{}).Check("packet", files, []*ast.File{idFile}, info); err != nil {
		return nil, fmt.Errorf("type-check gophertunnel packet IDs: %w", err)
	}
	constants := map[string]uint32{}
	for identifier, object := range info.Defs {
		value, ok := object.(*types.Const)
		if !ok || !strings.HasPrefix(identifier.Name, "ID") {
			continue
		}
		unsigned, exact := constant.Uint64Val(value.Val())
		if !exact || unsigned > math.MaxUint32 {
			return nil, fmt.Errorf("packet constant %s is not uint32", identifier.Name)
		}
		constants[identifier.Name] = uint32(unsigned)
	}
	ids := map[string]uint32{}
	for _, file := range pkg.Files {
		for _, declaration := range file.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if !ok || function.Name.Name != "ID" || function.Recv == nil || len(function.Recv.List) != 1 || function.Body == nil || len(function.Body.List) != 1 {
				continue
			}
			receiver := function.Recv.List[0].Type
			if pointer, ok := receiver.(*ast.StarExpr); ok {
				receiver = pointer.X
			}
			typeName, ok := receiver.(*ast.Ident)
			statement, returns := function.Body.List[0].(*ast.ReturnStmt)
			if !ok || !returns || len(statement.Results) != 1 {
				continue
			}
			constantName, ok := statement.Results[0].(*ast.Ident)
			if !ok {
				continue
			}
			id, constantFound := constants[constantName.Name]
			if !constantFound {
				continue
			}
			if previous, exists := ids[typeName.Name]; exists && previous != id {
				return nil, fmt.Errorf("packet.%s has IDs %d and %d", typeName.Name, previous, id)
			}
			ids[typeName.Name] = id
		}
	}
	return ids, nil
}

func embeddedPacketFieldName(expression ast.Expr) string {
	if pointer, ok := expression.(*ast.StarExpr); ok {
		expression = pointer.X
	}
	switch value := expression.(type) {
	case *ast.Ident:
		return value.Name
	case *ast.SelectorExpr:
		return value.Sel.Name
	default:
		return ""
	}
}

func packetCSharpField(expression ast.Expr) (string, packetFieldKind) {
	switch value := expression.(type) {
	case *ast.Ident:
		switch value.Name {
		case "bool":
			return "bool", packetFieldBool
		case "int8":
			return "sbyte", packetFieldSigned
		case "int16":
			return "short", packetFieldSigned
		case "int", "int32":
			return "int", packetFieldSigned
		case "int64":
			return "long", packetFieldSigned
		case "byte", "uint8":
			return "byte", packetFieldUnsigned
		case "uint16":
			return "ushort", packetFieldUnsigned
		case "uint", "uint32":
			return "uint", packetFieldUnsigned
		case "uint64":
			return "ulong", packetFieldUnsigned
		case "float32":
			return "float", packetFieldFloat
		case "float64":
			return "double", packetFieldFloat
		case "string":
			return "string", packetFieldString
		default:
			return "Value", packetFieldValue
		}
	case *ast.SelectorExpr:
		pkg, ok := value.X.(*ast.Ident)
		if !ok {
			return "Value", packetFieldValue
		}
		switch pkg.Name + "." + value.Sel.Name {
		case "mgl32.Vec2":
			return "Vector2", packetFieldVec2
		case "mgl32.Vec3":
			return "Vector3", packetFieldVec3
		case "uuid.UUID":
			return "Guid", packetFieldUUID
		default:
			return "Value", packetFieldValue
		}
	case *ast.ArrayType:
		if value.Len == nil {
			if element, ok := value.Elt.(*ast.Ident); ok && (element.Name == "byte" || element.Name == "uint8") {
				return "byte[]", packetFieldBytes
			}
		}
		return "Value", packetFieldValue
	default:
		return "Value", packetFieldValue
	}
}

func inspectInterceptHandler(path string) error {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return err
	}
	want := map[string]string{
		"HandleClientPacket": "ctx *Context, pk packet.Packet",
		"HandleServerPacket": "ctx *Context, pk packet.Packet",
	}
	for _, declaration := range file.Decls {
		generated, ok := declaration.(*ast.GenDecl)
		if !ok || generated.Tok != token.TYPE {
			continue
		}
		for _, raw := range generated.Specs {
			typeSpec, ok := raw.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != "Handler" {
				continue
			}
			contract, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				return fmt.Errorf("intercept.Handler is not an interface")
			}
			found := map[string]string{}
			for _, field := range contract.Methods.List {
				if len(field.Names) == 1 {
					function, ok := field.Type.(*ast.FuncType)
					if ok {
						found[field.Names[0].Name] = formatFieldList(function.Params)
					}
				}
			}
			if !reflect.DeepEqual(found, want) {
				return fmt.Errorf("intercept.Handler changed: %v", found)
			}
			return nil
		}
	}
	return fmt.Errorf("intercept.Handler not found")
}

func formatFieldList(fields *ast.FieldList) string {
	if fields == nil {
		return ""
	}
	parts := make([]string, 0, len(fields.List))
	for _, field := range fields.List {
		for _, name := range field.Names {
			parts = append(parts, name.Name+" "+formatGoExpression(field.Type))
		}
	}
	return strings.Join(parts, ", ")
}

func generatePacketTypes(types []packetTypeSpec) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from gophertunnel minecraft/protocol/packet Go AST. DO NOT EDIT.\n")
	output.WriteString("#nullable enable\nusing System;\n\nnamespace Dragonfly.Packet;\n\n")
	for _, definition := range types {
		fmt.Fprintf(&output, "public sealed class %s : Packet\n{\n", definition.Name)
		output.WriteString("    private readonly ulong _handle;\n")
		fmt.Fprintf(&output, "    internal %s(ulong handle) => _handle = handle;\n", definition.Name)
		fmt.Fprintf(&output, "    public uint ID() => %du;\n", definition.ID)
		for _, field := range definition.Fields {
			generatePacketField(&output, definition.Name, field)
		}
		output.WriteString("}\n\n")
	}
	output.WriteString("internal static class PacketCodec\n{\n")
	output.WriteString("    internal static Packet Decode(ulong handle, uint id) => id switch\n    {\n")
	for _, definition := range types {
		fmt.Fprintf(&output, "        %du => new %s(handle),\n", definition.ID, definition.Name)
	}
	output.WriteString("        _ => new Unknown(handle, id),\n    };\n}\n")
	return output.Bytes()
}

func generatePacketField(output *bytes.Buffer, owner string, field packetFieldSpec) {
	name := csharpIdentifier(field.Name)
	if strings.TrimPrefix(name, "@") == owner {
		name += "Value"
	}
	switch field.Kind {
	case packetFieldBool:
		fmt.Fprintf(output, "    public bool %s { get => PacketBridge.Bool(_handle, %d); set => PacketBridge.SetBool(_handle, %d, value); }\n", name, field.Index, field.Index)
	case packetFieldSigned:
		fmt.Fprintf(output, "    public %s %s { get => checked((%s)PacketBridge.Signed(_handle, %d)); set => PacketBridge.SetSigned(_handle, %d, value); }\n", field.Type, name, field.Type, field.Index, field.Index)
	case packetFieldUnsigned:
		fmt.Fprintf(output, "    public %s %s { get => checked((%s)PacketBridge.Unsigned(_handle, %d)); set => PacketBridge.SetUnsigned(_handle, %d, value); }\n", field.Type, name, field.Type, field.Index, field.Index)
	case packetFieldFloat:
		fmt.Fprintf(output, "    public %s %s { get => (%s)PacketBridge.Number(_handle, %d); set => PacketBridge.SetNumber(_handle, %d, value); }\n", field.Type, name, field.Type, field.Index, field.Index)
	case packetFieldString:
		fmt.Fprintf(output, "    public string %s { get => PacketBridge.String(_handle, %d); set => PacketBridge.SetString(_handle, %d, value); }\n", name, field.Index, field.Index)
	case packetFieldBytes:
		fmt.Fprintf(output, "    public byte[] %s { get => PacketBridge.Bytes(_handle, %d); set => PacketBridge.SetBytes(_handle, %d, value); }\n", name, field.Index, field.Index)
	case packetFieldVec2:
		fmt.Fprintf(output, "    public Vector2 %s { get => PacketBridge.Vector2(_handle, %d); set => PacketBridge.SetVector2(_handle, %d, value); }\n", name, field.Index, field.Index)
	case packetFieldVec3:
		fmt.Fprintf(output, "    public Vector3 %s { get => PacketBridge.Vector3(_handle, %d); set => PacketBridge.SetVector3(_handle, %d, value); }\n", name, field.Index, field.Index)
	case packetFieldUUID:
		fmt.Fprintf(output, "    public Guid %s { get => PacketBridge.Guid(_handle, %d); set => PacketBridge.SetGuid(_handle, %d, value); }\n", name, field.Index, field.Index)
	default:
		fmt.Fprintf(output, "    public Value %s => new(_handle, %d);\n", name, field.Index)
	}
}

func generatePacketHandler() []byte {
	return []byte(`// Code generated from bedrock-gophers/intercept Handler Go AST. DO NOT EDIT.
#nullable enable
namespace Dragonfly;

public abstract partial class Plugin
{
    [HandlerSubscription(18014398509481984UL)]
    public virtual void HandleClientPacket(Packet.Context ctx, Packet.Packet packet) { }

    [HandlerSubscription(36028797018963968UL)]
    public virtual void HandleServerPacket(Packet.Context ctx, Packet.Packet packet) { }
}
`)
}
