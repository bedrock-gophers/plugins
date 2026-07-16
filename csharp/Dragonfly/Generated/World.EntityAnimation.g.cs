// Code generated from Dragonfly server/world/entity_animation.go Go AST. DO NOT EDIT.
#nullable enable

namespace Dragonfly;

public sealed partial class World
{
    public static EntityAnimation NewEntityAnimation(string name) => new(name);

    public readonly struct EntityAnimation
    {
        private readonly string? _name;
        private readonly string? _nextState;
        private readonly string? _controller;
        private readonly string? _stopCondition;

        internal EntityAnimation(string name, string nextState = "", string controller = "", string stopCondition = "") =>
            (_name, _nextState, _controller, _stopCondition) =
                (name ?? throw new ArgumentNullException(nameof(name)), nextState, controller, stopCondition);

        public string Name() => _name ?? string.Empty;

        public string Controller() => _controller ?? string.Empty;

        public EntityAnimation WithController(string controller) =>
            new(Name(), NextState(), controller ?? throw new ArgumentNullException(nameof(controller)), StopCondition());

        public string NextState() => _nextState ?? string.Empty;

        public EntityAnimation WithNextState(string state) =>
            new(Name(), state ?? throw new ArgumentNullException(nameof(state)), Controller(), StopCondition());

        public string StopCondition() => _stopCondition ?? string.Empty;

        public EntityAnimation WithStopCondition(string condition) =>
            new(Name(), NextState(), Controller(), condition ?? throw new ArgumentNullException(nameof(condition)));
    }
}
