namespace Dragonfly;

public sealed partial class World
{
    public interface Item { }

    public partial class Tx
    {
        internal Tx(ulong invocation) => Invocation = invocation;
        internal ulong Invocation { get; }
    }

    public class Context : Tx
    {
        private bool _cancelled;

        internal Context(ulong invocation, bool cancelled) : base(invocation) =>
            _cancelled = cancelled;

        public bool Cancelled() => _cancelled;
        public void Cancel() => _cancelled = true;
    }
}
