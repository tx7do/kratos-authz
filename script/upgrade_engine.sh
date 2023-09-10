#!/bin/sh

root_path=$(pwd)
sub_path=../engine
for folder in `find $sub_path/* -type d`
do
    cur_path=`realpath $root_path/$folder`
    echo $cur_path
    cd $cur_path
    go get all
    go mod tidy
done

cd $root_path
cd $root_path/middleware
go get all
go mod tidy
