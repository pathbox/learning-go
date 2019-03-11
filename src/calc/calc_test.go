package calc

import (
	"testing"
)

func TestFinal(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		res := Calc(tt.input)
		if res != tt.expected {
			t.Errorf("Wrong answer, got=%d, want=%d", res, tt.expected)
		}
	}
}

func TestTokenizer(t *testing.T) {
	input := `(5 + -10 * 2 + 15 / 3) * 2`
	tests := []struct {
		expectedType    string
		expectedLiteral string
	}{
		{LPAREN, "("},
		{INT, "5"},
		{PLUS, "+"},
		{MINUS, "-"},
		{INT, "10"},
		{ASTERISK, "*"},
		{INT, "2"},
		{PLUS, "+"},
		{INT, "15"},
		{SLASH, "/"},
		{INT, "3"},
		{RPAREN, ")"},
		{ASTERISK, "*"},
		{INT, "2"},
	}

	l := NewLex(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}

}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "250"
	var expectValue int64 = 250

	l := NewLex(input)
	p := NewParser(l)

	checkParseErrors(t, p)
	expression := p.ParseExpression(LOWEST)
	testInterLiteral(t, expression, expectValue)
}

func TestParsingPrefixExpression(t *testing.T) {
	input := "-15"
	expectedOp := "-"
	var expectedValue int64 = 15

	l := NewLex(input)
	p := NewParser(l)
	checkParseErrors(t, p)

	expression := p.ParseExpression(LOWEST)
	exp, ok := expression.(*PrefixExpression)

	if !ok {
		t.Fatalf("stmt is not PrefixExpression, got=%T", exp)
	}

	if exp.Operator != expectedOp {
		t.Fatalf("exp.Operator is not %s, go=%s", expectedOp, exp.Operator)
	}

	testInterLiteral(t, exp.Right, expectedValue)
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
	}

	for _, tt := range infixTests {
		l := NewLex(tt.input)
		p := NewParser(l)
		checkParseErrors(t, p)

		expression := p.ParseExpression(LOWEST)
		exp, ok := expression.(*InfixExpression)

		if !ok {
			t.Fatalf("exp is not InfixExpression, got=%T", exp)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s, go=%s", tt.operator, exp.Operator)
		}

		testInterLiteral(t, exp.Left, tt.leftValue)
		testInterLiteral(t, exp.Right, tt.rightValue)
	}
}

func testInterLiteral(t *testing.T, il Expression, value int64) bool {
	integ, ok := il.(*IntegerLiteralExpression)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	return true
}

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
