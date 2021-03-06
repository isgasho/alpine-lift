package lift

import (
	"strconv"
)

// AlpineData is the main alpine-data yaml specification
type AlpineData struct {
	RootPasswd  string            `yaml:"password"`
	MOTD        string            `yaml:"motd"`
	Network     *NetworkSettings  `yaml:"network"`
	Packages    *PackagesConfig   `yaml:"packages"`
	DRP         *DRProvision      `yaml:"dr_provision"`
	SSHDConfig  *SSHD             `yaml:"sshd"`
	Groups      MultiString       `yaml:"groups"`
	Users       []User            `yaml:"users"`
	RunCMD      []MultiString     `yaml:"runcmd"`
	WriteFiles  []WriteFile       `yaml:"write_files"`
	TimeZone    string            `yaml:"timezone"`
	Keymap      string            `yaml:"keymap"`
	UnLift      bool              `yaml:"unlift"`
	ScratchDisk string            `yaml:"scratch_disk"`
	Disks       []Disk            `yaml:"disks"`
	MTA         *MTAConfiguration `yaml:"mta"`
}

// User specifies a specific OS user
type User struct {
	Name              string      `yaml:"name"`
	Description       string      `yaml:"gecos"`
	HomeDir           string      `yaml:"homedir"`
	Shell             string      `yaml:"shell"`
	NoCreateHomeDir   bool        `yaml:"no_create_homedir"`
	PrimaryGroup      string      `yaml:"primary_group"`
	Groups            MultiString `yaml:"groups"`
	System            bool        `yaml:"system"`
	SSHAuthorizedKeys []string    `yaml:"ssh_authorized_keys"`
	Password          string      `yaml:"passwd"`
}

// SSHD specifies the `sshd` entry
type SSHD struct {
	Port                   int      `yaml:"port"`
	ListenAddress          string   `yaml:"listen_address"`
	AuthorizedKeys         []string `yaml:"authorized_keys"`
	PermitRootLogin        bool     `yaml:"permit_root_login"`
	PermitEmptyPasswords   bool     `yaml:"permit_empty_passwords"`
	PasswordAuthentication bool     `yaml:"password_authentication"`
}

// DRProvision is used for installing and configuring drpcli
type DRProvision struct {
	InstallRunner bool   `yaml:"install_runner"`
	AssetsURL     string `yaml:"assets_url"`
	Token         string `yaml:"token"`
	Endpoint      string `yaml:"endpoint"`
	UUID          string `yaml:"uuid"`
}

// NetworkSettings contains all network settings lift should apply
type NetworkSettings struct {
	HostName      string               `yaml:"hostname"`
	InterfaceOpts string               `yaml:"interfaces"`
	ResolvConf    *ResolvConfiguration `yaml:"resolv_conf"`
	Proxy         string               `yaml:"proxy"`
	NTP           *NTPConfiguration    `yaml:"ntp"`
}

// ResolvConfiguration contains the DNS spec
type ResolvConfiguration struct {
	NameServers   MultiString `yaml:"nameservers"`
	SearchDomains MultiString `yaml:"search_domains"`
	Domain        string      `yaml:"domain"`
}

// NTPConfiguration is used for configuring chronyd
type NTPConfiguration struct {
	Pools   MultiString `yaml:"pools"`
	Servers MultiString `yaml:"servers"`
}

// MTAConfiguration contains all information for setting up a
// mail transfer agent (mail forwarding)
type MTAConfiguration struct {
	Root             string `yaml:"root"`
	Server           string `yaml:"server"`
	UseTLS           bool   `yaml:"use_tls"`
	UseSTARTTLS      bool   `yaml:"use_starttls"`
	User             string `yaml:"user"`
	Password         string `yaml:"password"`
	AuthMethod       string `yaml:"authmethod"`
	RewriteDomain    string `yaml:"rewrite_domain"`
	FromLineOverride bool   `yaml:"fromline_override"`
}

// PackagesConfig contains specification for the `packages:` block.
type PackagesConfig struct {
	Repositories MultiString `yaml:"repositories"`
	Update       bool        `yaml:"update"`
	Upgrade      bool        `yaml:"upgrade"`
	Install      MultiString `yaml:"install"`
	Uninstall    MultiString `yaml:"uninstall"`
}

// WriteFile allows for specifying files and their content
// that should be created on first boot.
type WriteFile struct {
	Encoding    string `yaml:"encoding"`
	Content     string `yaml:"content"`
	ContentURL  string `yaml:"content-url"`
	Path        string `yaml:"path"`
	Owner       string `yaml:"owner"`
	Permissions string `yaml:"permissions"`
}

// Disk specifies a disk that should be formatted and mounted
// (without partitioning, LUKS encrypted).
type Disk struct {
	Device         string `yaml:"device"`
	FileSystemType string `yaml:"filesystem"`
	MountPoint     string `yaml:"mountpoint"`
}

// MultiString is a type alias, needed for unmarshalling
type MultiString []string

// UnmarshalYAML is a custom unmarshalling function for parsing yaml values that
// contains one or more string (either string or array of strings)
// but always returning []string (aliased with MultiString)
func (ms *MultiString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var values []string
	err := unmarshal(&values)
	if err != nil {
		var s string
		err := unmarshal(&s)
		if err != nil {
			return err
		}
		*ms = []string{s}
	} else {
		*ms = values
	}
	return nil
}

var silent bool

// InitAlpineData initializes alpine-data with sane defaults
func InitAlpineData() *AlpineData {
	return &AlpineData{
		UnLift:   true,
		TimeZone: "UTC",
		Keymap:   "us us",
		Network: &NetworkSettings{
			HostName: "alpine",
		},
		SSHDConfig: &SSHD{
			Port:                   22,
			ListenAddress:          "0.0.0.0",
			PermitRootLogin:        true,
			PermitEmptyPasswords:   false,
			PasswordAuthentication: false,
		},
		DRP: &DRProvision{
			InstallRunner: true,
		},
		Packages: &PackagesConfig{
			Repositories: []string{
				"http://dl-cdn.alpinelinux.org/alpine/v3.8/main",
				"http://dl-cdn.alpinelinux.org/alpine/v3.8/community",
			},
		},
	}
}

// Returns a key-value map with SSH settings from alpine-data
func (l *Lift) getSSHDKVMap() map[string]string {
	return map[string]string{
		"Port":                   strconv.Itoa(l.Data.SSHDConfig.Port),
		"ListenAddress":          l.Data.SSHDConfig.ListenAddress,
		"PermitRootLogin":        boolToYesNo(l.Data.SSHDConfig.PermitRootLogin),
		"PermitEmptyPasswords":   boolToYesNo(l.Data.SSHDConfig.PermitEmptyPasswords),
		"PasswordAuthentication": boolToYesNo(l.Data.SSHDConfig.PasswordAuthentication),
	}
}

// Converts bool values to either "yes" or "no"
func boolToYesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
