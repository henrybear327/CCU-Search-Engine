#ifndef READER_H
#define READER_H

#include <string>

using namespace std;

class Reader
{
private:
    string input;

public:
    Reader(string _input) : input(_input) {}; // constructor must have {}
};

#endif