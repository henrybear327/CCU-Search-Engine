#include <bits/stdc++.h>

using namespace std;

#include "json.hpp"
using json = nlohmann::json;

int main()
{
    json j;
    j["body"] = "This is a test";

    cout << j.dump() << endl;
    return 0;
}