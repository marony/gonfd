package main

import (
	"fmt"
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/text/unicode/norm"
)

// まいんちゃん
func main() {
	// コマンド引数
	var (
		dir = flag.String("dir", "~/", "検索ディレクトリ")
		nfd = flag.Bool("nfd", true, "NFDファイル名を検索")
		nfc = flag.Bool("nfc", false, "NFCファイル名を検索")
		conv = flag.Bool("conv", false, "ファイル名を変更する")
	)
	flag.Parse()

	// [【Go言語】Ctrl\+cなどによるSIGINTの捕捉とdeferの実行 \- DRYな備忘録](http://otiai10.hatenablog.com/entry/2018/02/19/165228)
	defer teardown()

	// シグナル用チャネル
	c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)

	// 終了検知用チャネル
    done := make(chan error, 1)
    go do(done, *dir, *nfd, *nfc, *conv)

    select {
    case sig := <-c:
        fmt.Println("シグナル来た:", sig)
        /*
         teardown中に再度SIGINTが来る場合を考慮し、
         send on closed channelのpanicを避ける。
       */
        // close(c)
        return
    case err := <-done:
        fmt.Println("loopの終了:", err)
    }

	fmt.Println("終了")
}

// 実際の処理
func do(done chan<- error, dir string, nfd, nfc, conv bool) {
	// ディレクトリパスの"~"を展開する
	dir, err := homedir.Expand(dir)
	if err != nil {
		panic(err)
	}

	fmt.Println("検索ディレクトリ: ", dir, "NFD: ", nfd, "NFC: ", nfc, "CONV: ", conv)

	walk(dir, nfd, nfc, conv)

	// 終了
	done <- nil
    close(done)
}

func walk(dir string, nfd, nfc, conv bool) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// if info.IsDir() {
		// 	fmt.Println("検索中: ", path)
		// }
		dir, base := filepath.Split(path)
		if nfd {
			// NFCに変換したら異なるのであればNFD
			s := norm.NFC.String(base)
			if base != s {
				fmt.Printf("NFD: %s != %s\n", base, s)
				if conv {
					fmt.Printf("rename: %s/%s -> %s\n", dir, base, s)
					src := filepath.Join(dir, base)
					dst := filepath.Join(dir, s)
					// fmt.Printf("src = '%v', dst = '%v'\n", ([]byte(src)), ([]byte(dst)))
					if err := os.Rename(src, dst); err != nil {
						panic(err)
					}
				}
			}
		}
		if nfc {
			// NFDに変換したら異なるのであればNFC
			s := norm.NFD.String(base)
			if base != s {
				fmt.Printf("NFC: %s != %s\n", base, s)
				if conv {
					fmt.Printf("rename: %s/%s -> %s\n", dir, base, s)
					src := filepath.Join(dir, base)
					dst := filepath.Join(dir, s)
					// fmt.Printf("src = '%v', dst = '%v'\n", ([]byte(src)), ([]byte(dst)))
					if err := os.Rename(src, dst); err != nil {
						panic(err)
					}
				}
			}
		}
		return nil
	})
}

func teardown() {
	fmt.Println("データのあとかたづけ")
}
