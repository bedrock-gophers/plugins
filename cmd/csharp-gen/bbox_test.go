package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const bboxSource = `package cube
type BBox struct { min, max mgl64.Vec3 }
func Box(x0, y0, z0, x1, y1, z1 float64) BBox { return BBox{} }
func AnyIntersections(boxes []BBox, search BBox) bool { return false }
func (box BBox) Grow(x float64) BBox { return box }
func (box BBox) GrowVec3(vec mgl64.Vec3) BBox { return box }
func (box BBox) Min() mgl64.Vec3 { return box.min }
func (box BBox) Max() mgl64.Vec3 { return box.max }
func (box BBox) Width() float64 { return 0 }
func (box BBox) Length() float64 { return 0 }
func (box BBox) Height() float64 { return 0 }
func (box BBox) Extend(vec mgl64.Vec3) BBox { return box }
func (box BBox) ExtendTowards(f Face, x float64) BBox { return box }
func (box BBox) Stretch(a Axis, x float64) BBox { return box }
func (box BBox) Translate(vec mgl64.Vec3) BBox { return box }
func (box BBox) TranslateTowards(f Face, x float64) BBox { return box }
func (box BBox) IntersectsWith(other BBox) bool { return false }
func (box BBox) intersectsWith(other BBox, epsilon float64) bool { return false }
func (box BBox) Vec3Within(vec mgl64.Vec3) bool { return false }
func (box BBox) Vec3WithinYZ(vec mgl64.Vec3) bool { return false }
func (box BBox) Vec3WithinXZ(vec mgl64.Vec3) bool { return false }
func (box BBox) Vec3WithinXY(vec mgl64.Vec3) bool { return false }
func (box BBox) XOffset(nearby BBox, deltaX float64) float64 { return deltaX }
func (box BBox) YOffset(nearby BBox, deltaY float64) float64 { return deltaY }
func (box BBox) ZOffset(nearby BBox, deltaZ float64) float64 { return deltaZ }
func (box BBox) Corners() []mgl64.Vec3 { return nil }
func (box BBox) Mul(val float64) BBox { return box }
func (box BBox) Volume() float64 { return 0 }
`

const bboxAxisSource = `package cube
type Axis int
const (
	Y Axis = iota
	Z
	X
)
`

func writeBBoxFixture(t *testing.T, bbox, axis string) (string, string) {
	t.Helper()
	directory := t.TempDir()
	bboxPath := filepath.Join(directory, "bbox.go")
	axisPath := filepath.Join(directory, "axis.go")
	if err := os.WriteFile(bboxPath, []byte(bbox), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(axisPath, []byte(axis), 0o600); err != nil {
		t.Fatal(err)
	}
	return bboxPath, axisPath
}

func TestBBoxUsesGoAST(t *testing.T) {
	bboxPath, axisPath := writeBBoxFixture(t, bboxSource, bboxAxisSource)
	spec, err := inspectBBox(bboxPath, axisPath)
	if err != nil {
		t.Fatal(err)
	}
	output := string(generateBBox(spec))
	for _, expected := range []string{
		"public enum Axis",
		"Y = 0,",
		"Z = 1,",
		"X = 2,",
		"public static BBox Box(double x0, double y0, double z0, double x1, double y1, double z1)",
		"public static bool AnyIntersections(BBox[] boxes, BBox search)",
		"public readonly record struct BBox",
		"internal BBox(Vector3 min, Vector3 max)",
		"public BBox Grow(double x)",
		"public BBox GrowVec3(Vector3 vec)",
		"public Vector3 Min()",
		"public Vector3 Max()",
		"public double Width()",
		"public double Length()",
		"public double Height()",
		"public BBox Extend(Vector3 vec)",
		"public BBox ExtendTowards(Face f, double x)",
		"public BBox Stretch(Axis a, double x)",
		"public BBox Translate(Vector3 vec)",
		"public BBox TranslateTowards(Face f, double x)",
		"public bool IntersectsWith(BBox other)",
		"public bool Vec3Within(Vector3 vec)",
		"public bool Vec3WithinYZ(Vector3 vec)",
		"public bool Vec3WithinXZ(Vector3 vec)",
		"public bool Vec3WithinXY(Vector3 vec)",
		"public double XOffset(BBox nearby, double deltaX)",
		"public double YOffset(BBox nearby, double deltaY)",
		"public double ZOffset(BBox nearby, double deltaZ)",
		"public Vector3[] Corners()",
		"public BBox Mul(double val)",
		"public double Volume()",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("generated BBox surface missing %q:\n%s", expected, output)
		}
	}
}

func TestBBoxRejectsSignatureDrift(t *testing.T) {
	tests := map[string][2]string{
		"field":       {"min, max mgl64.Vec3", "min, max [3]float64"},
		"constructor": {"z1 float64) BBox", "z1 float32) BBox"},
		"method":      {"Grow(x float64) BBox", "Grow(x float32) BBox"},
		"param name":  {"Grow(x float64) BBox", "Grow(amount float64) BBox"},
		"receiver":    {"func (box BBox) Grow(x float64)", "func (box *BBox) Grow(x float64)"},
		"new export":  {"func (box BBox) Volume() float64", "func (box BBox) Future() float64 { return 0 }\nfunc (box BBox) Volume() float64"},
	}
	for name, replacement := range tests {
		t.Run(name, func(t *testing.T) {
			bboxPath, axisPath := writeBBoxFixture(t,
				strings.Replace(bboxSource, replacement[0], replacement[1], 1), bboxAxisSource)
			if _, err := inspectBBox(bboxPath, axisPath); err == nil || !strings.Contains(err.Error(), "changed") {
				t.Fatalf("expected signature drift error, got %v", err)
			}
		})
	}
}

