package parglare

type Associativity int

const (
	// Associativities
	AssocNone Associativity = iota
	AssocLeft
	AssocRight
)

// Priority
const DefaultPriority int = 10

type Multiplicity int

const (
	MultOne Multiplicity = iota
	MultOptional
	MultOneOrMore
	MultZeroOrMore
)

var multiplicities = [...]string{"1", "0..1", "1..*", "0..*"}

func (m Multiplicity) String() string {
	return multiplicities[m]
}

// Specials
var ReservedSymbolNames = [...]string{"EOF", "STOP", "EMPTY"}
var SpecialSymbolNames = [...]string{"KEYWORD", "LAYOUT"}

type ActionFunc func(context Context, subresults []interface{}) (interface{}, error)

type Grammar struct {
	Productions  []*Production
	RootSymbol   GrammarSymbol
	Terminals    map[string]Terminal
	NonTerminals map[string]NonTerminal
}

// Parsing context used when calling ActionFunc
type Context struct {
}

type GrammarParams struct {
}

// GrammarSymbol represents an abstract grammar symbol, either Terminal Or
// NonTerminal
type GrammarSymbol interface {
	// the name of this grammar symbol
	Name() string
	// name of common/user action given in the grammar.
	ActionName() string
	// Resolved action given by the user. Overrides grammar action if
	// provided. If not provided by the user defaults to grammarAction.
	Action() ActionFunc
	// Action provided in the grammar.
	GrammarAction() ActionFunc
}

type grammarSymbol struct {
	Name          string
	ActionName    string
	Action        ActionFunc
	GrammarAction ActionFunc
}

type Terminal struct {
	grammarSymbol

	// Priority used for lexical disambiguation.
	Prior int

	// Recognizer for this terminal.
	TermRecognizer Recognizer

	// Used for scanning optimization. If this terminal is `finish` no other
	// recognizers will be checked if this succeeds. If not provided in the
	// grammar implicit rules will be used during table construction.
	Finish bool

	// If true prefer this recognizer in case of multiple recognizers match
	// at the same place and implicit disambiguation doesn't resolve.
	Prefer bool

	// Should dynamic disambiguation be called to resolve conflict involving
	// this terminal.
	Dynamic bool

	// true if this terminal represents a keyword.
	Keyword bool
}

type NonTerminal struct {
	grammarSymbol
}

func (t *Terminal) Name() string {
	return t.grammarSymbol.Name

}

func (t *Terminal) ActionName() string {
	return t.grammarSymbol.ActionName

}

func (t *Terminal) Action() ActionFunc {
	return t.grammarSymbol.Action

}

func (t *Terminal) GrammarAction() ActionFunc {
	return t.grammarSymbol.GrammarAction

}

func (nt *NonTerminal) Name() string {
	return nt.grammarSymbol.Name

}

func (nt *NonTerminal) ActionName() string {
	return nt.grammarSymbol.ActionName

}

func (nt *NonTerminal) Action() ActionFunc {
	return nt.grammarSymbol.Action

}

func (nt *NonTerminal) GrammarAction() ActionFunc {
	return nt.grammarSymbol.GrammarAction

}

func (g *Grammar) initGrammar() {

}

func newNT(name string) *NonTerminal {
	return &NonTerminal{grammarSymbol{name, "", nil, nil}}
}

func newT(name, regex string) *Terminal {
	t := new(Terminal)
	t.grammarSymbol.Name = name
	t.TermRecognizer = NewRegExRecognizer(name, regex, 0, false)
	t.Prior = DefaultPriority
	return t
}

func newST(name, str string) *Terminal {
	t := new(Terminal)
	t.grammarSymbol.Name = name
	t.TermRecognizer = &StringRecognizer{name, str, false}
	t.Prior = DefaultPriority
	return t

}

func GrammarFromString(gstr string) *Grammar {
	g := new(Grammar)
	return g
}

func GrammarFromFile(gfile string) *Grammar {
	g := new(Grammar)
	return g
}

func grammarParser() Parser {
	g := &Grammar{Productions: pgProductions, RootSymbol: PGGrammar}
	g.initGrammar()
	p := NewLRParser(g)
	return p
}

