all: segmentation jieba
	clang-format -style=LLVM -i *.cpp *.hpp
	astyle --style=linux *.cpp *.hpp 
	rm *.orig 
	rm *.gch

	time ./jieba
jieba: segmentation jieba.cpp 
	g++ -Wall -Wextra -Wshadow -DLOCAL -O3 -std=c++14 *.o jieba.cpp -o jieba

segmentation: segmentation.hpp segmentation.cpp 
	g++ -Wall -Wextra -Wshadow -DLOCAL -O3 -std=c++14 -c segmentation.hpp

test: test.cpp
	clang-format -style=LLVM -i test.cpp
	astyle --style=linux test.cpp 
	g++ -Wall -Wextra -Wshadow -DLOCAL -std=c++14 test.cpp 
	time ./a.out