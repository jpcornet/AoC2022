#!/bin/sh
# perl golfing it

time perl -ple '/(.)((?!\1).)((?!\1)(?!\2).)((?!\1)(?!\2)(?!\3).)/g;$_=pos' $*
