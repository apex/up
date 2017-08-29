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
	rulePegText
	ruleAction29
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
	"PegText",
	"Action29",
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
	rules  [104]func() bool
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
			p.AddField(text)
		case ruleAction4:
			p.AddString(text)
		case ruleAction5:
			p.AddExpr()
		case ruleAction6:
			p.AddTupleValue()
		case ruleAction7:
			p.AddTupleValue()
		case ruleAction8:
			p.AddTuple()
		case ruleAction9:
			p.AddBinary(ast.IN)
		case ruleAction10:
			p.AddTuple()
		case ruleAction11:
			p.AddBinary(ast.IN)
			p.AddUnary(ast.LNOT)
		case ruleAction12:
			p.AddMember(text)
		case ruleAction13:
			p.AddSubscript(text)
		case ruleAction14:
			p.AddUnary(ast.NOT)
		case ruleAction15:
			p.AddBinary(ast.GE)
		case ruleAction16:
			p.AddBinary(ast.GT)
		case ruleAction17:
			p.AddBinary(ast.LE)
		case ruleAction18:
			p.AddBinary(ast.LT)
		case ruleAction19:
			p.AddBinary(ast.EQ)
		case ruleAction20:
			p.AddBinary(ast.NE)
		case ruleAction21:
			p.AddBinary(ast.EQ)
		case ruleAction22:
			p.AddBinaryContains()
		case ruleAction23:
			p.AddBinary(ast.AND)
		case ruleAction24:
			p.AddBinary(ast.AND)
		case ruleAction25:
			p.AddBinary(ast.AND)
		case ruleAction26:
			p.AddBinary(ast.OR)
		case ruleAction27:
			p.AddBinary(ast.OR)
		case ruleAction28:
			p.AddUnary(ast.LNOT)
		case ruleAction29:
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
		/* 1 PrimaryExpr <- <((Numbers Unit _ Action0) / (Severity Action2) / ((&('(') (LPAR Expr RPAR Action5)) | (&('"') (String Action4)) | (&('\t' | '\n' | '\r' | ' ' | '.' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') (Numbers _ Action1)) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (Id Action3))))> */
		nil,
		/* 2 TupleExpr <- <(LPAR Expr Action6 (COMMA Expr Action7)* RPAR)> */
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
					add(ruleAction6, position)
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
						add(ruleAction7, position)
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
		/* 3 InExpr <- <(IN Action8 TupleExpr Action9)> */
		nil,
		/* 4 NotInExpr <- <(NOT IN Action10 TupleExpr Action11)> */
		nil,
		/* 5 PostfixExpr <- <(PrimaryExpr ((&('n') NotInExpr) | (&('i') InExpr) | (&('[') (LBRK Number _ RBRK Action13)) | (&('.') (DOT Id Action12)))*)> */
		nil,
		/* 6 UnaryExpr <- <(PostfixExpr / (BANG RelationalExpr Action14))> */
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
											add(ruleAction5, position)
										}
										break
									case '"':
										{
											position71 := position
											if buffer[position] != rune('"') {
												goto l18
											}
											position++
											{
												position74 := position
											l75:
												{
													position76, tokenIndex76 := position, tokenIndex
													{
														position77 := position
														{
															position78, tokenIndex78 := position, tokenIndex
															{
																position80 := position
																{
																	position81, tokenIndex81 := position, tokenIndex
																	{
																		position83 := position
																		if buffer[position] != rune('\\') {
																			goto l82
																		}
																		position++
																		{
																			switch buffer[position] {
																			case 'v':
																				if buffer[position] != rune('v') {
																					goto l82
																				}
																				position++
																				break
																			case 't':
																				if buffer[position] != rune('t') {
																					goto l82
																				}
																				position++
																				break
																			case 'r':
																				if buffer[position] != rune('r') {
																					goto l82
																				}
																				position++
																				break
																			case 'n':
																				if buffer[position] != rune('n') {
																					goto l82
																				}
																				position++
																				break
																			case 'f':
																				if buffer[position] != rune('f') {
																					goto l82
																				}
																				position++
																				break
																			case 'b':
																				if buffer[position] != rune('b') {
																					goto l82
																				}
																				position++
																				break
																			case 'a':
																				if buffer[position] != rune('a') {
																					goto l82
																				}
																				position++
																				break
																			case '\\':
																				if buffer[position] != rune('\\') {
																					goto l82
																				}
																				position++
																				break
																			case '?':
																				if buffer[position] != rune('?') {
																					goto l82
																				}
																				position++
																				break
																			case '"':
																				if buffer[position] != rune('"') {
																					goto l82
																				}
																				position++
																				break
																			default:
																				if buffer[position] != rune('\'') {
																					goto l82
																				}
																				position++
																				break
																			}
																		}

																		add(ruleSimpleEscape, position83)
																	}
																	goto l81
																l82:
																	position, tokenIndex = position81, tokenIndex81
																	{
																		position86 := position
																		if buffer[position] != rune('\\') {
																			goto l85
																		}
																		position++
																		if c := buffer[position]; c < rune('0') || c > rune('7') {
																			goto l85
																		}
																		position++
																		{
																			position87, tokenIndex87 := position, tokenIndex
																			if c := buffer[position]; c < rune('0') || c > rune('7') {
																				goto l87
																			}
																			position++
																			goto l88
																		l87:
																			position, tokenIndex = position87, tokenIndex87
																		}
																	l88:
																		{
																			position89, tokenIndex89 := position, tokenIndex
																			if c := buffer[position]; c < rune('0') || c > rune('7') {
																				goto l89
																			}
																			position++
																			goto l90
																		l89:
																			position, tokenIndex = position89, tokenIndex89
																		}
																	l90:
																		add(ruleOctalEscape, position86)
																	}
																	goto l81
																l85:
																	position, tokenIndex = position81, tokenIndex81
																	{
																		position92 := position
																		if buffer[position] != rune('\\') {
																			goto l91
																		}
																		position++
																		if buffer[position] != rune('x') {
																			goto l91
																		}
																		position++
																		if !_rules[ruleHexDigit]() {
																			goto l91
																		}
																	l93:
																		{
																			position94, tokenIndex94 := position, tokenIndex
																			if !_rules[ruleHexDigit]() {
																				goto l94
																			}
																			goto l93
																		l94:
																			position, tokenIndex = position94, tokenIndex94
																		}
																		add(ruleHexEscape, position92)
																	}
																	goto l81
																l91:
																	position, tokenIndex = position81, tokenIndex81
																	{
																		position95 := position
																		{
																			position96, tokenIndex96 := position, tokenIndex
																			if buffer[position] != rune('\\') {
																				goto l97
																			}
																			position++
																			if buffer[position] != rune('u') {
																				goto l97
																			}
																			position++
																			if !_rules[ruleHexQuad]() {
																				goto l97
																			}
																			goto l96
																		l97:
																			position, tokenIndex = position96, tokenIndex96
																			if buffer[position] != rune('\\') {
																				goto l79
																			}
																			position++
																			if buffer[position] != rune('U') {
																				goto l79
																			}
																			position++
																			if !_rules[ruleHexQuad]() {
																				goto l79
																			}
																			if !_rules[ruleHexQuad]() {
																				goto l79
																			}
																		}
																	l96:
																		add(ruleUniversalCharacter, position95)
																	}
																}
															l81:
																add(ruleEscape, position80)
															}
															goto l78
														l79:
															position, tokenIndex = position78, tokenIndex78
															{
																position98, tokenIndex98 := position, tokenIndex
																{
																	switch buffer[position] {
																	case '\\':
																		if buffer[position] != rune('\\') {
																			goto l98
																		}
																		position++
																		break
																	case '\n':
																		if buffer[position] != rune('\n') {
																			goto l98
																		}
																		position++
																		break
																	default:
																		if buffer[position] != rune('"') {
																			goto l98
																		}
																		position++
																		break
																	}
																}

																goto l76
															l98:
																position, tokenIndex = position98, tokenIndex98
															}
															if !matchDot() {
																goto l76
															}
														}
													l78:
														add(ruleStringChar, position77)
													}
													goto l75
												l76:
													position, tokenIndex = position76, tokenIndex76
												}
												add(rulePegText, position74)
											}
											if buffer[position] != rune('"') {
												goto l18
											}
											position++
											if !_rules[rule_]() {
												goto l18
											}
										l72:
											{
												position73, tokenIndex73 := position, tokenIndex
												if buffer[position] != rune('"') {
													goto l73
												}
												position++
												{
													position100 := position
												l101:
													{
														position102, tokenIndex102 := position, tokenIndex
														{
															position103 := position
															{
																position104, tokenIndex104 := position, tokenIndex
																{
																	position106 := position
																	{
																		position107, tokenIndex107 := position, tokenIndex
																		{
																			position109 := position
																			if buffer[position] != rune('\\') {
																				goto l108
																			}
																			position++
																			{
																				switch buffer[position] {
																				case 'v':
																					if buffer[position] != rune('v') {
																						goto l108
																					}
																					position++
																					break
																				case 't':
																					if buffer[position] != rune('t') {
																						goto l108
																					}
																					position++
																					break
																				case 'r':
																					if buffer[position] != rune('r') {
																						goto l108
																					}
																					position++
																					break
																				case 'n':
																					if buffer[position] != rune('n') {
																						goto l108
																					}
																					position++
																					break
																				case 'f':
																					if buffer[position] != rune('f') {
																						goto l108
																					}
																					position++
																					break
																				case 'b':
																					if buffer[position] != rune('b') {
																						goto l108
																					}
																					position++
																					break
																				case 'a':
																					if buffer[position] != rune('a') {
																						goto l108
																					}
																					position++
																					break
																				case '\\':
																					if buffer[position] != rune('\\') {
																						goto l108
																					}
																					position++
																					break
																				case '?':
																					if buffer[position] != rune('?') {
																						goto l108
																					}
																					position++
																					break
																				case '"':
																					if buffer[position] != rune('"') {
																						goto l108
																					}
																					position++
																					break
																				default:
																					if buffer[position] != rune('\'') {
																						goto l108
																					}
																					position++
																					break
																				}
																			}

																			add(ruleSimpleEscape, position109)
																		}
																		goto l107
																	l108:
																		position, tokenIndex = position107, tokenIndex107
																		{
																			position112 := position
																			if buffer[position] != rune('\\') {
																				goto l111
																			}
																			position++
																			if c := buffer[position]; c < rune('0') || c > rune('7') {
																				goto l111
																			}
																			position++
																			{
																				position113, tokenIndex113 := position, tokenIndex
																				if c := buffer[position]; c < rune('0') || c > rune('7') {
																					goto l113
																				}
																				position++
																				goto l114
																			l113:
																				position, tokenIndex = position113, tokenIndex113
																			}
																		l114:
																			{
																				position115, tokenIndex115 := position, tokenIndex
																				if c := buffer[position]; c < rune('0') || c > rune('7') {
																					goto l115
																				}
																				position++
																				goto l116
																			l115:
																				position, tokenIndex = position115, tokenIndex115
																			}
																		l116:
																			add(ruleOctalEscape, position112)
																		}
																		goto l107
																	l111:
																		position, tokenIndex = position107, tokenIndex107
																		{
																			position118 := position
																			if buffer[position] != rune('\\') {
																				goto l117
																			}
																			position++
																			if buffer[position] != rune('x') {
																				goto l117
																			}
																			position++
																			if !_rules[ruleHexDigit]() {
																				goto l117
																			}
																		l119:
																			{
																				position120, tokenIndex120 := position, tokenIndex
																				if !_rules[ruleHexDigit]() {
																					goto l120
																				}
																				goto l119
																			l120:
																				position, tokenIndex = position120, tokenIndex120
																			}
																			add(ruleHexEscape, position118)
																		}
																		goto l107
																	l117:
																		position, tokenIndex = position107, tokenIndex107
																		{
																			position121 := position
																			{
																				position122, tokenIndex122 := position, tokenIndex
																				if buffer[position] != rune('\\') {
																					goto l123
																				}
																				position++
																				if buffer[position] != rune('u') {
																					goto l123
																				}
																				position++
																				if !_rules[ruleHexQuad]() {
																					goto l123
																				}
																				goto l122
																			l123:
																				position, tokenIndex = position122, tokenIndex122
																				if buffer[position] != rune('\\') {
																					goto l105
																				}
																				position++
																				if buffer[position] != rune('U') {
																					goto l105
																				}
																				position++
																				if !_rules[ruleHexQuad]() {
																					goto l105
																				}
																				if !_rules[ruleHexQuad]() {
																					goto l105
																				}
																			}
																		l122:
																			add(ruleUniversalCharacter, position121)
																		}
																	}
																l107:
																	add(ruleEscape, position106)
																}
																goto l104
															l105:
																position, tokenIndex = position104, tokenIndex104
																{
																	position124, tokenIndex124 := position, tokenIndex
																	{
																		switch buffer[position] {
																		case '\\':
																			if buffer[position] != rune('\\') {
																				goto l124
																			}
																			position++
																			break
																		case '\n':
																			if buffer[position] != rune('\n') {
																				goto l124
																			}
																			position++
																			break
																		default:
																			if buffer[position] != rune('"') {
																				goto l124
																			}
																			position++
																			break
																		}
																	}

																	goto l102
																l124:
																	position, tokenIndex = position124, tokenIndex124
																}
																if !matchDot() {
																	goto l102
																}
															}
														l104:
															add(ruleStringChar, position103)
														}
														goto l101
													l102:
														position, tokenIndex = position102, tokenIndex102
													}
													add(rulePegText, position100)
												}
												if buffer[position] != rune('"') {
													goto l73
												}
												position++
												if !_rules[rule_]() {
													goto l73
												}
												goto l72
											l73:
												position, tokenIndex = position73, tokenIndex73
											}
											add(ruleString, position71)
										}
										{
											add(ruleAction4, position)
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
										if !_rules[ruleId]() {
											goto l18
										}
										{
											add(ruleAction3, position)
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
											add(ruleAction10, position)
										}
										if !_rules[ruleTupleExpr]() {
											goto l130
										}
										{
											add(ruleAction11, position)
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
											add(ruleAction8, position)
										}
										if !_rules[ruleTupleExpr]() {
											goto l130
										}
										{
											add(ruleAction9, position)
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
										add(ruleAction13, position)
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
										add(ruleAction12, position)
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
						add(ruleAction14, position)
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
		/* 7 RelationalExpr <- <(UnaryExpr ((GE UnaryExpr Action15) / (GT UnaryExpr Action16) / (LE UnaryExpr Action17) / (LT UnaryExpr Action18))*)> */
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
							add(ruleAction15, position)
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
							add(ruleAction16, position)
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
							add(ruleAction17, position)
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
							add(ruleAction18, position)
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
		/* 8 EqualityExpr <- <(RelationalExpr ((EQEQ RelationalExpr Action19) / ((&('c') (CONTAINS RelationalExpr Action22)) | (&('=') (EQ RelationalExpr Action21)) | (&('!') (NE RelationalExpr Action20))))*)> */
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
							add(ruleAction19, position)
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
									add(ruleAction22, position)
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
									add(ruleAction21, position)
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
									add(ruleAction20, position)
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
		/* 9 LogicalAndExpr <- <(EqualityExpr ((AND EqualityExpr Action23) / (ANDAND EqualityExpr Action24) / (_ EqualityExpr Action25))*)> */
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
							add(ruleAction23, position)
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
							add(ruleAction24, position)
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
							add(ruleAction25, position)
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
		/* 10 LogicalOrExpr <- <(LogicalAndExpr ((OR LogicalAndExpr Action26) / (OROR LogicalAndExpr Action27))*)> */
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
							add(ruleAction26, position)
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
							add(ruleAction27, position)
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
		/* 11 LowNotExpr <- <(LogicalOrExpr / (NOT LogicalOrExpr Action28))> */
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
							add(ruleAction28, position)
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
		/* 13 String <- <('"' <StringChar*> '"' _)+> */
		nil,
		/* 14 StringChar <- <(Escape / (!((&('\\') '\\') | (&('\n') '\n') | (&('"') '"')) .))> */
		nil,
		/* 15 Escape <- <(SimpleEscape / OctalEscape / HexEscape / UniversalCharacter)> */
		nil,
		/* 16 SimpleEscape <- <('\\' ((&('v') 'v') | (&('t') 't') | (&('r') 'r') | (&('n') 'n') | (&('f') 'f') | (&('b') 'b') | (&('a') 'a') | (&('\\') '\\') | (&('?') '?') | (&('"') '"') | (&('\'') '\'')))> */
		nil,
		/* 17 OctalEscape <- <('\\' [0-7] [0-7]? [0-7]?)> */
		nil,
		/* 18 HexEscape <- <('\\' 'x' HexDigit+)> */
		nil,
		/* 19 UniversalCharacter <- <(('\\' 'u' HexQuad) / ('\\' 'U' HexQuad HexQuad))> */
		nil,
		/* 20 HexQuad <- <(HexDigit HexDigit HexDigit HexDigit)> */
		func() bool {
			position218, tokenIndex218 := position, tokenIndex
			{
				position219 := position
				if !_rules[ruleHexDigit]() {
					goto l218
				}
				if !_rules[ruleHexDigit]() {
					goto l218
				}
				if !_rules[ruleHexDigit]() {
					goto l218
				}
				if !_rules[ruleHexDigit]() {
					goto l218
				}
				add(ruleHexQuad, position219)
			}
			return true
		l218:
			position, tokenIndex = position218, tokenIndex218
			return false
		},
		/* 21 HexDigit <- <((&('A' | 'B' | 'C' | 'D' | 'E' | 'F') [A-F]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f') [a-f]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]))> */
		func() bool {
			position220, tokenIndex220 := position, tokenIndex
			{
				position221 := position
				{
					switch buffer[position] {
					case 'A', 'B', 'C', 'D', 'E', 'F':
						if c := buffer[position]; c < rune('A') || c > rune('F') {
							goto l220
						}
						position++
						break
					case 'a', 'b', 'c', 'd', 'e', 'f':
						if c := buffer[position]; c < rune('a') || c > rune('f') {
							goto l220
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l220
						}
						position++
						break
					}
				}

				add(ruleHexDigit, position221)
			}
			return true
		l220:
			position, tokenIndex = position220, tokenIndex220
			return false
		},
		/* 22 Numbers <- <(Number Action29)> */
		func() bool {
			position223, tokenIndex223 := position, tokenIndex
			{
				position224 := position
				if !_rules[ruleNumber]() {
					goto l223
				}
				{
					add(ruleAction29, position)
				}
				add(ruleNumbers, position224)
			}
			return true
		l223:
			position, tokenIndex = position223, tokenIndex223
			return false
		},
		/* 23 Number <- <(<Float> / <Integer>)> */
		func() bool {
			{
				position227 := position
				{
					position228, tokenIndex228 := position, tokenIndex
					{
						position230 := position
						{
							position231 := position
							{
								position232, tokenIndex232 := position, tokenIndex
								{
									position234 := position
									{
										position235, tokenIndex235 := position, tokenIndex
									l237:
										{
											position238, tokenIndex238 := position, tokenIndex
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l238
											}
											position++
											goto l237
										l238:
											position, tokenIndex = position238, tokenIndex238
										}
										if buffer[position] != rune('.') {
											goto l236
										}
										position++
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l236
										}
										position++
									l239:
										{
											position240, tokenIndex240 := position, tokenIndex
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l240
											}
											position++
											goto l239
										l240:
											position, tokenIndex = position240, tokenIndex240
										}
										goto l235
									l236:
										position, tokenIndex = position235, tokenIndex235
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l233
										}
										position++
									l241:
										{
											position242, tokenIndex242 := position, tokenIndex
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l242
											}
											position++
											goto l241
										l242:
											position, tokenIndex = position242, tokenIndex242
										}
										if buffer[position] != rune('.') {
											goto l233
										}
										position++
									}
								l235:
									add(ruleFraction, position234)
								}
								{
									position243, tokenIndex243 := position, tokenIndex
									if !_rules[ruleExponent]() {
										goto l243
									}
									goto l244
								l243:
									position, tokenIndex = position243, tokenIndex243
								}
							l244:
								goto l232
							l233:
								position, tokenIndex = position232, tokenIndex232
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l229
								}
								position++
							l245:
								{
									position246, tokenIndex246 := position, tokenIndex
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l246
									}
									position++
									goto l245
								l246:
									position, tokenIndex = position246, tokenIndex246
								}
								if !_rules[ruleExponent]() {
									goto l229
								}
							}
						l232:
							add(ruleFloat, position231)
						}
						add(rulePegText, position230)
					}
					goto l228
				l229:
					position, tokenIndex = position228, tokenIndex228
					{
						position247 := position
						{
							position248 := position
						l249:
							{
								position250, tokenIndex250 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l250
								}
								position++
								goto l249
							l250:
								position, tokenIndex = position250, tokenIndex250
							}
							add(ruleInteger, position248)
						}
						add(rulePegText, position247)
					}
				}
			l228:
				add(ruleNumber, position227)
			}
			return true
		},
		/* 24 Integer <- <[0-9]*> */
		nil,
		/* 25 Float <- <((Fraction Exponent?) / ([0-9]+ Exponent))> */
		nil,
		/* 26 Fraction <- <(([0-9]* '.' [0-9]+) / ([0-9]+ '.'))> */
		nil,
		/* 27 Exponent <- <(('e' / 'E') ('+' / '-')? [0-9]+)> */
		func() bool {
			position254, tokenIndex254 := position, tokenIndex
			{
				position255 := position
				{
					position256, tokenIndex256 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l257
					}
					position++
					goto l256
				l257:
					position, tokenIndex = position256, tokenIndex256
					if buffer[position] != rune('E') {
						goto l254
					}
					position++
				}
			l256:
				{
					position258, tokenIndex258 := position, tokenIndex
					{
						position260, tokenIndex260 := position, tokenIndex
						if buffer[position] != rune('+') {
							goto l261
						}
						position++
						goto l260
					l261:
						position, tokenIndex = position260, tokenIndex260
						if buffer[position] != rune('-') {
							goto l258
						}
						position++
					}
				l260:
					goto l259
				l258:
					position, tokenIndex = position258, tokenIndex258
				}
			l259:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l254
				}
				position++
			l262:
				{
					position263, tokenIndex263 := position, tokenIndex
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l263
					}
					position++
					goto l262
				l263:
					position, tokenIndex = position263, tokenIndex263
				}
				add(ruleExponent, position255)
			}
			return true
		l254:
			position, tokenIndex = position254, tokenIndex254
			return false
		},
		/* 28 Unit <- <(Bytes / Duration)> */
		nil,
		/* 29 Duration <- <(S / MS)> */
		nil,
		/* 30 S <- <(<'s'> !IdChar _)> */
		nil,
		/* 31 MS <- <(<('m' 's')> !IdChar _)> */
		nil,
		/* 32 Bytes <- <((&('g') GB) | (&('m') MB) | (&('k') KB) | (&('b') B))> */
		nil,
		/* 33 B <- <(<'b'> !IdChar _)> */
		nil,
		/* 34 KB <- <(<('k' 'b')> !IdChar _)> */
		nil,
		/* 35 MB <- <(<('m' 'b')> !IdChar _)> */
		nil,
		/* 36 GB <- <(<('g' 'b')> !IdChar _)> */
		nil,
		/* 37 Id <- <(!Keyword <(IdCharNoDigit IdChar*)> _)> */
		func() bool {
			position273, tokenIndex273 := position, tokenIndex
			{
				position274 := position
				{
					position275, tokenIndex275 := position, tokenIndex
					{
						position276 := position
						{
							position277, tokenIndex277 := position, tokenIndex
							if buffer[position] != rune('i') {
								goto l278
							}
							position++
							if buffer[position] != rune('n') {
								goto l278
							}
							position++
							if buffer[position] != rune('f') {
								goto l278
							}
							position++
							if buffer[position] != rune('o') {
								goto l278
							}
							position++
							goto l277
						l278:
							position, tokenIndex = position277, tokenIndex277
							if buffer[position] != rune('m') {
								goto l279
							}
							position++
							if buffer[position] != rune('b') {
								goto l279
							}
							position++
							goto l277
						l279:
							position, tokenIndex = position277, tokenIndex277
							{
								switch buffer[position] {
								case 's':
									if buffer[position] != rune('s') {
										goto l275
									}
									position++
									break
								case 'm':
									if buffer[position] != rune('m') {
										goto l275
									}
									position++
									if buffer[position] != rune('s') {
										goto l275
									}
									position++
									break
								case 'b':
									if buffer[position] != rune('b') {
										goto l275
									}
									position++
									break
								case 'k':
									if buffer[position] != rune('k') {
										goto l275
									}
									position++
									if buffer[position] != rune('b') {
										goto l275
									}
									position++
									break
								case 'g':
									if buffer[position] != rune('g') {
										goto l275
									}
									position++
									if buffer[position] != rune('b') {
										goto l275
									}
									position++
									break
								case 'i':
									if buffer[position] != rune('i') {
										goto l275
									}
									position++
									if buffer[position] != rune('n') {
										goto l275
									}
									position++
									break
								case 'f':
									if buffer[position] != rune('f') {
										goto l275
									}
									position++
									if buffer[position] != rune('a') {
										goto l275
									}
									position++
									if buffer[position] != rune('t') {
										goto l275
									}
									position++
									if buffer[position] != rune('a') {
										goto l275
									}
									position++
									if buffer[position] != rune('l') {
										goto l275
									}
									position++
									break
								case 'e':
									if buffer[position] != rune('e') {
										goto l275
									}
									position++
									if buffer[position] != rune('r') {
										goto l275
									}
									position++
									if buffer[position] != rune('r') {
										goto l275
									}
									position++
									if buffer[position] != rune('o') {
										goto l275
									}
									position++
									if buffer[position] != rune('r') {
										goto l275
									}
									position++
									break
								case 'w':
									if buffer[position] != rune('w') {
										goto l275
									}
									position++
									if buffer[position] != rune('a') {
										goto l275
									}
									position++
									if buffer[position] != rune('r') {
										goto l275
									}
									position++
									if buffer[position] != rune('n') {
										goto l275
									}
									position++
									break
								case 'd':
									if buffer[position] != rune('d') {
										goto l275
									}
									position++
									if buffer[position] != rune('e') {
										goto l275
									}
									position++
									if buffer[position] != rune('b') {
										goto l275
									}
									position++
									if buffer[position] != rune('u') {
										goto l275
									}
									position++
									if buffer[position] != rune('g') {
										goto l275
									}
									position++
									break
								case 'c':
									if buffer[position] != rune('c') {
										goto l275
									}
									position++
									if buffer[position] != rune('o') {
										goto l275
									}
									position++
									if buffer[position] != rune('n') {
										goto l275
									}
									position++
									if buffer[position] != rune('t') {
										goto l275
									}
									position++
									if buffer[position] != rune('a') {
										goto l275
									}
									position++
									if buffer[position] != rune('i') {
										goto l275
									}
									position++
									if buffer[position] != rune('n') {
										goto l275
									}
									position++
									if buffer[position] != rune('s') {
										goto l275
									}
									position++
									break
								case 'n':
									if buffer[position] != rune('n') {
										goto l275
									}
									position++
									if buffer[position] != rune('o') {
										goto l275
									}
									position++
									if buffer[position] != rune('t') {
										goto l275
									}
									position++
									break
								case 'a':
									if buffer[position] != rune('a') {
										goto l275
									}
									position++
									if buffer[position] != rune('n') {
										goto l275
									}
									position++
									if buffer[position] != rune('d') {
										goto l275
									}
									position++
									break
								default:
									if buffer[position] != rune('o') {
										goto l275
									}
									position++
									if buffer[position] != rune('r') {
										goto l275
									}
									position++
									break
								}
							}

						}
					l277:
						{
							position281, tokenIndex281 := position, tokenIndex
							if !_rules[ruleIdChar]() {
								goto l281
							}
							goto l275
						l281:
							position, tokenIndex = position281, tokenIndex281
						}
						add(ruleKeyword, position276)
					}
					goto l273
				l275:
					position, tokenIndex = position275, tokenIndex275
				}
				{
					position282 := position
					{
						position283 := position
						{
							switch buffer[position] {
							case '_':
								if buffer[position] != rune('_') {
									goto l273
								}
								position++
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l273
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l273
								}
								position++
								break
							}
						}

						add(ruleIdCharNoDigit, position283)
					}
				l285:
					{
						position286, tokenIndex286 := position, tokenIndex
						if !_rules[ruleIdChar]() {
							goto l286
						}
						goto l285
					l286:
						position, tokenIndex = position286, tokenIndex286
					}
					add(rulePegText, position282)
				}
				if !_rules[rule_]() {
					goto l273
				}
				add(ruleId, position274)
			}
			return true
		l273:
			position, tokenIndex = position273, tokenIndex273
			return false
		},
		/* 38 IdChar <- <((&('_') '_') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position287, tokenIndex287 := position, tokenIndex
			{
				position288 := position
				{
					switch buffer[position] {
					case '_':
						if buffer[position] != rune('_') {
							goto l287
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l287
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l287
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l287
						}
						position++
						break
					}
				}

				add(ruleIdChar, position288)
			}
			return true
		l287:
			position, tokenIndex = position287, tokenIndex287
			return false
		},
		/* 39 IdCharNoDigit <- <((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 40 Severity <- <((&('f') FATAL) | (&('e') ERROR) | (&('w') WARN) | (&('i') INFO) | (&('d') DEBUG))> */
		nil,
		/* 41 IN <- <('i' 'n' !IdChar _)> */
		func() bool {
			position292, tokenIndex292 := position, tokenIndex
			{
				position293 := position
				if buffer[position] != rune('i') {
					goto l292
				}
				position++
				if buffer[position] != rune('n') {
					goto l292
				}
				position++
				{
					position294, tokenIndex294 := position, tokenIndex
					if !_rules[ruleIdChar]() {
						goto l294
					}
					goto l292
				l294:
					position, tokenIndex = position294, tokenIndex294
				}
				if !_rules[rule_]() {
					goto l292
				}
				add(ruleIN, position293)
			}
			return true
		l292:
			position, tokenIndex = position292, tokenIndex292
			return false
		},
		/* 42 OR <- <('o' 'r' !IdChar _)> */
		nil,
		/* 43 AND <- <('a' 'n' 'd' !IdChar _)> */
		nil,
		/* 44 NOT <- <('n' 'o' 't' !IdChar _)> */
		func() bool {
			position297, tokenIndex297 := position, tokenIndex
			{
				position298 := position
				if buffer[position] != rune('n') {
					goto l297
				}
				position++
				if buffer[position] != rune('o') {
					goto l297
				}
				position++
				if buffer[position] != rune('t') {
					goto l297
				}
				position++
				{
					position299, tokenIndex299 := position, tokenIndex
					if !_rules[ruleIdChar]() {
						goto l299
					}
					goto l297
				l299:
					position, tokenIndex = position299, tokenIndex299
				}
				if !_rules[rule_]() {
					goto l297
				}
				add(ruleNOT, position298)
			}
			return true
		l297:
			position, tokenIndex = position297, tokenIndex297
			return false
		},
		/* 45 CONTAINS <- <('c' 'o' 'n' 't' 'a' 'i' 'n' 's' !IdChar _)> */
		nil,
		/* 46 DEBUG <- <(<('d' 'e' 'b' 'u' 'g')> !IdChar _)> */
		nil,
		/* 47 INFO <- <(<('i' 'n' 'f' 'o')> !IdChar _)> */
		nil,
		/* 48 WARN <- <(<('w' 'a' 'r' 'n')> !IdChar _)> */
		nil,
		/* 49 ERROR <- <(<('e' 'r' 'r' 'o' 'r')> !IdChar _)> */
		nil,
		/* 50 FATAL <- <(<('f' 'a' 't' 'a' 'l')> !IdChar _)> */
		nil,
		/* 51 Keyword <- <((('i' 'n' 'f' 'o') / ('m' 'b') / ((&('s') 's') | (&('m') ('m' 's')) | (&('b') 'b') | (&('k') ('k' 'b')) | (&('g') ('g' 'b')) | (&('i') ('i' 'n')) | (&('f') ('f' 'a' 't' 'a' 'l')) | (&('e') ('e' 'r' 'r' 'o' 'r')) | (&('w') ('w' 'a' 'r' 'n')) | (&('d') ('d' 'e' 'b' 'u' 'g')) | (&('c') ('c' 'o' 'n' 't' 'a' 'i' 'n' 's')) | (&('n') ('n' 'o' 't')) | (&('a') ('a' 'n' 'd')) | (&('o') ('o' 'r')))) !IdChar)> */
		nil,
		/* 52 EQ <- <('=' _)> */
		nil,
		/* 53 LBRK <- <('[' _)> */
		nil,
		/* 54 RBRK <- <(']' _)> */
		nil,
		/* 55 LPAR <- <('(' _)> */
		func() bool {
			position310, tokenIndex310 := position, tokenIndex
			{
				position311 := position
				if buffer[position] != rune('(') {
					goto l310
				}
				position++
				if !_rules[rule_]() {
					goto l310
				}
				add(ruleLPAR, position311)
			}
			return true
		l310:
			position, tokenIndex = position310, tokenIndex310
			return false
		},
		/* 56 RPAR <- <(')' _)> */
		func() bool {
			position312, tokenIndex312 := position, tokenIndex
			{
				position313 := position
				if buffer[position] != rune(')') {
					goto l312
				}
				position++
				if !_rules[rule_]() {
					goto l312
				}
				add(ruleRPAR, position313)
			}
			return true
		l312:
			position, tokenIndex = position312, tokenIndex312
			return false
		},
		/* 57 DOT <- <('.' _)> */
		nil,
		/* 58 BANG <- <('!' !'=' _)> */
		nil,
		/* 59 LT <- <('<' !'=' _)> */
		nil,
		/* 60 GT <- <('>' !'=' _)> */
		nil,
		/* 61 LE <- <('<' '=' _)> */
		nil,
		/* 62 EQEQ <- <('=' '=' _)> */
		nil,
		/* 63 GE <- <('>' '=' _)> */
		nil,
		/* 64 NE <- <('!' '=' _)> */
		nil,
		/* 65 ANDAND <- <('&' '&' _)> */
		nil,
		/* 66 OROR <- <('|' '|' _)> */
		nil,
		/* 67 COMMA <- <(',' _)> */
		nil,
		/* 68 _ <- <Whitespace*> */
		func() bool {
			{
				position326 := position
			l327:
				{
					position328, tokenIndex328 := position, tokenIndex
					{
						position329 := position
						{
							switch buffer[position] {
							case '\t':
								if buffer[position] != rune('\t') {
									goto l328
								}
								position++
								break
							case ' ':
								if buffer[position] != rune(' ') {
									goto l328
								}
								position++
								break
							default:
								{
									position331 := position
									{
										position332, tokenIndex332 := position, tokenIndex
										if buffer[position] != rune('\r') {
											goto l333
										}
										position++
										if buffer[position] != rune('\n') {
											goto l333
										}
										position++
										goto l332
									l333:
										position, tokenIndex = position332, tokenIndex332
										if buffer[position] != rune('\n') {
											goto l334
										}
										position++
										goto l332
									l334:
										position, tokenIndex = position332, tokenIndex332
										if buffer[position] != rune('\r') {
											goto l328
										}
										position++
									}
								l332:
									add(ruleEOL, position331)
								}
								break
							}
						}

						add(ruleWhitespace, position329)
					}
					goto l327
				l328:
					position, tokenIndex = position328, tokenIndex328
				}
				add(rule_, position326)
			}
			return true
		},
		/* 69 Whitespace <- <((&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') EOL))> */
		nil,
		/* 70 EOL <- <(('\r' '\n') / '\n' / '\r')> */
		nil,
		/* 71 EOF <- <!.> */
		nil,
		/* 73 Action0 <- <{ p.AddNumber(text) }> */
		nil,
		/* 74 Action1 <- <{ p.AddNumber("")   }> */
		nil,
		/* 75 Action2 <- <{ p.AddLevel(text)  }> */
		nil,
		/* 76 Action3 <- <{ p.AddField(text)  }> */
		nil,
		/* 77 Action4 <- <{ p.AddString(text) }> */
		nil,
		/* 78 Action5 <- <{ p.AddExpr()       }> */
		nil,
		/* 79 Action6 <- <{ p.AddTupleValue() }> */
		nil,
		/* 80 Action7 <- <{ p.AddTupleValue() }> */
		nil,
		/* 81 Action8 <- <{ p.AddTuple() }> */
		nil,
		/* 82 Action9 <- <{ p.AddBinary(ast.IN) }> */
		nil,
		/* 83 Action10 <- <{ p.AddTuple() }> */
		nil,
		/* 84 Action11 <- <{ p.AddBinary(ast.IN); p.AddUnary(ast.LNOT) }> */
		nil,
		/* 85 Action12 <- <{ p.AddMember(text)    }> */
		nil,
		/* 86 Action13 <- <{ p.AddSubscript(text) }> */
		nil,
		/* 87 Action14 <- <{ p.AddUnary(ast.NOT) }> */
		nil,
		/* 88 Action15 <- <{ p.AddBinary(ast.GE) }> */
		nil,
		/* 89 Action16 <- <{ p.AddBinary(ast.GT) }> */
		nil,
		/* 90 Action17 <- <{ p.AddBinary(ast.LE) }> */
		nil,
		/* 91 Action18 <- <{ p.AddBinary(ast.LT) }> */
		nil,
		/* 92 Action19 <- <{ p.AddBinary(ast.EQ)   }> */
		nil,
		/* 93 Action20 <- <{ p.AddBinary(ast.NE)   }> */
		nil,
		/* 94 Action21 <- <{ p.AddBinary(ast.EQ)   }> */
		nil,
		/* 95 Action22 <- <{ p.AddBinaryContains() }> */
		nil,
		/* 96 Action23 <- <{ p.AddBinary(ast.AND) }> */
		nil,
		/* 97 Action24 <- <{ p.AddBinary(ast.AND) }> */
		nil,
		/* 98 Action25 <- <{ p.AddBinary(ast.AND) }> */
		nil,
		/* 99 Action26 <- <{ p.AddBinary(ast.OR) }> */
		nil,
		/* 100 Action27 <- <{ p.AddBinary(ast.OR) }> */
		nil,
		/* 101 Action28 <- <{ p.AddUnary(ast.LNOT) }> */
		nil,
		nil,
		/* 103 Action29 <- <{ p.SetNumber(text) }> */
		nil,
	}
	p.rules = _rules
}
