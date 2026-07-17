// Command csharp-build publishes and stages the C# NativeAOT runtime and example plugins.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

const runtimeProject = "csharp/Dragonfly.Runtime/Dragonfly.Runtime.csproj"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "csharp-build:", err)
		os.Exit(1)
	}
}

func run() error {
	root := flag.String("root", ".", "repository root")
	action := flag.String("action", "publish", "publish, stage, prepare, or clean")
	rid := flag.String("rid", "", ".NET runtime identifier; inferred from the host when empty")
	flag.Parse()

	absRoot, err := filepath.Abs(*root)
	if err != nil {
		return fmt.Errorf("resolve root: %w", err)
	}
	projects, err := exampleProjects(absRoot)
	if err != nil {
		return err
	}

	switch *action {
	case "publish":
		return publish(absRoot, projects, *rid)
	case "stage":
		return stageExamples(absRoot)
	case "prepare":
		return os.MkdirAll(filepath.Join(absRoot, "build"), 0o755)
	case "clean":
		return clean(absRoot, projects)
	default:
		return fmt.Errorf("unknown action %q", *action)
	}
}

func publish(root string, projects []string, rid string) error {
	if rid == "" {
		var err error
		rid, err = dotnetRID(runtime.GOOS, runtime.GOARCH)
		if err != nil {
			return err
		}
	}
	targetOS, err := targetOS(rid)
	if err != nil {
		return err
	}
	extension := nativeExtension(targetOS)
	runtimeOutput := filepath.Join(root, "build", "dotnet", "runtime")
	exampleOutputs := filepath.Join(root, "build", "dotnet", "examples")
	if err := os.RemoveAll(runtimeOutput); err != nil {
		return err
	}
	if err := os.RemoveAll(exampleOutputs); err != nil {
		return err
	}

	if err := dotnet(root, "publish", filepath.Join(root, filepath.FromSlash(runtimeProject)), "-c", "Release", "-r", rid, "--self-contained", "true", "-o", runtimeOutput); err != nil {
		return err
	}
	for _, project := range projects {
		name := filepath.Base(filepath.Dir(project))
		output := filepath.Join(exampleOutputs, name)
		if err := dotnet(root, "publish", project, "-c", "Release", "-r", rid, "--self-contained", "true", "-o", output); err != nil {
			return err
		}

		files, err := nativeFiles(output, extension)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			return fmt.Errorf("%s produced no %s library", project, extension)
		}
	}

	libDirectory := filepath.Join(root, "build", "lib")
	pluginDirectory := filepath.Join(root, "build", "plugins")
	if err := clearNativeLibraries(libDirectory); err != nil {
		return err
	}
	if err := clearNativeLibraries(pluginDirectory); err != nil {
		return err
	}
	runtimeSource := filepath.Join(runtimeOutput, "Dragonfly.Runtime"+extension)
	if err := copyFile(runtimeSource, filepath.Join(libDirectory, runtimeLibraryName(targetOS))); err != nil {
		return fmt.Errorf("stage runtime: %w", err)
	}
	for _, project := range projects {
		name := filepath.Base(filepath.Dir(project))
		output := filepath.Join(exampleOutputs, name)
		if _, err := copyNativeLibraries(output, pluginDirectory, extension); err != nil {
			return fmt.Errorf("stage %s: %w", name, err)
		}
	}
	return nil
}

func stageExamples(root string) error {
	libDirectory := filepath.Join(root, "examples", "lib")
	pluginDirectory := filepath.Join(root, "examples", "plugins")
	if err := clearNativeLibraries(libDirectory); err != nil {
		return err
	}
	if err := clearNativeLibraries(pluginDirectory); err != nil {
		return err
	}
	runtimes, err := copyAllNativeLibraries(filepath.Join(root, "build", "lib"), libDirectory)
	if err != nil {
		return err
	}
	if runtimes == 0 {
		return errors.New("no built runtime library found; run the publish action first")
	}
	plugins, err := copyAllNativeLibraries(filepath.Join(root, "build", "plugins"), pluginDirectory)
	if err != nil {
		return err
	}
	if plugins == 0 {
		return errors.New("no built example plugins found; run the publish action first")
	}
	return nil
}