// Parglare grammar symbols
var (
	// Non-terminals
	PGGrammar                           = newNT("Grammar")
	PGRules                             = newNT("Rules")
	PGRule                              = newNT("Rule")
	PGProductionRule                    = newNT("ProductionRule")
	PGProductionRuleRHS                 = newNT("ProductionRuleRHS")
	PGProduction                        = newNT("Production")
	PGTerminalRule                      = newNT("TerminalRule")
	PGProductionDisambiguationRule      = newNT("ProductionDisambiguationRule")
	PGProductionDisambiguationRules     = newNT("ProductionDisambiguationRules")
	PGTerminalDisambiguationRule        = newNT("TerminalDisambiguationRule")
	PGTerminalDisambiguationRules       = newNT("TerminalDisambiguationRules")
	PGAssignment                        = newNT("Assignment")
	PGAssignments                       = newNT("Assignments")
	PGPlainAssignment                   = newNT("PlainAssignment")
	PGBoolAssignment                    = newNT("BoolAssignment")
	PGRepeatableGrammarSymbol           = newNT("RepeatableGrammarSymbol")
	PGRepeatableGrammarSymbols          = newNT("RepeatableGrammarSymbols")
	PGOptRepeatOperator                 = newNT("OptRepeatOperator")
	PGRepeatOperatorZero                = newNT("RepeatOperatorZero")
	PGRepeatOperatorOne                 = newNT("RepeatOperatorOne")
	PGRepeatOperatorOptional            = newNT("RepeatOperatorOptional")
	PGOptionalRepeatModifiersExpression = newNT("OptionalRepeatModifiersExpression")
	PGOptionalRepeatModifiers           = newNT("OptionalRepeatModifiers")
	PGOptionalRepeatModifier            = newNT("OptionalRepeatModifier")
	PGGrammarSymbol                     = newNT("GrammarSymbol")
	PGRecognizer                        = newNT("Recognizer")
	PGLayout                            = newNT("LAYOUT")
	PGLayoutItem                        = newNT("LAYOUT_ITEM")
	PGComment                           = newNT("Comment")
	PGCORNC                             = newNT("CORNC")
	PGCORNCS                            = newNT("CORNCS")

	// Terminals
	PGName        = newT("Name", `"[a-zA-Z0-9_]+`)
	PGStrTerm     = newT("StrTerm", `(?s)('[^'\\]*(?:\\.[^'\\]*)*')|("[^"\\]*(?:\\.[^"\\]*)*")`)
	PGRegExTerm   = newT("RegExTerm", `\/((\\/)|[^/])*\/`)
	PGPrior       = newT("Prior", `\d+`)
	PGAction      = newT("Action", `@[a-zA-Z0-9_]+`)
	PGWS          = newT("WS", `\s+`)
	PGCommentLine = newT("CommentLine", `\/\/.*`)
	PGNotComment  = newT("NotComment", `((\*[^\/])|[^\s*\/]|\/[^\*])+`)

	PGSemiColon  = newST("SemiColon", ";")
	PGColon      = newST("Colon", ":")
	PGComma      = newST("Comma", ",")
	PGPipe       = newST("Pipe", "|")
	PGOpenBrace  = newST("OpenBrace", "{")
	PGCloseBrace = newST("CloseBrace", "}")
	PGLeft       = newST("Left", "left")
	PGRight      = newST("Right", "right")
	PGReduce     = newST("Reduce", "reduce")
	PGShift      = newST("Shift", "shift")
	PGDynamic    = newST("Dynamic", "dynamic")
	PGNOPS       = newST("NOPS", "nops")
	PGNOPSE      = newST("NOPSE", "nopse")
	PGPrefer     = newST("Prefer", "prefer")
	PGFinish     = newST("Finish", "finish")
	PGNoFinish   = newST("NoFinish", "nofinish")

	PGEquals    = newST("Equals", "=")
	PGOptEquals = newST("OptEquals", "?=")
	PGAsterisk  = newST("Asterisk", "*")
	PGPlus      = newST("Plus", "+")
	PGQuestion  = newST("Question", "?")

	PGOpenSquare  = newST("OpenSquare", "[")
	PGCloseSquare = newST("CloseSquare", "]")
	PGBlockCStart = newST("BlockCStart", "/*")
	PGBlockCEnd   = newST("BlockCEnd", "*/")
)

