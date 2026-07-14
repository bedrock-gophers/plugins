package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

type bboxSpec struct {
	Axes []string
}

func inspectBBox(path, axisPath string) (bboxSpec, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return bboxSpec{}, err
	}
	if err := inspectBBoxType(file); err != nil {
		return bboxSpec{}, err
	}

	wantFunctions := map[string]goSignature{
		"Box":              {Parameters: "float64, float64, float64, float64, float64, float64", Results: "BBox"},
		"AnyIntersections": {Parameters: "[]BBox, BBox", Results: "bool"},
	}
	wantFunctionParameters := map[string]string{
		"Box":              "x0,y0,z0,x1,y1,z1",
		"AnyIntersections": "boxes,search",
	}
	wantMethods := map[string]goSignature{
		"Grow":             {Parameters: "float64", Results: "BBox"},
		"GrowVec3":         {Parameters: "mgl64.Vec3", Results: "BBox"},
		"Min":              {Results: "mgl64.Vec3"},
		"Max":              {Results: "mgl64.Vec3"},
		"Width":            {Results: "float64"},
		"Length":           {Results: "float64"},
		"Height":           {Results: "float64"},
		"Extend":           {Parameters: "mgl64.Vec3", Results: "BBox"},
		"ExtendTowards":    {Parameters: "Face, float64", Results: "BBox"},
		"Stretch":          {Parameters: "Axis, float64", Results: "BBox"},
		"Translate":        {Parameters: "mgl64.Vec3", Results: "BBox"},
		"TranslateTowards": {Parameters: "Face, float64", Results: "BBox"},
		"IntersectsWith":   {Parameters: "BBox", Results: "bool"},
		"Vec3Within":       {Parameters: "mgl64.Vec3", Results: "bool"},
		"Vec3WithinYZ":     {Parameters: "mgl64.Vec3", Results: "bool"},
		"Vec3WithinXZ":     {Parameters: "mgl64.Vec3", Results: "bool"},
		"Vec3WithinXY":     {Parameters: "mgl64.Vec3", Results: "bool"},
		"XOffset":          {Parameters: "BBox, float64", Results: "float64"},
		"YOffset":          {Parameters: "BBox, float64", Results: "float64"},
		"ZOffset":          {Parameters: "BBox, float64", Results: "float64"},
		"Corners":          {Results: "[]mgl64.Vec3"},
		"Mul":              {Parameters: "float64", Results: "BBox"},
		"Volume":           {Results: "float64"},
	}
	wantMethodParameters := map[string]string{
		"Grow": "x", "GrowVec3": "vec", "Min": "", "Max": "", "Width": "", "Length": "", "Height": "",
		"Extend": "vec", "ExtendTowards": "f,x", "Stretch": "a,x", "Translate": "vec", "TranslateTowards": "f,x",
		"IntersectsWith": "other", "Vec3Within": "vec", "Vec3WithinYZ": "vec", "Vec3WithinXZ": "vec", "Vec3WithinXY": "vec",
		"XOffset": "nearby,deltaX", "YOffset": "nearby,deltaY", "ZOffset": "nearby,deltaZ", "Corners": "", "Mul": "val", "Volume": "",
	}
	functions := map[string]goSignature{}
	methods := map[string]goSignature{}
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || !function.Name.IsExported() {
			continue
		}
		if function.Recv == nil {
			if got, want := bboxParameterNames(function.Type.Params), wantFunctionParameters[function.Name.Name]; got != want {
				return bboxSpec{}, fmt.Errorf("Dragonfly cube.%s parameter names changed: %s", function.Name.Name, got)
			}
			functions[function.Name.Name] = goFunctionSignature(function)
			continue
		}
		receiver, ok := receiverName(function)
		if !ok || receiver != "BBox" {
			continue
		}
		if !valueReceiver(function, "BBox") {
			return bboxSpec{}, fmt.Errorf("Dragonfly cube.BBox.%s receiver changed", function.Name.Name)
		}
		if got, want := bboxParameterNames(function.Type.Params), wantMethodParameters[function.Name.Name]; got != want {
			return bboxSpec{}, fmt.Errorf("Dragonfly cube.BBox.%s parameter names changed: %s", function.Name.Name, got)
		}
		methods[function.Name.Name] = goFunctionSignature(function)
	}
	if !reflect.DeepEqual(functions, wantFunctions) {
		return bboxSpec{}, fmt.Errorf("Dragonfly cube BBox functions changed: %v", functions)
	}
	if !reflect.DeepEqual(methods, wantMethods) {
		return bboxSpec{}, fmt.Errorf("Dragonfly cube.BBox methods changed: %v", methods)
	}
	axes, err := inspectBBoxAxes(axisPath)
	if err != nil {
		return bboxSpec{}, err
	}
	return bboxSpec{Axes: axes}, nil
}

