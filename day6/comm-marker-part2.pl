#!/usr/bin/env perl
$"="";$_=<>;/@{[map qq<(@{[map"(?!\\$_)",1..$_-1]}.)>,1..14]}/g;print pos
