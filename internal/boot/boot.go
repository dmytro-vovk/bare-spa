package boot

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/sysboard"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/webserver/middleware"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"

	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/webserver"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/webserver/auth"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/webserver/handlers/ws"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/webserver/handlers/ws/client"
)

type Boot struct {
	*Container
	configFile string
}

func New(configFile string) (*Boot, func()) {
	c, fn := NewContainer()
	boot := Boot{
		Container:  c,
		configFile: configFile,
	}
	// todo: we don't need to return shutdowner
	// these functions will be automatically called when a signal received
	return &boot, fn
}

func (c *Boot) Config() *Config {
	id := "Config"
	if s, ok := c.Get(id).(*Config); ok {
		return s
	}

	f, err := os.Open(c.configFile)
	if err != nil {
		panic(err)
	}

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}

	// сейчас в конфиге ничего не стоит, поэтому будет использоваться дефолтная борда
	if cfg.SystemBoard.Type == "" {
		cfg.SystemBoard.Type = sysboard.Omega
	}

	log.Printf("Config read from %s", c.configFile)
	c.Set(id, &cfg, nil)

	return &cfg
}

func (c *Boot) Application() *app.Application {
	id := "Application"
	if s, ok := c.Get(id).(*app.Application); ok {
		return s
	}

	a := app.New()
	a.UseBoard(c.Config().SystemBoard.Type)
	c.Set(id, a, nil)
	return a
}

func (c *Boot) Auth() *auth.PlainAuth {
	id := "Auth"
	if s, ok := c.Get(id).(*auth.PlainAuth); ok {
		return s
	}

	a := auth.New(
		c.Config().Auth.UserName,
		c.Config().Auth.Password,
	)
	c.Set(id, a, nil)

	return a
}

func (c *Boot) WSClient() *client.Client {
	id := "WS Client"
	if s, ok := c.Get(id).(*client.Client); ok {
		return s
	}

	s := client.New().
		NS("datetime",
			client.NSMethod("get", c.Application().GetDatetime),
			client.NSMethod("set", c.Application().SetDatetime),
			client.NSMethod("getNTP", c.Application().GetNTP),
			client.NSMethod("setNTP", c.Application().SetNTP),
		).
		NS("network",
			client.NSMethod("getWiFiAp", c.Application().GetWiFiAp),
			client.NSMethod("setWiFiAp", c.Application().SetWiFiAp),
			client.NSMethod("getWiFiCl", c.Application().GetWiFiCl),
			client.NSMethod("setWiFiCl", c.Application().SetWiFiCl),
			client.NSMethod("getWAN", c.Application().GetWAN),
			client.NSMethod("setWAN", c.Application().SetWAN),
			client.NSMethod("getDNS", c.Application().GetDNS),
			client.NSMethod("setDNS", c.Application().SetDNS),
		).
		NS("gpio",
			client.NSMethod("get", c.Application().GetGPIO),
			client.NSMethod("setDirection", c.Application().SetGPIODirection),
			client.NSMethod("setInversion", c.Application().SetGPIOInversion),
			client.NSMethod("setLevel", c.Application().SetGPIOLevel),
		).
		NS("uart",
			client.NSMethod("getConfig", c.Application().GetUARTConfig),
			client.NSMethod("setConfig", c.Application().SetUARTConfig),
		).
		NS("spi",
			client.NSMethod("getConfig", c.Application().GetSPIConfig),
			client.NSMethod("setConfig", c.Application().SetSPIConfig),
		).
		NS("i2c",
			client.NSMethod("getConfig", c.Application().GetI2CConfig),
			client.NSMethod("setConfig", c.Application().SetI2CConfig),
		).
		NS("settings",
			client.NSMethod("getSystemBoardTypes", c.Application().GetSystemBoardTypesList),
		).
		NS("system",
			client.NSMethod("getInterfaces", c.Application().GetInterfaces),
			client.NSMethod("getTemperatures", c.Application().GetTemperatures),
			client.NSMethod("logRead", c.Application().LogRead),
		).
		NS("sysboard",
			client.NSMethod("get", c.Application().GetSysboard),
		).
		NS("devices",
			client.NSMethod("create", c.Application().CreateDevice),
			client.NSMethod("read", c.Application().ReadDevice),
			client.NSMethod("update", c.Application().UpdateDevice),
			client.NSMethod("delete", c.Application().DeleteDevice),
			client.NSMethod("supportedTypes", c.Application().SupportedDevicesTypes),
			client.NSMethod("supportedInterfaces", c.Application().SupportedDeviceInterfaces),
			client.NSMethod("supportedModules", c.Application().SupportedDeviceModules),
			client.NSMethod("availableAddresses", c.Application().AvailableDevicesAddresses),
			client.NSMethod("list", c.Application().ConnectedDevicesList),
			client.NSMethod("count", c.Application().ConnectedDevicesCount),
			client.NSMethod("createModule", c.Application().CreateModule),
			client.NSMethod("readModule", c.Application().ReadModule),
			client.NSMethod("updateModule", c.Application().UpdateModule),
			client.NSMethod("deleteModule", c.Application().DeleteModule),
		).
		NS("aliases",
			client.NSMethod("list", c.Application().AliasesList),
			client.NSMethod("count", c.Application().AliasesCount),
		)

	c.Application().SetNotifier(s)
	go c.Application().LaunchObservers()
	c.Set(id, s, nil)

	return s
}

