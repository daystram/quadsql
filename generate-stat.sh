#!/bin/bash
set -e

mkdir -p stat/datasets

SEED=7
RUNS=10

run() {
    TOTAL=$((${#DISTS[@]}*${#SORTED_PS[@]}*${#SIZES[@]}*${#DIMS[@]}))
    STEP=0
    mkdir -p "stat/$EXP"
    for DIST in ${DISTS[@]}; do
        for SORTED_P in ${SORTED_PS[@]}; do
            IFS="," read SORTED SORTED_F <<< "$SORTED_P"
            CSV="stat/$EXP/$DIST-$SORTED.csv"
            echo "dim,row,index_type,node_count,max_depth,build_time,avg_exec_time,avg_index_access,avg_table_access,avg_point_comp,runs" > $CSV
            for DIM in ${DIMS[@]}; do
                for SIZE in ${SIZES[@]}; do
                    ((STEP++))
                    DB="stat/datasets/$DIST-$SORTED.db"
                    echo "[Exp: $EXP] ----------------------------- [$STEP/$TOTAL] $DIST $SORTED, $SIZE points @ ${DIM}D"
                    ./quadsql --db $DB generate $DIM $SIZE $DIST $SORTED_F --seed $SEED > /dev/null 2> /dev/null
                    ./quadsql --db $DB statistic --mode $MODE --runs $RUNS --no-head 2>> $CSV
                done
            done
        done
    done
}

# Experiment 1
EXP=1
DISTS=(uniform line-strict normal)
SORTED_PS=(random,false sorted,true)
SIZES=(100 250 500 750 1000 2500 5000 7500 10000 20000 40000 60000 80000 100000)
DIMS=(2)
MODE=point,none
run

# Experiment 2
EXP=2
DISTS=(uniform)
SORTED_PS=(random,false)
SIZES=(10 50 100 250 500 750 1000 2500 5000 7500 10000 20000 40000 60000 80000 100000 200000 400000 600000 800000 1000000 2000000 3000000 4000000 5000000)
DIMS=(2)
MODE=point
run

# Experiment 3
EXP=3
DISTS=(uniform)
SORTED_PS=(random,false)
SIZES=(100 250 500 750 1000 2500 5000 7500 10000 20000 40000 60000 80000 100000 200000 400000 600000 800000 1000000)
DIMS=(2 3 4 5 6)
MODE=point
run

# Experiment 4
EXP=4
DISTS=(uniform)
SORTED_PS=(random,false sorted,true)
SIZES=(100 250 500 750 1000 2500 5000 7500 10000 20000 40000 60000 80000 100000)
DIMS=(2)
MODE=point,region
run
