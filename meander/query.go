package meander

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var APIKey string

type Place struct {
	*GoogleGeometry `json:"geometry"`
	Name            string         `json:"name"`
	Icon            string         `json:"icon"`
	Photos          []*GooglePhoto `json:"photos"`
	Vicinity        string         `json:"vicinity"`
}

type GoogleResponse struct {
	Results []*Place `json:"results"`
}

type GoogleGeometry struct {
	*GoogleLocation `json:"location"`
}

type GoogleLocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type GooglePhoto struct {
	PhotoRef string `json:"photo_reference"`
	URL      string `json:"url"`
}

func (p *Place) Public() interface{} {
	return map[string]interface{}{
		"name":     p.Name,
		"icon":     p.Icon,
		"photos":   p.Photos,
		"vicinity": p.Vicinity,
		"lat":      p.Lat,
		"lng":      p.Lng,
	}
}

type Query struct {
	Lat          float64
	Lng          float64
	Journey      []string
	Radius       int
	CostRangeStr string
}

func (q *Query) find(types string) (*GoogleResponse, error) {
	u := "https://maps.googleapis.com/maps/api/place/nearbysearch/json"
	vals := make(url.Values)
	vals.Set("location", fmt.Sprintf("%g,%g", q.Lat, q.Lng))
	vals.Set("radius", fmt.Sprintf("%d", q.Radius))
	vals.Set("type", types)
	vals.Set("key", APIKey)
	if len(q.CostRangeStr) > 0 {
		r := ParseCostRange(q.CostRangeStr)
		vals.Set("minprice", fmt.Sprintf("%d", int(r.From)-1))
		vals.Set("maxprice", fmt.Sprintf("%d", int(r.To)-1))
	}

	res, err := http.Get(u + "?" + vals.Encode())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response GoogleResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (q *Query) Run() []interface{} {
	rand.Seed(time.Now().UnixNano())
	var w sync.WaitGroup
	places := make([]interface{}, len(q.Journey))

	for i, r := range q.Journey {
		w.Add(1)
		go func(types string, i int) {
			defer w.Done()
			response, err := q.find(types)

			if err != nil {
				log.Println("????????????????????????????????????: ", err)
				return
			}
			if len(response.Results) == 0 {
				log.Println("???????????????????????????????????????: ", types)
				return
			}

			for _, result := range response.Results {
				for _, photo := range result.Photos {
					photo.URL = "https://maps.googleapis.com/maps/api/place/photo?maxwidth=1000&photoreference=" + photo.PhotoRef + "&key=" + APIKey
				}
			}

			randI := rand.Intn(len(response.Results))
			places[i] = response.Results[randI]
		}(r, i)
	}

	w.Wait()
	return places
}
