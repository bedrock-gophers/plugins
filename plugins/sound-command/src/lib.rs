use dragonfly::{Context, Player, Plugin, Vec3, World, block, item, plugin, sound};

#[derive(Default)]
struct SoundCommand;

#[plugin]
#[command("sound")]
impl Plugin for SoundCommand {
    #[command]
    fn root(&self, context: &mut Context<'_, Player>) {
        context.source().message(
            "Use /sound player, explosion, door, note, equip, bucket, disc, horn, attack, or crossbow.",
        );
    }

    #[subcommand("player")]
    fn player(&self, context: &mut Context<'_, Player>) {
        context.source().play_sound(sound::LevelUp);
    }

    #[subcommand("explosion")]
    fn explosion(&self, context: &mut Context<'_, Player>) {
        self.play(context, sound::Explosion);
    }

    #[subcommand("door")]
    fn door(&self, context: &mut Context<'_, Player>) {
        self.play(
            context,
            sound::DoorOpen::new(
                block::WoodenDoor::new().with_cardinal_direction(block::CardinalDirection::North),
            ),
        );
    }

    #[subcommand("note")]
    fn note(&self, context: &mut Context<'_, Player>) {
        self.play(context, sound::Note::new(sound::Instrument::Bell, 12));
    }

    #[subcommand("equip")]
    fn equip(&self, context: &mut Context<'_, Player>) {
        self.play(
            context,
            sound::EquipItem::new(item::Sword::new(item::ToolTier::Diamond)),
        );
    }

    #[subcommand("bucket")]
    fn bucket(&self, context: &mut Context<'_, Player>) {
        self.play(context, sound::BucketEmpty::new(sound::Liquid::Water));
    }

    #[subcommand("disc")]
    fn disc(&self, context: &mut Context<'_, Player>) {
        self.play(
            context,
            sound::MusicDiscPlay::new(sound::DiscType::DiscPigstep),
        );
    }

    #[subcommand("horn")]
    fn horn(&self, context: &mut Context<'_, Player>) {
        self.play(context, sound::GoatHorn::new(sound::Horn::Dream));
    }

    #[subcommand("attack")]
    fn attack(&self, context: &mut Context<'_, Player>) {
        self.play(context, sound::Attack::new(true));
    }

    #[subcommand("crossbow")]
    fn crossbow(&self, context: &mut Context<'_, Player>) {
        self.play(
            context,
            sound::CrossbowLoad::new(sound::CrossbowStage::End, true),
        );
    }
}

impl SoundCommand {
    fn play(&self, context: &mut Context<'_, Player>, value: impl sound::Sound) {
        let player = context.source();
        let Some((world, position)) = world_and_position(player) else {
            player.message("World or position unavailable.");
            return;
        };
        world.play_sound(position, value);
    }
}

fn world_and_position(player: Player) -> Option<(World, Vec3)> {
    let entity = player.entity();
    Some((entity.world()?, entity.position()?))
}
