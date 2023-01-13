#!/bin/bash

configs=`ls *default.toml`

# echo "${configs[*]}"

for def_conf in $configs
do
    conf=`echo "${def_conf%%_default.toml*}"`
    cat $def_conf > "${conf}.toml"
done