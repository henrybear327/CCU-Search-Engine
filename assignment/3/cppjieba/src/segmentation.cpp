#include "segmentation.hpp"

void Segmentation::printSegmentationResult(std::vector<std::string> &res)
{
    for (auto i : res) {
        std::cout << i << " ";
    }
    std::cout << std::endl;
}

void Segmentation::performSegmentation(std::string &s,
                                       std::vector<std::string> &res)
{
    std::cout << "Cut With HMM" << std::endl;
    std::cout << s << std::endl;
    std::cout << "=====================================" << std::endl;
    jieba.Cut(s, res, true);
}