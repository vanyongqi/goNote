package main

/*
func IpContent(ip string) (queries []elastic.Query) {
	//step1: 将一串ip解析成单个ip，
	//10.10.14-16.*\n10.10.11.20-100\n10.10.10.10/24\n10.10.10.10
	//==> [10.10.14-16.* , 10.10.11.20-100 , 10.10.10.10/24，10.10.10.10]
	var ipBlock []string
	switch ip != "" {
	case strings.IndexByte(ip, ',') > 0:
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
	//var ipBlocks []string
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
			//ipBlocksPart, _ := cidrx.IPv4RangeToCIDRRange(ipStart, ipEnd)
			//ipBlocks = append(ipBlocks, ipBlocksPart...)
			queries = append(queries, elastic.NewRangeQuery("ip").Gte(ipStart).Lte(ipEnd))

		case strings.IndexByte(ipBlock[i], '-') > 0 && strings.IndexByte(ipBlock[i], '*') > 0: //10.10.14-16.*     //CIDR
			split := strings.Split(ipBlock[i], ".") //10.10.14-16.* ==> 10 10 1-11 *
			tem := strings.Split(split[2], "-")     //1 11

			ipStart := split[0] + "." + split[1] + "." + tem[0] + "." + strconv.Itoa(0)
			ipEnd := split[0] + "." + split[1] + "." + tem[1] + "." + strconv.Itoa(255)
			//ipBlocksPart, _ := cidrx.IPv4RangeToCIDRRange(ipStart, ipEnd) //10.10.14.0, 10.10.16.255
			//ipBlocks = append(ipBlocks, ipBlocksPart...)
			queries = append(queries, elastic.NewRangeQuery("ip").Gte(ipStart).Lte(ipEnd))

		case strings.IndexByte(ipBlock[i], '/') > 0: // 10.10.10.10/24
			split := strings.Split(ipBlock[i], "/")
			s1 := strings.Split(split[0], ".") //i= [10 ,10,10,10]
			s1[len(s1)-1] = strconv.Itoa(255)  //i=[10,10,10,20]
			queries = append(queries, elastic.NewTermQuery("ip", ipBlock[i]))

		default:
			queries = append(queries, elastic.NewTermQuery("ip", ipBlock[i]))

		}

	}

	return
}



*/
