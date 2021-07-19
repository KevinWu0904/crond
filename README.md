# CronD（WIP）
[![Actions Status](https://github.com/KevinWu0904/crond/actions/workflows/go.yml/badge.svg)](https://github.com/KevinWu0904/crond/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/KevinWu0904/crond)](https://goreportcard.com/report/github.com/KevinWu0904/crond)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://github.com/KevinWu0904/crond/blob/main/LICENSE)

CronD is a **Cloud Native** golang distributed cron scheduling service.

CronD serves a distributed unified job dispatcher for offline periodic tasks. It is recommended running in
a cluster with 3 or 5 nodes, peer nodes communicate by **Raft Consensus**.

## Architecture
![CronD Architecture](./docs/images/architecture.png)