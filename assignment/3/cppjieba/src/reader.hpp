#ifndef READER_H
#define READER_H

#include <fstream> // field ‘myfile’ has incomplete type ‘std::ifstream {aka std::basic_ifstream<char>}’
#include <iostream>
#include <string>
#include <vector>

using namespace std;

static const int HEADINGCOUNT = 19;

// https://stackoverflow.com/questions/2328671/constant-variables-not-working-in-header
static const char *recordHeading[HEADINGCOUNT] = {
    "Gais_REC",  "url",       "MainTextMD5", "UntagMD5", "SiteCode",
    "UrlCode",   "title",     "Size",        "keyword",  "image_links",
    "Fetchtime", "post_time", "Ref",         "BodyMD5",  "Lang",
    "IP",        "body",      "botVer",      "Time"
};

struct Record {
    bool hasData;
    string data[HEADINGCOUNT]; // heading type, data
    Record()
    {
        for (int i = 0; i < HEADINGCOUNT; i++)
            data[i] = "";
        hasData = false;
    }
};

class Reader
{
private:
    string inputFile;
    string line;
    ifstream myfile;

public:
    Reader(string _inputFile) : inputFile(_inputFile), myfile(_inputFile)
    {
        cerr << "Init Reader" << endl;
        if (myfile.is_open()) {
            getline(myfile, line);
            if (line.length() == 1 && line[0] == '@')
                cerr << "valid starting" << endl;
            else {
                cerr << "invalid data" << endl;
                exit(1);
            }
        }
    }; // constructor must have {}

    void testRun(int row);
    Record getRecord();
    void printRecord(const Record &rec, vector<int> &selection);

    ~Reader()
    {
        myfile.close();
    }
};

#endif