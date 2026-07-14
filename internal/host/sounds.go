package host

import (
	"math"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/block"
	dfitem "github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	dfsound "github.com/df-mc/dragonfly/server/world/sound"
)

const soundFlagPrimary uint32 = 1

// ValidSound validates payload fields that do not depend on a world registry.
func ValidSound(value native.WorldSound) bool {
	if value.Kind > native.SoundGoatHorn || math.IsNaN(value.Scalar) || math.IsInf(value.Scalar, 0) {
		return false
	}
	switch value.Kind {
	case native.SoundAttack:
		return value.Flags <= soundFlagPrimary
	case native.SoundBlockPlace, native.SoundBlockBreaking, native.SoundDoorOpen, native.SoundDoorClose,
		native.SoundTrapdoorOpen, native.SoundTrapdoorClose, native.SoundFenceGateOpen, native.SoundFenceGateClose,
		native.SoundItemUseOn:
		return validSoundBlock(value.Block)
	case native.SoundNote:
		return value.Data < 16
	case native.SoundMusicDiscPlay:
		return value.Data < 21
	case native.SoundEquipItem:
		return value.Item != nil
	case native.SoundBucketFill, native.SoundBucketEmpty:
		return value.Data < 2
	case native.SoundCrossbowLoad:
		return value.Integer >= dfsound.CrossbowLoadingStart && value.Integer <= dfsound.CrossbowLoadingEnd && value.Flags <= soundFlagPrimary
	case native.SoundGoatHorn:
		return value.Data < 8
	default:
		return true
	}
}

func validSoundBlock(value *native.WorldBlock) bool {
	if value == nil || value.Identifier == "" {
		return false
	}
	_, ok := DecodeBlockProperties(value.PropertiesNBT)
	return ok
}

