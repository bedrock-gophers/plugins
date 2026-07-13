package framework

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/pelletier/go-toml"
)

func TestCoreWorldConfigTOMLRoundTrip(t *testing.T) {
	config := DefaultConfig()
	config.Worlds.Core = CoreWorldConfig{
		ReadOnly:       true,
		RandomTicks:    WorldRandomTicksDisabled,
		RandomTickRate: 0,
		Time:           WorldTimeFixed,
		FixedTime:      6000,
		Weather:        WorldWeatherClear,
	}
	data, err := toml.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range []string{
		`random-ticks = "disabled"`,
		`time = "fixed"`,
		`weather = "clear"`,
	} {
		if !strings.Contains(string(data), value) {
			t.Fatalf("TOML missing %q:\n%s", value, data)
		}
	}

	loaded := DefaultConfig()
	if err := toml.Unmarshal(data, &loaded); err != nil {
		t.Fatal(err)
	}
	if loaded.Worlds.Core != config.Worlds.Core {
		t.Fatalf("core policy = %+v, want %+v", loaded.Worlds.Core, config.Worlds.Core)
	}
}

func TestCoreWorldConfigRejectsUnknownTOMLEnums(t *testing.T) {
	tests := []string{
		"random-ticks = \"sometimes\"",
		"time = \"noonish\"",
		"weather = \"sunny-ish\"",
	}
	for _, field := range tests {
		t.Run(field, func(t *testing.T) {
			config := DefaultConfig()
			data := []byte("[worlds.core]\n" + field + "\n")
			if err := toml.Unmarshal(data, &config); err == nil {
				t.Fatal("unknown policy accepted")
			}
		})
	}
}

func TestLoadConfigRejectsInvalidCoreWorldCombination(t *testing.T) {
	path := filepath.Join(t.TempDir(), "server.toml")
	data := []byte(`[plugins]
directory = "plugins"

[worlds]
directory = "worlds"

[worlds.core]
random-ticks = "disabled"
random-tick-rate = 3
`)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadConfig(path); err == nil || !strings.Contains(err.Error(), "worlds.core") {
		t.Fatalf("LoadConfig error = %v", err)
	}
}

func TestLoadConfigRetainsDefaultCorePolicyWhenTableIsAbsent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "server.toml")
	data := []byte(`[plugins]
directory = "plugins"

[worlds]
directory = "worlds"
`)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	config, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if config.Worlds.Core != defaultCoreWorldConfig() {
		t.Fatalf("core policy = %+v, want %+v", config.Worlds.Core, defaultCoreWorldConfig())
	}
}

func TestCoreWorldConfigDefaultsMatchDragonfly(t *testing.T) {
	policy := DefaultConfig().Worlds.Core
	if policy.ReadOnly || policy.RandomTicks != WorldRandomTicksPerSubchunk ||
		policy.RandomTickRate != 3 || policy.Time != WorldTimePreserve ||
		policy.FixedTime != 0 || policy.Weather != WorldWeatherPreserve {
		t.Fatalf("default core policy = %+v", policy)
	}
}

func TestCoreWorldConfigRejectsInvalidCombinations(t *testing.T) {
	tests := map[string]func(*CoreWorldConfig){
		"unknown random ticks": func(policy *CoreWorldConfig) { policy.RandomTicks = 99 },
		"disabled ticks with rate": func(policy *CoreWorldConfig) {
			policy.RandomTicks = WorldRandomTicksDisabled
		},
		"ticks with zero rate":  func(policy *CoreWorldConfig) { policy.RandomTickRate = 0 },
		"ticks above max int32": func(policy *CoreWorldConfig) { policy.RandomTickRate = math.MaxInt32 + 1 },
		"unknown time":          func(policy *CoreWorldConfig) { policy.Time = 99 },
		"preserved fixed time":  func(policy *CoreWorldConfig) { policy.FixedTime = 6000 },
		"unknown weather":       func(policy *CoreWorldConfig) { policy.Weather = 99 },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			policy := defaultCoreWorldConfig()
			mutate(&policy)
			if err := policy.validate(); err == nil {
				t.Fatal("invalid policy accepted")
			}
		})
	}
}

func TestApplyCoreWorldPolicyMapsDragonflyConfigAndSettings(t *testing.T) {
	settings := &world.Settings{
		Time: 42, TimeCycle: true,
		RainTime: 20, Raining: true,
		ThunderTime: 30, Thundering: true,
		WeatherCycle: true,
	}
	provider := &recordingCoreProvider{NopProvider: world.NopProvider{Set: settings}}
	config := server.Config{WorldProvider: provider, RandomTickSpeed: 7}
	policy := CoreWorldConfig{
		ReadOnly: true, RandomTicks: WorldRandomTicksDisabled,
		Time: WorldTimeFixed, FixedTime: 6000, Weather: WorldWeatherClear,
	}
	if err := applyCoreWorldPolicy(&config, policy); err != nil {
		t.Fatal(err)
	}
	if !config.ReadOnlyWorld || config.RandomTickSpeed != -1 {
		t.Fatalf("Dragonfly config = %#v", config)
	}
	if settings.Time != 6000 || settings.TimeCycle || settings.RainTime != 0 ||
		settings.Raining || settings.ThunderTime != 0 || settings.Thundering ||
		settings.WeatherCycle {
		t.Fatalf("settings = %#v", settings)
	}
}

