package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/namedotcom/go/v4/namecom"
)

type Config struct {
	Domain   string        `env:"DOMAIN,required"`
	Host     string        `env:"HOST"`
	User     string        `env:"USERNAME,required"`
	Token    string        `env:"TOKEN,required,unset"`
	Interval time.Duration `env:"INTERVAL" envDefault:"5m"`
}

func main() {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln(err)
	}

	client := namecom.New(cfg.User, cfg.Token)

	for {
		if err := process(cfg, client); err != nil {
			log.Println(err)
		}
		time.Sleep(cfg.Interval)
	}
}

func process(cfg Config, client *namecom.NameCom) error {
	ip, err := getPublicIP()
	if err != nil {
		return fmt.Errorf("failed to get public IP: %w", err)
	}

	if err := updateRecord(cfg, client, ip); err != nil {
		return fmt.Errorf("failed to update dns record: %w", err)
	}

	return nil
}

func updateRecord(cfg Config, client *namecom.NameCom, publicIP string) error {
	listRecordsResponse, err := client.ListRecords(&namecom.ListRecordsRequest{
		DomainName: cfg.Domain,
	})
	if err != nil {
		return fmt.Errorf("failed to list records: %w", err)
	}

	for _, record := range listRecordsResponse.Records {
		if record.Type == "A" && record.Host == cfg.Host {
			if record.Answer == publicIP {
				log.Println("IP is already set")
				return nil
			}

			log.Printf("Updating %s to %s", record.Answer, publicIP)

			record.Answer = publicIP
			_, err := client.UpdateRecord(record)
			if err != nil {
				return fmt.Errorf("failed to update record: %w", err)
			}
			return nil
		}
	}

	return createRecord(cfg, client, publicIP)
}

func createRecord(cfg Config, client *namecom.NameCom, publicIP string) error {
	log.Printf("Creating A record with IP: %s", publicIP)
	_, err := client.CreateRecord(&namecom.Record{
		DomainName: cfg.Domain,
		Host:       cfg.Host,
		Type:       "A",
		Answer:     publicIP,
		TTL:        300,
	})
	if err != nil {
		return fmt.Errorf("failed to create record: %w", err)
	}
	return nil
}

func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ipBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ip := string(ipBytes)

	if net.ParseIP(ip) == nil {
		return "", fmt.Errorf("invalid IP is detected: %s", ip)
	}

	return ip, nil
}