// SoundFromNative resolves a validated descriptor through the transaction's registries.
func SoundFromNative(tx *world.Tx, value native.WorldSound) (world.Sound, bool) {
	if tx == nil || !ValidSound(value) {
		return nil, false
	}
	switch value.Kind {
	case native.SoundAnvilBreak:
		return dfsound.AnvilBreak{}, true
	case native.SoundAnvilLand:
		return dfsound.AnvilLand{}, true
	case native.SoundAnvilUse:
		return dfsound.AnvilUse{}, true
	case native.SoundArrowHit:
		return dfsound.ArrowHit{}, true
	case native.SoundBarrelClose:
		return dfsound.BarrelClose{}, true
	case native.SoundBarrelOpen:
		return dfsound.BarrelOpen{}, true
	case native.SoundBlastFurnaceCrackle:
		return dfsound.BlastFurnaceCrackle{}, true
	case native.SoundBowShoot:
		return dfsound.BowShoot{}, true
	case native.SoundBurning:
		return dfsound.Burning{}, true
	case native.SoundBurp:
		return dfsound.Burp{}, true
	case native.SoundCampfireCrackle:
		return dfsound.CampfireCrackle{}, true
	case native.SoundChestClose:
		return dfsound.ChestClose{}, true
	case native.SoundChestOpen:
		return dfsound.ChestOpen{}, true
	case native.SoundClick:
		return dfsound.Click{}, true
	case native.SoundComposterEmpty:
		return dfsound.ComposterEmpty{}, true
	case native.SoundComposterFill:
		return dfsound.ComposterFill{}, true
	case native.SoundComposterFillLayer:
		return dfsound.ComposterFillLayer{}, true
	case native.SoundComposterReady:
		return dfsound.ComposterReady{}, true
	case native.SoundCopperScraped:
		return dfsound.CopperScraped{}, true
	case native.SoundCrossbowShoot:
		return dfsound.CrossbowShoot{}, true
	case native.SoundDecoratedPotInsertFailed:
		return dfsound.DecoratedPotInsertFailed{}, true
	case native.SoundDeny:
		return dfsound.Deny{}, true
	case native.SoundDoorCrash:
		return dfsound.DoorCrash{}, true
	case native.SoundDrowning:
		return dfsound.Drowning{}, true
	case native.SoundEnderChestClose:
		return dfsound.EnderChestClose{}, true
	case native.SoundEnderChestOpen:
		return dfsound.EnderChestOpen{}, true
	case native.SoundExperience:
		return dfsound.Experience{}, true
	case native.SoundExplosion:
		return dfsound.Explosion{}, true
	case native.SoundFireCharge:
		return dfsound.FireCharge{}, true
	case native.SoundFireExtinguish:
		return dfsound.FireExtinguish{}, true
	case native.SoundFireworkBlast:
		return dfsound.FireworkBlast{}, true
	case native.SoundFireworkHugeBlast:
		return dfsound.FireworkHugeBlast{}, true
	case native.SoundFireworkLaunch:
		return dfsound.FireworkLaunch{}, true
	case native.SoundFireworkTwinkle:
		return dfsound.FireworkTwinkle{}, true
	case native.SoundFizz:
		return dfsound.Fizz{}, true
	case native.SoundFurnaceCrackle:
		return dfsound.FurnaceCrackle{}, true
	case native.SoundGhastShoot:
		return dfsound.GhastShoot{}, true
	case native.SoundGhastWarning:
		return dfsound.GhastWarning{}, true
	case native.SoundGlassBreak:
		return dfsound.GlassBreak{}, true
	case native.SoundIgnite:
		return dfsound.Ignite{}, true
	case native.SoundItemAdd:
		return dfsound.ItemAdd{}, true
	case native.SoundItemBreak:
		return dfsound.ItemBreak{}, true
	case native.SoundItemFrameRemove:
		return dfsound.ItemFrameRemove{}, true
	case native.SoundItemFrameRotate:
		return dfsound.ItemFrameRotate{}, true
	case native.SoundItemThrow:
		return dfsound.ItemThrow{}, true
	case native.SoundLecternBookPlace:
		return dfsound.LecternBookPlace{}, true
	case native.SoundLevelUp:
		return dfsound.LevelUp{}, true
	case native.SoundLightningExplode:
		return dfsound.LightningExplode{}, true
	case native.SoundLightningThunder:
		return dfsound.LightningThunder{}, true
	case native.SoundMusicDiscEnd:
		return dfsound.MusicDiscEnd{}, true
	case native.SoundPop:
		return dfsound.Pop{}, true
	case native.SoundPotionBrewed:
		return dfsound.PotionBrewed{}, true
	case native.SoundPowerOff:
		return dfsound.PowerOff{}, true
	case native.SoundPowerOn:
		return dfsound.PowerOn{}, true
	case native.SoundSignWaxed:
		return dfsound.SignWaxed{}, true
	case native.SoundSmokerCrackle:
		return dfsound.SmokerCrackle{}, true
	case native.SoundStopUsingSpyglass:
		return dfsound.StopUsingSpyglass{}, true
	case native.SoundTnt:
		return dfsound.TNT{}, true
	case native.SoundTeleport:
		return dfsound.Teleport{}, true
	case native.SoundThunder:
		return dfsound.Thunder{}, true
	case native.SoundTotem:
		return dfsound.Totem{}, true
	case native.SoundUseSpyglass:
		return dfsound.UseSpyglass{}, true
	case native.SoundWaxRemoved:
		return dfsound.WaxRemoved{}, true
	case native.SoundWaxedSignFailedInteraction:
		return dfsound.WaxedSignFailedInteraction{}, true
	case native.SoundShulkerBoxOpen:
		return dfsound.ShulkerBoxOpen{}, true
	case native.SoundShulkerBoxClose:
		return dfsound.ShulkerBoxClose{}, true
	case native.SoundEnderEyePlaced:
		return dfsound.EnderEyePlaced{}, true
	case native.SoundEndPortalCreated:
		return dfsound.EndPortalCreated{}, true
	case native.SoundAttack:
		return dfsound.Attack{Damage: value.Flags&soundFlagPrimary != 0}, true
	case native.SoundFall:
		return dfsound.Fall{Distance: value.Scalar}, true
	case native.SoundBlockPlace:
		return blockSound(tx, value.Block, func(value world.Block) world.Sound { return dfsound.BlockPlace{Block: value} })
	case native.SoundBlockBreaking:
		return blockSound(tx, value.Block, func(value world.Block) world.Sound { return dfsound.BlockBreaking{Block: value} })
	case native.SoundDoorOpen:
		return blockSound(tx, value.Block, func(value world.Block) world.Sound { return dfsound.DoorOpen{Block: value} })
	case native.SoundDoorClose:
		return blockSound(tx, value.Block, func(value world.Block) world.Sound { return dfsound.DoorClose{Block: value} })
	case native.SoundTrapdoorOpen:
		return blockSound(tx, value.Block, func(value world.Block) world.Sound { return dfsound.TrapdoorOpen{Block: value} })
	case native.SoundTrapdoorClose:
		return blockSound(tx, value.Block, func(value world.Block) world.Sound { return dfsound.TrapdoorClose{Block: value} })
	case native.SoundFenceGateOpen:
		return blockSound(tx, value.Block, func(value world.Block) world.Sound { return dfsound.FenceGateOpen{Block: value} })
	case native.SoundFenceGateClose:
		return blockSound(tx, value.Block, func(value world.Block) world.Sound { return dfsound.FenceGateClose{Block: value} })
	case native.SoundNote:
		instrument, ok := soundInstrument(value.Data)
		return dfsound.Note{Instrument: instrument, Pitch: int(value.Integer)}, ok
	case native.SoundMusicDiscPlay:
		disc, ok := soundDisc(value.Data)
		return dfsound.MusicDiscPlay{DiscType: disc}, ok
	case native.SoundDecoratedPotInserted:
		return dfsound.DecoratedPotInserted{Progress: value.Scalar}, true
	case native.SoundItemUseOn:
		return blockSound(tx, value.Block, func(value world.Block) world.Sound { return dfsound.ItemUseOn{Block: value} })
	case native.SoundEquipItem:
		stack, ok := ItemStackFromNative(*value.Item)
		if !ok || stack.Empty() {
			return nil, false
		}
		return dfsound.EquipItem{Item: stack.Item()}, true
	case native.SoundBucketFill:
		liquid, ok := soundLiquid(value.Data)
		return dfsound.BucketFill{Liquid: liquid}, ok
	case native.SoundBucketEmpty:
		liquid, ok := soundLiquid(value.Data)
		return dfsound.BucketEmpty{Liquid: liquid}, ok
	case native.SoundCrossbowLoad:
		return dfsound.CrossbowLoad{Stage: int(value.Integer), QuickCharge: value.Flags&soundFlagPrimary != 0}, true
	case native.SoundGoatHorn:
		horn, ok := soundHorn(value.Data)
		return dfsound.GoatHorn{Horn: horn}, ok
	default:
		return nil, false
	}
}

