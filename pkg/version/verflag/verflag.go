// Package verflag defines utility functions to handle command line flags
// related to version of IAM.
package verflag

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rose839/IAM/pkg/version"
	flag "github.com/spf13/pflag"
)

type versionValue int

// Define const.
const (
	VersionFalse versionValue = 0
	VersionTrue  versionValue = 1
	VersionRaw   versionValue = 2
)

const (
	// Equal to "--version=raw"
	strRawVersion string = "raw"

	// Define version flag type to "version"
	strVersionFlagType string = "version"

	// Define version flag name to "--version"
	strVersionFlagName string = "version"
)

// Implemet pflag.Value interface.
func (v *versionValue) Set(s string) error {
	if s == strRawVersion {
		*v = VersionRaw
		return nil
	}

	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = VersionTrue
	} else {
		*v = VersionFalse
	}

	return err
}

func (v *versionValue) String() string {
	if *v == VersionRaw {
		return strRawVersion
	}

	return fmt.Sprintf("%v", bool(*v == VersionTrue))
}

// The type of the flag as required by the pflag.Value interface.
func (v *versionValue) Type() string {
	return strVersionFlagType
}

// VersionVar defines a flag with the specified name and usage string.
func VersionVar(p *versionValue, name string, value versionValue, usage string) {
	*p = value
	flag.Var(p, name, usage)

	// "--version" will be treated as "--version=true"
	flag.Lookup(name).NoOptDefVal = "true"
}

// Version wraps the VersionVar function.
func Version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)
	VersionVar(p, name, value, usage)
	return p
}

var versionFlag = Version(strVersionFlagName, VersionFalse, "Print version information and quit.")

// AddFlags registers this package's flags on arbitrary FlagSets, such that they point to the
// same value as the global flags.
func AddFlags(fs *flag.FlagSet) {
	fs.AddFlag(flag.Lookup(strVersionFlagName))
}

// PrintAndExitIfRequested will check if the -version flag was passed
// and, if so, print the version and exit.
func PrintAndExitIfRequested() {
	if *versionFlag == VersionRaw {
		fmt.Printf("%#v\n", version.Get())
		os.Exit(0)
	} else if *versionFlag == VersionTrue {
		fmt.Printf("%s\n", version.Get())
		os.Exit(0)
	}
}
