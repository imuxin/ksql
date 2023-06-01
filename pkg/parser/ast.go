package parser

/*
SELECT Statement Definition in [mysql reference](https://dev.mysql.com/doc/refman/5.7/en/select.html).
```
SELECT
    [ALL | DISTINCT | DISTINCTROW ]
    [HIGH_PRIORITY]
    [STRAIGHT_JOIN]
    [SQL_SMALL_RESULT] [SQL_BIG_RESULT] [SQL_BUFFER_RESULT]
    [SQL_CACHE | SQL_NO_CACHE] [SQL_CALC_FOUND_ROWS]
    select_expr [, select_expr] ...
    [into_option]
    [FROM table_references
      [PARTITION partition_list]]
    [WHERE where_condition]
    [GROUP BY {col_name | expr | position}
      [ASC | DESC], ... [WITH ROLLUP]]
    [HAVING where_condition]
    [ORDER BY {col_name | expr | position}
      [ASC | DESC], ...]
    [LIMIT {[offset,] row_count | row_count OFFSET offset}]
    [PROCEDURE procedure_name(argument_list)]
    [into_option]
    [FOR UPDATE | LOCK IN SHARE MODE]

into_option: {
    INTO OUTFILE 'file_name'
        [CHARACTER SET charset_name]
        export_options
  | INTO DUMPFILE 'file_name'
  | INTO var_name [, var_name] ...
}
```
*/

// The grammar syntax definition is referenced by alecthomas/participle library
// You can visit [participle#grammar-syntax](https://github.com/alecthomas/participle#grammar-syntax) to get more details.

type Statement interface{}

type KSQL struct {
	Use    *UseStat    `parser:" @@* "`
	Select *SelectStat `parser:" @@* "`

	// TODO
	// Delete DeleteStat `parser:" @@* "`
	// Update UpdateStat `parser:" @@* "`
}

type UseStat struct {
	// Database, we consider each cluster is a database.
	Database string `parser:" 'USE' @Ident "`
}

type SelectStat struct {
	Select            SelectExpr          `parser:" 'SELECT' @@ "`
	From              FromExpr            `parser:" 'FROM' @@ "`
	Where             *WhereExpr          `parser:" ( 'WHERE' @@ )? "`
	Namespace         string              `parser:" ( ( 'NAMESPACE' | 'NS' ) ( @Ident | @String ) )? "`
	KubernetesFilters []*KubernetesFilter `parser:" @@* "`
}

// TODO
type DeleteStat struct{}
type UpdateStat struct{}

type SelectExpr struct {
	ALL     bool      `parser:" @'*' "`
	Columns []*Column `parser:" | @@ ( ',' @@ )* "`
}

type Column struct {
	Name  string `parser:"( @Ident | @String )"`
	Alias string `parser:" ( 'AS' ( @Ident | @String | @'NAMESPACE' | @'NS' | @'NAME' | @'SELECT' | @'LABEL' ) )? "`
}

type FromExpr struct {
	Table string `parser:" ( @Ident | @'NAMESPACE' | @'NS' | @'NAME' | @'SELECT' | @'LABEL' ) "`
	DB    string `parser:" ( '@' @Ident )? "`
}

type KubernetesFilter struct {
	Label *Compare `parser:"   'LABEL' @@ "`
	Name  *string  `parser:" | 'NAME' (@Ident | @String) "`
}

type WhereExpr struct {
	First      Compare      `parser:" @@ "`
	Conditions []*Condition `parser:" @@* "`
}

type Condition struct {
	Type    string  `parser:" @( 'AND' | 'OR' ) "`
	Compare Compare `parser:" @@ "`
}

type Compare struct {
	NOT       bool      `parser:" @'NOT'? "`
	LHS       string    `parser:" ( @Ident | @String ) "`
	Operation Operation `parser:" @@ "`
}

type Operation struct {
	Exists string `parser:" ( @( 'NOT'? 'EXISTS') "`
	Op     string `parser:" | @( '<>' | '<=' | '>=' | '=' | '==' | '<' | '>' | '!=' | 'NOT'? 'IN' ) "`
	RHS    Value  `parser:" @@ ) "`
}

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "TRUE"
	return nil
}

type Value struct {
	Number  *float64 ` parser:" ( @Number "`
	String  *string  ` parser:" | @String | @Ident "`
	Boolean *Boolean ` parser:" | @('TRUE' | 'FALSE') "`
	Null    bool     ` parser:" | @'NULL' "`
	Array   *Array   ` parser:" | @@ )"`
}

type Array struct {
	Value []*Value `parser:" '(' @@ ( ',' @@ )* ')' "`
}
