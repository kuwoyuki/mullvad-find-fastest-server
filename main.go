package main

import (
	"flag"
	"fmt"
	"sort"
	"time"

	"github.com/biter777/countries"
	probing "github.com/prometheus-community/pro-bing"
	servers "github.com/victorb/mullvad-find-fastest-server/servers"
)

var results = map[string]int64{}
var mullvadAddr = ".relays.mullvad.net"

func main() {
	vpnType := flag.String("type", "wireguard", "vpn type")
	vpnRegion := flag.String("region", "", "servers region")
	vpnCountry := flag.String("country", "", "servers country")
	pingCount := flag.Int("pingCount", 3, "ping count")
	flag.Parse()

	mullvadServers := servers.GetServers()

	var servers []string

	for _, s := range mullvadServers {
		if !s.Active || s.Type != *vpnType {
			continue
		}
		if *vpnCountry != "" && s.CountryCode != *vpnCountry {
			continue
		}
		if *vpnRegion != "" && countries.ByName(s.CountryCode).Region().String() != *vpnRegion {
			continue
		}
		servers = append(servers, s.Hostname+mullvadAddr)
	}

	for _, server := range servers {
		pinger, err := probing.NewPinger(server)
		if err != nil {
			panic(err)
		}
		// run as privileged as most distros don't allow access to unprivileged ICMP sockets
		pinger.SetPrivileged(true)
		// timeout event if https://github.com/go-ping/ping/pull/176 is merged to skipo server
		pinger.Timeout = 5 * 1000 * 1000 * 1000 // 5s
		pinger.Count = *pingCount
		pinger.OnRecv = func(pkt *probing.Packet) {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		}
		err = pinger.Run() // Blocks until finished.
		if err != nil {
			panic(err)
		}
		stats := pinger.Statistics() // get send/receive/rtt stats
		fmt.Printf("%s = %s\n", server, stats.AvgRtt.String())
		results[server] = stats.AvgRtt.Nanoseconds()
	}
	type kv struct {
		Key   string
		Value int64
	}

	var ss []kv
	for k, v := range results {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value < ss[j].Value
	})

	fmt.Println("## Final Results (least latency first):")
	for _, kv := range ss {
		durr, _ := time.ParseDuration(fmt.Sprintf("%dns", kv.Value))
		fmt.Printf("%s = %s\n", kv.Key, durr.String())
	}
}
