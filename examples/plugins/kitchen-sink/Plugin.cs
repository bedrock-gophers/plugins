#nullable enable
using System;
using System.Collections.Generic;
using System.Threading;
using Dragonfly;

public sealed class KitchenSink : Plugin
{
    private long _jumps;
    private long _punches;
    private long _quits;
    private long _sneaks;
    private long _sprints;

    public override void OnEnable()
    {
        Cmd.Register(Cmd.New(
            "kitchen",
            "Exercises the reflected C# command API.",
            ["ks"],
            new KitchenStatus(this),
            new KitchenEcho(),
            new KitchenMode(),
            new KitchenPing(),
            new KitchenPosition(),
            new KitchenDestination(),
            new KitchenText(),
            new KitchenBlock(),
            new KitchenBiome(),
            new KitchenTick(),
            new KitchenParticle()));
        Console.WriteLine("kitchen-sink enabled");
    }

    public override void OnDisable() => Console.WriteLine(
        $"kitchen-sink disabled: jumps={_jumps}, punches={_punches}, " +
        $"sprints={_sprints}, sneaks={_sneaks}, quits={_quits}");

    public override void HandleMove(Player.Context ctx, Vector3 newPos, Rotation newRot)
    {
        if (!Finite(newPos) || !double.IsFinite(newRot.Yaw) || !double.IsFinite(newRot.Pitch))
            ctx.Cancel();
    }

    public override void HandleJump(Player player) => Increment(ref _jumps);

    public override void HandleTeleport(Player.Context ctx, Vector3 pos)
    {
        if (!Finite(pos)) ctx.Cancel();
    }

    public override void HandleToggleSprint(Player.Context ctx, bool sprinting)
    {
        if (sprinting) Increment(ref _sprints);
    }

    public override void HandleToggleSneak(Player.Context ctx, bool sneaking)
    {
        if (sneaking) Increment(ref _sneaks);
    }

    public override void HandleChat(Player.Context ctx, ref string message) =>
        message = message.Trim();

    public override void HandleFoodLoss(Player.Context ctx, int from, ref int to) =>
        to = Math.Clamp(to, 0, 20);

    public override void HandlePunchAir(Player.Context ctx) => Increment(ref _punches);
    public override void HandleQuit(Player player) => Increment(ref _quits);

    private static void Increment(ref long counter) => Interlocked.Increment(ref counter);

    private static bool Finite(Vector3 value) =>
        double.IsFinite(value.X) && double.IsFinite(value.Y) && double.IsFinite(value.Z);

    internal sealed class KitchenStatus(KitchenSink plugin) : Cmd.Runnable
    {
        public void Run(Cmd.Source source, Cmd.Output output, World.Tx? tx) => output.Printf(
            "jumps={0}, punches={1}, sprints={2}, sneaks={3}, quits={4}",
            plugin._jumps,
            plugin._punches,
            plugin._sprints,
            plugin._sneaks,
            plugin._quits);
    }

    internal sealed class KitchenEcho : Cmd.Runnable
    {
        public Cmd.SubCommand Echo;
        public Cmd.Varargs Message;

        public void Run(Cmd.Source source, Cmd.Output output, World.Tx? tx) => output.Print(Message);
    }

    internal enum Mode
    {
        Survival,
        Creative,
        Adventure,
        Spectator,
    }

    internal sealed class KitchenMode : Cmd.Runnable
    {
        public Cmd.SubCommand Mode;
        public Mode Value;

        public void Run(Cmd.Source source, Cmd.Output output, World.Tx? tx) => output.Printf("mode={0}", Value);
    }

    internal sealed class KitchenPing : Cmd.Runnable
    {
        public Cmd.SubCommand Ping;
        public Cmd.Optional<Player> Target;

        public void Run(Cmd.Source source, Cmd.Output output, World.Tx? tx)
        {
            var (target, _) = Target.Load();
            target ??= source as Player;
            if (target is null)
            {
                output.Error("A player target is required from the console.");
                return;
            }
            output.Printf("{0}'s ping: {1}ms", target.Name(), target.Latency().TotalMilliseconds);
        }
    }

    internal sealed class KitchenPosition : Cmd.Runnable
    {
        public Cmd.SubCommand Position;
        public Vector3 Value;

        public void Run(Cmd.Source source, Cmd.Output output, World.Tx? tx) =>
            output.Printf("position={0},{1},{2}", Value.X, Value.Y, Value.Z);
    }

