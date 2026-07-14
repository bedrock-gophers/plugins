# Plugins

Each source plugin owns its `.csproj`. `make run` publishes NativeAOT shared libraries here before starting the server. Compatible precompiled `.so` files may also be placed here.

`kitchen-sink` uses every exposed C# API and grows with each parity slice.

`vanilla-commands` keeps each command in its own file and grows as more of Dragonfly's gameplay API is exposed.
