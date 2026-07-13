package framework

import (
	"errors"
	"fmt"
	"math"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/world"
)

// CoreWorldConfig controls creation-time policies shared by Dragonfly's
// overworld, Nether, and End worlds. Their provider remains configured through
// dragonfly.World.
type CoreWorldConfig struct {
	ReadOnly       bool                  `toml:"read-only"`
	RandomTicks    WorldRandomTickPolicy `toml:"random-ticks"`
	RandomTickRate uint32                `toml:"random-tick-rate"`
	Time           WorldTimePolicy       `toml:"time"`
	FixedTime      int64                 `toml:"fixed-time"`
	Weather        WorldWeatherPolicy    `toml:"weather"`
}

func defaultCoreWorldConfig() CoreWorldConfig {
	return CoreWorldConfig{
		RandomTicks:    WorldRandomTicksPerSubchunk,
		RandomTickRate: 3,
		Time:           WorldTimePreserve,
		Weather:        WorldWeatherPreserve,
	}
}

func (config CoreWorldConfig) validate() error {
	if config.RandomTicks > WorldRandomTicksPerSubchunk {
		return fmt.Errorf("invalid random tick policy %d", config.RandomTicks)
	}
	switch config.RandomTicks {
	case WorldRandomTicksDisabled:
		if config.RandomTickRate != 0 {
			return errors.New("disabled random ticks require a zero rate")
		}
	case WorldRandomTicksPerSubchunk:
		if config.RandomTickRate == 0 || config.RandomTickRate > math.MaxInt32 {
			return fmt.Errorf("random tick rate must be between 1 and %d", math.MaxInt32)
		}
	}
	if config.Time > WorldTimeFixed {
		return fmt.Errorf("invalid time policy %d", config.Time)
	}
	if config.Time != WorldTimeFixed && config.FixedTime != 0 {
		return errors.New("non-fixed time policy requires a zero fixed time")
	}
	if config.Weather > WorldWeatherClear {
		return fmt.Errorf("invalid weather policy %d", config.Weather)
	}
	return nil
}

// applyCoreWorldPolicy mutates config before server.Config.New creates the
// three core worlds. Dragonfly gives all three worlds the same provider and
// Settings pointer.
func applyCoreWorldPolicy(config *server.Config, policy CoreWorldConfig) error {
	if config == nil {
		return errors.New("Dragonfly config is nil")
	}
	if err := policy.validate(); err != nil {
		return err
	}
	if config.WorldProvider == nil {
		settings := world.NopProvider{}.Settings()
		config.WorldProvider = world.NopProvider{Set: settings}
	}
	settings := config.WorldProvider.Settings()
	if settings == nil {
		return errors.New("world provider returned nil settings")
	}

	config.ReadOnlyWorld = policy.ReadOnly
	config.RandomTickSpeed = -1
	if policy.RandomTicks == WorldRandomTicksPerSubchunk {
		config.RandomTickSpeed = int(policy.RandomTickRate)
	}

	settings.Lock()
	defer settings.Unlock()
	switch policy.Time {
	case WorldTimeCycle:
		settings.TimeCycle = true
	case WorldTimeFixed:
		settings.Time = policy.FixedTime
		settings.TimeCycle = false
	}
	switch policy.Weather {
	case WorldWeatherCycle:
		settings.WeatherCycle = true
	case WorldWeatherClear:
		settings.RainTime = 0
		settings.Raining = false
		settings.ThunderTime = 0
		settings.Thundering = false
		settings.WeatherCycle = false
	}
	return nil
}

func (policy WorldRandomTickPolicy) MarshalText() ([]byte, error) {
	switch policy {
	case WorldRandomTicksDisabled:
		return []byte("disabled"), nil
	case WorldRandomTicksPerSubchunk:
		return []byte("per-subchunk"), nil
	default:
		return nil, fmt.Errorf("invalid random tick policy %d", policy)
	}
}

func (policy *WorldRandomTickPolicy) UnmarshalText(text []byte) error {
	switch string(text) {
	case "disabled":
		*policy = WorldRandomTicksDisabled
	case "per-subchunk":
		*policy = WorldRandomTicksPerSubchunk
	default:
		return fmt.Errorf("invalid random tick policy %q", text)
	}
	return nil
}

func (policy WorldTimePolicy) MarshalText() ([]byte, error) {
	switch policy {
	case WorldTimePreserve:
		return []byte("preserve"), nil
	case WorldTimeCycle:
		return []byte("cycle"), nil
	case WorldTimeFixed:
		return []byte("fixed"), nil
	default:
		return nil, fmt.Errorf("invalid time policy %d", policy)
	}
}

func (policy *WorldTimePolicy) UnmarshalText(text []byte) error {
	switch string(text) {
	case "preserve":
		*policy = WorldTimePreserve
	case "cycle":
		*policy = WorldTimeCycle
	case "fixed":
		*policy = WorldTimeFixed
	default:
		return fmt.Errorf("invalid time policy %q", text)
	}
	return nil
}

func (policy WorldWeatherPolicy) MarshalText() ([]byte, error) {
	switch policy {
	case WorldWeatherPreserve:
		return []byte("preserve"), nil
	case WorldWeatherCycle:
		return []byte("cycle"), nil
	case WorldWeatherClear:
		return []byte("clear"), nil
	default:
		return nil, fmt.Errorf("invalid weather policy %d", policy)
	}
}

func (policy *WorldWeatherPolicy) UnmarshalText(text []byte) error {
	switch string(text) {
	case "preserve":
		*policy = WorldWeatherPreserve
	case "cycle":
		*policy = WorldWeatherCycle
	case "clear":
		*policy = WorldWeatherClear
	default:
		return fmt.Errorf("invalid weather policy %q", text)
	}
	return nil
}
