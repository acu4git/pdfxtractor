package main

import (
	"fmt"
	"log"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func main() {
	// ImageMagic を初期化する。
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// 解像度を設定する。
	err := mw.SetResolution(200, 200)
	if err != nil {
		log.Fatal("failed at SetResolution", err)
	}

	// 変換元のPDFを読み込む。
	filename1 := "R6hpsyougakukinitiran20240624_page_1.pdf"
	// filename2 := "SQL1_2014.pdf"
	err = mw.ReadImage(filename1)
	if err != nil {
		log.Fatal("failed at ReadImage", err)
	}

	// ページ数を取得する。
	n := mw.GetNumberImages()
	log.Println("number image: ", n)

	// 出力フォーマットをPNGに設定する。
	err = mw.SetImageFormat("jpg")
	if err != nil {
		log.Fatal("failed at SetImageFormat")
	}

	// １ページずつ変換して出力する。
	for i := 0; i < int(n); i++ {
		// ページ番号を設定する。
		if ret := mw.SetIteratorIndex(i); !ret {
			break
		}

		// 画像を出力する。
		err = mw.WriteImage(fmt.Sprintf("test_%02d.jpg", i))
		if err != nil {
			log.Fatal("failed at WriteImage")
		}
	}
}
