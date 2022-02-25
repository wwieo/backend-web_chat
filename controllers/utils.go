package controllers

import (
	"net"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gopkg.in/olahol/melody.v1"
)

type UtilsController struct {
}

func NewUtilsController() *UtilsController {
	return &UtilsController{}
}

func (utilsController *UtilsController) GetMsgID() string {
	return uuid.New().String()
}

func (utilsController *UtilsController) GetUsername(session *melody.Session) string {
	return session.Request.URL.Query().Get("username")
}

func (utilsController *UtilsController) GetUserIP() string {
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

func (utilsController *UtilsController) GetConfig() (config *viper.Viper) {

	config = viper.New()
	config.AddConfigPath("./config")
	config.SetConfigName("config")
	config.SetConfigType("yaml")

	if err := config.ReadInConfig(); err != nil {
		panic("config err: " + err.Error())
	}
	return config
}
