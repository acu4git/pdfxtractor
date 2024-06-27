import pdfplumber
import pandas as pd

# PDFファイルのパス
pdf_path = './sample.pdf'
csv_path = './csv_plum/'

# PDFファイルを開く
with pdfplumber.open(pdf_path) as pdf:
    tables = []
    for page in pdf.pages:
        # ページ内のテーブルを抽出
        extracted_tables = page.extract_tables()
        for table in extracted_tables:
            tables.append(table)

# テーブルをデータフレームに変換し、CSVに保存
for i, table in enumerate(tables):
    df = pd.DataFrame(table)
    df.to_csv(f'{csv_path}_table_{i}.csv', index=False)
