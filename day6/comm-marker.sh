#!/bin/sh
# perl golfing it

perl -ple '/(.)((?!\1).)((?!\1)(?!\2).)((?!\1)(?!\2)(?!\3).)/g;$_=pos' $*
