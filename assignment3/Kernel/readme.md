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

* Page analyzer
    * `page` to `page index tuples`
    * 斷詞
        * Segmentation (cppJieba)
    * bigram
        * take high frquency terms as segmentation
    * scoring mechanism
        * score = term - document relation
        * options
            * documents 
                * unit occurrence score
                * position score
                    * early occurrence > late occurrence
                * tag score (webpage specific, e.g. `<span>`)
                * user-defined special terms
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
* Master slave
    * Two slaves
        * Query slave
        * Indexer slave