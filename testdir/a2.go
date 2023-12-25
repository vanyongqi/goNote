package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func main2() {
	IpContent("10.10.14-16.*\n10.10.11.20-100\n10.10.10.10/24\n10.10.10.10\n11.11.11.11")
}

func IpContent(ip string) (queries string) {
	//step1: 将一串ip解析成单个ip，
	//10.10.14-16.*\n10.10.11.20-100\n10.10.10.10/24\n10.10.10.10
	//==> [10.10.14-16.* , 10.10.11.20-100 , 10.10.10.10/24，10.10.10.10]
	var ipBlock []string
	switch ip != "" {
	case strings.IndexByte(ip, ',') > 0: //查找指定字节的位置。如果找到该字节，则返回其在字符串中的索引；如果没有找到，则返回 -1。
		split := strings.Split(ip, ",")
		for _, value := range split {
			ipBlock = append(ipBlock, value)
		}
	case strings.Contains(ip, "\r\n"):
		split := strings.Split(ip, "\r\n")

		for _, value := range split {
			ipBlock = append(ipBlock, value)
		}
	case strings.Contains(ip, "\n"):
		split := strings.Split(ip, "\n") //============>[10.10.14-16.* , 10.10.11.20-100 , 10.10.10.10/24, 10.10.10.10]
		for _, value := range split {    //
			ipBlock = append(ipBlock, value)
		}

	default:

		return
	}
	//ipBlock []string  [10.10.14-16.* , 10.10.11.20-100 , 10.10.10.10/24，10.10.10.10]
	//step2： 根据第四种情况进行判断分析
	var ipBlocks []string
	for i := 0; i < len(ipBlock); i++ {
		fmt.Println(ipBlock[i])
	}
	fmt.Println("------------------------------------")
	for i := 0; i < len(ipBlock); i++ {
		switch ipBlock[i] != " " {
		case strings.IndexByte(ipBlock[i], '-') > 0 && strings.IndexByte(ipBlock[i], '*') <= 0: //10.10.11.20-100  //CIDR
			split := strings.Split(ipBlock[i], "-") //10.10.11.20-100 ===> 10.10.10.20 split[1] = 100
			if len(split) != 2 {
				return
			}
			i := strings.Split(split[0], ".") //i= [10 ,10,10,10]
			i[len(i)-1] = split[1]            //i=[10,10,10,20]
			ipStart := split[0]
			ipEnd := strings.Join(i, ".")
			ipBlocksPart, _ := IPv4RangeToCIDRRange(ipStart, ipEnd)
			ipBlocks = append(ipBlocks, ipBlocksPart...)

		case strings.IndexByte(ipBlock[i], '-') > 0 && strings.IndexByte(ipBlock[i], '*') > 0: //10.10.14-16.*     //CIDR
			split := strings.Split(ipBlock[i], ".") //10.10.14-16.* ==> 10 10 1-11 *
			tem := strings.Split(split[2], "-")     //1 11

			ipStart := split[0] + "." + split[1] + "." + tem[0] + "." + strconv.Itoa(0)
			//fmt.Print("222", ipStart, "3333\n")
			ipEnd := split[0] + "." + split[1] + "." + tem[1] + "." + strconv.Itoa(255)
			ipBlocksPart, _ := IPv4RangeToCIDRRange(ipStart, ipEnd) //10.10.14.0, 10.10.16.255
			ipBlocks = append(ipBlocks, ipBlocksPart...)
		case strings.IndexByte(ipBlock[i], '/') > 0: // 10.10.10.10/24    //CIDR
			split := strings.Split(ipBlock[i], "/")
			i := strings.Split(split[0], ".") //i= [10,10,10,10]
			i[len(i)-1] = strconv.Itoa(255)   //i=[10,10,10,255]
			ipStart := split[0]
			ipEnd := strings.Join(i, ".")
			ipBlocksPart, _ := IPv4RangeToCIDRRange(ipStart, ipEnd)
			ipBlocks = append(ipBlocks, ipBlocksPart...)

		default:
			fmt.Print("----------")
			ipBlocksPart, _ := IPv4RangeToCIDRRange("11.11.11.10", ipBlocks[i])
			ipBlocks = append(ipBlocks, ipBlocksPart...)
		}
		for i := 0; i < len(ipBlocks); i++ {
			fmt.Println(ipBlocks[i])
		}

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
