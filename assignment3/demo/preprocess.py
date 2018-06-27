import json
import pathlib
import time

import requests
from tqdm import tqdm

DOC_DIR = pathlib.Path('../Kernel2/testcase/ettoday').resolve()
DOC_PATHS = sorted(DOC_DIR.glob('ettoday*.rec'))


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
    # Reading data back
    with path.open('r', encoding='UTF-8') as f:
        data = f.readlines()

    for record in extract_record(data):
        cnt += 1

        url = 'http://localhost:8001/insert'
        payload = {
            "title": record["title"],
            "body": record['body'],
            "url": record['url']
        }
        r = requests.post(url, data=json.dumps(payload))
        # print(r.json())

        # time.sleep(1)

print('total:', cnt)
