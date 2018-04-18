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
    cout << "Cut With HMM" << endl;
    cout << s << endl;
    cout << "=====================================" << endl;
    jieba.Cut(s, res, true);
}