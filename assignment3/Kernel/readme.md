# Search Kernel

Simple inverted page searching kernel

Let's start with simple things. Compressions and the like are not the main priority here. 

Start with building a prototype that is working first, and then try to improve it!

[Presentation slides](https://docs.google.com/presentation/d/1YRgBzzz5Y6f5qyWeQRvchEjwKv-QTGtD_-AugWRorFE/edit?usp=sharing)

## Usage

* Insertion
```bash
curl -X POST \
-d '{"title":"title", "body":"hello world", "url":"google.com"}' \
http://localhost:8001/insert
```
* Search
```bash
curl -X POST \
-d '{"query":"ettoday"}' \
http://localhost:8001/search
```

## Package

* gse
    * segmentation, `go get -u github.com/go-ego/gse`
* gopsutil
    * system stat, `go get -u github.com/shirou/gopsutil/mem`
    * code
```go
v, _ := mem.VirtualMemory()
    megabyte := uint64(1024 * 1024)
    // almost every return value is a struct
    log.Printf("Total: %v MB, Free:%v MB, UsedPercent:%f%%\n", v.Total/megabyte, v.Free/megabyte, v.UsedPercent)
    // convert to JSON. String() is also implemented
    // log.Println(v)
```

## Things to consider

* [index](https://www.elastic.co/guide/en/elasticsearch/guide/current/inverted-index.html)  
    * case-sensitive or not
* [analyzer](https://www.elastic.co/guide/en/elasticsearch/guide/current/analysis-intro.html)
    * Use whitespace version first

## Main components

Design 好 interface! e.g. Segmentation 就要有 interface ，這樣才能輕易地換底層元件

Logging by service, don't only use general log!

* Page analyzer 
    * `page` to `page index tuples`
    * preprocessing
        * stopword
            * issue: you will have a LOT of returning results that is simply useless!
            * phrase: `the who` (a band)
            * blacklist, whitelist (phrase-level terms, higher priority than blacklist)
            * auto detection: use frequency (set threshold)
        * case sensitivity
            * issue: `she` vs `SHE`
            * simple solution: all to lower case
            * hard solution: redundant index, build both original and lowercased
                * query processing determines what to use
                * space issue
            * best solution:    
                * 帶 prefix 的 search query e.g. _SHE|she (不是很懂)
    * 斷詞 (segmentation)
        * Segmentation (cppJieba)
        * n-gram (e.g. 杜斯妥也夫斯基, 7-gram is sufficient)
        * Do-it-yourself 
            * Data structure
                * Trie, using hash/skiplist for node data
                * Hash
                    * if max = 5
                        * 中正大學資工系
                            * n = 5 中正大學資 no match
                            * n = 4 中正大學 match, skip to 資工系
                            * n = 5 ... cont.
            * Chinese
                * 長詞優先
                    * 正向
                    * 正向 + 反向
                        * conflict resolution
                            * by score
                        * e.g. 研究生物理論
                            * 正向: 研究生 物理 論
                            * 反向: 研究 生物 理論
                            * 比詞頻 / 分數
                                * 落單通常不會高分 (論)
                * 容錯：音典
                    * e.g. 畢業典禮 打成 閉業典禮
                * 辭典很重要
            * English
                * Spelling error
                    * Transposition error, e.g. that -> taht
                * Case issue
    * bigram
        * take high frquency terms as segmentation
    * scoring mechanism
        * score = term - document relation
        * options
            * documents 
                * unit occurrence score
                * position score
                    * early occurrence > late occurrence
                    * title 
                * tag score (webpage specific, e.g. `<span>`)
                * user-defined special terms
                * key term score
                * near score
                    * for terms that is not in the term
            * site
                * meta info score
                * site score
                * link score
                * format score
                    * Responsive Web Design (RWD), https, etc. add bonus score
                * behaviour score (user score)
* Posting file manager 
    * Fine grain
    * Merge page index files
        * in memory - use map (hashing)
        * on disk - (sorting)
            * output on memory limit reached
            * tuple (term, docID, socre) 
            * merge (term, total count, (docID, score), (...), ...)
    * Save run by time interval
        * search latest run first
            * search second run if first run has insufficent results
* Query processing 
    * Segmentation
    * Perform searching
    * Listen to specific port for indexing and querying
    * Highlighting

# Great resources

* [Learn and code](https://www.rosettacode.org/wiki/Inverted_index)
* [http networking for decoding JSON payload](https://gist.github.com/aodin/9493190)
* Lectures 
    * [Cornell](http://www.cs.cornell.edu/courses/cs4300/2013fa/lectures.htm)
    * [Stanford](https://nlp.stanford.edu/IR-book/html/htmledition/contents-1.html)
        * [Textbook (free)](https://nlp.stanford.edu/IR-book/)