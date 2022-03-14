package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	arnList = []string{
		"24234261344000020134087",
		"24234261344000020133907",
		"24234261344000020135332",
	}

	orderIDList = []string{
		"CAQACAQL1M02wTd7CWNSstiaNc2PKjG7WsN2VlVKz83JS-LekjoBY85WwsJOo7Noq001c001",
		"CAQACAQM_S6U8zvyN2rV8cWXwnQ_JrI1Yzvl-jyGqi-TAPTGLyoDai348zqrYMpAq001c001",
		"CAQACAQMZUg6bB0F48rmde2_GMQMAnF04QgzrDU3tYMUBXRvH8IDyYOubQkfIs6Yq001c001",
		"CAQACAQM6nsI27Pgput-DOZ9AqIJ-cSL7YlWUR6r61qc5bC30F4C61pQ2Yqv7bxUq001c001",
		"CAQACAQLp55TRDDjC8pMSfQW9t8JyS_Vt4hBf7fkwXgufDkz5ZIDyXl_R4hRNiWcq001c001",
		"CAQACAQNYU-gcJAXo2B3zFBmVEEipmcJ-78bjYamu6OeWeKCr_oDY6OMc78jlIAgq001c001",
		"CAQACAQLEzkBsmy-r5HbvJfJ7Xi_StOhqR9zKK3UBl_5JiG5SDYDkl8psR3926K4q001c001",
		"CAQACAQLKdqltahP74sEYtUZqh__STVzTMyyqXTsWIGbBZ39nsoDiIKptM0oroSIq001c001",
	}
)

const (
	//baseURL = "http://production-api.acq-visa-clearing.melifrontends.com"
	baseURL    = "https://internal-api.mercadopago.com/acq/staging/visa-clearing"
	xAuthToken = "YOUR_TOKEN"
)

type (
	ARN struct {
		ID    string `json:"id"`
		State string `json:"state"`
	}

	ClearingRequest struct {
		IsSafetyNet bool `json:"is_safety_net"`
	}

	Order struct {
		ID              string          `json:"id"`
		State           string          `json:"state"`
		ARN             string          `json:"arn"`
		TCRecord        []string        `json:"tc_record"`
		SettlementDates SettlementDates `json:"settlement_dates"`
		Revision        Revision        `json:"revision"`
		BatchKey        BatchKey        `json:"batch_key"`
		IsSafetyNet     bool            `json:"is_safety_net"`
		ClearingRequest ClearingRequest `json:"clearing_request"`
	}

	SettlementDates struct {
		Reconciliation string `json:"reconciliation"`
		Settlement     string `json:"settlement"`
		Value          string `json:"value"`
		Merchant       string `json:"merchant"`
		WorkingDays    int    `json:"working_days"`
		CalendarDays   int    `json:"calendar_days"`
		ValidToUTC     string `json:"valid_to_utc"`
	}

	Revision struct {
		UpdatedAt string `json:"updated_at"`
	}

	BatchKey struct {
		UniqueFileID string `json:"unique_file_id"`
		BatchNumber  int    `json:"batch_number"`
	}

	GetStateResponse struct {
		ID    string `json:"transaction_id"`
		State string `json:"state"`
	}

	Authtrx struct {
		ID         string    `json:"transaction_id"`
		ModifiedAt time.Time `json:"modified_at"`
		State      string    `json:"state"`
	}
)

/*
func main() {
   timeNow := time.Now().Add(time.Hour)
   timeLimit := time.Now().Add(365 * 24 * time.Hour)

   fmt.Printf("time now: %s, time limit: %s\n", timeNow, timeLimit)
   fmt.Println(!timeNow.Before(timeLimit))
}
*/

func main() {
	fmt.Printf("Reprocessing %d orders\n", len(orderIDList))

	var count = 0

	wg := sync.WaitGroup{}
	for _, value := range orderIDList {
		wg.Add(1)
		go executeReviewAndReprocess(&wg, value)

		count++
		if count%49 == 0 {
			time.Sleep(1 * time.Second)
		}
	}

	wg.Wait()
	fmt.Printf("%d orders reprocessed", count)
}

