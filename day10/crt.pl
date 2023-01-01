#!/usr/bin/env perl

use strict;
use warnings;
use Time::HiRes qw(gettimeofday tv_interval);

my $x = 1;
my $cycle = 0;
# ss is the total "signal strength" at cycle 20, 60, ... for part 1
my $ss = 0;
# at what moment we should do part 1 signal strength
my $at = 20;
# collect CRT lines here (as one line)
my $crt = '';

my $start = [gettimeofday()];
while ( <> ) {
    my ($op, $val) = split;
    my $newx = $x;
    # number of cycles this instruction takes
    my $takes;
    if ( $op eq 'noop' ) {
        $takes = 1;
    } elsif ( $op eq 'addx' ) {
        $takes = 2;
        $newx = $x + $val;
    } else {
        die "Invalid op $op\n";
    }
    for ( 1 .. $takes ) {
        # if the sprite is at the CRT position ( = cycle num ) then we draw a lit pixel
        my $pixel = abs($x - $cycle % 40) <= 1 ? '#' : '.';
        $crt .= $pixel;
        $cycle++;
    }
    # part 1 signal strength
    if ( $cycle >= $at ) {
        chomp;
        $ss += $at * $x;
        print "At cycle $at, op = $_, x = $x, ss = $ss\n";
        $at += 40;
    }
    $x = $newx;
}

print "Part 2, CRT result:\n";
print substr($crt, $_ * 40, 40), "\n" for 0 .. 5;
print "Total runtime: ", tv_interval($start) * 1e6, "µs\n";
