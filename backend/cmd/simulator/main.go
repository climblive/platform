package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	_ "embed"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-faker/faker/v4"
)

type Config struct {
	APIUrl            string   `json:"apiUrl"`
	RegistrationCodes []string `json:"registrationCodes"`
	Iterations        int      `json:"iterations"`
	MaxSleep          int      `json:"maxSleep"`
}

type SimulatorEvent int

const (
	EventReceived SimulatorEvent = iota
	RequestSent
)

func main() {
	configPath := flag.String("config", "config.json", "path to simulator config JSON file")
	flag.Parse()

	cfg := Config{}

	f, err := os.Open(*configPath)
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}
	defer func() {
		_ = f.Close()
	}()

	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}

	apiURL := cfg.APIUrl
	maxSleep := time.Duration(cfg.MaxSleep) * time.Millisecond
	registrationCodes := cfg.RegistrationCodes

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
		runner := ContenderRunner{
			RegistrationCode: code,
			apiURL:           apiURL,
			maxSleep:         maxSleep,
			contender:        domain.Contender{},
			ticks:            nil,
			events:           nil,
		}

		wg.Add(1)
		go runner.Run(cfg.Iterations, &wg, events)
	}

	wg.Wait()
}

type tickOutcome int

const (
	outcomeFlash tickOutcome = iota
	outcomeTop
	outcomeZoneBoth
	outcomeZone1Only
)

func pickWeightedProblem(problems []domain.Problem) domain.Problem {
	weights := make([]float64, len(problems))
	var total float64

	for i, p := range problems {
		pts := p.PointsTop
		if pts <= 0 {
			pts = 1
		}
		w := 1.0 / float64(pts*pts)
		weights[i] = w
		total += w
	}

	r := rand.Float64() * total
	for i, w := range weights {
		r -= w
		if r <= 0 {
			return problems[i]
		}
	}

	return problems[len(problems)-1]
}

func buildTick(problem domain.Problem) domain.Tick {
	outcomes := []tickOutcome{outcomeFlash, outcomeTop}
	if problem.Zone1Enabled && problem.Zone2Enabled {
		outcomes = append(outcomes, outcomeZoneBoth)
	}
	if problem.Zone1Enabled {
		outcomes = append(outcomes, outcomeZone1Only)
	}

	outcome := outcomes[rand.Int()%len(outcomes)]

	tick := domain.Tick{
		ProblemID: problem.ID,
	}

	attempts := 1 + rand.Int()%4
	tick.AttemptsTop = attempts
	tick.AttemptsZone1 = attempts
	tick.AttemptsZone2 = attempts

	switch outcome {
	case outcomeFlash:
		tick.Top = true
		tick.AttemptsTop = 1
		tick.Zone1 = true
		tick.AttemptsZone1 = 1
		tick.Zone2 = true
		tick.AttemptsZone2 = 1
	case outcomeTop:
		tick.Top = true
		tick.Zone1 = true
		tick.Zone2 = true
	case outcomeZoneBoth:
		tick.Zone1 = true
		tick.Zone2 = true
	case outcomeZone1Only:
		tick.Zone1 = true
	}

	return tick
}

type ContenderRunner struct {
	RegistrationCode string
	apiURL           string
	maxSleep         time.Duration
	contender        domain.Contender
	ticks            map[domain.ProblemID]domain.Tick
	events           chan<- SimulatorEvent
}

