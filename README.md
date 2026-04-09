# VPMBuilder

VRChat Package Manager (VPM) のパッケージとリポジトリマニフェストをビルドするための CLI ツールです。
Unity プロジェクト内の VPM パッケージを zip にまとめ、GitHub Releases から `repo.json` を自動生成できます。

## 特長

- **パッケージビルド**: 指定した Unity パッケージディレクトリを zip に圧縮し、`packages.json` を生成
- **差分ビルド**: 既存リポジトリの最新リリースと比較し、バージョンが同一のパッケージはスキップ
- **並列処理**: `sourcegraph/conc` によるワーカープールで複数パッケージを並列ビルド
- **リポジトリマニフェスト生成**: GitHub Releases のアセットを走査して VPM `repo.json` を自動生成
- **マニフェストマージ**: 複数の VPM リポジトリ URL を 1 つの `repo.json` に統合

## ビルド
```bash
go build -o VPMBuilder .
```

必要環境: Go 1.26 以上

## 使い方

### 1. 設定ファイル

`config.json` をプロジェクトルートに配置します:

```json
{
  "logLevel": "info",
  "outputPath": "./output/",
  "repoTemplate": "./repo_template.json",
  "downloadRepo": "owner/repo",
  "packagePaths": [
    "D:\\VR\\Project\\MyAvatar\\Packages\\com.example.tool"
  ],
  "repoUrls": []
}
```

| フィールド       | 説明                                                                 |
| ---------------- | -------------------------------------------------------------------- |
| `logLevel`       | ログレベル (`debug` / `info` / `warn` / `error`)                     |
| `outputPath`     | ビルド成果物の出力先ディレクトリ                                     |
| `repoTemplate`   | `repo.json` のテンプレートファイルパス                               |
| `downloadRepo`   | 差分ビルド用の GitHub リポジトリ (`owner/repo` 形式、空なら全ビルド) |
| `packagePaths`   | ビルド対象の Unity パッケージディレクトリ一覧                        |
| `repoUrls`       | `merge` コマンドで統合する VPM リポジトリ URL 一覧                   |

### 2. リポジトリテンプレート

`repo_template.json`:

```json
{
  "name": "My VPM Repository",
  "id": "com.example.vpm",
  "url": "https://example.com/vpm/repo.json",
  "author": "Example Author",
  "packages": {}
}
```

### 3. コマンド

#### パッケージをビルド

```bash
VPMBuilder build -t package
```

`config.json` の `packagePaths` に指定された各パッケージを zip 化し、`outputPath` に
`<name>-<version>.zip` と `packages.json` を出力します。

`downloadRepo` が設定されている場合、その GitHub リポジトリの最新リリースに含まれる
`packages.json` と比較し、同一バージョンのパッケージはスキップします。

#### `repo.json` をビルド

```bash
VPMBuilder build -t repo
```

`downloadRepo` で指定された GitHub リポジトリのリリースを走査し、各リリースの zip
アセットから VPM リポジトリマニフェスト (`repo.json`) を生成します。

#### 複数リポジトリをマージ

```bash
VPMBuilder merge
```

`config.json` の `repoUrls` に指定されたリモート VPM リポジトリの `packages` を、
`repoTemplate` をベースに 1 つのマニフェストへ統合します。

### 4. 共通フラグ

| フラグ              | 説明                                |
| ------------------- | ----------------------------------- |
| `-c`, `--config`    | 設定ファイルのパス (既定: `config.json`) |
| `-t`, `--type`      | `build` のタイプ (`package` / `repo`) |

## プロジェクト構成

```
VPMBuilder/
├── main.go              # エントリポイント
├── cmd/                 # cobra コマンド定義 (build / merge)
├── internal/
│   ├── builder/         # パッケージ・リポジトリのビルドロジック
│   ├── config/          # 設定ファイルのパース
│   └── log/             # zap ベースのロガー
└── pkg/
    ├── github/          # GitHub Releases API クライアント
    ├── vpm/manifest/    # VPM パッケージ / リポジトリのマニフェスト型定義
    └── zip/             # ディレクトリ圧縮ユーティリティ
```

## 依存ライブラリ

- [spf13/cobra](https://github.com/spf13/cobra) — CLI フレームワーク
- [go-resty/resty](https://github.com/go-resty/resty) — HTTP クライアント
- [sourcegraph/conc](https://github.com/sourcegraph/conc) — 並列実行プール
- [uber-go/zap](https://github.com/uber-go/zap) — 構造化ロガー

## ライセンス

このリポジトリのライセンス表記に従います。
