/*
Copyright 2021 kqzh.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package elastic

import (
	"bytes"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8"
)

type esClient interface {
	DeleteLogs(instanceId string) error
}

type EsClient struct {
	client *elasticsearch.Client
}

var _ esClient = (*EsClient)(nil)

func New(config elasticsearch.Config) (*EsClient, error) {
	client, err := elasticsearch.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &EsClient{client}, nil
}

func (es *EsClient) DeleteLogs(name string) error {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"kubernetes.labels.job-name": name,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return err
	}

	res, err := es.client.DeleteByQuery([]string{"filebeat-*"}, &buf)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