func TestApplyCoreWorldPolicyPreservesRequestedSettings(t *testing.T) {
	settings := &world.Settings{
		Time: 42, TimeCycle: false,
		RainTime: 20, Raining: true,
		ThunderTime: 30, Thundering: true,
		WeatherCycle: false,
	}
	config := server.Config{WorldProvider: world.NopProvider{Set: settings}}
	if err := applyCoreWorldPolicy(&config, defaultCoreWorldConfig()); err != nil {
		t.Fatal(err)
	}
	if config.ReadOnlyWorld || config.RandomTickSpeed != 3 {
		t.Fatalf("Dragonfly config = %#v", config)
	}
	if settings.Time != 42 || settings.TimeCycle || settings.RainTime != 20 ||
		!settings.Raining || settings.ThunderTime != 30 || !settings.Thundering ||
		settings.WeatherCycle {
		t.Fatalf("preserved settings changed: %#v", settings)
	}
}

func TestApplyCoreWorldPolicyInstallsConfiguredNopProvider(t *testing.T) {
	config := server.Config{}
	policy := defaultCoreWorldConfig()
	policy.Time, policy.FixedTime = WorldTimeFixed, 6000
	policy.Weather = WorldWeatherClear
	if err := applyCoreWorldPolicy(&config, policy); err != nil {
		t.Fatal(err)
	}
	if config.WorldProvider == nil {
		t.Fatal("world provider remains nil")
	}
	settings := config.WorldProvider.Settings()
	if settings.Name != "World" || settings.DefaultGameMode == nil ||
		settings.Difficulty == nil || settings.TickRange != 6 {
		t.Fatalf("Dragonfly defaults lost: %#v", settings)
	}
	if settings.Time != 6000 || settings.TimeCycle || settings.WeatherCycle {
		t.Fatalf("policy not applied: %#v", settings)
	}
}

func TestApplyCoreWorldPolicyRejectsBeforeMutation(t *testing.T) {
	settings := &world.Settings{Time: 42, TimeCycle: true}
	provider := world.NopProvider{Set: settings}
	config := server.Config{WorldProvider: provider, RandomTickSpeed: 7}
	policy := defaultCoreWorldConfig()
	policy.RandomTickRate = 0
	if err := applyCoreWorldPolicy(&config, policy); err == nil {
		t.Fatal("invalid policy accepted")
	}
	if config.RandomTickSpeed != 7 || config.WorldProvider != provider ||
		settings.Time != 42 || !settings.TimeCycle {
		t.Fatal("invalid policy partially mutated Dragonfly config")
	}
}

func TestApplyCoreWorldPolicyRejectsNilProviderSettings(t *testing.T) {
	config := server.Config{WorldProvider: nilSettingsCoreProvider{}}
	if err := applyCoreWorldPolicy(&config, defaultCoreWorldConfig()); err == nil {
		t.Fatal("nil provider settings accepted")
	}
	if config.RandomTickSpeed != 0 || config.ReadOnlyWorld {
		t.Fatal("invalid provider partially mutated Dragonfly config")
	}
}

func TestCoreWorldPolicyAppliesToAllDragonflyDimensionsAndClosesProviderOnce(t *testing.T) {
	settings := world.NopProvider{}.Settings()
	provider := &recordingCoreProvider{NopProvider: world.NopProvider{Set: settings}}
	config := server.Config{WorldProvider: provider}
	policy := CoreWorldConfig{
		ReadOnly: true, RandomTicks: WorldRandomTicksDisabled,
		Time: WorldTimeFixed, FixedTime: 6000, Weather: WorldWeatherClear,
	}
	if err := applyCoreWorldPolicy(&config, policy); err != nil {
		t.Fatal(err)
	}
	srv := config.New()
	for name, current := range map[string]*world.World{
		"overworld": srv.World(), "nether": srv.Nether(), "end": srv.End(),
	} {
		if current.Time() != 6000 || current.TimeCycle() {
			t.Fatalf("%s time = %d, cycle = %t", name, current.Time(), current.TimeCycle())
		}
	}
	if err := srv.End().Close(); err != nil {
		t.Fatal(err)
	}
	if err := srv.Nether().Close(); err != nil {
		t.Fatal(err)
	}
	if provider.closes.Load() != 0 {
		t.Fatalf("provider closed before overworld: %d", provider.closes.Load())
	}
	if err := srv.World().Close(); err != nil {
		t.Fatal(err)
	}
	if provider.closes.Load() != 1 {
		t.Fatalf("provider close count = %d, want 1", provider.closes.Load())
	}
}

type recordingCoreProvider struct {
	world.NopProvider
	closes atomic.Int32
}

type nilSettingsCoreProvider struct{ world.NopProvider }

func (nilSettingsCoreProvider) Settings() *world.Settings { return nil }

func (provider *recordingCoreProvider) Close() error {
	provider.closes.Add(1)
	return nil
}
