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

	_ "embed"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-faker/faker/v4"
)

//go:embed codes.txt
var codes string

const (
	APIURL     = "http://localhost:8090"
	ITERATIONS = 10
	MAX_SLEEP  = 10_000 * time.Millisecond
	CONTENDERS = 200
)

type SimulatorEvent int

const (
	EventReceived SimulatorEvent = iota
	RequestSent
)

func main() {
	var registrationCodes []string = strings.Split(codes, "\n")

	var wg sync.WaitGroup
	var metricsMutex sync.Mutex
	metrics := map[SimulatorEvent]int{
		EventReceived: 0,
		RequestSent:   0,
	}
	events := make(chan SimulatorEvent)

	go func() {
		for {
			event := <-events

			metricsMutex.Lock()
			metrics[event] = metrics[event] + 1
			metricsMutex.Unlock()
		}
	}()

	go func() {
		for {
			start := make(map[SimulatorEvent]int)

			metricsMutex.Lock()
			for k, v := range metrics {
				start[k] = v
			}
			metricsMutex.Unlock()

			time.Sleep(time.Second)

			metricsMutex.Lock()
			log.Printf("%d requests/s; %d events/s",
				metrics[RequestSent]-start[RequestSent],
				metrics[EventReceived]-start[EventReceived])
			metricsMutex.Unlock()
		}
	}()

	for _, code := range registrationCodes {
		runner := ContenderRunner{RegistrationCode: code}

		wg.Add(1)
		go runner.Run(ITERATIONS, &wg, events)
	}

	wg.Wait()
}

type ContenderRunner struct {
	RegistrationCode string
	contender        domain.Contender
	ticks            map[domain.ProblemID]domain.Tick
	events           chan<- SimulatorEvent
}

func (r *ContenderRunner) Run(requests int, wg *sync.WaitGroup, events chan<- SimulatorEvent) {
	defer wg.Done()

	r.ticks = make(map[domain.ProblemID]domain.Tick)
	r.events = events

	r.contender = r.GetContender()
	compClasses := r.GetCompClasses(r.contender.ContestID)

	selectedCompClass := compClasses[rand.Int()%len(compClasses)]

	switch selectedCompClass.Name {
	case "Males":
		r.contender.Name = fmt.Sprintf("%s %s (%d)", faker.FirstNameMale(), faker.LastName(), r.contender.ID)
	case "Females":
		r.contender.Name = fmt.Sprintf("%s %s (%d)", faker.FirstNameFemale(), faker.LastName(), r.contender.ID)
	}

	r.contender.PublicName = r.contender.Name
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

		time.Sleep(time.Duration(rand.Int() % int(MAX_SLEEP)))
	}
}

func (r *ContenderRunner) GetContender() domain.Contender {
	resp, err := http.Get(fmt.Sprintf("%s/codes/%s/contender", APIURL, r.RegistrationCode))
	if err != nil {
		panic(err)
	}

	r.events <- RequestSent

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
	err := json.NewEncoder(buf).Encode(contender)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/contenders/%d", APIURL, contender.ID), buf)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Regcode %s", r.RegistrationCode))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	r.events <- RequestSent

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&contender)
	if err != nil {
		panic(err)
	}

	return contender
}

func (r *ContenderRunner) GetCompClasses(contestID domain.ContestID) []domain.CompClass {
	resp, err := http.Get(fmt.Sprintf("%s/contests/%d/comp-classes", APIURL, contestID))
	if err != nil {
		panic(err)
	}

	r.events <- RequestSent

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

	r.events <- RequestSent

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

	req.Header.Set("Authorization", fmt.Sprintf("Regcode %s", r.RegistrationCode))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	r.events <- RequestSent

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

	req.Header.Set("Authorization", fmt.Sprintf("Regcode %s", r.RegistrationCode))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	r.events <- RequestSent

	defer resp.Body.Close()
}

func (r *ContenderRunner) AddTick(contenderID domain.ContenderID, tick domain.Tick) domain.Tick {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(tick)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/contenders/%d/ticks", APIURL, contenderID), buf)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Regcode %s", r.RegistrationCode))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	r.events <- RequestSent

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

	r.events <- RequestSent

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
			r.events <- EventReceived
		}
	}
}
