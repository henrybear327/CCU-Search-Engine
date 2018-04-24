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

const int BATCH_SIZE = 100000;

void testSegmentation()
{
    Segmentation segmentation;

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

void performSegmentation()
{
    IO io(inputFile, outputFolder);

    // url, title, keyword, image link, body
    vector<int> selectionColumns{1, 6, 8, 9, 16};
    vector<int> segmentationColumns{16};
    vector<string> ret;

    while (1) {
        io.getBatchRecordsInJson(BATCH_SIZE, ret, selectionColumns,
                                 segmentationColumns);
        if (ret.size() == 0)
            break;
    }

    cerr << "Total records = " << io.getRecordCount() << endl;
    cerr << "Valid records = " << io.getRecordCount() << endl;
}

int main()
{
    // fast IO
    std::ios::sync_with_stdio(false);
    cin.tie(NULL);

    // init

    // testSegmentation();

    performSegmentation();

    return EXIT_SUCCESS;
}
