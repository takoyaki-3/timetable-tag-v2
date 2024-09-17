# GTFS時刻表タグv2

**このレポジトリのコンセプトは https://github.com/takoyaki-3/butter へ引き継がれています**

## 概要

このリポジトリは、日本のGTFS (General Transit Feed Specification) データを効率的に配信するためのシステムです。複数の事業者から提供されるGTFSデータを収集、加工、配信するためのツールとスクリプトを含んでいます。

## 特徴

* GTFSデータの一括ダウンロードと更新
* データの分割と圧縮による配信効率の向上
* 電子署名によるデータの整合性保証
* APIエンドポイントによるデータアクセス
* S3へのデータアップロード

## ファイルツリー

```
├── 1_py.py
├── 2_2obj.py
├── 3_downloader.py
├── 4_unzip.py
├── 5_split.go
├── 6_add_gtfs_id_info.go
├── 7_add_datalist.go
├── 999_upload.go
├── get.go
├── go.mod
├── go.sum
└── n-gram.go
```

### ファイルの説明

* **1_py.py**: GTFS-data.jp から事業者とGTFSフィードの情報をスクレイピングし、テキストファイルに出力します。
* **2_2obj.py**: スクレイピング結果のテキストファイルを読み込み、JSON形式に変換します。各事業者にはGTFS APIのIDが追加されます。
* **3_downloader.py**: JSONファイルの情報に基づき、各事業者のGTFSデータをダウンロードします。
* **4_unzip.py**: ダウンロードしたGTFS zipファイルを解凍します。
* **5_split.go**: 解凍されたGTFSデータを停留所IDと旅程IDで分割し、tar.gz形式で圧縮します。また、電子署名を追加します。
* **6_add_gtfs_id_info.go**: 各GTFSデータにバージョン情報とハッシュ値サイズ情報を追加します。
* **7_add_datalist.go**: 配信するGTFSデータの一覧を作成し、JSONファイルに出力します。
* **999_upload.go**: 生成されたデータをS3にアップロードします。
* **get.go**: GTFS-data.jp から情報を取得するサンプルスクリプトです。
* **n-gram.go**: 停留所名のn-gramを生成するサンプルスクリプトです。


## 環境変数

* **AWS_ACCESS_KEY_ID**: AWSアクセスキーID
* **AWS_SECRET_ACCESS_KEY**: AWSシークレットアクセスキー
* **AWS_REGION**: AWSリージョン
* **S3_BUCKET_NAME**: S3バケット名

## APIエンドポイント

* **`/datalist.json`**: 配信データの一覧を取得
* **`/\{gtfs_id\}/info.json`**: 特定のGTFSデータのバージョン情報を取得
* **`/\{gtfs_id\}/\{version_id\}/byStops/\{hash\}.tar.gz`**: 停留所IDで分割されたGTFSデータを取得
* **`/\{gtfs_id\}/\{version_id\}/byTrips/\{hash\}.tar.gz`**: 旅程IDで分割されたGTFSデータを取得
* **`/\{gtfs_id\}/\{version_id\}/stops.txt`**: 停留所情報のデータを取得
* **`/\{gtfs_id\}/\{version_id\}/GTFS/\{filename\}`**: その他のGTFSデータを取得


## 設定ファイル

### s3-conf.json

```json
{
  "access_key_id": "YOUR_AWS_ACCESS_KEY_ID",
  "secret_access_key": "YOUR_AWS_SECRET_ACCESS_KEY",
  "region": "YOUR_AWS_REGION",
  "bucket": "YOUR_S3_BUCKET_NAME"
}
```

## インストール方法

1. Go 1.18以上をインストールしてください。
2. 必要な依存パッケージをインストールします。

```bash
go mod download
```

3. 環境変数を設定します。

## 使い方

1. `1_py.py` を実行してGTFS-data.jpから事業者情報をスクレイピングします。
2. `2_2obj.py` を実行してスクレイピング結果をJSONに変換します。
3. `3_downloader.py` を実行してGTFSデータをダウンロードします。
4. `4_unzip.py` を実行してダウンロードしたGTFSデータを解凍します。
5. `5_split.go` を実行してGTFSデータを分割・圧縮します。
6. `6_add_gtfs_id_info.go` を実行してバージョン情報などを追加します。
7. `7_add_datalist.go` を実行してデータリストを作成します。
8. `999_upload.go` を実行してS3にデータをアップロードします。

## コマンド実行例

```bash
# GTFS-data.jp から事業者情報をスクレイピング
python 1_py.py

# スクレイピング結果をJSONに変換
python 2_2obj.py

# GTFSデータをダウンロード
python 3_downloader.py

# ダウンロードしたGTFSデータを解凍
python 4_unzip.py

# GTFSデータを分割・圧縮・署名
go run 5_split.go

# バージョン情報などを追加
go run 6_add_gtfs_id_info.go

# データリストを作成
go run 7_add_datalist.go

# S3にデータをアップロード
go run 999_upload.go
```

## ライセンス

このプロジェクトは MIT ライセンスの下で公開されています。詳細については、LICENSE ファイルを参照してください。