// SoundToNative encodes a Dragonfly sound into the native descriptor used by plugins.
func SoundToNative(tx *world.Tx, sound world.Sound) (native.WorldSound, bool) {
	if tx == nil || sound == nil {
		return native.WorldSound{}, false
	}
	switch value := sound.(type) {
	case dfsound.AnvilBreak:
		return native.WorldSound{Kind: native.SoundAnvilBreak}, true
	case dfsound.AnvilLand:
		return native.WorldSound{Kind: native.SoundAnvilLand}, true
	case dfsound.AnvilUse:
		return native.WorldSound{Kind: native.SoundAnvilUse}, true
	case dfsound.ArrowHit:
		return native.WorldSound{Kind: native.SoundArrowHit}, true
	case dfsound.BarrelClose:
		return native.WorldSound{Kind: native.SoundBarrelClose}, true
	case dfsound.BarrelOpen:
		return native.WorldSound{Kind: native.SoundBarrelOpen}, true
	case dfsound.BlastFurnaceCrackle:
		return native.WorldSound{Kind: native.SoundBlastFurnaceCrackle}, true
	case dfsound.BowShoot:
		return native.WorldSound{Kind: native.SoundBowShoot}, true
	case dfsound.Burning:
		return native.WorldSound{Kind: native.SoundBurning}, true
	case dfsound.Burp:
		return native.WorldSound{Kind: native.SoundBurp}, true
	case dfsound.CampfireCrackle:
		return native.WorldSound{Kind: native.SoundCampfireCrackle}, true
	case dfsound.ChestClose:
		return native.WorldSound{Kind: native.SoundChestClose}, true
	case dfsound.ChestOpen:
		return native.WorldSound{Kind: native.SoundChestOpen}, true
	case dfsound.Click:
		return native.WorldSound{Kind: native.SoundClick}, true
	case dfsound.ComposterEmpty:
		return native.WorldSound{Kind: native.SoundComposterEmpty}, true
	case dfsound.ComposterFill:
		return native.WorldSound{Kind: native.SoundComposterFill}, true
	case dfsound.ComposterFillLayer:
		return native.WorldSound{Kind: native.SoundComposterFillLayer}, true
	case dfsound.ComposterReady:
		return native.WorldSound{Kind: native.SoundComposterReady}, true
	case dfsound.CopperScraped:
		return native.WorldSound{Kind: native.SoundCopperScraped}, true
	case dfsound.CrossbowShoot:
		return native.WorldSound{Kind: native.SoundCrossbowShoot}, true
	case dfsound.DecoratedPotInsertFailed:
		return native.WorldSound{Kind: native.SoundDecoratedPotInsertFailed}, true
	case dfsound.Deny:
		return native.WorldSound{Kind: native.SoundDeny}, true
	case dfsound.DoorCrash:
		return native.WorldSound{Kind: native.SoundDoorCrash}, true
	case dfsound.Drowning:
		return native.WorldSound{Kind: native.SoundDrowning}, true
	case dfsound.EnderChestClose:
		return native.WorldSound{Kind: native.SoundEnderChestClose}, true
	case dfsound.EnderChestOpen:
		return native.WorldSound{Kind: native.SoundEnderChestOpen}, true
	case dfsound.Experience:
		return native.WorldSound{Kind: native.SoundExperience}, true
	case dfsound.Explosion:
		return native.WorldSound{Kind: native.SoundExplosion}, true
	case dfsound.FireCharge:
		return native.WorldSound{Kind: native.SoundFireCharge}, true
	case dfsound.FireExtinguish:
		return native.WorldSound{Kind: native.SoundFireExtinguish}, true
	case dfsound.FireworkBlast:
		return native.WorldSound{Kind: native.SoundFireworkBlast}, true
	case dfsound.FireworkHugeBlast:
		return native.WorldSound{Kind: native.SoundFireworkHugeBlast}, true
	case dfsound.FireworkLaunch:
		return native.WorldSound{Kind: native.SoundFireworkLaunch}, true
	case dfsound.FireworkTwinkle:
		return native.WorldSound{Kind: native.SoundFireworkTwinkle}, true
	case dfsound.Fizz:
		return native.WorldSound{Kind: native.SoundFizz}, true
	case dfsound.FurnaceCrackle:
		return native.WorldSound{Kind: native.SoundFurnaceCrackle}, true
	case dfsound.GhastShoot:
		return native.WorldSound{Kind: native.SoundGhastShoot}, true
	case dfsound.GhastWarning:
		return native.WorldSound{Kind: native.SoundGhastWarning}, true
	case dfsound.GlassBreak:
		return native.WorldSound{Kind: native.SoundGlassBreak}, true
	case dfsound.Ignite:
		return native.WorldSound{Kind: native.SoundIgnite}, true
	case dfsound.ItemAdd:
		return native.WorldSound{Kind: native.SoundItemAdd}, true
	case dfsound.ItemBreak:
		return native.WorldSound{Kind: native.SoundItemBreak}, true
	case dfsound.ItemFrameRemove:
		return native.WorldSound{Kind: native.SoundItemFrameRemove}, true
	case dfsound.ItemFrameRotate:
		return native.WorldSound{Kind: native.SoundItemFrameRotate}, true
	case dfsound.ItemThrow:
		return native.WorldSound{Kind: native.SoundItemThrow}, true
	case dfsound.LecternBookPlace:
		return native.WorldSound{Kind: native.SoundLecternBookPlace}, true
	case dfsound.LevelUp:
		return native.WorldSound{Kind: native.SoundLevelUp}, true
	case dfsound.LightningExplode:
		return native.WorldSound{Kind: native.SoundLightningExplode}, true
	case dfsound.LightningThunder:
		return native.WorldSound{Kind: native.SoundLightningThunder}, true
	case dfsound.MusicDiscEnd:
		return native.WorldSound{Kind: native.SoundMusicDiscEnd}, true
	case dfsound.Pop:
		return native.WorldSound{Kind: native.SoundPop}, true
	case dfsound.PotionBrewed:
		return native.WorldSound{Kind: native.SoundPotionBrewed}, true
	case dfsound.PowerOff:
		return native.WorldSound{Kind: native.SoundPowerOff}, true
	case dfsound.PowerOn:
		return native.WorldSound{Kind: native.SoundPowerOn}, true
	case dfsound.SignWaxed:
		return native.WorldSound{Kind: native.SoundSignWaxed}, true
	case dfsound.SmokerCrackle:
		return native.WorldSound{Kind: native.SoundSmokerCrackle}, true
	case dfsound.StopUsingSpyglass:
		return native.WorldSound{Kind: native.SoundStopUsingSpyglass}, true
	case dfsound.TNT:
		return native.WorldSound{Kind: native.SoundTnt}, true
	case dfsound.Teleport:
		return native.WorldSound{Kind: native.SoundTeleport}, true
	case dfsound.Thunder:
		return native.WorldSound{Kind: native.SoundThunder}, true
	case dfsound.Totem:
		return native.WorldSound{Kind: native.SoundTotem}, true
	case dfsound.UseSpyglass:
		return native.WorldSound{Kind: native.SoundUseSpyglass}, true
	case dfsound.WaxRemoved:
		return native.WorldSound{Kind: native.SoundWaxRemoved}, true
	case dfsound.WaxedSignFailedInteraction:
		return native.WorldSound{Kind: native.SoundWaxedSignFailedInteraction}, true
	case dfsound.ShulkerBoxOpen:
		return native.WorldSound{Kind: native.SoundShulkerBoxOpen}, true
	case dfsound.ShulkerBoxClose:
		return native.WorldSound{Kind: native.SoundShulkerBoxClose}, true
	case dfsound.EnderEyePlaced:
		return native.WorldSound{Kind: native.SoundEnderEyePlaced}, true
	case dfsound.EndPortalCreated:
		return native.WorldSound{Kind: native.SoundEndPortalCreated}, true
	case dfsound.Attack:
		var flags uint32
		if value.Damage {
			flags = soundFlagPrimary
		}
		return native.WorldSound{Kind: native.SoundAttack, Flags: flags}, true
	case dfsound.Fall:
		if math.IsNaN(value.Distance) || math.IsInf(value.Distance, 0) {
			return native.WorldSound{}, false
		}
		return native.WorldSound{Kind: native.SoundFall, Scalar: value.Distance}, true
	case dfsound.BlockPlace:
		return soundBlockToNative(tx, native.SoundBlockPlace, value.Block)
	case dfsound.BlockBreaking:
		return soundBlockToNative(tx, native.SoundBlockBreaking, value.Block)
	case dfsound.DoorOpen:
		return soundBlockToNative(tx, native.SoundDoorOpen, value.Block)
	case dfsound.DoorClose:
		return soundBlockToNative(tx, native.SoundDoorClose, value.Block)
	case dfsound.TrapdoorOpen:
		return soundBlockToNative(tx, native.SoundTrapdoorOpen, value.Block)
	case dfsound.TrapdoorClose:
		return soundBlockToNative(tx, native.SoundTrapdoorClose, value.Block)
	case dfsound.FenceGateOpen:
		return soundBlockToNative(tx, native.SoundFenceGateOpen, value.Block)
	case dfsound.FenceGateClose:
		return soundBlockToNative(tx, native.SoundFenceGateClose, value.Block)
	case dfsound.Note:
		instrument := value.Instrument.Int32()
		if instrument < 0 || instrument >= 16 || value.Pitch < math.MinInt32 || value.Pitch > math.MaxInt32 {
			return native.WorldSound{}, false
		}
		return native.WorldSound{Kind: native.SoundNote, Data: uint32(instrument), Integer: int32(value.Pitch)}, true
	case dfsound.MusicDiscPlay:
		disc := uint32(value.DiscType.Uint8())
		if disc >= 21 {
			return native.WorldSound{}, false
		}
		return native.WorldSound{Kind: native.SoundMusicDiscPlay, Data: disc}, true
	case dfsound.DecoratedPotInserted:
		if math.IsNaN(value.Progress) || math.IsInf(value.Progress, 0) {
			return native.WorldSound{}, false
		}
		return native.WorldSound{Kind: native.SoundDecoratedPotInserted, Scalar: value.Progress}, true
	case dfsound.ItemUseOn:
		return soundBlockToNative(tx, native.SoundItemUseOn, value.Block)
	case dfsound.EquipItem:
		if value.Item == nil {
			return native.WorldSound{}, false
		}
		stack, ok := itemStackToNative(dfitem.NewStack(value.Item, 1))
		if !ok || stack.Count != 1 {
			return native.WorldSound{}, false
		}
		return native.WorldSound{Kind: native.SoundEquipItem, Item: &stack}, true
	case dfsound.BucketFill:
		liquid, ok := soundLiquidToNative(value.Liquid)
		if !ok {
			return native.WorldSound{}, false
		}
		encoded, ok := soundBlockToNative(tx, native.SoundBucketFill, value.Liquid)
		encoded.Data = liquid
		return encoded, ok
	case dfsound.BucketEmpty:
		liquid, ok := soundLiquidToNative(value.Liquid)
		if !ok {
			return native.WorldSound{}, false
		}
		encoded, ok := soundBlockToNative(tx, native.SoundBucketEmpty, value.Liquid)
		encoded.Data = liquid
		return encoded, ok
	case dfsound.CrossbowLoad:
		if value.Stage < dfsound.CrossbowLoadingStart || value.Stage > dfsound.CrossbowLoadingEnd {
			return native.WorldSound{}, false
		}
		var flags uint32
		if value.QuickCharge {
			flags = soundFlagPrimary
		}
		return native.WorldSound{Kind: native.SoundCrossbowLoad, Integer: int32(value.Stage), Flags: flags}, true
	case dfsound.GoatHorn:
		horn := uint32(value.Horn.Uint8())
		if horn >= 8 {
			return native.WorldSound{}, false
		}
		return native.WorldSound{Kind: native.SoundGoatHorn, Data: horn}, true
	case *dfsound.AnvilBreak:
		return encodeSoundPointer(tx, value)
	case *dfsound.AnvilLand:
		return encodeSoundPointer(tx, value)
	case *dfsound.AnvilUse:
		return encodeSoundPointer(tx, value)
	case *dfsound.ArrowHit:
		return encodeSoundPointer(tx, value)
	case *dfsound.BarrelClose:
		return encodeSoundPointer(tx, value)
	case *dfsound.BarrelOpen:
		return encodeSoundPointer(tx, value)
	case *dfsound.BlastFurnaceCrackle:
		return encodeSoundPointer(tx, value)
	case *dfsound.BowShoot:
		return encodeSoundPointer(tx, value)
	case *dfsound.Burning:
		return encodeSoundPointer(tx, value)
	case *dfsound.Burp:
		return encodeSoundPointer(tx, value)
	case *dfsound.CampfireCrackle:
		return encodeSoundPointer(tx, value)
	case *dfsound.ChestClose:
		return encodeSoundPointer(tx, value)
	case *dfsound.ChestOpen:
		return encodeSoundPointer(tx, value)
	case *dfsound.Click:
		return encodeSoundPointer(tx, value)
	case *dfsound.ComposterEmpty:
		return encodeSoundPointer(tx, value)
	case *dfsound.ComposterFill:
		return encodeSoundPointer(tx, value)
	case *dfsound.ComposterFillLayer:
		return encodeSoundPointer(tx, value)
	case *dfsound.ComposterReady:
		return encodeSoundPointer(tx, value)
	case *dfsound.CopperScraped:
		return encodeSoundPointer(tx, value)
	case *dfsound.CrossbowShoot:
		return encodeSoundPointer(tx, value)
	case *dfsound.DecoratedPotInsertFailed:
		return encodeSoundPointer(tx, value)
	case *dfsound.Deny:
		return encodeSoundPointer(tx, value)
	case *dfsound.DoorCrash:
		return encodeSoundPointer(tx, value)
	case *dfsound.Drowning:
		return encodeSoundPointer(tx, value)
	case *dfsound.EnderChestClose:
		return encodeSoundPointer(tx, value)
	case *dfsound.EnderChestOpen:
		return encodeSoundPointer(tx, value)
	case *dfsound.Experience:
		return encodeSoundPointer(tx, value)
	case *dfsound.Explosion:
		return encodeSoundPointer(tx, value)
	case *dfsound.FireCharge:
		return encodeSoundPointer(tx, value)
	case *dfsound.FireExtinguish:
		return encodeSoundPointer(tx, value)
	case *dfsound.FireworkBlast:
		return encodeSoundPointer(tx, value)
	case *dfsound.FireworkHugeBlast:
		return encodeSoundPointer(tx, value)
	case *dfsound.FireworkLaunch:
		return encodeSoundPointer(tx, value)
	case *dfsound.FireworkTwinkle:
		return encodeSoundPointer(tx, value)
	case *dfsound.Fizz:
		return encodeSoundPointer(tx, value)
	case *dfsound.FurnaceCrackle:
		return encodeSoundPointer(tx, value)
	case *dfsound.GhastShoot:
		return encodeSoundPointer(tx, value)
	case *dfsound.GhastWarning:
		return encodeSoundPointer(tx, value)
	case *dfsound.GlassBreak:
		return encodeSoundPointer(tx, value)
	case *dfsound.Ignite:
		return encodeSoundPointer(tx, value)
	case *dfsound.ItemAdd:
		return encodeSoundPointer(tx, value)
	case *dfsound.ItemBreak:
		return encodeSoundPointer(tx, value)
	case *dfsound.ItemFrameRemove:
		return encodeSoundPointer(tx, value)
	case *dfsound.ItemFrameRotate:
		return encodeSoundPointer(tx, value)
	case *dfsound.ItemThrow:
		return encodeSoundPointer(tx, value)
	case *dfsound.LecternBookPlace:
		return encodeSoundPointer(tx, value)
	case *dfsound.LevelUp:
		return encodeSoundPointer(tx, value)
	case *dfsound.LightningExplode:
		return encodeSoundPointer(tx, value)
	case *dfsound.LightningThunder:
		return encodeSoundPointer(tx, value)
	case *dfsound.MusicDiscEnd:
		return encodeSoundPointer(tx, value)
	case *dfsound.Pop:
		return encodeSoundPointer(tx, value)
	case *dfsound.PotionBrewed:
		return encodeSoundPointer(tx, value)
	case *dfsound.PowerOff:
		return encodeSoundPointer(tx, value)
	case *dfsound.PowerOn:
		return encodeSoundPointer(tx, value)
	case *dfsound.SignWaxed:
		return encodeSoundPointer(tx, value)
	case *dfsound.SmokerCrackle:
		return encodeSoundPointer(tx, value)
	case *dfsound.StopUsingSpyglass:
		return encodeSoundPointer(tx, value)
	case *dfsound.TNT:
		return encodeSoundPointer(tx, value)
	case *dfsound.Teleport:
		return encodeSoundPointer(tx, value)
	case *dfsound.Thunder:
		return encodeSoundPointer(tx, value)
	case *dfsound.Totem:
		return encodeSoundPointer(tx, value)
	case *dfsound.UseSpyglass:
		return encodeSoundPointer(tx, value)
	case *dfsound.WaxRemoved:
		return encodeSoundPointer(tx, value)
	case *dfsound.WaxedSignFailedInteraction:
		return encodeSoundPointer(tx, value)
	case *dfsound.ShulkerBoxOpen:
		return encodeSoundPointer(tx, value)
	case *dfsound.ShulkerBoxClose:
		return encodeSoundPointer(tx, value)
	case *dfsound.EnderEyePlaced:
		return encodeSoundPointer(tx, value)
	case *dfsound.EndPortalCreated:
		return encodeSoundPointer(tx, value)
	case *dfsound.Attack:
		return encodeSoundPointer(tx, value)
	case *dfsound.Fall:
		return encodeSoundPointer(tx, value)
	case *dfsound.BlockPlace:
		return encodeSoundPointer(tx, value)
	case *dfsound.BlockBreaking:
		return encodeSoundPointer(tx, value)
	case *dfsound.DoorOpen:
		return encodeSoundPointer(tx, value)
	case *dfsound.DoorClose:
		return encodeSoundPointer(tx, value)
	case *dfsound.TrapdoorOpen:
		return encodeSoundPointer(tx, value)
	case *dfsound.TrapdoorClose:
		return encodeSoundPointer(tx, value)
	case *dfsound.FenceGateOpen:
		return encodeSoundPointer(tx, value)
	case *dfsound.FenceGateClose:
		return encodeSoundPointer(tx, value)
	case *dfsound.Note:
		return encodeSoundPointer(tx, value)
	case *dfsound.MusicDiscPlay:
		return encodeSoundPointer(tx, value)
	case *dfsound.DecoratedPotInserted:
		return encodeSoundPointer(tx, value)
	case *dfsound.ItemUseOn:
		return encodeSoundPointer(tx, value)
	case *dfsound.EquipItem:
		return encodeSoundPointer(tx, value)
	case *dfsound.BucketFill:
		return encodeSoundPointer(tx, value)
	case *dfsound.BucketEmpty:
		return encodeSoundPointer(tx, value)
	case *dfsound.CrossbowLoad:
		return encodeSoundPointer(tx, value)
	case *dfsound.GoatHorn:
		return encodeSoundPointer(tx, value)
	default:
		return native.WorldSound{}, false
	}
}

