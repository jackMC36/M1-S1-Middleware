package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go"
	"github.com/zhashkevych/scheduler"

	"middleware/example/internal/models"
	"github.com/joho/godotenv"

)

var jsc nats.JetStreamContext
var nc *nats.Conn

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	initStream()

	ctx := context.Background()
	sc := scheduler.NewScheduler()

	sc.Add(ctx, funcTest, time.Second*2)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	sc.Stop()
}

func funcTest(ctx context.Context) {
    configURL := os.Getenv("CONFIG_API_BASE_URL") 
    if configURL == "" {
        log.Printf("CONFIG_API_BASE_URL is not set")
        return
    }

    agendas, err := fetchAgendas(configURL)
    if err != nil {
        log.Printf("Error fetching agendas: %v", err)
        return
    }

    for _, a := range agendas {
        if a.Id == nil {
            log.Printf("Skipping agenda %q: id is nil", a.Name)
            continue
        }
        if strings.TrimSpace(a.UcaId) == "" {
            log.Printf("Skipping agenda %q (%s): ucaid is empty", a.Name, a.Id.String())
            continue
        }

        // This is the next step: fetch calendar for each agenda
        if err := fetchAndProcessCalendar(a.UcaId, *a.Id); err != nil {
            log.Printf("Error processing agenda %q (ucaid=%s): %v", a.Name, a.UcaId, err)
        }
    }
}


func fetchAgendas(configAPIBaseURL string) ([]models.Agenda, error) {
    url := strings.TrimRight(configAPIBaseURL, "/") + "/agendas/"

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        b, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("GET %s -> %d: %s", url, resp.StatusCode, string(b))
    }

	var agendas []models.Agenda
    if err := json.NewDecoder(resp.Body).Decode(&agendas); err != nil {
        return nil, err
    }
    return agendas, nil
}


func fetchAndProcessCalendar(ucaID string, agendaID uuid.UUID) error {
	rawDate := "20251228T152000Z"

	// 2006 = année ; 01 = mois ; 02 = jour ; 15 = heure ; 04 = minute ; 05 = seconde
	d, _ := time.Parse("20060102T150405Z", rawDate)

	if d.Before(time.Now()) {
		fmt.Println("Avant !")
	} else {
		fmt.Println("Après !")
	}

	fmt.Println(d)

	url := fmt.Sprintf("https://edt.uca.fr/jsp/custom/modules/plannings/anonymous_cal.jsp?resources=%s&projectId=3&calType=ical&nbWeeks=8&displayConfigId=128", ucaID)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	lines := strings.Split(string(body), "\n")

	currentlyParsing := false
	tmpObj := map[string]interface{}{}

	for _, line := range lines {
		if strings.HasPrefix(line, "BEGIN:VEVENT") {
			currentlyParsing = true
		} else {
			if currentlyParsing {
				if strings.HasPrefix(line, "END:VEVENT") {
					fmt.Println(tmpObj)
					err := publishEvent(tmpObj)
					if err != nil {
						log.Printf("Error publishing event: %v", err)
					}
					tmpObj = map[string]interface{}{}
					currentlyParsing = false
				} else {
					if strings.HasPrefix(line, "SUMMARY:") {
						// Attention, le dernier caractère est un "carriage return" (\r). On le supprime sinon ça fait échouer toute notre logique.
						tmpObj["summary"] = strings.Replace(strings.Replace(line, "SUMMARY:", "", 1), "\r", "", 1)
					}
					if strings.HasPrefix(line, "DTSTART:") {
						tmpObj["start"], _ = time.Parse("20060102T150405Z", strings.Replace(strings.Replace(line, "DTSTART:", "", 1), "\r", "", 1))
					}
				}
			} else {
				continue
			}
		}

	}

	return nil
}

func publishEvent(event map[string]interface{}) error {
	messageBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	pubAckFuture, err := jsc.PublishAsync("EVENTS.create", messageBytes)
	if err != nil {
		return err
	}
	// Pour info, les channels en Go permettent de lier les go routines (threads) entre elles : https://gobyexample.com/channels
	select {
	case <-pubAckFuture.Ok():
		return nil
	case <-pubAckFuture.Err():
		return errors.New(string(pubAckFuture.Msg().Data))
	}
}

func initStream() {
	var err error

	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	jsc, err = nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	_, err = jsc.AddStream(&nats.StreamConfig{
		Name:     "EVENTS",             // nom du stream
		Subjects: []string{"EVENTS.>"}, // tous les sujets sont sous le format "EVENTS.*"
	})
	if err != nil {
		log.Fatal(err)
	}
}
