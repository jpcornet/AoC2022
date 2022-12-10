#!/usr/bin/env perl

use strict;
use warnings;

my $x = 1;
my $cycle = 0;
my $ss = 0;
my $at = 20;

while ( <> ) {
    my ($op, $val) = split;
    my $newx = $x;
    if ( $op eq 'noop' ) {
        $cycle += 1;
    } elsif ( $op eq 'addx' ) {
        $cycle += 2;
        $newx = $x + $val;
    } else {
        die "Invalid op $op\n";
    }
    if ( $cycle >= $at ) {
        chomp;
        $ss += $at * $x;
        print "At cycle $at, op = $_, x = $x, ss = $ss\n";
        $at += 40;
    }
    $x = $newx;
}
