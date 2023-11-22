package client

import (
	"context"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/carlmjohnson/requests"
	"github.com/lonelycode/montag-cli/models"
)

const (
	EP_RunAIFunc      = "/api/aifunctions/call/{slug}"
	EP_VectorSearch   = "/api/bots/{id}/search"
	EP_VectorSearchNS = "/api/bots/{id}/search-namespace"
	EP_SnippetsList   = "/api//snippets"
)

type Client struct {
	key    string
	server string
}

func NewClient(key string, server string) *Client {
	return &Client{
		key:    key,
		server: server,
	}
}

func (c *Client) BaseRequest(url string) *requests.Builder {
	return requests.URL(url).Header("Authorization", c.key)
}

func (c *Client) ResourceURI(resource string) string {
	ep, err := url.JoinPath(c.server, resource)
	if err != nil {
		log.Fatal(err)
	}

	return ep
}

func (c *Client) URLWithVal(resource, key, value string) string {
	fixedResource := strings.Replace(resource, key, value, 1)
	return c.ResourceURI(fixedResource)
}

func (c *Client) RunAIFunc(name string, inputs *models.AIFuncCall) (*models.AIFuncResponse, error) {
	ep := c.URLWithVal(EP_RunAIFunc, "{slug}", name)
	var res models.AIFuncResponse
	ctx := context.Background()
	err := c.BaseRequest(ep).
		BodyJSON(inputs).
		ToJSON(&res).
		Fetch(ctx)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) VectorSearch(botID int, query string, numResults int) ([]*models.QueryMatch, error) {
	ep := c.URLWithVal(EP_VectorSearch, "{id}", strconv.Itoa(botID))
	var res []*models.QueryMatch
	ctx := context.Background()
	err := c.BaseRequest(ep).
		Param("query", query).
		Param("numResults", strconv.Itoa(numResults)).
		ToJSON(&res).
		Fetch(ctx)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) VectorSearchNS(botID int, namespace, query string, numResults int) ([]*models.QueryMatch, error) {
	ep := c.URLWithVal(EP_VectorSearchNS, "{id}", strconv.Itoa(botID))
	var res []*models.QueryMatch
	ctx := context.Background()
	err := c.BaseRequest(ep).
		Param("query", query).
		Param("namespace", namespace).
		Param("numResults", strconv.Itoa(numResults)).
		ToJSON(&res).
		Fetch(ctx)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) GetSnippet(slug string) (string, error) {
	ep := c.ResourceURI(EP_SnippetsList)
	var res []*models.Snippet
	ctx := context.Background()
	err := c.BaseRequest(ep).
		ToJSON(&res).
		Fetch(ctx)

	if err != nil {
		return "", err
	}

	for _, s := range res {
		if s.Slug == slug {
			return s.Content, nil
		}
	}

	return "", nil
}
