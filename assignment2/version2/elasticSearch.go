package main

import (
	"context"
	"log"
	"time"

	"github.com/olivere/elastic"
)

type pageRecord struct {
	Tld      string `json:"tld"`
	Link     string `json:"link"`
	Title    string `json:"title,omitempty"`
	MainText string `json:"mainText,omitempty"`

	Created time.Time             `json:"created,omitempty"`
	Suggest *elastic.SuggestField `json:"suggest_field,omitempty"`
}

type elasticSearchStorage struct {
	client *elastic.Client
	ctx    *context.Context
}

// TODO: to string?
const mapping = `
{
	"mappings":{
		"crawler" :{
			"properties": {
				"tld":{
					"type":"text"
				},
				"link":{
					"type":"text"
				},
				"title":{
					"type":"text",
					"analyzer": "ik_max_word",
					"search_analyzer": "ik_max_word"
				},
				"mainText":{
					"type":"text",
					"analyzer": "ik_max_word",
					"search_analyzer": "ik_max_word"
				}
			}
		}
	}
}`

func (es *elasticSearchStorage) init() {
	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()
	es.ctx = &ctx

	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	client, err := elastic.NewClient()
	if err != nil {
		// Handle error
		panic(err)
	}
	es.client = client

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := es.client.ElasticsearchVersion("http://127.0.0.1:9200")
	if err != nil {
		// Handle error
		panic(err)
	}
	log.Printf("Elasticsearch version %s\n", esversion)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := es.client.IndexExists(conf.MongoDB.Database).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := es.client.CreateIndex(conf.MongoDB.Database).BodyString(mapping).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
}

func (es *elasticSearchStorage) insert(tld, link, title, mainText string) {
	record := pageRecord{
		Tld:      tld,
		Link:     link,
		Title:    title,
		MainText: mainText,
	}
	put, err := es.client.Index().
		Index(conf.MongoDB.Database).
		Type("crawler").
		BodyJson(record).
		Do(*es.ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	log.Printf("Link %s: Indexed record %s to index %s, type %s\n", link, put.Id, put.Index, put.Type)
}
