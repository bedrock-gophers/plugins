// Code generated from Dragonfly damage/healing source Go AST. DO NOT EDIT.
#nullable enable

namespace Dragonfly;

public sealed partial class World
{
    public interface DamageSource
    {
        bool ReducedByArmour();
        bool ReducedByResistance();
        bool Fire();
        bool IgnoreTotem();
    }

    // Go's HealingSource() marker cannot be named on this C# interface: C# reserves
    // a member matching its enclosing type name for constructors.
    public interface HealingSource { }
}

public static partial class Entity
{
    public readonly record struct AttackDamageSource(World.Entity? Attacker = null) : World.DamageSource
    {
        public bool ReducedByArmour() => true;
        public bool ReducedByResistance() => true;
        public bool Fire() => false;
        public bool IgnoreTotem() => false;
    }

    public readonly record struct VoidDamageSource : World.DamageSource
    {
        public bool ReducedByArmour() => false;
        public bool ReducedByResistance() => false;
        public bool Fire() => false;
        public bool IgnoreTotem() => true;
    }

    public readonly record struct SuffocationDamageSource : World.DamageSource
    {
        public bool ReducedByArmour() => false;
        public bool ReducedByResistance() => false;
        public bool Fire() => false;
        public bool IgnoreTotem() => false;
    }

    public readonly record struct DrowningDamageSource : World.DamageSource
    {
        public bool ReducedByArmour() => false;
        public bool ReducedByResistance() => false;
        public bool Fire() => false;
        public bool IgnoreTotem() => false;
    }

    public readonly record struct FallDamageSource : Enchantment.AffectedDamageSource
    {
        public bool ReducedByArmour() => false;
        public bool ReducedByResistance() => true;
        public bool Fire() => false;
        public bool IgnoreTotem() => false;
        public bool AffectedByEnchantment(Item.EnchantmentType e) => object.Equals(e, Item.FeatherFalling);
    }

    public readonly record struct GlideDamageSource : World.DamageSource
    {
        public bool ReducedByArmour() => false;
        public bool ReducedByResistance() => true;
        public bool Fire() => false;
        public bool IgnoreTotem() => false;
    }

    public readonly record struct LightningDamageSource : World.DamageSource
    {
        public bool ReducedByArmour() => true;
        public bool ReducedByResistance() => true;
        public bool Fire() => false;
        public bool IgnoreTotem() => false;
    }

    public readonly record struct ProjectileDamageSource(World.Entity? Projectile = null, World.Entity? Owner = null) : Enchantment.AffectedDamageSource
    {
        public bool ReducedByArmour() => true;
        public bool ReducedByResistance() => true;
        public bool Fire() => false;
        public bool IgnoreTotem() => false;
        public bool AffectedByEnchantment(Item.EnchantmentType e) => object.Equals(e, Item.ProjectileProtection);
    }

    public readonly record struct ExplosionDamageSource : Enchantment.AffectedDamageSource
    {
        public bool ReducedByArmour() => true;
        public bool ReducedByResistance() => true;
        public bool Fire() => false;
        public bool IgnoreTotem() => false;
        public bool AffectedByEnchantment(Item.EnchantmentType e) => object.Equals(e, Item.BlastProtection);
    }

    public readonly record struct FoodHealingSource(bool QuickRegeneration = false) : World.HealingSource;
}

public static partial class Effect
{
    public readonly record struct WitherDamageSource : World.DamageSource
    {
        public bool ReducedByArmour() => false;
        public bool ReducedByResistance() => true;
        public bool Fire() => false;
        public bool IgnoreTotem() => false;
    }

    public readonly record struct InstantDamageSource : World.DamageSource
    {
        public bool ReducedByArmour() => false;
        public bool ReducedByResistance() => true;
        public bool Fire() => false;
        public bool IgnoreTotem() => false;
    }

    public readonly record struct PoisonDamageSource(bool Fatal = false) : World.DamageSource
    {
        public bool ReducedByArmour() => false;
        public bool ReducedByResistance() => true;
        public bool Fire() => false;
        public bool IgnoreTotem() => false;
    }

    public readonly record struct InstantHealingSource : World.HealingSource;

    public readonly record struct RegenerationHealingSource : World.HealingSource;
}

public sealed partial class Player
{
    public readonly record struct StarvationDamageSource : World.DamageSource
    {
        public bool ReducedByArmour() => false;
        public bool ReducedByResistance() => false;
        public bool Fire() => false;
        public bool IgnoreTotem() => false;
    }
}

public static partial class Block
{
    public readonly record struct DamageSource(World.Block? Block = null) : World.DamageSource
    {
        public bool ReducedByArmour() => true;
        public bool ReducedByResistance() => true;
        public bool Fire() => false;
        public bool IgnoreTotem() => false;
    }

    public readonly record struct MagmaDamageSource : Enchantment.AffectedDamageSource
    {
        public bool ReducedByArmour() => true;
        public bool ReducedByResistance() => true;
        public bool Fire() => true;
        public bool IgnoreTotem() => false;
        public bool AffectedByEnchantment(Item.EnchantmentType e) => object.Equals(e, Item.FireProtection);
    }

    public readonly record struct LavaDamageSource : World.DamageSource
    {
        public bool ReducedByArmour() => true;
        public bool ReducedByResistance() => true;
        public bool Fire() => true;
        public bool IgnoreTotem() => false;
    }

    public readonly record struct FireDamageSource : Enchantment.AffectedDamageSource
    {
        public bool ReducedByArmour() => true;
        public bool ReducedByResistance() => true;
        public bool Fire() => true;
        public bool IgnoreTotem() => false;
        public bool AffectedByEnchantment(Item.EnchantmentType e) => object.Equals(e, Item.FireProtection);
    }
}

public static partial class Enchantment
{
    public interface AffectedDamageSource : World.DamageSource
    {
        bool AffectedByEnchantment(Item.EnchantmentType e);
    }

    public readonly record struct ThornsDamageSource(World.Entity? Owner = null) : World.DamageSource
    {
        public bool ReducedByArmour() => false;
        public bool ReducedByResistance() => true;
        public bool Fire() => false;
        public bool IgnoreTotem() => false;
    }
}
