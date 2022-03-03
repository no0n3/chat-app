package main

import (
	"net"
)

func IsLocalIp(ip string) (bool, error) {
	localIp, err := GetIp()
	if err != nil {
		return false, err
	}

	return localIp == ip, nil
}

func GetIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", nil
}
