all:
	clang-format -style=LLVM -i *.cpp 
	astyle --style=linux *.cpp 
	clang++ -Wall -Wextra -Wshadow -DLOCAL -std=c++14 test.cpp 
	time ./a.out