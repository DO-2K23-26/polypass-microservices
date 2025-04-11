package infrastructure

import "github.com/elastic/go-elasticsearch/v8"

type ElasticAdapter struct {
	Client *elasticsearch.Client
}

func NewElasticAdapter(host string, password string) (*ElasticAdapter, error) {
	esConfig := elasticsearch.Config{Addresses: []string{host}, Password: password}
	Client, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		return nil, err
	}
	return &ElasticAdapter{
		Client,
	}, nil
}

// Ping checks the connection to the Elasticsearch server.
// False = failed connection , True = successful connection
func (e *ElasticAdapter) Ping() bool {
	_ , err := e.Client.Ping()
	if err != nil {
		return false
	}
	return true
}