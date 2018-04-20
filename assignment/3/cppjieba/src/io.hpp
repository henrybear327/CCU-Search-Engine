#ifndef READER_H
#define READER_H

#include "segmentation.hpp"
#include <fstream> // field ‘myfile’ has incomplete type ‘std::ifstream {aka std::basic_ifstream<char>}’
#include <iostream>
#include <string>
#include <vector>

using namespace std;

static const int HEADINGCOUNT = 19;

// https://stackoverflow.com/questions/2328671/constant-variables-not-working-in-header
extern const char *recordHeading[HEADINGCOUNT];

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

class IO
{
private:
    Segmentation segmentation;

    string inputFile;
    string outputFolder;

    string line;
    ifstream myfile;

    int filenameCounter;
    int recordCounter;
    int validRecordCounter;

    void writeToFile(vector<string> &batchData);

    // get raw record
    Record getRecord();
    // convert struct Record to json string
    string getRecordInJson(const Record &rec, vector<int> &selectionColumns);

public:
    IO(string _inputFile, string _outputFolder) : myfile(_inputFile)
    {
        cerr << "Init Reader" << endl;

        inputFile = _inputFile;
        outputFolder = _outputFolder;

        if (myfile.is_open()) {
            getline(myfile, line);
            if (line.length() == 1 && line[0] == '@')
                cerr << "valid starting" << endl;
            else {
                cerr << "invalid data" << endl;
                exit(1);
            }
        }

        filenameCounter = 0;
        validRecordCounter = 0;
    }; // constructor must have {}

    void testRun(int row);

    void debugPrintRecord(const Record &rec, vector<int> &selection);

    void getBatchRecordsInJson(int batchSize, vector<string> &ret,
                               vector<int> &selectionColumns,
                               vector<int> &segmentationColumns);

    int getRecordCount();
    int getValidRecordCount();

    ~IO()
    {
        myfile.close();
    }
};

#endif