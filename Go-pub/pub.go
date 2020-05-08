package main

import (
	"log"
	"fmt"	
	"github.com/gorilla/mux"
	"net/http"	
	"encoding/json"
	"flag"
	"os"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

var usageStr = `
Hola xd
`

func usage(){
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

type Caso struct {
	Nombre string `json:"nombre"`
	Departamento string `json:"departamento"`
	Edad string `json:"edad"`
	FormadeContagio string `json:"forma"`
	Estado string `json:"estado"`
}

var(
	clusterID 	string
	clientID 	string
	URL			string
	async		bool
)

func envio_datos(caso Caso) {
	cadena := "{\"nombre\":\""+caso.Nombre+"\",\"departamento\":\""+caso.Departamento+"\",\"edad\":\""+caso.Edad+"\",\"forma\":\""+caso.FormadeContagio+"\",\"estado\":\""+caso.Estado+"\"}"

	

	opts:= []nats.Option{nats.Name("Nats Publisher")}

	nc, err := nats.Connect(URL,opts...)
	if err!=nil{
		log.Fatal(err)
	}

	defer nc.Close()

	sc, err := stan.Connect("test-cluster","device-3F45",stan.NatsConn(nc))
	if err!=nil{
		log.Fatal(err)
	}
	//sc.Close()
	if !async {
		subj, msg := "lab", []byte(cadena)
		err = sc.Publish(subj,msg)
		if err != nil {
			sc.Close()
			log.Fatalf("Error en la publicacion: %v\n",err)
		}
		log.Printf("Prublicado [%s]: '%s' \n",subj,msg)
	}
	sc.Close()
}

func ingreso(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "ingreso.html")
}

func createEntrada(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var caso Caso
	_ = json.NewDecoder(r.Body).Decode(&caso)
	envio_datos(caso)
	json.NewEncoder(w).Encode(&caso)
}


var router = mux.NewRouter()

func main() {

	flag.StringVar(&URL, "s","nats://104.197.208.242:4222","")
	flag.StringVar(&URL, "server","nats://104.197.208.242:4222","")

	flag.StringVar(&clusterID,"c","test-cluster","")
	flag.StringVar(&clusterID,"cluster","test-cluster","")

	flag.StringVar(&clientID,"id","stan-pub","")
	flag.StringVar(&clientID,"clientid","stan-pub","")

	flag.BoolVar(&async,"a",false,"")
	flag.BoolVar(&async,"async",false,"")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()


	router.HandleFunc("/ingreso", ingreso).Methods("GET")
	router.HandleFunc("/ingreso", createEntrada).Methods("POST")
    http.Handle("/", router)
	fmt.Println("Servidor corriendo en 0.0.0.0:8082")
	http.ListenAndServe(":8082", nil)
}

