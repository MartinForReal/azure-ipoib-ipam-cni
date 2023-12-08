package ipoibconfig

import (
	"fmt"
	"net"
	"strings"

	"github.com/Azure/hyperkv"
)

const (
	IPoIBDataKey = "IPoIB_Data"
)

func GetIPoIBAddr(rawContent []byte, macAddr string) (*net.IPNet, error) {
	list := hyperkv.Parse(rawContent)
	var ipoibData string
	for _, item := range list {
		if strings.EqualFold(item.Key, IPoIBDataKey) {
			fmt.Printf("Found IPoIB_Data: %s", item.Value)
			ipoibData = item.Value
		}
	}
	if ipoibData == "" {
		return nil, fmt.Errorf("IPoIB_Data not found in config file")
	}

	configMap, err := ParseIPOIBConfig([]byte(ipoibData))
	if err != nil {
		return nil, err
	}
	if configMap == nil {
		return nil, fmt.Errorf("invalid config map")
	}
	if ipAddr, ok := configMap[macAddr]; ok {
		return ipAddr, nil
	}
	return nil, nil
}
