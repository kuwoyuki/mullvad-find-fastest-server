package servers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type MullvadServer struct {
	Hostname string `json:"hostname"`
	// CountryCode      string `json:"country_code"`
	// CountryName      string `json:"country_name"`
	// CityCode         string `json:"city_code"`
	// CityName         string `json:"city_name"`
	Active bool `json:"active"`
	// Owned            bool   `json:"owned"`
	// Provider         string `json:"provider"`
	// IPV4AddrIn       string `json:"ipv4_addr_in"`
	// IPV6AddrIn       string `json:"ipv6_addr_in"`
	// NetworkPortSpeed uint   `json:"network_port_speed"`
	Type string `json:"type"`
	// StatusMessages   []string `json:"status_messages"`
}

var mullvadAddr = ".mullvad.net"

func GetServers() []string {
	resp, err := http.Get("https://api.mullvad.net/www/relays/all/")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var mullvadServers []MullvadServer
	var servers []string
	err = json.Unmarshal(body, &mullvadServers)
	if err != nil {
		log.Fatalln(err)
	}

	for _, s := range mullvadServers {
		if s.Active == true && s.Type == "wireguard" {
			servers = append(servers, s.Hostname+mullvadAddr)
		}
	}

	return servers
}
