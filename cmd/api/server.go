package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type user struct {
	Name string `json:"name"`
	Age  string `json:"age"`
	City string `json:"city"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello Root Route")
	w.Write([]byte("Hello Root Router"))
	fmt.Println("Hello Root Router")
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {
	// teachers/{id}
	// teachers/9
	switch r.Method {
	case http.MethodGet:
		fmt.Println(r.URL.Path)
		path := strings.TrimPrefix(r.URL.Path, "/teachers/")
		userID := strings.TrimSuffix(path, "/")

		fmt.Println("The ID is:", userID)

		fmt.Println("Query Params:", r.URL.Query())
		queryParams := r.URL.Query()
		typeP := queryParams.Get("type")
		nameP := queryParams.Get("name")
		ageP := queryParams.Get("age")

		fmt.Printf("Nombre de la mascota: %v, Edad: %v, Categoria de la mascota %v\n", nameP, ageP, typeP)

		w.Write([]byte("Hello GET Method on Teachers Route"))
		// fmt.Println("Hello GET Method on Teachers Route")
	case http.MethodPost:
		w.Write([]byte("Hello POST Method on Teachers Route"))
		fmt.Println("Hello POST Method on Teachers Route")
	case http.MethodPut:
		w.Write([]byte("Hello PUT Method on Teachers Route"))
		fmt.Println("Hello PUT Method on Teachers Route")
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH Method on Teachers Route"))
		fmt.Println("Hello PATCH Method on Teachers Route")
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE Method on Teachers Route"))
		fmt.Println("Hello DELETE Method on Teachers Route")
	}

	// w.Write([]byte("Hello Teachers Router"))
	// fmt.Println("Hello Teachers Router")

}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET Method on Students Route"))
		fmt.Println("Hello GET Method on Students Route")
	case http.MethodPost:
		w.Write([]byte("Hello POST Method on Students Route"))
		fmt.Println("Hello POST Method on Students Route")
	case http.MethodPut:
		w.Write([]byte("Hello PUT Method on Students Route"))
		fmt.Println("Hello PUT Method on Students Route")
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH Method on Students Route"))
		fmt.Println("Hello PATCH Method on Students Route")
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE Method on Students Route"))
		fmt.Println("Hello DELETE Method on Students Route")
	}
	w.Write([]byte("Hello Students Router"))
	fmt.Println("Hello Students Router")
}

func execsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET Method on Execs Route"))
		fmt.Println("Hello GET Method on Execs Route")
	case http.MethodPost:
		w.Write([]byte("Hello POST Method on Execs Route"))
		fmt.Println("Hello POST Method on Execs Route")
	case http.MethodPut:
		w.Write([]byte("Hello PUT Method on Execs Route"))
		fmt.Println("Hello PUT Method on Execs Route")
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH Method on Execs Route"))
		fmt.Println("Hello PATCH Method on Execs Route")
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE Method on Execs Route"))
		fmt.Println("Hello DELETE Method on Execs Route")
	}
	w.Write([]byte("Hello Execs Router"))
	fmt.Println("Hello Execs Router")
}

func main() {

	port := ":3000"

	http.HandleFunc("/", rootHandler)

	http.HandleFunc("/teachers/", teachersHandler)

	http.HandleFunc("/students/", studentsHandler)

	http.HandleFunc("/execs/", execsHandler)

	fmt.Println("Server is running on port:", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("Error staring the server:", err)
	}

}