func bboxParameterNames(fields *ast.FieldList) string {
	if fields == nil {
		return ""
	}
	var names []string
	for _, field := range fields.List {
		for _, name := range field.Names {
			names = append(names, name.Name)
		}
	}
	return strings.Join(names, ",")
}

func inspectBBoxType(file *ast.File) error {
	for _, declaration := range file.Decls {
		general, ok := declaration.(*ast.GenDecl)
		if !ok || general.Tok != token.TYPE {
			continue
		}
		for _, raw := range general.Specs {
			spec, ok := raw.(*ast.TypeSpec)
			if !ok || spec.Name.Name != "BBox" {
				continue
			}
			structure, ok := spec.Type.(*ast.StructType)
			if !ok || len(structure.Fields.List) != 1 {
				return fmt.Errorf("Dragonfly cube.BBox shape changed")
			}
			field := structure.Fields.List[0]
			if len(field.Names) != 2 || field.Names[0].Name != "min" || field.Names[1].Name != "max" ||
				formatGoExpression(field.Type) != "mgl64.Vec3" {
				return fmt.Errorf("Dragonfly cube.BBox fields changed")
			}
			return nil
		}
	}
	return fmt.Errorf("Dragonfly cube.BBox not found")
}

func inspectBBoxAxes(path string) ([]string, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return nil, err
	}
	if !hasNamedType(file, "Axis", "int") {
		return nil, fmt.Errorf("Dragonfly cube.Axis is not backed by int")
	}
	var axes []string
	for _, declaration := range file.Decls {
		general, ok := declaration.(*ast.GenDecl)
		if !ok || general.Tok != token.CONST {
			continue
		}
		for _, raw := range general.Specs {
			spec, ok := raw.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, name := range spec.Names {
				axes = append(axes, name.Name)
			}
		}
	}
	want := []string{"Y", "Z", "X"}
	if !reflect.DeepEqual(axes, want) {
		return nil, fmt.Errorf("Dragonfly cube.Axis values changed: %v", axes)
	}
	return axes, nil
}