var EOF = &Terminal{grammarSymbol: grammarSymbol{"EOF", "", nil, nil},
	Prior: DefaultPriority}
var EMPTY = &Terminal{grammarSymbol: grammarSymbol{"EMPTY", "", nil, nil},
	Prior: DefaultPriority}

var pgProductions = []*Production{
	NewP(PGGrammar, []GrammarSymbol{PGRules, EOF}),
	NewP(PGRules, []GrammarSymbol{PGRules, PGRule}),
	NewP(PGRules, []GrammarSymbol{PGRule}),
	NewP(PGRule, []GrammarSymbol{PGProductionRule}),
	NewP(PGRule, []GrammarSymbol{PGAction, PGProductionRule}),
	NewP(PGRule, []GrammarSymbol{PGTerminalRule}),
	NewP(PGRule, []GrammarSymbol{PGAction, PGTerminalRule}),

	NewP(PGProductionRule,
		[]GrammarSymbol{PGName, PGColon, PGProductionRuleRHS, PGSemiColon}),
	NewPAP(PGProductionRuleRHS,
		[]GrammarSymbol{PGProductionRuleRHS, PGPipe, PGProduction},
		AssocLeft, 5),
	NewPAP(PGProductionRuleRHS,
		[]GrammarSymbol{PGProduction}, AssocLeft, 5),
	NewP(PGProduction, []GrammarSymbol{PGAssignments}),
	NewP(PGProduction, []GrammarSymbol{PGAssignments, PGOpenBrace,
		PGProductionDisambiguationRules, PGCloseBrace}),

	NewPAP(PGTerminalRule,
		[]GrammarSymbol{PGName, PGColon, PGRecognizer, PGSemiColon},
		AssocLeft, 15),
	NewPAP(PGTerminalRule,
		[]GrammarSymbol{PGName, PGColon, PGSemiColon}, AssocLeft, 15),
	NewPAP(PGTerminalRule, []GrammarSymbol{PGName, PGColon, PGRecognizer,
		PGOpenBrace, PGTerminalDisambiguationRules, PGCloseBrace,
		PGSemiColon}, AssocLeft, 15),
	NewPAP(PGTerminalRule,
		[]GrammarSymbol{PGName, PGColon,
			PGOpenBrace, PGTerminalDisambiguationRules, PGCloseBrace,
			PGSemiColon}, AssocLeft, 15),

	NewP(PGProductionDisambiguationRule, []GrammarSymbol{PGLeft}),
	NewP(PGProductionDisambiguationRule, []GrammarSymbol{PGReduce}),
	NewP(PGProductionDisambiguationRule, []GrammarSymbol{PGRight}),
	NewP(PGProductionDisambiguationRule, []GrammarSymbol{PGShift}),
	NewP(PGProductionDisambiguationRule, []GrammarSymbol{PGDynamic}),
	NewP(PGProductionDisambiguationRule, []GrammarSymbol{PGNOPS}),
	NewP(PGProductionDisambiguationRule, []GrammarSymbol{PGNOPSE}),
	NewP(PGProductionDisambiguationRule, []GrammarSymbol{PGPrior}),
	NewPAP(PGProductionDisambiguationRules,
		[]GrammarSymbol{PGProductionDisambiguationRules, PGComma,
			PGProductionDisambiguationRule},
		AssocLeft, DefaultPriority),
	NewP(PGProductionDisambiguationRules,
		[]GrammarSymbol{PGProductionDisambiguationRule}),

	NewP(PGTerminalDisambiguationRule, []GrammarSymbol{PGPrefer}),
	NewP(PGTerminalDisambiguationRule, []GrammarSymbol{PGFinish}),
	NewP(PGTerminalDisambiguationRule, []GrammarSymbol{PGNoFinish}),
	NewP(PGTerminalDisambiguationRule, []GrammarSymbol{PGDynamic}),
	NewP(PGTerminalDisambiguationRule, []GrammarSymbol{PGPrior}),
	NewP(PGTerminalDisambiguationRules,
		[]GrammarSymbol{PGTerminalDisambiguationRules, PGComma,
			PGTerminalDisambiguationRule}),
	NewP(PGTerminalDisambiguationRules,
		[]GrammarSymbol{PGTerminalDisambiguationRule}),

	// Assignments
	NewP(PGAssignment, []GrammarSymbol{PGPlainAssignment}),
	NewP(PGAssignment, []GrammarSymbol{PGBoolAssignment}),
	NewP(PGAssignment, []GrammarSymbol{PGRepeatableGrammarSymbol}),
	NewP(PGAssignments, []GrammarSymbol{PGAssignments, PGAssignment}),
	NewP(PGAssignments, []GrammarSymbol{PGAssignment}),
	NewP(PGPlainAssignment,
		[]GrammarSymbol{PGName, PGEquals, PGRepeatableGrammarSymbol}),
	NewP(PGBoolAssignment,
		[]GrammarSymbol{PGName, PGOptEquals, PGRepeatableGrammarSymbol}),

	// Regex-like repeat operators
	NewP(PGRepeatableGrammarSymbol,
		[]GrammarSymbol{PGGrammarSymbol, PGOptRepeatOperator}),
	NewP(PGOptRepeatOperator, []GrammarSymbol{PGRepeatOperatorZero}),
	NewP(PGOptRepeatOperator, []GrammarSymbol{PGRepeatOperatorOne}),
	NewP(PGOptRepeatOperator, []GrammarSymbol{PGRepeatOperatorOptional}),
	NewP(PGOptRepeatOperator, []GrammarSymbol{EMPTY}),
	NewP(PGRepeatOperatorZero,
		[]GrammarSymbol{PGAsterisk, PGOptionalRepeatModifiersExpression}),
	NewP(PGRepeatOperatorOne,
		[]GrammarSymbol{PGPlus, PGOptionalRepeatModifiersExpression}),
	NewP(PGRepeatOperatorOptional,
		[]GrammarSymbol{PGQuestion, PGOptionalRepeatModifiersExpression}),
	NewP(PGOptionalRepeatModifiersExpression,
		[]GrammarSymbol{PGOpenSquare, PGOptionalRepeatModifiers,
			PGCloseSquare}),
	NewP(PGOptionalRepeatModifiersExpression, []GrammarSymbol{EMPTY}),
	NewP(PGOptionalRepeatModifiers,
		[]GrammarSymbol{PGOptionalRepeatModifiers,
			PGComma, PGOptionalRepeatModifier}),
	NewP(PGOptionalRepeatModifiers, []GrammarSymbol{PGOptionalRepeatModifier}),

	NewP(PGGrammarSymbol, []GrammarSymbol{PGName}),
	NewP(PGGrammarSymbol, []GrammarSymbol{PGRecognizer}),
	NewP(PGRecognizer, []GrammarSymbol{PGStrTerm}),
	NewP(PGRecognizer, []GrammarSymbol{PGRegExTerm}),

	// Support for comments
	NewP(PGLayout, []GrammarSymbol{PGLayoutItem}),
	NewP(PGLayout, []GrammarSymbol{PGLayout, PGLayoutItem}),
	NewP(PGLayoutItem, []GrammarSymbol{PGWS}),
	NewP(PGLayoutItem, []GrammarSymbol{PGComment}),
	NewP(PGLayoutItem, []GrammarSymbol{EMPTY}),
	NewP(PGComment, []GrammarSymbol{PGBlockCStart, PGCORNCS, PGBlockCEnd}),
	NewP(PGComment, []GrammarSymbol{PGCommentLine}),
	NewP(PGCORNCS, []GrammarSymbol{PGCORNC}),
	NewP(PGCORNCS, []GrammarSymbol{PGCORNCS, PGCORNC}),
	NewP(PGCORNCS, []GrammarSymbol{EMPTY}),
	NewP(PGCORNC, []GrammarSymbol{PGComment}),
	NewP(PGCORNC, []GrammarSymbol{PGNotComment}),
	NewP(PGCORNC, []GrammarSymbol{PGWS}),
}
