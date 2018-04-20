COMPILER = g++
OPTIONS = -Wall -Wextra -O3 -std=c++14 
# -Wshadow

all: format jieba
	rm *.orig *.gch

format: *.cpp *.hpp
	clang-format -style=LLVM -i *.cpp *.hpp
	astyle --style=linux *.cpp *.hpp 

jieba: segmentation io jieba.cpp 
	$(COMPILER) $(OPTIONS) segmentation.o io.o jieba.cpp -o jieba

segmentation: segmentation.hpp segmentation.cpp 
	$(COMPILER) $(OPTIONS) -c segmentation.hpp segmentation.cpp

io: io.hpp io.cpp
	$(COMPILER) $(OPTIONS) -c io.hpp io.cpp

run: all
	time ./jieba

test: test.cpp
	clang-format -style=LLVM -i test.cpp
	astyle --style=linux test.cpp 
	g++ -Wall -Wextra -Wshadow -std=c++14 test.cpp 
	time ./a.out