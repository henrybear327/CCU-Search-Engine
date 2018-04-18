#ifndef READER_H
#define READER_H

#include <iostream>
#include <string>

using namespace std;

class Reader
{
private:
    string inputFile;

public:
    Reader(string _inputFile) : inputFile(_inputFile)
    {
        cerr << "Init Reader" << endl;
    }; // constructor must have {}

    void testRun(int row);
};

#endif