func generateBBox(spec bboxSpec) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/block/cube/bbox.go Go AST. DO NOT EDIT.\n")
	output.WriteString("namespace Dragonfly;\n\n")
	output.WriteString("public static partial class Cube\n{\n")
	output.WriteString("    public enum Axis\n    {\n")
	for index, axis := range spec.Axes {
		fmt.Fprintf(&output, "        %s = %d,\n", axis, index)
	}
	output.WriteString("    }\n\n")
	output.WriteString(`    public static BBox Box(double x0, double y0, double z0, double x1, double y1, double z1)
    {
        if (x0 > x1) (x0, x1) = (x1, x0);
        if (y0 > y1) (y0, y1) = (y1, y0);
        if (z0 > z1) (z0, z1) = (z1, z0);
        return new(new(x0, y0, z0), new(x1, y1, z1));
    }

    public static bool AnyIntersections(BBox[] boxes, BBox search)
    {
        foreach (var box in boxes)
        {
            if (box.IntersectsWith(search, 0)) return true;
        }
        return false;
    }

    public readonly record struct BBox
    {
        private readonly Vector3 _min;
        private readonly Vector3 _max;

        internal BBox(Vector3 min, Vector3 max) => (_min, _max) = (min, max);

        public BBox Grow(double x) => new(
            new(_min.X - x, _min.Y - x, _min.Z - x),
            new(_max.X + x, _max.Y + x, _max.Z + x));

        public BBox GrowVec3(Vector3 vec) => new(
            new(_min.X - vec.X, _min.Y - vec.Y, _min.Z - vec.Z),
            new(_max.X + vec.X, _max.Y + vec.Y, _max.Z + vec.Z));

        public Vector3 Min() => _min;
        public Vector3 Max() => _max;
        public double Width() => _max.X - _min.X;
        public double Length() => _max.Z - _min.Z;
        public double Height() => _max.Y - _min.Y;

        public BBox Extend(Vector3 vec)
        {
            var min = _min;
            var max = _max;
            if (vec.X < 0) min = min with { X = min.X + vec.X };
            else if (vec.X > 0) max = max with { X = max.X + vec.X };
            if (vec.Y < 0) min = min with { Y = min.Y + vec.Y };
            else if (vec.Y > 0) max = max with { Y = max.Y + vec.Y };
            if (vec.Z < 0) min = min with { Z = min.Z + vec.Z };
            else if (vec.Z > 0) max = max with { Z = max.Z + vec.Z };
            return new(min, max);
        }

        public BBox ExtendTowards(Face f, double x) => f switch
        {
            Face.Down => new(new(_min.X, _min.Y - x, _min.Z), _max),
            Face.Up => new(_min, new(_max.X, _max.Y + x, _max.Z)),
            Face.North => new(new(_min.X, _min.Y, _min.Z - x), _max),
            Face.South => new(_min, new(_max.X, _max.Y, _max.Z + x)),
            Face.West => new(new(_min.X - x, _min.Y, _min.Z), _max),
            Face.East => new(_min, new(_max.X + x, _max.Y, _max.Z)),
            _ => this,
        };

        public BBox Stretch(Axis a, double x) => a switch
        {
            Axis.Y => new(new(_min.X, _min.Y - x, _min.Z), new(_max.X, _max.Y + x, _max.Z)),
            Axis.Z => new(new(_min.X, _min.Y, _min.Z - x), new(_max.X, _max.Y, _max.Z + x)),
            Axis.X => new(new(_min.X - x, _min.Y, _min.Z), new(_max.X + x, _max.Y, _max.Z)),
            _ => this,
        };

        public BBox Translate(Vector3 vec) => new(
            new(_min.X + vec.X, _min.Y + vec.Y, _min.Z + vec.Z),
            new(_max.X + vec.X, _max.Y + vec.Y, _max.Z + vec.Z));

        public BBox TranslateTowards(Face f, double x) => f switch
        {
            Face.Down => Translate(new(0, -x, 0)),
            Face.Up => Translate(new(0, x, 0)),
            Face.North => Translate(new(0, 0, -x)),
            Face.South => Translate(new(0, 0, x)),
            Face.West => Translate(new(-x, 0, 0)),
            Face.East => Translate(new(x, 0, 0)),
            _ => this,
        };

        public bool IntersectsWith(BBox other) => IntersectsWith(other, 1e-5);

        internal bool IntersectsWith(BBox other, double epsilon) =>
            other._max.X - _min.X > epsilon && _max.X - other._min.X > epsilon &&
            other._max.Y - _min.Y > epsilon && _max.Y - other._min.Y > epsilon &&
            other._max.Z - _min.Z > epsilon && _max.Z - other._min.Z > epsilon;

        public bool Vec3Within(Vector3 vec) =>
            vec.X > _min.X && vec.X < _max.X &&
            vec.Z > _min.Z && vec.Z < _max.Z &&
            vec.Y > _min.Y && vec.Y < _max.Y;

        public bool Vec3WithinYZ(Vector3 vec) =>
            vec.Z >= _min.Z && vec.Z <= _max.Z &&
            vec.Y >= _min.Y && vec.Y <= _max.Y;

        public bool Vec3WithinXZ(Vector3 vec) =>
            vec.X >= _min.X && vec.X <= _max.X &&
            vec.Z >= _min.Z && vec.Z <= _max.Z;

        public bool Vec3WithinXY(Vector3 vec) =>
            vec.X >= _min.X && vec.X <= _max.X &&
            vec.Y >= _min.Y && vec.Y <= _max.Y;

        public double XOffset(BBox nearby, double deltaX)
        {
            if (_max.Y <= nearby._min.Y || _min.Y >= nearby._max.Y ||
                _max.Z <= nearby._min.Z || _min.Z >= nearby._max.Z) return deltaX;
            if (deltaX > 0 && _max.X <= nearby._min.X)
                deltaX = Math.Min(deltaX, nearby._min.X - _max.X);
            else if (deltaX < 0 && _min.X >= nearby._max.X)
                deltaX = Math.Max(deltaX, nearby._max.X - _min.X);
            return deltaX;
        }

        public double YOffset(BBox nearby, double deltaY)
        {
            if (_max.X <= nearby._min.X || _min.X >= nearby._max.X ||
                _max.Z <= nearby._min.Z || _min.Z >= nearby._max.Z) return deltaY;
            if (deltaY > 0 && _max.Y <= nearby._min.Y)
                deltaY = Math.Min(deltaY, nearby._min.Y - _max.Y);
            if (deltaY < 0 && _min.Y >= nearby._max.Y)
                deltaY = Math.Max(deltaY, nearby._max.Y - _min.Y);
            return deltaY;
        }

        public double ZOffset(BBox nearby, double deltaZ)
        {
            if (_max.X <= nearby._min.X || _min.X >= nearby._max.X ||
                _max.Y <= nearby._min.Y || _min.Y >= nearby._max.Y) return deltaZ;
            if (deltaZ > 0 && _max.Z <= nearby._min.Z)
                deltaZ = Math.Min(deltaZ, nearby._min.Z - _max.Z);
            if (deltaZ < 0 && _min.Z >= nearby._max.Z)
                deltaZ = Math.Max(deltaZ, nearby._max.Z - _min.Z);
            return deltaZ;
        }

        public Vector3[] Corners() =>
        [
            _min,
            _max,
            new(_min.X, _min.Y, _max.Z),
            new(_min.X, _max.Y, _min.Z),
            new(_min.X, _max.Y, _max.Z),
            new(_max.X, _max.Y, _min.Z),
            new(_max.X, _min.Y, _max.Z),
            new(_max.X, _min.Y, _min.Z),
        ];

        public BBox Mul(double val) => new(
            new(_min.X * val, _min.Y * val, _min.Z * val),
            new(_max.X * val, _max.Y * val, _max.Z * val));

        public double Volume() => Height() * Length() * Width();
    }
}
`)
	return output.Bytes()
}
