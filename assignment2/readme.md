# TODO

Prototype first! Improvements come later!

* master
* fetcher
    - [x] get parameters filtering
    - [x] how do get restart working
    - [ ] multi-thread (subprocess)
    - [ ] how to auto switch between requests/chrome headless
    - [ ] solve weird terminating issue on Chrome headless after 4 hours
* parser
    - [x] fix redirect issue (183726490 awards issue)
    - [ ] use readability.js
    - [ ] write my own readability package
* storage
    - [ ] use MongoDB
    - [ ] indexing
    - [ ] simhash
* url pool manager (urlManager)
    - [ ] use bloom filter
    - [ ] space complexity improvement
    - [ ] check inSet before or after enqueuing?

# Packages

## pip

* selenium
* beautifulsoup4
* lxml
* pymongo
    * ubuntu
        * `sudo service mongod start`
        * `systemctl enable mongod.service`

## manual installation

* chrome
* [chrome driver](https://chromedriver.storage.googleapis.com/index.html?path=2.38/)
    * Extract
    * `chmod +x chromedriver` 
    * `sudo mv chromedriver /usr/local/bin`
   
# Good articles

* [如何入门 Python 爬虫？](https://www.zhihu.com/question/20899988)
* [DRIVING HEADLESS CHROME WITH PYTHON](https://duo.com/decipher/driving-headless-chrome-with-python)
* [Resolving a relative url path to its absolute path Ask Question](https://stackoverflow.com/questions/476511/resolving-a-relative-url-path-to-its-absolute-path?utm_medium=organic&utm_source=google_rich_qa&utm_campaign=google_rich_qa)
    * 3 hours of time wasted

# Notes

* `pkill -f chrome`
    * If you forgot to `quit`... [read me](https://stackoverflow.com/questions/15067107/difference-between-webdriver-dispose-close-and-quit?utm_medium=organic&utm_source=google_rich_qa&utm_campaign=google_rich_qa)
* backup
    * dump `mongodump -h 127.0.0.1 -d npr_test -o ./data.dump`
    * store `mongorestore -h 127.0.0.1 -d npr_test ./data.dump/npr_test/data.bson`
    
# Sites

* English
    * News
        * [NPR](https://www.npr.org/)
            * Static
        * CNN
            * Dynamic
        * [Hacker news](https://news.ycombinator.com/)
            * Static / Dynamic
* Chinese
    * News
        * Ettoday