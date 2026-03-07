package operations

type Options struct {
	Directory    string
	Recursive    bool
	DryRun       bool
	Verbose      bool
	KeepOriginal bool
	Backup       bool
	Strict       bool
}
