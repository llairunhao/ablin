package main

import (
	"database/sql"
	"ggfly/controller"
	"ggfly/logger"
	"ggfly/service"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"net/http"
)

func main() {

	done := make(chan error, 2)
	stop := make(chan struct{})
	//go func() {
	//	done <- serveDebug(stop)
	//}()
	go func() {
		done <- serveApp(stop)
	}()

	go func() {
		db, err := connectMySql()
		if err != nil {
			done <- err
			return
		}
		service.SetDB(db)
	}()

	var stopped bool
	for i := 0; i < cap(done); i++ {
		if err := <-done; err != nil {
			logger.ErrorLog("error: %v", err)
		}
		if !stopped {
			stopped = true
			close(stop)
		}
	}
}

func serve(addr string, handler http.Handler, stop <-chan struct{}) error {
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		<-stop // wait for stop signal
		e := s.Shutdown(context.Background())
		if e != nil {
			logrus.Error(e)
		}
	}()

	return s.ListenAndServe()
}

func serveApp(stop <-chan struct{}) error {
	return serve("0.0.0.0:8080", controller.ServeRouter(), stop)
}

func serveDebug(stop <-chan struct{}) error {
	return serve("127.0.0.1:8001", controller.DebugRouter(), stop)
}

func connectMySql() (*sql.DB, error) {
	return sql.Open("mysql", "root:12345678@/ggfly")
}

//func checkErr(err error)  {
//	if err != nil {
//		panic(err)
//	}
//}
//
