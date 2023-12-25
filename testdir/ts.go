package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main3() {
	//res := ConvertIpContent("10.10.14-16.*\n10.10.11.20-100\n10.10.10.10/24\n10.10.10.10\n11.11.11.11")
	res := ConvertIpContent("10.10.10.91,10.10.10.58")
	fmt.Println(res)
}

// ConvertIpContent 处理前端传的ip段信息
func ConvertIpContent(ip string) (queries []string) {
	if len(ip) == 0 {
		return
	}

	if strings.IndexByte(ip, ',') > 0 {
		split := strings.Split(ip, ",")
		if len(split) < 2 {
			return
		}
		for _, value := range split {
			queries = append(queries, ConvertIpContent(value)...)
		}
	} else if strings.Contains(ip, "\r\n") {
		split := strings.Split(ip, "\r\n")
		if len(split) < 2 {
			return
		}
		for _, value := range split {
			queries = append(queries, ConvertIpContent(value)...)
		}
	} else if strings.Contains(ip, "\n") {
		split := strings.Split(ip, "\n")
		if len(split) < 2 {
			return
		}
		for _, value := range split {
			queries = append(queries, ConvertIpContent(value)...)
		}
	} else if strings.IndexByte(ip, '-') > 0 {
		if strings.IndexByte(ip, '*') > 0 { //存在- 并且存在*
			//ip = "10.10.1-11.*"
			split := strings.Split(ip, ".")     //10 10 1-11 *
			tem := strings.Split(split[2], "-") //1 11
			var ipstart, ipend string
			ipstart = split[0] + "." + split[1] + "." + tem[0] + "." + strconv.Itoa(0)
			ipend = split[0] + "." + split[1] + "." + tem[1] + "." + strconv.Itoa(255)
			queries = append(queries, ipstart)
			queries = append(queries, ipend)
			//log.Info("Que:2^2:---->", ipstart, "   ", ipend, "   ", queries)
		} else { // 存在 - 但不存在*
			split := strings.Split(ip, "-") //split[0] = 10.10.10.10 split[1] = 20
			if len(split) != 2 {
				return
			}
			i := strings.Split(split[0], ".") //i= [10 ,10,10,10]
			i[len(i)-1] = split[1]            //i=[10,10,10,20]
			//NewRangeQuery(ip).Gte(10.10.10.10).Lte(10.10.10.20)
			queries = append(queries, split[0])
			queries = append(queries, strings.Join(i, "."))
		}
	} else {
		queries = append(queries, ip)
	}
	return
}
