from flask import Flask, render_template, url_for, request, session, redirect
from elasticsearch import Elasticsearch

# import os

app = Flask(__name__)


# app.config.from_object(__name__)
# app.config.update(dict(
#     DATABASE=os.path.join(app.root_path, 'code/mining.db'),
# ))
# app.config['DATABASE']

@app.route('/test/<int:ha>', methods=['GET'])
def test(ha):
    return "fuck" + str(ha)


@app.route('/test', methods=['POST'])
def test1():
    return "POST FUCK"


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
            "match": {
                # "title": session['queryString']
                "body": session['queryString']
            }
        }
    }, from_=param_from, size=5)

    # print(res)
    data = []
    for record in res['hits']['hits']:
        # print(record)

        # highlight words that was search
        record['_source']['title'] = str(record['_source']['title']).replace(session['queryString'],
                                                                             "<mark>" + session[
                                                                                 'queryString'] + "</mark>")

        record['_source']['body'] = str(record['_source']['body']).replace(session['queryString'],
                                                                           "<mark>" + session[
                                                                               'queryString'] + "</mark>")
        data.append(record)
    # print(json.dumps(res, sort_keys=True, indent=4))

    return render_template('index.html', data=data, searchString=session['queryString'], page=page)


@app.route('/')
def homepage():
    print(app.root_path)
    return render_template('index.html')


if __name__ == '__main__':
    app.secret_key = 'super secret key'
    # app.config['SESSION_TYPE'] = 'filesystem'
    app.run(debug=True)  # development
