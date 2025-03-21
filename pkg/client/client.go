package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/voidshard/faction/pkg/kind"
	"github.com/voidshard/faction/pkg/structs/api"
	v1 "github.com/voidshard/faction/pkg/structs/v1"
	"github.com/voidshard/faction/pkg/util/log"
)

var (
	kindWorld = kind.KindOf(&v1.World{})
	kindActor = kind.KindOf(&v1.Actor{})
)

type Client struct {
	cfg        *Config
	httpClient *http.Client
}

func New(cfg *Config) (*Client, error) {
	if cfg == nil {
		cfg = NewConfig()
	}
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}
	httpClient := &http.Client{Transport: tr}
	return &Client{
		cfg:        cfg,
		httpClient: httpClient,
	}, nil
}

func (c *Client) Search(world, kind string, limit int64) *searchBuilder {
	q := v1.NewQuery()
	q.Kind = kind
	q.Limit = limit
	return &searchBuilder{
		client: c,
		world:  world,
		Req:    &api.SearchRequest{Query: *q},
	}
}

func (c *Client) search(world string, req *api.SearchRequest) (*api.SearchResponse, error) {
	resp, err := c.doRequest(fmt.Sprintf("%s/search", world), "GET", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	searchresp := &api.SearchResponse{}
	err = json.NewDecoder(resp.Body).Decode(searchresp)
	if err != nil {
		return nil, err
	}

	if searchresp.Error != nil {
		if searchresp.Error.Code != 0 {
			return nil, fmt.Errorf("error code: %d, message: %s", searchresp.Error.Code, searchresp.Error.Message)
		}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return searchresp, nil
}

func (c *Client) Watch() *watchBuilder {
	return &watchBuilder{
		client: c,
		Req:    &api.StreamEvents{},
	}
}

func (c *Client) Defer(kind, world, id string) *deferEventBuilder {
	return &deferEventBuilder{
		client: c,
		Req: &api.DeferEventRequest{
			Kind:  kind,
			World: world,
			Id:    id,
		},
	}
}

func (c *Client) doDefer(req *api.DeferEventRequest) (*api.DeferEventResponse, error) {
	deferresp := &api.DeferEventResponse{}

	resp, err := c.doRequest("event", "POST", req)
	if err != nil {
		return deferresp, err
	}

	err = json.NewDecoder(resp.Body).Decode(deferresp)
	if err != nil {
		return deferresp, err
	}

	if deferresp.Error != nil {
		if deferresp.Error.Code != 0 {
			return deferresp, fmt.Errorf("error code: %d, message: %s", deferresp.Error.Code, deferresp.Error.Message)
		}
	}
	if resp.StatusCode != http.StatusOK {
		return deferresp, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return deferresp, nil
}

func (c *Client) Delete(kind, world string, ids []string) error {
	return c.delete(kind, &api.DeleteRequest{
		Ids:   ids,
		World: world,
	})
}

func (c *Client) delete(k string, req *api.DeleteRequest) error {
	resp, err := c.doRequest(k, "DELETE", req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	delresp := &api.DeleteResponse{}
	err = json.NewDecoder(resp.Body).Decode(delresp)
	if err != nil {
		return err
	}

	if delresp.Error != nil {
		if delresp.Error.Code != 0 {
			return fmt.Errorf("error code: %d, message: %s", delresp.Error.Code, delresp.Error.Message)
		}
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) Set(in []v1.Object) error {
	byKind := map[string]map[string][]interface{}{}
	for _, obj := range in {
		k := obj.GetKind()
		byWorld, ok := byKind[k]
		if !ok {
			byWorld = map[string][]interface{}{}
		}

		objects, ok := byWorld[obj.GetWorld()]
		if !ok {
			objects = []interface{}{}
		}

		objects = append(objects, obj)
		byWorld[obj.GetWorld()] = objects
		byKind[k] = byWorld
	}
	for k, byWorld := range byKind {
		for world, objects := range byWorld {
			req := api.NewSetRequest()
			req.Data = objects
			if !kind.IsGlobal(k) {
				req.World = world
			}
			err := c.set(k, req)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Client) set(k string, req *api.SetRequest) error {
	resp, err := c.doRequest(k, "POST", req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	setresp := &api.SetResponse{}
	err = json.NewDecoder(resp.Body).Decode(setresp)
	if err != nil {
		return err
	}

	if setresp.Error != nil {
		if setresp.Error.Code != 0 {
			return fmt.Errorf("error code: %d, message: %s", setresp.Error.Code, setresp.Error.Message)
		}
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) Get() *getBuilder {
	return &getBuilder{
		client: c,
		Req:    api.NewGetRequest(),
	}
}

func (c *Client) doGet(k string, req *api.GetRequest) (*api.GetResponse, error) {
	resp, err := c.doRequest(k, "GET", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	getresp := &api.GetResponse{}
	err = json.NewDecoder(resp.Body).Decode(getresp)
	if err != nil {
		return nil, err
	}

	if getresp.Error != nil {
		if getresp.Error.Code != 0 {
			return nil, fmt.Errorf("error code: %d, message: %s", getresp.Error.Code, getresp.Error.Message)
		}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return getresp, nil
}

func (c *Client) doRequest(k, method string, req interface{}) (*http.Response, error) {
	// prepare request data
	data, err := json.Marshal(req)
	log.Debug().Str("kind", k).Str("method", method).Str("data", string(data)).Msg("sending request")
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(data)
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port),
		Path:   fmt.Sprintf("/v1/%s", k),
	}

	// do the request
	httpreq, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	httpreq.Header.Set("Content-Type", "application/json")

	log.Debug().Str("kind", k).Str("method", method).Str("url", u.String()).Msg("sending request")
	resp, err := c.httpClient.Do(httpreq)
	if err != nil {
		return nil, err
	}

	if resp.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}

	return resp, nil
}
