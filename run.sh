#!bin/sh
# start from project-path
cd scraper
go run main.go

#venvを使っている場合は予めactivateにするように
cd ../extractor
python extract.py

cd ../tabular
go run main.go

cd ..