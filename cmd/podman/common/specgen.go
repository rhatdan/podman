package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/containers/image/v5/manifest"
	"github.com/containers/libpod/cmd/podman/parse"
	"github.com/containers/libpod/libpod/define"
	ann "github.com/containers/libpod/pkg/annotations"
	envLib "github.com/containers/libpod/pkg/env"
	ns "github.com/containers/libpod/pkg/namespaces"
	"github.com/containers/libpod/pkg/specgen"
	systemdGen "github.com/containers/libpod/pkg/systemd/generate"
	"github.com/containers/libpod/pkg/util"
	"github.com/docker/go-units"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
)

func getCPULimits(s *specgen.SpecGenerator, c *ContainerCLIOpts, args []string) (*specs.LinuxCPU, error) {
	cpu := &specs.LinuxCPU{}
	hasLimits := false

	if c.CPUShares > 0 {
		cpu.Shares = &c.CPUShares
		hasLimits = true
	}
	if c.CPUPeriod > 0 {
		cpu.Period = &c.CPUPeriod
		hasLimits = true
	}
	if c.CPUSetCPUs != "" {
		cpu.Cpus = c.CPUSetCPUs
		hasLimits = true
	}
	if c.CPUSetMems != "" {
		cpu.Mems = c.CPUSetMems
		hasLimits = true
	}
	if c.CPUQuota > 0 {
		cpu.Quota = &c.CPUQuota
		hasLimits = true
	}
	if c.CPURTPeriod > 0 {
		cpu.RealtimePeriod = &c.CPURTPeriod
		hasLimits = true
	}
	if c.CPURTRuntime > 0 {
		cpu.RealtimeRuntime = &c.CPURTRuntime
		hasLimits = true
	}

	if !hasLimits {
		return nil, nil
	}
	return cpu, nil
}

func getIOLimits(s *specgen.SpecGenerator, c *ContainerCLIOpts, args []string) (*specs.LinuxBlockIO, error) {
	var err error
	io := &specs.LinuxBlockIO{}
	hasLimits := false
	if b := c.BlkIOWeight; len(b) > 0 {
		u, err := strconv.ParseUint(b, 10, 16)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid value for blkio-weight")
		}
		nu := uint16(u)
		io.Weight = &nu
		hasLimits = true
	}

	if len(c.BlkIOWeightDevice) > 0 {
		if err := parseWeightDevices(c.BlkIOWeightDevice, s); err != nil {
			return nil, err
		}
		hasLimits = true
	}

	if bps := c.DeviceReadBPs; len(bps) > 0 {
		if s.ThrottleReadBpsDevice, err = parseThrottleBPSDevices(bps); err != nil {
			return nil, err
		}
		hasLimits = true
	}

	if bps := c.DeviceWriteBPs; len(bps) > 0 {
		if s.ThrottleWriteBpsDevice, err = parseThrottleBPSDevices(bps); err != nil {
			return nil, err
		}
		hasLimits = true
	}

	if iops := c.DeviceReadIOPs; len(iops) > 0 {
		if s.ThrottleReadIOPSDevice, err = parseThrottleIOPsDevices(iops); err != nil {
			return nil, err
		}
		hasLimits = true
	}

	if iops := c.DeviceWriteIOPs; len(iops) > 0 {
		if s.ThrottleWriteIOPSDevice, err = parseThrottleIOPsDevices(iops); err != nil {
			return nil, err
		}
		hasLimits = true
	}

	if !hasLimits {
		return nil, nil
	}
	return io, nil
}

func getPidsLimits(s *specgen.SpecGenerator, c *ContainerCLIOpts, args []string) (*specs.LinuxPids, error) {
	pids := &specs.LinuxPids{}
	hasLimits := false
	if c.CGroupsMode == "disabled" && c.PIDsLimit > 0 {
		return nil, nil
	}
	if c.PIDsLimit > 0 {
		pids.Limit = c.PIDsLimit
		hasLimits = true
	}
	if !hasLimits {
		return nil, nil
	}
	return pids, nil
}

