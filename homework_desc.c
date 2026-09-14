#include <stdio.h>



// int main(void){
//     int c=0;
//     for (int i=200;i<7000;i++){
//         if(i%33==0 && i%7!=0){
//             c++;
//         }
//     }
//     printf("%d\n", c);

//     return 0;
// }
int main(void){
    int c=1;
    for (int i=2;i<125;i++){
        if(125%i==0 && 765%i==0){
            c=c*i;
            //printf("%d\n",i);
        }
    }
    printf("%d\n",c);

    return 0;
}