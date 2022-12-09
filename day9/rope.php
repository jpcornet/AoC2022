#!/usr/bin/env php
<?php

function main() {
    global $argv, $argc;
    if ( $argc != 2 ) {
        print("Provide input filename\n");
        exit(-1);
    }
    $starttime = hrtime(true);
    $input = file($argv[1]);
    $parsetime = hrtime(true);
    $part1 = follow_tail($input, 1);
    $part1time = hrtime(true);
    $part2 = follow_tail($input, 9);
    $part2time = hrtime(true);
    echo "using 1 knot the tail touches $part1 positions\n";
    echo "using 9 knots, the tail touches $part2 positions\n";
    echo "Loading input took: ", ($parsetime - $starttime) / 1000, "µs\n";
    echo "part 1 took: ", ($part1time - $parsetime) / 1000, "µs\n";
    echo "part 2 took: ", ($part2time - $part1time) / 1000, "µs\n";
}

function one_step(&$object, $direction) {
    foreach ([0, 1] as $i) {
        $object[$i] += $direction[$i];
    }
}

function distance($o1, $o2) {
    $dist = [];
    foreach ([0, 1] as $i) {
        $dist[$i] = $o1[$i] - $o2[$i];
    }
    return $dist;
}

function sign ($x) {
    return ($x > 0) - ($x < 0);
}

function follow_tail($lines, $knots) {
    $head = [0, 0];
    foreach (range(0, $knots - 1) as $i) {
        $knot[$i] = [0, 0];
    }
    $tailpos = [];
    $direction = [
        "U" => [0, -1],
        "D" => [0, 1],
        "L" => [-1, 0],
        "R" => [1, 0],
    ];
    foreach ($lines as $line) {
        list($dirstr, $steps) = explode(" ", rtrim($line));
        $dir = $direction[$dirstr];
        if ( !$dir ) {
            print("Invalid direction $dirstr\n");
            exit(1);
        }
        foreach (range(1,$steps) as $dummy) {
            # update position of head
            one_step($head, $dir);
            # make knots follow head. prev points to the previous knot, and starts with the head.
            $prev = $head;
            foreach (range(0, $knots - 1) as $i) {
                $dist = distance($prev, $knot[$i]);
                # is any distance > 1 ?
                if ( abs($dist[0]) > 1 or abs($dist[1]) > 1 ) {
                    # move tail in direction of head, for both x and y
                    one_step($knot[$i], [ sign($dist[0]), sign($dist[1]) ]);
                } else {
                    # no need to move anything, so no knots further down move either
                    break;
                }
                $prev = $knot[$i];
            }
            @$tailpos[implode(",", $knot[$knots - 1])]++;
        }
    }
    return count($tailpos);
}

main();