func getMemoryLimits(s *specgen.SpecGenerator, c *ContainerCLIOpts, args []string) (*specs.LinuxMemory, error) {
	var err error
	memory := &specs.LinuxMemory{}
	hasLimits := false
	if m := c.Memory; len(m) > 0 {
		ml, err := units.RAMInBytes(m)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid value for memory")
		}
		memory.Limit = &ml
		hasLimits = true
	}
	if m := c.MemoryReservation; len(m) > 0 {
		mr, err := units.RAMInBytes(m)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid value for memory")
		}
		memory.Reservation = &mr
		hasLimits = true
	}
	if m := c.MemorySwap; len(m) > 0 {
		var ms int64
		if m == "-1" {
			ms = int64(-1)
			s.ResourceLimits.Memory.Swap = &ms
		} else {
			ms, err = units.RAMInBytes(m)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid value for memory")
			}
		}
		memory.Swap = &ms
		hasLimits = true
	}
	if m := c.KernelMemory; len(m) > 0 {
		mk, err := units.RAMInBytes(m)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid value for kernel-memory")
		}
		memory.Kernel = &mk
		hasLimits = true
	}
	if c.MemorySwappiness >= 0 {
		swappiness := uint64(c.MemorySwappiness)
		memory.Swappiness = &swappiness
		hasLimits = true
	}
	if c.OOMKillDisable {
		memory.DisableOOMKiller = &c.OOMKillDisable
		hasLimits = true
	}
	if !hasLimits {
		return nil, nil
	}
	return memory, nil
}

