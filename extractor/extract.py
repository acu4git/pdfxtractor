import pdfplumber
import pandas as pd
import os

# PDFファイルのパス
pdf_path = './target.pdf'
csv_path = './csv/'

# PDFファイルを開く
os.mkdir("./csv")
with pdfplumber.open(pdf_path) as pdf:
    tables = []
    for page in pdf.pages:
        # ページ内のテーブルを抽出
        extracted_tables = page.extract_tables()
        for table in extracted_tables:
            tables.append(table)

os.remove(pdf_path)

# テーブルをデータフレームに変換し、CSVに保存
for i, table in enumerate(tables):
    df = pd.DataFrame(table)
    df.to_csv(f'{csv_path}_table_{i}.csv', index=False)
