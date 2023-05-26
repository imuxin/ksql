package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	const EBNF = `KSQL = UseStat* SelectStat* .
UseStat = "USE" <ident> .
SelectStat = "SELECT" SelectExpr "FROM" FromExpr ("WHERE" WhereExpr)? ("LABEL" LabelExpr)? .
SelectExpr = "*" | (Column ("," Column)*) .
Column = (<ident> | <string>) ("AS" <ident>)? .
FromExpr = <ident> ("@" <ident>)? .
WhereExpr = Compare Condition* .
Compare = (<ident> | <string>) Operation .
Operation = (("NOT"? "EXISTS") | (("<>" | "<=" | ">=" | "=" | "==" | "<" | ">" | "!=" | ("NOT"? "IN")) Value)) .
Value = (<number> | <string> | ("TRUE" | "FALSE") | "NULL" | Array) .
Array = "(" Value ("," Value)* ")" .
Condition = ("AND" | "OR" | "NOT") Compare .
LabelExpr = Compare ("," Compare)* .`
	assert.Equal(t, EBNF, parser.String())
}

func TestParse(t *testing.T) {
	const demoSQLStr = `
	select a as aa, b, "spec.name"
	    from te.st@cluster1
		where x = 1.1
		    and 'in' = 'abc'
		    and 'in' == 'abc'
			and xx = TRUE
			or abc in (1,2,3)
			or abc not in (1,2,3) # dfadf
		label cluster exists , cluster not exists , k8s.io/proj = "sample"
	`
	if _, err := Parse(demoSQLStr); err != nil {
		t.Error(err)
	}
}
