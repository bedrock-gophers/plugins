package main

import "testing"

func TestDotnetRID(t *testing.T) {
	tests := []struct {
		goos   string
		goarch string
		want   string
	}{
		{goos: "windows", goarch: "amd64", want: "win-x64"},
		{goos: "windows", goarch: "arm64", want: "win-arm64"},
		{goos: "linux", goarch: "amd64", want: "linux-x64"},
		{goos: "darwin", goarch: "arm64", want: "osx-arm64"},
	}
	for _, test := range tests {
		got, err := dotnetRID(test.goos, test.goarch)
		if err != nil {
			t.Fatalf("dotnetRID(%q, %q): %v", test.goos, test.goarch, err)
		}
		if got != test.want {
			t.Fatalf("dotnetRID(%q, %q) = %q, want %q", test.goos, test.goarch, got, test.want)
		}
	}
}

func TestTargetPlatformNames(t *testing.T) {
	tests := []struct {
		rid       string
		extension string
		runtime   string
	}{
		{rid: "win-x64", extension: ".dll", runtime: "dragonfly_plugin_runtime.dll"},
		{rid: "linux-arm64", extension: ".so", runtime: "libdragonfly_plugin_runtime.so"},
		{rid: "osx-x64", extension: ".dylib", runtime: "libdragonfly_plugin_runtime.dylib"},
	}
	for _, test := range tests {
		goos, err := targetOS(test.rid)
		if err != nil {
			t.Fatalf("targetOS(%q): %v", test.rid, err)
		}
		if got := nativeExtension(goos); got != test.extension {
			t.Fatalf("nativeExtension(%q) = %q, want %q", goos, got, test.extension)
		}
		if got := runtimeLibraryName(goos); got != test.runtime {
			t.Fatalf("runtimeLibraryName(%q) = %q, want %q", goos, got, test.runtime)
		}
	}
}
