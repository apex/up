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
	ruleNumber
	ruleInteger
	ruleFloat
	ruleFraction
	ruleExponent
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
	rulePegText
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
	"Number",
	"Integer",
	"Float",
	"Fraction",
	"Exponent",
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
	"PegText",
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
	stack []ast.Node

	Buffer string
	buffer []rune
	rules  [92]func() bool
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
			p.AddLevel(text)
		case ruleAction2:
			p.AddField(text)
		case ruleAction3:
			p.AddString(text)
		case ruleAction4:
			p.AddExpr()
		case ruleAction5:
			p.AddTupleValue()
		case ruleAction6:
			p.AddTupleValue()
		case ruleAction7:
			p.AddTuple()
		case ruleAction8:
			p.AddBinary(ast.IN)
		case ruleAction9:
			p.AddTuple()
		case ruleAction10:
			p.AddBinary(ast.IN)
			p.AddUnary(ast.LNOT)
		case ruleAction11:
			p.AddMember(text)
		case ruleAction12:
			p.AddSubscript(text)
		case ruleAction13:
			p.AddUnary(ast.NOT)
		case ruleAction14:
			p.AddBinary(ast.GE)
		case ruleAction15:
			p.AddBinary(ast.GT)
		case ruleAction16:
			p.AddBinary(ast.LE)
		case ruleAction17:
			p.AddBinary(ast.LT)
		case ruleAction18:
			p.AddBinary(ast.EQ)
		case ruleAction19:
			p.AddBinary(ast.NE)
		case ruleAction20:
			p.AddBinary(ast.EQ)
		case ruleAction21:
			p.AddBinaryContains()
		case ruleAction22:
			p.AddBinary(ast.AND)
		case ruleAction23:
			p.AddBinary(ast.AND)
		case ruleAction24:
			p.AddBinary(ast.AND)
		case ruleAction25:
			p.AddBinary(ast.OR)
		case ruleAction26:
			p.AddBinary(ast.OR)
		case ruleAction27:
			p.AddUnary(ast.LNOT)

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
		/* 1 PrimaryExpr <- <((Severity Action1) / ((&('(') (LPAR Expr RPAR Action4)) | (&('"') (String Action3)) | (&('\t' | '\n' | '\r' | ' ' | '.' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') (Number Action0)) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (Id Action2))))> */
		nil,
		/* 2 TupleExpr <- <(LPAR Expr Action5 (COMMA Expr Action6)* RPAR)> */
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
					add(ruleAction5, position)
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
						add(ruleAction6, position)
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
		/* 3 InExpr <- <(IN Action7 TupleExpr Action8)> */
		nil,
		/* 4 NotInExpr <- <(NOT IN Action9 TupleExpr Action10)> */
		nil,
		/* 5 PostfixExpr <- <(PrimaryExpr ((&('n') NotInExpr) | (&('i') InExpr) | (&('[') (LBRK Number RBRK Action12)) | (&('.') (DOT Id Action11)))*)> */
		nil,
		/* 6 UnaryExpr <- <(PostfixExpr / (BANG RelationalExpr Action13))> */
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
								{
									position23 := position
									{
										switch buffer[position] {
										case 'f':
											{
												position25 := position
												{
													position26 := position
													if buffer[position] != rune('f') {
														goto l22
													}
													position++
													if buffer[position] != rune('a') {
														goto l22
													}
													position++
													if buffer[position] != rune('t') {
														goto l22
													}
													position++
													if buffer[position] != rune('a') {
														goto l22
													}
													position++
													if buffer[position] != rune('l') {
														goto l22
													}
													position++
													add(rulePegText, position26)
												}
												{
													position27, tokenIndex27 := position, tokenIndex
													if !_rules[ruleIdChar]() {
														goto l27
													}
													goto l22
												l27:
													position, tokenIndex = position27, tokenIndex27
												}
												if !_rules[rule_]() {
													goto l22
												}
												add(ruleFATAL, position25)
											}
											break
										case 'e':
											{
												position28 := position
												{
													position29 := position
													if buffer[position] != rune('e') {
														goto l22
													}
													position++
													if buffer[position] != rune('r') {
														goto l22
													}
													position++
													if buffer[position] != rune('r') {
														goto l22
													}
													position++
													if buffer[position] != rune('o') {
														goto l22
													}
													position++
													if buffer[position] != rune('r') {
														goto l22
													}
													position++
													add(rulePegText, position29)
												}
												{
													position30, tokenIndex30 := position, tokenIndex
													if !_rules[ruleIdChar]() {
														goto l30
													}
													goto l22
												l30:
													position, tokenIndex = position30, tokenIndex30
												}
												if !_rules[rule_]() {
													goto l22
												}
												add(ruleERROR, position28)
											}
											break
										case 'w':
											{
												position31 := position
												{
													position32 := position
													if buffer[position] != rune('w') {
														goto l22
													}
													position++
													if buffer[position] != rune('a') {
														goto l22
													}
													position++
													if buffer[position] != rune('r') {
														goto l22
													}
													position++
													if buffer[position] != rune('n') {
														goto l22
													}
													position++
													add(rulePegText, position32)
												}
												{
													position33, tokenIndex33 := position, tokenIndex
													if !_rules[ruleIdChar]() {
														goto l33
													}
													goto l22
												l33:
													position, tokenIndex = position33, tokenIndex33
												}
												if !_rules[rule_]() {
													goto l22
												}
												add(ruleWARN, position31)
											}
											break
										case 'i':
											{
												position34 := position
												{
													position35 := position
													if buffer[position] != rune('i') {
														goto l22
													}
													position++
													if buffer[position] != rune('n') {
														goto l22
													}
													position++
													if buffer[position] != rune('f') {
														goto l22
													}
													position++
													if buffer[position] != rune('o') {
														goto l22
													}
													position++
													add(rulePegText, position35)
												}
												{
													position36, tokenIndex36 := position, tokenIndex
													if !_rules[ruleIdChar]() {
														goto l36
													}
													goto l22
												l36:
													position, tokenIndex = position36, tokenIndex36
												}
												if !_rules[rule_]() {
													goto l22
												}
												add(ruleINFO, position34)
											}
											break
										default:
											{
												position37 := position
												{
													position38 := position
													if buffer[position] != rune('d') {
														goto l22
													}
													position++
													if buffer[position] != rune('e') {
														goto l22
													}
													position++
													if buffer[position] != rune('b') {
														goto l22
													}
													position++
													if buffer[position] != rune('u') {
														goto l22
													}
													position++
													if buffer[position] != rune('g') {
														goto l22
													}
													position++
													add(rulePegText, position38)
												}
												{
													position39, tokenIndex39 := position, tokenIndex
													if !_rules[ruleIdChar]() {
														goto l39
													}
													goto l22
												l39:
													position, tokenIndex = position39, tokenIndex39
												}
												if !_rules[rule_]() {
													goto l22
												}
												add(ruleDEBUG, position37)
											}
											break
										}
									}

									add(ruleSeverity, position23)
								}
								{
									add(ruleAction1, position)
								}
								goto l21
							l22:
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
											add(ruleAction4, position)
										}
										break
									case '"':
										{
											position43 := position
											if buffer[position] != rune('"') {
												goto l18
											}
											position++
											{
												position46 := position
											l47:
												{
													position48, tokenIndex48 := position, tokenIndex
													{
														position49 := position
														{
															position50, tokenIndex50 := position, tokenIndex
															{
																position52 := position
																{
																	position53, tokenIndex53 := position, tokenIndex
																	{
																		position55 := position
																		if buffer[position] != rune('\\') {
																			goto l54
																		}
																		position++
																		{
																			switch buffer[position] {
																			case 'v':
																				if buffer[position] != rune('v') {
																					goto l54
																				}
																				position++
																				break
																			case 't':
																				if buffer[position] != rune('t') {
																					goto l54
																				}
																				position++
																				break
																			case 'r':
																				if buffer[position] != rune('r') {
																					goto l54
																				}
																				position++
																				break
																			case 'n':
																				if buffer[position] != rune('n') {
																					goto l54
																				}
																				position++
																				break
																			case 'f':
																				if buffer[position] != rune('f') {
																					goto l54
																				}
																				position++
																				break
																			case 'b':
																				if buffer[position] != rune('b') {
																					goto l54
																				}
																				position++
																				break
																			case 'a':
																				if buffer[position] != rune('a') {
																					goto l54
																				}
																				position++
																				break
																			case '\\':
																				if buffer[position] != rune('\\') {
																					goto l54
																				}
																				position++
																				break
																			case '?':
																				if buffer[position] != rune('?') {
																					goto l54
																				}
																				position++
																				break
																			case '"':
																				if buffer[position] != rune('"') {
																					goto l54
																				}
																				position++
																				break
																			default:
																				if buffer[position] != rune('\'') {
																					goto l54
																				}
																				position++
																				break
																			}
																		}

																		add(ruleSimpleEscape, position55)
																	}
																	goto l53
																l54:
																	position, tokenIndex = position53, tokenIndex53
																	{
																		position58 := position
																		if buffer[position] != rune('\\') {
																			goto l57
																		}
																		position++
																		if c := buffer[position]; c < rune('0') || c > rune('7') {
																			goto l57
																		}
																		position++
																		{
																			position59, tokenIndex59 := position, tokenIndex
																			if c := buffer[position]; c < rune('0') || c > rune('7') {
																				goto l59
																			}
																			position++
																			goto l60
																		l59:
																			position, tokenIndex = position59, tokenIndex59
																		}
																	l60:
																		{
																			position61, tokenIndex61 := position, tokenIndex
																			if c := buffer[position]; c < rune('0') || c > rune('7') {
																				goto l61
																			}
																			position++
																			goto l62
																		l61:
																			position, tokenIndex = position61, tokenIndex61
																		}
																	l62:
																		add(ruleOctalEscape, position58)
																	}
																	goto l53
																l57:
																	position, tokenIndex = position53, tokenIndex53
																	{
																		position64 := position
																		if buffer[position] != rune('\\') {
																			goto l63
																		}
																		position++
																		if buffer[position] != rune('x') {
																			goto l63
																		}
																		position++
																		if !_rules[ruleHexDigit]() {
																			goto l63
																		}
																	l65:
																		{
																			position66, tokenIndex66 := position, tokenIndex
																			if !_rules[ruleHexDigit]() {
																				goto l66
																			}
																			goto l65
																		l66:
																			position, tokenIndex = position66, tokenIndex66
																		}
																		add(ruleHexEscape, position64)
																	}
																	goto l53
																l63:
																	position, tokenIndex = position53, tokenIndex53
																	{
																		position67 := position
																		{
																			position68, tokenIndex68 := position, tokenIndex
																			if buffer[position] != rune('\\') {
																				goto l69
																			}
																			position++
																			if buffer[position] != rune('u') {
																				goto l69
																			}
																			position++
																			if !_rules[ruleHexQuad]() {
																				goto l69
																			}
																			goto l68
																		l69:
																			position, tokenIndex = position68, tokenIndex68
																			if buffer[position] != rune('\\') {
																				goto l51
																			}
																			position++
																			if buffer[position] != rune('U') {
																				goto l51
																			}
																			position++
																			if !_rules[ruleHexQuad]() {
																				goto l51
																			}
																			if !_rules[ruleHexQuad]() {
																				goto l51
																			}
																		}
																	l68:
																		add(ruleUniversalCharacter, position67)
																	}
																}
															l53:
																add(ruleEscape, position52)
															}
															goto l50
														l51:
															position, tokenIndex = position50, tokenIndex50
															{
																position70, tokenIndex70 := position, tokenIndex
																{
																	switch buffer[position] {
																	case '\\':
																		if buffer[position] != rune('\\') {
																			goto l70
																		}
																		position++
																		break
																	case '\n':
																		if buffer[position] != rune('\n') {
																			goto l70
																		}
																		position++
																		break
																	default:
																		if buffer[position] != rune('"') {
																			goto l70
																		}
																		position++
																		break
																	}
																}

																goto l48
															l70:
																position, tokenIndex = position70, tokenIndex70
															}
															if !matchDot() {
																goto l48
															}
														}
													l50:
														add(ruleStringChar, position49)
													}
													goto l47
												l48:
													position, tokenIndex = position48, tokenIndex48
												}
												add(rulePegText, position46)
											}
											if buffer[position] != rune('"') {
												goto l18
											}
											position++
											if !_rules[rule_]() {
												goto l18
											}
										l44:
											{
												position45, tokenIndex45 := position, tokenIndex
												if buffer[position] != rune('"') {
													goto l45
												}
												position++
												{
													position72 := position
												l73:
													{
														position74, tokenIndex74 := position, tokenIndex
														{
															position75 := position
															{
																position76, tokenIndex76 := position, tokenIndex
																{
																	position78 := position
																	{
																		position79, tokenIndex79 := position, tokenIndex
																		{
																			position81 := position
																			if buffer[position] != rune('\\') {
																				goto l80
																			}
																			position++
																			{
																				switch buffer[position] {
																				case 'v':
																					if buffer[position] != rune('v') {
																						goto l80
																					}
																					position++
																					break
																				case 't':
																					if buffer[position] != rune('t') {
																						goto l80
																					}
																					position++
																					break
																				case 'r':
																					if buffer[position] != rune('r') {
																						goto l80
																					}
																					position++
																					break
																				case 'n':
																					if buffer[position] != rune('n') {
																						goto l80
																					}
																					position++
																					break
																				case 'f':
																					if buffer[position] != rune('f') {
																						goto l80
																					}
																					position++
																					break
																				case 'b':
																					if buffer[position] != rune('b') {
																						goto l80
																					}
																					position++
																					break
																				case 'a':
																					if buffer[position] != rune('a') {
																						goto l80
																					}
																					position++
																					break
																				case '\\':
																					if buffer[position] != rune('\\') {
																						goto l80
																					}
																					position++
																					break
																				case '?':
																					if buffer[position] != rune('?') {
																						goto l80
																					}
																					position++
																					break
																				case '"':
																					if buffer[position] != rune('"') {
																						goto l80
																					}
																					position++
																					break
																				default:
																					if buffer[position] != rune('\'') {
																						goto l80
																					}
																					position++
																					break
																				}
																			}

																			add(ruleSimpleEscape, position81)
																		}
																		goto l79
																	l80:
																		position, tokenIndex = position79, tokenIndex79
																		{
																			position84 := position
																			if buffer[position] != rune('\\') {
																				goto l83
																			}
																			position++
																			if c := buffer[position]; c < rune('0') || c > rune('7') {
																				goto l83
																			}
																			position++
																			{
																				position85, tokenIndex85 := position, tokenIndex
																				if c := buffer[position]; c < rune('0') || c > rune('7') {
																					goto l85
																				}
																				position++
																				goto l86
																			l85:
																				position, tokenIndex = position85, tokenIndex85
																			}
																		l86:
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
																			add(ruleOctalEscape, position84)
																		}
																		goto l79
																	l83:
																		position, tokenIndex = position79, tokenIndex79
																		{
																			position90 := position
																			if buffer[position] != rune('\\') {
																				goto l89
																			}
																			position++
																			if buffer[position] != rune('x') {
																				goto l89
																			}
																			position++
																			if !_rules[ruleHexDigit]() {
																				goto l89
																			}
																		l91:
																			{
																				position92, tokenIndex92 := position, tokenIndex
																				if !_rules[ruleHexDigit]() {
																					goto l92
																				}
																				goto l91
																			l92:
																				position, tokenIndex = position92, tokenIndex92
																			}
																			add(ruleHexEscape, position90)
																		}
																		goto l79
																	l89:
																		position, tokenIndex = position79, tokenIndex79
																		{
																			position93 := position
																			{
																				position94, tokenIndex94 := position, tokenIndex
																				if buffer[position] != rune('\\') {
																					goto l95
																				}
																				position++
																				if buffer[position] != rune('u') {
																					goto l95
																				}
																				position++
																				if !_rules[ruleHexQuad]() {
																					goto l95
																				}
																				goto l94
																			l95:
																				position, tokenIndex = position94, tokenIndex94
																				if buffer[position] != rune('\\') {
																					goto l77
																				}
																				position++
																				if buffer[position] != rune('U') {
																					goto l77
																				}
																				position++
																				if !_rules[ruleHexQuad]() {
																					goto l77
																				}
																				if !_rules[ruleHexQuad]() {
																					goto l77
																				}
																			}
																		l94:
																			add(ruleUniversalCharacter, position93)
																		}
																	}
																l79:
																	add(ruleEscape, position78)
																}
																goto l76
															l77:
																position, tokenIndex = position76, tokenIndex76
																{
																	position96, tokenIndex96 := position, tokenIndex
																	{
																		switch buffer[position] {
																		case '\\':
																			if buffer[position] != rune('\\') {
																				goto l96
																			}
																			position++
																			break
																		case '\n':
																			if buffer[position] != rune('\n') {
																				goto l96
																			}
																			position++
																			break
																		default:
																			if buffer[position] != rune('"') {
																				goto l96
																			}
																			position++
																			break
																		}
																	}

																	goto l74
																l96:
																	position, tokenIndex = position96, tokenIndex96
																}
																if !matchDot() {
																	goto l74
																}
															}
														l76:
															add(ruleStringChar, position75)
														}
														goto l73
													l74:
														position, tokenIndex = position74, tokenIndex74
													}
													add(rulePegText, position72)
												}
												if buffer[position] != rune('"') {
													goto l45
												}
												position++
												if !_rules[rule_]() {
													goto l45
												}
												goto l44
											l45:
												position, tokenIndex = position45, tokenIndex45
											}
											add(ruleString, position43)
										}
										{
											add(ruleAction3, position)
										}
										break
									case '\t', '\n', '\r', ' ', '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if !_rules[ruleNumber]() {
											goto l18
										}
										{
											add(ruleAction0, position)
										}
										break
									default:
										if !_rules[ruleId]() {
											goto l18
										}
										{
											add(ruleAction2, position)
										}
										break
									}
								}

							}
						l21:
							add(rulePrimaryExpr, position20)
						}
					l101:
						{
							position102, tokenIndex102 := position, tokenIndex
							{
								switch buffer[position] {
								case 'n':
									{
										position104 := position
										if !_rules[ruleNOT]() {
											goto l102
										}
										if !_rules[ruleIN]() {
											goto l102
										}
										{
											add(ruleAction9, position)
										}
										if !_rules[ruleTupleExpr]() {
											goto l102
										}
										{
											add(ruleAction10, position)
										}
										add(ruleNotInExpr, position104)
									}
									break
								case 'i':
									{
										position107 := position
										if !_rules[ruleIN]() {
											goto l102
										}
										{
											add(ruleAction7, position)
										}
										if !_rules[ruleTupleExpr]() {
											goto l102
										}
										{
											add(ruleAction8, position)
										}
										add(ruleInExpr, position107)
									}
									break
								case '[':
									{
										position110 := position
										if buffer[position] != rune('[') {
											goto l102
										}
										position++
										if !_rules[rule_]() {
											goto l102
										}
										add(ruleLBRK, position110)
									}
									if !_rules[ruleNumber]() {
										goto l102
									}
									{
										position111 := position
										if buffer[position] != rune(']') {
											goto l102
										}
										position++
										if !_rules[rule_]() {
											goto l102
										}
										add(ruleRBRK, position111)
									}
									{
										add(ruleAction12, position)
									}
									break
								default:
									{
										position113 := position
										if buffer[position] != rune('.') {
											goto l102
										}
										position++
										if !_rules[rule_]() {
											goto l102
										}
										add(ruleDOT, position113)
									}
									if !_rules[ruleId]() {
										goto l102
									}
									{
										add(ruleAction11, position)
									}
									break
								}
							}

							goto l101
						l102:
							position, tokenIndex = position102, tokenIndex102
						}
						add(rulePostfixExpr, position19)
					}
					goto l17
				l18:
					position, tokenIndex = position17, tokenIndex17
					{
						position115 := position
						if buffer[position] != rune('!') {
							goto l15
						}
						position++
						{
							position116, tokenIndex116 := position, tokenIndex
							if buffer[position] != rune('=') {
								goto l116
							}
							position++
							goto l15
						l116:
							position, tokenIndex = position116, tokenIndex116
						}
						if !_rules[rule_]() {
							goto l15
						}
						add(ruleBANG, position115)
					}
					if !_rules[ruleRelationalExpr]() {
						goto l15
					}
					{
						add(ruleAction13, position)
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
		/* 7 RelationalExpr <- <(UnaryExpr ((GE UnaryExpr Action14) / (GT UnaryExpr Action15) / (LE UnaryExpr Action16) / (LT UnaryExpr Action17))*)> */
		func() bool {
			position118, tokenIndex118 := position, tokenIndex
			{
				position119 := position
				if !_rules[ruleUnaryExpr]() {
					goto l118
				}
			l120:
				{
					position121, tokenIndex121 := position, tokenIndex
					{
						position122, tokenIndex122 := position, tokenIndex
						{
							position124 := position
							if buffer[position] != rune('>') {
								goto l123
							}
							position++
							if buffer[position] != rune('=') {
								goto l123
							}
							position++
							if !_rules[rule_]() {
								goto l123
							}
							add(ruleGE, position124)
						}
						if !_rules[ruleUnaryExpr]() {
							goto l123
						}
						{
							add(ruleAction14, position)
						}
						goto l122
					l123:
						position, tokenIndex = position122, tokenIndex122
						{
							position127 := position
							if buffer[position] != rune('>') {
								goto l126
							}
							position++
							{
								position128, tokenIndex128 := position, tokenIndex
								if buffer[position] != rune('=') {
									goto l128
								}
								position++
								goto l126
							l128:
								position, tokenIndex = position128, tokenIndex128
							}
							if !_rules[rule_]() {
								goto l126
							}
							add(ruleGT, position127)
						}
						if !_rules[ruleUnaryExpr]() {
							goto l126
						}
						{
							add(ruleAction15, position)
						}
						goto l122
					l126:
						position, tokenIndex = position122, tokenIndex122
						{
							position131 := position
							if buffer[position] != rune('<') {
								goto l130
							}
							position++
							if buffer[position] != rune('=') {
								goto l130
							}
							position++
							if !_rules[rule_]() {
								goto l130
							}
							add(ruleLE, position131)
						}
						if !_rules[ruleUnaryExpr]() {
							goto l130
						}
						{
							add(ruleAction16, position)
						}
						goto l122
					l130:
						position, tokenIndex = position122, tokenIndex122
						{
							position133 := position
							if buffer[position] != rune('<') {
								goto l121
							}
							position++
							{
								position134, tokenIndex134 := position, tokenIndex
								if buffer[position] != rune('=') {
									goto l134
								}
								position++
								goto l121
							l134:
								position, tokenIndex = position134, tokenIndex134
							}
							if !_rules[rule_]() {
								goto l121
							}
							add(ruleLT, position133)
						}
						if !_rules[ruleUnaryExpr]() {
							goto l121
						}
						{
							add(ruleAction17, position)
						}
					}
				l122:
					goto l120
				l121:
					position, tokenIndex = position121, tokenIndex121
				}
				add(ruleRelationalExpr, position119)
			}
			return true
		l118:
			position, tokenIndex = position118, tokenIndex118
			return false
		},
		/* 8 EqualityExpr <- <(RelationalExpr ((EQEQ RelationalExpr Action18) / ((&('c') (CONTAINS RelationalExpr Action21)) | (&('=') (EQ RelationalExpr Action20)) | (&('!') (NE RelationalExpr Action19))))*)> */
		func() bool {
			position136, tokenIndex136 := position, tokenIndex
			{
				position137 := position
				if !_rules[ruleRelationalExpr]() {
					goto l136
				}
			l138:
				{
					position139, tokenIndex139 := position, tokenIndex
					{
						position140, tokenIndex140 := position, tokenIndex
						{
							position142 := position
							if buffer[position] != rune('=') {
								goto l141
							}
							position++
							if buffer[position] != rune('=') {
								goto l141
							}
							position++
							if !_rules[rule_]() {
								goto l141
							}
							add(ruleEQEQ, position142)
						}
						if !_rules[ruleRelationalExpr]() {
							goto l141
						}
						{
							add(ruleAction18, position)
						}
						goto l140
					l141:
						position, tokenIndex = position140, tokenIndex140
						{
							switch buffer[position] {
							case 'c':
								{
									position145 := position
									if buffer[position] != rune('c') {
										goto l139
									}
									position++
									if buffer[position] != rune('o') {
										goto l139
									}
									position++
									if buffer[position] != rune('n') {
										goto l139
									}
									position++
									if buffer[position] != rune('t') {
										goto l139
									}
									position++
									if buffer[position] != rune('a') {
										goto l139
									}
									position++
									if buffer[position] != rune('i') {
										goto l139
									}
									position++
									if buffer[position] != rune('n') {
										goto l139
									}
									position++
									if buffer[position] != rune('s') {
										goto l139
									}
									position++
									{
										position146, tokenIndex146 := position, tokenIndex
										if !_rules[ruleIdChar]() {
											goto l146
										}
										goto l139
									l146:
										position, tokenIndex = position146, tokenIndex146
									}
									if !_rules[rule_]() {
										goto l139
									}
									add(ruleCONTAINS, position145)
								}
								if !_rules[ruleRelationalExpr]() {
									goto l139
								}
								{
									add(ruleAction21, position)
								}
								break
							case '=':
								{
									position148 := position
									if buffer[position] != rune('=') {
										goto l139
									}
									position++
									if !_rules[rule_]() {
										goto l139
									}
									add(ruleEQ, position148)
								}
								if !_rules[ruleRelationalExpr]() {
									goto l139
								}
								{
									add(ruleAction20, position)
								}
								break
							default:
								{
									position150 := position
									if buffer[position] != rune('!') {
										goto l139
									}
									position++
									if buffer[position] != rune('=') {
										goto l139
									}
									position++
									if !_rules[rule_]() {
										goto l139
									}
									add(ruleNE, position150)
								}
								if !_rules[ruleRelationalExpr]() {
									goto l139
								}
								{
									add(ruleAction19, position)
								}
								break
							}
						}

					}
				l140:
					goto l138
				l139:
					position, tokenIndex = position139, tokenIndex139
				}
				add(ruleEqualityExpr, position137)
			}
			return true
		l136:
			position, tokenIndex = position136, tokenIndex136
			return false
		},
		/* 9 LogicalAndExpr <- <(EqualityExpr ((AND EqualityExpr Action22) / (ANDAND EqualityExpr Action23) / (_ EqualityExpr Action24))*)> */
		func() bool {
			position152, tokenIndex152 := position, tokenIndex
			{
				position153 := position
				if !_rules[ruleEqualityExpr]() {
					goto l152
				}
			l154:
				{
					position155, tokenIndex155 := position, tokenIndex
					{
						position156, tokenIndex156 := position, tokenIndex
						{
							position158 := position
							if buffer[position] != rune('a') {
								goto l157
							}
							position++
							if buffer[position] != rune('n') {
								goto l157
							}
							position++
							if buffer[position] != rune('d') {
								goto l157
							}
							position++
							{
								position159, tokenIndex159 := position, tokenIndex
								if !_rules[ruleIdChar]() {
									goto l159
								}
								goto l157
							l159:
								position, tokenIndex = position159, tokenIndex159
							}
							if !_rules[rule_]() {
								goto l157
							}
							add(ruleAND, position158)
						}
						if !_rules[ruleEqualityExpr]() {
							goto l157
						}
						{
							add(ruleAction22, position)
						}
						goto l156
					l157:
						position, tokenIndex = position156, tokenIndex156
						{
							position162 := position
							if buffer[position] != rune('&') {
								goto l161
							}
							position++
							if buffer[position] != rune('&') {
								goto l161
							}
							position++
							if !_rules[rule_]() {
								goto l161
							}
							add(ruleANDAND, position162)
						}
						if !_rules[ruleEqualityExpr]() {
							goto l161
						}
						{
							add(ruleAction23, position)
						}
						goto l156
					l161:
						position, tokenIndex = position156, tokenIndex156
						if !_rules[rule_]() {
							goto l155
						}
						if !_rules[ruleEqualityExpr]() {
							goto l155
						}
						{
							add(ruleAction24, position)
						}
					}
				l156:
					goto l154
				l155:
					position, tokenIndex = position155, tokenIndex155
				}
				add(ruleLogicalAndExpr, position153)
			}
			return true
		l152:
			position, tokenIndex = position152, tokenIndex152
			return false
		},
		/* 10 LogicalOrExpr <- <(LogicalAndExpr ((OR LogicalAndExpr Action25) / (OROR LogicalAndExpr Action26))*)> */
		func() bool {
			position165, tokenIndex165 := position, tokenIndex
			{
				position166 := position
				if !_rules[ruleLogicalAndExpr]() {
					goto l165
				}
			l167:
				{
					position168, tokenIndex168 := position, tokenIndex
					{
						position169, tokenIndex169 := position, tokenIndex
						{
							position171 := position
							if buffer[position] != rune('o') {
								goto l170
							}
							position++
							if buffer[position] != rune('r') {
								goto l170
							}
							position++
							{
								position172, tokenIndex172 := position, tokenIndex
								if !_rules[ruleIdChar]() {
									goto l172
								}
								goto l170
							l172:
								position, tokenIndex = position172, tokenIndex172
							}
							if !_rules[rule_]() {
								goto l170
							}
							add(ruleOR, position171)
						}
						if !_rules[ruleLogicalAndExpr]() {
							goto l170
						}
						{
							add(ruleAction25, position)
						}
						goto l169
					l170:
						position, tokenIndex = position169, tokenIndex169
						{
							position174 := position
							if buffer[position] != rune('|') {
								goto l168
							}
							position++
							if buffer[position] != rune('|') {
								goto l168
							}
							position++
							if !_rules[rule_]() {
								goto l168
							}
							add(ruleOROR, position174)
						}
						if !_rules[ruleLogicalAndExpr]() {
							goto l168
						}
						{
							add(ruleAction26, position)
						}
					}
				l169:
					goto l167
				l168:
					position, tokenIndex = position168, tokenIndex168
				}
				add(ruleLogicalOrExpr, position166)
			}
			return true
		l165:
			position, tokenIndex = position165, tokenIndex165
			return false
		},
		/* 11 LowNotExpr <- <(LogicalOrExpr / (NOT LogicalOrExpr Action27))> */
		nil,
		/* 12 Expr <- <LowNotExpr> */
		func() bool {
			position177, tokenIndex177 := position, tokenIndex
			{
				position178 := position
				{
					position179 := position
					{
						position180, tokenIndex180 := position, tokenIndex
						if !_rules[ruleLogicalOrExpr]() {
							goto l181
						}
						goto l180
					l181:
						position, tokenIndex = position180, tokenIndex180
						if !_rules[ruleNOT]() {
							goto l177
						}
						if !_rules[ruleLogicalOrExpr]() {
							goto l177
						}
						{
							add(ruleAction27, position)
						}
					}
				l180:
					add(ruleLowNotExpr, position179)
				}
				add(ruleExpr, position178)
			}
			return true
		l177:
			position, tokenIndex = position177, tokenIndex177
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
			position190, tokenIndex190 := position, tokenIndex
			{
				position191 := position
				if !_rules[ruleHexDigit]() {
					goto l190
				}
				if !_rules[ruleHexDigit]() {
					goto l190
				}
				if !_rules[ruleHexDigit]() {
					goto l190
				}
				if !_rules[ruleHexDigit]() {
					goto l190
				}
				add(ruleHexQuad, position191)
			}
			return true
		l190:
			position, tokenIndex = position190, tokenIndex190
			return false
		},
		/* 21 HexDigit <- <((&('A' | 'B' | 'C' | 'D' | 'E' | 'F') [A-F]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f') [a-f]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]))> */
		func() bool {
			position192, tokenIndex192 := position, tokenIndex
			{
				position193 := position
				{
					switch buffer[position] {
					case 'A', 'B', 'C', 'D', 'E', 'F':
						if c := buffer[position]; c < rune('A') || c > rune('F') {
							goto l192
						}
						position++
						break
					case 'a', 'b', 'c', 'd', 'e', 'f':
						if c := buffer[position]; c < rune('a') || c > rune('f') {
							goto l192
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l192
						}
						position++
						break
					}
				}

				add(ruleHexDigit, position193)
			}
			return true
		l192:
			position, tokenIndex = position192, tokenIndex192
			return false
		},
		/* 22 Number <- <((<Float> _) / (<Integer> _))> */
		func() bool {
			position195, tokenIndex195 := position, tokenIndex
			{
				position196 := position
				{
					position197, tokenIndex197 := position, tokenIndex
					{
						position199 := position
						{
							position200 := position
							{
								position201, tokenIndex201 := position, tokenIndex
								{
									position203 := position
									{
										position204, tokenIndex204 := position, tokenIndex
									l206:
										{
											position207, tokenIndex207 := position, tokenIndex
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l207
											}
											position++
											goto l206
										l207:
											position, tokenIndex = position207, tokenIndex207
										}
										if buffer[position] != rune('.') {
											goto l205
										}
										position++
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l205
										}
										position++
									l208:
										{
											position209, tokenIndex209 := position, tokenIndex
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l209
											}
											position++
											goto l208
										l209:
											position, tokenIndex = position209, tokenIndex209
										}
										goto l204
									l205:
										position, tokenIndex = position204, tokenIndex204
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l202
										}
										position++
									l210:
										{
											position211, tokenIndex211 := position, tokenIndex
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l211
											}
											position++
											goto l210
										l211:
											position, tokenIndex = position211, tokenIndex211
										}
										if buffer[position] != rune('.') {
											goto l202
										}
										position++
									}
								l204:
									add(ruleFraction, position203)
								}
								{
									position212, tokenIndex212 := position, tokenIndex
									if !_rules[ruleExponent]() {
										goto l212
									}
									goto l213
								l212:
									position, tokenIndex = position212, tokenIndex212
								}
							l213:
								goto l201
							l202:
								position, tokenIndex = position201, tokenIndex201
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l198
								}
								position++
							l214:
								{
									position215, tokenIndex215 := position, tokenIndex
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l215
									}
									position++
									goto l214
								l215:
									position, tokenIndex = position215, tokenIndex215
								}
								if !_rules[ruleExponent]() {
									goto l198
								}
							}
						l201:
							add(ruleFloat, position200)
						}
						add(rulePegText, position199)
					}
					if !_rules[rule_]() {
						goto l198
					}
					goto l197
				l198:
					position, tokenIndex = position197, tokenIndex197
					{
						position216 := position
						{
							position217 := position
						l218:
							{
								position219, tokenIndex219 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l219
								}
								position++
								goto l218
							l219:
								position, tokenIndex = position219, tokenIndex219
							}
							add(ruleInteger, position217)
						}
						add(rulePegText, position216)
					}
					if !_rules[rule_]() {
						goto l195
					}
				}
			l197:
				add(ruleNumber, position196)
			}
			return true
		l195:
			position, tokenIndex = position195, tokenIndex195
			return false
		},
		/* 23 Integer <- <[0-9]*> */
		nil,
		/* 24 Float <- <((Fraction Exponent?) / ([0-9]+ Exponent))> */
		nil,
		/* 25 Fraction <- <(([0-9]* '.' [0-9]+) / ([0-9]+ '.'))> */
		nil,
		/* 26 Exponent <- <(('e' / 'E') ('+' / '-')? [0-9]+)> */
		func() bool {
			position223, tokenIndex223 := position, tokenIndex
			{
				position224 := position
				{
					position225, tokenIndex225 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l226
					}
					position++
					goto l225
				l226:
					position, tokenIndex = position225, tokenIndex225
					if buffer[position] != rune('E') {
						goto l223
					}
					position++
				}
			l225:
				{
					position227, tokenIndex227 := position, tokenIndex
					{
						position229, tokenIndex229 := position, tokenIndex
						if buffer[position] != rune('+') {
							goto l230
						}
						position++
						goto l229
					l230:
						position, tokenIndex = position229, tokenIndex229
						if buffer[position] != rune('-') {
							goto l227
						}
						position++
					}
				l229:
					goto l228
				l227:
					position, tokenIndex = position227, tokenIndex227
				}
			l228:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l223
				}
				position++
			l231:
				{
					position232, tokenIndex232 := position, tokenIndex
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l232
					}
					position++
					goto l231
				l232:
					position, tokenIndex = position232, tokenIndex232
				}
				add(ruleExponent, position224)
			}
			return true
		l223:
			position, tokenIndex = position223, tokenIndex223
			return false
		},
		/* 27 Id <- <(!Keyword <(IdCharNoDigit IdChar*)> _)> */
		func() bool {
			position233, tokenIndex233 := position, tokenIndex
			{
				position234 := position
				{
					position235, tokenIndex235 := position, tokenIndex
					{
						position236 := position
						{
							position237, tokenIndex237 := position, tokenIndex
							if buffer[position] != rune('i') {
								goto l238
							}
							position++
							if buffer[position] != rune('n') {
								goto l238
							}
							position++
							if buffer[position] != rune('f') {
								goto l238
							}
							position++
							if buffer[position] != rune('o') {
								goto l238
							}
							position++
							goto l237
						l238:
							position, tokenIndex = position237, tokenIndex237
							{
								switch buffer[position] {
								case 'i':
									if buffer[position] != rune('i') {
										goto l235
									}
									position++
									if buffer[position] != rune('n') {
										goto l235
									}
									position++
									break
								case 'f':
									if buffer[position] != rune('f') {
										goto l235
									}
									position++
									if buffer[position] != rune('a') {
										goto l235
									}
									position++
									if buffer[position] != rune('t') {
										goto l235
									}
									position++
									if buffer[position] != rune('a') {
										goto l235
									}
									position++
									if buffer[position] != rune('l') {
										goto l235
									}
									position++
									break
								case 'e':
									if buffer[position] != rune('e') {
										goto l235
									}
									position++
									if buffer[position] != rune('r') {
										goto l235
									}
									position++
									if buffer[position] != rune('r') {
										goto l235
									}
									position++
									if buffer[position] != rune('o') {
										goto l235
									}
									position++
									if buffer[position] != rune('r') {
										goto l235
									}
									position++
									break
								case 'w':
									if buffer[position] != rune('w') {
										goto l235
									}
									position++
									if buffer[position] != rune('a') {
										goto l235
									}
									position++
									if buffer[position] != rune('r') {
										goto l235
									}
									position++
									if buffer[position] != rune('n') {
										goto l235
									}
									position++
									break
								case 'd':
									if buffer[position] != rune('d') {
										goto l235
									}
									position++
									if buffer[position] != rune('e') {
										goto l235
									}
									position++
									if buffer[position] != rune('b') {
										goto l235
									}
									position++
									if buffer[position] != rune('u') {
										goto l235
									}
									position++
									if buffer[position] != rune('g') {
										goto l235
									}
									position++
									break
								case 'c':
									if buffer[position] != rune('c') {
										goto l235
									}
									position++
									if buffer[position] != rune('o') {
										goto l235
									}
									position++
									if buffer[position] != rune('n') {
										goto l235
									}
									position++
									if buffer[position] != rune('t') {
										goto l235
									}
									position++
									if buffer[position] != rune('a') {
										goto l235
									}
									position++
									if buffer[position] != rune('i') {
										goto l235
									}
									position++
									if buffer[position] != rune('n') {
										goto l235
									}
									position++
									if buffer[position] != rune('s') {
										goto l235
									}
									position++
									break
								case 'n':
									if buffer[position] != rune('n') {
										goto l235
									}
									position++
									if buffer[position] != rune('o') {
										goto l235
									}
									position++
									if buffer[position] != rune('t') {
										goto l235
									}
									position++
									break
								case 'a':
									if buffer[position] != rune('a') {
										goto l235
									}
									position++
									if buffer[position] != rune('n') {
										goto l235
									}
									position++
									if buffer[position] != rune('d') {
										goto l235
									}
									position++
									break
								default:
									if buffer[position] != rune('o') {
										goto l235
									}
									position++
									if buffer[position] != rune('r') {
										goto l235
									}
									position++
									break
								}
							}

						}
					l237:
						{
							position240, tokenIndex240 := position, tokenIndex
							if !_rules[ruleIdChar]() {
								goto l240
							}
							goto l235
						l240:
							position, tokenIndex = position240, tokenIndex240
						}
						add(ruleKeyword, position236)
					}
					goto l233
				l235:
					position, tokenIndex = position235, tokenIndex235
				}
				{
					position241 := position
					{
						position242 := position
						{
							switch buffer[position] {
							case '_':
								if buffer[position] != rune('_') {
									goto l233
								}
								position++
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l233
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l233
								}
								position++
								break
							}
						}

						add(ruleIdCharNoDigit, position242)
					}
				l244:
					{
						position245, tokenIndex245 := position, tokenIndex
						if !_rules[ruleIdChar]() {
							goto l245
						}
						goto l244
					l245:
						position, tokenIndex = position245, tokenIndex245
					}
					add(rulePegText, position241)
				}
				if !_rules[rule_]() {
					goto l233
				}
				add(ruleId, position234)
			}
			return true
		l233:
			position, tokenIndex = position233, tokenIndex233
			return false
		},
		/* 28 IdChar <- <((&('_') '_') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position246, tokenIndex246 := position, tokenIndex
			{
				position247 := position
				{
					switch buffer[position] {
					case '_':
						if buffer[position] != rune('_') {
							goto l246
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l246
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l246
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l246
						}
						position++
						break
					}
				}

				add(ruleIdChar, position247)
			}
			return true
		l246:
			position, tokenIndex = position246, tokenIndex246
			return false
		},
		/* 29 IdCharNoDigit <- <((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 30 Severity <- <((&('f') FATAL) | (&('e') ERROR) | (&('w') WARN) | (&('i') INFO) | (&('d') DEBUG))> */
		nil,
		/* 31 IN <- <('i' 'n' !IdChar _)> */
		func() bool {
			position251, tokenIndex251 := position, tokenIndex
			{
				position252 := position
				if buffer[position] != rune('i') {
					goto l251
				}
				position++
				if buffer[position] != rune('n') {
					goto l251
				}
				position++
				{
					position253, tokenIndex253 := position, tokenIndex
					if !_rules[ruleIdChar]() {
						goto l253
					}
					goto l251
				l253:
					position, tokenIndex = position253, tokenIndex253
				}
				if !_rules[rule_]() {
					goto l251
				}
				add(ruleIN, position252)
			}
			return true
		l251:
			position, tokenIndex = position251, tokenIndex251
			return false
		},
		/* 32 OR <- <('o' 'r' !IdChar _)> */
		nil,
		/* 33 AND <- <('a' 'n' 'd' !IdChar _)> */
		nil,
		/* 34 NOT <- <('n' 'o' 't' !IdChar _)> */
		func() bool {
			position256, tokenIndex256 := position, tokenIndex
			{
				position257 := position
				if buffer[position] != rune('n') {
					goto l256
				}
				position++
				if buffer[position] != rune('o') {
					goto l256
				}
				position++
				if buffer[position] != rune('t') {
					goto l256
				}
				position++
				{
					position258, tokenIndex258 := position, tokenIndex
					if !_rules[ruleIdChar]() {
						goto l258
					}
					goto l256
				l258:
					position, tokenIndex = position258, tokenIndex258
				}
				if !_rules[rule_]() {
					goto l256
				}
				add(ruleNOT, position257)
			}
			return true
		l256:
			position, tokenIndex = position256, tokenIndex256
			return false
		},
		/* 35 CONTAINS <- <('c' 'o' 'n' 't' 'a' 'i' 'n' 's' !IdChar _)> */
		nil,
		/* 36 DEBUG <- <(<('d' 'e' 'b' 'u' 'g')> !IdChar _)> */
		nil,
		/* 37 INFO <- <(<('i' 'n' 'f' 'o')> !IdChar _)> */
		nil,
		/* 38 WARN <- <(<('w' 'a' 'r' 'n')> !IdChar _)> */
		nil,
		/* 39 ERROR <- <(<('e' 'r' 'r' 'o' 'r')> !IdChar _)> */
		nil,
		/* 40 FATAL <- <(<('f' 'a' 't' 'a' 'l')> !IdChar _)> */
		nil,
		/* 41 Keyword <- <((('i' 'n' 'f' 'o') / ((&('i') ('i' 'n')) | (&('f') ('f' 'a' 't' 'a' 'l')) | (&('e') ('e' 'r' 'r' 'o' 'r')) | (&('w') ('w' 'a' 'r' 'n')) | (&('d') ('d' 'e' 'b' 'u' 'g')) | (&('c') ('c' 'o' 'n' 't' 'a' 'i' 'n' 's')) | (&('n') ('n' 'o' 't')) | (&('a') ('a' 'n' 'd')) | (&('o') ('o' 'r')))) !IdChar)> */
		nil,
		/* 42 EQ <- <('=' _)> */
		nil,
		/* 43 LBRK <- <('[' _)> */
		nil,
		/* 44 RBRK <- <(']' _)> */
		nil,
		/* 45 LPAR <- <('(' _)> */
		func() bool {
			position269, tokenIndex269 := position, tokenIndex
			{
				position270 := position
				if buffer[position] != rune('(') {
					goto l269
				}
				position++
				if !_rules[rule_]() {
					goto l269
				}
				add(ruleLPAR, position270)
			}
			return true
		l269:
			position, tokenIndex = position269, tokenIndex269
			return false
		},
		/* 46 RPAR <- <(')' _)> */
		func() bool {
			position271, tokenIndex271 := position, tokenIndex
			{
				position272 := position
				if buffer[position] != rune(')') {
					goto l271
				}
				position++
				if !_rules[rule_]() {
					goto l271
				}
				add(ruleRPAR, position272)
			}
			return true
		l271:
			position, tokenIndex = position271, tokenIndex271
			return false
		},
		/* 47 DOT <- <('.' _)> */
		nil,
		/* 48 BANG <- <('!' !'=' _)> */
		nil,
		/* 49 LT <- <('<' !'=' _)> */
		nil,
		/* 50 GT <- <('>' !'=' _)> */
		nil,
		/* 51 LE <- <('<' '=' _)> */
		nil,
		/* 52 EQEQ <- <('=' '=' _)> */
		nil,
		/* 53 GE <- <('>' '=' _)> */
		nil,
		/* 54 NE <- <('!' '=' _)> */
		nil,
		/* 55 ANDAND <- <('&' '&' _)> */
		nil,
		/* 56 OROR <- <('|' '|' _)> */
		nil,
		/* 57 COMMA <- <(',' _)> */
		nil,
		/* 58 _ <- <Whitespace*> */
		func() bool {
			{
				position285 := position
			l286:
				{
					position287, tokenIndex287 := position, tokenIndex
					{
						position288 := position
						{
							switch buffer[position] {
							case '\t':
								if buffer[position] != rune('\t') {
									goto l287
								}
								position++
								break
							case ' ':
								if buffer[position] != rune(' ') {
									goto l287
								}
								position++
								break
							default:
								{
									position290 := position
									{
										position291, tokenIndex291 := position, tokenIndex
										if buffer[position] != rune('\r') {
											goto l292
										}
										position++
										if buffer[position] != rune('\n') {
											goto l292
										}
										position++
										goto l291
									l292:
										position, tokenIndex = position291, tokenIndex291
										if buffer[position] != rune('\n') {
											goto l293
										}
										position++
										goto l291
									l293:
										position, tokenIndex = position291, tokenIndex291
										if buffer[position] != rune('\r') {
											goto l287
										}
										position++
									}
								l291:
									add(ruleEOL, position290)
								}
								break
							}
						}

						add(ruleWhitespace, position288)
					}
					goto l286
				l287:
					position, tokenIndex = position287, tokenIndex287
				}
				add(rule_, position285)
			}
			return true
		},
		/* 59 Whitespace <- <((&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') EOL))> */
		nil,
		/* 60 EOL <- <(('\r' '\n') / '\n' / '\r')> */
		nil,
		/* 61 EOF <- <!.> */
		nil,
		/* 63 Action0 <- <{ p.AddNumber(text) }> */
		nil,
		/* 64 Action1 <- <{ p.AddLevel(text)  }> */
		nil,
		/* 65 Action2 <- <{ p.AddField(text)  }> */
		nil,
		/* 66 Action3 <- <{ p.AddString(text) }> */
		nil,
		/* 67 Action4 <- <{ p.AddExpr()       }> */
		nil,
		/* 68 Action5 <- <{ p.AddTupleValue() }> */
		nil,
		/* 69 Action6 <- <{ p.AddTupleValue() }> */
		nil,
		/* 70 Action7 <- <{ p.AddTuple() }> */
		nil,
		/* 71 Action8 <- <{ p.AddBinary(ast.IN) }> */
		nil,
		/* 72 Action9 <- <{ p.AddTuple() }> */
		nil,
		/* 73 Action10 <- <{ p.AddBinary(ast.IN); p.AddUnary(ast.LNOT) }> */
		nil,
		/* 74 Action11 <- <{ p.AddMember(text)    }> */
		nil,
		/* 75 Action12 <- <{ p.AddSubscript(text) }> */
		nil,
		/* 76 Action13 <- <{ p.AddUnary(ast.NOT) }> */
		nil,
		/* 77 Action14 <- <{ p.AddBinary(ast.GE) }> */
		nil,
		/* 78 Action15 <- <{ p.AddBinary(ast.GT) }> */
		nil,
		/* 79 Action16 <- <{ p.AddBinary(ast.LE) }> */
		nil,
		/* 80 Action17 <- <{ p.AddBinary(ast.LT) }> */
		nil,
		/* 81 Action18 <- <{ p.AddBinary(ast.EQ) }> */
		nil,
		/* 82 Action19 <- <{ p.AddBinary(ast.NE) }> */
		nil,
		/* 83 Action20 <- <{ p.AddBinary(ast.EQ) }> */
		nil,
		/* 84 Action21 <- <{ p.AddBinaryContains()   }> */
		nil,
		/* 85 Action22 <- <{ p.AddBinary(ast.AND) }> */
		nil,
		/* 86 Action23 <- <{ p.AddBinary(ast.AND) }> */
		nil,
		/* 87 Action24 <- <{ p.AddBinary(ast.AND) }> */
		nil,
		/* 88 Action25 <- <{ p.AddBinary(ast.OR) }> */
		nil,
		/* 89 Action26 <- <{ p.AddBinary(ast.OR) }> */
		nil,
		/* 90 Action27 <- <{ p.AddUnary(ast.LNOT) }> */
		nil,
		nil,
	}
	p.rules = _rules
}
