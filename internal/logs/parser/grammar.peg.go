package parser

import (
	"github.com/apex/up/internal/logs/parser/ast"
	"fmt"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleQuery
	rulePrimaryExpr
	ruleTupleExpr
	ruleInExpr
	ruleNotInExpr
	rulePostfixExpr
	ruleUnaryExpr
	ruleRelationalExpr
	ruleEqualityExpr
	ruleLogicalAndExpr
	ruleLogicalOrExpr
	ruleLowNotExpr
	ruleExpr
	ruleString
	ruleStringChar
	ruleUnquotedString
	ruleUnquotedStringStartChar
	ruleUnquotedStringChar
	ruleEscape
	ruleSimpleEscape
	ruleOctalEscape
	ruleHexEscape
	ruleUniversalCharacter
	ruleHexQuad
	ruleHexDigit
	ruleNumbers
	ruleNumber
	ruleInteger
	ruleFloat
	ruleFraction
	ruleExponent
	ruleStage
	ruleDEVELOPMENT
	ruleSTAGING
	rulePRODUCTION
	ruleUnit
	ruleDuration
	ruleS
	ruleMS
	ruleBytes
	ruleB
	ruleKB
	ruleMB
	ruleGB
	ruleId
	ruleIdChar
	ruleIdCharNoDigit
	ruleSeverity
	ruleIN
	ruleOR
	ruleAND
	ruleNOT
	ruleCONTAINS
	ruleDEBUG
	ruleINFO
	ruleWARN
	ruleERROR
	ruleFATAL
	ruleKeyword
	ruleEQ
	ruleLBRK
	ruleRBRK
	ruleLPAR
	ruleRPAR
	ruleDOT
	ruleBANG
	ruleLT
	ruleGT
	ruleLE
	ruleEQEQ
	ruleGE
	ruleNE
	ruleANDAND
	ruleOROR
	ruleCOMMA
	rule_
	ruleWhitespace
	ruleEOL
	ruleEOF
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9
	ruleAction10
	ruleAction11
	ruleAction12
	ruleAction13
	ruleAction14
	ruleAction15
	ruleAction16
	ruleAction17
	ruleAction18
	ruleAction19
	ruleAction20
	ruleAction21
	ruleAction22
	ruleAction23
	ruleAction24
	ruleAction25
	ruleAction26
	ruleAction27
	ruleAction28
	ruleAction29
	ruleAction30
	rulePegText
	ruleAction31
)

