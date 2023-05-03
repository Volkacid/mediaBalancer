package service

import (
	"math/rand"
	"net"
	"time"
)

type AmazonServices struct {
	activeServers map[string]struct{}
}

func ConnectToCloud() *AmazonServices {
	rand.Seed(time.Now().Unix())
	return &AmazonServices{activeServers: make(map[string]struct{}, 10)}
}

func (as *AmazonServices) CreateNewServer() string {
	ipBuf := make([]byte, 4)

	for { //Генерация нового уникального адреса
		rand.Read(ipBuf)
		ipStr := net.IP(ipBuf).String()
		if _, ok := as.activeServers[ipStr]; !ok {
			as.activeServers[ipStr] = struct{}{}
			return ipStr
		}
	}

}

func (as *AmazonServices) DestroyServer(ip string) {

	delete(as.activeServers, ip)
}
