echo off

::指定起始文件夹
set DIR="%cd%\..\engine"

for /d %%i in ("%DIR%\*") do (
    echo %%i
    pushd "%%i"
    go get all
    go mod tidy
    popd
)

cd %cd%\..\middleware\
go get all
go mod tidy
