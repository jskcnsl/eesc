# Easy Elasticsearch Client

An client to debug Elasticsearch data easy **WITHOUT JSON**.

[![Build Status](https://github.com/jskcnsl/eesc/actions/workflows/go.yml/badge.svg)](https://github.com/jskcnsl/eesc/actions)

## Usage

```shell
Easy Elasticsearch Client

Usage:
  eesc [command]

Available Commands:
  count       get count result with query
  help        Help about any command
  search      search elasticsearch with query

Flags:
  -f, --file string     collection of exporession
  -h, --help            help for eesc
  -x, --idx string      Index to work around
  -j, --join string     query which will join to each expression
  -o, --output string   output file
  -q, --query string    single query expression
  -s, --server string   Address of database server
  -v, --verbose         show details

Use "eesc [command] --help" for more information about a command.
```

### Expression

`eesc` create elastic.Query in BoolQuery, and push every query in `filter`. This means relationship between each query you need is **AND**.

#### One expression

`field_name term abc`

#### More expression

`field_name term abc timestamp range 1621777167 1621777170`

#### Support operation

- term
- terms
- range

### Query multi expression from file

- `eesc search -f 'e.txt'` each query return one result
- `eesc search -f 'e.txt' -j 'field_name2 term bbb'` each query run with `join` query and return one result
