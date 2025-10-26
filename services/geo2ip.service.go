package services

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/IdrisAkintobi/go-basic-crud/config"
	"github.com/IncSW/geoip2"
)

const (
	TEN_MEGABYTE     = 10 << 20
	UNKNOWN_LOCATION = "unknown location"
)

type Geo2IPService struct {
	db           *geoip2.CityReader
	dataFilePath string
}

func NewGeo2IPService() *Geo2IPService {
	cfg := config.Get()
	// Ensure geo2ip database exists
	fileInfo, err := os.Stat(cfg.DataFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("geo2ip database does not exist :\n%v", err)
		}
		log.Fatalf("error checking geo2ip database:\n%v", err)
	}

	if fileInfo.Size() < TEN_MEGABYTE {
		log.Fatalf("geo2ip database is less than 10mb. File size is %dKB", fileInfo.Size()/1024)
	}

	db, err := geoip2.NewCityReaderFromFile(cfg.DataFilePath)
	if err != nil {
		log.Fatalf("error creating geo2ip database: \n%v", err)
	}
	return &Geo2IPService{
		db:           db,
		dataFilePath: cfg.DataFilePath,
	}
}

func (g *Geo2IPService) GetIPLocation(ipAddress string) string {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return UNKNOWN_LOCATION
	}

	record, err := g.db.Lookup(ip)
	if err != nil {
		log.Printf("error getting IP location: %v", err)
		return UNKNOWN_LOCATION
	}

	city := record.City.Names["en"]
	country := record.Country.Names["en"]

	if city != "" && country != "" {
		return fmt.Sprintf("%s, %s", city, country)
	} else if city != "" {
		return city
	} else if country != "" {
		return country
	}

	return UNKNOWN_LOCATION
}