func FillOutSpecGen(s *specgen.SpecGenerator, c *ContainerCLIOpts, args []string) error {
	var (
		err error
		//namespaces map[string]string
	)

	// validate flags as needed
	if err := c.validate(); err != nil {
		return nil
	}

	s.User = c.User
	inputCommand := args[1:]
	if len(c.HealthCmd) > 0 {
		if c.NoHealthCheck {
			return errors.New("Cannot specify both --no-healthcheck and --health-cmd")
		}
		s.HealthConfig, err = makeHealthCheckFromCli(c.HealthCmd, c.HealthInterval, c.HealthRetries, c.HealthTimeout, c.HealthStartPeriod)
		if err != nil {
			return err
		}
	} else if c.NoHealthCheck {
		s.HealthConfig = &manifest.Schema2HealthConfig{
			Test: []string{"NONE"},
		}
	}

	userNS := ns.UsernsMode(c.UserNS)
	s.IDMappings, err = util.ParseIDMapping(userNS, c.UIDMap, c.GIDMap, c.SubUIDName, c.SubGIDName)
	if err != nil {
		return err
	}
	// If some mappings are specified, assume a private user namespace
	if userNS.IsDefaultValue() && (!s.IDMappings.HostUIDMapping || !s.IDMappings.HostGIDMapping) {
		s.UserNS.NSMode = specgen.Private
	}

	s.Terminal = c.TTY

	if err := verifyExpose(c.Expose); err != nil {
		return err
	}
	// We are not handling the Expose flag yet.
	// s.PortsExpose = c.Expose
	s.PortMappings = c.Net.PublishPorts
	s.PublishImagePorts = c.PublishAll
	s.Pod = c.Pod

	for k, v := range map[string]*specgen.Namespace{
		c.IPC:       &s.IpcNS,
		c.PID:       &s.PidNS,
		c.UTS:       &s.UtsNS,
		c.CGroupsNS: &s.CgroupNS,
	} {
		if k != "" {
			*v, err = specgen.ParseNamespace(k)
			if err != nil {
				return err
			}
		}
	}
	// userns must be treated differently
	if c.UserNS != "" {
		s.UserNS, err = specgen.ParseUserNamespace(c.UserNS)
		if err != nil {
			return err
		}
	}
	if c.Net != nil {
		s.NetNS = c.Net.Network
	}

	// STOP SIGNAL
	signalString := "TERM"
	if sig := c.StopSignal; len(sig) > 0 {
		signalString = sig
	}
	stopSignal, err := util.ParseSignal(signalString)
	if err != nil {
		return err
	}
	s.StopSignal = &stopSignal

	// ENVIRONMENT VARIABLES
	//
	// Precedence order (higher index wins):
	//  1) env-host, 2) image data, 3) env-file, 4) env
	env := map[string]string{
		"container": "podman",
	}

	// First transform the os env into a map. We need it for the labels later in
	// any case.
	osEnv, err := envLib.ParseSlice(os.Environ())
	if err != nil {
		return errors.Wrap(err, "error parsing host environment variables")
	}

	if c.EnvHost {
		env = envLib.Join(env, osEnv)
	} else if c.HTTPProxy {
		for _, envSpec := range []string{
			"http_proxy",
			"HTTP_PROXY",
			"https_proxy",
			"HTTPS_PROXY",
			"ftp_proxy",
			"FTP_PROXY",
			"no_proxy",
			"NO_PROXY",
		} {
			if v, ok := osEnv[envSpec]; ok {
				env[envSpec] = v
			}
		}
	}

	// env-file overrides any previous variables
	for _, f := range c.EnvFile {
		fileEnv, err := envLib.ParseFile(f)
		if err != nil {
			return err
		}
		// File env is overridden by env.
		env = envLib.Join(env, fileEnv)
	}

	// env overrides any previous variables
	if cmdLineEnv := c.env; len(cmdLineEnv) > 0 {
		parsedEnv, err := envLib.ParseSlice(cmdLineEnv)
		if err != nil {
			return err
		}
		env = envLib.Join(env, parsedEnv)
	}
	s.Env = env

	// LABEL VARIABLES
	labels, err := parse.GetAllLabels(c.LabelFile, c.Label)
	if err != nil {
		return errors.Wrapf(err, "unable to process labels")
	}

	if systemdUnit, exists := osEnv[systemdGen.EnvVariable]; exists {
		labels[systemdGen.EnvVariable] = systemdUnit
	}

	s.Labels = labels

	// ANNOTATIONS
	annotations := make(map[string]string)

	// First, add our default annotations
	annotations[ann.TTY] = "false"
	if c.TTY {
		annotations[ann.TTY] = "true"
	}

	// Last, add user annotations
	for _, annotation := range c.Annotation {
		splitAnnotation := strings.SplitN(annotation, "=", 2)
		if len(splitAnnotation) < 2 {
			return errors.Errorf("Annotations must be formatted KEY=VALUE")
		}
		annotations[splitAnnotation[0]] = splitAnnotation[1]
	}
	s.Annotations = annotations

	workDir := "/"
	if wd := c.Workdir; len(wd) > 0 {
		workDir = wd
	}
	s.WorkDir = workDir
	entrypoint := []string{}
	userCommand := []string{}
	if ep := c.Entrypoint; len(ep) > 0 {
		// Check if entrypoint specified is json
		if err := json.Unmarshal([]byte(c.Entrypoint), &entrypoint); err != nil {
			entrypoint = append(entrypoint, ep)
		}
	}

	var command []string

	s.Entrypoint = entrypoint

	// Build the command
	// If we have an entry point, it goes first
	if len(entrypoint) > 0 {
		command = entrypoint
	}
	if len(inputCommand) > 0 {
		// User command overrides data CMD
		command = append(command, inputCommand...)
		userCommand = append(userCommand, inputCommand...)
	}

	if len(inputCommand) > 0 {
		s.Command = userCommand
	} else {
		s.Command = command
	}

	// SHM Size
	if c.ShmSize != "" {
		shmSize, err := units.FromHumanSize(c.ShmSize)
		if err != nil {
			return errors.Wrapf(err, "unable to translate --shm-size")
		}
		s.ShmSize = &shmSize
	}
	s.HostAdd = c.Net.AddHosts
	s.UseImageResolvConf = c.Net.UseImageResolvConf
	s.DNSServers = c.Net.DNSServers
	s.DNSSearch = c.Net.DNSSearch
	s.DNSOptions = c.Net.DNSOptions
	s.StaticIP = c.Net.StaticIP
	s.StaticMAC = c.Net.StaticMAC
	s.UseImageHosts = c.Net.NoHosts

	s.ImageVolumeMode = c.ImageVolume
	if s.ImageVolumeMode == "bind" {
		s.ImageVolumeMode = "anonymous"
	}

	systemd := c.SystemdD == "always"
	if !systemd && command != nil {
		x, err := strconv.ParseBool(c.SystemdD)
		if err != nil {
			return errors.Wrapf(err, "cannot parse bool %s", c.SystemdD)
		}
		if x && (command[0] == "/usr/sbin/init" || command[0] == "/sbin/init" || (filepath.Base(command[0]) == "systemd")) {
			systemd = true
		}
	}
	if systemd {
		if s.StopSignal == nil {
			stopSignal, err = util.ParseSignal("RTMIN+3")
			if err != nil {
				return errors.Wrapf(err, "error parsing systemd signal")
			}
			s.StopSignal = &stopSignal
		}
	}
	if s.ResourceLimits == nil {
		s.ResourceLimits = &specs.LinuxResources{}
	}
	s.ResourceLimits.Memory, err = getMemoryLimits(s, c, args)
	if err != nil {
		return err
	}
	s.ResourceLimits.BlockIO, err = getIOLimits(s, c, args)
	if err != nil {
		return err
	}
	s.ResourceLimits.Pids, err = getPidsLimits(s, c, args)
	if err != nil {
		return err
	}
	s.ResourceLimits.CPU, err = getCPULimits(s, c, args)
	if err != nil {
		return err
	}
	if s.ResourceLimits.CPU == nil && s.ResourceLimits.Pids == nil && s.ResourceLimits.BlockIO == nil && s.ResourceLimits.Memory == nil {
		s.ResourceLimits = nil
	}

	if s.LogConfiguration == nil {
		s.LogConfiguration = &specgen.LogConfig{}
	}
	s.LogConfiguration.Driver = define.KubernetesLogging
	if ld := c.LogDriver; len(ld) > 0 {
		s.LogConfiguration.Driver = ld
	}
	s.CgroupParent = c.CGroupParent
	s.CgroupsMode = c.CGroupsMode
	s.Groups = c.GroupAdd
	// TODO WTF
	//cgroup := &cc.CgroupConfig{
	//	Cgroupns:     c.String("cgroupns"),
	//}
	//
	//userns := &cc.UserConfig{
	//	GroupAdd:   c.StringSlice("group-add"),
	//	IDMappings: idmappings,
	//	UsernsMode: usernsMode,
	//	User:       user,
	//}
	//
	//uts := &cc.UtsConfig{
	//	UtsMode:  utsMode,
	//	NoHosts:  c.Bool("no-hosts"),
	//	HostAdd:  c.StringSlice("add-host"),
	//	Hostname: c.String("hostname"),
	//}

	s.Hostname = c.Hostname
	sysctl := map[string]string{}
	if ctl := c.Sysctl; len(ctl) > 0 {
		sysctl, err = util.ValidateSysctls(ctl)
		if err != nil {
			return err
		}
	}
	s.Sysctl = sysctl

	s.CapAdd = c.CapAdd
	s.CapDrop = c.CapDrop
	s.Privileged = c.Privileged
	s.ReadOnlyFilesystem = c.ReadOnly

	// TODO
	// ouitside of specgen and oci though
	// defaults to true, check spec/storage
	//s.readon = c.ReadOnlyTmpFS
	//  TODO convert to map?
	// check if key=value and convert
	sysmap := make(map[string]string)
	for _, ctl := range c.Sysctl {
		splitCtl := strings.SplitN(ctl, "=", 2)
		if len(splitCtl) < 2 {
			return errors.Errorf("invalid sysctl value %q", ctl)
		}
		sysmap[splitCtl[0]] = splitCtl[1]
	}
	s.Sysctl = sysmap

	for _, opt := range c.SecurityOpt {
		if opt == "no-new-privileges" {
			s.ContainerSecurityConfig.NoNewPrivileges = true
		} else {
			con := strings.SplitN(opt, "=", 2)
			if len(con) != 2 {
				return fmt.Errorf("invalid --security-opt 1: %q", opt)
			}

			switch con[0] {
			case "label":
				// TODO selinux opts and label opts are the same thing
				s.ContainerSecurityConfig.SelinuxOpts = append(s.ContainerSecurityConfig.SelinuxOpts, con[1])
			case "apparmor":
				s.ContainerSecurityConfig.ApparmorProfile = con[1]
			case "seccomp":
				s.SeccompProfilePath = con[1]
			default:
				return fmt.Errorf("invalid --security-opt 2: %q", opt)
			}
		}
	}

	s.SeccompPolicy = c.SeccompPolicy

	// TODO: should parse out options
	s.VolumesFrom = c.VolumesFrom

	// Only add read-only tmpfs mounts in case that we are read-only and the
	// read-only tmpfs flag has been set.
	mounts, volumes, err := parseVolumes(c.Volume, c.Mount, c.TmpFS, (c.ReadOnlyTmpFS && c.ReadOnly))
	if err != nil {
		return err
	}
	s.Mounts = mounts
	s.Volumes = volumes

	// TODO any idea why this was done
	//devices := rtc.Containers.Devices
	// TODO conflict on populate?
	//
	//if c.Changed("device") {
	//	devices = append(devices, c.StringSlice("device")...)
	//}

	for _, dev := range c.Devices {
		s.Devices = append(s.Devices, specs.LinuxDevice{Path: dev})
	}

	// TODO things i cannot find in spec
	// we dont think these are in the spec
	// init - initbinary
	// initpath
	s.Stdin = c.Interactive
	// quiet
	//DeviceCgroupRules: c.StringSlice("device-cgroup-rule"),

	// Rlimits/Ulimits
	for _, u := range c.Ulimit {
		if u == "host" {
			s.Rlimits = nil
			break
		}
		ul, err := units.ParseUlimit(u)
		if err != nil {
			return errors.Wrapf(err, "ulimit option %q requires name=SOFT:HARD, failed to be parsed", u)
		}
		rl := specs.POSIXRlimit{
			Type: ul.Name,
			Hard: uint64(ul.Hard),
			Soft: uint64(ul.Soft),
		}
		s.Rlimits = append(s.Rlimits, rl)
	}

	//Tmpfs:         c.StringArray("tmpfs"),

	// TODO how to handle this?
	//Syslog:        c.Bool("syslog"),

	logOpts := make(map[string]string)
	for _, o := range c.LogOptions {
		split := strings.SplitN(o, "=", 2)
		if len(split) < 2 {
			return errors.Errorf("invalid log option %q", o)
		}
		switch {
		case split[0] == "driver":
			s.LogConfiguration.Driver = split[1]
		case split[0] == "path":
			s.LogConfiguration.Path = split[1]
		default:
			logOpts[split[0]] = split[1]
		}
	}
	s.LogConfiguration.Options = logOpts
	s.Name = c.Name

	s.OOMScoreAdj = &c.OOMScoreAdj
	s.RestartPolicy = c.Restart
	s.Remove = c.Rm
	s.StopTimeout = &c.StopTimeout

	// TODO where should we do this?
	//func verifyContainerResources(config *cc.CreateConfig, update bool) ([]string, error) {
	return nil
}

