// Code generated from Dragonfly server/world/entity.go Go AST. DO NOT EDIT.
namespace Dragonfly;

public sealed partial class World
{
    public interface Entity
    {
        void Close();
        EntityHandle H();
        Vector3 Position();
        Rotation Rotation();
    }
}
