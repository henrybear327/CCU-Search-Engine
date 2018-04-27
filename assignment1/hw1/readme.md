# 網際網路資料檢索 作業一

## Specification

利用 Elastic Search 作為 search engine 的 backend，做一個搜尋引擎。

## Elastic Search 

### 概覽

Elastic Search 是用 Java 編寫而成，需要 Java 8 或以上的板本才能運行。

對於資料的查詢，elastic search 有 restful API 可以用，資料都是使用 JSON 來進行傳遞。

對於python而言，可以使用官方的 elasticsearch 套件，或是使用社群開發的 elasticsearch-dsl。

網路上有人拿傳統 SQL database 來打比方，簡單的說明的 elastic search 的架構上的用詞與我們已知的概念大概是呈現什麼樣的關係。

| Elastic Search    | Relational Database |
| ------------------- |-----------------|
| node                | database server | 
| index               | database        |
| type                | table           | 
| document            | row             | 
| field               | column          | 

### Restful API

預設安裝完並啟動的 elastic search 會使用 port 9200 來對外溝通。 利用 `curl` 打 request 到此可以獲得 json 資料。但是因為作業用倒的都是透過 python ，所以這部分就不深究了。

### elasticsearch-dsl

這是社群將 elastic search 官方提供的 low-level API 進行包裝後所發布的套件。

#### Basic search example

對於基本的搜尋，可以簡單的利用API呈現出來。

```python
from elasticsearch import Elasticsearch


def go():
    es = Elasticsearch([{'host': 'localhost', 'port': 9200}])

    # get result based on exact match
    res = es.search(index="ettoday", body={
        "query": {
            "match": {
                "body": "中正大學"
            }
        }
    })
    for record in res['hits']['hits']:
        print(record)


if __name__ == "__main__":
    go()
```

#### Multiple field search and highlighting

這兩個功能真的很方便，尤其highlighting的部分！因為不用自己做excerpt又可以拿到highlight!

```python
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
```

## Flask

網頁的部分，我完全沒有使用到 javascript，而是全部都在 flask + template engine 處理完，直接呈現。

### backend code

```python
from flask import Flask, render_template, url_for, request, session, redirect
from elasticsearch import Elasticsearch

import os, json

app = Flask(__name__)

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

    data = []
    for record in res['hits']['hits']:
        data.append(record)

    return render_template('index.html', data=data, searchString=session['queryString'], page=page)


@app.route('/')
def homepage():
    print(app.root_path)
    return render_template('index.html')


if __name__ == '__main__':
    app.secret_key = 'super secret key'
    # app.config['SESSION_TYPE'] = 'filesystem'
    app.run(debug=True)  # development
```

### frontend (template engine) code

```html
{% extends "layout.html" %}

{% block head %}

{% endblock %}

{% block body %}

    <div class="container">
        <div class="row">
            <!--https://getbootstrap.com/docs/4.0/utilities/spacing/-->
            <div class="col-md-8 mx-auto p-md-8 my-5">
                <div class="card">
                    <div class="card-body">
                        <form action="{{ url_for('query') }}" method="post">
                            <div class="form-group">
                                <div class="input-group">
                                    <input type="text" name="queryString" class="form-control" placeholder="Search!"
                                           aria-label="search" value="{{ searchString }}">
                                    <div class="input-group-append">
                                        <button class="btn btn-outline-primary" type="submit"><i
                                                class="fas fa-search"></i>
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </form>
                    </div>
                </div>

                {% if data is defined %}
                    <hr>
                    {% if searchString|length == 0 %}

                        <h5 class="text-center">
                            Please enter some text to search...
                        </h5>

                    {% elif data|length == 0 %}

                        <h5 class="text-center">
                            No result about {{ searchString }} is found...
                        </h5>

                    {% else %}

                        {% for record in data %}
                            <div class="card">
                                <h5 class="card-header">
                                    <span class="badge badge-primary">{{ loop.index }}</span>

                                    {% if record['highlight']['title'] is defined %}
                                        <a href="{{ record['_source']['url'] }}" target="_blank">
                                            {{ record['highlight']['title'][0]|safe }}
                                        </a>
                                    {% else %}
                                        <a href="{{ record['_source']['url'] }}" target="_blank">
                                            {{ record['_source']['title']|safe }}
                                        </a>
                                    {% endif %}
                                </h5>
                                <div class="card-body">
                                    <h5 class="card-title"></h5>
                                    <p class="card-text">
                                        {% if record['highlight']['body'] is defined %}
                                            {{ record['highlight']['body'][0]|safe }}
                                        {% else %}
                                            {{ record['_source']['body']|safe }}
                                        {% endif %}
                                    </p>
                                </div>
                            </div>
                            <br>
                        {% endfor %}

                        <nav aria-label="searchPagination">
                            <ul class="pagination">
                                <li class="page-item {% if page == 1 %} disabled {% endif %}">
                                    <a class="page-link" href="{{ url_for('query') }}/{{ page - 1 }}">Previous</a>
                                </li>

                                {% for cnt in range(1, 11, 1) %}
                                    <li class="page-item {% if cnt == page %} active {% endif %}">
                                        <a class="page-link" href="{{ url_for('query') }}/{{ cnt }}">
                                            {{ cnt }} {% if cnt == page %}
                                            <span class="sr-only">(current)</span> {% endif %}
                                        </a>
                                    </li>
                                {% endfor %}

                                <li class="page-item {% if page == 10 %} disabled {% endif %}">
                                    <a class="page-link" href="{{ url_for('query') }}/{{ page + 1 }}">Next</a>
                                </li>
                            </ul>
                        </nav>

                    {% endif %}

                {% endif %}

            </div>
        </div>
    </div>
{% endblock %}
```