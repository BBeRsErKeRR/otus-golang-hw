package integration_test

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/goccy/go-json"
	. "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/require"
)

var _ = Describe("Calendar HTTP", func() {
	now := time.Now()

	tr := &http.Transport{
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
	}

	client := http.Client{
		Transport: tr,
	}

	Describe("CreateEvent Error", func() {
		var checkError func(payload storage.EventDTO)
		BeforeEach(func() {
			checkError = func(payload storage.EventDTO) {
				data, err := json.Marshal(&payload)
				require.NoError(GinkgoT(), err)
				resp, err := client.Post(rootHTTPURL+"/v1/event", "application/json", bytes.NewReader(data))
				require.NoError(GinkgoT(), err)
				defer resp.Body.Close()
				var response struct {
					Error string `json:"error"`
				}
				require.Equal(GinkgoT(), http.StatusBadRequest, resp.StatusCode)
				err = json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(GinkgoT(), err)
				require.NotEmpty(GinkgoT(), response.Error)
			}
		})
		It("add bad date event", func() {
			badEvent := storage.EventDTO{
				Title:      gofakeit.Hobby(),
				Date:       now.Add(1 * time.Minute),
				EndDate:    now,
				Desc:       gofakeit.Phrase(),
				UserID:     gofakeit.UUID(),
				RemindDate: now.Add(1 * time.Minute),
			}
			checkError(badEvent)
		})
		It("add bad remind date event", func() {
			badEvent := storage.EventDTO{
				Title:      gofakeit.Hobby(),
				Date:       now.Add(1 * time.Minute),
				EndDate:    now.Add(1 * time.Minute),
				Desc:       gofakeit.Phrase(),
				UserID:     gofakeit.UUID(),
				RemindDate: now,
			}
			checkError(badEvent)
		})
		It("add bad userId event", func() {
			badEvent := storage.EventDTO{
				Title:      gofakeit.Hobby(),
				Date:       now,
				EndDate:    now.Add(2 * time.Minute),
				Desc:       gofakeit.Phrase(),
				UserID:     "bad",
				RemindDate: now.Add(1 * time.Minute),
			}
			checkError(badEvent)
		})
	})
	Describe("Id Error", func() {
		It("del bad id", func() {
			req, err := http.NewRequest(http.MethodDelete, rootHTTPURL+"/v1/event/222", nil)
			require.NoError(GinkgoT(), err)
			resp, err := client.Do(req)
			require.NoError(GinkgoT(), err)
			defer resp.Body.Close()
			var response struct {
				Error string `json:"error"`
			}
			require.Equal(GinkgoT(), http.StatusBadRequest, resp.StatusCode)
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(GinkgoT(), err)
			require.Equal(GinkgoT(), "bad event id", response.Error)
		})

		It("update bad args", func() {
			req, err := http.NewRequest(http.MethodPut, rootHTTPURL+"/v1/event/222", nil)
			require.NoError(GinkgoT(), err)
			resp, err := client.Do(req)
			require.NoError(GinkgoT(), err)
			defer resp.Body.Close()
			var response struct {
				Error string `json:"error"`
			}
			require.Equal(GinkgoT(), http.StatusBadRequest, resp.StatusCode)
			// b, _ := io.ReadAll(resp.Body)
			// fmt.Println(string(b))
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(GinkgoT(), err)
			require.Equal(GinkgoT(), "bad args", response.Error)
			fmt.Println(response)
		})
		It("update bad id", func() {
			req, err := http.NewRequest(http.MethodPut, rootHTTPURL+"/v1/event/222", bytes.NewBuffer([]byte("{}")))
			require.NoError(GinkgoT(), err)
			resp, err := client.Do(req)
			require.NoError(GinkgoT(), err)
			defer resp.Body.Close()
			var response struct {
				Error string `json:"error"`
			}
			require.Equal(GinkgoT(), http.StatusBadRequest, resp.StatusCode)
			// b, _ := io.ReadAll(resp.Body)
			// fmt.Println(string(b))
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(GinkgoT(), err)
			require.Equal(GinkgoT(), "bad event id", response.Error)
			fmt.Println(response)
		})
	})
})
