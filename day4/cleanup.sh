#!/usr/bin/env bash
# pure bash, no external programs except bash itself

error() {
    printf "%s\n" "$@" 1>&2
    exit -1
}

[ $# -eq 1 ] || error Specify input file

in="$1"

part12() {
    local contained=0
    local overlaps=0
    local r1start r1end r2start r2end
    while read line; do
        # make sure input line is in the correct format
        [[ $line =~ ^[0-9]+-[0-9]+,[0-9]+-[0-9]+$ ]] || error "Invalid input line $line"
        # split ranges into start and end
        IFS="$IFS-" read r1start r1end <<<${line%,*}
        IFS="$IFS-" read r2start r2end <<<${line#*,}
        # one range is fully contained in the other if start and end are contained in the other
        if [ $r1start -ge $r2start -a $r1start -le $r2end -a $r1end -ge $r2start -a $r1end -le $r2end ]; then
            contained=$(($contained + 1))
        elif [ $r2start -ge $r1start -a $r2start -le $r1end -a $r2end -ge $r1start -a $r2end -le $r1end ]; then
            contained=$(($contained + 1))
        # there is overlap if start or end of a range is within the other range
        elif [ \( $r1start -ge $r2start -a $r1start -le $r2end \) -o \( $r1end -ge $r2start -a $r1end -le $r2end \) ]; then
            overlaps=$(($overlaps + 1))
        fi
    done < "$1"
    echo part 1, total pairs where one contains the other: $contained
    echo part 2, total pairs where there is overlap: $(( $overlaps + $contained ))
}

part12 "$in"
