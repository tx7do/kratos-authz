echo off

::指定起始文件夹
set DIR="%cd%\..\engine"

for /R %DIR% /d %%i in (*) do (
    echo %%i
    cd %%i
    go get all
    go mod tidy
)


cd ..\..\middleware\
go get all
go mod tidy

pause
