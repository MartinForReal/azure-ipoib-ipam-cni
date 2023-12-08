// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	type100 "github.com/containernetworking/cni/pkg/types/100"
	"github.com/containernetworking/cni/pkg/version"
	bv "github.com/containernetworking/plugins/pkg/utils/buildversion"
	"github.com/vishvananda/netlink"

	"github.com/Azure/azure-ipoib-ipam-cni/pkg/ipoibconfig"
)

const (
	IPAMPluginName = "azure-ipoib-ipam-cni"
	KVPStorePath   = "/var/lib/hyperv/.kvp_pool_0"
)

func main() {
	skel.PluginMain(cmdAdd, cmdCheck, cmdDel, version.All, bv.BuildString(IPAMPluginName))
}

func cmdAdd(args *skel.CmdArgs) error {
	config := &types.NetConf{}
	err := json.Unmarshal(args.StdinData, config)
	if err != nil {
		return err
	}
	if config.IPAM.Type != IPAMPluginName {
		return fmt.Errorf("unsupported IPAM type: %s", config.IPAM.Type)
	}
	content, err := os.ReadFile(KVPStorePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v", KVPStorePath, err)
	}

	ibLink, err := netlink.LinkByName(args.IfName)
	if err != nil {
		return err
	}

	macAddr := ibLink.Attrs().HardwareAddr
	fmt.Printf("MAC address: %s", macAddr.String())

	ipAddr, err := ipoibconfig.GetIPoIBAddr(content, macAddr.String())
	if err != nil {
		return err
	}
	if ipAddr == nil {
		return fmt.Errorf("IPoIB address not found")
	}
	result := &type100.Result{
		CNIVersion: type100.ImplementedSpecVersion,
		IPs: []*type100.IPConfig{
			{
				Address: *ipAddr,
			},
		},
	}
	// outputCmdArgs(args)
	return types.PrintResult(result, config.CNIVersion)
}

func cmdDel(args *skel.CmdArgs) error {
	// get cni config
	config := &types.NetConf{}
	err := json.Unmarshal(args.StdinData, config)
	if err != nil {
		return err
	}
	if config.IPAM.Type != IPAMPluginName {
		return fmt.Errorf("unsupported IPAM type: %s", config.IPAM.Type)
	}
	return types.PrintResult(&type100.Result{}, config.CNIVersion)
}

func cmdCheck(args *skel.CmdArgs) error {
	// get cni config
	config := &types.NetConf{}
	err := json.Unmarshal(args.StdinData, config)
	if err != nil {
		return err
	}
	if config.IPAM.Type != IPAMPluginName {
		return fmt.Errorf("unsupported IPAM type: %s", config.IPAM.Type)
	}
	return types.PrintResult(&type100.Result{}, config.CNIVersion)
}
