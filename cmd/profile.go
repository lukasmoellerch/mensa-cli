package cmd

import (
	"os"
	"runtime/pprof"
)

func profile() (func(), error) {
	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			return nil, err
		}
		pprof.StartCPUProfile(f)
		return func() { pprof.StopCPUProfile() }, nil
	}
	return func() {}, nil
}
