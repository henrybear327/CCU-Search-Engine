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
    cerr << "Cut With HMM" << endl;
    cerr << s << endl;
    cerr << "=====================================" << endl;
    jieba.Cut(s, res, true);
}