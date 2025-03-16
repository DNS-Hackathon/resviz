package engine

import (
	"crypto/tls"
	"os"

	"fmt"
	"log"
	"net/http"

	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/yaml.v3"

	dig "resviz/dig"
)

const (
	ContentTypeBinary = "application/octet-stream"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
)

type httpconf struct {
	LogLevel         int      `yaml:"LogLevel" json:"LogLevel"`
	TLS              string   `yaml:"TLS" json:"TLS"`
	Address          string   `yaml:"Address" json:"Address"`
	Port             string   `yaml:"Port" json:"Port"`
	Hostnames        []string `yaml:"Hostnames" json:"Hostnames"`
	AutoTlsCertdir   string   `yaml:"AutoTlsCertdir" json:"AutoTlsCertdir"`
	LocalCertfile    string   `yaml:"LocalCertfile" json:"LocalCertfile"`
	LocalKeyfile     string   `yaml:"LocalKeyfile" json:"LocalKeyfile"`
	AllowOrigins     []string `yaml:"AllowOrigins" json:"AllowOrigins"`
	AllowMethods     []string `yaml:"AllowMethods" json:"AllowMethods"`
	AllowHeaders     []string `yaml:"AllowHeaders" json:"AllowHeaders"`
	ExposeHeaders    []string `yaml:"ExposeHeaders" json:"ExposeHeaders"`
	AllowCredentials bool     `yaml:"AllowCredentials" json:"AllowCredentials"`
}

func ping(c *gin.Context) {

	outstr := "pong"
	c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))

}

func Run() {

	// http conf
	cf, err := os.ReadFile("conf.yaml")
	if err != nil {
		fmt.Printf("ReadFile error: %v\n", err)
	}
	var hc httpconf
	yaml.Unmarshal(cf, &hc)
	if err != nil {
		fmt.Printf("YAML error: %v\n", err)
	}

	if hc.LogLevel > 1 {
		fmt.Printf("Conf:%v\n", hc)

	}
	// Create a Gin router
	router := gin.Default()

	// Fixes issues with browsers diallowing POSTing (No 'Access-Control-Allow-Origin' header is present on the requested resource)
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		//AllowAllOrigins:  false,
		//AllowOrigins:     hc.AllowOrigins,
		AllowMethods:     hc.AllowMethods,
		AllowHeaders:     hc.AllowHeaders,
		ExposeHeaders:    hc.ExposeHeaders,
		AllowCredentials: hc.AllowCredentials,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(static.Serve("/", static.LocalFile("html", false)))

	router.POST("/resviz", func(c *gin.Context) {
		domain := c.PostForm("domain")
		fmt.Printf("Got domain: %s", domain)
		outstr := dig.OSdig(domain)
		c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))

	})
	router.POST("/dig/webclient", func(c *gin.Context) {
		var wq dig.WebQuery
		if err := c.ShouldBindJSON(&wq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := wq.Parse()
		fmt.Printf("Query: %v", query)

		outstr := dig.OSdig(query.Qname)
		//out := dig.Dig(query)

		// update the #terminal div
		//outstr := out.ToHTML()

		if hc.LogLevel > 2 {
			fmt.Printf("\n\n%s\n\n", outstr)
		}

		c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))
	})

	router.GET("/info/*name", func(c *gin.Context) {
		// trim any leading slash (applies when no 'name' is provided)
		name := strings.TrimLeft(c.Param("name"), "/")
		outstr := name
		c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))

	})

	router.GET("/ping", ping)

	switch hc.TLS {
	case "auto":
		hosts := strings.Join(hc.Hostnames, ",")
		if hc.LogLevel > 1 {
			fmt.Printf("\nHostname Whitelist: %s\n", hosts)
		}

		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(hosts),
			Cache:      autocert.DirCache(hc.AutoTlsCertdir),
		}

		log.Fatal(autotls.RunWithManager(router, &m))
	case "local":
		// Load Certificates
		cer, err := tls.LoadX509KeyPair(hc.LocalCertfile, hc.LocalKeyfile)
		if err != nil {
			log.Println(err)
			return
		}

		config := &tls.Config{Certificates: []tls.Certificate{cer}}

		server := &http.Server{
			Addr:      hc.Address + ":" + hc.Port,
			TLSConfig: config,
			Handler:   router,
		}

		// Start the HTTPS server
		log.Fatal(server.ListenAndServeTLS("", ""))
	default:
		server := &http.Server{
			Addr:    hc.Address + ":" + hc.Port,
			Handler: router,
		}
		log.Fatal(server.ListenAndServe())

	}
}
