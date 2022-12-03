#include <stdio.h>
#include <stdlib.h>
#include <string.h>
/* for apple high resolution timer functions. Yes, this is apple/Darwin specific. */
#define HAVE_MACH_TIMER
#include <mach/mach_time.h>

#define MAXLINE 1000

/* return time in nanoseconds */
static uint64_t ns() {
    static uint64_t is_init = 0;
    static mach_timebase_info_data_t info;
    if (0 == is_init) {
        mach_timebase_info(&info);
        is_init = 1;
    }
    uint64_t now;
    now = mach_absolute_time();
    now *= info.numer;
    now /= info.denom;
    return now;
}

void part1(FILE* f) {
    char linebuf[MAXLINE];
    int prio_sum = 0;
    /* variables for high res timer */
    uint64_t start, end;

    start = ns();

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
    end  = ns();
    printf("part 1, total prio: %d\n", prio_sum);
    printf("part 1 took: %llu ns\n", end - start);
}

/* a constant with 54 bits set, to initialize "what is common" mask. */
#define BITS54 0xFFFFFFFFFFFFFF

void part2(FILE* f) {
    char linebuf[MAXLINE];
    int prio_sum = 0;
    int lineno = 0; /* need to count lines to get groups of 3 */
    /* groupshared contains all items shared between the group, initialiased as "all present" */
    uint64_t groupshared = (uint64_t) BITS54;

    /* variables for high res timer */
    uint64_t start, end;

    start = ns();

    while ( fgets(linebuf, sizeof(linebuf), f) != NULL ) {
        uint64_t vec = 0;
        lineno++;

        /* all characters, process until EOL */
        for ( int i = 0; linebuf[i] != '\n'; i++ ) {
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
        /* calculate which characters are shared between this group */
        groupshared &= vec;

        /* if we processed the 3rd line, get the group badge: only item all have in common */
        if ( lineno % 3 == 0 ) {
            if ( ! groupshared ) {
                fprintf(stderr, "No chars in common at line %d\n", lineno);
            } else {
                int prio = __builtin_ctzll(groupshared);
                /* sanity check, make sure this prio is only bit set in groupshared */
                if ( (uint64_t) 1 << prio != groupshared ) {
                    fprintf(stderr, "Warning, multiple chars in common at line %d, groupshared is %llx\n", lineno, groupshared);
                }
                prio_sum += prio;
            }
            groupshared = BITS54;
        }
    }
    end  = ns();
    printf("part 2, total prio: %d\n", prio_sum);
    printf("part 2 took: %llu ns\n", end - start);
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
    part1(f);
    fseek(f, 0, SEEK_SET);
    part2(f);
    return 0;
}
