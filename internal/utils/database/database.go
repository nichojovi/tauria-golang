package database

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nichojovi/tauria-test/cmd/config"
)

//Db object
var (
	Master   *DB
	Slave    *DB
	dbTicker *time.Ticker
)

type (
	//DB configuration
	DB struct {
		DBConnection  *sqlx.DB
		DBString      string
		RetryInterval int
		MaxIdleConn   int
		MaxConn       int
		doneChannel   chan bool
	}

	Store struct {
		Master *sqlx.DB
		Slave  *sqlx.DB
	}
)

func (s *Store) GetMaster() *sqlx.DB {
	return s.Master
}

func (s *Store) GetSlave() *sqlx.DB {
	return s.Slave
}

func New(cfg config.MainConfig, dbDriver string) *Store {
	masterDSN := cfg.DBConfig.MasterDSN
	slaveDSN := cfg.DBConfig.SlaveDSN

	Master = &DB{
		DBString:      masterDSN,
		RetryInterval: cfg.DBConfig.RetryInterval,
		MaxIdleConn:   cfg.DBConfig.MaxIdleConn,
		MaxConn:       cfg.DBConfig.MaxConn,
		doneChannel:   make(chan bool),
	}

	err := Master.ConnectAndMonitor(dbDriver)
	if err != nil {
		log.Fatal("Could not initiate Master DB connection: " + err.Error())
		return &Store{}
	}
	Slave = &DB{
		DBString:      slaveDSN,
		RetryInterval: cfg.DBConfig.RetryInterval,
		MaxIdleConn:   cfg.DBConfig.MaxIdleConn,
		MaxConn:       cfg.DBConfig.MaxConn,
		doneChannel:   make(chan bool),
	}
	err = Slave.ConnectAndMonitor(dbDriver)
	if err != nil {
		log.Fatal("Could not initiate Slave DB connection: " + err.Error())
		return &Store{}
	}

	dbTicker = time.NewTicker(time.Second * 2)

	return &Store{Master: Master.DBConnection, Slave: Slave.DBConnection}
}

// Connect to database
func (d *DB) Connect(driver string) error {
	var db *sqlx.DB
	var err error
	db, err = sqlx.Open(driver, d.DBString)

	if err != nil {
		log.Println("[Error]: DB open connection error", err.Error())
	} else {
		d.DBConnection = db
		err = db.Ping()
		if err != nil {
			log.Println("[Error]: DB connection error", err.Error())
		}
		return err
	}

	db.SetMaxOpenConns(d.MaxConn)
	db.SetMaxIdleConns(d.MaxIdleConn)

	return err
}

// ConnectAndMonitor to database
func (d *DB) ConnectAndMonitor(driver string) error {
	err := d.Connect(driver)

	if err != nil {
		log.Printf("Not connected to database %s, trying", d.DBString)
		return err
	}

	ticker := time.NewTicker(time.Duration(d.RetryInterval) * time.Second)
	go func() error {
		for {
			select {
			case <-ticker.C:
				if d.DBConnection == nil {
					d.Connect(driver)
				} else {
					err := d.DBConnection.Ping()
					if err != nil {
						log.Println("[Error]: DB reconnect error", err.Error())
						return err
					}
				}
			case <-d.doneChannel:
				return nil
			}
		}
	}()
	return nil
}
