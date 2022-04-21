package main

// Symbols of the theme
type SymbolTemplate struct {
	Lock                 string
	Network              string
	NetworkAlternate     string
	Separator            string
	SeparatorThin        string
	SeparatorReverse     string
	SeparatorReverseThin string

	RepoDetached   string
	RepoBranch     string
	RepoAhead      string
	RepoBehind     string
	RepoStaged     string
	RepoNotStaged  string
	RepoUntracked  string
	RepoConflicted string
	RepoStashed    string

	VenvIndicator string
	NodeIndicator string
	RvmIndicator  string
}

// Theme definitions
type Theme struct {
	BoldForeground bool

	Reset uint8

	DefaultFg uint8
	DefaultBg uint8

	UsernameFg     uint8
	UsernameBg     uint8
	UsernameRootBg uint8

	HostnameFg uint8
	HostnameBg uint8

	// The foreground-background mapping is precomputed and stored in a map for improved performance
	// The old script used to brute-force this at runtime
	HostnameColorizedFgMap map[uint8]uint8

	HomeSpecialDisplay bool
	HomeFg             uint8
	HomeBg             uint8
	AliasFg            uint8
	AliasBg            uint8
	PathFg             uint8
	PathBg             uint8
	CwdFg              uint8
	SeparatorFg        uint8

	ReadonlyFg uint8
	ReadonlyBg uint8

	SSHFg uint8
	SSHBg uint8

	DockerMachineFg uint8
	DockerMachineBg uint8

	KubeClusterFg   uint8
	KubeClusterBg   uint8
	KubeNamespaceFg uint8
	KubeNamespaceBg uint8

	WSLMachineFg uint8
	WSLMachineBg uint8

	DotEnvFg uint8
	DotEnvBg uint8

	AWSFg uint8
	AWSBg uint8

	RepoCleanFg uint8
	RepoCleanBg uint8
	RepoDirtyFg uint8
	RepoDirtyBg uint8

	JobsFg uint8
	JobsBg uint8

	CmdPassedFg uint8
	CmdPassedBg uint8
	CmdFailedFg uint8
	CmdFailedBg uint8

	SvnChangesFg uint8
	SvnChangesBg uint8

	GCPFg uint8
	GCPBg uint8

	GitAheadFg      uint8
	GitAheadBg      uint8
	GitBehindFg     uint8
	GitBehindBg     uint8
	GitStagedFg     uint8
	GitStagedBg     uint8
	GitNotStagedFg  uint8
	GitNotStagedBg  uint8
	GitUntrackedFg  uint8
	GitUntrackedBg  uint8
	GitConflictedFg uint8
	GitConflictedBg uint8
	GitStashedFg    uint8
	GitStashedBg    uint8

	GoenvFg uint8
	GoenvBg uint8

	VirtualEnvFg uint8
	VirtualEnvBg uint8

	VirtualGoFg uint8
	VirtualGoBg uint8

	PerlbrewFg uint8
	PerlbrewBg uint8

	PlEnvFg uint8
	PlEnvBg uint8

	TFWsFg uint8
	TFWsBg uint8

	TimeFg uint8
	TimeBg uint8

	ShellVarFg uint8
	ShellVarBg uint8

	ShEnvFg uint8
	ShEnvBg uint8

	NodeFg        uint8
	NodeBg        uint8
	NodeVersionFg uint8
	NodeVersionBg uint8

	RvmFg        uint8
	RvmBg        uint8

	LoadFg           uint8
	LoadBg           uint8
	LoadHighBg       uint8
	LoadAvgValue     byte
	LoadThresholdBad float64

	NixShellFg uint8
	NixShellBg uint8

	DurationFg uint8
	DurationBg uint8

	ViModeCommandFg uint8
	ViModeCommandBg uint8
	ViModeInsertFg uint8
	ViModeInsertBg uint8
}
