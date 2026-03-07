package naming

import "fmt"

func GenerateFileName(blockType string, labels []string, alias string) string {
	switch blockType {
	case "resource", "data":
		if len(labels) >= 2 {
			return fmt.Sprintf("%s.%s.%s.tf", blockType, labels[0], labels[1])
		}
	case "provider":
		if len(labels) >= 1 {
			if alias != "" {
				return fmt.Sprintf("provider.%s.%s.tf", labels[0], alias)
			}
			return fmt.Sprintf("provider.%s.tf", labels[0])
		}
	case "variable", "module", "output":
		if len(labels) >= 1 {
			return fmt.Sprintf("%s.%s.tf", blockType, labels[0])
		}
	case "locals":
		return "locals.tf"
	case "terraform":
		return "terraform.tf"
	}
	return ""
}

func IsCompliant(filename string, blockType string, labels []string, alias string) bool {
	expected := GenerateFileName(blockType, labels, alias)
	return filename == expected
}
