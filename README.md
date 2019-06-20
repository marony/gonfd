# ファイル名 NFD⇔NFC変換

Mac ⇔ Windows, Linux間でファイルをコピーした際に、ユニコードのNFD正規化(Mac)とNFC正規化(Windows, Linux)のファイル名をどちらかに統一する。

ん～、はずなんだけどGo言語のos.Renameを実行してもファイル名が変更されない？？？

## Mac(NFD)に合わせる

```
go run gonfd.go -dir ディレクトリ -nfd -conv
```

## Windows, Linux(NFC)に合わせる

```
go run gonfd.go -dir ディレクトリ -nfc -conv
```
