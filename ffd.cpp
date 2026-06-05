#include <algorithm>
#include <iostream>
#include <vector>

using namespace std;

int ffd(int weight[], int c, int n)
{
    sort(weight, weight + n, greater<int>());

    vector<int> buckets;
    buckets.push_back(c);

    for (int i = 0; i < n; i++)
    {
        bool flag = false;

        for (int j = 0; j < buckets.size(); j++)
        {
            if (buckets[j] >= weight[i])
            {
                flag = true;
                buckets[j] -= weight[i];
                break;
            }
        }

        if (!flag)
        {
            buckets.push_back(c - weight[i]);
        }
    }

    return buckets.size();
}

int main()
{
    int weight[] = {9, 8, 2, 2, 2, 2};
    int c = 10;

    int n = sizeof(weight) / sizeof(weight[0]);

    int res = ffd(weight, c, n);

    cout << res;

    return 0;
}