var rul3s = [...]string{
	"Unknown",
	"Query",
	"PrimaryExpr",
	"TupleExpr",
	"InExpr",
	"NotInExpr",
	"PostfixExpr",
	"UnaryExpr",
	"RelationalExpr",
	"EqualityExpr",
	"LogicalAndExpr",
	"LogicalOrExpr",
	"LowNotExpr",
	"Expr",
	"String",
	"StringChar",
	"UnquotedString",
	"UnquotedStringStartChar",
	"UnquotedStringChar",
	"Escape",
	"SimpleEscape",
	"OctalEscape",
	"HexEscape",
	"UniversalCharacter",
	"HexQuad",
	"HexDigit",
	"Numbers",
	"Number",
	"Integer",
	"Float",
	"Fraction",
	"Exponent",
	"Stage",
	"DEVELOPMENT",
	"STAGING",
	"PRODUCTION",
	"Unit",
	"Duration",
	"S",
	"MS",
	"Bytes",
	"B",
	"KB",
	"MB",
	"GB",
	"Id",
	"IdChar",
	"IdCharNoDigit",
	"Severity",
	"IN",
	"OR",
	"AND",
	"NOT",
	"CONTAINS",
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"FATAL",
	"Keyword",
	"EQ",
	"LBRK",
	"RBRK",
	"LPAR",
	"RPAR",
	"DOT",
	"BANG",
	"LT",
	"GT",
	"LE",
	"EQEQ",
	"GE",
	"NE",
	"ANDAND",
	"OROR",
	"COMMA",
	"_",
	"Whitespace",
	"EOL",
	"EOF",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
	"Action10",
	"Action11",
	"Action12",
	"Action13",
	"Action14",
	"Action15",
	"Action16",
	"Action17",
	"Action18",
	"Action19",
	"Action20",
	"Action21",
	"Action22",
	"Action23",
	"Action24",
	"Action25",
	"Action26",
	"Action27",
	"Action28",
	"Action29",
	"Action30",
	"PegText",
	"Action31",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Printf(" ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Printf("%v %v\n", rule, quote)
			} else {
				fmt.Printf("\x1B[34m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(buffer string) {
	node.print(false, buffer)
}

func (node *node32) PrettyPrint(buffer string) {
	node.print(true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	if tree := t.tree; int(index) >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	t.tree[index] = token32{
		pegRule: rule,
		begin:   begin,
		end:     end,
	}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type parser struct {
	stack  []ast.Node
	number string

	Buffer string
	buffer []rune
	rules  [113]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *parser) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *parser) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *parser
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *parser) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *parser) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for _, token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:
			p.AddNumber(text)
		case ruleAction1:
			p.AddNumber("")
		case ruleAction2:
			p.AddLevel(text)
		case ruleAction3:
			p.AddStage(text)
		case ruleAction4:
			p.AddField(text)
		case ruleAction5:
			p.AddString(text)
		case ruleAction6:
			p.AddString(text)
		case ruleAction7:
			p.AddExpr()
		case ruleAction8:
			p.AddTupleValue()
		case ruleAction9:
			p.AddTupleValue()
		case ruleAction10:
			p.AddTuple()
		case ruleAction11:
			p.AddBinary(ast.IN)
		case ruleAction12:
			p.AddTuple()
		case ruleAction13:
			p.AddBinary(ast.IN)
			p.AddUnary(ast.LNOT)
		case ruleAction14:
			p.AddMember(text)
		case ruleAction15:
			p.AddSubscript(text)
		case ruleAction16:
			p.AddUnary(ast.NOT)
		case ruleAction17:
			p.AddBinary(ast.GE)
		case ruleAction18:
			p.AddBinary(ast.GT)
		case ruleAction19:
			p.AddBinary(ast.LE)
		case ruleAction20:
			p.AddBinary(ast.LT)
		case ruleAction21:
			p.AddBinary(ast.EQ)
		case ruleAction22:
			p.AddBinary(ast.NE)
		case ruleAction23:
			p.AddBinary(ast.EQ)
		case ruleAction24:
			p.AddBinaryContains()
		case ruleAction25:
			p.AddBinary(ast.AND)
		case ruleAction26:
			p.AddBinary(ast.AND)
		case ruleAction27:
			p.AddBinary(ast.AND)
		case ruleAction28:
			p.AddBinary(ast.OR)
		case ruleAction29:
			p.AddBinary(ast.OR)
		case ruleAction30:
			p.AddUnary(ast.LNOT)
		case ruleAction31:
			p.SetNumber(text)

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *parser) Init() {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 Query <- <(_ Expr _ EOF)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				if !_rules[rule_]() {
					goto l0
				}
				if !_rules[ruleExpr]() {
					goto l0
				}
				if !_rules[rule_]() {
					goto l0
				}
				{
					position2 := position
					{
						position3, tokenIndex3 := position, tokenIndex
						if !matchDot() {
							goto l3
						}
						goto l0
					l3:
						position, tokenIndex = position3, tokenIndex3
					}
					add(ruleEOF, position2)
				}
				add(ruleQuery, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 PrimaryExpr <- <((Numbers Unit _ Action0) / (Severity Action2) / (Stage Action3) / (Id Action4) / ((&('(') (LPAR Expr RPAR Action7)) | (&('"') (String Action5)) | (&('\t' | '\n' | '\r' | ' ' | '.' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') (Numbers _ Action1)) | (&('/' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (UnquotedString Action6))))> */
		nil,
		/* 2 TupleExpr <- <(LPAR Expr Action8 (COMMA Expr Action9)* RPAR)> */
		func() bool {
			position5, tokenIndex5 := position, tokenIndex
			{
				position6 := position
				if !_rules[ruleLPAR]() {
					goto l5
				}
				if !_rules[ruleExpr]() {
					goto l5
				}
				{
					add(ruleAction8, position)
				}
			l8:
				{
					position9, tokenIndex9 := position, tokenIndex
					{
						position10 := position
						if buffer[position] != rune(',') {
							goto l9
						}
						position++
						if !_rules[rule_]() {
							goto l9
						}
						add(ruleCOMMA, position10)
					}
					if !_rules[ruleExpr]() {
						goto l9
					}
					{
						add(ruleAction9, position)
					}
					goto l8
				l9:
					position, tokenIndex = position9, tokenIndex9
				}
				if !_rules[ruleRPAR]() {
					goto l5
				}
				add(ruleTupleExpr, position6)
			}
			return true
		l5:
			position, tokenIndex = position5, tokenIndex5
			return false
		},
		/* 3 InExpr <- <(IN Action10 TupleExpr Action11)> */
		nil,
		/* 4 NotInExpr <- <(NOT IN Action12 TupleExpr Action13)> */
		nil,
		/* 5 PostfixExpr <- <(PrimaryExpr ((&('n') NotInExpr) | (&('i') InExpr) | (&('[') (LBRK Number _ RBRK Action15)) | (&('.') (DOT Id Action14)))*)> */
		nil,
		/* 6 UnaryExpr <- <(PostfixExpr / (BANG RelationalExpr Action16))> */
		func() bool {
			position15, tokenIndex15 := position, tokenIndex
			{
				position16 := position
				{
					position17, tokenIndex17 := position, tokenIndex
					{
						position19 := position
						{
							position20 := position
							{
								position21, tokenIndex21 := position, tokenIndex
								if !_rules[ruleNumbers]() {
									goto l22
								}
								{
									position23 := position
									{
										position24, tokenIndex24 := position, tokenIndex
										{
											position26 := position
											{
												switch buffer[position] {
												case 'g':
													{
														position28 := position
														{
															position29 := position
															if buffer[position] != rune('g') {
																goto l25
															}
															position++
															if buffer[position] != rune('b') {
																goto l25
															}
															position++
															add(rulePegText, position29)
														}
														{
															position30, tokenIndex30 := position, tokenIndex
															if !_rules[ruleIdChar]() {
																goto l30
															}
															goto l25
														l30:
															position, tokenIndex = position30, tokenIndex30
														}
														if !_rules[rule_]() {
															goto l25
														}
														add(ruleGB, position28)
													}
													break
												case 'm':
													{
														position31 := position
														{
															position32 := position
															if buffer[position] != rune('m') {
																goto l25
															}
															position++
															if buffer[position] != rune('b') {
																goto l25
															}
															position++
															add(rulePegText, position32)
														}
														{
															position33, tokenIndex33 := position, tokenIndex
															if !_rules[ruleIdChar]() {
																goto l33
															}
															goto l25
														l33:
															position, tokenIndex = position33, tokenIndex33
														}
														if !_rules[rule_]() {
															goto l25
														}
														add(ruleMB, position31)
													}
													break
												case 'k':
													{
														position34 := position
														{
															position35 := position
															if buffer[position] != rune('k') {
																goto l25
															}
															position++
															if buffer[position] != rune('b') {
																goto l25
															}
															position++
															add(rulePegText, position35)
														}
														{
															position36, tokenIndex36 := position, tokenIndex
															if !_rules[ruleIdChar]() {
																goto l36
															}
															goto l25
														l36:
															position, tokenIndex = position36, tokenIndex36
														}
														if !_rules[rule_]() {
															goto l25
														}
														add(ruleKB, position34)
													}
													break
												default:
													{
														position37 := position
														{
															position38 := position
															if buffer[position] != rune('b') {
																goto l25
															}
															position++
															add(rulePegText, position38)
														}
														{
															position39, tokenIndex39 := position, tokenIndex
															if !_rules[ruleIdChar]() {
																goto l39
															}
															goto l25
														l39:
															position, tokenIndex = position39, tokenIndex39
														}
														if !_rules[rule_]() {
															goto l25
														}
														add(ruleB, position37)
													}
													break
												}
											}

											add(ruleBytes, position26)
										}
										goto l24
									l25:
										position, tokenIndex = position24, tokenIndex24
										{
											position40 := position
											{
												position41, tokenIndex41 := position, tokenIndex
												{
													position43 := position
													{
														position44 := position
														if buffer[position] != rune('s') {
															goto l42
														}
														position++
														add(rulePegText, position44)
													}
													{
														position45, tokenIndex45 := position, tokenIndex
														if !_rules[ruleIdChar]() {
															goto l45
														}
														goto l42
													l45:
														position, tokenIndex = position45, tokenIndex45
													}
													if !_rules[rule_]() {
														goto l42
													}
													add(ruleS, position43)
												}
												goto l41
											l42:
												position, tokenIndex = position41, tokenIndex41
												{
													position46 := position
													{
														position47 := position
														if buffer[position] != rune('m') {
															goto l22
														}
														position++
														if buffer[position] != rune('s') {
															goto l22
														}
														position++
														add(rulePegText, position47)
													}
													{
														position48, tokenIndex48 := position, tokenIndex
														if !_rules[ruleIdChar]() {
															goto l48
														}
														goto l22
													l48:
														position, tokenIndex = position48, tokenIndex48
													}
													if !_rules[rule_]() {
														goto l22
													}
													add(ruleMS, position46)
												}
											}
										l41:
											add(ruleDuration, position40)
										}
									}
								l24:
									add(ruleUnit, position23)
								}
								if !_rules[rule_]() {
									goto l22
								}
								{
									add(ruleAction0, position)
								}
								goto l21
							l22:
								position, tokenIndex = position21, tokenIndex21
								{
									position51 := position
									{
										switch buffer[position] {
										case 'f':
											{
												position53 := position
												{
													position54 := position
													if buffer[position] != rune('f') {
														goto l50
													}
													position++
													if buffer[position] != rune('a') {
														goto l50
													}
													position++
													if buffer[position] != rune('t') {
														goto l50
													}
													position++
													if buffer[position] != rune('a') {
														goto l50
													}
													position++
													if buffer[position] != rune('l') {
														goto l50
													}
													position++
													add(rulePegText, position54)
												}
												{
													position55, tokenIndex55 := position, tokenIndex
													if !_rules[ruleIdChar]() {
														goto l55
													}
													goto l50
												l55:
													position, tokenIndex = position55, tokenIndex55
												}
												if !_rules[rule_]() {
													goto l50
												}
												add(ruleFATAL, position53)
											}
											break
										case 'e':
											{
												position56 := position
												{
													position57 := position
													if buffer[position] != rune('e') {
														goto l50
													}
													position++
													if buffer[position] != rune('r') {
														goto l50
													}
													position++
													if buffer[position] != rune('r') {
														goto l50
													}
													position++
													if buffer[position] != rune('o') {
														goto l50
													}
													position++
													if buffer[position] != rune('r') {
														goto l50
													}
													position++
													add(rulePegText, position57)
												}
												{
													position58, tokenIndex58 := position, tokenIndex
													if !_rules[ruleIdChar]() {
														goto l58
													}
													goto l50
												l58:
													position, tokenIndex = position58, tokenIndex58
												}
												if !_rules[rule_]() {
													goto l50
												}
												add(ruleERROR, position56)
											}
											break
										case 'w':
											{
												position59 := position
												{
													position60 := position
													if buffer[position] != rune('w') {
														goto l50
													}
													position++
													if buffer[position] != rune('a') {
														goto l50
													}
													position++
													if buffer[position] != rune('r') {
														goto l50
													}
													position++
													if buffer[position] != rune('n') {
														goto l50
													}
													position++
													add(rulePegText, position60)
												}
												{
													position61, tokenIndex61 := position, tokenIndex
													if !_rules[ruleIdChar]() {
														goto l61
													}
													goto l50
												l61:
													position, tokenIndex = position61, tokenIndex61
												}
												if !_rules[rule_]() {
													goto l50
												}
												add(ruleWARN, position59)
											}
											break
										case 'i':
											{
												position62 := position
												{
													position63 := position
													if buffer[position] != rune('i') {
														goto l50
													}
													position++
													if buffer[position] != rune('n') {
														goto l50
													}
													position++
													if buffer[position] != rune('f') {
														goto l50
													}
													position++
													if buffer[position] != rune('o') {
														goto l50
													}
													position++
													add(rulePegText, position63)
												}
												{
													position64, tokenIndex64 := position, tokenIndex
													if !_rules[ruleIdChar]() {
														goto l64
													}
													goto l50
												l64:
													position, tokenIndex = position64, tokenIndex64
												}
												if !_rules[rule_]() {
													goto l50
												}
												add(ruleINFO, position62)
											}
											break
										default:
											{
												position65 := position
												{
													position66 := position
													if buffer[position] != rune('d') {
														goto l50
													}
													position++
													if buffer[position] != rune('e') {
														goto l50
													}
													position++
													if buffer[position] != rune('b') {
														goto l50
													}
													position++
													if buffer[position] != rune('u') {
														goto l50
													}
													position++
													if buffer[position] != rune('g') {
														goto l50
													}
													position++
													add(rulePegText, position66)
												}
												{
													position67, tokenIndex67 := position, tokenIndex
													if !_rules[ruleIdChar]() {
														goto l67
													}
													goto l50
												l67:
													position, tokenIndex = position67, tokenIndex67
												}
												if !_rules[rule_]() {
													goto l50
												}
												add(ruleDEBUG, position65)
											}
											break
										}
									}

									add(ruleSeverity, position51)
								}
								{
									add(ruleAction2, position)
								}
								goto l21
							l50:
								position, tokenIndex = position21, tokenIndex21
								{
									position70 := position
									{
										switch buffer[position] {
										case 'p':
											{
												position72 := position
												{
													position73 := position
													if buffer[position] != rune('p') {
														goto l69
													}
													position++
													if buffer[position] != rune('r') {
														goto l69
													}
													position++
													if buffer[position] != rune('o') {
														goto l69
													}
													position++
													if buffer[position] != rune('d') {
														goto l69
													}
													position++
													if buffer[position] != rune('u') {
														goto l69
													}
													position++
													if buffer[position] != rune('c') {
														goto l69
													}
													position++
													if buffer[position] != rune('t') {
														goto l69
													}
													position++
													if buffer[position] != rune('i') {
														goto l69
													}
													position++
													if buffer[position] != rune('o') {
														goto l69
													}
													position++
													if buffer[position] != rune('n') {
														goto l69
													}
													position++
													add(rulePegText, position73)
												}
												{
													position74, tokenIndex74 := position, tokenIndex
													if !_rules[ruleIdChar]() {
														goto l74
													}
													goto l69
												l74:
													position, tokenIndex = position74, tokenIndex74
												}
												if !_rules[rule_]() {
													goto l69
												}
												add(rulePRODUCTION, position72)
											}
											break
										case 's':
											{
												position75 := position
												{
													position76 := position
													if buffer[position] != rune('s') {
														goto l69
													}
													position++
													if buffer[position] != rune('t') {
														goto l69
													}
													position++
													if buffer[position] != rune('a') {
														goto l69
													}
													position++
													if buffer[position] != rune('g') {
														goto l69
													}
													position++
													if buffer[position] != rune('i') {
														goto l69
													}
													position++
													if buffer[position] != rune('n') {
														goto l69
													}
													position++
													if buffer[position] != rune('g') {
														goto l69
													}
													position++
													add(rulePegText, position76)
												}
												{
													position77, tokenIndex77 := position, tokenIndex
													if !_rules[ruleIdChar]() {
														goto l77
													}
													goto l69
												l77:
													position, tokenIndex = position77, tokenIndex77
												}
												if !_rules[rule_]() {
													goto l69
												}
												add(ruleSTAGING, position75)
											}
											break
										default:
											{
												position78 := position
												{
													position79 := position
													if buffer[position] != rune('d') {
														goto l69
													}
													position++
													if buffer[position] != rune('e') {
														goto l69
													}
													position++
													if buffer[position] != rune('v') {
														goto l69
													}
													position++
													if buffer[position] != rune('e') {
														goto l69
													}
													position++
													if buffer[position] != rune('l') {
														goto l69
													}
													position++
													if buffer[position] != rune('o') {
														goto l69
													}
													position++
													if buffer[position] != rune('p') {
														goto l69
													}
													position++
													if buffer[position] != rune('m') {
														goto l69
													}
													position++
													if buffer[position] != rune('e') {
														goto l69
													}
													position++
													if buffer[position] != rune('n') {
														goto l69
													}
													position++
													if buffer[position] != rune('t') {
														goto l69
													}
													position++
													add(rulePegText, position79)
												}
												{
													position80, tokenIndex80 := position, tokenIndex
													if !_rules[ruleIdChar]() {
														goto l80
													}
													goto l69
												l80:
													position, tokenIndex = position80, tokenIndex80
												}
												if !_rules[rule_]() {
													goto l69
												}
												add(ruleDEVELOPMENT, position78)
											}
											break
										}
									}

									add(ruleStage, position70)
								}
								{
									add(ruleAction3, position)
								}
								goto l21
							l69:
								position, tokenIndex = position21, tokenIndex21
								if !_rules[ruleId]() {
									goto l82
								}
								{
									add(ruleAction4, position)
								}
								goto l21
							l82:
								position, tokenIndex = position21, tokenIndex21
								{
									switch buffer[position] {
									case '(':
										if !_rules[ruleLPAR]() {
											goto l18
										}
										if !_rules[ruleExpr]() {
											goto l18
										}
										if !_rules[ruleRPAR]() {
											goto l18
										}
										{
											add(ruleAction7, position)
										}
										break
									case '"':
										{
											position86 := position
											if buffer[position] != rune('"') {
												goto l18
											}
											position++
											{
												position87 := position
											l88:
												{
													position89, tokenIndex89 := position, tokenIndex
													{
														position90 := position
														{
															position91, tokenIndex91 := position, tokenIndex
															{
																position93 := position
																{
																	position94, tokenIndex94 := position, tokenIndex
																	{
																		position96 := position
																		if buffer[position] != rune('\\') {
																			goto l95
																		}
																		position++
																		{
																			switch buffer[position] {
																			case 'v':
																				if buffer[position] != rune('v') {
																					goto l95
																				}
																				position++
																				break
																			case 't':
																				if buffer[position] != rune('t') {
																					goto l95
																				}
																				position++
																				break
																			case 'r':
																				if buffer[position] != rune('r') {
																					goto l95
																				}
																				position++
																				break
																			case 'n':
																				if buffer[position] != rune('n') {
																					goto l95
																				}
																				position++
																				break
																			case 'f':
																				if buffer[position] != rune('f') {
																					goto l95
																				}
																				position++
																				break
																			case 'b':
																				if buffer[position] != rune('b') {
																					goto l95
																				}
																				position++
																				break
																			case 'a':
																				if buffer[position] != rune('a') {
																					goto l95
																				}
																				position++
																				break
																			case '\\':
																				if buffer[position] != rune('\\') {
																					goto l95
																				}
																				position++
																				break
																			case '?':
																				if buffer[position] != rune('?') {
																					goto l95
																				}
																				position++
																				break
																			case '"':
																				if buffer[position] != rune('"') {
																					goto l95
																				}
																				position++
																				break
																			default:
																				if buffer[position] != rune('\'') {
																					goto l95
																				}
																				position++
																				break
																			}
																		}

																		add(ruleSimpleEscape, position96)
																	}
																	goto l94
																l95:
																	position, tokenIndex = position94, tokenIndex94
																	{
																		position99 := position
																		if buffer[position] != rune('\\') {
																			goto l98
																		}
																		position++
																		if c := buffer[position]; c < rune('0') || c > rune('7') {
																			goto l98
																		}
																		position++
																		{
																			position100, tokenIndex100 := position, tokenIndex
																			if c := buffer[position]; c < rune('0') || c > rune('7') {
																				goto l100
																			}
																			position++
																			goto l101
																		l100:
																			position, tokenIndex = position100, tokenIndex100
																		}
																	l101:
																		{
																			position102, tokenIndex102 := position, tokenIndex
																			if c := buffer[position]; c < rune('0') || c > rune('7') {
																				goto l102
																			}
																			position++
																			goto l103
																		l102:
																			position, tokenIndex = position102, tokenIndex102
																		}
																	l103:
																		add(ruleOctalEscape, position99)
																	}
																	goto l94
																l98:
																	position, tokenIndex = position94, tokenIndex94
																	{
																		position105 := position
																		if buffer[position] != rune('\\') {
																			goto l104
																		}
																		position++
																		if buffer[position] != rune('x') {
																			goto l104
																		}
																		position++
																		if !_rules[ruleHexDigit]() {
																			goto l104
																		}
																	l106:
																		{
																			position107, tokenIndex107 := position, tokenIndex
																			if !_rules[ruleHexDigit]() {
																				goto l107
																			}
																			goto l106
																		l107:
																			position, tokenIndex = position107, tokenIndex107
																		}
																		add(ruleHexEscape, position105)
																	}
																	goto l94
																l104:
																	position, tokenIndex = position94, tokenIndex94
																	{
																		position108 := position
																		{
																			position109, tokenIndex109 := position, tokenIndex
																			if buffer[position] != rune('\\') {
																				goto l110
																			}
																			position++
																			if buffer[position] != rune('u') {
																				goto l110
																			}
																			position++
																			if !_rules[ruleHexQuad]() {
																				goto l110
																			}
																			goto l109
																		l110:
																			position, tokenIndex = position109, tokenIndex109
																			if buffer[position] != rune('\\') {
																				goto l92
																			}
																			position++
																			if buffer[position] != rune('U') {
																				goto l92
																			}
																			position++
																			if !_rules[ruleHexQuad]() {
																				goto l92
																			}
																			if !_rules[ruleHexQuad]() {
																				goto l92
																			}
																		}
																	l109:
																		add(ruleUniversalCharacter, position108)
																	}
																}
															l94:
																add(ruleEscape, position93)
															}
															goto l91
														l92:
															position, tokenIndex = position91, tokenIndex91
															{
																position111, tokenIndex111 := position, tokenIndex
																{
																	switch buffer[position] {
																	case '\\':
																		if buffer[position] != rune('\\') {
																			goto l111
																		}
																		position++
																		break
																	case '\n':
																		if buffer[position] != rune('\n') {
																			goto l111
																		}
																		position++
																		break
																	default:
																		if buffer[position] != rune('"') {
																			goto l111
																		}
																		position++
																		break
																	}
																}

																goto l89
															l111:
																position, tokenIndex = position111, tokenIndex111
															}
															if !matchDot() {
																goto l89
															}
														}
													l91:
														add(ruleStringChar, position90)
													}
													goto l88
												l89:
													position, tokenIndex = position89, tokenIndex89
												}
												add(rulePegText, position87)
											}
											if buffer[position] != rune('"') {
												goto l18
											}
											position++
											if !_rules[rule_]() {
												goto l18
											}
											add(ruleString, position86)
										}
										{
											add(ruleAction5, position)
										}
										break
									case '\t', '\n', '\r', ' ', '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if !_rules[ruleNumbers]() {
											goto l18
										}
										if !_rules[rule_]() {
											goto l18
										}
										{
											add(ruleAction1, position)
										}
										break
									default:
										{
											position115 := position
											{
												position116, tokenIndex116 := position, tokenIndex
												if !_rules[ruleKeyword]() {
													goto l116
												}
												goto l18
											l116:
												position, tokenIndex = position116, tokenIndex116
											}
											{
												position117 := position
												{
													position118 := position
													{
														switch buffer[position] {
														case '/', '_':
															{
																position120, tokenIndex120 := position, tokenIndex
																if buffer[position] != rune('/') {
																	goto l121
																}
																position++
																goto l120
															l121:
																position, tokenIndex = position120, tokenIndex120
																if buffer[position] != rune('_') {
																	goto l18
																}
																position++
															}
														l120:
															break
														case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
															if c := buffer[position]; c < rune('A') || c > rune('Z') {
																goto l18
															}
															position++
															break
														default:
															if c := buffer[position]; c < rune('a') || c > rune('z') {
																goto l18
															}
															position++
															break
														}
													}

													add(ruleUnquotedStringStartChar, position118)
												}
											l122:
												{
													position123, tokenIndex123 := position, tokenIndex
													{
														position124 := position
														{
															switch buffer[position] {
															case '/', '_':
																{
																	position126, tokenIndex126 := position, tokenIndex
																	if buffer[position] != rune('/') {
																		goto l127
																	}
																	position++
																	goto l126
																l127:
																	position, tokenIndex = position126, tokenIndex126
																	if buffer[position] != rune('_') {
																		goto l123
																	}
																	position++
																}
															l126:
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l123
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l123
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l123
																}
																position++
																break
															}
														}

														add(ruleUnquotedStringChar, position124)
													}
													goto l122
												l123:
													position, tokenIndex = position123, tokenIndex123
												}
												add(rulePegText, position117)
											}
											if !_rules[rule_]() {
												goto l18
											}
											add(ruleUnquotedString, position115)
										}
										{
											add(ruleAction6, position)
										}
										break
									}
								}

							}
						l21:
							add(rulePrimaryExpr, position20)
						}
					l129:
						{
							position130, tokenIndex130 := position, tokenIndex
							{
								switch buffer[position] {
								case 'n':
									{
										position132 := position
										if !_rules[ruleNOT]() {
											goto l130
										}
										if !_rules[ruleIN]() {
											goto l130
										}
										{
											add(ruleAction12, position)
										}
										if !_rules[ruleTupleExpr]() {
											goto l130
										}
										{
											add(ruleAction13, position)
										}
										add(ruleNotInExpr, position132)
									}
									break
								case 'i':
									{
										position135 := position
										if !_rules[ruleIN]() {
											goto l130
										}
										{
											add(ruleAction10, position)
										}
										if !_rules[ruleTupleExpr]() {
											goto l130
										}
										{
											add(ruleAction11, position)
										}
										add(ruleInExpr, position135)
									}
									break
								case '[':
									{
										position138 := position
										if buffer[position] != rune('[') {
											goto l130
										}
										position++
										if !_rules[rule_]() {
											goto l130
										}
										add(ruleLBRK, position138)
									}
									if !_rules[ruleNumber]() {
										goto l130
									}
									if !_rules[rule_]() {
										goto l130
									}
									{
										position139 := position
										if buffer[position] != rune(']') {
											goto l130
										}
										position++
										if !_rules[rule_]() {
											goto l130
										}
										add(ruleRBRK, position139)
									}
									{
										add(ruleAction15, position)
									}
									break
								default:
									{
										position141 := position
										if buffer[position] != rune('.') {
											goto l130
										}
										position++
										if !_rules[rule_]() {
											goto l130
										}
										add(ruleDOT, position141)
									}
									if !_rules[ruleId]() {
										goto l130
									}
									{
										add(ruleAction14, position)
									}
									break
								}
							}

							goto l129
						l130:
							position, tokenIndex = position130, tokenIndex130
						}
						add(rulePostfixExpr, position19)
					}
					goto l17
				l18:
					position, tokenIndex = position17, tokenIndex17
					{
						position143 := position
						if buffer[position] != rune('!') {
							goto l15
						}
						position++
						{
							position144, tokenIndex144 := position, tokenIndex
							if buffer[position] != rune('=') {
								goto l144
							}
							position++
							goto l15
						l144:
							position, tokenIndex = position144, tokenIndex144
						}
						if !_rules[rule_]() {
							goto l15
						}
						add(ruleBANG, position143)
					}
					if !_rules[ruleRelationalExpr]() {
						goto l15
					}
					{
						add(ruleAction16, position)
					}
				}
			l17:
				add(ruleUnaryExpr, position16)
			}
			return true
		l15:
			position, tokenIndex = position15, tokenIndex15
			return false
		},
		/* 7 RelationalExpr <- <(UnaryExpr ((GE UnaryExpr Action17) / (GT UnaryExpr Action18) / (LE UnaryExpr Action19) / (LT UnaryExpr Action20))*)> */
		func() bool {
			position146, tokenIndex146 := position, tokenIndex
			{
				position147 := position
				if !_rules[ruleUnaryExpr]() {
					goto l146
				}
			l148:
				{
					position149, tokenIndex149 := position, tokenIndex
					{
						position150, tokenIndex150 := position, tokenIndex
						{
							position152 := position
							if buffer[position] != rune('>') {
								goto l151
							}
							position++
							if buffer[position] != rune('=') {
								goto l151
							}
							position++
							if !_rules[rule_]() {
								goto l151
							}
							add(ruleGE, position152)
						}
						if !_rules[ruleUnaryExpr]() {
							goto l151
						}
						{
							add(ruleAction17, position)
						}
						goto l150
					l151:
						position, tokenIndex = position150, tokenIndex150
						{
							position155 := position
							if buffer[position] != rune('>') {
								goto l154
							}
							position++
							{
								position156, tokenIndex156 := position, tokenIndex
								if buffer[position] != rune('=') {
									goto l156
								}
								position++
								goto l154
							l156:
								position, tokenIndex = position156, tokenIndex156
							}
							if !_rules[rule_]() {
								goto l154
							}
							add(ruleGT, position155)
						}
						if !_rules[ruleUnaryExpr]() {
							goto l154
						}
						{
							add(ruleAction18, position)
						}
						goto l150
					l154:
						position, tokenIndex = position150, tokenIndex150
						{
							position159 := position
							if buffer[position] != rune('<') {
								goto l158
							}
							position++
							if buffer[position] != rune('=') {
								goto l158
							}
							position++
							if !_rules[rule_]() {
								goto l158
							}
							add(ruleLE, position159)
						}
						if !_rules[ruleUnaryExpr]() {
							goto l158
						}
						{
							add(ruleAction19, position)
						}
						goto l150
					l158:
						position, tokenIndex = position150, tokenIndex150
						{
							position161 := position
							if buffer[position] != rune('<') {
								goto l149
							}
							position++
							{
								position162, tokenIndex162 := position, tokenIndex
								if buffer[position] != rune('=') {
									goto l162
								}
								position++
								goto l149
							l162:
								position, tokenIndex = position162, tokenIndex162
							}
							if !_rules[rule_]() {
								goto l149
							}
							add(ruleLT, position161)
						}
						if !_rules[ruleUnaryExpr]() {
							goto l149
						}
						{
							add(ruleAction20, position)
						}
					}
				l150:
					goto l148
				l149:
					position, tokenIndex = position149, tokenIndex149
				}
				add(ruleRelationalExpr, position147)
			}
			return true
		l146:
			position, tokenIndex = position146, tokenIndex146
			return false
		},
		/* 8 EqualityExpr <- <(RelationalExpr ((EQEQ RelationalExpr Action21) / ((&('c') (CONTAINS RelationalExpr Action24)) | (&('=') (EQ RelationalExpr Action23)) | (&('!') (NE RelationalExpr Action22))))*)> */
		func() bool {
			position164, tokenIndex164 := position, tokenIndex
			{
				position165 := position
				if !_rules[ruleRelationalExpr]() {
					goto l164
				}
			l166:
				{
					position167, tokenIndex167 := position, tokenIndex
					{
						position168, tokenIndex168 := position, tokenIndex
						{
							position170 := position
							if buffer[position] != rune('=') {
								goto l169
							}
							position++
							if buffer[position] != rune('=') {
								goto l169
							}
							position++
							if !_rules[rule_]() {
								goto l169
							}
							add(ruleEQEQ, position170)
						}
						if !_rules[ruleRelationalExpr]() {
							goto l169
						}
						{
							add(ruleAction21, position)
						}
						goto l168
					l169:
						position, tokenIndex = position168, tokenIndex168
						{
							switch buffer[position] {
							case 'c':
								{
									position173 := position
									if buffer[position] != rune('c') {
										goto l167
									}
									position++
									if buffer[position] != rune('o') {
										goto l167
									}
									position++
									if buffer[position] != rune('n') {
										goto l167
									}
									position++
									if buffer[position] != rune('t') {
										goto l167
									}
									position++
									if buffer[position] != rune('a') {
										goto l167
									}
									position++
									if buffer[position] != rune('i') {
										goto l167
									}
									position++
									if buffer[position] != rune('n') {
										goto l167
									}
									position++
									if buffer[position] != rune('s') {
										goto l167
									}
									position++
									{
										position174, tokenIndex174 := position, tokenIndex
										if !_rules[ruleIdChar]() {
											goto l174
										}
										goto l167
									l174:
										position, tokenIndex = position174, tokenIndex174
									}
									if !_rules[rule_]() {
										goto l167
									}
									add(ruleCONTAINS, position173)
								}
								if !_rules[ruleRelationalExpr]() {
									goto l167
								}
								{
									add(ruleAction24, position)
								}
								break
							case '=':
								{
									position176 := position
									if buffer[position] != rune('=') {
										goto l167
									}
									position++
									if !_rules[rule_]() {
										goto l167
									}
									add(ruleEQ, position176)
								}
								if !_rules[ruleRelationalExpr]() {
									goto l167
								}
								{
									add(ruleAction23, position)
								}
								break
							default:
								{
									position178 := position
									if buffer[position] != rune('!') {
										goto l167
									}
									position++
									if buffer[position] != rune('=') {
										goto l167
									}
									position++
									if !_rules[rule_]() {
										goto l167
									}
									add(ruleNE, position178)
								}
								if !_rules[ruleRelationalExpr]() {
									goto l167
								}
								{
									add(ruleAction22, position)
								}
								break
							}
						}

					}
				l168:
					goto l166
				l167:
					position, tokenIndex = position167, tokenIndex167
				}
				add(ruleEqualityExpr, position165)
			}
			return true
		l164:
			position, tokenIndex = position164, tokenIndex164
			return false
		},
		/* 9 LogicalAndExpr <- <(EqualityExpr ((AND EqualityExpr Action25) / (ANDAND EqualityExpr Action26) / (_ EqualityExpr Action27))*)> */
		func() bool {
			position180, tokenIndex180 := position, tokenIndex
			{
				position181 := position
				if !_rules[ruleEqualityExpr]() {
					goto l180
				}
			l182:
				{
					position183, tokenIndex183 := position, tokenIndex
					{
						position184, tokenIndex184 := position, tokenIndex
						{
							position186 := position
							if buffer[position] != rune('a') {
								goto l185
							}
							position++
							if buffer[position] != rune('n') {
								goto l185
							}
							position++
							if buffer[position] != rune('d') {
								goto l185
							}
							position++
							{
								position187, tokenIndex187 := position, tokenIndex
								if !_rules[ruleIdChar]() {
									goto l187
								}
								goto l185
							l187:
								position, tokenIndex = position187, tokenIndex187
							}
							if !_rules[rule_]() {
								goto l185
							}
							add(ruleAND, position186)
						}
						if !_rules[ruleEqualityExpr]() {
							goto l185
						}
						{
							add(ruleAction25, position)
						}
						goto l184
					l185:
						position, tokenIndex = position184, tokenIndex184
						{
							position190 := position
							if buffer[position] != rune('&') {
								goto l189
							}
							position++
							if buffer[position] != rune('&') {
								goto l189
							}
							position++
							if !_rules[rule_]() {
								goto l189
							}
							add(ruleANDAND, position190)
						}
						if !_rules[ruleEqualityExpr]() {
							goto l189
						}
						{
							add(ruleAction26, position)
						}
						goto l184
					l189:
						position, tokenIndex = position184, tokenIndex184
						if !_rules[rule_]() {
							goto l183
						}
						if !_rules[ruleEqualityExpr]() {
							goto l183
						}
						{
							add(ruleAction27, position)
						}
					}
				l184:
					goto l182
				l183:
					position, tokenIndex = position183, tokenIndex183
				}
				add(ruleLogicalAndExpr, position181)
			}
			return true
		l180:
			position, tokenIndex = position180, tokenIndex180
			return false
		},
		/* 10 LogicalOrExpr <- <(LogicalAndExpr ((OR LogicalAndExpr Action28) / (OROR LogicalAndExpr Action29))*)> */
		func() bool {
			position193, tokenIndex193 := position, tokenIndex
			{
				position194 := position
				if !_rules[ruleLogicalAndExpr]() {
					goto l193
				}
			l195:
				{
					position196, tokenIndex196 := position, tokenIndex
					{
						position197, tokenIndex197 := position, tokenIndex
						{
							position199 := position
							if buffer[position] != rune('o') {
								goto l198
							}
							position++
							if buffer[position] != rune('r') {
								goto l198
							}
							position++
							{
								position200, tokenIndex200 := position, tokenIndex
								if !_rules[ruleIdChar]() {
									goto l200
								}
								goto l198
							l200:
								position, tokenIndex = position200, tokenIndex200
							}
							if !_rules[rule_]() {
								goto l198
							}
							add(ruleOR, position199)
						}
						if !_rules[ruleLogicalAndExpr]() {
							goto l198
						}
						{
							add(ruleAction28, position)
						}
						goto l197
					l198:
						position, tokenIndex = position197, tokenIndex197
						{
							position202 := position
							if buffer[position] != rune('|') {
								goto l196
							}
							position++
							if buffer[position] != rune('|') {
								goto l196
							}
							position++
							if !_rules[rule_]() {
								goto l196
							}
							add(ruleOROR, position202)
						}
						if !_rules[ruleLogicalAndExpr]() {
							goto l196
						}
						{
							add(ruleAction29, position)
						}
					}
				l197:
					goto l195
				l196:
					position, tokenIndex = position196, tokenIndex196
				}
				add(ruleLogicalOrExpr, position194)
			}
			return true
		l193:
			position, tokenIndex = position193, tokenIndex193
			return false
		},
		/* 11 LowNotExpr <- <(LogicalOrExpr / (NOT LogicalOrExpr Action30))> */
		nil,
		/* 12 Expr <- <LowNotExpr> */
		func() bool {
			position205, tokenIndex205 := position, tokenIndex
			{
				position206 := position
				{
					position207 := position
					{
						position208, tokenIndex208 := position, tokenIndex
						if !_rules[ruleLogicalOrExpr]() {
							goto l209
						}
						goto l208
					l209:
						position, tokenIndex = position208, tokenIndex208
						if !_rules[ruleNOT]() {
							goto l205
						}
						if !_rules[ruleLogicalOrExpr]() {
							goto l205
						}
						{
							add(ruleAction30, position)
						}
					}
				l208:
					add(ruleLowNotExpr, position207)
				}
				add(ruleExpr, position206)
			}
			return true
		l205:
			position, tokenIndex = position205, tokenIndex205
			return false
		},
		/* 13 String <- <('"' <StringChar*> '"' _)> */
		nil,
		/* 14 StringChar <- <(Escape / (!((&('\\') '\\') | (&('\n') '\n') | (&('"') '"')) .))> */
		nil,
		/* 15 UnquotedString <- <(!Keyword <(UnquotedStringStartChar UnquotedStringChar*)> _)> */
		nil,
		/* 16 UnquotedStringStartChar <- <((&('/' | '_') ('/' / '_')) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 17 UnquotedStringChar <- <((&('/' | '_') ('/' / '_')) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 18 Escape <- <(SimpleEscape / OctalEscape / HexEscape / UniversalCharacter)> */
		nil,
		/* 19 SimpleEscape <- <('\\' ((&('v') 'v') | (&('t') 't') | (&('r') 'r') | (&('n') 'n') | (&('f') 'f') | (&('b') 'b') | (&('a') 'a') | (&('\\') '\\') | (&('?') '?') | (&('"') '"') | (&('\'') '\'')))> */
		nil,
		/* 20 OctalEscape <- <('\\' [0-7] [0-7]? [0-7]?)> */
		nil,
		/* 21 HexEscape <- <('\\' 'x' HexDigit+)> */
		nil,
		/* 22 UniversalCharacter <- <(('\\' 'u' HexQuad) / ('\\' 'U' HexQuad HexQuad))> */
		nil,
		/* 23 HexQuad <- <(HexDigit HexDigit HexDigit HexDigit)> */
		func() bool {
			position221, tokenIndex221 := position, tokenIndex
			{
				position222 := position
				if !_rules[ruleHexDigit]() {
					goto l221
				}
				if !_rules[ruleHexDigit]() {
					goto l221
				}
				if !_rules[ruleHexDigit]() {
					goto l221
				}
				if !_rules[ruleHexDigit]() {
					goto l221
				}
				add(ruleHexQuad, position222)
			}
			return true
		l221:
			position, tokenIndex = position221, tokenIndex221
			return false
		},
		/* 24 HexDigit <- <((&('A' | 'B' | 'C' | 'D' | 'E' | 'F') [A-F]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f') [a-f]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]))> */
		func() bool {
			position223, tokenIndex223 := position, tokenIndex
			{
				position224 := position
				{
					switch buffer[position] {
					case 'A', 'B', 'C', 'D', 'E', 'F':
						if c := buffer[position]; c < rune('A') || c > rune('F') {
							goto l223
						}
						position++
						break
					case 'a', 'b', 'c', 'd', 'e', 'f':
						if c := buffer[position]; c < rune('a') || c > rune('f') {
							goto l223
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l223
						}
						position++
						break
					}
				}

				add(ruleHexDigit, position224)
			}
			return true
		l223:
			position, tokenIndex = position223, tokenIndex223
			return false
		},
		/* 25 Numbers <- <(Number Action31)> */
		func() bool {
			position226, tokenIndex226 := position, tokenIndex
			{
				position227 := position
				if !_rules[ruleNumber]() {
					goto l226
				}
				{
					add(ruleAction31, position)
				}
				add(ruleNumbers, position227)
			}
			return true
		l226:
			position, tokenIndex = position226, tokenIndex226
			return false
		},
		/* 26 Number <- <(<Float> / <Integer>)> */
		func() bool {
			{
				position230 := position
				{
					position231, tokenIndex231 := position, tokenIndex
					{
						position233 := position
						{
							position234 := position
							{
								position235, tokenIndex235 := position, tokenIndex
								{
									position237 := position
									{
										position238, tokenIndex238 := position, tokenIndex
									l240:
										{
											position241, tokenIndex241 := position, tokenIndex
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l241
											}
											position++
											goto l240
										l241:
											position, tokenIndex = position241, tokenIndex241
										}
										if buffer[position] != rune('.') {
											goto l239
										}
										position++
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l239
										}
										position++
									l242:
										{
											position243, tokenIndex243 := position, tokenIndex
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l243
											}
											position++
											goto l242
										l243:
											position, tokenIndex = position243, tokenIndex243
										}
										goto l238
									l239:
										position, tokenIndex = position238, tokenIndex238
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l236
										}
										position++
									l244:
										{
											position245, tokenIndex245 := position, tokenIndex
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l245
											}
											position++
											goto l244
										l245:
											position, tokenIndex = position245, tokenIndex245
										}
										if buffer[position] != rune('.') {
											goto l236
										}
										position++
									}
								l238:
									add(ruleFraction, position237)
								}
								{
									position246, tokenIndex246 := position, tokenIndex
									if !_rules[ruleExponent]() {
										goto l246
									}
									goto l247
								l246:
									position, tokenIndex = position246, tokenIndex246
								}
							l247:
								goto l235
							l236:
								position, tokenIndex = position235, tokenIndex235
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l232
								}
								position++
							l248:
								{
									position249, tokenIndex249 := position, tokenIndex
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l249
									}
									position++
									goto l248
								l249:
									position, tokenIndex = position249, tokenIndex249
								}
								if !_rules[ruleExponent]() {
									goto l232
								}
							}
						l235:
							add(ruleFloat, position234)
						}
						add(rulePegText, position233)
					}
					goto l231
				l232:
					position, tokenIndex = position231, tokenIndex231
					{
						position250 := position
						{
							position251 := position
						l252:
							{
								position253, tokenIndex253 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l253
								}
								position++
								goto l252
							l253:
								position, tokenIndex = position253, tokenIndex253
							}
							add(ruleInteger, position251)
						}
						add(rulePegText, position250)
					}
				}
			l231:
				add(ruleNumber, position230)
			}
			return true
		},
		/* 27 Integer <- <[0-9]*> */
		nil,
		/* 28 Float <- <((Fraction Exponent?) / ([0-9]+ Exponent))> */
		nil,
		/* 29 Fraction <- <(([0-9]* '.' [0-9]+) / ([0-9]+ '.'))> */
		nil,
		/* 30 Exponent <- <(('e' / 'E') ('+' / '-')? [0-9]+)> */
		func() bool {
			position257, tokenIndex257 := position, tokenIndex
			{
				position258 := position
				{
					position259, tokenIndex259 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l260
					}
					position++
					goto l259
				l260:
					position, tokenIndex = position259, tokenIndex259
					if buffer[position] != rune('E') {
						goto l257
					}
					position++
				}
			l259:
				{
					position261, tokenIndex261 := position, tokenIndex
					{
						position263, tokenIndex263 := position, tokenIndex
						if buffer[position] != rune('+') {
							goto l264
						}
						position++
						goto l263
					l264:
						position, tokenIndex = position263, tokenIndex263
						if buffer[position] != rune('-') {
							goto l261
						}
						position++
					}
				l263:
					goto l262
				l261:
					position, tokenIndex = position261, tokenIndex261
				}
			l262:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l257
				}
				position++
			l265:
				{
					position266, tokenIndex266 := position, tokenIndex
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l266
					}
					position++
					goto l265
				l266:
					position, tokenIndex = position266, tokenIndex266
				}
				add(ruleExponent, position258)
			}
			return true
		l257:
			position, tokenIndex = position257, tokenIndex257
			return false
		},
		/* 31 Stage <- <((&('p') PRODUCTION) | (&('s') STAGING) | (&('d') DEVELOPMENT))> */
		nil,
		/* 32 DEVELOPMENT <- <(<('d' 'e' 'v' 'e' 'l' 'o' 'p' 'm' 'e' 'n' 't')> !IdChar _)> */
		nil,
		/* 33 STAGING <- <(<('s' 't' 'a' 'g' 'i' 'n' 'g')> !IdChar _)> */
		nil,
		/* 34 PRODUCTION <- <(<('p' 'r' 'o' 'd' 'u' 'c' 't' 'i' 'o' 'n')> !IdChar _)> */
		nil,
		/* 35 Unit <- <(Bytes / Duration)> */
		nil,
		/* 36 Duration <- <(S / MS)> */
		nil,
		/* 37 S <- <(<'s'> !IdChar _)> */
		nil,
		/* 38 MS <- <(<('m' 's')> !IdChar _)> */
		nil,
		/* 39 Bytes <- <((&('g') GB) | (&('m') MB) | (&('k') KB) | (&('b') B))> */
		nil,
		/* 40 B <- <(<'b'> !IdChar _)> */
		nil,
		/* 41 KB <- <(<('k' 'b')> !IdChar _)> */
		nil,
		/* 42 MB <- <(<('m' 'b')> !IdChar _)> */
		nil,
		/* 43 GB <- <(<('g' 'b')> !IdChar _)> */
		nil,
		/* 44 Id <- <(!Keyword <(IdCharNoDigit IdChar*)> _)> */
		func() bool {
			position280, tokenIndex280 := position, tokenIndex
			{
				position281 := position
				{
					position282, tokenIndex282 := position, tokenIndex
					if !_rules[ruleKeyword]() {
						goto l282
					}
					goto l280
				l282:
					position, tokenIndex = position282, tokenIndex282
				}
				{
					position283 := position
					{
						position284 := position
						{
							switch buffer[position] {
							case '_':
								if buffer[position] != rune('_') {
									goto l280
								}
								position++
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l280
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l280
								}
								position++
								break
							}
						}

						add(ruleIdCharNoDigit, position284)
					}
				l286:
					{
						position287, tokenIndex287 := position, tokenIndex
						if !_rules[ruleIdChar]() {
							goto l287
						}
						goto l286
					l287:
						position, tokenIndex = position287, tokenIndex287
					}
					add(rulePegText, position283)
				}
				if !_rules[rule_]() {
					goto l280
				}
				add(ruleId, position281)
			}
			return true
		l280:
			position, tokenIndex = position280, tokenIndex280
			return false
		},
		/* 45 IdChar <- <((&('_') '_') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position288, tokenIndex288 := position, tokenIndex
			{
				position289 := position
				{
					switch buffer[position] {
					case '_':
						if buffer[position] != rune('_') {
							goto l288
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l288
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l288
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l288
						}
						position++
						break
					}
				}

				add(ruleIdChar, position289)
			}
			return true
		l288:
			position, tokenIndex = position288, tokenIndex288
			return false
		},
		/* 46 IdCharNoDigit <- <((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 47 Severity <- <((&('f') FATAL) | (&('e') ERROR) | (&('w') WARN) | (&('i') INFO) | (&('d') DEBUG))> */
		nil,
		/* 48 IN <- <('i' 'n' !IdChar _)> */
		func() bool {
			position293, tokenIndex293 := position, tokenIndex
			{
				position294 := position
				if buffer[position] != rune('i') {
					goto l293
				}
				position++
				if buffer[position] != rune('n') {
					goto l293
				}
				position++
				{
					position295, tokenIndex295 := position, tokenIndex
					if !_rules[ruleIdChar]() {
						goto l295
					}
					goto l293
				l295:
					position, tokenIndex = position295, tokenIndex295
				}
				if !_rules[rule_]() {
					goto l293
				}
				add(ruleIN, position294)
			}
			return true
		l293:
			position, tokenIndex = position293, tokenIndex293
			return false
		},
		/* 49 OR <- <('o' 'r' !IdChar _)> */
		nil,
		/* 50 AND <- <('a' 'n' 'd' !IdChar _)> */
		nil,
		/* 51 NOT <- <('n' 'o' 't' !IdChar _)> */
		func() bool {
			position298, tokenIndex298 := position, tokenIndex
			{
				position299 := position
				if buffer[position] != rune('n') {
					goto l298
				}
				position++
				if buffer[position] != rune('o') {
					goto l298
				}
				position++
				if buffer[position] != rune('t') {
					goto l298
				}
				position++
				{
					position300, tokenIndex300 := position, tokenIndex
					if !_rules[ruleIdChar]() {
						goto l300
					}
					goto l298
				l300:
					position, tokenIndex = position300, tokenIndex300
				}
				if !_rules[rule_]() {
					goto l298
				}
				add(ruleNOT, position299)
			}
			return true
		l298:
			position, tokenIndex = position298, tokenIndex298
			return false
		},
		/* 52 CONTAINS <- <('c' 'o' 'n' 't' 'a' 'i' 'n' 's' !IdChar _)> */
		nil,
		/* 53 DEBUG <- <(<('d' 'e' 'b' 'u' 'g')> !IdChar _)> */
		nil,
		/* 54 INFO <- <(<('i' 'n' 'f' 'o')> !IdChar _)> */
		nil,
		/* 55 WARN <- <(<('w' 'a' 'r' 'n')> !IdChar _)> */
		nil,
		/* 56 ERROR <- <(<('e' 'r' 'r' 'o' 'r')> !IdChar _)> */
		nil,
		/* 57 FATAL <- <(<('f' 'a' 't' 'a' 'l')> !IdChar _)> */
		nil,
		/* 58 Keyword <- <((('s' 't' 'a' 'g' 'i' 'n' 'g') / ('d' 'e' 'v' 'e' 'l' 'o' 'p' 'm' 'e' 'n' 't') / ('i' 'n' 'f' 'o') / ('m' 'b') / ((&('s') 's') | (&('m') ('m' 's')) | (&('b') 'b') | (&('k') ('k' 'b')) | (&('g') ('g' 'b')) | (&('i') ('i' 'n')) | (&('f') ('f' 'a' 't' 'a' 'l')) | (&('e') ('e' 'r' 'r' 'o' 'r')) | (&('w') ('w' 'a' 'r' 'n')) | (&('d') ('d' 'e' 'b' 'u' 'g')) | (&('c') ('c' 'o' 'n' 't' 'a' 'i' 'n' 's')) | (&('n') ('n' 'o' 't')) | (&('a') ('a' 'n' 'd')) | (&('o') ('o' 'r')) | (&('p') ('p' 'r' 'o' 'd' 'u' 'c' 't' 'i' 'o' 'n')))) !IdChar)> */
		func() bool {
			position307, tokenIndex307 := position, tokenIndex
			{
				position308 := position
				{
					position309, tokenIndex309 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l310
					}
					position++
					if buffer[position] != rune('t') {
						goto l310
					}
					position++
					if buffer[position] != rune('a') {
						goto l310
					}
					position++
					if buffer[position] != rune('g') {
						goto l310
					}
					position++
					if buffer[position] != rune('i') {
						goto l310
					}
					position++
					if buffer[position] != rune('n') {
						goto l310
					}
					position++
					if buffer[position] != rune('g') {
						goto l310
					}
					position++
					goto l309
				l310:
					position, tokenIndex = position309, tokenIndex309
					if buffer[position] != rune('d') {
						goto l311
					}
					position++
					if buffer[position] != rune('e') {
						goto l311
					}
					position++
					if buffer[position] != rune('v') {
						goto l311
					}
					position++
					if buffer[position] != rune('e') {
						goto l311
					}
					position++
					if buffer[position] != rune('l') {
						goto l311
					}
					position++
					if buffer[position] != rune('o') {
						goto l311
					}
					position++
					if buffer[position] != rune('p') {
						goto l311
					}
					position++
					if buffer[position] != rune('m') {
						goto l311
					}
					position++
					if buffer[position] != rune('e') {
						goto l311
					}
					position++
					if buffer[position] != rune('n') {
						goto l311
					}
					position++
					if buffer[position] != rune('t') {
						goto l311
					}
					position++
					goto l309
				l311:
					position, tokenIndex = position309, tokenIndex309
					if buffer[position] != rune('i') {
						goto l312
					}
					position++
					if buffer[position] != rune('n') {
						goto l312
					}
					position++
					if buffer[position] != rune('f') {
						goto l312
					}
					position++
					if buffer[position] != rune('o') {
						goto l312
					}
					position++
					goto l309
				l312:
					position, tokenIndex = position309, tokenIndex309
					if buffer[position] != rune('m') {
						goto l313
					}
					position++
					if buffer[position] != rune('b') {
						goto l313
					}
					position++
					goto l309
				l313:
					position, tokenIndex = position309, tokenIndex309
					{
						switch buffer[position] {
						case 's':
							if buffer[position] != rune('s') {
								goto l307
							}
							position++
							break
						case 'm':
							if buffer[position] != rune('m') {
								goto l307
							}
							position++
							if buffer[position] != rune('s') {
								goto l307
							}
							position++
							break
						case 'b':
							if buffer[position] != rune('b') {
								goto l307
							}
							position++
							break
						case 'k':
							if buffer[position] != rune('k') {
								goto l307
							}
							position++
							if buffer[position] != rune('b') {
								goto l307
							}
							position++
							break
						case 'g':
							if buffer[position] != rune('g') {
								goto l307
							}
							position++
							if buffer[position] != rune('b') {
								goto l307
							}
							position++
							break
						case 'i':
							if buffer[position] != rune('i') {
								goto l307
							}
							position++
							if buffer[position] != rune('n') {
								goto l307
							}
							position++
							break
						case 'f':
							if buffer[position] != rune('f') {
								goto l307
							}
							position++
							if buffer[position] != rune('a') {
								goto l307
							}
							position++
							if buffer[position] != rune('t') {
								goto l307
							}
							position++
							if buffer[position] != rune('a') {
								goto l307
							}
							position++
							if buffer[position] != rune('l') {
								goto l307
							}
							position++
							break
						case 'e':
							if buffer[position] != rune('e') {
								goto l307
							}
							position++
							if buffer[position] != rune('r') {
								goto l307
							}
							position++
							if buffer[position] != rune('r') {
								goto l307
							}
							position++
							if buffer[position] != rune('o') {
								goto l307
							}
							position++
							if buffer[position] != rune('r') {
								goto l307
							}
							position++
							break
						case 'w':
							if buffer[position] != rune('w') {
								goto l307
							}
							position++
							if buffer[position] != rune('a') {
								goto l307
							}
							position++
							if buffer[position] != rune('r') {
								goto l307
							}
							position++
							if buffer[position] != rune('n') {
								goto l307
							}
							position++
							break
						case 'd':
							if buffer[position] != rune('d') {
								goto l307
							}
							position++
							if buffer[position] != rune('e') {
								goto l307
							}
							position++
							if buffer[position] != rune('b') {
								goto l307
							}
							position++
							if buffer[position] != rune('u') {
								goto l307
							}
							position++
							if buffer[position] != rune('g') {
								goto l307
							}
							position++
							break
						case 'c':
							if buffer[position] != rune('c') {
								goto l307
							}
							position++
							if buffer[position] != rune('o') {
								goto l307
							}
							position++
							if buffer[position] != rune('n') {
								goto l307
							}
							position++
							if buffer[position] != rune('t') {
								goto l307
							}
							position++
							if buffer[position] != rune('a') {
								goto l307
							}
							position++
							if buffer[position] != rune('i') {
								goto l307
							}
							position++
							if buffer[position] != rune('n') {
								goto l307
							}
							position++
							if buffer[position] != rune('s') {
								goto l307
							}
							position++
							break
						case 'n':
							if buffer[position] != rune('n') {
								goto l307
							}
							position++
							if buffer[position] != rune('o') {
								goto l307
							}
							position++
							if buffer[position] != rune('t') {
								goto l307
							}
							position++
							break
						case 'a':
							if buffer[position] != rune('a') {
								goto l307
							}
							position++
							if buffer[position] != rune('n') {
								goto l307
							}
							position++
							if buffer[position] != rune('d') {
								goto l307
							}
							position++
							break
						case 'o':
							if buffer[position] != rune('o') {
								goto l307
							}
							position++
							if buffer[position] != rune('r') {
								goto l307
							}
							position++
							break
						default:
							if buffer[position] != rune('p') {
								goto l307
							}
							position++
							if buffer[position] != rune('r') {
								goto l307
							}
							position++
							if buffer[position] != rune('o') {
								goto l307
							}
							position++
							if buffer[position] != rune('d') {
								goto l307
							}
							position++
							if buffer[position] != rune('u') {
								goto l307
							}
							position++
							if buffer[position] != rune('c') {
								goto l307
							}
							position++
							if buffer[position] != rune('t') {
								goto l307
							}
							position++
							if buffer[position] != rune('i') {
								goto l307
							}
							position++
							if buffer[position] != rune('o') {
								goto l307
							}
							position++
							if buffer[position] != rune('n') {
								goto l307
							}
							position++
							break
						}
					}

				}
			l309:
				{
					position315, tokenIndex315 := position, tokenIndex
					if !_rules[ruleIdChar]() {
						goto l315
					}
					goto l307
				l315:
					position, tokenIndex = position315, tokenIndex315
				}
				add(ruleKeyword, position308)
			}
			return true
		l307:
			position, tokenIndex = position307, tokenIndex307
			return false
		},
		/* 59 EQ <- <('=' _)> */
		nil,
		/* 60 LBRK <- <('[' _)> */
		nil,
		/* 61 RBRK <- <(']' _)> */
		nil,
		/* 62 LPAR <- <('(' _)> */
		func() bool {
			position319, tokenIndex319 := position, tokenIndex
			{
				position320 := position
				if buffer[position] != rune('(') {
					goto l319
				}
				position++
				if !_rules[rule_]() {
					goto l319
				}
				add(ruleLPAR, position320)
			}
			return true
		l319:
			position, tokenIndex = position319, tokenIndex319
			return false
		},
		/* 63 RPAR <- <(')' _)> */
		func() bool {
			position321, tokenIndex321 := position, tokenIndex
			{
				position322 := position
				if buffer[position] != rune(')') {
					goto l321
				}
				position++
				if !_rules[rule_]() {
					goto l321
				}
				add(ruleRPAR, position322)
			}
			return true
		l321:
			position, tokenIndex = position321, tokenIndex321
			return false
		},
		/* 64 DOT <- <('.' _)> */
		nil,
		/* 65 BANG <- <('!' !'=' _)> */
		nil,
		/* 66 LT <- <('<' !'=' _)> */
		nil,
		/* 67 GT <- <('>' !'=' _)> */
		nil,
		/* 68 LE <- <('<' '=' _)> */
		nil,
		/* 69 EQEQ <- <('=' '=' _)> */
		nil,
		/* 70 GE <- <('>' '=' _)> */
		nil,
		/* 71 NE <- <('!' '=' _)> */
		nil,
		/* 72 ANDAND <- <('&' '&' _)> */
		nil,
		/* 73 OROR <- <('|' '|' _)> */
		nil,
		/* 74 COMMA <- <(',' _)> */
		nil,
		/* 75 _ <- <Whitespace*> */
		func() bool {
			{
				position335 := position
			l336:
				{
					position337, tokenIndex337 := position, tokenIndex
					{
						position338 := position
						{
							switch buffer[position] {
							case '\t':
								if buffer[position] != rune('\t') {
									goto l337
								}
								position++
								break
							case ' ':
								if buffer[position] != rune(' ') {
									goto l337
								}
								position++
								break
							default:
								{
									position340 := position
									{
										position341, tokenIndex341 := position, tokenIndex
										if buffer[position] != rune('\r') {
											goto l342
										}
										position++
										if buffer[position] != rune('\n') {
											goto l342
										}
										position++
										goto l341
									l342:
										position, tokenIndex = position341, tokenIndex341
										if buffer[position] != rune('\n') {
											goto l343
										}
										position++
										goto l341
									l343:
										position, tokenIndex = position341, tokenIndex341
										if buffer[position] != rune('\r') {
											goto l337
										}
										position++
									}
								l341:
									add(ruleEOL, position340)
								}
								break
							}
						}

						add(ruleWhitespace, position338)
					}
					goto l336
				l337:
					position, tokenIndex = position337, tokenIndex337
				}
				add(rule_, position335)
			}
			return true
		},
		/* 76 Whitespace <- <((&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') EOL))> */
		nil,
		/* 77 EOL <- <(('\r' '\n') / '\n' / '\r')> */
		nil,
		/* 78 EOF <- <!.> */
		nil,
		/* 80 Action0 <- <{ p.AddNumber(text) }> */
		nil,
		/* 81 Action1 <- <{ p.AddNumber("")   }> */
		nil,
		/* 82 Action2 <- <{ p.AddLevel(text)  }> */
		nil,
		/* 83 Action3 <- <{ p.AddStage(text)  }> */
		nil,
		/* 84 Action4 <- <{ p.AddField(text)  }> */
		nil,
		/* 85 Action5 <- <{ p.AddString(text) }> */
		nil,
		/* 86 Action6 <- <{ p.AddString(text) }> */
		nil,
		/* 87 Action7 <- <{ p.AddExpr()       }> */
		nil,
		/* 88 Action8 <- <{ p.AddTupleValue() }> */
		nil,
		/* 89 Action9 <- <{ p.AddTupleValue() }> */
		nil,
		/* 90 Action10 <- <{ p.AddTuple() }> */
		nil,
		/* 91 Action11 <- <{ p.AddBinary(ast.IN) }> */
		nil,
		/* 92 Action12 <- <{ p.AddTuple() }> */
		nil,
		/* 93 Action13 <- <{ p.AddBinary(ast.IN); p.AddUnary(ast.LNOT) }> */
		nil,
		/* 94 Action14 <- <{ p.AddMember(text)    }> */
		nil,
		/* 95 Action15 <- <{ p.AddSubscript(text) }> */
		nil,
		/* 96 Action16 <- <{ p.AddUnary(ast.NOT) }> */
		nil,
		/* 97 Action17 <- <{ p.AddBinary(ast.GE) }> */
		nil,
		/* 98 Action18 <- <{ p.AddBinary(ast.GT) }> */
		nil,
		/* 99 Action19 <- <{ p.AddBinary(ast.LE) }> */
		nil,
		/* 100 Action20 <- <{ p.AddBinary(ast.LT) }> */
		nil,
		/* 101 Action21 <- <{ p.AddBinary(ast.EQ)   }> */
		nil,
		/* 102 Action22 <- <{ p.AddBinary(ast.NE)   }> */
		nil,
		/* 103 Action23 <- <{ p.AddBinary(ast.EQ)   }> */
		nil,
		/* 104 Action24 <- <{ p.AddBinaryContains() }> */
		nil,
		/* 105 Action25 <- <{ p.AddBinary(ast.AND) }> */
		nil,
		/* 106 Action26 <- <{ p.AddBinary(ast.AND) }> */
		nil,
		/* 107 Action27 <- <{ p.AddBinary(ast.AND) }> */
		nil,
		/* 108 Action28 <- <{ p.AddBinary(ast.OR) }> */
		nil,
		/* 109 Action29 <- <{ p.AddBinary(ast.OR) }> */
		nil,
		/* 110 Action30 <- <{ p.AddUnary(ast.LNOT) }> */
		nil,
		nil,
		/* 112 Action31 <- <{ p.SetNumber(text) }> */
		nil,
	}
	p.rules = _rules
}
