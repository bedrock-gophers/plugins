// Code generated from Dragonfly server/block/cube Go AST. DO NOT EDIT.
namespace Dragonfly;

public static partial class Cube
{
    public enum Face
    {
        Down = 0,
        Up = 1,
        North = 2,
        South = 3,
        West = 4,
        East = 5,
    }

    public readonly record struct Range(int Minimum, int Maximum)
    {
        public int Min() => Minimum;
        public int Max() => Maximum;
        public int Height() => Maximum - Minimum;
    }

    public readonly record struct Pos
    {
        private readonly int _x;
        private readonly int _y;
        private readonly int _z;

        public Pos(int x, int y, int z) => (_x, _y, _z) = (x, y, z);

        public int X() => _x;
        public int Y() => _y;
        public int Z() => _z;
        public bool OutOfBounds(Range range) => _y > range.Max() || _y < range.Min();
        public bool Within(Pos min, Pos max) =>
            _x >= min._x && _x <= max._x &&
            _y >= min._y && _y <= max._y &&
            _z >= min._z && _z <= max._z;
        public Pos Add(Pos other) => new(_x + other._x, _y + other._y, _z + other._z);
        public Pos Sub(Pos other) => new(_x - other._x, _y - other._y, _z - other._z);
        public Vector3 Vec3() => new(_x, _y, _z);
        public Vector3 Vec3Middle() => new(_x + 0.5, _y, _z + 0.5);
        public Vector3 Vec3Centre() => new(_x + 0.5, _y + 0.5, _z + 0.5);

        public Pos Side(Face face) => face switch
        {
            Cube.Face.Up => new(_x, _y + 1, _z),
            Cube.Face.Down => new(_x, _y - 1, _z),
            Cube.Face.North => new(_x, _y, _z - 1),
            Cube.Face.South => new(_x, _y, _z + 1),
            Cube.Face.West => new(_x - 1, _y, _z),
            Cube.Face.East => new(_x + 1, _y, _z),
            _ => this,
        };

        public Face Face(Pos other) => NeighbourFace(other).Face;

        public (Face Face, bool Ok) NeighbourFace(Pos other) => other.Sub(this) switch
        {
            Pos { _x: 0, _y: 1, _z: 0 } => (Cube.Face.Up, true),
            Pos { _x: 0, _y: -1, _z: 0 } => (Cube.Face.Down, true),
            Pos { _x: 0, _y: 0, _z: -1 } => (Cube.Face.North, true),
            Pos { _x: 0, _y: 0, _z: 1 } => (Cube.Face.South, true),
            Pos { _x: -1, _y: 0, _z: 0 } => (Cube.Face.West, true),
            Pos { _x: 1, _y: 0, _z: 0 } => (Cube.Face.East, true),
            _ => (Cube.Face.Up, false),
        };

        public override string ToString() => $"({_x},{_y},{_z})";
    }

    public static Pos PosFromVec3(Vector3 value) => new(
        checked((int)Math.Floor(value.X)),
        checked((int)Math.Floor(value.Y)),
        checked((int)Math.Floor(value.Z)));

    public static Pos Min(Pos first, Pos second) => new(
        Math.Min(first.X(), second.X()),
        Math.Min(first.Y(), second.Y()),
        Math.Min(first.Z(), second.Z()));

    public static Pos Max(Pos first, Pos second) => new(
        Math.Max(first.X(), second.X()),
        Math.Max(first.Y(), second.Y()),
        Math.Max(first.Z(), second.Z()));
}
