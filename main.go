package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	_ "github.com/denisenkom/go-mssqldb"
)

type Configuracion struct {
	BD               string
	BDPath           string
	BDUsuario        string
	BDPassword       string
	SRVPuerto        int
	SRVCantDiasAtras int
}

type Item struct {
	Id_jerarquia  int       `json:"id_jerarquia"`
	Id_producto   string    `json:"id_producto"`
	Desc_producto string    `json:"desc_producto"`
	Id_sucursal   string    `json:"id_sucursal"`
	Fecha         time.Time `json:"fecha"`
	Cantidad      float32   `json:"cantidad"`
	vendedor      string
	rubro         string
}

var db *sql.DB
var err error
var config Configuracion

func leerconfig() {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	//config := Configuracion{}
	err := decoder.Decode(&config)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(config.BDPath)
}

func main() {

	leerconfig()
	db, err = sql.Open("sqlserver", fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", config.BDUsuario, config.BDPassword, config.BDPath, config.BD))

	if err != nil {
		fmt.Println(" Error open db:", err.Error())
	}

	router := mux.NewRouter()
	router.HandleFunc("/ventas", ventas).Methods("GET")

	defer db.Close()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.SRVPuerto), router))
}

func ventas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var items []Item
	//fecha := config.SRVCantDiasAtras * -1
	result, err := db.Query("select '9790213' AS id_jerarquia,cve_FEmision AS fecha,apr_codequiv as id_producto, ive_Desc as desc_producto , cve_CodPvt as id_sucursal ,ive_CantUM1 as cantidad, COALESCE(cveven_Cod, '') as vendedor, artda1_Cod as rubro from itemvta  join cabventa on cve_id = ivecve_id join articulos on iveart_codgen = art_codgen join artprov on iveart_codgen = aprart_codgen where iveart_CodGen is not null and cvetco_Cod IN ('FRE','FAE','FEL','FX','DX','NCX','CED','CEL') AND cve_FEmision >= '04/01/2021' and artda1_Cod in ('0001', '0002') order by fecha desc")

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()
	for result.Next() {
		var item Item
		err := result.Scan(&item.Id_jerarquia, &item.Fecha, &item.Id_producto, &item.Desc_producto, &item.Id_sucursal, &item.Cantidad, &item.vendedor, &item.rubro)
		if err != nil {
			panic(err.Error())
		}
		switch item.Id_sucursal {
		case "0008":
			item.Id_sucursal = "9732133"
		case "0007", "0009":
			item.Id_sucursal = "9732131"
		case "0010":
			item.Id_sucursal = "9732132"
		case "0002", "0004":
			if item.vendedor == "0012" {
				item.Id_sucursal = "9752134"
			} else {
				item.Id_sucursal = "9732134"
			}

		}
		items = append(items, item)
	}
	json.NewEncoder(w).Encode(items)
}
