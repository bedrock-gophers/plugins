package host

import (
	"math"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/block"
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
