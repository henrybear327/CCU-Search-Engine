import requests
import json
from elasticsearch import Elasticsearch


def test_es_connection():
    res = requests.get('http://localhost:9200')
    print(res.text)


def go():
    es = Elasticsearch([{'host': 'localhost', 'port': 9200}])

    # get result based on id
    # res = es.get(index='ettoday', doc_type='news', id=5)
    # print(json.dumps(res, sort_keys=True, indent=4))

    # get result based on exact match
    res = es.search(index="ettoday", body={
        "query": {
            "match": {
                "body": "中正大學"
            }
        }
    })
    print(res)
    for record in res['hits']['hits']:
        print(record)
    # print(json.dumps(res, sort_keys=True, indent=4))


if __name__ == "__main__":
    # test_es_connection()

    go()