    internal sealed class KitchenDestination : Cmd.Runnable
    {
        public Cmd.SubCommand Destination;
        public Destination Value = new();

        public void Run(Cmd.Source source, Cmd.Output output, World.Tx? tx) => output.Printf("destination={0}", Value);
    }

    internal sealed class Destination(string value = "spawn") : Cmd.Enum
    {
        public string Type() => "kitchen_destination";
        public IReadOnlyList<string> Options(Cmd.Source source) => ["spawn", "source"];
        public override string ToString() => value;
    }

    internal enum TextAction
    {
        Message,
        Popup,
        Tip,
        Jukebox,
        NameTag,
        Disconnect,
    }

    internal sealed class KitchenText : Cmd.Runnable
    {
        [Cmd.Tag("text")]
        public Cmd.SubCommand Text;
        [Cmd.Tag("action")]
        public TextAction Action;
        public Cmd.Varargs Content;

        public void Run(Cmd.Source source, Cmd.Output output, World.Tx? tx)
        {
            if (source is not Player player)
            {
                output.Error("This command can only be used by a player.");
                return;
            }
            switch (Action)
            {
                case TextAction.Message:
                    player.Message(Content, true, 12, 1.5, null);
                    break;
                case TextAction.Popup:
                    player.SendPopup(Content);
                    break;
                case TextAction.Tip:
                    player.SendTip(Content);
                    break;
                case TextAction.Jukebox:
                    player.SendJukeboxPopup(Content);
                    break;
                case TextAction.NameTag:
                    player.SetNameTag(Content);
                    break;
                case TextAction.Disconnect:
                    player.Disconnect(Content);
                    break;
            }
        }
    }

    internal sealed class KitchenBlock : Cmd.Runnable
    {
        public Cmd.SubCommand Block;

        public void Run(Cmd.Source source, Cmd.Output output, World.Tx? tx)
        {
            if (tx is null)
            {
                output.Error("A world transaction is required.");
                return;
            }
            var position = Cube.PosFromVec3(source.Position()).Side(Cube.Face.Down);
            var range = tx.Range();
            var (block, loaded) = tx.BlockLoaded(position);
            var previous = loaded ? block : tx.Block(position);
            var wasSand = previous is Block.Sand;
            var nearbySand = tx.BlocksWithin(position, 2, new Block.Sand());
            var highestLightBlocker = tx.HighestLightBlocker(position.X(), position.Z());
            var highestBlock = tx.HighestBlock(position.X(), position.Z());
            var light = tx.Light(position);
            var skyLight = tx.SkyLight(position);
            var (_, liquidBeforeFound) = tx.Liquid(position);
            tx.SetBlock(position, new Block.Sand(), new World.SetOpts
            {
                DisableBlockUpdates = true,
                DisableLiquidDisplacement = true,
                DisableRedstoneUpdates = true,
            });
            var scheduledWater = new Block.Water(Still: true, Depth: 8, Falling: false);
            tx.SetLiquid(position, scheduledWater);
            var (liquid, liquidFound) = tx.Liquid(position);
            var liquidState = liquid is Block.Water water
                ? $"Water(still={(water.Still ? "true" : "false")},depth={water.Depth}," +
                  $"falling={(water.Falling ? "true" : "false")})"
                : liquid?.GetType().Name ?? "none";
            tx.SetLiquid(position, null);
            tx.SetLiquid(position, scheduledWater);
            var blockUpdateDelay = TimeSpan.FromMilliseconds(250);
            tx.ScheduleBlockUpdate(position, scheduledWater, blockUpdateDelay);
            var firstNearbySand = "none";
            foreach (var nearbyPosition in nearbySand)
            {
                firstNearbySand = nearbyPosition.ToString();
                break;
            }
            output.Printf(
                "block={0}, range={1}..{2}, loaded={3}, was_sand={4}, nearby_sand={5}, " +
                "highest_light_blocker={6}, highest_block={7}, light={8}, sky_light={9}, " +
                "liquid_before={10}, liquid={11}:{12}, scheduled_update=water:{13}ms",
                position,
                range.Min(),
                range.Max(),
                loaded ? "true" : "false",
                wasSand ? "true" : "false",
                firstNearbySand,
                highestLightBlocker,
                highestBlock,
                light,
                skyLight,
                liquidBeforeFound ? "true" : "false",
                liquidFound ? "true" : "false",
                liquidState,
                blockUpdateDelay.TotalMilliseconds);
        }
    }

