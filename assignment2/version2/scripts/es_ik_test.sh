#!/bin/bash

curl -XPUT http://localhost:9200/ik

echo "\n\n"

curl -XPOST http://localhost:9200/ik/fulltext/_mapping -H 'Content-Type:application/json' -d'
{
        "properties": {
            "content": {
                "type": "text",
                "analyzer": "ik_max_word",
                "search_analyzer": "ik_max_word"
            }
        }

}'

echo "\n\n"

curl -XPOST http://localhost:9200/ik/fulltext/1 -H 'Content-Type:application/json' -d'
{"content":"美国留给伊拉克的是个烂摊子吗"}
'
echo "\n\n"

curl -XPOST http://localhost:9200/ik/fulltext/2 -H 'Content-Type:application/json' -d'
{"content":"公安部：各地校车将享最高路权"}
'
echo "\n\n"

curl -XPOST http://localhost:9200/ik/fulltext/3 -H 'Content-Type:application/json' -d'
{"content":"中韩渔警冲突调查：韩警平均每天扣1艘中国渔船"}
'
echo "\n\n"

curl -XPOST http://localhost:9200/ik/fulltext/4 -H 'Content-Type:application/json' -d'
{"content":"中國驻洛杉矶领事馆遭亚裔男子枪击 嫌犯已自首"}
'
echo "\n\n"

curl -XPOST http://localhost:9200/ik/fulltext/_search?pretty  -H 'Content-Type:application/json' -d'
{
    "query" : { "match" : { "content" : "中國" }},
    "highlight" : {
        "pre_tags" : ["<tag1>", "<tag2>"],
        "post_tags" : ["</tag1>", "</tag2>"],
        "fields" : {
            "content" : {}
        }
    }
}
'
echo "\n\n"

# curl -XDELETE 'localhost:9200/ik?pretty'