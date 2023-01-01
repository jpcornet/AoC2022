# AoC2022
🎄 AdventOfCode 2022 🎄

My first attempts to participate in AoC. Be gentle. Trying as many different languages as possible, one for each day.

Given timings are rough and are from my 2020 intel macbook with a 2 GHz Quad-Core Intel Core i5. Times are from the programs themselves so do not include compilation time, unless the language used doesn't really allow for good timings.

* day 1 - SQL

Using PostgreSQL. Need to create the database "aoc2022_1" before running this.

Runtime: using wall clock via the "time" command:

    both parts: 50ms

* day 2 - Go

No table lookups...

Runtime:

    part1: 125µs
    part2: 80µs
    Total: 205µs

* day 3 - C

Plain old C, with a dash of apple and gcc specific stuff, to get the timing and for the __builtin_XXX functions to get the rightmost 1 bit in a word.

Runtime:

    part1: 85µs
    part2: 60µs
    Total: 145µs

* day 4 - Bash

Plain bash, no other programs used

Runtime (measured with "time")

    both parts: 88ms

* day 5 - Lua

I used lua 5.4, with the "lpeg" and "chronos" luarocks additions installed. lpeg for parsing, chronos just for displaying the time it takes.

Runtime:

    part1: 1.8ms
    part2: 1.3ms
    Total: 3.1ms

* day 6 - perl golf

I consider this different from regular perl. Normally when I write perl I want it to be easy to maintain. This type of perl is write-only and may hurt your eyes when you look at it.

Actually, the language that is used is "perl regexes". It could have been done with pcregrep except perl is easier to actually get the correct output. And the 14-character sequence would have been very tedious in pcregrep.

Wrapped in a shell script to be able to set some perl runtime flags.

Runtime (measured with "time" before the oneliners)

    part1: 7ms
    part2: 9ms
    Total: 16ms

* day 7 - Javascript

Using node to turn it into a cli

Runtime:

    part1: 15ms
    part2: 16ms
    Total: 31ms

* day 8 - Python

Python 3.10

Runtime:

    part1: 25ms
    part2: 37ms
    Total: 62ms

* day 9 - PHP

Plain old PHP (version 8) without any fancy extensions

Runtime:

    part1: 9ms
    part2: 24ms
    Total: 33ms

* day 10 -  Perl

Plain old perl. No meticulously trimmed grass fields in sight.

Runtime (both parts are calculated at the same time):

    Total: 250µs

* day 11 - Python

Yep, again python. No time to dive into another obscure language

Runtime:

    part1: 950µs
    part2: 469ms
    Total: 470ms

* day 12 - Python

... and again. Adding some OO but that's it.

Runtime:

    part1: 25ms
    part2: 18ms
    Total: 43ms

* day 13 - Python

again. Partly because the data structures could be parsed by python eval().

Runtime:

    parsing: 21ms
    part1: 0.6ms
    part2: 9ms
    total: 31ms

* day 14 - Go

again. Needed speed. Biggest speed improvent was using the recursive strategy as suggested by Cor.

Runtime:

    parsing: 430µs
    part1 & 2: 350µs

* day 15 - Go

Needed even more speed.

Runtime:

    part1: 20µs
    part2: 146ms
    Total: 146ms

* day 16 - Go

This took some time to get right... and needed the speed. The best speedup suggestion was to remap the vulcano to a graph where each valve is connected
to every other (non-zero or starting) valve, with a distance. Afterwards the problem space is a lot smaller, and trying the most promising paths first,
and dropping any paths that cannot possibly get better than the current best solution.

Runtime:

    part1: 5ms
    part2: 60ms
    Total: 65ms

* day 17 - Perl and Go

Started in perl because it looked simple enough, switched to Go for more speed (but that wasn't really needed, after the proper optimizations).

Runtime:

    part1: 26ms
    part2: 450µs
    Total: 26ms

* day 18 - Go

Using the useful feature in Go that you can use most data structures as a hash key.

Runtime:

    part1: 1ms
    part2: 8ms
    Total: 9ms

* day 19 - Go

This took a while to get right. After 2 false starts where I didn't even get the correct answers on the examples, or it took too long, I realised I misread
the description. Then I got the runtime to below a second using some hand-waving optimizations that did not feel right. In particular, I sorted the possible
solutions based on a "score" that tries to incorporate amount of materials available and materials needed, then took only the X best scoring solutions, and
dropped the rest. This gave the correct answer as long as I made "X" large enough. Experimentally, about 150 solutions were needed.

Then I got some simple suggestions from the reddit solutions thread, incorporated those, and that greatly reduced the problem space, so I could remove the
"best scoring" cutoff completely, and it made it faster, too.

Runtime:

    part1: 305ms
    part2: 117ms
    Total: 422ms

* day 20 - Python

As I thought this was going to be easy and would not take too long.

Also, I completely misjudged what is and isn't fast in python. For part 1 I had a solution that performed reasonably, using lots of array/slice copying.
Then I had the (failing) suggestion that instead of actually moving those slices around, I would just update a virtual position index. This turned out to be a mistake, but it did allow me to complete part 2, and I didn't feel like rewriting it all.

So, these runtimes are rather sub-optimal.

Runtime:

    part1: 2.2s
    part2: 23s
    Total: 25s

* day 21 - Python

(Mainly because I didn't feel like actually doing it in a spreadsheet program, but I guess that is possible too).

Runtime:

    Parsing: 10ms
    part1: 1.7ms
    part2: 1.8ms
    Total: 13ms

* day 22 - Go

I suspected I needed the speed, but part 2 went in a different direction than I thought (I expected part 2 to be "find the starting position where you can finish every run without running into a wall").

Anyway, for part 2 I expected that there could be different cube nets in the input, and I wanted to make it generic... so my code detects the cube net and
dynamically determines the mapping. Made me cut out several cube nets and fold them. Oh, and the dimensions of the cube are autodetected, too.

In case anyone wants to test this oin their own code, I added two test inputs:

day22/input/othernet.txt

    endpos part 1: {[4 1] 1}. Password: 2021
    endpos part 2: {[9 3] 0}, Password: 4040

day22/input/othernet2.txt

    endpos part 1: {[11 17] 1}. Password: 18049
    endpos part 2: {[20 4] 2}, Password: 5086

Runtime:

    parsing: 190µs
    part1: 370µs
    part2: 380µs
    Total: 1ms

* day 23 - Go

Again for data structures as keys in hash maps.

Runtime:

    part1: 4ms
    part2: 360ms
    total: 365ms

* day 24 - Go

More data structures as keys in has maps. There might be a nice speedup considering that the blizzard patterns repeat after a while, saving the trouble of
having to recalculate the position of all the blizzards. But I haven't implemented that (yet).

Runtime:

    part1: 270ms
    part2: 560ms
    total: 830ms

* day 25 - Python

Using operator overloading to implement the SNAFU numbers.

Runtime:

    total: 850µs

