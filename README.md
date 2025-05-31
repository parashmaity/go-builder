# Go Builder

A flexible and powerful build automation tool for Go projects that simplifies the process of building, testing, and releasing Go applications across multiple platforms.

## Features

- Multi-platform build support (Linux, Windows, macOS)
- Customizable build configurations via YAML
- Built-in dependency management
- Automated pre-build and post-build tasks
- Integrated testing and linting
- Version management and release automation

## Installation

```bash
go get github.com/parashmaity/go-builder
```


# Installation Options

## Using go install (requires Go)
```bash
go install github.com/parashmaity/go-builder@latest
```

## From Pre-built Binaries

1. Download the appropriate binary for your system from the [Releases page](https://github.com/parashmaity/go-builder/releases)
2. Make it executable:
   ```bash
   chmod +x go-builder-*
   ```
3. Move it to your PATH:
   ```bash
   sudo mv go-builder-* /usr/local/bin/go-builder
   ```

## Verify Installation
```bash
go-builder --version
```

## Quick Start

1. Create a `build-config.yaml` in your project root:

```yaml
project: "your-project"
version: "1.0.0"
zip: false

build:
  default:
    os: "darwin"
    arch: "arm64"
    output: "./bin/{{.Project}}-{{.Version}}-darwin-arm64"
```

2. Run the builder:

```bash
go-builder build
```

## Configuration

### Build Configuration

The `build-config.yaml` file supports the following sections:

#### Project Settings
- `project`: Project name
- `version`: Project version

#### Build Targets
```yaml
build:
  default:
    os: "darwin"    # Target operating system
    arch: "arm64"   # Target architecture
    output: "./bin/{{.Project}}-{{.Version}}-darwin-arm64"
    ldflags: "-X main.version={{.Version}}"
    tags: []
    cgo: false
```

#### Custom Build Targets
```yaml
  targets:
    linux-amd64:
      os: "linux"
      arch: "amd64"
      output: "./bin/{{.Project}}-{{.Version}}-linux-amd64"
```

### Dependencies
```yaml
dependencies:
  check:
    - "golang.org/x/tools/cmd/stringer"
    - "github.com/golang/mock/mockgen"
```

### Tasks
```yaml
tasks:
  pre-build:
    - "go generate ./..."
    - "go mod tidy"

  post-build:
    - "chmod +x ./bin/*"

  test:
    - "go test -v -cover ./..."

  lint:
    - "golangci-lint run ./..."

  release:
    - "git tag v{{.Version}}"
    - "git push origin v{{.Version}}"
```

## Commands

- Build the project:
  ```bash
  # Build default target
  go-builder build

  # Build specific target
  go-builder build --target linux-amd64

  # Build all targets
  go-builder build --target all
  ```

- Run tests:
  ```bash
  go-builder test
  ```

- Run linting:
  ```bash
  go-builder lint
  ```

- Create and push a release:
  ```bash
  go-builder release
  ```

## Template Variables

The following variables are available in the configuration file:

- `{{.Project}}`: Project name from config
- `{{.Version}}`: Project version from config
- `{{.OutputDir}}`: Output directory path

## Supported Platforms

The builder supports the following platform combinations:

- Linux: 386, amd64, arm, arm64
- Windows: 386, amd64, arm, arm64
- macOS: amd64, arm64
- FreeBSD: 386, amd64
- OpenBSD: amd64
- NetBSD: amd64
- DragonFly BSD: amd64
- Solaris: amd64

## Complete build-config.yaml
```yaml
# Project settings
project: "example-project"
version: "1.0.0"
zip: false

build:
  # Default build settings
  default:
    os: "darwin"
    arch: "arm64"
    output: "./bin/{{.Project}}-{{.Version}}-darwin-arm64"
    ldflags: "-X main.version={{.Version}}"
    tags: []
    cgo: false

  # Custom build targets
  targets:
    # Linux targets
    linux-amd64:
      os: "linux"
      arch: "amd64"
      output: "./bin/{{.Project}}-{{.Version}}-linux-amd64"
    linux-386:
      os: "linux"
      arch: "386"
      output: "./bin/{{.Project}}-{{.Version}}-linux-386"
    linux-arm:
      os: "linux"
      arch: "arm"
      output: "./bin/{{.Project}}-{{.Version}}-linux-arm"
    linux-arm64:
      os: "linux"
      arch: "arm64"
      output: "./bin/{{.Project}}-{{.Version}}-linux-arm64"

    # Windows targets
    windows-amd64:
      os: "windows"
      arch: "amd64"
      output: "./bin/{{.Project}}-{{.Version}}-windows-amd64.exe"
    windows-386:
      os: "windows"
      arch: "386"
      output: "./bin/{{.Project}}-{{.Version}}-windows-386.exe"
    windows-arm:
      os: "windows"
      arch: "arm"
      output: "./bin/{{.Project}}-{{.Version}}-windows-arm.exe"
    windows-arm64:
      os: "windows"
      arch: "arm64"
      output: "./bin/{{.Project}}-{{.Version}}-windows-arm64.exe"

    # macOS targets
    darwin-amd64:
      os: "darwin"
      arch: "amd64"
      output: "./bin/{{.Project}}-{{.Version}}-darwin-amd64"
    darwin-arm64:
      os: "darwin"
      arch: "arm64"
      output: "./bin/{{.Project}}-{{.Version}}-darwin-arm64"

    # FreeBSD targets
    freebsd-386:
      os: "freebsd"
      arch: "386"
      output: "./bin/{{.Project}}-{{.Version}}-freebsd-386"
    freebsd-amd64:
      os: "freebsd"
      arch: "amd64"
      output: "./bin/{{.Project}}-{{.Version}}-freebsd-amd64"

    # OpenBSD targets
    openbsd-amd64:
      os: "openbsd"
      arch: "amd64"
      output: "./bin/{{.Project}}-{{.Version}}-openbsd-amd64"

    # NetBSD targets
    netbsd-amd64:
      os: "netbsd"
      arch: "amd64"
      output: "./bin/{{.Project}}-{{.Version}}-netbsd-amd64"

    # DragonFly BSD targets
    dragonfly-amd64:
      os: "dragonfly"
      arch: "amd64"
      output: "./bin/{{.Project}}-{{.Version}}-dragonfly-amd64"

    # Solaris targets
    solaris-amd64:
      os: "solaris"
      arch: "amd64"
      output: "./bin/{{.Project}}-{{.Version}}-solaris-amd64"

    # Build for all platforms
    all:
      platforms:
        - linux-386
        - linux-amd64
        - linux-arm
        - linux-arm64
        - windows-386
        - windows-amd64
        - windows-arm
        - windows-arm64
        - darwin-amd64
        - darwin-arm64
        - freebsd-386
        - freebsd-amd64
        - openbsd-amd64
        - netbsd-amd64
        - dragonfly-amd64
        - solaris-amd64

dependencies:
  check:
    - "golang.org/x/tools/cmd/stringer"
    - "github.com/golang/mock/mockgen"

tasks:
  pre-build:
    - "go generate ./..."
    - "go mod tidy"

  post-build:
    - "chmod +x ./bin/*"

  test:
    - "go test -v -cover ./..."

  lint:
    - "golangci-lint run ./..."

  release:
    - "git tag v{{.Version}}"
    - "git push origin v{{.Version}}"
```
## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support, please open an issue in the GitHub repository.
