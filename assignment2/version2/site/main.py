import json
import pymongo
import pprint
from pymongo import MongoClient
from elasticsearch import Elasticsearch
from flask import Flask, render_template, url_for, request, session, redirect

app = Flask(__name__)

app.config.from_object(__name__)
app.config.update(dict(
    DATABASE="crawler",
))


@app.route('/query', methods=['POST'])
@app.route('/query/<int:page>', methods=['GET'])
def query(page=1):
    # print(request.form)

    es = Elasticsearch([{'host': 'localhost', 'port': 9200}])

    if request.method == 'GET':
        if session.get('queryString') is None:  # special case
            return redirect(url_for('homepage'))
    else:  # POST
        session['queryString'] = request.form['queryString']

    if page < 1:
        page = 1
    if page > 10:
        page = 10
    param_from = (page - 1) * 10

    # print(session)

    # get result based on exact match
    res = es.search(index="ettoday", body={
        "query": {
            "multi_match": {
                "query": session['queryString'],
                "fields": ["title", "body"]
            }
        },
        "highlight": {
            "number_of_fragments": 5,
            "fragment_size": 200,
            "fields": {
                "title": {"pre_tags": ["<mark>"], "post_tags": ["</mark>"]},
                "body": {"pre_tags": ["<mark>"], "post_tags": ["</mark>"]}
            }
        }
    }, from_=param_from, size=5)
    print(json.dumps(res, sort_keys=True, indent=4, ensure_ascii=False))

    # print(res)
    data = []
    for record in res['hits']['hits']:
        data.append(record)

    return render_template('index.html', data=data, searchString=session['queryString'], page=page)


@app.route('/')
def homepage():
    print(app.root_path)
    return render_template('index.html')


@app.route('/report/<string:top_level_domain>/<int:start>/<int:end>')
def mongo_db_query(top_level_domain: str, start: int, end: int):
    client = MongoClient('localhost', 27017)
    db = client[app.config['DATABASE']]
    collection = db["sitePage"]

    res = collection.find({"tld": top_level_domain}).sort('fetchTime', pymongo.DESCENDING).skip(start).limit(
        end - start)  # [start, end)
    data = []
    for post in res:
        # pprint.pprint(post)
        data.append(post)

    # 100-199
    prev_left = start - 100
    prev_right = end - 100
    next_left = start + 100
    next_right = end + 100
    if prev_left < 0:
        prev_left = 0
        prev_right = 100

    return render_template('TLDreport.html', data=data, tld=top_level_domain, prev_left=prev_left, prev_right=prev_right
                           ,next_left=next_left, next_right=next_right)


@app.route('/report')
def report():
    client = MongoClient('localhost', 27017)
    db = client[app.config['DATABASE']]
    collection = db["sitePage"]

    tld = []
    for rec in collection.distinct("tld"):
        # pprint.pprint(rec)
        tld.append(rec)

    cnt = []
    last_fetched = []
    for rec in tld:
        ret = collection.find({"tld": rec})
        cnt.append(ret.count())
        last = ret.sort("fetchTime", pymongo.DESCENDING).limit(1)
        for tmp in last:
            last_fetched.append(tmp['fetchTime'])
            break

    return render_template('TLDListing.html', tld=tld, cnt=cnt, last_fetched=last_fetched)


if __name__ == '__main__':
    app.secret_key = 'super secret key'
    # app.config['SESSION_TYPE'] = 'filesystem'
    app.run(debug=True)  # development
