package main

import (
	qrv2 "github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func main() {
	Qrcode2()
}

func Qrcode2() {
	qrc, err := qrv2.NewWith("https://wiki.eryajf.net",
		qrv2.WithErrorCorrectionLevel(qrv2.ErrorCorrectionHighest),
	)
	if err != nil {
		panic(err)
	}

	w, err := standard.New("v2.jpeg",
		standard.WithCircleShape(),
		standard.WithFgColorRGBHex("#0000ff"),
		standard.WithBgColorRGBHex("#ffffff"),
		standard.WithLogoImage(),
		standard.WithQRWidth(20),
		standard.WithBorderWidth(20),
	)
	if err != nil {
		panic(err)
	}

	err = qrc.Save(w)
	if err != nil {
		panic(err)
	}
}

// func Qrcode1() {
// 	qr, err := qrcode.New("https://wiki.eryajf.net", qrcode.Medium)
// 	if err != nil {
// 		log.Fatal(err)
// 	} else {
// 		qr.BackgroundColor = color.RGBA{50, 205, 50, 255}
// 		qr.ForegroundColor = color.White
// 		qr.WriteFile(256, "./zidingy.png")
// 	}
// }
