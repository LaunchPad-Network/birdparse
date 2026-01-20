# birdparse

A Go library for parsing BIRD command output.

## Features

- Parse BGP protocol information including session state, routes, and neighbor details
- Parse routing table data with BGP attributes
- Support for standard and large BGP communities
- Extract AS paths, next hops, and other BGP path attributes

## Installation

```bash
go get github.com/LaunchPad-Network/birdparse
```

## Usage

```go
import "github.com/LaunchPad-Network/birdparse"

// Parse BGP protocols
protocols := birdparse.ParseBGPProtocols(birdOutput)

// Parse routes
routes := birdparse.ParseRoutes(birdOutput)
```

## License

See [LICENSE](LICENSE) file.
