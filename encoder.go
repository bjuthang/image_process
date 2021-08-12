package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math/rand"
	"path"
	"strconv"
	"strings"
	"time"

	// "github.com/gofiber/fiber"
	"github.com/chai2010/webp"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/bmp"
)

// p1: 图片原路径
// p2: 图片处理后的存储路径
func webpEncoder(p1, p2 string, quality float32, c chan int) (err error) {
	now := time.Now()
	log.Info("webpEncode start:", now)
	log.Debugf("target: %s with quality of %f", path.Base(p1), quality)
	var buf bytes.Buffer
	var img image.Image

	data, err := ioutil.ReadFile(p1)
	if err != nil {
		chanErr(c)
		return err
	}

	contentType := getFileContentType(data[:512])
	if strings.Contains(contentType, "jpeg") {
		img, _ = jpeg.Decode(bytes.NewReader(data))
	} else if strings.Contains(contentType, "png") {
		img, _ = png.Decode(bytes.NewReader(data))
	} else if strings.Contains(contentType, "bmp") {
		img, _ = bmp.Decode(bytes.NewReader(data))
	} else if strings.Contains(contentType, "gif") {
		log.Warn("Gif support is not perfect!")
		img, _ = gif.Decode(bytes.NewReader(data))
	}

	if img == nil {
		msg := "image file " + path.Base(p1) + " is corrupted or not supported"
		log.Debug(msg)
		err = errors.New(msg)
		chanErr(c)
		return err
	}

	if err = webp.Encode(&buf, img, &webp.Options{Lossless: false, Quality: quality}); err != nil {
		log.Error(err)
		chanErr(c)
		return err
	}
	if err = ioutil.WriteFile(p2, buf.Bytes(), 0644); err != nil {
		log.Error(err)
		chanErr(c)
		return err
	}

	log.Info("Save to " + p2 + " ok!\n")

	log.Info("webpEncode span:", time.Since(now))
	chanErr(c)

	return nil
}

// 输入url，输出压缩图片路径
// func imageFetchAndWebpEncoder(url string, quality float32, Log bool, c chan int) interface{} {
func imageFetchAndWebpEncoder(url string, quality float32) (err error, imagePath string) {
	//var randFile = string(rune(rand.Intn(100)))
	var randFile = strconv.Itoa(rand.Intn(100))
	var rawImage = config.ImgPath + "/" + randFile
	var webpImage = config.ExhaustPath + "/" + randFile

	err = fetchRemoteImage(rawImage, url)
	if err != nil {
		msg := fmt.Sprint("fetch remote image error, err:", err)
		log.Error(msg)
		return err, ""
	}

	err = webpEncoder(rawImage, webpImage, quality, nil)
	if err != nil {
		log.Error(err)
		return err, ""
	}
	return nil, webpImage
}

// 请求，及响应压缩图片
func handlerFunc(c *fiber.Ctx) error { /// 127.0.0.1:3333/abc/big.jpg
	// var reqURI, _ = url.QueryUnescape(c.Path()) // /abc/big.jpg
	// var url = "host" + reqURI
	var url = c.Get("url")

	var ua = c.Get("User-Agent")
	var debug = c.Get("debug")

	if debug != "" {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("Incoming connection from %s@%s", ua, c.IP())

	err, webpImage := imageFetchAndWebpEncoder(url, config.Quality)
	if err != nil {
		msg := fmt.Sprint("fetch remote image error, err:", err)
		log.Error(msg)
		c.SendStatus(503)
		return c.Send([]byte(msg))
	}

	return c.SendFile(webpImage)
}
