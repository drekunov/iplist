# iplist

A command-line utility for converting IP addresses and CIDR blocks into MikroTik RouterOS route configuration commands.

## Overview

`iplist` reads a list of IP addresses and CIDR notations from a file and generates MikroTik RouterOS `/ip route add` commands with a specified gateway and comment. This is useful for managing blocklists, lists, or custom routing rules on MikroTik devices.

## Features

- Supports both single IP addresses and CIDR notation
- Automatically deduplicates entries
- Filters out IPv6 addresses (processes only IPv4)
- Handles large files efficiently with buffering (4MB buffer)
- Supports comments in input files (lines starting with `#`)
- Generates MikroTik-compatible route commands

## Installation

### Prerequisites
- Go 1.24.4 or later

### Build
```bash
go build
```

This will create an executable named `iplist` in the current directory.

## Usage

```bash
./iplist <ip_list_file> <gateway> <comment>
```

### Arguments
- `<ip_list_file>`: Path to a text file containing IP addresses or CIDR blocks (one per line)
- `<gateway>`: Gateway IP address for the routes
- `<comment>`: Comment to add to all generated routes

### Example
```bash
./iplist list.txt 192.168.1.1 mylist
```

## Input File Format

The input file should contain one IP address or CIDR block per line. Comments are supported:

```
# List
1.179.201.18
103.224.182.250
103.224.212.221/24
104.16.10.0/16

# IPv6 addresses are automatically filtered out
2001:4860:4860::8888

# Empty lines and comments are ignored
10.0.0.0/8
```

## Output

The tool outputs MikroTik RouterOS `/ip route` commands suitable for import:

```
/ip route
add dst-address=1.179.201.18/32 gateway=192.168.1.1 comment="mylist"
add dst-address=103.224.182.250/32 gateway=192.168.1.1 comment="mylist"
add dst-address=103.224.212.221/24 gateway=192.168.1.1 comment="mylist"
add dst-address=104.16.10.0/16 gateway=192.168.1.1 comment="mylist"
add dst-address=10.0.0.0/8 gateway=192.168.1.1 comment="mylist"
```

You can redirect this output to a file and import it into MikroTik:

```bash
./iplist list.txt 192.168.1.1 mylist > routes.rsc
```

Then import in MikroTik:
```
[admin@MikroTik] > import file-name=routes.rsc
```

## How It Works

1. Reads the input file line by line
2. Skips empty lines and comments (lines starting with `#`)
3. Parses each line as either:
   - A CIDR block (e.g., `192.168.1.0/24`)
   - A single IP address (automatically converted to `/32` for IPv4)
4. Filters out duplicate entries
5. Removes IPv6 addresses
6. Generates and outputs MikroTik `/ip route add` commands

## Performance

The utility is optimized for large files:
- Uses a 4MB output buffer for efficient writing
- Scans input with a 1MB line buffer to handle large entries
- Deduplicates entries in-memory using a map for fast lookup

## Building from Source

```bash
git clone <repository>
cd iplist
go build -v
```
