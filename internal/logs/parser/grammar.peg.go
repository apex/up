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
	rulePegText
	ruleAction30
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
	"PegText",
	"Action30",
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
	rules  [109]func() bool
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
			p.AddExpr()
		case ruleAction7:
			p.AddTupleValue()
		case ruleAction8:
			p.AddTupleValue()
		case ruleAction9:
			p.AddTuple()
		case ruleAction10:
			p.AddBinary(ast.IN)
		case ruleAction11:
			p.AddTuple()
		case ruleAction12:
			p.AddBinary(ast.IN)
			p.AddUnary(ast.LNOT)
		case ruleAction13:
			p.AddMember(text)
		case ruleAction14:
			p.AddSubscript(text)
		case ruleAction15:
			p.AddUnary(ast.NOT)
		case ruleAction16:
			p.AddBinary(ast.GE)
		case ruleAction17:
			p.AddBinary(ast.GT)
		case ruleAction18:
			p.AddBinary(ast.LE)
		case ruleAction19:
			p.AddBinary(ast.LT)
		case ruleAction20:
			p.AddBinary(ast.EQ)
		case ruleAction21:
			p.AddBinary(ast.NE)
		case ruleAction22:
			p.AddBinary(ast.EQ)
		case ruleAction23:
			p.AddBinaryContains()
		case ruleAction24:
			p.AddBinary(ast.AND)
		case ruleAction25:
			p.AddBinary(ast.AND)
		case ruleAction26:
			p.AddBinary(ast.AND)
		case ruleAction27:
			p.AddBinary(ast.OR)
		case ruleAction28:
			p.AddBinary(ast.OR)
		case ruleAction29:
			p.AddUnary(ast.LNOT)
		case ruleAction30:
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
		/* 1 PrimaryExpr <- <((Numbers Unit _ Action0) / (Severity Action2) / (Stage Action3) / ((&('(') (LPAR Expr RPAR Action6)) | (&('"') (String Action5)) | (&('\t' | '\n' | '\r' | ' ' | '.' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') (Numbers _ Action1)) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (Id Action4))))> */
		nil,
		/* 2 TupleExpr <- <(LPAR Expr Action7 (COMMA Expr Action8)* RPAR)> */
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
					add(ruleAction7, position)
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
						add(ruleAction8, position)
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
		/* 3 InExpr <- <(IN Action9 TupleExpr Action10)> */
		nil,
		/* 4 NotInExpr <- <(NOT IN Action11 TupleExpr Action12)> */
		nil,
		/* 5 PostfixExpr <- <(PrimaryExpr ((&('n') NotInExpr) | (&('i') InExpr) | (&('[') (LBRK Number _ RBRK Action14)) | (&('.') (DOT Id Action13)))*)> */
		nil,
		/* 6 UnaryExpr <- <(PostfixExpr / (BANG RelationalExpr Action15))> */
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
											add(ruleAction6, position)
										}
										break
									case '"':
										{
											position84 := position
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
										l85:
											{
												position86, tokenIndex86 := position, tokenIndex
												if buffer[position] != rune('"') {
													goto l86
												}
												position++
												{
													position113 := position
												l114:
													{
														position115, tokenIndex115 := position, tokenIndex
														{
															position116 := position
															{
																position117, tokenIndex117 := position, tokenIndex
																{
																	position119 := position
																	{
																		position120, tokenIndex120 := position, tokenIndex
																		{
																			position122 := position
																			if buffer[position] != rune('\\') {
																				goto l121
																			}
																			position++
																			{
																				switch buffer[position] {
																				case 'v':
																					if buffer[position] != rune('v') {
																						goto l121
																					}
																					position++
																					break
																				case 't':
																					if buffer[position] != rune('t') {
																						goto l121
																					}
																					position++
																					break
																				case 'r':
																					if buffer[position] != rune('r') {
																						goto l121
																					}
																					position++
																					break
																				case 'n':
																					if buffer[position] != rune('n') {
																						goto l121
																					}
																					position++
																					break
																				case 'f':
																					if buffer[position] != rune('f') {
																						goto l121
																					}
																					position++
																					break
																				case 'b':
																					if buffer[position] != rune('b') {
																						goto l121
																					}
																					position++
																					break
																				case 'a':
																					if buffer[position] != rune('a') {
																						goto l121
																					}
																					position++
																					break
																				case '\\':
																					if buffer[position] != rune('\\') {
																						goto l121
																					}
																					position++
																					break
																				case '?':
																					if buffer[position] != rune('?') {
																						goto l121
																					}
																					position++
																					break
																				case '"':
																					if buffer[position] != rune('"') {
																						goto l121
																					}
																					position++
																					break
																				default:
																					if buffer[position] != rune('\'') {
																						goto l121
																					}
																					position++
																					break
																				}
																			}

																			add(ruleSimpleEscape, position122)
																		}
																		goto l120
																	l121:
																		position, tokenIndex = position120, tokenIndex120
																		{
																			position125 := position
																			if buffer[position] != rune('\\') {
																				goto l124
																			}
																			position++
																			if c := buffer[position]; c < rune('0') || c > rune('7') {
																				goto l124
																			}
																			position++
																			{
																				position126, tokenIndex126 := position, tokenIndex
																				if c := buffer[position]; c < rune('0') || c > rune('7') {
																					goto l126
																				}
																				position++
																				goto l127
																			l126:
																				position, tokenIndex = position126, tokenIndex126
																			}
																		l127:
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
																			add(ruleOctalEscape, position125)
																		}
																		goto l120
																	l124:
																		position, tokenIndex = position120, tokenIndex120
																		{
																			position131 := position
																			if buffer[position] != rune('\\') {
																				goto l130
																			}
																			position++
																			if buffer[position] != rune('x') {
																				goto l130
																			}
																			position++
																			if !_rules[ruleHexDigit]() {
																				goto l130
																			}
																		l132:
																			{
																				position133, tokenIndex133 := position, tokenIndex
																				if !_rules[ruleHexDigit]() {
																					goto l133
																				}
																				goto l132
																			l133:
																				position, tokenIndex = position133, tokenIndex133
																			}
																			add(ruleHexEscape, position131)
																		}
																		goto l120
																	l130:
																		position, tokenIndex = position120, tokenIndex120
																		{
																			position134 := position
																			{
																				position135, tokenIndex135 := position, tokenIndex
																				if buffer[position] != rune('\\') {
																					goto l136
																				}
																				position++
																				if buffer[position] != rune('u') {
																					goto l136
																				}
																				position++
																				if !_rules[ruleHexQuad]() {
																					goto l136
																				}
																				goto l135
																			l136:
																				position, tokenIndex = position135, tokenIndex135
																				if buffer[position] != rune('\\') {
																					goto l118
																				}
																				position++
																				if buffer[position] != rune('U') {
																					goto l118
																				}
																				position++
																				if !_rules[ruleHexQuad]() {
																					goto l118
																				}
																				if !_rules[ruleHexQuad]() {
																					goto l118
																				}
																			}
																		l135:
																			add(ruleUniversalCharacter, position134)
																		}
																	}
																l120:
																	add(ruleEscape, position119)
																}
																goto l117
															l118:
																position, tokenIndex = position117, tokenIndex117
																{
																	position137, tokenIndex137 := position, tokenIndex
																	{
																		switch buffer[position] {
																		case '\\':
																			if buffer[position] != rune('\\') {
																				goto l137
																			}
																			position++
																			break
																		case '\n':
																			if buffer[position] != rune('\n') {
																				goto l137
																			}
																			position++
																			break
																		default:
																			if buffer[position] != rune('"') {
																				goto l137
																			}
																			position++
																			break
																		}
																	}

																	goto l115
																l137:
																	position, tokenIndex = position137, tokenIndex137
																}
																if !matchDot() {
																	goto l115
																}
															}
														l117:
															add(ruleStringChar, position116)
														}
														goto l114
													l115:
														position, tokenIndex = position115, tokenIndex115
													}
													add(rulePegText, position113)
												}
												if buffer[position] != rune('"') {
													goto l86
												}
												position++
												if !_rules[rule_]() {
													goto l86
												}
												goto l85
											l86:
												position, tokenIndex = position86, tokenIndex86
											}
											add(ruleString, position84)
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
										if !_rules[ruleId]() {
											goto l18
										}
										{
											add(ruleAction4, position)
										}
										break
									}
								}

							}
						l21:
							add(rulePrimaryExpr, position20)
						}
					l142:
						{
							position143, tokenIndex143 := position, tokenIndex
							{
								switch buffer[position] {
								case 'n':
									{
										position145 := position
										if !_rules[ruleNOT]() {
											goto l143
										}
										if !_rules[ruleIN]() {
											goto l143
										}
										{
											add(ruleAction11, position)
										}
										if !_rules[ruleTupleExpr]() {
											goto l143
										}
										{
											add(ruleAction12, position)
										}
										add(ruleNotInExpr, position145)
									}
									break
								case 'i':
									{
										position148 := position
										if !_rules[ruleIN]() {
											goto l143
										}
										{
											add(ruleAction9, position)
										}
										if !_rules[ruleTupleExpr]() {
											goto l143
										}
										{
											add(ruleAction10, position)
										}
										add(ruleInExpr, position148)
									}
									break
								case '[':
									{
										position151 := position
										if buffer[position] != rune('[') {
											goto l143
										}
										position++
										if !_rules[rule_]() {
											goto l143
										}
										add(ruleLBRK, position151)
									}
									if !_rules[ruleNumber]() {
										goto l143
									}
									if !_rules[rule_]() {
										goto l143
									}
									{
										position152 := position
										if buffer[position] != rune(']') {
											goto l143
										}
										position++
										if !_rules[rule_]() {
											goto l143
										}
										add(ruleRBRK, position152)
									}
									{
										add(ruleAction14, position)
									}
									break
								default:
									{
										position154 := position
										if buffer[position] != rune('.') {
											goto l143
										}
										position++
										if !_rules[rule_]() {
											goto l143
										}
										add(ruleDOT, position154)
									}
									if !_rules[ruleId]() {
										goto l143
									}
									{
										add(ruleAction13, position)
									}
									break
								}
							}

							goto l142
						l143:
							position, tokenIndex = position143, tokenIndex143
						}
						add(rulePostfixExpr, position19)
					}
					goto l17
				l18:
					position, tokenIndex = position17, tokenIndex17
					{
						position156 := position
						if buffer[position] != rune('!') {
							goto l15
						}
						position++
						{
							position157, tokenIndex157 := position, tokenIndex
							if buffer[position] != rune('=') {
								goto l157
							}
							position++
							goto l15
						l157:
							position, tokenIndex = position157, tokenIndex157
						}
						if !_rules[rule_]() {
							goto l15
						}
						add(ruleBANG, position156)
					}
					if !_rules[ruleRelationalExpr]() {
						goto l15
					}
					{
						add(ruleAction15, position)
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
		/* 7 RelationalExpr <- <(UnaryExpr ((GE UnaryExpr Action16) / (GT UnaryExpr Action17) / (LE UnaryExpr Action18) / (LT UnaryExpr Action19))*)> */
		func() bool {
			position159, tokenIndex159 := position, tokenIndex
			{
				position160 := position
				if !_rules[ruleUnaryExpr]() {
					goto l159
				}
			l161:
				{
					position162, tokenIndex162 := position, tokenIndex
					{
						position163, tokenIndex163 := position, tokenIndex
						{
							position165 := position
							if buffer[position] != rune('>') {
								goto l164
							}
							position++
							if buffer[position] != rune('=') {
								goto l164
							}
							position++
							if !_rules[rule_]() {
								goto l164
							}
							add(ruleGE, position165)
						}
						if !_rules[ruleUnaryExpr]() {
							goto l164
						}
						{
							add(ruleAction16, position)
						}
						goto l163
					l164:
						position, tokenIndex = position163, tokenIndex163
						{
							position168 := position
							if buffer[position] != rune('>') {
								goto l167
							}
							position++
							{
								position169, tokenIndex169 := position, tokenIndex
								if buffer[position] != rune('=') {
									goto l169
								}
								position++
								goto l167
							l169:
								position, tokenIndex = position169, tokenIndex169
							}
							if !_rules[rule_]() {
								goto l167
							}
							add(ruleGT, position168)
						}
						if !_rules[ruleUnaryExpr]() {
							goto l167
						}
						{
							add(ruleAction17, position)
						}
						goto l163
					l167:
						position, tokenIndex = position163, tokenIndex163
						{
							position172 := position
							if buffer[position] != rune('<') {
								goto l171
							}
							position++
							if buffer[position] != rune('=') {
								goto l171
							}
							position++
							if !_rules[rule_]() {
								goto l171
							}
							add(ruleLE, position172)
						}
						if !_rules[ruleUnaryExpr]() {
							goto l171
						}
						{
							add(ruleAction18, position)
						}
						goto l163
					l171:
						position, tokenIndex = position163, tokenIndex163
						{
							position174 := position
							if buffer[position] != rune('<') {
								goto l162
							}
							position++
							{
								position175, tokenIndex175 := position, tokenIndex
								if buffer[position] != rune('=') {
									goto l175
								}
								position++
								goto l162
							l175:
								position, tokenIndex = position175, tokenIndex175
							}
							if !_rules[rule_]() {
								goto l162
							}
							add(ruleLT, position174)
						}
						if !_rules[ruleUnaryExpr]() {
							goto l162
						}
						{
							add(ruleAction19, position)
						}
					}
				l163:
					goto l161
				l162:
					position, tokenIndex = position162, tokenIndex162
				}
				add(ruleRelationalExpr, position160)
			}
			return true
		l159:
			position, tokenIndex = position159, tokenIndex159
			return false
		},
		/* 8 EqualityExpr <- <(RelationalExpr ((EQEQ RelationalExpr Action20) / ((&('c') (CONTAINS RelationalExpr Action23)) | (&('=') (EQ RelationalExpr Action22)) | (&('!') (NE RelationalExpr Action21))))*)> */
		func() bool {
			position177, tokenIndex177 := position, tokenIndex
			{
				position178 := position
				if !_rules[ruleRelationalExpr]() {
					goto l177
				}
			l179:
				{
					position180, tokenIndex180 := position, tokenIndex
					{
						position181, tokenIndex181 := position, tokenIndex
						{
							position183 := position
							if buffer[position] != rune('=') {
								goto l182
							}
							position++
							if buffer[position] != rune('=') {
								goto l182
							}
							position++
							if !_rules[rule_]() {
								goto l182
							}
							add(ruleEQEQ, position183)
						}
						if !_rules[ruleRelationalExpr]() {
							goto l182
						}
						{
							add(ruleAction20, position)
						}
						goto l181
					l182:
						position, tokenIndex = position181, tokenIndex181
						{
							switch buffer[position] {
							case 'c':
								{
									position186 := position
									if buffer[position] != rune('c') {
										goto l180
									}
									position++
									if buffer[position] != rune('o') {
										goto l180
									}
									position++
									if buffer[position] != rune('n') {
										goto l180
									}
									position++
									if buffer[position] != rune('t') {
										goto l180
									}
									position++
									if buffer[position] != rune('a') {
										goto l180
									}
									position++
									if buffer[position] != rune('i') {
										goto l180
									}
									position++
									if buffer[position] != rune('n') {
										goto l180
									}
									position++
									if buffer[position] != rune('s') {
										goto l180
									}
									position++
									{
										position187, tokenIndex187 := position, tokenIndex
										if !_rules[ruleIdChar]() {
											goto l187
										}
										goto l180
									l187:
										position, tokenIndex = position187, tokenIndex187
									}
									if !_rules[rule_]() {
										goto l180
									}
									add(ruleCONTAINS, position186)
								}
								if !_rules[ruleRelationalExpr]() {
									goto l180
								}
								{
									add(ruleAction23, position)
								}
								break
							case '=':
								{
									position189 := position
									if buffer[position] != rune('=') {
										goto l180
									}
									position++
									if !_rules[rule_]() {
										goto l180
									}
									add(ruleEQ, position189)
								}
								if !_rules[ruleRelationalExpr]() {
									goto l180
								}
								{
									add(ruleAction22, position)
								}
								break
							default:
								{
									position191 := position
									if buffer[position] != rune('!') {
										goto l180
									}
									position++
									if buffer[position] != rune('=') {
										goto l180
									}
									position++
									if !_rules[rule_]() {
										goto l180
									}
									add(ruleNE, position191)
								}
								if !_rules[ruleRelationalExpr]() {
									goto l180
								}
								{
									add(ruleAction21, position)
								}
								break
							}
						}

					}
				l181:
					goto l179
				l180:
					position, tokenIndex = position180, tokenIndex180
				}
				add(ruleEqualityExpr, position178)
			}
			return true
		l177:
			position, tokenIndex = position177, tokenIndex177
			return false
		},
		/* 9 LogicalAndExpr <- <(EqualityExpr ((AND EqualityExpr Action24) / (ANDAND EqualityExpr Action25) / (_ EqualityExpr Action26))*)> */
		func() bool {
			position193, tokenIndex193 := position, tokenIndex
			{
				position194 := position
				if !_rules[ruleEqualityExpr]() {
					goto l193
				}
			l195:
				{
					position196, tokenIndex196 := position, tokenIndex
					{
						position197, tokenIndex197 := position, tokenIndex
						{
							position199 := position
							if buffer[position] != rune('a') {
								goto l198
							}
							position++
							if buffer[position] != rune('n') {
								goto l198
							}
							position++
							if buffer[position] != rune('d') {
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
							add(ruleAND, position199)
						}
						if !_rules[ruleEqualityExpr]() {
							goto l198
						}
						{
							add(ruleAction24, position)
						}
						goto l197
					l198:
						position, tokenIndex = position197, tokenIndex197
						{
							position203 := position
							if buffer[position] != rune('&') {
								goto l202
							}
							position++
							if buffer[position] != rune('&') {
								goto l202
							}
							position++
							if !_rules[rule_]() {
								goto l202
							}
							add(ruleANDAND, position203)
						}
						if !_rules[ruleEqualityExpr]() {
							goto l202
						}
						{
							add(ruleAction25, position)
						}
						goto l197
					l202:
						position, tokenIndex = position197, tokenIndex197
						if !_rules[rule_]() {
							goto l196
						}
						if !_rules[ruleEqualityExpr]() {
							goto l196
						}
						{
							add(ruleAction26, position)
						}
					}
				l197:
					goto l195
				l196:
					position, tokenIndex = position196, tokenIndex196
				}
				add(ruleLogicalAndExpr, position194)
			}
			return true
		l193:
			position, tokenIndex = position193, tokenIndex193
			return false
		},
		/* 10 LogicalOrExpr <- <(LogicalAndExpr ((OR LogicalAndExpr Action27) / (OROR LogicalAndExpr Action28))*)> */
		func() bool {
			position206, tokenIndex206 := position, tokenIndex
			{
				position207 := position
				if !_rules[ruleLogicalAndExpr]() {
					goto l206
				}
			l208:
				{
					position209, tokenIndex209 := position, tokenIndex
					{
						position210, tokenIndex210 := position, tokenIndex
						{
							position212 := position
							if buffer[position] != rune('o') {
								goto l211
							}
							position++
							if buffer[position] != rune('r') {
								goto l211
							}
							position++
							{
								position213, tokenIndex213 := position, tokenIndex
								if !_rules[ruleIdChar]() {
									goto l213
								}
								goto l211
							l213:
								position, tokenIndex = position213, tokenIndex213
							}
							if !_rules[rule_]() {
								goto l211
							}
							add(ruleOR, position212)
						}
						if !_rules[ruleLogicalAndExpr]() {
							goto l211
						}
						{
							add(ruleAction27, position)
						}
						goto l210
					l211:
						position, tokenIndex = position210, tokenIndex210
						{
							position215 := position
							if buffer[position] != rune('|') {
								goto l209
							}
							position++
							if buffer[position] != rune('|') {
								goto l209
							}
							position++
							if !_rules[rule_]() {
								goto l209
							}
							add(ruleOROR, position215)
						}
						if !_rules[ruleLogicalAndExpr]() {
							goto l209
						}
						{
							add(ruleAction28, position)
						}
					}
				l210:
					goto l208
				l209:
					position, tokenIndex = position209, tokenIndex209
				}
				add(ruleLogicalOrExpr, position207)
			}
			return true
		l206:
			position, tokenIndex = position206, tokenIndex206
			return false
		},
		/* 11 LowNotExpr <- <(LogicalOrExpr / (NOT LogicalOrExpr Action29))> */
		nil,
		/* 12 Expr <- <LowNotExpr> */
		func() bool {
			position218, tokenIndex218 := position, tokenIndex
			{
				position219 := position
				{
					position220 := position
					{
						position221, tokenIndex221 := position, tokenIndex
						if !_rules[ruleLogicalOrExpr]() {
							goto l222
						}
						goto l221
					l222:
						position, tokenIndex = position221, tokenIndex221
						if !_rules[ruleNOT]() {
							goto l218
						}
						if !_rules[ruleLogicalOrExpr]() {
							goto l218
						}
						{
							add(ruleAction29, position)
						}
					}
				l221:
					add(ruleLowNotExpr, position220)
				}
				add(ruleExpr, position219)
			}
			return true
		l218:
			position, tokenIndex = position218, tokenIndex218
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
			position231, tokenIndex231 := position, tokenIndex
			{
				position232 := position
				if !_rules[ruleHexDigit]() {
					goto l231
				}
				if !_rules[ruleHexDigit]() {
					goto l231
				}
				if !_rules[ruleHexDigit]() {
					goto l231
				}
				if !_rules[ruleHexDigit]() {
					goto l231
				}
				add(ruleHexQuad, position232)
			}
			return true
		l231:
			position, tokenIndex = position231, tokenIndex231
			return false
		},
		/* 21 HexDigit <- <((&('A' | 'B' | 'C' | 'D' | 'E' | 'F') [A-F]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f') [a-f]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]))> */
		func() bool {
			position233, tokenIndex233 := position, tokenIndex
			{
				position234 := position
				{
					switch buffer[position] {
					case 'A', 'B', 'C', 'D', 'E', 'F':
						if c := buffer[position]; c < rune('A') || c > rune('F') {
							goto l233
						}
						position++
						break
					case 'a', 'b', 'c', 'd', 'e', 'f':
						if c := buffer[position]; c < rune('a') || c > rune('f') {
							goto l233
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l233
						}
						position++
						break
					}
				}

				add(ruleHexDigit, position234)
			}
			return true
		l233:
			position, tokenIndex = position233, tokenIndex233
			return false
		},
		/* 22 Numbers <- <(Number Action30)> */
		func() bool {
			position236, tokenIndex236 := position, tokenIndex
			{
				position237 := position
				if !_rules[ruleNumber]() {
					goto l236
				}
				{
					add(ruleAction30, position)
				}
				add(ruleNumbers, position237)
			}
			return true
		l236:
			position, tokenIndex = position236, tokenIndex236
			return false
		},
		/* 23 Number <- <(<Float> / <Integer>)> */
		func() bool {
			{
				position240 := position
				{
					position241, tokenIndex241 := position, tokenIndex
					{
						position243 := position
						{
							position244 := position
							{
								position245, tokenIndex245 := position, tokenIndex
								{
									position247 := position
									{
										position248, tokenIndex248 := position, tokenIndex
									l250:
										{
											position251, tokenIndex251 := position, tokenIndex
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l251
											}
											position++
											goto l250
										l251:
											position, tokenIndex = position251, tokenIndex251
										}
										if buffer[position] != rune('.') {
											goto l249
										}
										position++
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l249
										}
										position++
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
										goto l248
									l249:
										position, tokenIndex = position248, tokenIndex248
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l246
										}
										position++
									l254:
										{
											position255, tokenIndex255 := position, tokenIndex
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l255
											}
											position++
											goto l254
										l255:
											position, tokenIndex = position255, tokenIndex255
										}
										if buffer[position] != rune('.') {
											goto l246
										}
										position++
									}
								l248:
									add(ruleFraction, position247)
								}
								{
									position256, tokenIndex256 := position, tokenIndex
									if !_rules[ruleExponent]() {
										goto l256
									}
									goto l257
								l256:
									position, tokenIndex = position256, tokenIndex256
								}
							l257:
								goto l245
							l246:
								position, tokenIndex = position245, tokenIndex245
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l242
								}
								position++
							l258:
								{
									position259, tokenIndex259 := position, tokenIndex
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l259
									}
									position++
									goto l258
								l259:
									position, tokenIndex = position259, tokenIndex259
								}
								if !_rules[ruleExponent]() {
									goto l242
								}
							}
						l245:
							add(ruleFloat, position244)
						}
						add(rulePegText, position243)
					}
					goto l241
				l242:
					position, tokenIndex = position241, tokenIndex241
					{
						position260 := position
						{
							position261 := position
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
							add(ruleInteger, position261)
						}
						add(rulePegText, position260)
					}
				}
			l241:
				add(ruleNumber, position240)
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
			position267, tokenIndex267 := position, tokenIndex
			{
				position268 := position
				{
					position269, tokenIndex269 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l270
					}
					position++
					goto l269
				l270:
					position, tokenIndex = position269, tokenIndex269
					if buffer[position] != rune('E') {
						goto l267
					}
					position++
				}
			l269:
				{
					position271, tokenIndex271 := position, tokenIndex
					{
						position273, tokenIndex273 := position, tokenIndex
						if buffer[position] != rune('+') {
							goto l274
						}
						position++
						goto l273
					l274:
						position, tokenIndex = position273, tokenIndex273
						if buffer[position] != rune('-') {
							goto l271
						}
						position++
					}
				l273:
					goto l272
				l271:
					position, tokenIndex = position271, tokenIndex271
				}
			l272:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l267
				}
				position++
			l275:
				{
					position276, tokenIndex276 := position, tokenIndex
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l276
					}
					position++
					goto l275
				l276:
					position, tokenIndex = position276, tokenIndex276
				}
				add(ruleExponent, position268)
			}
			return true
		l267:
			position, tokenIndex = position267, tokenIndex267
			return false
		},
		/* 28 Stage <- <((&('p') PRODUCTION) | (&('s') STAGING) | (&('d') DEVELOPMENT))> */
		nil,
		/* 29 DEVELOPMENT <- <(<('d' 'e' 'v' 'e' 'l' 'o' 'p' 'm' 'e' 'n' 't')> !IdChar _)> */
		nil,
		/* 30 STAGING <- <(<('s' 't' 'a' 'g' 'i' 'n' 'g')> !IdChar _)> */
		nil,
		/* 31 PRODUCTION <- <(<('p' 'r' 'o' 'd' 'u' 'c' 't' 'i' 'o' 'n')> !IdChar _)> */
		nil,
		/* 32 Unit <- <(Bytes / Duration)> */
		nil,
		/* 33 Duration <- <(S / MS)> */
		nil,
		/* 34 S <- <(<'s'> !IdChar _)> */
		nil,
		/* 35 MS <- <(<('m' 's')> !IdChar _)> */
		nil,
		/* 36 Bytes <- <((&('g') GB) | (&('m') MB) | (&('k') KB) | (&('b') B))> */
		nil,
		/* 37 B <- <(<'b'> !IdChar _)> */
		nil,
		/* 38 KB <- <(<('k' 'b')> !IdChar _)> */
		nil,
		/* 39 MB <- <(<('m' 'b')> !IdChar _)> */
		nil,
		/* 40 GB <- <(<('g' 'b')> !IdChar _)> */
		nil,
		/* 41 Id <- <(!Keyword <(IdCharNoDigit IdChar*)> _)> */
		func() bool {
			position290, tokenIndex290 := position, tokenIndex
			{
				position291 := position
				{
					position292, tokenIndex292 := position, tokenIndex
					{
						position293 := position
						{
							position294, tokenIndex294 := position, tokenIndex
							if buffer[position] != rune('s') {
								goto l295
							}
							position++
							if buffer[position] != rune('t') {
								goto l295
							}
							position++
							if buffer[position] != rune('a') {
								goto l295
							}
							position++
							if buffer[position] != rune('g') {
								goto l295
							}
							position++
							if buffer[position] != rune('i') {
								goto l295
							}
							position++
							if buffer[position] != rune('n') {
								goto l295
							}
							position++
							if buffer[position] != rune('g') {
								goto l295
							}
							position++
							goto l294
						l295:
							position, tokenIndex = position294, tokenIndex294
							if buffer[position] != rune('d') {
								goto l296
							}
							position++
							if buffer[position] != rune('e') {
								goto l296
							}
							position++
							if buffer[position] != rune('v') {
								goto l296
							}
							position++
							if buffer[position] != rune('e') {
								goto l296
							}
							position++
							if buffer[position] != rune('l') {
								goto l296
							}
							position++
							if buffer[position] != rune('o') {
								goto l296
							}
							position++
							if buffer[position] != rune('p') {
								goto l296
							}
							position++
							if buffer[position] != rune('m') {
								goto l296
							}
							position++
							if buffer[position] != rune('e') {
								goto l296
							}
							position++
							if buffer[position] != rune('n') {
								goto l296
							}
							position++
							if buffer[position] != rune('t') {
								goto l296
							}
							position++
							goto l294
						l296:
							position, tokenIndex = position294, tokenIndex294
							if buffer[position] != rune('i') {
								goto l297
							}
							position++
							if buffer[position] != rune('n') {
								goto l297
							}
							position++
							if buffer[position] != rune('f') {
								goto l297
							}
							position++
							if buffer[position] != rune('o') {
								goto l297
							}
							position++
							goto l294
						l297:
							position, tokenIndex = position294, tokenIndex294
							if buffer[position] != rune('m') {
								goto l298
							}
							position++
							if buffer[position] != rune('b') {
								goto l298
							}
							position++
							goto l294
						l298:
							position, tokenIndex = position294, tokenIndex294
							{
								switch buffer[position] {
								case 's':
									if buffer[position] != rune('s') {
										goto l292
									}
									position++
									break
								case 'm':
									if buffer[position] != rune('m') {
										goto l292
									}
									position++
									if buffer[position] != rune('s') {
										goto l292
									}
									position++
									break
								case 'b':
									if buffer[position] != rune('b') {
										goto l292
									}
									position++
									break
								case 'k':
									if buffer[position] != rune('k') {
										goto l292
									}
									position++
									if buffer[position] != rune('b') {
										goto l292
									}
									position++
									break
								case 'g':
									if buffer[position] != rune('g') {
										goto l292
									}
									position++
									if buffer[position] != rune('b') {
										goto l292
									}
									position++
									break
								case 'i':
									if buffer[position] != rune('i') {
										goto l292
									}
									position++
									if buffer[position] != rune('n') {
										goto l292
									}
									position++
									break
								case 'f':
									if buffer[position] != rune('f') {
										goto l292
									}
									position++
									if buffer[position] != rune('a') {
										goto l292
									}
									position++
									if buffer[position] != rune('t') {
										goto l292
									}
									position++
									if buffer[position] != rune('a') {
										goto l292
									}
									position++
									if buffer[position] != rune('l') {
										goto l292
									}
									position++
									break
								case 'e':
									if buffer[position] != rune('e') {
										goto l292
									}
									position++
									if buffer[position] != rune('r') {
										goto l292
									}
									position++
									if buffer[position] != rune('r') {
										goto l292
									}
									position++
									if buffer[position] != rune('o') {
										goto l292
									}
									position++
									if buffer[position] != rune('r') {
										goto l292
									}
									position++
									break
								case 'w':
									if buffer[position] != rune('w') {
										goto l292
									}
									position++
									if buffer[position] != rune('a') {
										goto l292
									}
									position++
									if buffer[position] != rune('r') {
										goto l292
									}
									position++
									if buffer[position] != rune('n') {
										goto l292
									}
									position++
									break
								case 'd':
									if buffer[position] != rune('d') {
										goto l292
									}
									position++
									if buffer[position] != rune('e') {
										goto l292
									}
									position++
									if buffer[position] != rune('b') {
										goto l292
									}
									position++
									if buffer[position] != rune('u') {
										goto l292
									}
									position++
									if buffer[position] != rune('g') {
										goto l292
									}
									position++
									break
								case 'c':
									if buffer[position] != rune('c') {
										goto l292
									}
									position++
									if buffer[position] != rune('o') {
										goto l292
									}
									position++
									if buffer[position] != rune('n') {
										goto l292
									}
									position++
									if buffer[position] != rune('t') {
										goto l292
									}
									position++
									if buffer[position] != rune('a') {
										goto l292
									}
									position++
									if buffer[position] != rune('i') {
										goto l292
									}
									position++
									if buffer[position] != rune('n') {
										goto l292
									}
									position++
									if buffer[position] != rune('s') {
										goto l292
									}
									position++
									break
								case 'n':
									if buffer[position] != rune('n') {
										goto l292
									}
									position++
									if buffer[position] != rune('o') {
										goto l292
									}
									position++
									if buffer[position] != rune('t') {
										goto l292
									}
									position++
									break
								case 'a':
									if buffer[position] != rune('a') {
										goto l292
									}
									position++
									if buffer[position] != rune('n') {
										goto l292
									}
									position++
									if buffer[position] != rune('d') {
										goto l292
									}
									position++
									break
								case 'o':
									if buffer[position] != rune('o') {
										goto l292
									}
									position++
									if buffer[position] != rune('r') {
										goto l292
									}
									position++
									break
								default:
									if buffer[position] != rune('p') {
										goto l292
									}
									position++
									if buffer[position] != rune('r') {
										goto l292
									}
									position++
									if buffer[position] != rune('o') {
										goto l292
									}
									position++
									if buffer[position] != rune('d') {
										goto l292
									}
									position++
									if buffer[position] != rune('u') {
										goto l292
									}
									position++
									if buffer[position] != rune('c') {
										goto l292
									}
									position++
									if buffer[position] != rune('t') {
										goto l292
									}
									position++
									if buffer[position] != rune('i') {
										goto l292
									}
									position++
									if buffer[position] != rune('o') {
										goto l292
									}
									position++
									if buffer[position] != rune('n') {
										goto l292
									}
									position++
									break
								}
							}

						}
					l294:
						{
							position300, tokenIndex300 := position, tokenIndex
							if !_rules[ruleIdChar]() {
								goto l300
							}
							goto l292
						l300:
							position, tokenIndex = position300, tokenIndex300
						}
						add(ruleKeyword, position293)
					}
					goto l290
				l292:
					position, tokenIndex = position292, tokenIndex292
				}
				{
					position301 := position
					{
						position302 := position
						{
							switch buffer[position] {
							case '_':
								if buffer[position] != rune('_') {
									goto l290
								}
								position++
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l290
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l290
								}
								position++
								break
							}
						}

						add(ruleIdCharNoDigit, position302)
					}
				l304:
					{
						position305, tokenIndex305 := position, tokenIndex
						if !_rules[ruleIdChar]() {
							goto l305
						}
						goto l304
					l305:
						position, tokenIndex = position305, tokenIndex305
					}
					add(rulePegText, position301)
				}
				if !_rules[rule_]() {
					goto l290
				}
				add(ruleId, position291)
			}
			return true
		l290:
			position, tokenIndex = position290, tokenIndex290
			return false
		},
		/* 42 IdChar <- <((&('_') '_') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position306, tokenIndex306 := position, tokenIndex
			{
				position307 := position
				{
					switch buffer[position] {
					case '_':
						if buffer[position] != rune('_') {
							goto l306
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l306
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l306
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l306
						}
						position++
						break
					}
				}

				add(ruleIdChar, position307)
			}
			return true
		l306:
			position, tokenIndex = position306, tokenIndex306
			return false
		},
		/* 43 IdCharNoDigit <- <((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 44 Severity <- <((&('f') FATAL) | (&('e') ERROR) | (&('w') WARN) | (&('i') INFO) | (&('d') DEBUG))> */
		nil,
		/* 45 IN <- <('i' 'n' !IdChar _)> */
		func() bool {
			position311, tokenIndex311 := position, tokenIndex
			{
				position312 := position
				if buffer[position] != rune('i') {
					goto l311
				}
				position++
				if buffer[position] != rune('n') {
					goto l311
				}
				position++
				{
					position313, tokenIndex313 := position, tokenIndex
					if !_rules[ruleIdChar]() {
						goto l313
					}
					goto l311
				l313:
					position, tokenIndex = position313, tokenIndex313
				}
				if !_rules[rule_]() {
					goto l311
				}
				add(ruleIN, position312)
			}
			return true
		l311:
			position, tokenIndex = position311, tokenIndex311
			return false
		},
		/* 46 OR <- <('o' 'r' !IdChar _)> */
		nil,
		/* 47 AND <- <('a' 'n' 'd' !IdChar _)> */
		nil,
		/* 48 NOT <- <('n' 'o' 't' !IdChar _)> */
		func() bool {
			position316, tokenIndex316 := position, tokenIndex
			{
				position317 := position
				if buffer[position] != rune('n') {
					goto l316
				}
				position++
				if buffer[position] != rune('o') {
					goto l316
				}
				position++
				if buffer[position] != rune('t') {
					goto l316
				}
				position++
				{
					position318, tokenIndex318 := position, tokenIndex
					if !_rules[ruleIdChar]() {
						goto l318
					}
					goto l316
				l318:
					position, tokenIndex = position318, tokenIndex318
				}
				if !_rules[rule_]() {
					goto l316
				}
				add(ruleNOT, position317)
			}
			return true
		l316:
			position, tokenIndex = position316, tokenIndex316
			return false
		},
		/* 49 CONTAINS <- <('c' 'o' 'n' 't' 'a' 'i' 'n' 's' !IdChar _)> */
		nil,
		/* 50 DEBUG <- <(<('d' 'e' 'b' 'u' 'g')> !IdChar _)> */
		nil,
		/* 51 INFO <- <(<('i' 'n' 'f' 'o')> !IdChar _)> */
		nil,
		/* 52 WARN <- <(<('w' 'a' 'r' 'n')> !IdChar _)> */
		nil,
		/* 53 ERROR <- <(<('e' 'r' 'r' 'o' 'r')> !IdChar _)> */
		nil,
		/* 54 FATAL <- <(<('f' 'a' 't' 'a' 'l')> !IdChar _)> */
		nil,
		/* 55 Keyword <- <((('s' 't' 'a' 'g' 'i' 'n' 'g') / ('d' 'e' 'v' 'e' 'l' 'o' 'p' 'm' 'e' 'n' 't') / ('i' 'n' 'f' 'o') / ('m' 'b') / ((&('s') 's') | (&('m') ('m' 's')) | (&('b') 'b') | (&('k') ('k' 'b')) | (&('g') ('g' 'b')) | (&('i') ('i' 'n')) | (&('f') ('f' 'a' 't' 'a' 'l')) | (&('e') ('e' 'r' 'r' 'o' 'r')) | (&('w') ('w' 'a' 'r' 'n')) | (&('d') ('d' 'e' 'b' 'u' 'g')) | (&('c') ('c' 'o' 'n' 't' 'a' 'i' 'n' 's')) | (&('n') ('n' 'o' 't')) | (&('a') ('a' 'n' 'd')) | (&('o') ('o' 'r')) | (&('p') ('p' 'r' 'o' 'd' 'u' 'c' 't' 'i' 'o' 'n')))) !IdChar)> */
		nil,
		/* 56 EQ <- <('=' _)> */
		nil,
		/* 57 LBRK <- <('[' _)> */
		nil,
		/* 58 RBRK <- <(']' _)> */
		nil,
		/* 59 LPAR <- <('(' _)> */
		func() bool {
			position329, tokenIndex329 := position, tokenIndex
			{
				position330 := position
				if buffer[position] != rune('(') {
					goto l329
				}
				position++
				if !_rules[rule_]() {
					goto l329
				}
				add(ruleLPAR, position330)
			}
			return true
		l329:
			position, tokenIndex = position329, tokenIndex329
			return false
		},
		/* 60 RPAR <- <(')' _)> */
		func() bool {
			position331, tokenIndex331 := position, tokenIndex
			{
				position332 := position
				if buffer[position] != rune(')') {
					goto l331
				}
				position++
				if !_rules[rule_]() {
					goto l331
				}
				add(ruleRPAR, position332)
			}
			return true
		l331:
			position, tokenIndex = position331, tokenIndex331
			return false
		},
		/* 61 DOT <- <('.' _)> */
		nil,
		/* 62 BANG <- <('!' !'=' _)> */
		nil,
		/* 63 LT <- <('<' !'=' _)> */
		nil,
		/* 64 GT <- <('>' !'=' _)> */
		nil,
		/* 65 LE <- <('<' '=' _)> */
		nil,
		/* 66 EQEQ <- <('=' '=' _)> */
		nil,
		/* 67 GE <- <('>' '=' _)> */
		nil,
		/* 68 NE <- <('!' '=' _)> */
		nil,
		/* 69 ANDAND <- <('&' '&' _)> */
		nil,
		/* 70 OROR <- <('|' '|' _)> */
		nil,
		/* 71 COMMA <- <(',' _)> */
		nil,
		/* 72 _ <- <Whitespace*> */
		func() bool {
			{
				position345 := position
			l346:
				{
					position347, tokenIndex347 := position, tokenIndex
					{
						position348 := position
						{
							switch buffer[position] {
							case '\t':
								if buffer[position] != rune('\t') {
									goto l347
								}
								position++
								break
							case ' ':
								if buffer[position] != rune(' ') {
									goto l347
								}
								position++
								break
							default:
								{
									position350 := position
									{
										position351, tokenIndex351 := position, tokenIndex
										if buffer[position] != rune('\r') {
											goto l352
										}
										position++
										if buffer[position] != rune('\n') {
											goto l352
										}
										position++
										goto l351
									l352:
										position, tokenIndex = position351, tokenIndex351
										if buffer[position] != rune('\n') {
											goto l353
										}
										position++
										goto l351
									l353:
										position, tokenIndex = position351, tokenIndex351
										if buffer[position] != rune('\r') {
											goto l347
										}
										position++
									}
								l351:
									add(ruleEOL, position350)
								}
								break
							}
						}

						add(ruleWhitespace, position348)
					}
					goto l346
				l347:
					position, tokenIndex = position347, tokenIndex347
				}
				add(rule_, position345)
			}
			return true
		},
		/* 73 Whitespace <- <((&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') EOL))> */
		nil,
		/* 74 EOL <- <(('\r' '\n') / '\n' / '\r')> */
		nil,
		/* 75 EOF <- <!.> */
		nil,
		/* 77 Action0 <- <{ p.AddNumber(text) }> */
		nil,
		/* 78 Action1 <- <{ p.AddNumber("")   }> */
		nil,
		/* 79 Action2 <- <{ p.AddLevel(text)  }> */
		nil,
		/* 80 Action3 <- <{ p.AddStage(text)  }> */
		nil,
		/* 81 Action4 <- <{ p.AddField(text)  }> */
		nil,
		/* 82 Action5 <- <{ p.AddString(text) }> */
		nil,
		/* 83 Action6 <- <{ p.AddExpr()       }> */
		nil,
		/* 84 Action7 <- <{ p.AddTupleValue() }> */
		nil,
		/* 85 Action8 <- <{ p.AddTupleValue() }> */
		nil,
		/* 86 Action9 <- <{ p.AddTuple() }> */
		nil,
		/* 87 Action10 <- <{ p.AddBinary(ast.IN) }> */
		nil,
		/* 88 Action11 <- <{ p.AddTuple() }> */
		nil,
		/* 89 Action12 <- <{ p.AddBinary(ast.IN); p.AddUnary(ast.LNOT) }> */
		nil,
		/* 90 Action13 <- <{ p.AddMember(text)    }> */
		nil,
		/* 91 Action14 <- <{ p.AddSubscript(text) }> */
		nil,
		/* 92 Action15 <- <{ p.AddUnary(ast.NOT) }> */
		nil,
		/* 93 Action16 <- <{ p.AddBinary(ast.GE) }> */
		nil,
		/* 94 Action17 <- <{ p.AddBinary(ast.GT) }> */
		nil,
		/* 95 Action18 <- <{ p.AddBinary(ast.LE) }> */
		nil,
		/* 96 Action19 <- <{ p.AddBinary(ast.LT) }> */
		nil,
		/* 97 Action20 <- <{ p.AddBinary(ast.EQ)   }> */
		nil,
		/* 98 Action21 <- <{ p.AddBinary(ast.NE)   }> */
		nil,
		/* 99 Action22 <- <{ p.AddBinary(ast.EQ)   }> */
		nil,
		/* 100 Action23 <- <{ p.AddBinaryContains() }> */
		nil,
		/* 101 Action24 <- <{ p.AddBinary(ast.AND) }> */
		nil,
		/* 102 Action25 <- <{ p.AddBinary(ast.AND) }> */
		nil,
		/* 103 Action26 <- <{ p.AddBinary(ast.AND) }> */
		nil,
		/* 104 Action27 <- <{ p.AddBinary(ast.OR) }> */
		nil,
		/* 105 Action28 <- <{ p.AddBinary(ast.OR) }> */
		nil,
		/* 106 Action29 <- <{ p.AddUnary(ast.LNOT) }> */
		nil,
		nil,
		/* 108 Action30 <- <{ p.SetNumber(text) }> */
		nil,
	}
	p.rules = _rules
}