func encodeSoundPointer[T world.Sound](tx *world.Tx, value *T) (native.WorldSound, bool) {
	if value == nil {
		return native.WorldSound{}, false
	}
	return SoundToNative(tx, *value)
}

func soundBlockToNative(tx *world.Tx, kind native.SoundKind, value world.Block) (native.WorldSound, bool) {
	if value == nil {
		return native.WorldSound{}, false
	}
	identifier, properties := value.EncodeBlock()
	if identifier == "" {
		return native.WorldSound{}, false
	}
	if _, ok := tx.World().BlockRegistry().BlockByName(identifier, properties); !ok {
		return native.WorldSound{}, false
	}
	encoded, ok := EncodeBlockProperties(properties)
	if !ok {
		return native.WorldSound{}, false
	}
	return native.WorldSound{
		Kind:  kind,
		Block: &native.WorldBlock{Identifier: identifier, PropertiesNBT: encoded},
	}, true
}

func soundLiquidToNative(value world.Liquid) (uint32, bool) {
	switch value.(type) {
	case block.Water:
		return 0, true
	case block.Lava:
		return 1, true
	default:
		return 0, false
	}
}

func blockSound(tx *world.Tx, value *native.WorldBlock, build func(world.Block) world.Sound) (world.Sound, bool) {
	properties, ok := DecodeBlockProperties(value.PropertiesNBT)
	if !ok {
		return nil, false
	}
	resolved, ok := tx.World().BlockRegistry().BlockByName(value.Identifier, properties)
	if !ok {
		return nil, false
	}
	return build(resolved), true
}

