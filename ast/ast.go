package ast

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/olivere/elastic/v7"
)

const (
	termOper  string = "term"
	termsOper string = "terms"
	rangeOper string = "range"
)

const (
	fieldStage = iota
	operStage
)

type Ast struct {
	origin string
	parts  []string
	sep    string
	q      *elastic.BoolQuery
}

func NewAst(exp string) *Ast {
	return &Ast{
		origin: exp,
		sep:    " ",
		q:      elastic.NewBoolQuery(),
	}
}

func (a *Ast) Resolve() error {
	a.parts = strings.Split(a.origin, a.sep)
	var (
		err     error
		nowFn   string
		step    int
		stage   = fieldStage
		nowWord string
	)

	for i := 0; i < len(a.parts); i += step {
		nowWord, step, err = parseWord(a.parts[i:])
		if err != nil {
			return fmt.Errorf("parse fieldname failed: %s", err)
		}
		if nowWord == "" {
			continue
		}

		switch stage {
		case fieldStage:
			nowFn = nowWord
			stage = operStage
		case operStage:
			step = 1
			if aStep, err := a.handleOper(nowFn, nowWord, i+step); err != nil {
				return err
			} else {
				step += aStep
			}
			stage = fieldStage
		default:
			return errors.New("something happend while resolving")
		}
	}

	return nil
}

func (a *Ast) handleOper(fn, oper string, argsIdx int) (int, error) {
	var (
		q    elastic.Query
		step int
		err  error
	)
	switch oper {
	case termOper:
		q, step, err = a.handleTerm(fn, argsIdx)
	case termsOper:
		q, step, err = a.handleTerms(fn, argsIdx)
	case rangeOper:
		q, step, err = a.handleRange(fn, argsIdx)
	default:
		return 0, fmt.Errorf("unspport operation %s", oper)
	}

	if err != nil {
		return 0, err
	}

	a.q = a.q.Filter(q)
	return step, nil
}

func (a *Ast) handleTerm(fn string, argsIdx int) (elastic.Query, int, error) {
	if argsIdx >= len(a.parts) {
		return nil, 0, errors.New("no args for term")
	}
	return elastic.NewTermQuery(fn, a.parts[argsIdx]), 1, nil
}

func (a *Ast) handleTerms(fn string, argsIdx int) (elastic.Query, int, error) {
	if argsIdx >= len(a.parts) {
		return nil, 0, errors.New("no args for terms")
	}
	step := 0

	args := []interface{}{}
	for {
		arg, nowStep, err := parseWord(a.parts[argsIdx+step:])
		if err != nil {
			return nil, 0, err
		}
		step += nowStep
		args = append(args, arg)

		// can we have next args ?
		if (argsIdx + step) >= len(a.parts) {

			break
		} else if isOperation(a.parts[argsIdx+step]) {
			// next args is operation, so this one is field name
			step -= nowStep
			args = args[:len(args)-1]
			break
		}
	}

	return elastic.NewTermsQuery(fn, args...), step, nil
}

func (a *Ast) handleRange(fn string, argsIdx int) (elastic.Query, int, error) {
	if argsIdx+1 >= len(a.parts) {
		return nil, 0, errors.New("no args for range")
	}
	q := elastic.NewRangeQuery(fn)
	if strings.HasPrefix(a.parts[argsIdx], "=") {
		q = q.Gte(a.parts[argsIdx][1:])
	} else {
		q = q.Gt(a.parts[argsIdx])
	}

	if strings.HasPrefix(a.parts[argsIdx+1], "=") {
		q = q.Lte(a.parts[argsIdx+1][1:])
	} else {
		q = q.Lt(a.parts[argsIdx+1])
	}

	return q, 2, nil
}

func (a *Ast) Query() interface{} {
	return a.q
}

func (a *Ast) Details() string {
	if a.q == nil {
		return "no query in ast"
	}

	s, err := a.q.Source()
	if err != nil {
		return ""
	}

	b, _ := json.MarshalIndent(s, "", "    ")
	return string(b)
}
