#ifndef SEGMENTATION_H
#define SEGMENTATION_H

#include "../include/cppjieba/Jieba.hpp"

#include <iostream>
#include <string>
#include <vector>

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
        : jieba(DICT_PATH, HMM_PATH, USER_DICT_PATH, IDF_PATH, STOP_WORD_PATH)
    {
        cerr << "Init Segmentation" << endl;
    }

    void printSegmentationResult(vector<string> &res);
    void performSegmentation(string &s, vector<string> &res);
    string getSegmentationString(vector<string> &res);
};

#endif