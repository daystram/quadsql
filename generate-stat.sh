#!/bin/bash

mkdir -p stat/datasets

SEED=7

DISTS=(normal uniform line-strict)
SORTED_PS=(random,false sorted,true)
# SIZES=(10 50 100 250 500 1000 1500 2000 2500 3000 5000 10000 25000 50000 75000 100000 125000 150000 175000 200000)
SIZES=(5 10 50 100 250 500 750 1000 1500 2000 2500 3000 3500 4000 4500 5000 7500 10000 12500 15000 17500 20000 25000 50000 75000 100000)
DIMS=(2)

RUNS=1
MODE=point,region,none

for DIST in ${DISTS[@]}; do
    for SORTED_P in ${SORTED_PS[@]}; do
        IFS="," read SORTED SORTED_F <<< "$SORTED_P"
        CSV="stat/$DIST-$SORTED.csv"
        echo "dim,row,index_type,node_count,max_depth,build_time,avg_exec_time,avg_index_access,avg_table_access,avg_point_comp,runs" > $CSV
        for SIZE in ${SIZES[@]}; do
            for DIM in ${DIMS[@]}; do
                DB="stat/datasets/$DIST-$SORTED.db"
                echo "---------  $DIST $SORTED, $SIZE points @ ${DIM}D"
                ./quadsql --db $DB generate $DIM $SIZE $DIST $SORTED_F --seed $SEED > /dev/null 2> /dev/null
                ./quadsql --db $DB statistic --mode $MODE --runs $RUNS --no-head > /dev/null 2>> $CSV
            done
        done
    done
done
