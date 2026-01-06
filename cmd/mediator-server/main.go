package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"mediator/logger"
	"mediator/mediatorscript"
	"mediator/mediatorsettings"
	"mediator/totp"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

// @title Mediator Back-end API
// @version 1.0
// @description UquidIT.co back-end server
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.uquidit.co/support
// @contact.email support@suquidit.co

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host
// @BasePath /v1

var (
	Version = "develop"
)

func main() {
	// CLI flags
	// version
	versionPtr := flag.Bool("version", false, "Print version number and exit.")
	flag.BoolVar(versionPtr, "v", false, "Alias of --version")

	flag.Parse()

	if *versionPtr {
		fmt.Printf("uQuidIT Mediator Server version %s\n", Version)
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		logrus.Fatalf("WRONG NUMBER OF ARGUMENTS: 1 expected, got %d: %v", flag.NArg(), flag.Args())
	}
	fmt.Printf("Reading configuration file %s\n", flag.Arg(0))
	if err := ReadConfFromFile(flag.Arg(0)); err != nil {
		logrus.Fatalf("ERROR while reading configuration file %s: %v", flag.Arg(0), err)
	}

	// init traditional logger
	if Configuration.Server.Log.Error == "" || Configuration.Server.Log.Error == "-" {
		if _, err := logger.InitAppLogger(true, logrus.WarnLevel, true, false, false, true, false, true, "", ""); err != nil {
			logrus.Fatalf("error while initializing logger: %v", err)
		}
	} else {
		if _, err := logger.InitAppLogger(true, logrus.WarnLevel, false, true, false, true, true, true, "", Configuration.Server.Log.Error); err != nil {
			logrus.Fatalf("error while initializing logger: %v", err)
		}
	}
	defer logger.CloseLogFile()

	// initialize mediatorscript package
	if err := mediatorscript.Init(Configuration.Mediatorscript.ScriptStorage); err != nil {
		logrus.Warningf("error while loading scripts for mediator list: %v", err)
	}

	// Echo instance
	e := echo.New()
	e.HideBanner = true
	e.Pre(middleware.RemoveTrailingSlash())
	e.Pre(RemoveMultipleSlash())

	logformat := "${time_rfc3339} ${remote_ip} ${method} ${path} ${status} ${latency_human} ${bytes_in} ${bytes_out}\n"
	if Configuration.Server.Log.Access == "" || Configuration.Server.Log.Access == "-" {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: logformat,
		}))
	} else {
		logfile, err := os.OpenFile(Configuration.Server.Log.Access, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logrus.Fatalf("error while opening logfile '%s': %v", Configuration.Server.Log.Access, err)
		}
		defer logfile.Close()

		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: logformat,
			Output: logfile,
		}))
	}

	if errs := mediatorsettings.Init(
		Configuration.Mediatorscript.ClientConfiguration.SettingsFile,
		Configuration.Mediatorscript.ClientConfiguration.UploadScript,
		Configuration.Mediatorscript.ClientConfiguration.DownloadScript,
	); len(errs) != 0 {
		for _, err := range errs {
			logrus.Warning(err)
		}
	}

	// Middleware
	e.Use(middleware.Recover())
	// CORS default
	// Allows requests from any origin wth GET, HEAD, PUT, POST or DELETE method.
	e.Use(middleware.CORS())

	// Routes

	// API current version: all entry point must be behind a version number
	v1 := e.Group("/v1")
	v1.Use(NoCacheHeader)

	otp := v1.Group("/otp")
	otp.Use(middleware.KeyAuth(totp.CheckKey))

	// mediator-client entry points protected by otp
	mediatorscript.AddMediatorscriptAPI(otp)

	// upload and download settings
	otp.GET("/settings", mediatorsettings.GetSettings)
	otp.POST("/settings", mediatorsettings.SetSettings)
	otp.POST("/settings/workflows", mediatorsettings.SetWorkflowSettings)

	// auth := v1.Group("/-")
	// auth.Use(echojwt.JWT([]byte(Configuration.Server.Secret)))
	// auth.GET("/settings", mediatorsettings.GetSettings)
	// auth.POST("/settings", mediatorsettings.SetSettings)

	// Start server
	listen_address := fmt.Sprintf("%s:%d", Configuration.Server.Host, Configuration.Server.Port)
	if Configuration.Server.Ssl.Enabled {
		if Configuration.Server.Ssl.Certificate == "" {
			logrus.Fatal("cannot start server using SSL: empty certificate")
		}
		if Configuration.Server.Ssl.Key == "" {
			logrus.Fatal("cannot start server using SSL: empty key")
		}
		logrus.Fatal(e.StartTLS(listen_address, Configuration.Server.Ssl.Certificate, Configuration.Server.Ssl.Key))
	} else {
		logrus.Fatal(e.Start(listen_address))
	}
}

func NoCacheHeader(next echo.HandlerFunc) echo.HandlerFunc {
	// NoCache middleware adds a `Cache-Control: no-store` header to the response.
	// cf https://developer.mozilla.org/fr/docs/Web/HTTP/Headers/Cache-Control and https://developer.mozilla.org/fr/docs/Web/HTTP/Caching
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "no-store")
		return next(c)
	}
}

// RemoveMultipleSlash returns a root level (before router) middleware which replaces
// multiple slashes from the request URI by a unique slash
//
// Usage `Echo#Pre(RemoveMultipleSlash())`
func RemoveMultipleSlash() echo.MiddlewareFunc {
	return RemoveMultipleSlashWithConfig(middleware.TrailingSlashConfig{})
}

var reMultipleSlash = regexp.MustCompile(`(?m)\/+`)

// RemoveTrailingSlashWithConfig returns a RemoveTrailingSlash middleware with
// See `RemoveTrailingSlash()`.
func RemoveMultipleSlashWithConfig(config middleware.TrailingSlashConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultTrailingSlashConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			url := req.URL
			path := url.Path
			qs := c.QueryString()
			var substitution = "/"

			path = reMultipleSlash.ReplaceAllString(path, substitution)
			uri := path
			if qs != "" {
				uri += "?" + qs
			}
			// Redirect
			if config.RedirectCode != 0 {
				return c.Redirect(config.RedirectCode, uri)
			}

			// Forward
			req.RequestURI = uri
			url.Path = path

			return next(c)
		}
	}
}
