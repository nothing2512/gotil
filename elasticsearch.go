package gotil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/olivere/elastic/v7"
)

// Elastic Search Table Interface
type Elastable interface {
	TableName() string
}

// Base Elastic Client Struct
type ElasticClient struct {
	_client *elastic.Client
}

// Create new Elastic Client
func NewElasticSearch(uri string) (*ElasticClient, error) {
	client, err := elastic.NewClient(elastic.SetURL(uri), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}
	return &ElasticClient{client}, nil
}

// Delete Data By Model
func (ec *ElasticClient) Delete(data Elastable) error {
	body, err := ec.getObjValue(data)
	if err != nil {
		return err
	}
	return ec.DeleteById(data.TableName(), body["id"])
}

// Delete Data by Id
func (ec *ElasticClient) DeleteById(table string, id any) error {
	_, err := ec._client.Delete().
		Index(table).
		Id(fmt.Sprintf("%v", id)).
		Do(context.Background())
	return err
}

// Update Data by model
func (ec *ElasticClient) Update(data Elastable) error {
	body, err := ec.getObjValue(data)
	if err != nil {
		return err
	}
	_, err = ec._client.Update().
		Index(data.TableName()).
		Id(fmt.Sprintf("%v", body["id"])).
		Doc(body).
		Refresh("true").
		Do(context.Background())
	if err != nil {
		return err
	}
	return err
}

// search data
func (ec *ElasticClient) Search(obj any, table, value string, keys ...string) {
	sr, _ := ec._client.Search().
		Index(table).
		Query(ec.term(value, keys...)).
		Do(context.TODO())
	responses := []any{}
	for _, hit := range sr.Hits.Hits {
		var src JSON
		json.Unmarshal(hit.Source, &src)
		responses = append(responses, src)
	}
	b, _ := json.Marshal(responses)
	_ = json.Unmarshal(b, obj)
}

// term builder
func (ec *ElasticClient) term(value string, keys ...string) *elastic.BoolQuery {
	var queries []elastic.Query
	for _, v := range keys {
		q := elastic.NewFuzzyQuery(v, value).
			Fuzziness(2).
			Transpositions(true)
		queries = append(queries, q)
	}
	return elastic.NewBoolQuery().MinimumShouldMatch("1").Should(queries...)
}

// create new data
func (ec *ElasticClient) Save(data Elastable) error {
	if ec._client == nil {
		return errors.New("Disconnected")
	}
	body, err := ec.getObjValue(data)
	if err != nil {
		return err
	}
	_, err = ec._client.Index().
		Index(data.TableName()).
		Id(fmt.Sprintf("%v", body["id"])).
		BodyJson(body).
		Do(context.Background())
	return err
}

// getting object data to JSON
func (ec *ElasticClient) getObjValue(ptr any) (JSON, error) {

	v := reflect.ValueOf(ptr)

	res := JSON{}

	val := v.Elem()

	if val.Kind() != reflect.Struct {
		return nil, errors.New("invalid object")
	}

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("es")
		if tag != "" && tag != "-" {
			res[tag] = val.FieldByName(typ.Field(i).Name).Interface()
		}
	}

	if res["id"] == nil || res["id"] == "" {
		return nil, errors.New("column id not found")
	}

	return res, nil
}
