#!bin/sh
# start from project-path
cd scraper
go run main.go

cd ../extractor
source .venv/bin/activate
python extract.py
deactivate

cd ../tabular
go run main.go

cd ..