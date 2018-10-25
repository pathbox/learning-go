package main

import (
	"fmt"

	"github.com/shirou/gopsutil/net"
)

func main() {
	r, _ := net.IOCounters(true)
	fmt.Println(r)

}

// example of `netstat -ibdnW` output on yosemite
// Name  Mtu   Network       Address            Ipkts Ierrs     Ibytes    Opkts Oerrs     Obytes  Coll Drop
// lo0   16384 <Link#1>                        869107     0  169411755   869107     0  169411755     0   0
// lo0   16384 ::1/128     ::1                 869107     -  169411755   869107     -  169411755     -   -
// lo0   16384 127           127.0.0.1         869107     -  169411755   869107     -  169411755     -   -
