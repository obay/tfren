package hcl

type Block struct {
	Type       string
	Labels     []string
	Alias      string
	StartLine  int
	EndLine    int
	Comments   []string
	RawContent string
}

type ParsedFile struct {
	Path       string
	Blocks     []Block
	ParseError error
}

func (b *Block) ContentWithComments() string {
	result := ""
	for _, comment := range b.Comments {
		result += comment + "\n"
	}
	result += b.RawContent
	return result
}
