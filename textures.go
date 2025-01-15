package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/png"

	g "github.com/AllenDang/giu"
)

var (
	rgbaIcon16       *image.RGBA
	rgbaIcon         *image.RGBA
	headerSplashRgba *image.RGBA
	splashTexture    = &g.ReflectiveBoundTexture{}
	icon16Texture    = &g.ReflectiveBoundTexture{}
)

func DecodeEmbedded(data []byte) (*image.RGBA, error) {
	r := bytes.NewReader(data)
	img, err := png.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("LoadImage: error decoding png image: %w", err)
	}
	return g.ImageToRgba(img), nil
}

//go:embed winres/splash.png
var splashHeaderLogo []byte

func DecodeSplashHeaderLogo() (*image.RGBA, error) {
	return DecodeEmbedded(splashHeaderLogo)
}

//go:embed winres/icon16.png
var appIcon16 []byte

func DecodeAppIcon16() (*image.RGBA, error) {
	return DecodeEmbedded(appIcon16)
}

//go:embed winres/icon.png
var appIcon []byte

func DecodeAppIcon() (*image.RGBA, error) {
	return DecodeEmbedded(appIcon)
}

func initTextures() {
	rgbaIcon, _ = DecodeAppIcon()
	rgbaIcon16, _ = DecodeAppIcon16()
	headerSplashRgba, _ := DecodeSplashHeaderLogo()
	splashTexture.SetSurfaceFromRGBA(headerSplashRgba, false)
	icon16Texture.SetSurfaceFromRGBA(rgbaIcon16, false)
}
