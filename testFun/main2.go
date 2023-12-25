// ConvertIpContent 处理前端传的ip段信息
package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/3th1nk/cidr"
)

func main() {
	strs := IpContent("10.10.10.76,10.10.10.20-50,")
	fmt.Print(strs)
}
func IpContent(ip string) (queries []string) {
	if len(ip) == 0 {
		return
	}

	if strings.IndexByte(ip, ',') > 0 {
		split := strings.Split(ip, ",")
		if len(split) < 2 {
			return
		}
		for _, value := range split {
			queries = append(queries, IpContent(value)...)
		}
	} else if strings.Contains(ip, "\r\n") {
		split := strings.Split(ip, "\r\n")
		if len(split) < 2 {
			return
		}
		for _, value := range split {
			queries = append(queries, IpContent(value)...)
		}
	} else if strings.Contains(ip, "\n") {
		split := strings.Split(ip, "\n")
		if len(split) < 2 {
			return
		}
		for _, value := range split {
			queries = append(queries, IpContent(value)...)
		}
	} else if strings.IndexByte(ip, '-') > 0 {
		if strings.IndexByte(ip, '*') > 0 { //存在- 并且存在*
			//ip = "10.10.1-11.*"
			split := strings.Split(ip, ".")     //10 10 1-11 *
			tem := strings.Split(split[2], "-") //1 11
			var ipStart, ipEnd string
			ipStart = split[0] + "." + split[1] + "." + tem[0] + "." + strconv.Itoa(0)
			ipEnd = split[0] + "." + split[1] + "." + tem[1] + "." + strconv.Itoa(255)
			cidrRange, err := IPv4RangeToCIDRRange(ipStart, ipEnd)
			if err != nil {
				return nil
			}
			fmt.Print("threatFIx----787878783", cidrRange)
			queries = append(queries, cidrRange...)

		} else { // 存在 - 但不存在*
			split := strings.Split(ip, "-") //split[0] = 10.10.10.10 split[1] = 20
			if len(split) != 2 {
				return
			}

			tem1 := strings.Split(split[0], ".") //i= [10 ,10,10,10]
			tem2 := strings.Split(split[1], ".")
			var ipStart, ipEnd string
			if len(tem2) == 4 { //解析 10.10.10.10-10.10.10.格式
				ipStart = split[0]
				ipEnd = split[1]
			} else {
				tem1[len(tem1)-1] = split[1] //i=[10,10,10,20]
				ipStart = split[0]
				ipEnd = strings.Join(tem1, ".")
			}
			cidrRange, err := IPv4RangeToCIDRRange(ipStart, ipEnd)
			if err != nil {
				return nil
			}
			fmt.Print("threatFIx----787878781", cidrRange)
			queries = append(queries, cidrRange...)
		}
	} else {
		fmt.Print("threatFIx----787878782", ip)
		queries = append(queries, ip)
	}
	return
}

// IPv4RangeToCIDRRange Convert IPv4 range into CIDR
func IPv4RangeToCIDRRange(ipStart string, ipEnd string) (cidrs []string, err error) {
	cidr2mask := []uint32{
		0x00000000, 0x80000000, 0xC0000000,
		0xE0000000, 0xF0000000, 0xF8000000,
		0xFC000000, 0xFE000000, 0xFF000000,
		0xFF800000, 0xFFC00000, 0xFFE00000,
		0xFFF00000, 0xFFF80000, 0xFFFC0000,
		0xFFFE0000, 0xFFFF0000, 0xFFFF8000,
		0xFFFFC000, 0xFFFFE000, 0xFFFFF000,
		0xFFFFF800, 0xFFFFFC00, 0xFFFFFE00,
		0xFFFFFF00, 0xFFFFFF80, 0xFFFFFFC0,
		0xFFFFFFE0, 0xFFFFFFF0, 0xFFFFFFF8,
		0xFFFFFFFC, 0xFFFFFFFE, 0xFFFFFFFF,
	}

	ipStartUint32 := iPv4ToUint32(ipStart)
	ipEndUint32 := iPv4ToUint32(ipEnd)

	if ipStartUint32 > ipEndUint32 {
		return nil, fmt.Errorf("start IP:%s must be less than end IP:%s", ipStart, ipEnd)
	}

	for ipEndUint32 >= ipStartUint32 {
		maxSize := 32
		for maxSize > 0 {

			maskedBase := ipStartUint32 & cidr2mask[maxSize-1]

			if maskedBase != ipStartUint32 {
				break
			}
			maxSize--

		}

		x := math.Log(float64(ipEndUint32-ipStartUint32+1)) / math.Log(2)
		maxDiff := 32 - int(math.Floor(x))
		if maxSize < maxDiff {
			maxSize = maxDiff
		}

		cidrs = append(cidrs, uInt32ToIPv4(ipStartUint32)+"/"+strconv.Itoa(maxSize))

		ipStartUint32 += uint32(math.Exp2(float64(32 - maxSize)))
	}

	return cidrs, err
}

// iPv4ToUint32 Convert IPv4 to uint32
func iPv4ToUint32(iPv4 string) uint32 {

	ipOctets := [4]uint64{}

	for i, v := range strings.SplitN(iPv4, ".", 4) {
		ipOctets[i], _ = strconv.ParseUint(v, 10, 32)
	}

	result := (ipOctets[0] << 24) | (ipOctets[1] << 16) | (ipOctets[2] << 8) | ipOctets[3]

	return uint32(result)
}

// uInt32ToIPv4 Convert uint32 to IP
func uInt32ToIPv4(iPuInt32 uint32) (iP string) {
	iP = fmt.Sprintf("%d.%d.%d.%d",
		iPuInt32>>24,
		(iPuInt32&0x00FFFFFF)>>16,
		(iPuInt32&0x0000FFFF)>>8,
		iPuInt32&0x000000FF)
	return iP
}

func CIDRs2IPs(cidrs []string) []string {
	var r []string
	for _, v := range cidrs {
		c, err := cidr.ParseCIDR(v)
		if err != nil {
			continue
		}
		if err := c.ForEachIP(func(ip string) error {
			r = append(r, ip)
			return nil
		}); err != nil {
			continue
		}
	}
	return r
}
