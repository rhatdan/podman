package abi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	cniversion "github.com/containernetworking/cni/pkg/version"
	"github.com/containers/libpod/libpod"
	"github.com/containers/libpod/pkg/domain/entities"
	"github.com/containers/libpod/pkg/network"
	"github.com/containers/libpod/pkg/util"
	"github.com/pkg/errors"
)

func getCNIConfDir(r *libpod.Runtime) (string, error) {
	config, err := r.GetConfig()
	if err != nil {
		return "", err
	}
	configPath := config.Network.NetworkConfigDir

	if len(config.Network.NetworkConfigDir) < 1 {
		configPath = network.CNIConfigDir
	}
	return configPath, nil
}

func (ic *ContainerEngine) NetworkList(ctx context.Context, options entities.NetworkListOptions) ([]*entities.NetworkListReport, error) {
	var reports []*entities.NetworkListReport
	cniConfigPath, err := getCNIConfDir(ic.Libpod)
	if err != nil {
		return nil, err
	}
	networks, err := network.LoadCNIConfsFromDir(cniConfigPath)
	if err != nil {
		return nil, err
	}

	for _, n := range networks {
		reports = append(reports, &entities.NetworkListReport{NetworkConfigList: n})
	}
	return reports, nil
}

func (ic *ContainerEngine) NetworkInspect(ctx context.Context, namesOrIds []string, options entities.NetworkInspectOptions) ([]entities.NetworkInspectReport, error) {
	var (
		rawCNINetworks []entities.NetworkInspectReport
	)
	for _, name := range namesOrIds {
		rawList, err := network.InspectNetwork(name)
		if err != nil {
			return nil, err
		}
		rawCNINetworks = append(rawCNINetworks, rawList)
	}
	return rawCNINetworks, nil
}

func (ic *ContainerEngine) NetworkRm(ctx context.Context, namesOrIds []string, options entities.NetworkRmOptions) ([]*entities.NetworkRmReport, error) {
	var reports []*entities.NetworkRmReport
	for _, name := range namesOrIds {
		report := entities.NetworkRmReport{Name: name}
		containers, err := ic.Libpod.GetAllContainers()
		if err != nil {
			return reports, err
		}
		// We need to iterate containers looking to see if they belong to the given network
		for _, c := range containers {
			if util.StringInSlice(name, c.Config().Networks) {
				// if user passes force, we nuke containers
				if !options.Force {
					// Without the force option, we return an error
					return reports, errors.Errorf("%q has associated containers with it. Use -f to forcibly delete containers", name)
				}
				if err := ic.Libpod.RemoveContainer(ctx, c, true, true); err != nil {
					return reports, err
				}
			}
		}
		if err := network.RemoveNetwork(name); err != nil {
			report.Err = err
		}
		reports = append(reports, &report)
	}
	return reports, nil
}

func (ic *ContainerEngine) NetworkCreate(ctx context.Context, name string, options entities.NetworkCreateOptions) (*entities.NetworkCreateReport, error) {
	var (
		err      error
		fileName string
	)
	if len(options.MacVLAN) > 0 {
		fileName, err = createMacVLAN(ic.Libpod, name, options)
	} else {
		fileName, err = createBridge(ic.Libpod, name, options)
	}
	if err != nil {
		return nil, err
	}
	return &entities.NetworkCreateReport{Filename: fileName}, nil
}

