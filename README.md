<div align="center">

# Masjid's Droid

<img src="./assets/img/fares___blueprint.png" width="512" alt="M-Droid Blueprint"/>

![Latest commit](https://img.shields.io/github/last-commit/ccil-kbw/robot/master?style=flat-square)
![GitHub issues](https://img.shields.io/github/issues/ccil-kbw/robot)
[![Publish](https://github.com/ccil-kbw/robot/actions/workflows/publish.yml/badge.svg)](https://github.com/ccil-kbw/robot/actions/workflows/publish.yml)
![GitHub](https://img.shields.io/github/license/ccil-kbw/robot)

</div>

### Table of Contents
- [Introduction](#introduction)
- [info](#info)
- [License](#license)

# Introduction

The Masjid's Droid was initially created to serve the needs of the beloved Community of the Masjid Khaled Ben Walid but is intended to be Open and available for the Ummah, East to West.

# Info

Documentation about installation and usage are found at [https://ccil-kbw.github.io](https://ccil-kbw.github.io)

# Setup Development Environment

## Prerequisite
- Go 1.22
- [OBS Studio](https://obsproject.com/download)
- [Discord Bot](https://discord.com/developers/docs/quick-start/getting-started)
  - To run and test your Discord Commands in Developer Mode

## Sample Commands
### Verify Go Version

```bash
> go version
go version go1.22.1 windows/amd64
```

### Iquama CLI
```bash
> go run cmd/cli/main.go
+------------+---------+---------+---------+---------+---------+
| DATE       | FAJR    | DHUHUR  | ASR     | MAGHRIB | ISHA    |
+------------+---------+---------+---------+---------+---------+
| 04/16/2024 | 5:15 am | 1:15 pm | 6:15 pm | 7:45 pm | 9:10 pm |
+------------+---------+---------+---------+---------+---------+
```

# License

BSD-3, see LICENSE