func makeHealthCheckFromCli(inCmd, interval string, retries uint, timeout, startPeriod string) (*manifest.Schema2HealthConfig, error) {
	// Every healthcheck requires a command
	if len(inCmd) == 0 {
		return nil, errors.New("Must define a healthcheck command for all healthchecks")
	}

	// first try to parse option value as JSON array of strings...
	cmd := []string{}

	if inCmd == "none" {
		cmd = []string{"NONE"}
	} else {
		err := json.Unmarshal([]byte(inCmd), &cmd)
		if err != nil {
			// ...otherwise pass it to "/bin/sh -c" inside the container
			cmd = []string{"CMD-SHELL", inCmd}
		}
	}
	hc := manifest.Schema2HealthConfig{
		Test: cmd,
	}

	if interval == "disable" {
		interval = "0"
	}
	intervalDuration, err := time.ParseDuration(interval)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid healthcheck-interval %s ", interval)
	}

	hc.Interval = intervalDuration

	if retries < 1 {
		return nil, errors.New("healthcheck-retries must be greater than 0.")
	}
	hc.Retries = int(retries)
	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid healthcheck-timeout %s", timeout)
	}
	if timeoutDuration < time.Duration(1) {
		return nil, errors.New("healthcheck-timeout must be at least 1 second")
	}
	hc.Timeout = timeoutDuration

	startPeriodDuration, err := time.ParseDuration(startPeriod)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid healthcheck-start-period %s", startPeriod)
	}
	if startPeriodDuration < time.Duration(0) {
		return nil, errors.New("healthcheck-start-period must be 0 seconds or greater")
	}
	hc.StartPeriod = startPeriodDuration

	return &hc, nil
}

