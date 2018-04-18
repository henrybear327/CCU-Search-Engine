#include "reader.hpp"

#include <algorithm>
#include <iostream>
#include <string>

using namespace std;

void Reader::testRun(int row)
{
    // init their own reading variable
    string testline;
    ifstream testfile(inputFile);

    int cnt = 0;
    if (testfile.is_open()) {
        while (getline(testfile, testline)) {
            if (cnt == row)
                break;
            cnt++;

            cout << testline << '\n';

            // if (testline[0] == '@')
            //     cout << "WOW!" << endl;
        }
        testfile.close();
    } else {
        cerr << "Unable to open file" << endl;
        exit(1);
    }
}

Record Reader::getRecord()
{
    Record rec;

    if (myfile.is_open()) {
        int headerIndex = -1;
        string data = "";
        while (getline(myfile, line)) {
            rec.hasData = true;

            // cerr << line << '\n';

            if (line[0] == '@') {
                data = "";

                // cerr << "Header line" << endl;
                if (line.length() == 1) {
                    // cerr << "start of a record" << endl;
                    break;
                }

                int ending = -1;
                for (int i = 1; i < (int)line.length(); i++) {
                    if (i > 20)
                        break;
                    if (line[i] == ':') {
                        ending = i;
                        break;
                    }
                }
                if (ending == -1) {
                    // cerr << "no : found for @" << endl;
                    // cerr << line << endl;
                    // exit(1);

                    // GG: this might just be another normal line starting with @
                    // cat ettoday| grep -A 3 -B 3 "^@@" | more
                    rec.data[headerIndex] += "\n" + line;
                    continue;
                }

                string header = line.substr(1, ending - 1);
                // cerr << "header = " << header << endl;

                headerIndex = -1;
                for (int i = 0; i < HEADINGCOUNT; i++) {
                    if (header == recordHeading[i]) {
                        headerIndex = i;
                        break;
                    }
                }
                if (headerIndex == -1) {
                    cerr << "header not found!" << endl;
                    cerr << "got " << header << endl;
                    exit(1);
                }

                rec.data[headerIndex] = line.substr(ending + 1);
            } else {
                rec.data[headerIndex] += "\n" + line;
            }
        }
    } else {
        cerr << "Unable to open file" << endl;
        exit(1);
    }

    return rec;
}

void Reader::printRecord(const Record &rec, vector<int> &selection)
{
    cout << "===========================================" << endl;
    for (int i = 0; i < HEADINGCOUNT; i++) {
        if (find(selection.begin(), selection.end(), i) != selection.end()) {
            cout << recordHeading[i] << ": ";
            cout << rec.data[i] << endl;
        }
    }
    cout << "===========================================" << endl;
}
