package models

// Options Main config structure
type Options struct {
	General struct {
		Debug                bool   `long:"debug" description:"Debug log output"`
		DoNotRemoveTmpFolder bool   `long:"do-not-remove-tmp-folder" description:"Do not remove the tmp folder where all mutations are saved to"`
		Help                 bool   `long:"help" description:"Show this help message"`
		Verbose              bool   `long:"verbose" description:"Verbose log output"`
		Config               string `long:"config" description:"Path to config file"`
	} `group:"General options"`

	Files struct {
		Blacklist []string `long:"blacklist" description:"List of MD5 checksums of mutations which should be ignored. Each checksum must end with a new line character."`
		ListFiles bool     `long:"list-files" description:"List found files"`
		PrintAST  bool     `long:"print-ast" description:"Print the ASTs of all given files and exit"`
	} `group:"File options"`

	Mutator struct {
		DisableMutators []string `long:"disable" description:"Disable mutator by their name or using * as a suffix pattern (in order to check remaining enabled mutators use --verbose option)"`
		ListMutators    bool     `long:"list-mutators" description:"List all available mutators (including disabled)"`
	} `group:"Mutator options"`

	Filter struct {
		Match string `long:"match" description:"Only functions are mutated that confirm to the arguments regex"`
	} `group:"Filter options"`

	Exec struct {
		Exec    string `long:"exec" description:"Execute this command for every mutation (by default the built-in exec command is used)"`
		NoExec  bool   `long:"no-exec" description:"Skip the built-in exec command and just generate the mutations"`
		Timeout uint   `long:"exec-timeout" description:"Sets a timeout for the command execution (in seconds)" default:"10"`
	} `group:"Exec options"`

	Test struct {
		Recursive bool `long:"test-recursive" description:"Defines if the executer should test recursively"`
	} `group:"Test options"`

	Remaining struct {
		Targets []string `description:"Packages, directories and files even with patterns (by default the current directory)"`
	} `positional-args:"true" required:"true"`

	Config struct {
		SkipFileWithoutTest  bool     `yaml:"skip_without_test"`
		SkipFileWithBuildTag bool     `yaml:"skip_with_build_tags"`
		JSONOutput           bool     `yaml:"json_output"`
		SilentMode           bool     `yaml:"silent_mode"`
		ExcludeDirs          []string `yaml:"exclude_dirs"`
	}
}
