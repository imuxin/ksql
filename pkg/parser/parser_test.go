package parser

import (
	"testing"

	"github.com/alecthomas/repr"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	const EBNF = `KSQL = UseStat* SelectStat* .
UseStat = "USE" <ident> .
SelectStat = "SELECT" SelectExpr "FROM" FromExpr ("WHERE" WhereExpr)? (("NAMESPACE" | "NS") (<ident> | <string>))? KubernetesFilter* .
SelectExpr = "*" | (Column ("," Column)*) .
Column = (<ident> | <string>) ("AS" (<ident> | <string> | "NAMESPACE" | "NS" | "NAME" | "SELECT" | "LABEL"))? .
FromExpr = (<ident> | "NAMESPACE" | "NS" | "NAME" | "SELECT" | "LABEL") ("@" <ident>)? .
WhereExpr = Compare Condition* .
Compare = "NOT"? (<ident> | <string>) Operation .
Operation = (("NOT"? "EXISTS") | (("<>" | "<=" | ">=" | "=" | "==" | "<" | ">" | "!=" | ("NOT"? "IN")) Value)) .
Value = (<number> | <string> | <ident> | ("TRUE" | "FALSE") | "NULL" | Array) .
Array = "(" Value ("," Value)* ")" .
Condition = ("AND" | "OR") Compare .
KubernetesFilter = ("LABEL" Compare) | ("NAME" (<ident> | <string>)) .`
	// fmt.Println(parser.String())
	assert.Equal(t, EBNF, parser.String())
}

func TestParse(t *testing.T) {
	const demoSQLStr = `
	SELECT a AS aa, b, "spec.name"
	    FROM te.st@cluster1
		WHERE NOT x = 1.1
		    AND 'in' = 'abc'
		    AND NOT 'in' = 'abc'
		    AND 'in' == 'abc'
			AND xx = TRUE
			OR abc IN (1,2,3)
			OR abc NOT IN (1,2,3) # dfadf
		# NAMESPACE istiosystem
		NS istiosystem
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
