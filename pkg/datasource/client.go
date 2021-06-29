package datasource

import (
	"net/http"
	"context"
	"fmt"
	"io/ioutil"
	dac "github.com/xinsnake/go-http-digest-auth-client"

	"github.com/valiton/grafana-mongodb-atlas-datasource/pkg/models"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type MongoDBAtlasClient struct {
	settings *models.Settings
}

func (c *MongoDBAtlasClient) query(ctx context.Context, path string, query map[string]string) ([]byte, error) {
	apiType := c.settings.ApiType
	if apiType != "atlas" && apiType != "public" {
		apiType = "atlas"
	}
	
	var method = "GET"
	var baseURL = "https://cloud.mongodb.com/api/" + apiType + "/v1.0"
	var uri = baseURL + path

	log.DefaultLogger.Debug("MakeHttpRequest", "URL", uri)

	var t = dac.NewTransport(c.settings.PublicKey, c.settings.PrivateKey)
	req, err := http.NewRequest(method, uri, nil)

	if query != nil {
		q := req.URL.Query()
		for key, value := range query {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	log.DefaultLogger.Debug("MakeHttpRequest", "Full URL", req.URL.RequestURI())

	if err != nil {
		return nil, err
	}

	resp, err := t.RoundTrip(req)

	if err != nil {
		log.DefaultLogger.Info("MakeHttpRequest", "error", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.DefaultLogger.Debug("MakeHttpRequest", "io error", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("invalid status code. status: %v", resp.Status)
		log.DefaultLogger.Debug("MakeHttpRequest", "status code error", err)
		return nil, err
	}

	log.DefaultLogger.Debug("MakeHttpRequest", "body", string(body))

	return body, nil
}
