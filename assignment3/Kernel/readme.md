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
    * `page` to `page index`
    * 斷詞
        * Segmentation (cppJieba)
    * bigram
        * take high frquency terms as segmentation
    * scoring mechanism
        * baseline version: by term frequency count
* Posting file manager
    * Fine grain
    * Merge page index files
* Query processing
    * Segmentation
    * Perform searching