# mug

Super simple and personal package manager.

## Why?

Recently I've started to prefer installing tools simply by grabbing an executable from GitHub Releases and symlinking
them somewhere on my PATH. Mug automates that process. Mug is also a nice excuse to finally learn Go.

Idea inspired by [hysp](https://github.com/pwnwriter/hysp).

## Usage

> Mug only works on x86\_64 Linux, since I use it mostly on WSL.

```
mug install <package-name>
```

Package information goes in [`registry.json`](./resources/registry.json). Mug downloads the files mentioned in the
registry then symlinks them to `~/.local/bin/`. Mug will always download the latest available version.

Mug downloads all packages to `~/.mug/`

