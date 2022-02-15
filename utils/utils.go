package utils

import (
	"math/rand"
	"net"
	"time"

	"gopkg.in/olahol/melody.v1"
)

func GetRandom() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func GetRandomI64() int64 {
	return GetRandom().Int63()
}

func GetUsername(session *melody.Session) string {
	return session.Request.URL.Query().Get("username")
}

func GetUserIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
