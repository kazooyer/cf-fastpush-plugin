CloudFoundry Fast Push Plugin
==

本`cf-fastpush-plugin` は [cf-fast-push](https://github.com/xiwenc/fastpush) プラグインにいくつかの機能を拡張したものです。

- fork元との主な違い
  - .cfignore のパースおよび有効化処理を実装(カレントディレクトリの.cfigonoreを読み込む）
  - -d --dry オプションで実行時、ファイルのアップロードを実行しないで差分だけを表示する
  - 複数引数のアプリ名に対応
  - cf fp/fps の引数なしで実行した場合にCurrent dirのmanifest.ymlを読み込み、 `command:` に `./fp` が含まれるappを自動検出する
  - cf uninstall-plugin 処理のバグを修正（ログイン状態ではない場合にcf uninstall-plugin できないエラーを解消）
  - アプリ内環境変数に「FP_PROTOCOL」を設定することで「http」か「https」かを選択可能(指定しない場合はhttpsになります）
  - アプリ内環境変数に「FP_DOMAIN」を設定することでfast-push用のドメインを指定可能(指定しない場合は最後のドメインになります
  - 以下のlocaleに対応
    - de-DE
    - en-US
    - es-ES
    - fr-FR
    - it-IT
    - ja-JP
    - ko-KR
    - pt-BR
    - zh-Hans
    - zh-Hant

Usage
===

fast-push-plugin を利用するためには専用の manifest.yml でアプリをデプロイしておく必要があります。

manifest.yml Sample

```YAML
applications
  command: wget -q -O ./fp https://github.com/xiwenc/cf-fastpush-controller/releases/download/v1.1.0/cf-fastpush-controller_linux_amd64 && chmod +x ./fp && ./fp
  routes:
  - route app-space-org.cfapps.io
  env
    BACKEND_COMMAND: python .bp/bin/start # for php-buildpack
    # BACKEND_COMMAND: node server.js # for nodejs-buildpack
    BACKEND_DIRS: ./
    BACKEND_PORT: 8081
    RESTART_REGEX: ""
    # RESTART_REGEX: "^*.js$" # for nodejs-buildpack
    FP_PROTOCOL: https
    FP_DOMAIN: cfapps.io
```

上記項目を追加したmanifest.ymlを予めpushしておきます。

```bash
$ cf push
```

デプロイしたアプリ名を確認し、cf fast-pushを実行します。(App名省略可）
```bash
$ cf fp [App]
```

Requirements
===

- Cloudfoundry CLI 6.x or later

Installation
===

* build

```bash
$ godep get
$ go build -gcflags="-trimpath=${HOME}"
```

* install

```bash
$ cf install-plugin ./cf-fastpush
```

* unintall

```bash
$ cf uninstall-plugin cf-fastpush
```


Commands
===

| Command | Short-cut | Description |
| --- | --- | --- |
| `cf fast-push <app name>` | `cf fp <app name>` | Update application files and restart app if needed. |
| `cf fast-push-status <app name>` | `cf fps <app name>` | Get status of the app. |
