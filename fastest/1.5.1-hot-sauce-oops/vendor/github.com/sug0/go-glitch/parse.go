package glitch

import (
    "fmt"
    "strings"
    "unicode"
)

type Expression struct {
    infix string
    toks  []string
}

func (expr *Expression) String() string {
    return expr.infix
}

// shunting yard algorithm
func CompileExpression(input string) (exp *Expression, err error){
    defer func() {
        if r := recover(); r != nil {
            exp = nil
            err = fmt.Errorf("invalid expression: %s", input)
        }
    }()

    lastWasDigit := false
    output := ""
    opers := make([]rune, 0, len(input))
    nOpers, nOperands := 0, 0

    for _,tok := range []rune(input) {
        switch {
        default:
            return nil, fmt.Errorf("invalid expression: %s", input)
        case unicode.IsSpace(tok):
            continue
        case tok == '(':
            if lastWasDigit {
                lastWasDigit = false
                output += " "
            }
            opers = append(opers, tok)
        case tok == ')':
            if lastWasDigit {
                lastWasDigit = false
                output += " "
            }
            for {
                // pop front
                op := opers[len(opers)-1]
                opers = opers[:len(opers)-1]

                if op == '(' {
                    break
                } else {
                    output += string(op) + " "
                }
            }
        case operMap[tok] != nil:
            nOpers++
            op := operMap[tok]
            if lastWasDigit {
                lastWasDigit = false
                output += " "
            }
            for {
                if len(opers) == 0 {
                    break
                }

                opTok := opers[len(opers)-1]
                if op2, ok := operMap[opTok]; !ok || !op2.hasPrecedence(op) {
                    break
                }

                // pop front
                opers = opers[:len(opers)-1]
                output += string(opTok) + " "
            }
            opers = append(opers, tok)
        case validTok(tok):
            nOperands++
            if lastWasDigit {
                lastWasDigit = false
                output += " "
            }
            output += string(tok) + " "
        case unicode.IsDigit(rune(tok)):
            if !lastWasDigit {
                nOperands++
                lastWasDigit = true
            }
            output += string(tok)
        }
    }

    // insufficient number of operands
    if nOpers != nOperands - 1 {
        return nil, fmt.Errorf("invalid expression: %s", input)
    }

    // try to find unmatched parenthesis
    // while reversing oper order to insert
    // in postfix expression
    l := len(opers)
    rev := make([]string, l)
    l--

    for i, tok := range opers {
        if tok == '(' {
            rev = nil
            return nil, fmt.Errorf("invalid expression: %s", input)
        }
        rev[l - i] = string(tok)
    }

    // success
    return &Expression{
        infix: input,
        toks: append(strings.Split(output, " "), rev...),
    }, nil
}

func validTok(tok rune) bool {
    return tok == 'c' || tok == 's' || tok == 'Y' ||
           tok == 'r' || tok == 'x' || tok == 'y' ||
           tok == 'N' || tok == 'R' || tok == 'G' ||
           tok == 'B' || tok == 'e' || tok == 'b' ||
           tok == 'H' || tok == 'L' || tok == 'h' ||
           tok == 'v' || tok == 'd'
}
