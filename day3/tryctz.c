// C program to illustrate __builtin_clz(x)
#include <stdio.h>
#include <stdlib.h>

int main()
{
    uint64_t n = (uint64_t) 1 << 50;

    printf("Count of trailing zeros after 1 in 0x%llx is %d", n, __builtin_ctzll(n));
    return 0;
}
