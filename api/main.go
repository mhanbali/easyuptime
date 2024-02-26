// Author: Mehdi Hanbali
// Credits:
// 1. Helped with CORS:
// https://codeandlife.com/2022/04/03/golang-julienschmidt-httprouter-cors-middleware/
// 2. Helped with some structuring of GORM:
// https://github.com/cockroachdb/examples-orms/blob/master/go/gorm/server.go

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Server is an http server that handles REST requests.
type Server struct {
	db *gorm.DB
}

// NewServer creates a new instance of a Server.
func NewServer(db *gorm.DB) *Server {
	return &Server{db: db}
}

func (s *Server) RegisterRouter(router *httprouter.Router) {
	router.POST("/check/:domain", MiddleCORS(s.CheckDomain))
	router.POST("/add/:domain", MiddleCORS(s.AddDomain))
	router.GET("/view/:domain", MiddleCORS(s.ViewDomain))
	router.GET("/list", MiddleCORS(s.ListDomains))
}

// DB schema
// Keep track of the domains to check uptime
type Domain struct {
	gorm.Model
	Domain string 	`gorm:"unique" json:"domain"`
}

// DB schema
// History of uptime checked
type Uptime struct {
	gorm.Model
	Domain string 	`json:"domain"`
	Response int	`json:"response"`
	Duration int64	`json:"duration"`
}

// This API call would be used by a scheduler/cron to run on an interval
// by going through the Domains table and checking their status
func (s *Server) CheckDomain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// d will contain Uptime struct that we can use
	// to insert into the DB
	d := getDomainStatus(ps.ByName("domain"))
	s.db.Create(&Uptime{Domain: d.Domain, Response: d.Response, Duration: d.Duration})
}

func (s *Server) AddDomain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: verify if ps == a valid domain format

	// Not really necessary, but can help the frontend
	// determine if it was added or not
	t, _ := json.Marshal(map[string]int {"response": 200})
	f, _ := json.Marshal(map[string]int {"response": 500})

	result := s.db.Create(&Domain{Domain: ps.ByName("domain")})
	if result.Error != nil {
		fmt.Fprint(w, string(f))
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Fprint(w, string(t))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// View the entire uptime history of a domain
// TODO: paginate
func (s *Server) ViewDomain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var domains []Uptime

	// Since the Domain table stores the full url i.e. "http://domain.com"
	// LIKE is used with the wildcard in front of the requested domain
	s.db.Where("Domain LIKE ?", "%" + ps.ByName("domain")).Find(&domains)

	jsonData, _ := json.Marshal(domains)
	fmt.Fprint(w, string(jsonData))
}

func (s *Server) ListDomains(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {	
	var domains []Domain
	s.db.Find(&domains).Scan(&domains)
	jsonData, _ := json.Marshal(domains)

	fmt.Fprint(w, string(jsonData))
}

func main() {
	db := setupDB("easyuptime.db")

	router := httprouter.New()

	server := NewServer(db)
	server.RegisterRouter(router)

	fmt.Println("Server running...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func setupDB(f string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(f), &gorm.Config{})
	
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Domain{})
	db.AutoMigrate(&Uptime{})

	return db
}

func getDomainStatus(u string) *Uptime {
	s := time.Now()
    
	if !strings.HasPrefix(u, "http") {
		u = "http://" + u
	}

	client := http.Client{}
	req , err := http.NewRequest("GET", u, nil)
	if err != nil {
		fmt.Println("Error in HTTP request:", err)
	}

	req.Header = http.Header{
		"Content-Type": {"text/html"},
		"User-Agent": {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36"},
	}
	response , err := client.Do(req)

	if err != nil {
		fmt.Println("Error in HTTP request:", err)
	}
	defer response.Body.Close()

	e := time.Since(s).Milliseconds()
	c := response.StatusCode
	domain := &Uptime{Domain: u, Response: c, Duration: e}
	fmt.Printf("Response %v from %s: %d ms\n", c, u, e)
	
	return domain
}

func MiddleCORS(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter,
    r *http.Request, ps httprouter.Params) {
		
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next(w, r, ps)
	}
}