func parseWeightDevices(weightDevs []string, s *specgen.SpecGenerator) error {
	for _, val := range weightDevs {
		split := strings.SplitN(val, ":", 2)
		if len(split) != 2 {
			return fmt.Errorf("bad format: %s", val)
		}
		if !strings.HasPrefix(split[0], "/dev/") {
			return fmt.Errorf("bad format for device path: %s", val)
		}
		weight, err := strconv.ParseUint(split[1], 10, 0)
		if err != nil {
			return fmt.Errorf("invalid weight for device: %s", val)
		}
		if weight > 0 && (weight < 10 || weight > 1000) {
			return fmt.Errorf("invalid weight for device: %s", val)
		}
		w := uint16(weight)
		s.WeightDevice[split[0]] = specs.LinuxWeightDevice{
			Weight:     &w,
			LeafWeight: nil,
		}
	}
	return nil
}

func parseThrottleBPSDevices(bpsDevices []string) (map[string]specs.LinuxThrottleDevice, error) {
	td := make(map[string]specs.LinuxThrottleDevice)
	for _, val := range bpsDevices {
		split := strings.SplitN(val, ":", 2)
		if len(split) != 2 {
			return nil, fmt.Errorf("bad format: %s", val)
		}
		if !strings.HasPrefix(split[0], "/dev/") {
			return nil, fmt.Errorf("bad format for device path: %s", val)
		}
		rate, err := units.RAMInBytes(split[1])
		if err != nil {
			return nil, fmt.Errorf("invalid rate for device: %s. The correct format is <device-path>:<number>[<unit>]. Number must be a positive integer. Unit is optional and can be kb, mb, or gb", val)
		}
		if rate < 0 {
			return nil, fmt.Errorf("invalid rate for device: %s. The correct format is <device-path>:<number>[<unit>]. Number must be a positive integer. Unit is optional and can be kb, mb, or gb", val)
		}
		td[split[0]] = specs.LinuxThrottleDevice{Rate: uint64(rate)}
	}
	return td, nil
}

func parseThrottleIOPsDevices(iopsDevices []string) (map[string]specs.LinuxThrottleDevice, error) {
	td := make(map[string]specs.LinuxThrottleDevice)
	for _, val := range iopsDevices {
		split := strings.SplitN(val, ":", 2)
		if len(split) != 2 {
			return nil, fmt.Errorf("bad format: %s", val)
		}
		if !strings.HasPrefix(split[0], "/dev/") {
			return nil, fmt.Errorf("bad format for device path: %s", val)
		}
		rate, err := strconv.ParseUint(split[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid rate for device: %s. The correct format is <device-path>:<number>. Number must be a positive integer", val)
		}
		td[split[0]] = specs.LinuxThrottleDevice{Rate: rate}
	}
	return td, nil
}
