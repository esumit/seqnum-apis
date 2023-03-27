package main

import (
	"context"
	"github.com/esumit/seqnum-apis/pkg/config"
	"github.com/esumit/seqnum-apis/pkg/httprqrs"
	"github.com/esumit/seqnum-apis/pkg/mw"
	"github.com/esumit/seqnum-apis/pkg/seqnum"
	negronilogrus "github.com/esumit/seqnum-apis/pkg/third-party/negroni-logrus"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/urfave/negroni"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var (
	sqsc config.SeqNumServerConfig
)
var c = make(chan os.Signal, 1)

var rootCmd = &cobra.Command{
	Use:   "seqnum-apis",
	Short: "Run seqnum-apis as a microservice",
	Run:   SeqnumApiService,
}

func SeqnumApiService(cmd *cobra.Command, args []string) {
	sm := seqnum.NewSeqnumManager()
	smh := seqnum.NewSeqnumApiRqHandler(sm)
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(httprqrs.NotFoundHandler)

	r.HandleFunc("/seqnum", mw.HttpRqRsMiddleware(smh.Get)).Methods("GET")

	h := cors.AllowAll().Handler(r)
	n := negroni.New()
	n.Use(negronilogrus.NewMiddlewareFromLogger(log.StandardLogger(), "web"))
	n.UseHandler(h)

	srv := &http.Server{
		Addr:         sqsc.IPAddress + ":" + sqsc.Port,
		WriteTimeout: time.Second * time.Duration(sqsc.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(sqsc.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(sqsc.IdleTimeout),
		Handler:      n,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	signal.Notify(c, os.Interrupt)

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)
	log.Infoln("seqnum service server shutdown ...")
	os.Exit(0)
}

func initConfig() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	sqsc.Port = os.Getenv("SERVER_PORT")
	sqsc.IPAddress = os.Getenv("SERVER_IP_ADDRESS")
	sqsc.WriteTimeout, _ = strconv.Atoi(os.Getenv("HTTP_WRITE_TIMEOUT"))
	sqsc.ReadTimeout, _ = strconv.Atoi(os.Getenv("HTTP_READ_TIMEOUT"))
	sqsc.IdleTimeout, _ = strconv.Atoi(os.Getenv("HTTP_IDLE_TIMEOUT"))

	log.Println("Config Applied:")
	log.Println("Port: ", sqsc.Port)
	log.Println("IPAddress: ", sqsc.IPAddress)
	log.Println("HTTP WriteTimeout: ", sqsc.WriteTimeout)
	log.Println("HTTP ReadTimeout: ", sqsc.ReadTimeout)
	log.Println("HTTP IdleTimeout: ", sqsc.IdleTimeout)

	log.Println("All configs loaded")
}

func init() {
	cobra.OnInitialize(initConfig)
}

func SeqnumService() {
	if err := rootCmd.Execute(); err != nil {
		log.Infoln(err)
		os.Exit(1)
	}
}

func main() {
	SeqnumService()
}
