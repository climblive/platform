package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-faker/faker/v4"
)

const APIURL = "http://localhost:8090"

func main() {
	log.Println("Hello, World!")

	var registrationCodes []string

	for n := range 200 {
		registrationCodes = append(registrationCodes, fmt.Sprintf("ABCD%04d", n+1))
	}

	var wg sync.WaitGroup
	requests := make(chan struct{})
	events := make(chan struct{})

	stats := struct {
		requests int
		events   int
	}{}

	go func() {
		for {
			select {
			case <-requests:
				stats.requests++
			case <-events:
				stats.events++
			}
		}
	}()

	go func() {
		for {
			startRequests, startEvents := stats.requests, stats.events

			time.Sleep(time.Second)

			log.Printf("%d req/s; %d events/s", stats.requests-startRequests, stats.events-startEvents)
		}
	}()

	for _, code := range registrationCodes {
		runner := ContenderRunner{RegistrationCode: code}

		wg.Add(1)
		go runner.Run(1250, &wg, requests, events)
	}

	wg.Wait()
}

type ContenderRunner struct {
	RegistrationCode string
	contender        domain.Contender
	ticks            map[domain.ProblemID]domain.Tick
	requestsChannel  chan<- struct{}
	eventsChannel    chan<- struct{}
}

func (r *ContenderRunner) Run(requests int, wg *sync.WaitGroup, requestsChannel, eventsChannel chan<- struct{}) {
	defer wg.Done()

	r.ticks = make(map[domain.ProblemID]domain.Tick)
	r.requestsChannel = requestsChannel
	r.eventsChannel = eventsChannel

	r.contender = r.GetContender()
	compClasses := r.GetCompClasses(r.contender.ContestID)

	selectedCompClass := compClasses[rand.Int()%len(compClasses)]

	switch selectedCompClass.Name {
	case "Males":
		r.contender.Name = fmt.Sprintf("%s %s", faker.FirstNameMale(), faker.LastName())
	case "Females":
		r.contender.Name = fmt.Sprintf("%s %s", faker.FirstNameFemale(), faker.LastName())
	}

	r.contender.CompClassID = selectedCompClass.ID
	r.contender.ClubName = faker.ChineseName()

	r.UpdateContender(r.contender)

	problems := r.GetProblems(r.contender.ContestID)
	ticks := r.GetTicks(r.contender.ID)

	for tick := range slices.Values(ticks) {
		r.ticks[tick.ProblemID] = tick
	}

	go r.ReadEvents(r.contender.ID)

	for range requests {
		problem := problems[rand.Int()%len(problems)]

		tick, ok := r.ticks[problem.ID]
		if ok {
			r.DeleteTick(tick.ID)
			delete(r.ticks, problem.ID)
		} else {
			tick := domain.Tick{
				ProblemID:    problem.ID,
				AttemptsTop:  1,
				Top:          true,
				AttemptsZone: 1,
				Zone:         true,
			}

			tick = r.AddTick(r.contender.ID, tick)
			r.ticks[problem.ID] = tick
		}

		time.Sleep(time.Duration(rand.Int()%10000) * time.Millisecond)
	}
}

func (r *ContenderRunner) GetContender() domain.Contender {
	resp, err := http.Get(fmt.Sprintf("%s/codes/%s/contender", APIURL, r.RegistrationCode))
	if err != nil {
		panic(err)
	}

	r.requestsChannel <- struct{}{}

	defer resp.Body.Close()

	contender := domain.Contender{}

	err = json.NewDecoder(resp.Body).Decode(&contender)
	if err != nil {
		panic(err)
	}

	return contender
}

func (r *ContenderRunner) UpdateContender(contender domain.Contender) domain.Contender {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(contender)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/contenders/%d", APIURL, contender.ID), buf)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Regcode %s", r.RegistrationCode))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	r.requestsChannel <- struct{}{}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&contender)
	if err != nil {
		panic(err)
	}

	return contender
}

func (r *ContenderRunner) GetCompClasses(contestID domain.ContestID) []domain.CompClass {
	resp, err := http.Get(fmt.Sprintf("%s/contests/%d/compClasses", APIURL, contestID))
	if err != nil {
		panic(err)
	}

	r.requestsChannel <- struct{}{}

	defer resp.Body.Close()

	compClasses := []domain.CompClass{}

	err = json.NewDecoder(resp.Body).Decode(&compClasses)
	if err != nil {
		panic(err)
	}

	return compClasses
}

func (r *ContenderRunner) GetProblems(contestID domain.ContestID) []domain.Problem {
	resp, err := http.Get(fmt.Sprintf("%s/contests/%d/problems", APIURL, contestID))
	if err != nil {
		panic(err)
	}

	r.requestsChannel <- struct{}{}

	defer resp.Body.Close()

	problems := []domain.Problem{}

	err = json.NewDecoder(resp.Body).Decode(&problems)
	if err != nil {
		panic(err)
	}

	return problems
}

func (r *ContenderRunner) GetTicks(contenderID domain.ContenderID) []domain.Tick {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/contenders/%d/ticks", APIURL, contenderID), nil)
	if err != nil {
		panic(err)
	}

	r.requestsChannel <- struct{}{}

	req.Header.Set("Authorization", fmt.Sprintf("Regcode %s", r.RegistrationCode))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	ticks := []domain.Tick{}

	err = json.NewDecoder(resp.Body).Decode(&ticks)
	if err != nil {
		panic(err)
	}

	return ticks
}

func (r *ContenderRunner) DeleteTick(tickID domain.TickID) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/ticks/%d", APIURL, tickID), nil)
	if err != nil {
		panic(err)
	}

	r.requestsChannel <- struct{}{}

	req.Header.Set("Authorization", fmt.Sprintf("Regcode %s", r.RegistrationCode))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
}

func (r *ContenderRunner) AddTick(contenderID domain.ContenderID, tick domain.Tick) domain.Tick {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(tick)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/contenders/%d/ticks", APIURL, contenderID), buf)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Regcode %s", r.RegistrationCode))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	r.requestsChannel <- struct{}{}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tick)
	if err != nil {
		panic(err)
	}

	return tick
}

func (r *ContenderRunner) ReadEvents(contenderID domain.ContenderID) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/contenders/%d/events", APIURL, contenderID), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Regcode %s", r.RegistrationCode))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	r.requestsChannel <- struct{}{}

	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			panic(err)
		}

		if strings.HasPrefix(line, "event:") {
			r.eventsChannel <- struct{}{}
		}
	}
}
