package operations

import (
	"strings"

	"github.com/obay/tfren/internal/azure"
	"github.com/obay/tfren/internal/hcl"
	"github.com/obay/tfren/internal/output"
)

func AzureNaming(opts Options, result *output.Result, console *output.Console) error {
	files, err := GetTerraformFiles(opts.Directory, opts.Recursive)
	if err != nil {
		return err
	}

	for _, file := range files {
		result.FilesProcessed++
		if err := checkAzureNaming(file, opts, result, console); err != nil {
			result.AddError(err.Error())
		}
	}

	if len(result.NamingViolations) > 0 && opts.Strict {
		result.ExitCode = 2
	}

	return nil
}

func checkAzureNaming(path string, opts Options, result *output.Result, console *output.Console) error {
	parsed, err := hcl.ParseFile(path)
	if err != nil {
		return err
	}

	if parsed.ParseError != nil {
		result.AddError("Failed to parse: " + path)
		return nil
	}

	for _, block := range parsed.Blocks {
		if block.Type != "resource" {
			continue
		}
		if len(block.Labels) < 2 {
			continue
		}

		resourceType := block.Labels[0]
		resourceLabel := block.Labels[1]

		if !strings.HasPrefix(resourceType, "azurerm_") {
			continue
		}

		// Check if we have a CAF rule for this resource type
		_, known := azure.LookupRule(resourceType)
		if !known {
			result.ResourcesUnknown++
			if console != nil && opts.Verbose {
				console.Info("No CAF rule for " + resourceType + "." + resourceLabel)
			}
			continue
		}

		result.ResourcesChecked++

		// Check if name attribute exists
		if block.NameValue == "" && !block.NameDynamic {
			// No name attribute found - skip silently
			continue
		}

		// Check if name is dynamic
		if block.NameDynamic {
			result.ResourcesSkipped++
			if console != nil {
				console.Info(resourceType + "." + resourceLabel + ": skipped (dynamic name)")
			}
			continue
		}

		// Validate static name
		issues := azure.ValidateResourceName(resourceType, block.NameValue)
		if len(issues) == 0 {
			result.ResourcesCompliant++
			if console != nil && opts.Verbose {
				console.Success(resourceType + "." + resourceLabel + ": name \"" + block.NameValue + "\" OK")
			}
		} else {
			result.AddNamingViolation(
				path,
				resourceType,
				resourceLabel,
				block.NameValue,
				azure.GetExpectedPrefix(resourceType),
				azure.GetExpectedPattern(resourceType),
				issues,
			)
			if console != nil {
				console.Warning(resourceType + "." + resourceLabel + ": name \"" + block.NameValue + "\"")
				for _, issue := range issues {
					console.Print("    - " + issue)
				}
				console.Print("    Expected pattern: " + azure.GetExpectedPattern(resourceType))
			}
		}
	}

	return nil
}
