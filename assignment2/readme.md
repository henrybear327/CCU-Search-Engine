# TODO

Prototype first! Improvements come later!

* master
* fetcher
    - [x] get parameters filtering
    - [x] how do get restart working
    - [ ] multi-thread (subprocess)
* parser
    - [ ] use readability.js
    - [ ] write my own readability package
* storage
    - [ ] use SQLite
* url pool manager (urlManager)
    - [ ] use bloom filter
    - [ ] space complexity improvement
    - [ ] check inSet before or after enqueuing?

# Packages

## pip

* selenium
* chrome
* beautifulsoup4
* lxml
* pymongo

## manual installation

* [chrome driver](http://chromedriver.storage.googleapis.com/2.38/chromedriver_linux64.zip)
    * Extract
    * `chmod +x chromedriver` 
    * `sudo mv chromedriver /usr/local/bin`
   
# Good articles

* [如何入门 Python 爬虫？](https://www.zhihu.com/question/20899988)
* [DRIVING HEADLESS CHROME WITH PYTHON](https://duo.com/decipher/driving-headless-chrome-with-python)

# Notes

* `pkill -f chrome`
    * If you forgot to `quit`... [read me](https://stackoverflow.com/questions/15067107/difference-between-webdriver-dispose-close-and-quit?utm_medium=organic&utm_source=google_rich_qa&utm_campaign=google_rich_qa)