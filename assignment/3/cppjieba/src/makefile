jieba: jieba.cpp
	clang-format -style=LLVM -i jieba.cpp
	astyle --style=linux jieba.cpp
	clang++ -Wall -Wextra -Wshadow -DLOCAL -O3 -std=c++14 jieba.cpp
	time ./a.out
test: test.cpp
	clang-format -style=LLVM -i test.cpp
	astyle --style=linux test.cpp 
	clang++ -Wall -Wextra -Wshadow -DLOCAL -std=c++14 test.cpp 
	time ./a.out