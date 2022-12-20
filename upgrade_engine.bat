cd .\engine\casbin\
go get all
go mod tidy

cd ..\..\engine\opa\
go get all
go mod tidy

cd ..\..\engine\zanzibar\
go get all
go mod tidy

cd ..\..\middleware\
go get all
go mod tidy

pause