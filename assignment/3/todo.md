# Todo List

- [ ] Try out cpp jieba 
- [ ] Preprocess the ettoday data
- [ ] Pour data into Solr
- [ ] Pour data into ES

# Machine

* VMWare Workstation Ubuntu 16.04
* Intel® Core™ i5-7500 CPU @ 3.40GHz × 4
* 8GB RAM

# [Simplified Chinese to Traditional Chinese](http://linux-wiki.cn/wiki/zh-tw/%E7%AE%80%E7%B9%81%E8%BD%AC%E6%8D%A2)

* Install `zh-autoconvert`
* Run `iconv -f utf8 -t gbk sim.txt | autob5 | iconv -f big5 -t utf8 > fan.txt`

# Notes on cppjieba

1. Build time is around 35 seconds
```sh
git clone --depth=10 --branch=master git://github.com/yanyiwu/cppjieba.git
cd cppjieba
mkdir build
cd build
cmake ..
make
```
2. Put your source code under `src`
3. When compiling, you will encounter `fatal error: 'limonp/Logging.hpp' file not found`, the solution to it is copy `deps/limonp` folder to `include/cppjieba`. WTF