func clean(root string, projects []string) error {
	if err := dotnet(root, "clean", filepath.Join(root, filepath.FromSlash(runtimeProject))); err != nil {
		return err
	}
	for _, project := range projects {
		if err := dotnet(root, "clean", project); err != nil {
			return err
		}
	}
	if err := os.RemoveAll(filepath.Join(root, "build")); err != nil {
		return err
	}
	if err := os.RemoveAll(filepath.Join(root, "examples", "lib")); err != nil {
		return err
	}
	return clearNativeLibraries(filepath.Join(root, "examples", "plugins"))
}

func exampleProjects(root string) ([]string, error) {
	projects, err := filepath.Glob(filepath.Join(root, "examples", "plugins", "*", "*.csproj"))
	if err != nil {
		return nil, fmt.Errorf("find example projects: %w", err)
	}
	sort.Strings(projects)
	return projects, nil
}

func dotnet(root string, args ...string) error {
	command := exec.Command("dotnet", args...)
	command.Dir = root
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("dotnet %s: %w", strings.Join(args, " "), err)
	}
	return nil
}

func dotnetRID(goos, goarch string) (string, error) {
	osName, ok := map[string]string{"windows": "win", "linux": "linux", "darwin": "osx"}[goos]
	if !ok {
		return "", fmt.Errorf("no .NET RID for %s/%s", goos, goarch)
	}
	architecture, ok := map[string]string{"amd64": "x64", "arm64": "arm64"}[goarch]
	if !ok {
		return "", fmt.Errorf("no .NET RID for %s/%s", goos, goarch)
	}
	return osName + "-" + architecture, nil
}

func targetOS(rid string) (string, error) {
	switch {
	case strings.HasPrefix(rid, "win-"):
		return "windows", nil
	case strings.HasPrefix(rid, "linux-"):
		return "linux", nil
	case strings.HasPrefix(rid, "osx-"):
		return "darwin", nil
	default:
		return "", fmt.Errorf("cannot determine target OS from RID %q", rid)
	}
}

func nativeExtension(goos string) string {
	switch goos {
	case "windows":
		return ".dll"
	case "darwin":
		return ".dylib"
	default:
		return ".so"
	}
}

func runtimeLibraryName(goos string) string {
	switch goos {
	case "windows":
		return "dragonfly_plugin_runtime.dll"
	case "darwin":
		return "libdragonfly_plugin_runtime.dylib"
	default:
		return "libdragonfly_plugin_runtime.so"
	}
}

func clearNativeLibraries(directory string) error {
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() && isNativeLibrary(entry.Name()) {
			if err := os.Remove(filepath.Join(directory, entry.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyAllNativeLibraries(source, destination string) (int, error) {
	if err := os.MkdirAll(destination, 0o755); err != nil {
		return 0, err
	}
	entries, err := os.ReadDir(source)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, entry := range entries {
		if entry.IsDir() || !isNativeLibrary(entry.Name()) {
			continue
		}
		if err := copyFile(filepath.Join(source, entry.Name()), filepath.Join(destination, entry.Name())); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func copyNativeLibraries(source, destination, extension string) (int, error) {
	files, err := nativeFiles(source, extension)
	if err != nil {
		return 0, err
	}
	if err := os.MkdirAll(destination, 0o755); err != nil {
		return 0, err
	}
	for _, file := range files {
		if err := copyFile(file, filepath.Join(destination, filepath.Base(file))); err != nil {
			return 0, err
		}
	}
	return len(files), nil
}

func nativeFiles(directory, extension string) ([]string, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.EqualFold(filepath.Ext(entry.Name()), extension) {
			files = append(files, filepath.Join(directory, entry.Name()))
		}
	}
	sort.Strings(files)
	return files, nil
}

func isNativeLibrary(name string) bool {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".dll", ".so", ".dylib":
		return true
	default:
		return false
	}
}

func copyFile(source, destination string) error {
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()
	info, err := input.Stat()
	if err != nil {
		return err
	}
	output, err := os.OpenFile(destination, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode().Perm())
	if err != nil {
		return err
	}
	if _, err := io.Copy(output, input); err != nil {
		_ = output.Close()
		return err
	}
	return output.Close()
}
