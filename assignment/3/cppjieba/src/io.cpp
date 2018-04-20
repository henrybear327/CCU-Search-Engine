#include "io.hpp"

#include <algorithm>
#include <iostream>
#include <string>

#include "json.hpp"
using json = nlohmann::json;

using namespace std;

void IO::testRun(int row)
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

Record IO::getRecord()
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

void IO::debugPrintRecord(const Record &rec, vector<int> &selection)
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

string IO::getRecordInJson(const Record &rec, vector<int> &selection)
{
    // to JSON
    json j;
    for (auto i : selection) {
        try {
            j[recordHeading[i]] = rec.data[i];
        } catch (nlohmann::detail::type_error) {
            cerr << "json error while building" << endl;
            cerr << rec.data[i] << endl;
            return "";
        }
    }

    // print
    string ret = "";
    try {
        ret = j.dump();
        // cerr << ret << endl;
    } catch (nlohmann::detail::type_error) {
        cerr << "json error while printing (ignore this data)" << endl;
    }
    return ret;
}