    internal sealed class KitchenBiome : Cmd.Runnable
    {
        public Cmd.SubCommand Biome;

        public void Run(Cmd.Source source, Cmd.Output output, World.Tx? tx)
        {
            if (tx is null)
            {
                output.Error("A world transaction is required.");
                return;
            }
            var position = Cube.PosFromVec3(source.Position());
            var previous = tx.Biome(position);
            World.Biome current = previous;
            var temperature = 0.0;
            var rainingAt = false;
            var snowingAt = false;
            var thunderingAt = false;
            var raining = false;
            var thundering = false;
            tx.SetBiome(position, new Biome.Desert());
            try
            {
                current = tx.Biome(position);
                temperature = tx.Temperature(position);
                rainingAt = tx.RainingAt(position);
                snowingAt = tx.SnowingAt(position);
                thunderingAt = tx.ThunderingAt(position);
                raining = tx.Raining();
                thundering = tx.Thundering();
            }
            finally
            {
                tx.SetBiome(position, previous);
            }
            var restored = tx.Biome(position);
            output.Printf(
                "biome=Desert, applied={0}, temperature={1}, raining_at={2}, snowing_at={3}, " +
                "thundering_at={4}, raining={5}, thundering={6}, restored={7}",
                current is Biome.Desert ? "true" : "false",
                temperature,
                rainingAt ? "true" : "false",
                snowingAt ? "true" : "false",
                thunderingAt ? "true" : "false",
                raining ? "true" : "false",
                thundering ? "true" : "false",
                restored.Equals(previous) ? "true" : "false");
        }
    }

    internal sealed class KitchenTick : Cmd.Runnable
    {
        public Cmd.SubCommand Tick;

        public void Run(Cmd.Source source, Cmd.Output output, World.Tx? tx)
        {
            if (tx is null)
            {
                output.Error("A world transaction is required.");
                return;
            }
            output.Printf("tick={0}", tx.CurrentTick());
        }
    }

    internal sealed class KitchenParticle : Cmd.Runnable
    {
        public Cmd.SubCommand Particle;

        public void Run(Cmd.Source source, Cmd.Output output, World.Tx? tx)
        {
            if (tx is null)
            {
                output.Error("A world transaction is required.");
                return;
            }
            World.Particle[] particles =
            [
                new Particle.Flame(new Color.RGBA(1, 2, 3, 4)),
                new Particle.Dust(new Color.RGBA(5, 6, 7, 8)),
                new Particle.BlockBreak(new Block.Sand()),
                new Particle.PunchBlock(new Block.Sand(), Cube.Face.East),
                new Particle.BlockForceField(),
                new Particle.BoneMeal(true),
                new Particle.Note(Sound.Piano(), 24),
                new Particle.Note(Sound.BassDrum(), 24),
                new Particle.Note(Sound.Snare(), 24),
                new Particle.Note(Sound.ClicksAndSticks(), 24),
                new Particle.Note(Sound.Bass(), 24),
                new Particle.Note(Sound.Flute(), 24),
                new Particle.Note(Sound.Bell(), 24),
                new Particle.Note(Sound.Guitar(), 24),
                new Particle.Note(Sound.Chimes(), 24),
                new Particle.Note(Sound.Xylophone(), 24),
                new Particle.Note(Sound.IronXylophone(), 24),
                new Particle.Note(Sound.CowBell(), 24),
                new Particle.Note(Sound.Didgeridoo(), 24),
                new Particle.Note(Sound.Bit(), 24),
                new Particle.Note(Sound.Banjo(), 24),
                new Particle.Note(Sound.Pling(), 24),
                new Particle.DragonEggTeleport(new Cube.Pos(-3, 4, 5)),
                new Particle.Evaporate(),
                new Particle.WaterDrip(),
                new Particle.LavaDrip(),
                new Particle.Lava(),
                new Particle.DustPlume(),
                new Particle.HugeExplosion(),
                new Particle.EndermanTeleport(),
                new Particle.SnowballPoof(),
                new Particle.EggSmash(),
                new Particle.Splash(new Color.RGBA(9, 10, 11, 12)),
                new Particle.Effect(new Color.RGBA(13, 14, 15, 16)),
                new Particle.EntityFlame(),
            ];
            foreach (var particle in particles)
                tx.AddParticle(source.Position(), particle);
            output.Printf("particles={0}", particles.Length);
        }
    }
}
