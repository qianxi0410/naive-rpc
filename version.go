package naiverpc

import "fmt"

// the struct of rpc verison
type Version struct {
	Major, Minor, Patch string
	Metadata            string
}

// naive rpc current version
var CurrentVersion = Version{
	Major:    "0",
	Minor:    "0",
	Patch:    "1",
	Metadata: "dev",
}

func (v Version) String() string {
	version := fmt.Sprintf("Naive RPC Version: %s.%s.%s", v.Major, v.Minor, v.Patch)
	if v.Metadata != "" {
		version += "-" + v.Metadata
	}

	return version
}