// createBridge creates a CNI network
func createBridge(r *libpod.Runtime, name string, options entities.NetworkCreateOptions) (string, error) {
	isGateway := true
	ipMasq := true
	subnet := &options.Subnet
	ipRange := options.Range
	runtimeConfig, err := r.GetConfig()
	if err != nil {
		return "", err
	}
	// if range is provided, make sure it is "in" network
	if subnet.IP != nil {
		// if network is provided, does it conflict with existing CNI or live networks
		err = network.ValidateUserNetworkIsAvailable(subnet)
	} else {
		// if no network is provided, figure out network
		subnet, err = network.GetFreeNetwork()
	}
	if err != nil {
		return "", err
	}
	gateway := options.Gateway
	if gateway == nil {
		// if no gateway is provided, provide it as first ip of network
		gateway = network.CalcGatewayIP(subnet)
	}
	// if network is provided and if gateway is provided, make sure it is "in" network
	if options.Subnet.IP != nil && options.Gateway != nil {
		if !subnet.Contains(gateway) {
			return "", errors.Errorf("gateway %s is not in valid for subnet %s", gateway.String(), subnet.String())
		}
	}
	if options.Internal {
		isGateway = false
		ipMasq = false
	}

	// if a range is given, we need to ensure it is "in" the network range.
	if options.Range.IP != nil {
		if options.Subnet.IP == nil {
			return "", errors.New("you must define a subnet range to define an ip-range")
		}
		firstIP, err := network.FirstIPInSubnet(&options.Range)
		if err != nil {
			return "", err
		}
		lastIP, err := network.LastIPInSubnet(&options.Range)
		if err != nil {
			return "", err
		}
		if !subnet.Contains(firstIP) || !subnet.Contains(lastIP) {
			return "", errors.Errorf("the ip range %s does not fall within the subnet range %s", options.Range.String(), subnet.String())
		}
	}
	bridgeDeviceName, err := network.GetFreeDeviceName()
	if err != nil {
		return "", err
	}

	if len(name) > 0 {
		netNames, err := network.GetNetworkNamesFromFileSystem()
		if err != nil {
			return "", err
		}
		if util.StringInSlice(name, netNames) {
			return "", errors.Errorf("the network name %s is already used", name)
		}
	} else {
		// If no name is given, we give the name of the bridge device
		name = bridgeDeviceName
	}

	ncList := network.NewNcList(name, cniversion.Current())
	var plugins []network.CNIPlugins
	var routes []network.IPAMRoute

	defaultRoute, err := network.NewIPAMDefaultRoute()
	if err != nil {
		return "", err
	}
	routes = append(routes, defaultRoute)
	ipamConfig, err := network.NewIPAMHostLocalConf(subnet, routes, ipRange, gateway)
	if err != nil {
		return "", err
	}

	// TODO need to iron out the role of isDefaultGW and IPMasq
	bridge := network.NewHostLocalBridge(bridgeDeviceName, isGateway, false, ipMasq, ipamConfig)
	plugins = append(plugins, bridge)
	plugins = append(plugins, network.NewPortMapPlugin())
	plugins = append(plugins, network.NewFirewallPlugin())
	// if we find the dnsname plugin, we add configuration for it
	if network.HasDNSNamePlugin(runtimeConfig.Network.CNIPluginDirs) && !options.DisableDNS {
		// Note: in the future we might like to allow for dynamic domain names
		plugins = append(plugins, network.NewDNSNamePlugin(network.DefaultPodmanDomainName))
	}
	ncList["plugins"] = plugins
	b, err := json.MarshalIndent(ncList, "", "   ")
	if err != nil {
		return "", err
	}
	cniConfigPath, err := getCNIConfDir(r)
	if err != nil {
		return "", err
	}
	cniPathName := filepath.Join(cniConfigPath, fmt.Sprintf("%s.conflist", name))
	err = ioutil.WriteFile(cniPathName, b, 0644)
	return cniPathName, err
}

func createMacVLAN(r *libpod.Runtime, name string, options entities.NetworkCreateOptions) (string, error) {
	var (
		plugins []network.CNIPlugins
	)
	liveNetNames, err := network.GetLiveNetworkNames()
	if err != nil {
		return "", err
	}
	// Make sure the host-device exists
	if !util.StringInSlice(options.MacVLAN, liveNetNames) {
		return "", errors.Errorf("failed to find network interface %q", options.MacVLAN)
	}
	if len(name) > 0 {
		netNames, err := network.GetNetworkNamesFromFileSystem()
		if err != nil {
			return "", err
		}
		if util.StringInSlice(name, netNames) {
			return "", errors.Errorf("the network name %s is already used", name)
		}
	} else {
		name, err = network.GetFreeDeviceName()
		if err != nil {
			return "", err
		}
	}
	ncList := network.NewNcList(name, cniversion.Current())
	macvlan := network.NewMacVLANPlugin(options.MacVLAN)
	plugins = append(plugins, macvlan)
	ncList["plugins"] = plugins
	b, err := json.MarshalIndent(ncList, "", "   ")
	if err != nil {
		return "", err
	}
	cniConfigPath, err := getCNIConfDir(r)
	if err != nil {
		return "", err
	}
	cniPathName := filepath.Join(cniConfigPath, fmt.Sprintf("%s.conflist", name))
	err = ioutil.WriteFile(cniPathName, b, 0644)
	return cniPathName, err
}
