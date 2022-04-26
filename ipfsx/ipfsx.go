package ipfsx

import (
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
	"log"
	"os/exec"
	"strings"
)

func GetPeers() []string {

	peers, err := exec.Command("ipfs", "swarm", "peers").Output()
	if err != nil {
		log.Fatal("Error querying ipfs. Is your ipfs node up?", err)
	}

	peersStrs := strings.Split(string(peers), "\n")

	var ips []string

	for _, p := range peersStrs {
		if p == "" {
			continue
		}
		multiaddr, err := ma.NewMultiaddr(p)
		if err != nil {
			log.Fatal(err)
			return nil
		}

		ip, err := manet.ToIP(multiaddr)
		if err != nil {
			log.Fatal(err)
			return nil
		}

		if !stringInSlice(ip.String(), ips) {
			ips = append(ips, ip.String())
		}
	}

	return ips
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
