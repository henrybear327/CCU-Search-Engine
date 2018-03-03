from datetime import datetime
from elasticsearch import Elasticsearch

if __name__ == "__main__":
    es = Elasticsearch()

    print(es.info)
