package service

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"os"
	"path"
)

// ToDO: what is fprocess in openfaas
func prepareEnv(reqEnvVars map[string]string) []string { // envProcess string,
	envs := []string{}
	//fprocessFound := false
	//fprocess := "fprocess=" + envProcess
	//if len(envProcess) > 0 {
	//	fprocessFound = true
	//}

	for k, v := range reqEnvVars {
		//if k == "fprocess" {
		//	fprocessFound = true
		//	fprocess = v
		//} else {
		envs = append(envs, k+"="+v)
		//}
	}
	//if fprocessFound {
	//	envs = append(envs, fprocess)
	//}
	return envs
}

// getOSMounts provides a mount for os-specific files such
// as the hosts file and resolv.conf
func getOSMounts() []specs.Mount {
	// Prior to hosts_dir env-var, this value was set to
	// os.Getwd()
	hostsDir := HostDir
	if v, ok := os.LookupEnv("hosts_dir"); ok && len(v) > 0 {
		hostsDir = v
	}

	mounts := []specs.Mount{}
	mounts = append(mounts, specs.Mount{
		Destination: "/etc/resolv.conf",
		Type:        "bind",
		Source:      path.Join(hostsDir, "resolv.conf"),
		Options:     []string{"rbind", "ro"},
	})

	mounts = append(mounts, specs.Mount{
		Destination: "/etc/hosts",
		Type:        "bind",
		Source:      path.Join(hostsDir, "hosts"),
		Options:     []string{"rbind", "ro"},
	})
	return mounts
}

//func createMounts(volumeMounts  []specs.Mount) error {
//	for _, volumeMount := range volumeMounts {
//
//	}
//}
