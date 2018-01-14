package parglare

// Production represents a single production of a grammar.
type Production struct {

	// Left-Hand-Side grammar symbol.
	Symbol GrammarSymbol

	// Right-Hand-Side of a production.
	RHS []GrammarSymbol

	// TODO: Assignments
	Assoc Associativity

	// Priority used in disambiguation.
	Prior int

	// Should dynamic disambiguation be used. Set in the grammar.
	Dynamic bool

	// Disable "prefer shift" disambiguation strategy for this production.
	NOPS bool

	// Disable "prefer shift over empty" for this production.
	NOPSE bool
}

func NewP(symbol GrammarSymbol, rhs []GrammarSymbol) *Production {
	p := new(Production)
	p.Symbol = symbol
	p.RHS = rhs
	p.Assoc = AssocNone
	p.Prior = DefaultPriority
	return p
}

func NewPAP(symbol GrammarSymbol, rhs []GrammarSymbol,
	assoc Associativity, prior int) *Production {
	p := NewP(symbol, rhs)
	p.Assoc = assoc
	p.Prior = prior
	return p
}

// Assignment represents assignment from the grammar.
type Assignment struct {
	Name       string        // The name n the LHS of assignment
	Op         string        // Either a "=" or "?="
	Symbol     GrammarSymbol // A grammar symbol on the RHS
	OrigSymbol GrammarSymbol // A de-sugared grammar symbol on the RHS
	Mult       Multiplicity  // Multiplicity of RHS
	Index      int           // Position in production RHS
}
