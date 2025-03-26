package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"unicode"
)

/*
   What is a CFG?

       A CFG is a Context-Free-Grammar.

   How is it implemented?

       Every Expr has a value that can come next, and if it doesn't
           meet one of those values that can come next, there should be an error thrown.
           When that next value is correct, you add it to the AST (Abstract Syntax Tree),
           specifying for the later code-generation step. But in order to get to that, we
           must first parse the CFG rules to allow them to be understood in a
           "What comes next?" context.
*/

var keywords = []string{
	"package",
	"import",
}

var typeKeywords = []string{
	"MFunc",
}

type ASTNode struct {
	parent   *ASTNode
	children []*ASTNode
}

func (n *ASTNode) AppendChild(child *ASTNode) {
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
}

func NewExpr(parent *ASTNode) *Expr {
	return &Expr{ASTNode: ASTNode{parent: parent}}
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

func lex(source string) []Token {
	var tokens []Token
	inQuote := false
	var token string
	line := 0
	lastCharType := ""
	tokenStart := 0

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
		}

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

type MiddletonInterpreter struct {
	middletonian map[string]interface{}
}

func (mi *MiddletonInterpreter) ToBytecode(source string) []byte {
	var middlebytes []byte

	tokens := lex(source)
	// Parsing is not implemented here, so just output tokens
	fmt.Println("Tokens:", tokens)

	return middlebytes
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("MiddletonScript ready! Interpreting a 'Hello World' example.")
		middleton := MiddletonInterpreter{}
		middlebytes := middleton.ToBytecode(`
		swyk middleton

		notes (
			"middleton:fmt"
		)

		MFunc main() {
    			middleton::mout("Hello, MiddletonScript!");
		}
		`)
		fmt.Println(middlebytes)
		return
	}

	var flags []rune

	// Flags pass
	for _, arg := range args {
		if arg[0] == '-' {
			flags = append(flags, rune(arg[1]))
			continue
		}
	}

	for _, arg := range args {
		// Ignore flags as we have already parsed them
		if arg[0] == '-' {
			continue
		}

		data, err := os.ReadFile(arg)

		if err != nil {
			fmt.Printf("Error reading file with name '%s'!\n", arg)
			continue
		}

		middleton := MiddletonInterpreter{}
		middlebytes := middleton.ToBytecode(string(data))

		if slices.Contains(flags, 'b') {
			file, err := os.Create(arg + ".middlebytes")
			if err != nil {
				log.Fatal(fmt.Errorf("Error opening MiddletonScript bytes file! %s", err))
			}

			_, err = file.Write(middlebytes)
			if err != nil {
				log.Fatal(fmt.Errorf("Error writing to MiddletonScript bytes file! %s", err))
			}

			file.Close()
		}
	}
}
