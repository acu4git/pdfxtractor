# 概要

このアプリは，京都工芸繊維大学が掲載している奨学金一覧から，自分が対象に含まれるような項目を取り出してまとめる目的で作成しています．

# 開発環境

- go1.22.4 linux/amd64
- Python 3.10.12
- Ubuntu 22.04.4 LTS on Windows 10 x86_64

# 使い方

extractor では requirements.txt から pip install をお願いします．

```bash
pip install -r requirements.txt
```

また，Go の各モジュール(scraper, tabular)にて go mod tidy しておいてください

```bash
go mod tidy
```

プロジェクトのルートディレクトリにて実行

```bash
source run.sh
```

<!-- `-t`: 対象（undergraduate，graduate）(追加予定) -->

新たに result.txt が生成され，処理結果がそのファイルに保存されます．

# 実装予定

- 引数を指定し，必要な情報だけ抽出できるようにする（より高精度に）
- いちいちローカルにファイルを作成せずに，レスポンスから直接操作できるようにする（多分難しい）
- python または Go 言語 のみで実装（気が向いたら）
- Web 版の実装
- TUI 化
