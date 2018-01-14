package parglare

type GLRParser struct {
	LRParser
}

func NewGLRParser(g *Grammar) Parser {
	p := &GLRParser{}
	return p
}

func (p *GLRParser) Parse(input interface{}) interface{} {
	return nil
}

func (p *GLRParser) ParseFrom(input interface{}, pos int) interface{} {
	return nil
}
