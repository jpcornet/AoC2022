#include <stdio.h>
#include <stdlib.h>
#include <string.h>
/* for apple high resolution timer functions*/
#define HAVE_MACH_TIMER
#include <mach/mach_time.h>

#define MAXLINE 1000

int process_file(FILE* f) {
    char linebuf[MAXLINE];
    int prio_sum = 0;
    /* variable for high res timer */
    static uint64_t is_init = 0;
    static mach_timebase_info_data_t info;
    uint64_t start, now;

    if (0 == is_init) {
      mach_timebase_info(&info);
      is_init = 1;
    }
    start  = mach_absolute_time();
    start *= info.numer;
    start /= info.denom;

    while ( fgets(linebuf, sizeof(linebuf), f) != NULL ) {
        uint64_t vec = 0; /* 64 bit to keep track of characters */
        int prio = 0;
        int len;

        /* make sure buffer was large enough */
        len = strlen(linebuf);
        if ( len == 0 ) {
            /* we read an empty line, try next line. */
            continue;
        }
        if ( linebuf[--len] != '\n' ) {
            fprintf(stderr, "Overlong line [%s], abort\n", linebuf);
            exit(-1);
        }
        /* effectively remove \n from line */
        linebuf[len] = '\0';
        /* first half of line/rucksack, remember characters used. */
        len /= 2;
        for ( int i = 0; i < len; i++ ) {
            /* if character is uppercase, init prio at 26, 0 for lowercase. */
            char c = linebuf[i];
            int char_prio = 26 * (c <= 'Z');
            c |= 0x20; /* make sure c is lower case now */
            c -= 'a'; /* c is now between 0 and 25 */
            if ( c < 0 || c > 25 ) {
                fprintf(stderr, "Invalid character %c in line [%s], abort\n", linebuf[i], linebuf);
                exit(-1);
            }
            char_prio += c + 1;
            vec |= (uint64_t) 1 << char_prio;
        }
        /* second half of line, scan for char already used */
        for ( int i = len; linebuf[i]; i++ ) {
            /* if character is uppercase, init prio at 26, 0 for lowercase. */
            char c = linebuf[i];
            int char_prio = 26 * (c <= 'Z');
            c |= 0x20; /* make sure c is lower case now */
            c -= 'a'; /* c is now between 0 and 25 */
            if ( c < 0 || c > 25 ) {
                fprintf(stderr, "Invalid character '%c' in line [%s], abort\n", linebuf[i], linebuf);
                exit(-1);
            }
            char_prio += c + 1;
            if ( vec & (uint64_t) 1 << char_prio ) {
                if ( prio != 0 && prio != char_prio ) {
                    fprintf(stderr, "Hmm, more than 1 repeating character '%c' in string halves [%s]\n", linebuf[i], linebuf);
                }
                prio = char_prio;
            }
        }
        prio_sum += prio;
    }
    now  = mach_absolute_time();
    now *= info.numer;
    now /= info.denom;
    printf("part 1, total prio: %d\n", prio_sum);
    printf("part 1 took: %llu ns\n", now - start);
    return 0;
}

int main(int argc, char** argv) {
    FILE* f;
    if ( argc != 2 ) {
        fprintf(stderr, "Provide input file\n");
        exit(1);
    }
    if ( (f = fopen(argv[1], "r")) == NULL ) {
        perror("Error opening input file");
        exit(2);
    }
    return process_file(f);
}
