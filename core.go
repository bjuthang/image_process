package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	// "github.com/gofiber/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
)

func serverInit() {
	flag.StringVar(&configPath, "config", "config.json", "/path/to/config.json. (Default: ./config.json)")
	flag.IntVar(&jobs, "jobs", runtime.NumCPU(), "Prefetch thread, default is all.")
	flag.BoolVar(&verboseMode, "v", false, "Verbose, print out debug info.")
	flag.Parse()
}

func loadConfig(path string) Config {
	jsonObject, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(jsonObject)
	_ = decoder.Decode(&config)
	_ = jsonObject.Close()
	return config
}

func main() {
	rand.Seed(time.Now().UnixNano())

	serverInit()
	config = loadConfig(configPath)

	logfile, err := os.OpenFile("error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Print("create logfile error, err:", err)
		log.SetOutput(os.Stdout)
	} else {
		defer logfile.Close()
		log.SetOutput(logfile)
	}

	log.SetReportCaller(true)
	Formatter := &log.TextFormatter{
		ForceQuote:      true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return fmt.Sprintf("[%s()]", f.Function), ""
		},
	}
	log.SetFormatter(Formatter)

	if verboseMode {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug mode is enable!")
	} else {
		log.SetLevel(log.InfoLevel)
	}

	app := fiber.New(fiber.Config{
		ServerHeader:          "Webp-Proxy-Server",
		DisableStartupMessage: true,
	})
	app.Use(logger.New())

	listenAddress := config.Host + ":" + config.Port

	// app.Get("/clean/image", cleanHandlerFunc)
	app.Get("/hd1", handlerFunc)
	app.Get("/hd2", handlerFunc2)
	app.Get("/*", handlerFunc2)

	banner := fmt.Sprintf(`
	==============================
	==== image process server ====
	==============================
	`)
	fmt.Printf("\n %c[1;32m%s%c[0m\n\n", 0x1B, banner, 0x1B)
	fmt.Println("Webp-Proxy-Server is Running on http://" + listenAddress)

	_ = app.Listen(listenAddress)
}
