package output

import (
	"encoding/json"
	"os"
)

type Result struct {
	Version              string      `json:"version"`
	Command              string      `json:"command"`
	Directory            string      `json:"directory"`
	DryRun               bool        `json:"dry_run"`
	FilesProcessed       int         `json:"files_processed"`
	FilesSplit           int         `json:"files_split"`
	FilesRenamed         int         `json:"files_renamed"`
	FilesSkipped         int         `json:"files_skipped"`
	FilesAlreadyCompliant int               `json:"files_already_compliant"`
	Changes              []Change          `json:"changes"`
	Violations           []Violation       `json:"violations,omitempty"`
	NamingViolations     []NamingViolation `json:"naming_violations,omitempty"`
	ResourcesChecked     int               `json:"resources_checked,omitempty"`
	ResourcesCompliant   int               `json:"resources_compliant,omitempty"`
	ResourcesSkipped     int               `json:"resources_skipped,omitempty"`
	ResourcesUnknown     int               `json:"resources_unknown,omitempty"`
	Errors               []string          `json:"errors,omitempty"`
	ExitCode             int               `json:"exit_code"`
}

type Change struct {
	Action  string   `json:"action"`
	Source  string   `json:"source"`
	Target  string   `json:"target,omitempty"`
	Targets []string `json:"targets,omitempty"`
	Reason  string   `json:"reason,omitempty"`
}

type Violation struct {
	File         string `json:"file"`
	CurrentName  string `json:"current_name"`
	ExpectedName string `json:"expected_name"`
	BlockType    string `json:"block_type"`
}

type NamingViolation struct {
	File            string   `json:"file"`
	ResourceType    string   `json:"resource_type"`
	ResourceLabel   string   `json:"resource_label"`
	NameValue       string   `json:"name_value"`
	ExpectedPrefix  string   `json:"expected_prefix"`
	ExpectedPattern string   `json:"expected_pattern"`
	Issues          []string `json:"issues"`
}

func NewResult(version, command, directory string, dryRun bool) *Result {
	return &Result{
		Version:   version,
		Command:   command,
		Directory: directory,
		DryRun:    dryRun,
		Changes:   []Change{},
		Errors:    []string{},
	}
}

func (r *Result) AddSplit(source string, targets []string) {
	r.Changes = append(r.Changes, Change{
		Action:  "split",
		Source:  source,
		Targets: targets,
	})
	r.FilesSplit++
}

func (r *Result) AddRename(source, target string) {
	r.Changes = append(r.Changes, Change{
		Action: "rename",
		Source: source,
		Target: target,
	})
	r.FilesRenamed++
}

func (r *Result) AddSkip(source, reason string) {
	r.Changes = append(r.Changes, Change{
		Action: "skip",
		Source: source,
		Reason: reason,
	})
	r.FilesSkipped++
}

func (r *Result) AddCompliant(source string) {
	r.FilesAlreadyCompliant++
}

func (r *Result) AddViolation(file, currentName, expectedName, blockType string) {
	if r.Violations == nil {
		r.Violations = []Violation{}
	}
	r.Violations = append(r.Violations, Violation{
		File:         file,
		CurrentName:  currentName,
		ExpectedName: expectedName,
		BlockType:    blockType,
	})
}

func (r *Result) AddNamingViolation(file, resourceType, resourceLabel, nameValue, expectedPrefix, expectedPattern string, issues []string) {
	if r.NamingViolations == nil {
		r.NamingViolations = []NamingViolation{}
	}
	r.NamingViolations = append(r.NamingViolations, NamingViolation{
		File:            file,
		ResourceType:    resourceType,
		ResourceLabel:   resourceLabel,
		NameValue:       nameValue,
		ExpectedPrefix:  expectedPrefix,
		ExpectedPattern: expectedPattern,
		Issues:          issues,
	})
}

func (r *Result) AddError(err string) {
	r.Errors = append(r.Errors, err)
}

func (r *Result) Print() error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(r)
}
