#include <climits>
#include <iostream>
#include <string>

#include "io.hpp"
#include "segmentation.hpp"

using namespace std;

#include "json.hpp"
using json = nlohmann::json;

string inputFile = "../../data/ettoday";
string outputFolder = "../../data/jieba/";

const int PRINT_KEEP_ALIVE = 10000;
const int BATCH_SIZE = 100000;
int filenameCounter;

void testSegmentation(Segmentation &segmentation)
{
    // test input
    // https://udn.com/news/story/7266/3094115
    string s = "台灣虎航麻疹群聚疫情延燒，擴及其他航空公司。衛福部疾管署今公布"
               "，此群聚情"
               "再新增4名感染者，其中包含2名馬來西亞航空與長榮航空的地勤人員，"
               "以及2名台"
               "虎女性空服員，推估感染個案皆與4月初發病的台虎空服員於可傳染期在"
               "相同或鄰"
               "近的備勤室工作，累計此群聚疫情已達12人，接觸者共達2978人。";
    vector<string> res;

    segmentation.performSegmentation(s, res);
    segmentation.printSegmentationResult(res);
}

void printJson(vector<string> &batchData)
{
    if (batchData.size() == 0)
        return;

    string filename = outputFolder + "ettoday_" + to_string(batchData.size()) +
                      "_" + to_string(filenameCounter++) + ".txt";

    ofstream myfile;
    myfile.open(filename, ios::trunc);

    try {
        for (auto i : batchData) {
            // myfile << j.dump(4) << endl;
            myfile << i << endl;
        }
    } catch (nlohmann::detail::type_error) {
        cerr << "json error while dumping (batch data loss)" << endl;
    }

    batchData.clear();
    myfile.close();
}

void dealJson(vector<string> &batchData, IO &io, Record &rec,
              vector<int> &selection)
{
    if (batchData.size() == BATCH_SIZE) {
        printJson(batchData);
    }
    auto jsonString = io.getRecordInJson(rec, selection);
    if (jsonString != "")
        batchData.push_back(jsonString);
}

void performSegmentation(Segmentation &segmentation, IO &io, int n = INT_MAX)
{
    // url, title, keyword, image link, body
    vector<int> selection{1, 6, 8, 9, 16};
    int cnt = 0;
    vector<string> batchData;
    for (; cnt < n; cnt++) {
        auto rec = io.getRecord();
        if (rec.hasData == false) // end of file
            break;

        if (cnt % 10000 == 0)
            cerr << "Data cnt " << cnt << endl;

        // segmentation on body only
        vector<string> res;
        segmentation.performSegmentation(rec.data[16], res);
        rec.data[16] = segmentation.getSegmentationString(res);

        // reader.debugPrintRecord(rec, selection);
        dealJson(batchData, io, rec, selection);
    }

    printJson(batchData);
    cerr << "Done! " << cnt << " records" << endl;
}

int main()
{
    // fast IO
    std::ios::sync_with_stdio(false);
    cin.tie(NULL);

    // init
    Segmentation segmentation;
    IO reader(inputFile);
    filenameCounter = 0;

    testSegmentation(segmentation);
    // testReader(segmentation, reader, 3);

    // performSegmentation(segmentation, reader, 10);
    // performSegmentation(segmentation, reader);

    return EXIT_SUCCESS;
}
