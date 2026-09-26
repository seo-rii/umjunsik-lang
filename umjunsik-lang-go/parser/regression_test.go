package parser

import (
	"fmt"
	"strings"
	"testing"
	"umjunsik-lang/umjunsik-lang-go/ast"
	"umjunsik-lang/umjunsik-lang-go/lexer"
)

func TestSourceLineSlots(t *testing.T) {
	cases := []struct {
		name string
		body []string
	}{
		{"header-only", nil},
		{"consecutive-blanks", []string{"", "", "식.!", "", "식..!"}},
		{"empty-assignment", []string{"엄", "식어!"}},
		{"empty-assignment-before-blank", []string{"엄", "", "식어!"}},
		{"empty-indexed-assignment", []string{"어엄", "식어어!"}},
		{"empty-indexed-assignment-before-blank", []string{"어엄", "", "식어어!"}},
	}
	for _, tc := range cases {
		for _, sep := range []struct{ name, text string }{{"lf", "\n"}, {"crlf", "\r\n"}, {"tilde", "~"}} {
			t.Run(tc.name+"/"+sep.name, func(t *testing.T) {
				lines := append([]string{"어떻게"}, tc.body...)
				lines = append(lines, "이 사람이름이냐ㅋㅋ")
				program := New(lexer.New(strings.Join(lines, sep.text))).ParseProgram()
				// The footer is EOF; all preceding lines, including the header,
				// must retain their original positions in the program.
				if got, want := len(program.Lines), len(lines)-1; got != want {
					t.Fatalf("line slots = %d; want %d", got, want)
				}
				for i, raw := range lines[:len(lines)-1] {
					line, ok := program.Lines[i].(*ast.ExpressionLine)
					if !ok {
						t.Fatalf("source line %d has type %T", i+1, program.Lines[i])
					}
					blank := i == 0 || raw == ""
					if (line.Expression == nil) != blank {
						t.Fatalf("source line %d (%q): expression=%T; blank=%v", i+1, raw, line.Expression, blank)
					}
				}
			})
		}
	}
}

func TestInfixIntegerSuffixTokens(t *testing.T) {
	for _, operand := range []string{"어", "어어", "식?"} {
		t.Run(operand, func(t *testing.T) {
			for length := 1; length <= 8; length++ {
				for mask := 0; mask < 1<<length; mask++ {
					suffix := []byte(strings.Repeat(".", length))
					for bit := range suffix {
						if mask&(1<<bit) != 0 {
							suffix[bit] = ','
						}
					}
					text := string(suffix)
					p := New(lexer.New(fmt.Sprintf("어떻게\n%s%s\n이 사람이름이냐ㅋㅋ", operand, text)))
					program := p.ParseProgram()
					if len(p.Errors()) != 0 || len(program.Lines) != 2 {
						t.Fatalf("%s%s: errors=%v, line slots=%d", operand, text, p.Errors(), len(program.Lines))
					}
					line := program.Lines[1].(*ast.ExpressionLine)
					expr, ok := line.Expression.(*ast.InfixIntegerExpression)
					if !ok {
						t.Fatalf("%s%s: expression type %T", operand, text, line.Expression)
					}
					want := int64(strings.Count(text, ".") - strings.Count(text, ","))
					if expr.Right != want {
						t.Fatalf("%s%s: offset=%d; want %d", operand, text, expr.Right, want)
					}
				}
			}
		})
	}
}

func TestBlankLinesDoNotProduceParserErrors(t *testing.T) {
	for _, sep := range []struct{ name, text string }{{"lf", "\n"}, {"crlf", "\r\n"}, {"tilde", "~"}} {
		t.Run(sep.name, func(t *testing.T) {
			source := strings.Join([]string{"어떻게", "", "", "식.!", "", "이 사람이름이냐ㅋㅋ"}, sep.text)
			p := New(lexer.New(source))
			p.ParseProgram()
			if len(p.Errors()) != 0 {
				t.Fatalf("blank source lines must not produce parser errors: %v", p.Errors())
			}
		})
	}
}
