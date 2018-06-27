import json
import requests
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

    if request.method == 'GET':
        if session.get('queryString') is None:  # special case
            return redirect(url_for('homepage'))
    else:  # POST
        session['queryString'] = request.form['queryString']

    if page < 1:
        page = 1
    if page > 10:
        page = 10
    param_from = (page - 1) * 5

    # print(session)

    # get result based on exact match
    url = 'http://localhost:8001/search'
    payload = {
        "query": session['queryString'],
        "from": param_from,
        "to": param_from + 5
    }
    r = requests.post(url, data=json.dumps(payload))
    print(r.json())

    # res = es.search(index="crawler", body={
    #     "query": {
    #         "multi_match": {
    #             "query": session['queryString'],
    #             "fields": ["title", "mainText"]
    #         }
    #     },
    #     "highlight": {
    #         "number_of_fragments": 5,
    #         "fragment_size": 500,
    #         "fields": {
    #             "title": {"pre_tags": ["<mark>"], "post_tags": ["</mark>"]},
    #             "mainText": {"pre_tags": ["<mark>"], "post_tags": ["</mark>"]}
    #         }
    #     }
    # }, from_=param_from, size=5)
    # print(json.dumps(res, sort_keys=True, indent=4, ensure_ascii=False))

    # print(res)
    # data = []
    # for record in res['hits']['hits']:
    #     data.append(record)

    return render_template('index.html', data=r.json(), searchString=session['queryString'], page=page)


@app.route('/')
def homepage():
    print(app.root_path)
    return render_template('index.html')


if __name__ == '__main__':
    app.secret_key = 'super secret key'
    # app.config['SESSION_TYPE'] = 'filesystem'
    app.run(debug=True)  # development
