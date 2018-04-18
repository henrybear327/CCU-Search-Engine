COMPILER = g++
OPTIONS = -Wall -Wextra -Wshadow -O3 -std=c++14

all: format segmentation jieba
	rm *.orig *.gch
	time ./jieba

format: *.cpp *.hpp
	clang-format -style=LLVM -i *.cpp *.hpp
	astyle --style=linux *.cpp *.hpp 

jieba: segmentation jieba.cpp 
	$(COMPILER) $(OPTIONS) *.o jieba.cpp -o jieba

segmentation: segmentation.hpp segmentation.cpp 
	$(COMPILER) $(OPTIONS) -c segmentation.hpp

reader: reader.hpp reader.cpp
	$(COMPILER) $(OPTIONS) -c segmentation.hpp

test: test.cpp
	clang-format -style=LLVM -i test.cpp
	astyle --style=linux test.cpp 
	g++ -Wall -Wextra -Wshadow -std=c++14 test.cpp 
	time ./a.out