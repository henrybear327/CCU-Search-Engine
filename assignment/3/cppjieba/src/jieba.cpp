#include <iostream>
#include <string>

#include "reader.hpp"
#include "segmentation.hpp"

using namespace std;

string inputFile = "../../data/ettoday";

const int PRINT_KEEP_ALIVE = 10000;

int main()
{
    // fast IO
    std::ios::sync_with_stdio(false);
    cin.tie(NULL);

    // init
    Segmentation segmentation;
    Reader reader(inputFile);

    // {
    //     // test input
    //     // https://udn.com/news/story/7266/3094115
    //     string s =
    //     "台灣虎航麻疹群聚疫情延燒，擴及其他航空公司。衛福部疾管署今公布"
    //                "，此群聚情"
    //                "再新增4名感染者，其中包含2名馬來西亞航空與長榮航空的地勤人員，"
    //                "以及2名台"
    //                "虎女性空服員，推估感染個案皆與4月初發病的台虎空服員於可傳染期在"
    //                "相同或鄰"
    //                "近的備勤室工作，累計此群聚疫情已達12人，接觸者共達2978人。";
    //     vector<string> res;

    //     segmentation.performSegmentation(s, res);
    //     segmentation.printSegmentationResult(res);
    // }

    // reader.testRun(100);
    // {
    //     vector<int> selection{6, 16};
    //     for (int i = 0; i < 3; i++) {
    //         auto rec = reader.getRecord();
    //         reader.printRecord(rec, selection);
    //     }
    // }

    {
        vector<int> selection{6, 16};
        int cnt = 0;
        while (1) {
            auto rec = reader.getRecord();
            if (rec.hasData == false)
                break;
            cnt++;
            if (cnt % 10000 == 0)
                cerr << "Data cnt " << cnt << endl;
            reader.printRecord(rec, selection);
        }
        cerr << "Done! " << cnt << " records" << endl;
    }

    return EXIT_SUCCESS;
}
