package controllers

import (
	"net"

	"github.com/google/uuid"
	"gopkg.in/olahol/melody.v1"
)

func GetMsgID() string {
	return uuid.New().String()
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
