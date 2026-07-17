[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [ValidateSet("CheckRevision", "Prepare", "Build", "Clean")]
    [string]$Action,

    [Parameter(Mandatory = $true)]
    [ValidatePattern("^[0-9a-fA-F]{40}$")]
    [string]$FrameworkRevision,

    [string]$CachePath = ".cache/bedrock-gophers",
    [string]$BuildPath = ".build",
    [string]$DotnetRid = "win-x64"
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$projectRoot = Split-Path -Parent $PSScriptRoot
$cacheRoot = [IO.Path]::GetFullPath((Join-Path $projectRoot $CachePath))
$buildRoot = [IO.Path]::GetFullPath((Join-Path $projectRoot $BuildPath))
$pluginsRoot = Join-Path $projectRoot "plugins"
$libRoot = Join-Path $projectRoot "lib"

function Invoke-Checked {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Command,

        [Parameter(ValueFromRemainingArguments = $true)]
        [string[]]$Arguments
    )

    & $Command @Arguments
    if ($LASTEXITCODE -ne 0) {
        throw "Command failed with exit code ${LASTEXITCODE}: $Command $($Arguments -join ' ')"
    }
}

function Test-FrameworkRevision {
    $versionOutput = & go list -m -f "{{.Version}}" github.com/bedrock-gophers/plugins 2>&1
    if ($LASTEXITCODE -ne 0) {
        throw "Unable to read the Go framework revision: $($versionOutput -join [Environment]::NewLine)"
    }

    $moduleVersion = [string]($versionOutput | Select-Object -Last 1)
    $actualRevision = ($moduleVersion.Trim() -split "-")[-1]
    $expectedRevision = $FrameworkRevision.Substring(0, 12)

    if ($actualRevision -ne $expectedRevision) {
        throw "Go framework is pinned to $actualRevision, expected $expectedRevision"
    }

    Write-Host "[OK] Go framework revision: $actualRevision"
}

function Initialize-FrameworkCache {
    $cacheParent = Split-Path -Parent $cacheRoot
    New-Item -ItemType Directory -Force -Path $cacheParent | Out-Null

    if (-not (Test-Path -LiteralPath (Join-Path $cacheRoot ".git"))) {
        Invoke-Checked -Command "git" -Arguments @(
            "clone",
            "--quiet",
            "https://github.com/bedrock-gophers/plugins",
            $cacheRoot
        )
    }

    Invoke-Checked -Command "git" -Arguments @("-C", $cacheRoot, "fetch", "--quiet", "origin")
    Invoke-Checked -Command "git" -Arguments @("-C", $cacheRoot, "checkout", "--quiet", $FrameworkRevision)
    Write-Host "[OK] Framework cache: $cacheRoot"
}

function Build-WindowsArtifacts {
    $runtimeProject = Join-Path $cacheRoot "csharp/Dragonfly.Runtime/Dragonfly.Runtime.csproj"
    $runtimeOutput = Join-Path $buildRoot "dotnet/runtime"
    $pluginOutputRoot = Join-Path $buildRoot "dotnet/plugins"
    $binRoot = Join-Path $buildRoot "bin"

    New-Item -ItemType Directory -Force -Path $runtimeOutput, $pluginOutputRoot, $binRoot, $libRoot | Out-Null

    Invoke-Checked -Command "dotnet" -Arguments @(
        "publish",
        $runtimeProject,
        "-c", "Release",
        "-r", $DotnetRid,
        "--self-contained", "true",
        "-o", $runtimeOutput
    )

    Get-ChildItem -LiteralPath $pluginsRoot -Filter "*.dll" -File -ErrorAction SilentlyContinue |
        Remove-Item -Force

    $pluginProjects = Get-ChildItem -LiteralPath $pluginsRoot -Directory |
        ForEach-Object { Get-ChildItem -LiteralPath $_.FullName -Filter "*.csproj" -File } |
        Sort-Object FullName

    foreach ($project in $pluginProjects) {
        $pluginName = $project.Directory.Name
        $pluginOutput = Join-Path $pluginOutputRoot $pluginName
        New-Item -ItemType Directory -Force -Path $pluginOutput | Out-Null

        Invoke-Checked -Command "dotnet" -Arguments @(
            "publish",
            $project.FullName,
            "-c", "Release",
            "-r", $DotnetRid,
            "--self-contained", "true",
            "-p:DragonflyFrameworkRoot=$cacheRoot",
            "-o", $pluginOutput
        )

        $nativeLibraries = @(Get-ChildItem -LiteralPath $pluginOutput -Filter "*.dll" -File)
        if ($nativeLibraries.Count -eq 0) {
            throw "No native DLL was produced for $($project.FullName)"
        }

        $nativeLibraries | Copy-Item -Destination $pluginsRoot -Force
        Write-Host "[OK] Plugin: $pluginName"
    }

    $runtimeLibrary = Join-Path $runtimeOutput "Dragonfly.Runtime.dll"
    if (-not (Test-Path -LiteralPath $runtimeLibrary)) {
        throw "Native runtime not found: $runtimeLibrary"
    }

    Copy-Item -LiteralPath $runtimeLibrary -Destination (Join-Path $libRoot "dragonfly_plugin_runtime.dll") -Force
    Invoke-Checked -Command "go" -Arguments @("build", "-o", (Join-Path $binRoot "server.exe"), ".")

    Write-Host "[OK] Runtime: $(Join-Path $libRoot 'dragonfly_plugin_runtime.dll')"
    Write-Host "[OK] Server:  $(Join-Path $binRoot 'server.exe')"
}

function Remove-BuildArtifacts {
    foreach ($path in @(
        $cacheRoot,
        (Join-Path $projectRoot ".data"),
        $buildRoot,
        $libRoot
    )) {
        if (Test-Path -LiteralPath $path) {
            Remove-Item -LiteralPath $path -Recurse -Force
        }
    }

    Get-ChildItem -LiteralPath $pluginsRoot -Filter "*.dll" -File -ErrorAction SilentlyContinue |
        Remove-Item -Force

    Get-ChildItem -LiteralPath $pluginsRoot -Directory -Recurse |
        Where-Object { $_.Name -in @("bin", "obj") } |
        Sort-Object FullName -Descending |
        Remove-Item -Recurse -Force

    Write-Host "[OK] Windows build artifacts removed."
}

Push-Location $projectRoot
try {
    switch ($Action) {
        "CheckRevision" { Test-FrameworkRevision }
        "Prepare" { Initialize-FrameworkCache }
        "Build" { Build-WindowsArtifacts }
        "Clean" { Remove-BuildArtifacts }
    }
}
finally {
    Pop-Location
}
