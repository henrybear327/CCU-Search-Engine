# Search Kernel

Simple inverted page searching kernel

# TODO

Let's start with simple things. Compressions and the like are not the main priority here. 

Start with building a prototype that is working first, and then try to improve it!

## Things to consider

* [index](https://www.elastic.co/guide/en/elasticsearch/guide/current/inverted-index.html)  
    * case-sensitive or not
* [analyzer](https://www.elastic.co/guide/en/elasticsearch/guide/current/analysis-intro.html)
    * Use whitespace version first

## Main components

Design 好 interface! e.g. Segmentation 就要有 interface ，這樣才能輕易地換底層元件

Logging by service, don't only use general log!

* Page analyzer (parser)
    * `page` to `page index tuples`
    * 斷詞
        * Segmentation (cppJieba)
        * n-gram
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
            * site
                * meta info score
                * site score
                * link score
                * format score
                    * Responsive Web Design (RWD), https, etc. add bonus score
                * behaviour score (user score)
* Posting file manager (indexer)
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
* Query processing (interface)
    * Segmentation
    * Perform searching
    * Listen to specific port for indexing and querying