func soundInstrument(value uint32) (dfsound.Instrument, bool) {
	values := [...]dfsound.Instrument{dfsound.Piano(), dfsound.BassDrum(), dfsound.Snare(), dfsound.ClicksAndSticks(), dfsound.Bass(), dfsound.Flute(), dfsound.Bell(), dfsound.Guitar(), dfsound.Chimes(), dfsound.Xylophone(), dfsound.IronXylophone(), dfsound.CowBell(), dfsound.Didgeridoo(), dfsound.Bit(), dfsound.Banjo(), dfsound.Pling()}
	if value >= uint32(len(values)) {
		return dfsound.Instrument{}, false
	}
	return values[value], true
}

func soundDisc(value uint32) (dfsound.DiscType, bool) {
	values := [...]dfsound.DiscType{dfsound.Disc13(), dfsound.DiscCat(), dfsound.DiscBlocks(), dfsound.DiscChirp(), dfsound.DiscFar(), dfsound.DiscMall(), dfsound.DiscMellohi(), dfsound.DiscStal(), dfsound.DiscStrad(), dfsound.DiscWard(), dfsound.Disc11(), dfsound.DiscWait(), dfsound.DiscOtherside(), dfsound.DiscPigstep(), dfsound.Disc5(), dfsound.DiscRelic(), dfsound.DiscCreator(), dfsound.DiscCreatorMusicBox(), dfsound.DiscPrecipice(), dfsound.DiscTears(), dfsound.DiscLavaChicken()}
	if value >= uint32(len(values)) {
		return dfsound.DiscType{}, false
	}
	return values[value], true
}

