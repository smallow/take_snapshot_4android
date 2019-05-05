package main

import (
	"flag"
	"log"
	"net"
	"spirit2/websocket"
	"time"
)

var (
	ip          = "192.168.20.33"
	port        = "8080"
	wsHandleUrl = "/ws"
	netWorkName = "eth0"
)

var (
	mac                  = ""
	clientIp             = ""
	interval             = 10 //客户端心跳间隔(单位秒)
	takeSnapshotInterval = 3  //截图周期单位秒
)

func main() {
	p := flag.String("port", "8080", "服务端口号")
	i := flag.String("ip", "192.168.20.33", "服务端ip")
	n := flag.String("netWorkName", "wlan0", "网卡名称") //mac为en0 无线网卡地址为wlan0,有线为eth0
	flag.Parse()
	if *p != "" || *i != "" || *n != "" {
		port = *p
		ip = *i
		netWorkName = *n
	}
	interfaces, _ := net.InterfaceByName(netWorkName)
	if interfaces != nil && interfaces.HardwareAddr.String() != "" {
		mac = interfaces.HardwareAddr.String()
	}
	getClientIp()
	log.Printf("服务端ip:[%s],端口:[%s],网卡:[%s],Mac地址:[%s],客户端IP:[%s]", ip, port, netWorkName, mac, clientIp)
	go startTakeSnapshotJob()
	websocket.StartConnect(ip, port, wsHandleUrl, time.Duration(interval)*time.Second)

}

func getClientIp() {
	addrs, _ := net.InterfaceAddrs()
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				clientIp = ipnet.IP.String()
			}
		}
	}
}
