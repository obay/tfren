package hcl

import (
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func ParseFile(path string) (*ParsedFile, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(content, path)
	if diags.HasErrors() {
		return &ParsedFile{
			Path:       path,
			ParseError: diags,
		}, nil
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return &ParsedFile{
			Path:   path,
			Blocks: []Block{},
		}, nil
	}

	lines := strings.Split(string(content), "\n")
	blocks := extractBlocksWithComments(body, lines, content)

	return &ParsedFile{
		Path:   path,
		Blocks: blocks,
	}, nil
}

func extractBlocksWithComments(body *hclsyntax.Body, lines []string, source []byte) []Block {
	var blocks []Block

	for _, hclBlock := range body.Blocks {
		startLine := hclBlock.TypeRange.Start.Line
		endLine := hclBlock.Body.EndRange.Start.Line

		block := Block{
			Type:      hclBlock.Type,
			Labels:    hclBlock.Labels,
			StartLine: startLine,
			EndLine:   endLine,
		}

		if hclBlock.Type == "provider" {
			block.Alias = extractAlias(hclBlock)
		}

		if hclBlock.Type == "resource" {
			name, isDynamic := extractName(hclBlock)
			block.NameValue = name
			block.NameDynamic = isDynamic
		}

		block.Comments = findPrecedingComments(lines, startLine)

		block.RawContent = extractLines(lines, startLine-1, endLine)

		blocks = append(blocks, block)
	}

	return blocks
}

func extractAlias(block *hclsyntax.Block) string {
	for name, attr := range block.Body.Attributes {
		if name == "alias" {
			if val, diags := attr.Expr.Value(nil); !diags.HasErrors() {
				return val.AsString()
			}
		}
	}
	return ""
}

func extractName(block *hclsyntax.Block) (value string, isDynamic bool) {
	for name, attr := range block.Body.Attributes {
		if name == "name" {
			if val, diags := attr.Expr.Value(nil); !diags.HasErrors() {
				return val.AsString(), false
			}
			return "", true
		}
	}
	return "", false
}
