package parser

import (
	"testing"

	"github.com/alecthomas/repr"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	const EBNF = `KSQL = UseStat* SelectStat* .
UseStat = "USE" <ident> .
SelectStat = "SELECT" SelectExpr "FROM" FromExpr ("WHERE" WhereExpr)? ("LABEL" LabelExpr)? ("NAME" (<ident> | <string>) ("," (<ident> | <string>))*)? (("NAMESPACE" <ident>) | <string>)? .
SelectExpr = "*" | (Column ("," Column)*) .
Column = (<ident> | <string>) ("AS" <ident>)? .
FromExpr = <ident> ("@" <ident>)? .
WhereExpr = Compare Condition* .
Compare = "NOT"? (<ident> | <string>) Operation .
Operation = (("NOT"? "EXISTS") | (("<>" | "<=" | ">=" | "=" | "==" | "<" | ">" | "!=" | ("NOT"? "IN")) Value)) .
Value = (<number> | <string> | ("TRUE" | "FALSE") | "NULL" | Array) .
Array = "(" Value ("," Value)* ")" .
Condition = ("AND" | "OR" | "NOT") Compare .
LabelExpr = Compare ("," Compare)* .`
	assert.Equal(t, EBNF, parser.String())
}

func TestParse(t *testing.T) {
	const demoSQLStr = `
	SELECT a AS aa, b, "spec.name"
	    FROM te.st@cluster1
		WHERE NOT x = 1.1
		    AND 'in' = 'abc'
		    AND 'in' == 'abc'
			AND xx = TRUE
			OR abc IN (1,2,3)
			OR abc NOT IN (1,2,3) # dfadf
		NAMESPACE istiosystem
		LABEL cluster EXISTS
		LABEL cluster NOT EXISTS
		LABEL k8s.io/proj = "sample"
		NAME istiod-116
		NAME envoy
	`
	if ksql, err := Parse(demoSQLStr); err != nil {
		t.Error(err)
	} else {
		repr.Println(ksql)
	}
}
