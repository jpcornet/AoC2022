#!/usr/bin/env perl

use strict;
use warnings;
use Time::HiRes qw(gettimeofday tv_interval);

my @rocks = map { [ reverse split /\n/, $_ ] } split("\n\n", <<END_OF_ROCKS);
####

.#.
###
.#.

..#
..#
###

#
#
#
#

##
##
END_OF_ROCKS

my $infile = shift
    or die "Specify input file\n";
open my $fh, '<', $infile or die "Cannot open $infile: $!\n";
my $streams = join('', <$fh>);
chomp $streams;

my $rocknr = 0;
my $streamnr = 0;
# the stack contains lines just like in the example or rocks, with . and # chars
my @stack = ();
my $width = 7;

sub next_rock {
    # returns the rock, x and y coordinate
    return ( $rocks[$rocknr++ % @rocks], 2, 3 + @stack );
}

sub next_stream {
    # returns the next gas steam direction
    return substr($streams, $streamnr++ % length($streams), 1);
}

my %streamdir = (
    '>' => 1,
    '<' => -1,
);

# return true if the projected rock has overlap between what is in the stack and the rock
sub has_overlap {
    my ($rock, $x, $y) = @_;
    for my $ry ( 0 .. $#$rock ) {
        if ( $y + $ry > $#stack ) {
            # above the stack, so cannot have overlap
            return;
        }
        for my $rx ( 0 .. length($rock->[$ry])-1 ) {
            return 1 if substr($rock->[$ry], $rx, 1) eq '#' and substr($stack[$y+$ry], $x + $rx, 1) eq '#';
        }
    }
    return;
}

sub add_rock_to_stack {
    my ($rock, $x, $y) = @_;

    for my $ry ( 0 .. $#$rock ) {
        push @stack, '.' x $width while $y + $ry >= @stack;
        for my $rx ( 0 .. length($rock->[$ry]) ) {
            substr($stack[$y+$ry], $x + $rx, 1, '#') if substr($rock->[$ry], $rx, 1) eq '#';
        }
    }
}

sub show_stack {
    for my $y ( reverse 0..$#stack ) {
        print "|$stack[$y]|\n";
    }
    print "|", "-" x $width, "|\n";
}

sub drop_one_rock {
    my ($rock, $x, $y) = next_rock();

    while ( 1 ) {
        my $dir = next_stream();
        my $dx = $streamdir{$dir} or die "Unknown gas stream direction in input: $dir\n";
        # check if movement is possible
        if ( $x + $dx >= 0 and $x + $dx + length($rock->[0]) - 1 < $width and !has_overlap($rock, $x + $dx, $y) ) {
            $x += $dx;
        }
        if ( $y == 0 or has_overlap($rock, $x, $y-1) ) {
            # add rock to stack
            add_rock_to_stack($rock, $x, $y);
            last;
        }
        $y--;
    }
}

my $start = [gettimeofday()];
for ( 1..2022 ) {
    drop_one_rock();
}
my $took = tv_interval($start);
show_stack();
print "Height of tower: ", scalar @stack, "\n";
print "Took: ", $took * 1e3, "ms\n";
