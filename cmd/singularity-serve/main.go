package main

import (
	"fmt"
	"net/http"
	"path"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/brandur/singularity"
	"github.com/joeshaw/envdecode"
)

var llog = log.WithFields(log.Fields{
	"prefix": "serve",
})

// Conf contains configuration information for the command.
type Conf struct {
	// Port is the port on which the command will serve the site over HTTP.
	Port int `env:"PORT,default=5001"`
}

func main() {
	var conf Conf
	err := envdecode.Decode(&conf)
	if err != nil {
		llog.Fatal(err)
	}

	err = singularity.CreateTargetDir()
	if err != nil {
		llog.Fatal(err)
	}

	err = serve(conf.Port)
	if err != nil {
		llog.Fatal(err)
	}
}

func serve(port int) error {
	fmt.Printf("Serving '%v' on port %v\n", path.Clean(singularity.TargetDir), port)

	handler := http.FileServer(http.Dir(singularity.TargetDir))
	return http.ListenAndServe(":"+strconv.Itoa(port), handler)
}
