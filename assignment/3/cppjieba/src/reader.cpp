#include "reader.hpp"

#include <fstream>
#include <iostream>
#include <string>

using namespace std;

void Reader::testRun(int row)
{
    string line;
    ifstream myfile(inputFile);

    int cnt = 0;
    if (myfile.is_open()) {
        while (getline(myfile, line)) {
            if (cnt == row)
                break;
            cnt++;

            cout << line << '\n';
        }
        myfile.close();
    } else {
        cerr << "Unable to open file" << endl;
        exit(1);
    }
}