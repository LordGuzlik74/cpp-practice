#include <iostream>
using namespace std;

int main(){
    int n;
    cin>>n;
    int *arr=new int[n];
    for (int i=0;i<n;i++){
        int w;
        cin>>w;
        arr[i]=w;
    }
    int k;
    cin>>k;
    k=k%n;
    int *new_arr=new int[n];
    for (int i=0;i<n;i++){
        new_arr[(i + k) % n] = arr[i];
    }
    for (int i=0;i<n;i++){
        cout<<new_arr[i]<<" ";
    }
    cout<<"\n";
    delete[] arr;
    delete[] new_arr;
    return 0;
}