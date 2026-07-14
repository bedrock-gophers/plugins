// Code generated from Dragonfly server/cmd Go AST. DO NOT EDIT.
#nullable enable
namespace Dragonfly;

public static partial class Cmd
{
    public interface Runnable
    {
        void Run(Source src, Output o, World.Tx? tx);
    }

    public interface Allower
    {
        bool Allow(Source src);
    }

    public interface Target
    {
        Vector3 Position();
    }

    public interface NamedTarget : Target
    {
        string Name();
    }

    public interface Source : Target
    {
        void SendCommandOutput(Output o);
    }

    public interface Enum
    {
        string Type();
        IReadOnlyList<string> Options(Source source);
    }
}
