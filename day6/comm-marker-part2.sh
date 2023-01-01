#!/bin/sh

time perl -ple 'BEGIN{$"=""}/@{[map qq<(@{[map"(?!\\$_)",1..$_-1]}.)>,1..14]}/g;$_=pos' $*
