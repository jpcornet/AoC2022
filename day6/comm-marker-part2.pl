#!/usr/bin/env perl
$"="";$r=qq[@{[map qq<(@{[map"(?!\\$_)",1..$_-1]}.)>,1..14]}];$_=<>;/$r/g;print pos
