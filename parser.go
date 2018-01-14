package parglare

type Parser interface {
	Parse(input interface{}) interface{}
	ParseFrom(input interface{}, pos int) interface{}
}

type LRParser struct {
}

func NewLRParser(g *Grammar) Parser {
	p := &LRParser{}
	return p
}

func (p *LRParser) Parse(input interface{}) interface{} {
	return nil
}

func (p *LRParser) ParseFrom(input interface{}, pos int) interface{} {
	return nil
}
