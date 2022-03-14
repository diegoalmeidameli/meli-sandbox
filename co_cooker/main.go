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
		"CAQACAQNxrv7PegjQ4ISmvSuzBAssGu6frBoZrnzmR7AjiHAT84DgRxnPrHlcvosq001c001",
		"CAQACAQMxQSLfRi1kQEZHPS-sCh2vDqK7tbXYyEX7-fd1nFGbhYBA-djftbV6pVMq001c001",
		"CAQACAQPy8sbtuQLwDozF2vJrTxAbQsL7-ULbxrsqeLXurmhpYoAOeNvt-ajGfnQq001c001",
		"CAQACAQNyURCfVD4ABED53cCzZLJn1Cffhlu3K12pANvO-GoIgYAEALefhocrl-kq001c001",
		"CAQACAQPHh5TdgIiXcK6Gs7QN-nfHriNAev3t0cbyF7hHCnKCWYBwF-3depARAr8q001c001",
		"CAQACAQNihlpZkRBzMI1lsX5MxsdhnDK4md7AQ73q3EY9jjuQoIAw3MBZmQSpMGAq001c001",
		"CAQACAQOBh6vG7nDvRPvx-V-iEvTWs2dkOM9DMMoYv3R9oz3MroBEv0PGOP9iUmUq001c001",
		"CAQACAQMXnQuxyJQieLnLjMH9o-84HeLb4ycBZMjMtPCG7juX44B4tAGx49MZCJkq001c001",
		"CAQACAQMy2Z5F7kO8qDLQKDI0vu1usa9kWWzprsIQ1I4bgbnYEYCo1OlFWcDsJfwq001c001",
		"CAQACAQNrf3hv2aT_NpolOKti1LUPnjCPRrxwvK-qpJuqH8kxiYA2pHBvRsnWBysq001c001",
		"CAQACAQOSZAuahLdF0vEJ3w-JmAevevlVSgz8x2CRwubkmDaqrYDSwvyaStd8abEq002c001",
		"CAQACAQPJhKggEYFvlijfE4piab1DyLpPp02AQVkFR9ze72tKZ4CWR4AgpynQOYIq001c001",
		"CAQACAQOLj7uNz8ShBlh87JJTKeV3xs1bL_ISw0dQd4WZly2VpYAGdxKNLxtW9ugq001c001",
		"CAQACAQOhO-zDRERA4PLp8bGbsb83Q2SIREXtyHsVdHVHONHUaIDgdO3DRKybSFYq001c001",
		"CAQACAQMZT6gtvbcLCvKj22oQTUiUEpbS8w9G7Sw59vwoMOdxoIAK9kYt8zJyh2Mq001c001",
		"CAQACAQO4TpyXE6gWrF2Qh2s7JYlUA-twvpJX6ezN32o0DVuBK4Cs31eXvvWkqyEq001c001",
		"CAQACAQMdeW12lNym6pZkmIxs_c3FISaEIAJK9KKUjiPWU8QS8oDqjkp2IHwZ5Dkq001c001",
		"CAQACAQNgaV49YLD68pTBecs1vu1gfNGKylpcG2hcfAlGhQxRw4DyfFw9yiRb78Qq001c001",
		"CAQACAQPlzlhIM-kW6v9y22Sa8yYM8zD7T6fBZFyTRQ5mHQ19CoDqRcFITyzJxuEq001c001",
		"CAQACAQOQx8PN3s2y3pFBSh3VYqjlRE6dUwYB-1bxszjDCli5bYDeswHNU_sVfjMq001c001",
		"CAQACAQOYsGh3nyeVZkDlMMKKPZ0Nh94LFpW-GDoBBk55IymeNoBmBr53Fm2_FEQq001c001",
		"CAQACAQMqMpkk2tPujh4s2Xht0SZEJY3gZ_A71IPAQ5LM1TTdwYCOQzskZ373u1cq001c001",
		"CAQACAQPFZqeBxOg66nEyqj9SD4VdS0YS-CgKkQw_y1KJewmPWYDqywqB-K6Ni-kq001c001",
		"CAQACAQOlslIaEtd42OHK00vnqh5M80ULeK2mcLX1qp20y2HoC4DYqqYaeDTHdLEq001c001",
		"CAQACAQOu76VP-XAbAGUxMSB326UsjZTKXHLTkV3b5iAotzjlZoAA5tNPXBtlzxIq001c001",
		"CAQACAQN5p-YJhxLQjtKXPjNZ_RGFoTroYwsj90aoa5FQSmA7koCOayMJYwjlEyUq001c001",
		"CAQACAQNn5HggDkFNOJfZ19bI74Caox9kZhvRsZxxqca7ojfyNYA4qdEgZjGu9XIq001c001",
		"CAQACAQP6--IJnd4yMpNPKqN-mJ8isPvJ9al5CqvTcJ1AutMHe4AycHkJ9WJoE9cq001c001",
		"CAQACAQO5K0-w2lfXFkEzVebrzAPBQ8gjg0uNzsfUqGjImnd8b4AWqI2wg2CYi7sq001c001",
		"CAQACAQNfODcM1HIPfsRFZvm7OgNtLVhLKzGBs8s16NTE6gLI-4B-6IEMK5134AMq001c001",
		"CAQACAQOSsLnpk_C6yg5upA2VZGaCZtlFgeYOziKy3bdayhpfnIDK3Q7pgaTt1nEq001c001",
		"CAQACAQNCEyK80xgCTBmb0N23hR3V38zknp3LUsnJkLvRme5rHoBMkMu8nsS13qIq001c001",
		"CAQACAQOV8jrJnZskSGke0rnaYScs7t4Ww_gyvj_W-xlq1y_xboBI-zLJw3AC1Akq001c001",
		"CAQACAQPoIqx8QdWlRpLqoz2Vh8-T2A0C1aEnXZ2A3i1DLGatVIBG3id81eUW9Xgq001c001",
		"CAQACAQNVqz7wzgZzbqO34lL1Aj-_KrjUD-kRGZ8_ZKZj-CB2nYBuZBHwDyNXo1Eq001c001",
	}
)

const (
	// baseURL = "http://production-api.acq-visa-clearing.melifrontends.com"
	baseURL    = "https://internal-api.mercadopago.com/acq/staging/visa-clearing"
	xAuthToken = "4732e48314010f3fd22430109f48c08dfc8ea514407787e87994c09f9e75010d"
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

func main() {
	// fmt.Printf("Reprocessing %d orders\n", len(orderIDList))

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
	fmt.Printf("%d order(s) reprocessed", count)
}

func executeReviewAndReprocess(wg *sync.WaitGroup, value string) {
	defer wg.Done()

	order := getByOrderID(value)
	fmt.Printf(
		"ID: %s, St: %s, ARN: %s, Update: %s, FileID: %s, BatchNum: %d\n",
		value,
		order.State,
		order.ARN,
		order.Revision.UpdatedAt,
		order.BatchKey.UniqueFileID,
		order.BatchKey.BatchNumber,
	)

	//order := getByARN(value)
	//fmt.Printf("ARN: %s, St: %s, ID: %s\n", value, order.State, order.ID)
	/*
		if order.State == "presenting" || order.State == "presentation_in_review" || order.State == "sent_exception" {
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

/*
func main() {
   timeNow := time.Now().Add(time.Hour)
   timeLimit := time.Now().Add(365 * 24 * time.Hour)

   fmt.Printf("time now: %s, time limit: %s\n", timeNow, timeLimit)
   fmt.Println(!timeNow.Before(timeLimit))
}
*/