func TestBBoxRejectsAxisDrift(t *testing.T) {
	bboxPath, axisPath := writeBBoxFixture(t, bboxSource, strings.Replace(bboxAxisSource, "\tZ\n\tX", "\tX\n\tZ", 1))
	if _, err := inspectBBox(bboxPath, axisPath); err == nil || !strings.Contains(err.Error(), "Axis values changed") {
		t.Fatalf("expected Axis drift error, got %v", err)
	}
}

func TestGeneratedBBoxBehavior(t *testing.T) {
	if _, err := exec.LookPath("dotnet"); err != nil {
		t.Skip("dotnet is not installed")
	}
	directory := t.TempDir()
	files := map[string]string{
		"BBoxBehavior.csproj": `<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <OutputType>Exe</OutputType>
    <TargetFramework>net10.0</TargetFramework>
    <ImplicitUsings>enable</ImplicitUsings>
  </PropertyGroup>
</Project>
`,
		"Primitives.cs": `namespace Dragonfly;
public readonly record struct Vector3(double X, double Y, double Z);
public static partial class Cube
{
    public enum Face { Down, Up, North, South, West, East }
}
`,
		"Cube.BBox.g.cs": string(generateBBox(bboxSpec{Axes: []string{"Y", "Z", "X"}})),
		"Program.cs": `using Dragonfly;

static void Assert(bool condition, string message)
{
    if (!condition) throw new Exception(message);
}

var normalized = Cube.Box(3, 4, 5, 1, 2, -1);
Assert(normalized.Min() == new Vector3(1, 2, -1), "Box minimum normalization");
Assert(normalized.Max() == new Vector3(3, 4, 5), "Box maximum normalization");

var unit = Cube.Box(0, 0, 0, 1, 1, 1);
Assert(unit.Vec3Within(new Vector3(.5, .5, .5)), "strict interior");
Assert(!unit.Vec3Within(new Vector3(0, .5, .5)), "strict X boundary");
Assert(!unit.Vec3Within(new Vector3(.5, 1, .5)), "strict Y boundary");
Assert(unit.Vec3WithinYZ(new Vector3(99, 0, 1)), "inclusive YZ boundary");
Assert(unit.Vec3WithinXZ(new Vector3(0, 99, 1)), "inclusive XZ boundary");
Assert(unit.Vec3WithinXY(new Vector3(0, 1, 99)), "inclusive XY boundary");

var epsilonMiss = Cube.Box(.999995, 0, 0, 2, 1, 1);
var epsilonHit = Cube.Box(.99998, 0, 0, 2, 1, 1);
Assert(!unit.IntersectsWith(epsilonMiss), "intersection epsilon");
Assert(unit.IntersectsWith(epsilonHit), "intersection above epsilon");
Assert(Cube.AnyIntersections([unit], epsilonMiss), "zero-epsilon collection intersection");

var invertedGrow = Cube.Box(0, 0, 0, 2, 2, 2).Grow(-2);
Assert(invertedGrow.Min() == new Vector3(2, 2, 2), "negative Grow minimum");
Assert(invertedGrow.Max() == new Vector3(0, 0, 0), "negative Grow maximum");
Assert(invertedGrow.Width() == -2, "negative Grow remains non-normalized");

var negativeMul = Cube.Box(1, 2, 3, 4, 5, 6).Mul(-1);
Assert(negativeMul.Min() == new Vector3(-1, -2, -3), "negative Mul minimum");
Assert(negativeMul.Max() == new Vector3(-4, -5, -6), "negative Mul remains non-normalized");

var corners = Cube.Box(1, 2, 3, 4, 5, 6).Corners();
var wantCorners = new[] {
    new Vector3(1, 2, 3), new Vector3(4, 5, 6),
    new Vector3(1, 2, 6), new Vector3(1, 5, 3),
    new Vector3(1, 5, 6), new Vector3(4, 5, 3),
    new Vector3(4, 2, 6), new Vector3(4, 2, 3),
};
Assert(corners.SequenceEqual(wantCorners), "corner order");
`,
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(directory, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	command := exec.Command("dotnet", "run", "--project", filepath.Join(directory, "BBoxBehavior.csproj"), "--nologo")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("generated BBox behavior program failed: %v\n%s", err, output)
	}
}

func TestPinnedDragonflyBBoxHasExactSupportedSurface(t *testing.T) {
	command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly")
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	directory := filepath.Join(string(bytes.TrimSpace(output)), "server", "block", "cube")
	if _, err := inspectBBox(filepath.Join(directory, "bbox.go"), filepath.Join(directory, "axis.go")); err != nil {
		t.Fatal(err)
	}
}
