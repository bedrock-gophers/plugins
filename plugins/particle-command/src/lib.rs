use dragonfly::{
    BlockFace, BlockPos, Context, Player, Plugin, Vec3, World, block, particle, plugin, sound,
};

#[derive(Default)]
struct ParticleCommand;

#[plugin]
#[command("particle")]
impl Plugin for ParticleCommand {
    #[command]
    fn root(&self, context: &mut Context<'_, Player>) {
        context.source().message(
            "Use /particle flame, coloured-flame, dust, block-break, punch-block, bone-meal, note, dragon-egg, explosion, splash, or effect.",
        );
    }

    #[subcommand("flame")]
    fn flame(&self, context: &mut Context<'_, Player>) {
        self.add(context, particle::Flame::new());
    }

    #[subcommand("coloured-flame")]
    fn coloured_flame(&self, context: &mut Context<'_, Player>) {
        self.add(
            context,
            particle::Flame::coloured(particle::Colour::rgb(80, 180, 255)),
        );
    }

    #[subcommand("dust")]
    fn dust(&self, context: &mut Context<'_, Player>) {
        self.add(
            context,
            particle::Dust::new(particle::Colour::rgb(255, 90, 40)),
        );
    }

    #[subcommand("block-break")]
    fn block_break(&self, context: &mut Context<'_, Player>) {
        self.add(context, particle::BlockBreak::new(block::DiamondBlock));
    }

    #[subcommand("punch-block")]
    fn punch_block(&self, context: &mut Context<'_, Player>) {
        self.add(
            context,
            particle::PunchBlock::new(block::GrassBlock, BlockFace::Up),
        );
    }

    #[subcommand("bone-meal")]
    fn bone_meal(&self, context: &mut Context<'_, Player>) {
        self.add(context, particle::BoneMeal::new(true));
    }

    #[subcommand("note")]
    fn note(&self, context: &mut Context<'_, Player>) {
        self.add(context, particle::Note::new(sound::Instrument::Bell, 12));
    }

    #[subcommand("dragon-egg")]
    fn dragon_egg(&self, context: &mut Context<'_, Player>) {
        self.add(
            context,
            particle::DragonEggTeleport::new(BlockPos { x: 3, y: 1, z: 0 }),
        );
    }

    #[subcommand("explosion")]
    fn explosion(&self, context: &mut Context<'_, Player>) {
        self.add(context, particle::HugeExplosion);
    }

    #[subcommand("splash")]
    fn splash(&self, context: &mut Context<'_, Player>) {
        self.add(
            context,
            particle::Splash::coloured(particle::Colour::rgb(125, 70, 255)),
        );
    }

    #[subcommand("effect")]
    fn effect(&self, context: &mut Context<'_, Player>) {
        self.add(
            context,
            particle::Effect::new(particle::Colour::rgb(80, 255, 120)),
        );
    }
}

impl ParticleCommand {
    fn add(&self, context: &mut Context<'_, Player>, value: impl particle::Particle) {
        let player = context.source();
        let Some((world, position)) = world_and_position(player) else {
            player.message("World or position unavailable.");
            return;
        };
        world.add_particle(above(position), value);
    }
}

fn world_and_position(player: Player) -> Option<(World, Vec3)> {
    let entity = player.entity();
    Some((entity.world()?, entity.position()?))
}

fn above(mut position: Vec3) -> Vec3 {
    position.y += 1.5;
    position
}
