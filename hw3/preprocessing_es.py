# https://github.com/amoshyc/ccu-search-engine/blob/master/hw1/build-data.py

import pathlib
from tqdm import tqdm
from elasticsearch import Elasticsearch
from elasticsearch import helpers
import json

DOC_DIR = pathlib.Path('./data/jieba').resolve()
DOC_PATHS = sorted(DOC_DIR.glob('ettoday_*.txt'))
# DOC_PATHS = sorted(DOC_DIR.glob('example.rec'))
BATCH_SIZE = 2000

es = Elasticsearch()
# curl -XDELETE 'localhost:9200/ettoday?pretty'
es.indices.delete(index='ettoday', ignore_unavailable=True)

cnt = 0
for path in tqdm(DOC_PATHS):
    # Reading data back
    with path.open('r', encoding='UTF-8') as f:
        data = f.readlines()

    actions = []  # batch operation
    for record in data:
        actions.append({
            '_index': 'ettoday',
            '_type': 'news',
            '_id': cnt,
            '_source': record,
        })
        cnt += 1
        if len(actions) == BATCH_SIZE:
            helpers.bulk(es, actions)
            actions = []

    if len(actions) > 0:
        helpers.bulk(es, actions)
        pass

print('total:', cnt)
