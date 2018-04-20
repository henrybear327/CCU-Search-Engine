#include "segmentation.hpp"

void Segmentation::printSegmentationResult(vector<string> &res)
{
    for (auto i : res) {
        cout << i << " ";
    }
    cout << endl;
}

void Segmentation::performSegmentation(string &s, vector<string> &res)
{
    res.clear();

    // cerr << "Cut With HMM" << endl;
    // cerr << s << endl;
    // cerr << "=====================================" << endl;
    jieba.Cut(s, res, false);
}

string Segmentation::getSegmentationString(vector<string> &res)
{
    string ans = "";

    for (auto i : res)
        ans += i + " ";

    return ans;
}