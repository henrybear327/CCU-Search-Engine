# Todo List

- [x] Try out [cpp jieba](https://github.com/yanyiwu/cppjieba)
- [x] Try out Jsoncpp
- [x] Translate Jieba dictionary to Traditional Chinese
- [x] Preprocess the ettoday data
- [x] Pour data into ES
- [ ] Pour data into Solr

# Preface

If you can use Python, don't use c++, unless speed is what your care.

# Machine

* VMWare Workstation Ubuntu 16.04
* Intel® Core™ i5-7500 CPU @ 3.40GHz × 4
* 8GB RAM

# [Simplified Chinese to Traditional Chinese](https://github.com/BYVoid/OpenCC)

* `sudo apt install opencc`
* `opencc -i input.txt -o output.txt`

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
3. When compiling, you will encounter `fatal error: 'limonp/Logging.hpp' file not found`
	* the [solution](https://blog.csdn.net/xiaolong361/article/details/76640511) to it is copy `deps/limonp` folder to `include/cppjieba`. WTF