func (c *Boot) WebsocketHandler() *ws.Handler {
	id := "WS Handler"
	if s, ok := c.Get(id).(*ws.Handler); ok {
		return s
	}

	h := ws.NewHandler(c.WSClient())
	c.Set(id, h, nil)

	return h
}

func (c *Boot) Webserver() *webserver.Webserver {
	id := "Web Server"
	if s, ok := c.Get(id).(*webserver.Webserver); ok {
		return s
	}

	s := webserver.New(
		c.WebsocketHandler(),
		c.Config().Webserver.Listen,
	).Use(
		//middleware.BasicAuth(c.Auth()),
		middleware.Translator(),
	)

	fn := func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := s.Stop(ctx); err != nil {
			log.Printf("Error stopping web server: %s", err)
		}
	}
	c.Set(id, s, fn)

	return s
}

func (c *Boot) Logger() {
	cfg := c.Config().Logger

	if level, err := log.ParseLevel(cfg.Level); err == nil {
		log.SetLevel(level)
	}

	log.SetOutput(io.Discard)
	log.SetFormatter(cfg.Formatter(&log.TextFormatter{
		ForceColors:     true,
		DisableQuote:    true,
		FullTimestamp:   true,
		TimestampFormat: cfg.TimestampFormat,
		FieldMap:        cfg.FieldMap(),
	}))

	// Configure logs rotation
	rotor := &lumberjack.Logger{
		Filename:   cfg.Rotor.Filename,
		MaxSize:    cfg.Rotor.MaxSize,
		MaxAge:     cfg.Rotor.MaxAge,
		MaxBackups: cfg.Rotor.MaxBackups,
		LocalTime:  cfg.Rotor.LocalTime,
		Compress:   cfg.Rotor.Compress,
	}

	for _, hook := range []log.Hook{
		// Send logs with level higher than warning to stderr
		&writer.Hook{
			Writer: os.Stderr,
			LogLevels: []log.Level{
				log.PanicLevel,
				log.FatalLevel,
				log.ErrorLevel,
				log.WarnLevel,
			},
		},
		// Send info and debug logs to stdout
		&writer.Hook{
			Writer: os.Stdout,
			LogLevels: []log.Level{
				log.InfoLevel,
				log.DebugLevel,
			},
		},
		// Send all logs to file in JSON format with rotation
		lfshook.NewHook(rotor, cfg.Formatter(&log.JSONFormatter{
			TimestampFormat: cfg.TimestampFormat,
			FieldMap:        cfg.FieldMap(),
		})),
	} {
		log.AddHook(hook)
	}

	c.Set("logger", rotor, func() {
		if err := rotor.Rotate(); err != nil {
			log.Errorln("error rotating log files:", err)
		}

		if err := rotor.Close(); err != nil {
			log.Errorln("error closing log files rotator:", err)
		}
	})

	log.WithFields(log.Fields{
		"output": rotor.Filename,
		"grade":  log.GetLevel(),
	}).Info("logger was successfully configured")
}
