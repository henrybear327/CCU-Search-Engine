#ifndef SEGMENTATION_H
#define SEGMENTATION_H

#include "../include/cppjieba/Jieba.hpp"

#include <string>
#include <vector>

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

    void printSegmentationResult(std::vector<std::string> &res);
    void performSegmentation(std::string &s, std::vector<std::string> &res);
};

#endif