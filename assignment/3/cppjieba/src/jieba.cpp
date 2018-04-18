#include "../include/cppjieba/Jieba.hpp"

using namespace std;

const char *const DICT_PATH = "../dict/jieba.dict.fan.utf8";
const char *const HMM_PATH = "../dict/hmm_model.fan.utf8";
const char *const USER_DICT_PATH = "../dict/user.dict.fan.utf8";
const char *const IDF_PATH = "../dict/idf.fan.utf8";
const char *const STOP_WORD_PATH = "../dict/stop_words.fan.utf8";

class Segmentation
{
private:
    cppjieba::Jieba jieba;

public:
    Segmentation()
        : jieba(DICT_PATH, HMM_PATH, USER_DICT_PATH, IDF_PATH, STOP_WORD_PATH) {}
    void printSegmentationResult(vector<string> &res)
    {
        for (auto i : res) {
            cout << i << " ";
        }
        cout << endl;
    }

    void performSegmentation(string &s, vector<string> &res)
    {
        cout << "Cut With HMM" << endl;
        cout << s << endl;
        cout << "=====================================" << endl;
        jieba.Cut(s, res, true);
    }
};

int main(int argc, char **argv)
{
    // fast IO
    std::ios::sync_with_stdio(false);
    cin.tie(NULL);

    // init
    Segmentation segmentation;

    string s =
        "台灣虎航麻疹群聚疫情延燒，擴及其他航空公司。衛福部疾管署今公布，此群聚情"
        "再新增4名感染者，其中包含2名馬來西亞航空與長榮航空的地勤人員，以及2名台"
        "虎女性空服員，推估感染個案皆與4月初發病的台虎空服員於可傳染期在相同或鄰"
        "近的備勤室工作，累計此群聚疫情已達12人，接觸者共達2978人。";
    vector<string> res;

    segmentation.performSegmentation(s, res);
    segmentation.printSegmentationResult(res);

    return EXIT_SUCCESS;
}
