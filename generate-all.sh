#!/bin/bash

SEED=7
SIZE=500
SCALE=0.75

mkdir -p datasets
mkdir -p svg

for DISTRIBUTION in "uniform" "normal" "line" "line-strict" "exp"; do
    for SORTED_P in "sorted,true" "random,false"; do
        IFS="," read SORTED SORTED_F <<< "$SORTED_P"
        DB="datasets/$DISTRIBUTION-$SORTED.db"
        echo $DB
        ./quadsql --db $DB generate 2 $SIZE $DISTRIBUTION $SORTED_F --seed $SEED > /dev/null 2> /dev/null
        echo "/svg $SCALE $DISTRIBUTION-$SORTED-point" | ./quadsql --db $DB > /dev/null 2> /dev/null
        echo "/svg $SCALE $DISTRIBUTION-$SORTED-region" | ./quadsql --db $DB --region > /dev/null 2> /dev/null
    done
done

mv *.svg svg/