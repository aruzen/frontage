# Frontage

盤面とカード効果を扱うGo製ゲームエンジン・サーバーです。C/C++向けのクライアントAPIは `lib/client/io` にあります。

## 実行

開発用のTCPサーバーを起動します。

```sh
go run ./cmd/server -addr 127.0.0.1:8275
```

TLSを使う場合は証明書と秘密鍵を両方指定します。

```sh
go run ./cmd/server -addr :8275 -tls-cert cert.pem -tls-key key.pem
```

PvEは現在、メイン／サブデッキのドローとパスによるターン進行に対応しています。NPCはドローしてパスします。カード使用・スキル・複合ドローは実装途中です。クライアントはサーバーからの `SyncPacket` に同じ `sync_cause` を返してください。

## 検証

```sh
go test -race ./...
```

盤面・イベントの試行実行が元の状態を変更しないこと、TCP／TLSでの接続・6ターン進行・切断・再接続、同期タイムアウトを確認します。

`Board.Copy()` と `Sandbox()` はプレイヤー・カード・駒の可変データをコピーします。StructureとSkillの定義は共有するため、不変として扱います。過去の盤面を保持する場合はCHAIN方式または `Copy()` を使ってください。SWAP方式は2世代の盤面を再利用します。内部で別のActionStateを保持する状態型は `ActionStateCloner` を実装し、試行実行時の参照関係を保ちます。

## License

frontage © 2025 ahaha-craft.org  
Licensed under the MIT License with the Commons Clause restriction.  
Commercial use or resale is prohibited.
