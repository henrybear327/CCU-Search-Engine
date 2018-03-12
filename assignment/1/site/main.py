from flask import Flask, render_template, g, request
from elasticsearch import Elasticsearch
# import os

app = Flask(__name__)


# app.config.from_object(__name__)
# app.config.update(dict(
#     DATABASE=os.path.join(app.root_path, 'code/mining.db'),
# ))
# app.config['DATABASE']

@app.route('/query', methods=['POST'])
def query():
    print(request.form)

    es = Elasticsearch([{'host': 'localhost', 'port': 9200}])

    # get result based on exact match
    res = es.search(index="ettoday", body={
        "query": {
            "match": {
                # "title": request.form['queryString']
                "body": request.form['queryString']
            }
        }
    })

    # print(res)
    data = []
    for record in res['hits']['hits']:
        print(record)
        data.append(record)
    # print(json.dumps(res, sort_keys=True, indent=4))

    return render_template('index.html', data=data, searchString=request.form['queryString'])


@app.route('/')
def homepage():
    print(app.root_path)
    return render_template('index.html')


if __name__ == '__main__':
    app.run(debug=True)  # development
