package ast_test

import (
	"fmt"
	"testing"

	"github.com/jskcnsl/eesc/ast"
)

type Case struct {
	Ast     string
	Details string
}

var cases = []Case{}

func init() {
	cases = append(cases, Case{
		Ast: "field_name_1 term 1",
		Details: `{
    "bool": {
        "filter": {
            "term": {
                "field_name_1": "1"
            }
        }
    }
}`,
	}, Case{
		Ast: "field_name_1 term 1 field_name_2 term 2",
		Details: `{
    "bool": {
        "filter": [
            {
                "term": {
                    "field_name_1": "1"
                }
            },
            {
                "term": {
                    "field_name_2": "2"
                }
            }
        ]
    }
}`,
	}, Case{
		Ast: "field_name_3 terms 1 2 3",
		Details: `{
    "bool": {
        "filter": {
            "terms": {
                "field_name_3": [
                    "1",
                    "2",
                    "3"
                ]
            }
        }
    }
}`,
	}, Case{
		Ast: "field_name_4 range 100 200",
		Details: `{
    "bool": {
        "filter": {
            "range": {
                "field_name_4": {
                    "from": "100",
                    "include_lower": false,
                    "include_upper": false,
                    "to": "200"
                }
            }
        }
    }
}`,
	}, Case{
		Ast: "field_name_4 range =100 200",
		Details: `{
    "bool": {
        "filter": {
            "range": {
                "field_name_4": {
                    "from": "100",
                    "include_lower": true,
                    "include_upper": false,
                    "to": "200"
                }
            }
        }
    }
}`,
	}, Case{
		Ast: "field_name_3 terms 100 200 300 field_name_4 range 100 =200",
		Details: `{
    "bool": {
        "filter": [
            {
                "terms": {
                    "field_name_3": [
                        "100",
                        "200",
                        "300"
                    ]
                }
            },
            {
                "range": {
                    "field_name_4": {
                        "from": "100",
                        "include_lower": false,
                        "include_upper": true,
                        "to": "200"
                    }
                }
            }
        ]
    }
}`,
	})
}

func TestAst1(t *testing.T) {
	for i, c := range cases {
		a := ast.NewAst(c.Ast)
		if err := a.Resolve(); err != nil {
			fmt.Printf("case %d failed: %s\n", i, err)
			t.Fail()
			continue
		}
		if a.Details() != c.Details {
			fmt.Printf("case %d failed:\n - expect:\n%s\n + actual:\n%s\n", i, c.Details, a.Details())
			t.Fail()
			continue
		}
	}
}