func soundHorn(value uint32) (dfsound.Horn, bool) {
	values := [...]dfsound.Horn{dfsound.Ponder(), dfsound.Sing(), dfsound.Seek(), dfsound.Feel(), dfsound.Admire(), dfsound.Call(), dfsound.Yearn(), dfsound.Dream()}
	if value >= uint32(len(values)) {
		return dfsound.Horn{}, false
	}
	return values[value], true
}

func soundLiquid(value uint32) (world.Liquid, bool) {
	switch value {
	case 0:
		return block.Water{Depth: 8, Still: true}, true
	case 1:
		return block.Lava{Depth: 8, Still: true}, true
	default:
		return nil, false
	}
}

func (p *Players) PlayPlayerSound(invocation native.InvocationID, id native.PlayerID, value native.WorldSound) bool {
	if !ValidSound(value) {
		return false
	}
	play := func(tx *world.Tx, connected *player.Player) bool {
		decoded, ok := SoundFromNative(tx, value)
		if !ok {
			return false
		}
		connected.PlaySound(decoded)
		return true
	}
	if connected, ok := p.ResolveID(id, invocation); ok {
		tx, _ := p.InvocationTx(invocation)
		return play(tx, connected)
	}
	if invocation != 0 {
		if _, ok := p.InvocationTx(invocation); !ok {
			return false
		}
	}
	entry, ok := p.playerEntry(id)
	if !ok {
		return false
	}
	task := world.NewEntityRef[*player.Player](entry.handle).Do(func(tx *world.Tx, connected *player.Player) {
		play(tx, connected)
	})
	// Off-callback work is fire-and-forget. Success means the validated sound
	// was accepted for scheduling.
	return task.Err() == nil
}
