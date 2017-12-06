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
												position89 := position
											l90:
												{
													position91, tokenIndex91 := position, tokenIndex
													{
														position92 := position
														{
															position93, tokenIndex93 := position, tokenIndex
															{
																position95 := position
																{
																	position96, tokenIndex96 := position, tokenIndex
																	{
																		position98 := position
																		if buffer[position] != rune('\\') {
																			goto l97
																		}
																		position++
																		{
																			switch buffer[position] {
																			case 'v':
																				if buffer[position] != rune('v') {
																					goto l97
																				}
																				position++
																				break
																			case 't':
																				if buffer[position] != rune('t') {
																					goto l97
																				}
																				position++
																				break
																			case 'r':
																				if buffer[position] != rune('r') {
																					goto l97
																				}
																				position++
																				break
																			case 'n':
																				if buffer[position] != rune('n') {
																					goto l97
																				}
																				position++
																				break
																			case 'f':
																				if buffer[position] != rune('f') {
																					goto l97
																				}
																				position++
																				break
																			case 'b':
																				if buffer[position] != rune('b') {
																					goto l97
																				}
																				position++
																				break
																			case 'a':
																				if buffer[position] != rune('a') {
																					goto l97
																				}
																				position++
																				break
																			case '\\':
																				if buffer[position] != rune('\\') {
																					goto l97
																				}
																				position++
																				break
																			case '?':
																				if buffer[position] != rune('?') {
																					goto l97
																				}
																				position++
																				break
																			case '"':
																				if buffer[position] != rune('"') {
																					goto l97
																				}
																				position++
																				break
																			default:
																				if buffer[position] != rune('\'') {
																					goto l97
																				}
																				position++
																				break
																			}
																		}

																		add(ruleSimpleEscape, position98)
																	}
																	goto l96
																l97:
																	position, tokenIndex = position96, tokenIndex96
																	{
																		position101 := position
																		if buffer[position] != rune('\\') {
																			goto l100
																		}
																		position++
																		if c := buffer[position]; c < rune('0') || c > rune('7') {
																			goto l100
																		}
																		position++
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
																		{
																			position104, tokenIndex104 := position, tokenIndex
																			if c := buffer[position]; c < rune('0') || c > rune('7') {
																				goto l104
																			}
																			position++
																			goto l105
																		l104:
																			position, tokenIndex = position104, tokenIndex104
																		}
																	l105:
																		add(ruleOctalEscape, position101)
																	}
																	goto l96
																l100:
																	position, tokenIndex = position96, tokenIndex96
																	{
																		position107 := position
																		if buffer[position] != rune('\\') {
																			goto l106
																		}
																		position++
																		if buffer[position] != rune('x') {
																			goto l106
																		}
																		position++
																		if !_rules[ruleHexDigit]() {
																			goto l106
																		}
																	l108:
																		{
																			position109, tokenIndex109 := position, tokenIndex
																			if !_rules[ruleHexDigit]() {
																				goto l109
																			}
																			goto l108
																		l109:
																			position, tokenIndex = position109, tokenIndex109
																		}
																		add(ruleHexEscape, position107)
																	}
																	goto l96
																l106:
																	position, tokenIndex = position96, tokenIndex96
																	{
																		position110 := position
																		{
																			position111, tokenIndex111 := position, tokenIndex
																			if buffer[position] != rune('\\') {
																				goto l112
																			}
																			position++
																			if buffer[position] != rune('u') {
																				goto l112
																			}
																			position++
																			if !_rules[ruleHexQuad]() {
																				goto l112
																			}
																			goto l111
																		l112:
																			position, tokenIndex = position111, tokenIndex111
																			if buffer[position] != rune('\\') {
																				goto l94
																			}
																			position++
																			if buffer[position] != rune('U') {
																				goto l94
																			}
																			position++
																			if !_rules[ruleHexQuad]() {
																				goto l94
																			}
																			if !_rules[ruleHexQuad]() {
																				goto l94
																			}
																		}
																	l111:
																		add(ruleUniversalCharacter, position110)
																	}
																}
															l96:
																add(ruleEscape, position95)
															}
															goto l93
														l94:
															position, tokenIndex = position93, tokenIndex93
															{
																position113, tokenIndex113 := position, tokenIndex
																{
																	switch buffer[position] {
																	case '\\':
																		if buffer[position] != rune('\\') {
																			goto l113
																		}
																		position++
																		break
																	case '\n':
																		if buffer[position] != rune('\n') {
																			goto l113
																		}
																		position++
																		break
																	default:
																		if buffer[position] != rune('"') {
																			goto l113
																		}
																		position++
																		break
																	}
																}

																goto l91
															l113:
																position, tokenIndex = position113, tokenIndex113
															}
															if !matchDot() {
																goto l91
															}
														}
													l93:
														add(ruleStringChar, position92)
													}
													goto l90
												l91:
													position, tokenIndex = position91, tokenIndex91
												}
												add(rulePegText, position89)
											}
											if buffer[position] != rune('"') {
												goto l18
											}
											position++
											if !_rules[rule_]() {
												goto l18
											}
										l87:
											{
												position88, tokenIndex88 := position, tokenIndex
												if buffer[position] != rune('"') {
													goto l88
												}
												position++
												{
													position115 := position
												l116:
													{
														position117, tokenIndex117 := position, tokenIndex
														{
															position118 := position
															{
																position119, tokenIndex119 := position, tokenIndex
																{
																	position121 := position
																	{
																		position122, tokenIndex122 := position, tokenIndex
																		{
																			position124 := position
																			if buffer[position] != rune('\\') {
																				goto l123
																			}
																			position++
																			{
																				switch buffer[position] {
																				case 'v':
																					if buffer[position] != rune('v') {
																						goto l123
																					}
																					position++
																					break
																				case 't':
																					if buffer[position] != rune('t') {
																						goto l123
																					}
																					position++
																					break
																				case 'r':
																					if buffer[position] != rune('r') {
																						goto l123
																					}
																					position++
																					break
																				case 'n':
																					if buffer[position] != rune('n') {
																						goto l123
																					}
																					position++
																					break
																				case 'f':
																					if buffer[position] != rune('f') {
																						goto l123
																					}
																					position++
																					break
																				case 'b':
																					if buffer[position] != rune('b') {
																						goto l123
																					}
																					position++
																					break
																				case 'a':
																					if buffer[position] != rune('a') {
																						goto l123
																					}
																					position++
																					break
																				case '\\':
																					if buffer[position] != rune('\\') {
																						goto l123
																					}
																					position++
																					break
																				case '?':
																					if buffer[position] != rune('?') {
																						goto l123
																					}
																					position++
																					break
																				case '"':
																					if buffer[position] != rune('"') {
																						goto l123
																					}
																					position++
																					break
																				default:
																					if buffer[position] != rune('\'') {
																						goto l123
																					}
																					position++
																					break
																				}
																			}

																			add(ruleSimpleEscape, position124)
																		}
																		goto l122
																	l123:
																		position, tokenIndex = position122, tokenIndex122
																		{
																			position127 := position
																			if buffer[position] != rune('\\') {
																				goto l126
																			}
																			position++
																			if c := buffer[position]; c < rune('0') || c > rune('7') {
																				goto l126
																			}
																			position++
																			{
																				position128, tokenIndex128 := position, tokenIndex
																				if c := buffer[position]; c < rune('0') || c > rune('7') {
																					goto l128
																				}
																				position++
																				goto l129
																			l128:
																				position, tokenIndex = position128, tokenIndex128
																			}
																		l129:
																			{
																				position130, tokenIndex130 := position, tokenIndex
																				if c := buffer[position]; c < rune('0') || c > rune('7') {
																					goto l130
																				}
																				position++
																				goto l131
																			l130:
																				position, tokenIndex = position130, tokenIndex130
																			}
																		l131:
																			add(ruleOctalEscape, position127)
																		}
																		goto l122
																	l126:
																		position, tokenIndex = position122, tokenIndex122
																		{
																			position133 := position
																			if buffer[position] != rune('\\') {
																				goto l132
																			}
																			position++
																			if buffer[position] != rune('x') {
																				goto l132
																			}
																			position++
																			if !_rules[ruleHexDigit]() {
																				goto l132
																			}
																		l134:
																			{
																				position135, tokenIndex135 := position, tokenIndex
																				if !_rules[ruleHexDigit]() {
																					goto l135
																				}
																				goto l134
																			l135:
																				position, tokenIndex = position135, tokenIndex135
																			}
																			add(ruleHexEscape, position133)
																		}
																		goto l122
																	l132:
																		position, tokenIndex = position122, tokenIndex122
																		{
																			position136 := position
																			{
																				position137, tokenIndex137 := position, tokenIndex
																				if buffer[position] != rune('\\') {
																					goto l138
																				}
																				position++
																				if buffer[position] != rune('u') {
																					goto l138
																				}
																				position++
																				if !_rules[ruleHexQuad]() {
																					goto l138
																				}
																				goto l137
																			l138:
																				position, tokenIndex = position137, tokenIndex137
																				if buffer[position] != rune('\\') {
																					goto l120
																				}
																				position++
																				if buffer[position] != rune('U') {
																					goto l120
																				}
																				position++
																				if !_rules[ruleHexQuad]() {
																					goto l120
																				}
																				if !_rules[ruleHexQuad]() {
																					goto l120
																				}
																			}
																		l137:
																			add(ruleUniversalCharacter, position136)
																		}
																	}
																l122:
																	add(ruleEscape, position121)
																}
																goto l119
															l120:
																position, tokenIndex = position119, tokenIndex119
																{
																	position139, tokenIndex139 := position, tokenIndex
																	{
																		switch buffer[position] {
																		case '\\':
																			if buffer[position] != rune('\\') {
																				goto l139
																			}
																			position++
																			break
																		case '\n':
																			if buffer[position] != rune('\n') {
																				goto l139
																			}
																			position++
																			break
																		default:
																			if buffer[position] != rune('"') {
																				goto l139
																			}
																			position++
																			break
																		}
																	}

																	goto l117
																l139:
																	position, tokenIndex = position139, tokenIndex139
																}
																if !matchDot() {
																	goto l117
																}
															}
														l119:
															add(ruleStringChar, position118)
														}
														goto l116
													l117:
														position, tokenIndex = position117, tokenIndex117
													}
													add(rulePegText, position115)
												}
												if buffer[position] != rune('"') {
													goto l88
												}
												position++
												if !_rules[rule_]() {
													goto l88
												}
												goto l87
											l88:
												position, tokenIndex = position88, tokenIndex88
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
											position143 := position
											{
												position144, tokenIndex144 := position, tokenIndex
												if !_rules[ruleKeyword]() {
													goto l144
												}
												goto l18
											l144:
												position, tokenIndex = position144, tokenIndex144
											}
											{
												position145 := position
												{
													position146 := position
													{
														switch buffer[position] {
														case '/', '_':
															{
																position148, tokenIndex148 := position, tokenIndex
																if buffer[position] != rune('/') {
																	goto l149
																}
																position++
																goto l148
															l149:
																position, tokenIndex = position148, tokenIndex148
																if buffer[position] != rune('_') {
																	goto l18
																}
																position++
															}
														l148:
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

													add(ruleUnquotedStringStartChar, position146)
												}
											l150:
												{
													position151, tokenIndex151 := position, tokenIndex
													{
														position152 := position
														{
															switch buffer[position] {
															case '/', '_':
																{
																	position154, tokenIndex154 := position, tokenIndex
																	if buffer[position] != rune('/') {
																		goto l155
																	}
																	position++
																	goto l154
																l155:
																	position, tokenIndex = position154, tokenIndex154
																	if buffer[position] != rune('_') {
																		goto l151
																	}
																	position++
																}
															l154:
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l151
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l151
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l151
																}
																position++
																break
															}
														}

														add(ruleUnquotedStringChar, position152)
													}
													goto l150
												l151:
													position, tokenIndex = position151, tokenIndex151
												}
												add(rulePegText, position145)
											}
											if !_rules[rule_]() {
												goto l18
											}
											add(ruleUnquotedString, position143)
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
					l157:
						{
							position158, tokenIndex158 := position, tokenIndex
							{
								switch buffer[position] {
								case 'n':
									{
										position160 := position
										if !_rules[ruleNOT]() {
											goto l158
										}
										if !_rules[ruleIN]() {
											goto l158
										}
										{
											add(ruleAction12, position)
										}
										if !_rules[ruleTupleExpr]() {
											goto l158
										}
										{
											add(ruleAction13, position)
										}
										add(ruleNotInExpr, position160)
									}
									break
								case 'i':
									{
										position163 := position
										if !_rules[ruleIN]() {
											goto l158
										}
										{
											add(ruleAction10, position)
										}
										if !_rules[ruleTupleExpr]() {
											goto l158
										}
										{
											add(ruleAction11, position)
										}
										add(ruleInExpr, position163)
									}
									break
								case '[':
									{
										position166 := position
										if buffer[position] != rune('[') {
											goto l158
										}
										position++
										if !_rules[rule_]() {
											goto l158
										}
										add(ruleLBRK, position166)
									}
									if !_rules[ruleNumber]() {
										goto l158
									}
									if !_rules[rule_]() {
										goto l158
									}
									{
										position167 := position
										if buffer[position] != rune(']') {
											goto l158
										}
										position++
										if !_rules[rule_]() {
											goto l158
										}
										add(ruleRBRK, position167)
									}
									{
										add(ruleAction15, position)
									}
									break
								default:
									{
										position169 := position
										if buffer[position] != rune('.') {
											goto l158
										}
										position++
										if !_rules[rule_]() {
											goto l158
										}
										add(ruleDOT, position169)
									}
									if !_rules[ruleId]() {
										goto l158
									}
									{
										add(ruleAction14, position)
									}
									break
								}
							}

							goto l157
						l158:
							position, tokenIndex = position158, tokenIndex158
						}
						add(rulePostfixExpr, position19)
					}
					goto l17
				l18:
					position, tokenIndex = position17, tokenIndex17
					{
						position171 := position
						if buffer[position] != rune('!') {
							goto l15
						}
						position++
						{
							position172, tokenIndex172 := position, tokenIndex
							if buffer[position] != rune('=') {
								goto l172
							}
							position++
							goto l15
						l172:
							position, tokenIndex = position172, tokenIndex172
						}
						if !_rules[rule_]() {
							goto l15
						}
						add(ruleBANG, position171)
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
			position174, tokenIndex174 := position, tokenIndex
			{
				position175 := position
				if !_rules[ruleUnaryExpr]() {
					goto l174
				}
			l176:
				{
					position177, tokenIndex177 := position, tokenIndex
					{
						position178, tokenIndex178 := position, tokenIndex
						{
							position180 := position
							if buffer[position] != rune('>') {
								goto l179
							}
							position++
							if buffer[position] != rune('=') {
								goto l179
							}
							position++
							if !_rules[rule_]() {
								goto l179
							}
							add(ruleGE, position180)
						}
						if !_rules[ruleUnaryExpr]() {
							goto l179
						}
						{
							add(ruleAction17, position)
						}
						goto l178
					l179:
						position, tokenIndex = position178, tokenIndex178
						{
							position183 := position
							if buffer[position] != rune('>') {
								goto l182
							}
							position++
							{
								position184, tokenIndex184 := position, tokenIndex
								if buffer[position] != rune('=') {
									goto l184
								}
								position++
								goto l182
							l184:
								position, tokenIndex = position184, tokenIndex184
							}
							if !_rules[rule_]() {
								goto l182
							}
							add(ruleGT, position183)
						}
						if !_rules[ruleUnaryExpr]() {
							goto l182
						}
						{
							add(ruleAction18, position)
						}
						goto l178
					l182:
						position, tokenIndex = position178, tokenIndex178
						{
							position187 := position
							if buffer[position] != rune('<') {
								goto l186
							}
							position++
							if buffer[position] != rune('=') {
								goto l186
							}
							position++
							if !_rules[rule_]() {
								goto l186
							}
							add(ruleLE, position187)
						}
						if !_rules[ruleUnaryExpr]() {
							goto l186
						}
						{
							add(ruleAction19, position)
						}
						goto l178
					l186:
						position, tokenIndex = position178, tokenIndex178
						{
							position189 := position
							if buffer[position] != rune('<') {
								goto l177
							}
							position++
							{
								position190, tokenIndex190 := position, tokenIndex
								if buffer[position] != rune('=') {
									goto l190
								}
								position++
								goto l177
							l190:
								position, tokenIndex = position190, tokenIndex190
							}
							if !_rules[rule_]() {
								goto l177
							}
							add(ruleLT, position189)
						}
						if !_rules[ruleUnaryExpr]() {
							goto l177
						}
						{
							add(ruleAction20, position)
						}
					}
				l178:
					goto l176
				l177:
					position, tokenIndex = position177, tokenIndex177
				}
				add(ruleRelationalExpr, position175)
			}
			return true
		l174:
			position, tokenIndex = position174, tokenIndex174
			return false
		},
		/* 8 EqualityExpr <- <(RelationalExpr ((EQEQ RelationalExpr Action21) / ((&('c') (CONTAINS RelationalExpr Action24)) | (&('=') (EQ RelationalExpr Action23)) | (&('!') (NE RelationalExpr Action22))))*)> */
		func() bool {
			position192, tokenIndex192 := position, tokenIndex
			{
				position193 := position
				if !_rules[ruleRelationalExpr]() {
					goto l192
				}
			l194:
				{
					position195, tokenIndex195 := position, tokenIndex
					{
						position196, tokenIndex196 := position, tokenIndex
						{
							position198 := position
							if buffer[position] != rune('=') {
								goto l197
							}
							position++
							if buffer[position] != rune('=') {
								goto l197
							}
							position++
							if !_rules[rule_]() {
								goto l197
							}
							add(ruleEQEQ, position198)
						}
						if !_rules[ruleRelationalExpr]() {
							goto l197
						}
						{
							add(ruleAction21, position)
						}
						goto l196
					l197:
						position, tokenIndex = position196, tokenIndex196
						{
							switch buffer[position] {
							case 'c':
								{
									position201 := position
									if buffer[position] != rune('c') {
										goto l195
									}
									position++
									if buffer[position] != rune('o') {
										goto l195
									}
									position++
									if buffer[position] != rune('n') {
										goto l195
									}
									position++
									if buffer[position] != rune('t') {
										goto l195
									}
									position++
									if buffer[position] != rune('a') {
										goto l195
									}
									position++
									if buffer[position] != rune('i') {
										goto l195
									}
									position++
									if buffer[position] != rune('n') {
										goto l195
									}
									position++
									if buffer[position] != rune('s') {
										goto l195
									}
									position++
									{
										position202, tokenIndex202 := position, tokenIndex
										if !_rules[ruleIdChar]() {
											goto l202
										}
										goto l195
									l202:
										position, tokenIndex = position202, tokenIndex202
									}
									if !_rules[rule_]() {
										goto l195
									}
									add(ruleCONTAINS, position201)
								}
								if !_rules[ruleRelationalExpr]() {
									goto l195
								}
								{
									add(ruleAction24, position)
								}
								break
							case '=':
								{
									position204 := position
									if buffer[position] != rune('=') {
										goto l195
									}
									position++
									if !_rules[rule_]() {
										goto l195
									}
									add(ruleEQ, position204)
								}
								if !_rules[ruleRelationalExpr]() {
									goto l195
								}
								{
									add(ruleAction23, position)
								}
								break
							default:
								{
									position206 := position
									if buffer[position] != rune('!') {
										goto l195
									}
									position++
									if buffer[position] != rune('=') {
										goto l195
									}
									position++
									if !_rules[rule_]() {
										goto l195
									}
									add(ruleNE, position206)
								}
								if !_rules[ruleRelationalExpr]() {
									goto l195
								}
								{
									add(ruleAction22, position)
								}
								break
							}
						}

					}
				l196:
					goto l194
				l195:
					position, tokenIndex = position195, tokenIndex195
				}
				add(ruleEqualityExpr, position193)
			}
			return true
		l192:
			position, tokenIndex = position192, tokenIndex192
			return false
		},
		/* 9 LogicalAndExpr <- <(EqualityExpr ((AND EqualityExpr Action25) / (ANDAND EqualityExpr Action26) / (_ EqualityExpr Action27))*)> */
		func() bool {
			position208, tokenIndex208 := position, tokenIndex
			{
				position209 := position
				if !_rules[ruleEqualityExpr]() {
					goto l208
				}
			l210:
				{
					position211, tokenIndex211 := position, tokenIndex
					{
						position212, tokenIndex212 := position, tokenIndex
						{
							position214 := position
							if buffer[position] != rune('a') {
								goto l213
							}
							position++
							if buffer[position] != rune('n') {
								goto l213
							}
							position++
							if buffer[position] != rune('d') {
								goto l213
							}
							position++
							{
								position215, tokenIndex215 := position, tokenIndex
								if !_rules[ruleIdChar]() {
									goto l215
								}
								goto l213
							l215:
								position, tokenIndex = position215, tokenIndex215
							}
							if !_rules[rule_]() {
								goto l213
							}
							add(ruleAND, position214)
						}
						if !_rules[ruleEqualityExpr]() {
							goto l213
						}
						{
							add(ruleAction25, position)
						}
						goto l212
					l213:
						position, tokenIndex = position212, tokenIndex212
						{
							position218 := position
							if buffer[position] != rune('&') {
								goto l217
							}
							position++
							if buffer[position] != rune('&') {
								goto l217
							}
							position++
							if !_rules[rule_]() {
								goto l217
							}
							add(ruleANDAND, position218)
						}
						if !_rules[ruleEqualityExpr]() {
							goto l217
						}
						{
							add(ruleAction26, position)
						}
						goto l212
					l217:
						position, tokenIndex = position212, tokenIndex212
						if !_rules[rule_]() {
							goto l211
						}
						if !_rules[ruleEqualityExpr]() {
							goto l211
						}
						{
							add(ruleAction27, position)
						}
					}
				l212:
					goto l210
				l211:
					position, tokenIndex = position211, tokenIndex211
				}
				add(ruleLogicalAndExpr, position209)
			}
			return true
		l208:
			position, tokenIndex = position208, tokenIndex208
			return false
		},
		/* 10 LogicalOrExpr <- <(LogicalAndExpr ((OR LogicalAndExpr Action28) / (OROR LogicalAndExpr Action29))*)> */
		func() bool {
			position221, tokenIndex221 := position, tokenIndex
			{
				position222 := position
				if !_rules[ruleLogicalAndExpr]() {
					goto l221
				}
			l223:
				{
					position224, tokenIndex224 := position, tokenIndex
					{
						position225, tokenIndex225 := position, tokenIndex
						{
							position227 := position
							if buffer[position] != rune('o') {
								goto l226
							}
							position++
							if buffer[position] != rune('r') {
								goto l226
							}
							position++
							{
								position228, tokenIndex228 := position, tokenIndex
								if !_rules[ruleIdChar]() {
									goto l228
								}
								goto l226
							l228:
								position, tokenIndex = position228, tokenIndex228
							}
							if !_rules[rule_]() {
								goto l226
							}
							add(ruleOR, position227)
						}
						if !_rules[ruleLogicalAndExpr]() {
							goto l226
						}
						{
							add(ruleAction28, position)
						}
						goto l225
					l226:
						position, tokenIndex = position225, tokenIndex225
						{
							position230 := position
							if buffer[position] != rune('|') {
								goto l224
							}
							position++
							if buffer[position] != rune('|') {
								goto l224
							}
							position++
							if !_rules[rule_]() {
								goto l224
							}
							add(ruleOROR, position230)
						}
						if !_rules[ruleLogicalAndExpr]() {
							goto l224
						}
						{
							add(ruleAction29, position)
						}
					}
				l225:
					goto l223
				l224:
					position, tokenIndex = position224, tokenIndex224
				}
				add(ruleLogicalOrExpr, position222)
			}
			return true
		l221:
			position, tokenIndex = position221, tokenIndex221
			return false
		},
		/* 11 LowNotExpr <- <(LogicalOrExpr / (NOT LogicalOrExpr Action30))> */
		nil,
		/* 12 Expr <- <LowNotExpr> */
		func() bool {
			position233, tokenIndex233 := position, tokenIndex
			{
				position234 := position
				{
					position235 := position
					{
						position236, tokenIndex236 := position, tokenIndex
						if !_rules[ruleLogicalOrExpr]() {
							goto l237
						}
						goto l236
					l237:
						position, tokenIndex = position236, tokenIndex236
						if !_rules[ruleNOT]() {
							goto l233
						}
						if !_rules[ruleLogicalOrExpr]() {
							goto l233
						}
						{
							add(ruleAction30, position)
						}
					}
				l236:
					add(ruleLowNotExpr, position235)
				}
				add(ruleExpr, position234)
			}
			return true
		l233:
			position, tokenIndex = position233, tokenIndex233
			return false
		},
		/* 13 String <- <('"' <StringChar*> '"' _)+> */
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
			position249, tokenIndex249 := position, tokenIndex
			{
				position250 := position
				if !_rules[ruleHexDigit]() {
					goto l249
				}
				if !_rules[ruleHexDigit]() {
					goto l249
				}
				if !_rules[ruleHexDigit]() {
					goto l249
				}
				if !_rules[ruleHexDigit]() {
					goto l249
				}
				add(ruleHexQuad, position250)
			}
			return true
		l249:
			position, tokenIndex = position249, tokenIndex249
			return false
		},
		/* 24 HexDigit <- <((&('A' | 'B' | 'C' | 'D' | 'E' | 'F') [A-F]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f') [a-f]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]))> */
		func() bool {
			position251, tokenIndex251 := position, tokenIndex
			{
				position252 := position
				{
					switch buffer[position] {
					case 'A', 'B', 'C', 'D', 'E', 'F':
						if c := buffer[position]; c < rune('A') || c > rune('F') {
							goto l251
						}
						position++
						break
					case 'a', 'b', 'c', 'd', 'e', 'f':
						if c := buffer[position]; c < rune('a') || c > rune('f') {
							goto l251
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l251
						}
						position++
						break
					}
				}

				add(ruleHexDigit, position252)
			}
			return true
		l251:
			position, tokenIndex = position251, tokenIndex251
			return false
		},
		/* 25 Numbers <- <(Number Action31)> */
		func() bool {
			position254, tokenIndex254 := position, tokenIndex
			{
				position255 := position
				if !_rules[ruleNumber]() {
					goto l254
				}
				{
					add(ruleAction31, position)
				}
				add(ruleNumbers, position255)
			}
			return true
		l254:
			position, tokenIndex = position254, tokenIndex254
			return false
		},
		/* 26 Number <- <(<Float> / <Integer>)> */
		func() bool {
			{
				position258 := position
				{
					position259, tokenIndex259 := position, tokenIndex
					{
						position261 := position
						{
							position262 := position
							{
								position263, tokenIndex263 := position, tokenIndex
								{
									position265 := position
									{
										position266, tokenIndex266 := position, tokenIndex
									l268:
										{
											position269, tokenIndex269 := position, tokenIndex
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l269
											}
											position++
											goto l268
										l269:
											position, tokenIndex = position269, tokenIndex269
										}
										if buffer[position] != rune('.') {
											goto l267
										}
										position++
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l267
										}
										position++
									l270:
										{
											position271, tokenIndex271 := position, tokenIndex
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l271
											}
											position++
											goto l270
										l271:
											position, tokenIndex = position271, tokenIndex271
										}
										goto l266
									l267:
										position, tokenIndex = position266, tokenIndex266
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l264
										}
										position++
									l272:
										{
											position273, tokenIndex273 := position, tokenIndex
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l273
											}
											position++
											goto l272
										l273:
											position, tokenIndex = position273, tokenIndex273
										}
										if buffer[position] != rune('.') {
											goto l264
										}
										position++
									}
								l266:
									add(ruleFraction, position265)
								}
								{
									position274, tokenIndex274 := position, tokenIndex
									if !_rules[ruleExponent]() {
										goto l274
									}
									goto l275
								l274:
									position, tokenIndex = position274, tokenIndex274
								}
							l275:
								goto l263
							l264:
								position, tokenIndex = position263, tokenIndex263
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l260
								}
								position++
							l276:
								{
									position277, tokenIndex277 := position, tokenIndex
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l277
									}
									position++
									goto l276
								l277:
									position, tokenIndex = position277, tokenIndex277
								}
								if !_rules[ruleExponent]() {
									goto l260
								}
							}
						l263:
							add(ruleFloat, position262)
						}
						add(rulePegText, position261)
					}
					goto l259
				l260:
					position, tokenIndex = position259, tokenIndex259
					{
						position278 := position
						{
							position279 := position
						l280:
							{
								position281, tokenIndex281 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l281
								}
								position++
								goto l280
							l281:
								position, tokenIndex = position281, tokenIndex281
							}
							add(ruleInteger, position279)
						}
						add(rulePegText, position278)
					}
				}
			l259:
				add(ruleNumber, position258)
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
			position285, tokenIndex285 := position, tokenIndex
			{
				position286 := position
				{
					position287, tokenIndex287 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l288
					}
					position++
					goto l287
				l288:
					position, tokenIndex = position287, tokenIndex287
					if buffer[position] != rune('E') {
						goto l285
					}
					position++
				}
			l287:
				{
					position289, tokenIndex289 := position, tokenIndex
					{
						position291, tokenIndex291 := position, tokenIndex
						if buffer[position] != rune('+') {
							goto l292
						}
						position++
						goto l291
					l292:
						position, tokenIndex = position291, tokenIndex291
						if buffer[position] != rune('-') {
							goto l289
						}
						position++
					}
				l291:
					goto l290
				l289:
					position, tokenIndex = position289, tokenIndex289
				}
			l290:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l285
				}
				position++
			l293:
				{
					position294, tokenIndex294 := position, tokenIndex
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l294
					}
					position++
					goto l293
				l294:
					position, tokenIndex = position294, tokenIndex294
				}
				add(ruleExponent, position286)
			}
			return true
		l285:
			position, tokenIndex = position285, tokenIndex285
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
			position308, tokenIndex308 := position, tokenIndex
			{
				position309 := position
				{
					position310, tokenIndex310 := position, tokenIndex
					if !_rules[ruleKeyword]() {
						goto l310
					}
					goto l308
				l310:
					position, tokenIndex = position310, tokenIndex310
				}
				{
					position311 := position
					{
						position312 := position
						{
							switch buffer[position] {
							case '_':
								if buffer[position] != rune('_') {
									goto l308
								}
								position++
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l308
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l308
								}
								position++
								break
							}
						}

						add(ruleIdCharNoDigit, position312)
					}
				l314:
					{
						position315, tokenIndex315 := position, tokenIndex
						if !_rules[ruleIdChar]() {
							goto l315
						}
						goto l314
					l315:
						position, tokenIndex = position315, tokenIndex315
					}
					add(rulePegText, position311)
				}
				if !_rules[rule_]() {
					goto l308
				}
				add(ruleId, position309)
			}
			return true
		l308:
			position, tokenIndex = position308, tokenIndex308
			return false
		},
		/* 45 IdChar <- <((&('_') '_') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position316, tokenIndex316 := position, tokenIndex
			{
				position317 := position
				{
					switch buffer[position] {
					case '_':
						if buffer[position] != rune('_') {
							goto l316
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l316
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l316
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l316
						}
						position++
						break
					}
				}

				add(ruleIdChar, position317)
			}
			return true
		l316:
			position, tokenIndex = position316, tokenIndex316
			return false
		},
		/* 46 IdCharNoDigit <- <((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 47 Severity <- <((&('f') FATAL) | (&('e') ERROR) | (&('w') WARN) | (&('i') INFO) | (&('d') DEBUG))> */
		nil,
		/* 48 IN <- <('i' 'n' !IdChar _)> */
		func() bool {
			position321, tokenIndex321 := position, tokenIndex
			{
				position322 := position
				if buffer[position] != rune('i') {
					goto l321
				}
				position++
				if buffer[position] != rune('n') {
					goto l321
				}
				position++
				{
					position323, tokenIndex323 := position, tokenIndex
					if !_rules[ruleIdChar]() {
						goto l323
					}
					goto l321
				l323:
					position, tokenIndex = position323, tokenIndex323
				}
				if !_rules[rule_]() {
					goto l321
				}
				add(ruleIN, position322)
			}
			return true
		l321:
			position, tokenIndex = position321, tokenIndex321
			return false
		},
		/* 49 OR <- <('o' 'r' !IdChar _)> */
		nil,
		/* 50 AND <- <('a' 'n' 'd' !IdChar _)> */
		nil,
		/* 51 NOT <- <('n' 'o' 't' !IdChar _)> */
		func() bool {
			position326, tokenIndex326 := position, tokenIndex
			{
				position327 := position
				if buffer[position] != rune('n') {
					goto l326
				}
				position++
				if buffer[position] != rune('o') {
					goto l326
				}
				position++
				if buffer[position] != rune('t') {
					goto l326
				}
				position++
				{
					position328, tokenIndex328 := position, tokenIndex
					if !_rules[ruleIdChar]() {
						goto l328
					}
					goto l326
				l328:
					position, tokenIndex = position328, tokenIndex328
				}
				if !_rules[rule_]() {
					goto l326
				}
				add(ruleNOT, position327)
			}
			return true
		l326:
			position, tokenIndex = position326, tokenIndex326
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
			position335, tokenIndex335 := position, tokenIndex
			{
				position336 := position
				{
					position337, tokenIndex337 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l338
					}
					position++
					if buffer[position] != rune('t') {
						goto l338
					}
					position++
					if buffer[position] != rune('a') {
						goto l338
					}
					position++
					if buffer[position] != rune('g') {
						goto l338
					}
					position++
					if buffer[position] != rune('i') {
						goto l338
					}
					position++
					if buffer[position] != rune('n') {
						goto l338
					}
					position++
					if buffer[position] != rune('g') {
						goto l338
					}
					position++
					goto l337
				l338:
					position, tokenIndex = position337, tokenIndex337
					if buffer[position] != rune('d') {
						goto l339
					}
					position++
					if buffer[position] != rune('e') {
						goto l339
					}
					position++
					if buffer[position] != rune('v') {
						goto l339
					}
					position++
					if buffer[position] != rune('e') {
						goto l339
					}
					position++
					if buffer[position] != rune('l') {
						goto l339
					}
					position++
					if buffer[position] != rune('o') {
						goto l339
					}
					position++
					if buffer[position] != rune('p') {
						goto l339
					}
					position++
					if buffer[position] != rune('m') {
						goto l339
					}
					position++
					if buffer[position] != rune('e') {
						goto l339
					}
					position++
					if buffer[position] != rune('n') {
						goto l339
					}
					position++
					if buffer[position] != rune('t') {
						goto l339
					}
					position++
					goto l337
				l339:
					position, tokenIndex = position337, tokenIndex337
					if buffer[position] != rune('i') {
						goto l340
					}
					position++
					if buffer[position] != rune('n') {
						goto l340
					}
					position++
					if buffer[position] != rune('f') {
						goto l340
					}
					position++
					if buffer[position] != rune('o') {
						goto l340
					}
					position++
					goto l337
				l340:
					position, tokenIndex = position337, tokenIndex337
					if buffer[position] != rune('m') {
						goto l341
					}
					position++
					if buffer[position] != rune('b') {
						goto l341
					}
					position++
					goto l337
				l341:
					position, tokenIndex = position337, tokenIndex337
					{
						switch buffer[position] {
						case 's':
							if buffer[position] != rune('s') {
								goto l335
							}
							position++
							break
						case 'm':
							if buffer[position] != rune('m') {
								goto l335
							}
							position++
							if buffer[position] != rune('s') {
								goto l335
							}
							position++
							break
						case 'b':
							if buffer[position] != rune('b') {
								goto l335
							}
							position++
							break
						case 'k':
							if buffer[position] != rune('k') {
								goto l335
							}
							position++
							if buffer[position] != rune('b') {
								goto l335
							}
							position++
							break
						case 'g':
							if buffer[position] != rune('g') {
								goto l335
							}
							position++
							if buffer[position] != rune('b') {
								goto l335
							}
							position++
							break
						case 'i':
							if buffer[position] != rune('i') {
								goto l335
							}
							position++
							if buffer[position] != rune('n') {
								goto l335
							}
							position++
							break
						case 'f':
							if buffer[position] != rune('f') {
								goto l335
							}
							position++
							if buffer[position] != rune('a') {
								goto l335
							}
							position++
							if buffer[position] != rune('t') {
								goto l335
							}
							position++
							if buffer[position] != rune('a') {
								goto l335
							}
							position++
							if buffer[position] != rune('l') {
								goto l335
							}
							position++
							break
						case 'e':
							if buffer[position] != rune('e') {
								goto l335
							}
							position++
							if buffer[position] != rune('r') {
								goto l335
							}
							position++
							if buffer[position] != rune('r') {
								goto l335
							}
							position++
							if buffer[position] != rune('o') {
								goto l335
							}
							position++
							if buffer[position] != rune('r') {
								goto l335
							}
							position++
							break
						case 'w':
							if buffer[position] != rune('w') {
								goto l335
							}
							position++
							if buffer[position] != rune('a') {
								goto l335
							}
							position++
							if buffer[position] != rune('r') {
								goto l335
							}
							position++
							if buffer[position] != rune('n') {
								goto l335
							}
							position++
							break
						case 'd':
							if buffer[position] != rune('d') {
								goto l335
							}
							position++
							if buffer[position] != rune('e') {
								goto l335
							}
							position++
							if buffer[position] != rune('b') {
								goto l335
							}
							position++
							if buffer[position] != rune('u') {
								goto l335
							}
							position++
							if buffer[position] != rune('g') {
								goto l335
							}
							position++
							break
						case 'c':
							if buffer[position] != rune('c') {
								goto l335
							}
							position++
							if buffer[position] != rune('o') {
								goto l335
							}
							position++
							if buffer[position] != rune('n') {
								goto l335
							}
							position++
							if buffer[position] != rune('t') {
								goto l335
							}
							position++
							if buffer[position] != rune('a') {
								goto l335
							}
							position++
							if buffer[position] != rune('i') {
								goto l335
							}
							position++
							if buffer[position] != rune('n') {
								goto l335
							}
							position++
							if buffer[position] != rune('s') {
								goto l335
							}
							position++
							break
						case 'n':
							if buffer[position] != rune('n') {
								goto l335
							}
							position++
							if buffer[position] != rune('o') {
								goto l335
							}
							position++
							if buffer[position] != rune('t') {
								goto l335
							}
							position++
							break
						case 'a':
							if buffer[position] != rune('a') {
								goto l335
							}
							position++
							if buffer[position] != rune('n') {
								goto l335
							}
							position++
							if buffer[position] != rune('d') {
								goto l335
							}
							position++
							break
						case 'o':
							if buffer[position] != rune('o') {
								goto l335
							}
							position++
							if buffer[position] != rune('r') {
								goto l335
							}
							position++
							break
						default:
							if buffer[position] != rune('p') {
								goto l335
							}
							position++
							if buffer[position] != rune('r') {
								goto l335
							}
							position++
							if buffer[position] != rune('o') {
								goto l335
							}
							position++
							if buffer[position] != rune('d') {
								goto l335
							}
							position++
							if buffer[position] != rune('u') {
								goto l335
							}
							position++
							if buffer[position] != rune('c') {
								goto l335
							}
							position++
							if buffer[position] != rune('t') {
								goto l335
							}
							position++
							if buffer[position] != rune('i') {
								goto l335
							}
							position++
							if buffer[position] != rune('o') {
								goto l335
							}
							position++
							if buffer[position] != rune('n') {
								goto l335
							}
							position++
							break
						}
					}

				}
			l337:
				{
					position343, tokenIndex343 := position, tokenIndex
					if !_rules[ruleIdChar]() {
						goto l343
					}
					goto l335
				l343:
					position, tokenIndex = position343, tokenIndex343
				}
				add(ruleKeyword, position336)
			}
			return true
		l335:
			position, tokenIndex = position335, tokenIndex335
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
			position347, tokenIndex347 := position, tokenIndex
			{
				position348 := position
				if buffer[position] != rune('(') {
					goto l347
				}
				position++
				if !_rules[rule_]() {
					goto l347
				}
				add(ruleLPAR, position348)
			}
			return true
		l347:
			position, tokenIndex = position347, tokenIndex347
			return false
		},
		/* 63 RPAR <- <(')' _)> */
		func() bool {
			position349, tokenIndex349 := position, tokenIndex
			{
				position350 := position
				if buffer[position] != rune(')') {
					goto l349
				}
				position++
				if !_rules[rule_]() {
					goto l349
				}
				add(ruleRPAR, position350)
			}
			return true
		l349:
			position, tokenIndex = position349, tokenIndex349
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
				position363 := position
			l364:
				{
					position365, tokenIndex365 := position, tokenIndex
					{
						position366 := position
						{
							switch buffer[position] {
							case '\t':
								if buffer[position] != rune('\t') {
									goto l365
								}
								position++
								break
							case ' ':
								if buffer[position] != rune(' ') {
									goto l365
								}
								position++
								break
							default:
								{
									position368 := position
									{
										position369, tokenIndex369 := position, tokenIndex
										if buffer[position] != rune('\r') {
											goto l370
										}
										position++
										if buffer[position] != rune('\n') {
											goto l370
										}
										position++
										goto l369
									l370:
										position, tokenIndex = position369, tokenIndex369
										if buffer[position] != rune('\n') {
											goto l371
										}
										position++
										goto l369
									l371:
										position, tokenIndex = position369, tokenIndex369
										if buffer[position] != rune('\r') {
											goto l365
										}
										position++
									}
								l369:
									add(ruleEOL, position368)
								}
								break
							}
						}

						add(ruleWhitespace, position366)
					}
					goto l364
				l365:
					position, tokenIndex = position365, tokenIndex365
				}
				add(rule_, position363)
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