/*
func main() {
   urlGetByState := "https://internal-api.mercadopago.com/acq/internal/visa-authorization/v1/transactions/state/capture_requested"
   urLGetByID := "https://internal-api.mercadopago.com/acq/internal/visa-authorization/v1/transactions/"

   client := http.DefaultClient

   checkLateTransactions(client, urlGetByState, urLGetByID)
   t := time.Tick(time.Hour * 24)
   for {
      select {
      case <-t:
         checkLateTransactions(client, urlGetByState, urLGetByID)
      }
   }
}


*/
func executeReviewAndReprocess(wg *sync.WaitGroup, value string) {
	defer wg.Done()

	order := getByOrderID(value)
	fmt.Printf("ID: %s, St: %s, ARN: %s, Update: %s, FileID: %s, BatchNum: %d\n",
		value, order.State, order.ARN, order.Revision.UpdatedAt, order.BatchKey.UniqueFileID, order.BatchKey.BatchNumber)

	//order := getByARN(value)
	//fmt.Printf("ARN: %s, St: %s, ID: %s\n", value, order.State, order.ID)

	/*
	   if order.State == "presenting" || order.State == "presentation_in_review" || order.State == "sent_exception"  {
	      if order.State == "presenting" || order.State == "sent_exception" {
	         err := reviewOrder(order.ID)
	         if err != nil {
	            fmt.Println("Error review: %w", err)
	            return
	         }

	         err = reprocessOrder(order.ID)
	         if err != nil {
	            fmt.Println("Error reprocess: %w", err)
	            return
	         }
	      } else if order.State == "presentation_in_review" {
	         err := reprocessOrder(order.ID)
	         if err != nil {
	            fmt.Println("Error reprocess: %w", err)
	            return
	         }
	      }
	   }
	*/
}

func reviewOrder(orderID string) error {
	reviewURL := fmt.Sprintf("%s/maintenance/clearing/orders/%s/review", baseURL, orderID)

	err := doPost(reviewURL)
	if err != nil {
		return err
	}

	fmt.Printf("Order in review: %s\n", orderID)

	return nil
}

func reprocessOrder(orderID string) error {
	reprocessURL := fmt.Sprintf("%s/maintenance/clearing/orders/%s/reprocess", baseURL, orderID)

	err := doPost(reprocessURL)
	if err != nil {
		return err
	}

	fmt.Printf("Order in reprocess: %s\n", orderID)

	return nil
}

func getByOrderID(orderID string) Order {
	getByOrderID := fmt.Sprintf("%s/maintenance/clearing/orders/%s", baseURL, orderID)

	returnedData := doGet(getByOrderID)

	var result Order
	_ = json.Unmarshal(returnedData, &result)

	return result
}

func getByARN(arn string) ARN {
	getOrderByARN := fmt.Sprintf("%s/maintenance/clearing/orders?arn=%s", baseURL, arn)

	returnedData := doGet(getOrderByARN)

	var result ARN
	_ = json.Unmarshal(returnedData, &result)

	return result
}

func doPost(postURL string) error {
	reqURL, _ := url.Parse(postURL)

	req := &http.Request{
		Method: "PUT",
		URL:    reqURL,
		Header: map[string][]string{
			"Content-Type": {"application/json; charset=UTF-8"},
			"x-auth-token": {xAuthToken},
		},
	}

	res, err := http.DefaultClient.Do(req)

	res.Body.Close()

	return err
}

func doGet(getUrl string) []byte {
	reqURL, _ := url.Parse(getUrl)

	req := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Content-Type": {"application/json; charset=UTF-8"},
			"x-auth-token": {xAuthToken},
		},
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal("Error:", err)
	}

	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	return data
}

func checkLateTransactions(client *http.Client, urlGetByState string, urLGetByID string) {
	fmt.Println("checking late transactions")
	resp, err := client.Get(urlGetByState)
	if err != nil {
		panic(err)
	}

	resps := []GetStateResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&resps); err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	for _, r := range resps {
		wg.Add(1)
		go checkLateTrx(&wg, client, urLGetByID, r)
	}

	wg.Wait()
	fmt.Println("finished checking late transactions")
}

func checkLateTrx(wg *sync.WaitGroup, client *http.Client, urLGetByID string, r GetStateResponse) {
	defer wg.Done()
	getById, err := client.Get(urLGetByID + r.ID)
	if err != nil {
		panic(err)
	}

	trx := Authtrx{}
	if err := json.NewDecoder(getById.Body).Decode(&trx); err != nil {
		panic(err)
	}

	if trx.ModifiedAt.Before(time.Now().Add(-time.Hour * 18)) {
		fmt.Println(fmt.Sprintf("transaction is late: %s", trx.ID))

		req, _ := http.NewRequest(http.MethodPatch, urLGetByID+trx.ID+"/notify",
			&bytes.Buffer{})

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		req, _ = http.NewRequest(http.MethodPost, "https://hooks.slack.com/services/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",
			bytes.NewBufferString(fmt.Sprintf(`{"text": "transaction was late: %s notify response code: %d"}`, trx.ID, resp.StatusCode)))

		client.Do(req)
	}
}
