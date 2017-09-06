package utils

import (
	"database/sql"
	"sync"
	"time"

	"os"

	"github.com/MaxBosse/MakaBot/log"
	"github.com/bwmarrin/discordgo"
	client "github.com/influxdata/influxdb/client/v2"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func RemoveDuplicateMembers(list *[]*discordgo.Member) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *list {
		if !found[x.User.ID] {
			found[x.User.ID] = true
			(*list)[j] = (*list)[i]
			j++
		}
	}
	*list = (*list)[:j]
}

// INFLUXDB CONNECTION
type InfluxDB struct {
	DBConnection     *sql.DB
	influxDBHost     string
	influxDBDatabase string
	influxDBUser     string
	influxDBPassword string
	batchPoints      client.BatchPoints
	client           client.Client
	batchTicker      *time.Ticker
	appName          string
	version          string
	pointMutex       sync.Mutex
}

func (iDB *InfluxDB) Reconnect() error {
	var err error

	iDB.client.Close()

	iDB.client, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     iDB.influxDBHost,
		Username: iDB.influxDBUser,
		Password: iDB.influxDBPassword,
	})
	if err != nil {
		return err
	}

	// Create a new point batch
	iDB.batchPoints, err = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  iDB.influxDBDatabase,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	return nil
}

// New will create a database connection and return the sql.DB
func (iDB *InfluxDB) New(influxDBHost, influxDBDatabase, influxDBUser, influxDBPassword, appName, version string) error {
	var err error

	iDB.influxDBHost = influxDBHost
	iDB.influxDBDatabase = influxDBDatabase
	iDB.influxDBUser = influxDBUser
	iDB.influxDBPassword = influxDBPassword
	iDB.appName = appName
	iDB.version = version

	iDB.client, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     iDB.influxDBHost,
		Username: iDB.influxDBUser,
		Password: iDB.influxDBPassword,
	})
	if err != nil {
		return err
	}

	// Create a new point batch
	iDB.batchPoints, err = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  iDB.influxDBDatabase,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	// Start regular sending every 10 seconds
	iDB.batchTicker = time.NewTicker(time.Second * 10)
	go func() {
		for range iDB.batchTicker.C {
			iDB.Flush()
		}
	}()

	return nil
}

func (iDB *InfluxDB) AddMetric(name string, tags map[string]string, fields map[string]interface{}) error {
	hostname, _ := os.Hostname()
	tags["hostname"] = hostname
	tags["app"] = iDB.appName
	tags["version"] = iDB.version
	pt, err := client.NewPoint(name, tags, fields, time.Now())
	if err != nil {
		return err
	}
	iDB.pointMutex.Lock()
	iDB.batchPoints.AddPoint(pt)
	iDB.pointMutex.Unlock()

	return nil
}

func (iDB *InfluxDB) Flush() {
	iDB.pointMutex.Lock()
	if err := iDB.client.Write(iDB.batchPoints); err != nil {
		log.Errorln(err)
	}

	iDB.batchPoints, _ = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  iDB.influxDBDatabase,
		Precision: "s",
	})
	iDB.pointMutex.Unlock()
}

func (iDB *InfluxDB) Stop() {
	iDB.batchTicker.Stop()
}
