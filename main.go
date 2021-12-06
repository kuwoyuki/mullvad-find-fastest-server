package main

import (
	"fmt"
	"sort"
	"time"

	ping "github.com/go-ping/ping"
	servers "github.com/victorb/mullvad-find-fastest-server/servers"
)

var pingCount = 3
var results = map[string]int64{}

func main() {
	for _, server := range servers.GetServers() {
		pinger, err := ping.NewPinger(server)
		if err != nil {
			panic(err)
		}
		// run as privileged as most distros don't allow access to unprivileged ICMP sockets
		pinger.SetPrivileged(true)
		// timeout event if https://github.com/go-ping/ping/pull/176 is merged to skipo server
		pinger.Timeout = 5 * 1000 * 1000 * 1000 // 5s
		pinger.Count = pingCount
		pinger.OnRecv = func(pkt *ping.Packet) {
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
