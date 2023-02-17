package main

import (
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

var IpMap map[string]bool

func dnsquery(domain string, ip string, DnsServer string, OnlyIp bool, repeat int, v6 bool) {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}
	c := new(dns.Client)
	m := new(dns.Msg)
	if v6 {
		m.SetQuestion(domain, dns.TypeAAAA)
	} else {
		m.SetQuestion(domain, dns.TypeA)
	}

	o := new(dns.OPT)
	o.Hdr.Name = "."
	o.Hdr.Rrtype = dns.TypeOPT
	e := new(dns.EDNS0_SUBNET) //EDNS
	e.Code = dns.EDNS0SUBNET
	if v6 {
		e.Family = 2         // 1 IPv4 2 IPv6
		e.SourceNetmask = 56 //  地址掩码 ipv4 一般为 /24  ipv6为 /56
		e.Address = net.ParseIP(ip).To16()
	} else {
		e.Family = 1
		e.SourceNetmask = 24
		e.Address = net.ParseIP(ip).To4()
	}

	e.SourceScope = 0
	o.Option = append(o.Option, e)
	m.Extra = append(m.Extra, o)
	for i := 0; i < repeat; i++ {
		in, _, err := c.Exchange(m, DnsServer) //注意:要选择支持自定义EDNS的DNS 或者是 目标NS服务器  国内DNS大部分不支持自定义EDNS数据

		if err != nil {

		}
		for _, answer := range in.Answer {

			if answer.Header().Rrtype == dns.TypeA {
				if OnlyIp {
					IpMap[answer.(*dns.A).A.String()] = true
				} else {
					print("CDN:")
					httpGet(answer.(*dns.A).A.String())
					println(" ", answer.(*dns.A).A.String())
				}
			} else if answer.Header().Rrtype == dns.TypeAAAA {
				IpMap[answer.(*dns.AAAA).AAAA.String()] = true
			}
		}
	}

}

func main() {
	Initlist()
	var domain = flag.String("d", "www.taobao.com", "domain")
	var DnsServer = flag.String("s", "119.29.29.29:53", "dns server addr")
	var ip = flag.String("ip", "", "client ip")
	var OnlyIp = flag.Bool("i", false, "Only output ip addr")
	var repeat = flag.Int("r", 1, "repeat query rounds")
	var v6 = flag.Bool("6", false, "query AAAA (ipv6)")
	flag.Parse()
	fmt.Println("开始请求：domain", *domain, "DnsServer:", *DnsServer)
	IpMap = make(map[string]bool)
	if (*ip != "") || (*v6) {
		*OnlyIp = true
		dnsquery(*domain, *ip, *DnsServer, *OnlyIp, *repeat, *v6)
	} else {
		for city, ip := range CityMap {
			if !*OnlyIp {
				fmt.Println()
				fmt.Println("ECS:", city, ip)
			}
			dnsquery(*domain, ip, *DnsServer, *OnlyIp, *repeat, *v6)

		}
	}

	if *OnlyIp {
		for ip, _ := range IpMap {
			println(ip)
		}
	}

}

func httpGet(ip string) {
	resp, err := http.Get("https://ip.useragentinfo.com/json?ip=" + ip)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	str := string(body)
	province := gjson.Get(str, "province")
	city := gjson.Get(str, "city")
	isp := gjson.Get(str, "isp")
	fmt.Print(province, city, isp)
}
