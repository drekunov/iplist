package main

import (
	"bufio"
	"fmt"
	"net"
	"net/netip"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <ip list file> <gateway> <comment>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s list.txt 192.168.1.1 mylist\n", os.Args[0])
		os.Exit(1)
	}

	iplistFile := os.Args[1]
	gateway := os.Args[2]
	comment := os.Args[3]

	file, err := os.Open(iplistFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening %s: %v\n", iplistFile, err)
		os.Exit(1)
	}
	defer file.Close()

	seen := make(map[string]struct{})
	w := bufio.NewWriterSize(os.Stdout, 4*1024*1024)
	fmt.Fprintln(w, "/ip route")

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		cidr := parseIPorCIDR(line)
		if cidr == "" {
			continue
		}
		if _, ok := seen[cidr]; ok {
			continue
		}
		prefix, err := netip.ParsePrefix(cidr)
		if err != nil {
			continue
		}
		if prefix.Addr().Is6() {
			continue
		}

		seen[cidr] = struct{}{}
		fmt.Fprintf(w, "add dst-address=%s gateway=%s comment=\"%s\"\n", cidr, gateway, comment)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading: %v\n", err)
		os.Exit(1)
	}

	w.Flush()
}

func parseIPorCIDR(s string) string {
	if _, network, err := net.ParseCIDR(s); err == nil {
		return network.String()
	}
	if ip := net.ParseIP(s); ip != nil {
		if ip.To4() != nil {
			return s + "/32"
		}
		return s + "/128"
	}
	return ""
}
