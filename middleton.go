package main

import (
	"fmt"
	"log"
	"slices"
	"strings"
	"unicode"
)

var keywords = []string{
	"package",
	"import",
}

var typeKeywords = []string{
	"MFunc",
}

type ASTNode struct {
	parent   *ASTNode
	children []interface{}
}

func (n *ASTNode) AppendChild(child interface{}) {
	n.children = append(n.children, child)
}

type RootASTNode struct {
	ASTNode
}

func NewRootASTNode() *RootASTNode {
	return &RootASTNode{ASTNode: ASTNode{parent: nil}}
}

type Expr struct {
	ASTNode
	value interface{}
}

func NewExpr(parent *ASTNode, s interface{}) *Expr {
	return &Expr{ASTNode: ASTNode{parent: parent}, value: s}
}

func isNumericCharacter(char rune) bool {
	return unicode.IsDigit(char)
}

func isAlphabeticCharacter(char rune) bool {
	return unicode.IsLetter(char)
}

func getCharType(c rune) string {
	if unicode.IsSpace(c) {
		return "whitespace"
	} else if isNumericCharacter(c) || isAlphabeticCharacter(c) {
		return "alnum"
	} else {
		return "other"
	}
}

type Token struct {
	Value string
	Type  string
	At    int
	Start int
}

func (t Token) String() string {
    return fmt.Sprintf("<%s '%s' %d:%d>", t.Type, t.Value, t.At, t.Start)
}

var TOKEN_GROUPS = []string{"++", "--", "//", "==", "!=", "/=", "%=", "-=", "+=", "&&"}

func errtok(ref Token, msg string)  {
	log.Fatal(fmt.Errorf("Fatal error on line %d!\n%s", ref.At + 1, msg))
}

func lex(source string) []Token {
	if strings.TrimSpace(string) == "" { return []Token{} }

	var tokens []Token
	inQuote := false
	var token string
	line := 0
	lastCharType := ""
	tokenStart := 0
	ignoreLine := false
	
	pushCategorizedToken := func(token string) {
		if token[0] == '"' && token[len(token)-1] == '"' {
			tokens = append(tokens, Token{Value: token, Type: "string", At: line, Start: tokenStart})
		} else {
			ct := Token{Value: token, Type: getCharType(rune(token[0])), At: line, Start: tokenStart}

			if contains(keywords, ct.Value) {
				ct.Type = "keyword"
			} else if contains(typeKeywords, ct.Value) {
				ct.Type = "type"
			}

			tokens = append(tokens, ct)
		}
	}

	pushToken := func(also string) {
		if strings.TrimSpace(token) != "" {
			pushCategorizedToken(token)
		}

		token = ""

		if also != "" {
			pushCategorizedToken(also)
		}

		tokenStart = len(token)
	}

	for i, char := range source {
		if char == '\n' {
			line++
			ignoreLine = false
		}

		if ignoreLine { continue }

		if char == '"' {
			inQuote = !inQuote
			if inQuote {
				pushToken("")
				token += string(char)
				continue
			} else {
				token += string(char)
				pushToken("")
				continue
			}
		}

		if inQuote {
			token += string(char)
			continue
		}

		charType := getCharType(char)

		if i == 0 {
			lastCharType = charType
		}

		if charType != lastCharType || charType == "other" {
			pushToken("")
			if charType != "other" {
				token += string(char)
			} else {
				token = string(char)
				pushToken("")
			}
		} else {
			token += string(char)
		}

		lastCharType = charType
	}

	pushToken("") // Push the last token

	return tokens
}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

var OPERATOR_PRECEDENCE = map[string]int{
	"=":  1,
	"||": 2,
	"&&": 3,
	"<":  7, ">": 7, "<=": 7, ">=": 7, "==": 7, "!=": 7,
	"+": 10, "-": 10,
	"*": 20, "/": 20, "%": 20,
}

