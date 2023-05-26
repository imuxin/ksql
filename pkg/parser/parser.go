package parser

import (
	"regexp"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

const (
	KEYWORD = `(?i)
	\b(
		  USE
		| SELECT
		| AS
		| FROM
		| WHERE
		| LABEL
		| TRUE
		| FALSE
		| NULL
		| NOT
		| AND
		| OR
		| IN
		| EXISTS
	)\b`
	OPERATORS = `<> | != | <= | >= | == | @ | [-+*/%,.()=<>]`
	SPACE     = `\s+`
	IDENTITY  = `[a-zA-Z][a-zA-Z0-9_\.\/]*`
	NUMBER    = `[-+]?\d*\.?\d+([eE][-+]?\d+)?`
	STRING    = `'[^']*' | "[^"]*"`
	COMMENT   = `#[^\n]+`
)

var (
	format   = func(str string) string { return regexp.MustCompile(`[\n\t\s]*`).ReplaceAllString(str, "") }
	sqlLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: `Keyword`, Pattern: format(KEYWORD)},
		{Name: `Operators`, Pattern: format(OPERATORS)},
		{Name: "whitespace", Pattern: format(SPACE)},
		{Name: `Ident`, Pattern: format(IDENTITY)},
		{Name: `Number`, Pattern: format(NUMBER)},
		{Name: `String`, Pattern: format(STRING)},
		{Name: "Comment", Pattern: format(COMMENT)},
	})

	parser = participle.MustBuild[KSQL](
		participle.Lexer(sqlLexer),
		participle.Unquote("String"),
		participle.CaseInsensitive("Keyword"),
		participle.Elide("Comment"),
	)
)

func Parse(sql string) (*KSQL, error) {
	return parser.ParseString("", sql)
}
