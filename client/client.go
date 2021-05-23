package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/olivere/elastic/v7"
)

type Client interface {
	Count(context.Context, interface{}) (string, error)
	Search(context.Context, interface{}, int, int) (string, error)
}

type esClient struct {
	c *elastic.Client

	indexName string
	address   string
}

func NewEsClient(addr, idx string) (*esClient, error) {
	c := &esClient{
		address:   addr,
		indexName: idx,
	}

	options := []elastic.ClientOptionFunc{
		elastic.SetSniff(false),
	}
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	if u.User != nil {
		pwd, _ := u.User.Password()
		options = append(options, elastic.SetBasicAuth(u.User.Username(), pwd))
		u.User = nil
	}
	options = append(options, elastic.SetURL(u.String()))

	c.c, err = elastic.NewClient(options...)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *esClient) Count(ctx context.Context, query interface{}) (string, error) {
	eq, ok := query.(elastic.Query)
	if !ok {
		return "", fmt.Errorf("invalide query")
	}

	res, err := c.c.Count(c.indexName).Query(eq).Do(ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Count: %d\n", res), nil
}

func (c *esClient) Search(ctx context.Context, query interface{}, size, offset int) (string, error) {
	eq, ok := query.(elastic.Query)
	if !ok {
		return "", fmt.Errorf("invalide query")
	}

	res, err := c.c.Search(c.indexName).Query(eq).Size(size).From(offset).Do(ctx)
	if err != nil {
		return "", err
	}

	output := make([]string, 0, len(res.Hits.Hits))
	for _, h := range res.Hits.Hits {
		output = append(output, fmt.Sprintf("%s\n%s", string(h.Source), strings.Repeat("-", 64)))
	}
	return strings.Join(output, "\n"), nil
}