func parse(tokens []Token) *RootASTNode {
	root := NewRootASTNode()

	if len(tokens) == 0 { return root }

	// Choose whether to pop or just get the token
	// Pop tokens so that the Expr can be made until there are no tokens left
	
	tok := func() Token {
		return tokens[0]
	}
	
	poptok := func() Token {
		if len(tokens) == 0 {
			log.Fatal(fmt.Errorf("MiddletonScript parse error: no tokens to consume!"))
		}
		
		popped := tokens[0]
		tokens = tokens[1:]
		return popped
	}
	
	// pop expecting: this version can contain multiple "next-tokens"
	/*
	poptokexs := func(vs []string) Token {
		if !slices.Contains(vs, tok().Value) {
			errtok(tok(), fmt.Sprintf("Got '%s' expected %s", tok(), strings.Join(vs, " or ")))
		}
		
		return poptok()
	}
	*/

	poptokex := func(v string) Token {
		if v != tok().Value {
			errtok(tok(), fmt.Sprintf("Got '%s' expected '%s'", tok().Value, v))
		}
		
		return poptok()
	}

	nextExpression := func() *Expr { return NewExpr(&root.ASTNode, nil) }
	
	parseIf := func() ASTNode {
		poptokex("if")
		poptokex("(")
		nextExpression()
		poptokex(")")
		poptokex("{")
		nextExpression()
		poptokex("}")
		
		return ASTNode{}
	}
	
	parseSwyk := func() ASTNode {
		poptokex("swyk")

		popped := poptok()		
		if popped.Type != "alnum" {
			errtok(popped, fmt.Sprintf("Expected 'alnum', got '%s'!", popped.Type))
		}
		
		return ASTNode{}
	}

	parseNotes := func() ASTNode {
		poptokex("notes")
		poptokex("(")

		for tok().Value != ")" { 
			popped := poptok()
			
			if popped.Type != "string" {
				errtok(popped, fmt.Sprintf("Expected 'string', got '%s' (%s)!", popped.Value, popped.Type))
			}
		}
		
		poptokex(")")

		return ASTNode{}
	}

	parseMFunc := func() ASTNode {
		poptokex("MFunc")
		fname := poptok()
		if fname.Type != "alnum" {
			errtok(fname, fmt.Sprintf("Invalid function name '%s'", fname.Value))
		}
		poptokex("(")
		for tok().Value != ")" {
			ptype := poptok()
			if !slices.Contains(typeKeywords, ptype.Value) {
				errtok(ptype, fmt.Sprintf("Parameter type unknown: '%s'", ptype.Value))
			}
			
			pname := poptok()
			if pname.Type != "alnum" {
				errtok(ptype, fmt.Sprintf("'%s' found, expected 'alnum'!", ptype.Type))
			}
		}
		poptokex(")")
		poptokex("{")
		nextExpression()
		poptokex("}")
		return ASTNode{}
	}

	nextExpression = func() *Expr {
		var exprv ASTNode
		
		if tok().Value == "if" {
			exprv = parseIf()
		} else if tok().Value == "swyk" {
			exprv = parseSwyk()
		} else if tok().Value == "notes" {
			exprv = parseNotes()
		} else if tok().Value == "MFunc" {
			exprv = parseMFunc()
		} else {
			errtok(tok(), fmt.Sprintf("MiddletonScript global parse error: %s '%s' found.", tok().Type, tok().Value))
		}
		
		expr := NewExpr(&root.ASTNode, exprv)
		return expr
	}

	parseRoot := func() *RootASTNode {
		for len(tokens) > 0 {
			root.AppendChild(nextExpression())
		}

		return root
	}

	return parseRoot()
}

type MiddletonInterpreter struct {
	middletonian map[string]interface{}
}

func (mi *MiddletonInterpreter) ToBytecode(source string) []byte {
	var middlebytes []byte

	tokens := lex(source)
	fmt.Println("Tokens:", tokens)

	parsed := parse(tokens)
	fmt.Println("Parsed:", parsed)

	return middlebytes
}
