# https://github.com/amoshyc/ccu-search-engine/blob/master/hw1/build-data.py

import pathlib
from tqdm import tqdm
from elasticsearch import Elasticsearch
from elasticsearch import helpers

DOC_DIR = pathlib.Path('./ettoday/').resolve()
DOC_PATHS = sorted(DOC_DIR.glob('et*.rec'))
# DOC_PATHS = sorted(DOC_DIR.glob('example.rec'))
BATCH_SIZE = 2000

es = Elasticsearch()
# curl -XDELETE 'localhost:9200/ettoday?pretty'
es.indices.delete(index='ettoday', ignore_unavailable=True)


def extract_record(data):
    for line in data:
        line = line.strip()
        if line.startswith('@GAISRec:'):
            record = dict()
        elif line.startswith('@U:'):
            record['url'] = line[3:]
        elif line.startswith('@T:'):
            record['title'] = line[3:]
        elif line.startswith('@B:'):
            pass
        else:
            record['body'] = line
            yield record


cnt = 0
for path in tqdm(DOC_PATHS):
    with path.open('r', encoding='UTF-8') as f:
        data = f.readlines()

    actions = []  # batch operation
    for record in extract_record(data):
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

print('total:', cnt)