func mustStatus(resp *http.Response, expected int) {
	if resp.StatusCode != expected {
		panic(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
}

func (r *ContenderRunner) Run(requests int, wg *sync.WaitGroup, events chan<- SimulatorEvent) {
	defer wg.Done()

	r.ticks = make(map[domain.ProblemID]domain.Tick)
	r.events = events

	r.sleep()
	r.contender = r.GetContender()

	r.sleep()
	compClasses := r.GetCompClasses(r.contender.ContestID)

	patch := domain.ContenderPatch{}

	selectedCompClass := compClasses[rand.Int()%len(compClasses)]

	switch selectedCompClass.Name {
	case "Males":
		patch.Name = domain.NewPatch(fmt.Sprintf("%s %s", faker.FirstNameMale(), faker.LastName()))
	case "Females":
		fallthrough
	default:
		patch.Name = domain.NewPatch(fmt.Sprintf("%s %s", faker.FirstNameFemale(), faker.LastName()))
	}

	patch.CompClassID = domain.NewPatch(selectedCompClass.ID)

	r.sleep()
	r.PatchContender(r.contender.ID, patch)

	r.sleep()
	problems := r.GetProblems(r.contender.ContestID)

	r.sleep()
	ticks := r.GetTicks(r.contender.ID)

	for tick := range slices.Values(ticks) {
		r.ticks[tick.ProblemID] = tick
	}

	go r.ReadEvents(r.contender.ID)

	for range requests {
		r.sleep()

		problem := pickWeightedProblem(problems)

		tick, ok := r.ticks[problem.ID]
		if ok {
			r.DeleteTick(tick.ID)
			delete(r.ticks, problem.ID)
		} else {
			tick := buildTick(problem)

			tick = r.AddTick(r.contender.ID, tick)
			r.ticks[problem.ID] = tick
		}
	}
}

func (r *ContenderRunner) sleep() {
	time.Sleep(time.Duration(rand.Int() % int(r.maxSleep)))
}

func (r *ContenderRunner) GetContender() domain.Contender {
	resp, err := http.Get(fmt.Sprintf("%s/codes/%s/contender", r.apiURL, r.RegistrationCode))
	if err != nil {
		panic(err)
	}

	r.events <- RequestSent

	defer func() { _ = resp.Body.Close() }()

	mustStatus(resp, http.StatusOK)

	contender := domain.Contender{}

	err = json.NewDecoder(resp.Body).Decode(&contender)
	if err != nil {
		panic(err)
	}

	return contender
}

func (r *ContenderRunner) PatchContender(contenderID domain.ContenderID, patch domain.ContenderPatch) domain.Contender {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(patch)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/contenders/%d", r.apiURL, contenderID), buf)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Regcode %s", r.RegistrationCode))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	r.events <- RequestSent

	defer func() { _ = resp.Body.Close() }()

	mustStatus(resp, http.StatusOK)

	var contender domain.Contender

	err = json.NewDecoder(resp.Body).Decode(&contender)
	if err != nil {
		panic(err)
	}

	return contender
}

func (r *ContenderRunner) GetCompClasses(contestID domain.ContestID) []domain.CompClass {
	resp, err := http.Get(fmt.Sprintf("%s/contests/%d/comp-classes", r.apiURL, contestID))
	if err != nil {
		panic(err)
	}

	r.events <- RequestSent

	defer func() { _ = resp.Body.Close() }()

	mustStatus(resp, http.StatusOK)

	compClasses := []domain.CompClass{}

	err = json.NewDecoder(resp.Body).Decode(&compClasses)
	if err != nil {
		panic(err)
	}

	return compClasses
}

func (r *ContenderRunner) GetProblems(contestID domain.ContestID) []domain.Problem {
	resp, err := http.Get(fmt.Sprintf("%s/contests/%d/problems", r.apiURL, contestID))
	if err != nil {
		panic(err)
	}

	r.events <- RequestSent

	defer func() { _ = resp.Body.Close() }()

	mustStatus(resp, http.StatusOK)

	problems := []domain.Problem{}

	err = json.NewDecoder(resp.Body).Decode(&problems)
	if err != nil {
		panic(err)
	}

	return problems
}

func (r *ContenderRunner) GetTicks(contenderID domain.ContenderID) []domain.Tick {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/contenders/%d/ticks", r.apiURL, contenderID), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Regcode %s", r.RegistrationCode))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	r.events <- RequestSent

	defer func() { _ = resp.Body.Close() }()

	mustStatus(resp, http.StatusOK)

	ticks := []domain.Tick{}

	err = json.NewDecoder(resp.Body).Decode(&ticks)
	if err != nil {
		panic(err)
	}

	return ticks
}

func (r *ContenderRunner) DeleteTick(tickID domain.TickID) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/ticks/%d", r.apiURL, tickID), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Regcode %s", r.RegistrationCode))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	r.events <- RequestSent

	defer func() { _ = resp.Body.Close() }()

	mustStatus(resp, http.StatusNoContent)
}

func (r *ContenderRunner) AddTick(contenderID domain.ContenderID, tick domain.Tick) domain.Tick {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(tick)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/contenders/%d/ticks", r.apiURL, contenderID), buf)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Regcode %s", r.RegistrationCode))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	r.events <- RequestSent

	defer func() { _ = resp.Body.Close() }()

	mustStatus(resp, http.StatusCreated)

	err = json.NewDecoder(resp.Body).Decode(&tick)
	if err != nil {
		panic(err)
	}

	return tick
}

func (r *ContenderRunner) ReadEvents(contenderID domain.ContenderID) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/contenders/%d/events", r.apiURL, contenderID), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Regcode %s", r.RegistrationCode))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	r.events <- RequestSent

	defer func() { _ = resp.Body.Close() }()

	mustStatus(resp, http.StatusOK)

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
