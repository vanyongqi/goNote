package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func main() {
	// connection.Connect()
	res := IsMultiNetworkSegment("10.10.10-30.*")
	fmt.Println(res)
}

// IsMultiNetworkSegment multi network segment
func IsMultiNetworkSegment(segment string) bool {
	suffix := strings.HasSuffix(segment, ".*") //10.10.10-30.*
	if suffix {
		address := strings.Split(segment, "-") //10.10.10  20.*
		if len(address) != 2 {
			return false
		}

		first := address[0]  //10.10.10
		second := address[1] //30.*

		firstAddress := strings.Split(first, ".") //[10,10,10]
		if len(firstAddress) != 3 {
			return false
		}

		secondAddress := strings.Split(second, ".") //[30,*]
		if len(secondAddress) != 2 {
			return false
		}

		start := firstAddress[len(firstAddress)-1] //[10]
		end := secondAddress[0]                    //[30]

		startNumber, err := strconv.Atoi(start) //10
		if err != nil {
			return false
		}

		endNumber, err := strconv.Atoi(end) //30
		if err != nil {
			return false
		}

		if startNumber > endNumber {
			return false
		}

		startIPAddr := fmt.Sprintf("%s.0", strings.Join(firstAddress, ".")) //10.10.10.0
		fmt.Println(startIPAddr)
		if net.ParseIP(startIPAddr) == nil {
			return false
		}

		endIPAddr := func(firstAddress []string, endNumber int) string {
			prefix := fmt.Sprintf("%s.%d",
				strings.Join(firstAddress[:len(firstAddress)-1], "."), //10.10.30.0
				endNumber,
			)
			fmt.Println(prefix + ".0")
			return fmt.Sprintf("%s.255", prefix) //10.10.30.0
		}(firstAddress, endNumber)
		if net.ParseIP(endIPAddr) == nil {
			return false
		}
		fmt.Println(startIPAddr + "-" + endIPAddr)
		return true
	}

	return false
}
