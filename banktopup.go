package banktopup

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	EndPoint = "https://api-v1.banktopup.com"
)

type (
	Client struct {
		client *http.Client

		deviceID      string
		accountNumber string
		pin           string
		license       string
	}
)

func NewClient(deviceID, accountNumber, pin, license string) *Client {
	return &Client{
		client: http.DefaultClient,

		deviceID:      deviceID,
		accountNumber: accountNumber,
		pin:           pin,
		license:       license,
	}
}

type (
	RegisterParam struct {
		Identification string `json:"identification"`
		AccountNo      string `json:"account_no"`
		PIN            string `json:"pin"`
		Phone          string `json:"mobile_phone_no"`
		DeviceBrand    string `json:"device_brand"`
		DeviceCode     string `json:"device_code"`

		Year  string `json:"year"`
		Month string `json:"month"`
		Day   string `json:"day"`
	}

	RegisterResponse struct {
		Error struct {
			Code  int    `json:"code"`
			MsgTH string `json:"msg_th"`
		} `json:"error"`
		Result struct {
			Message  string `json:"msg"`
			DeviceID string `json:"deviceid"`
		} `json:"result"`
	}
)

func (c *Client) Register(param RegisterParam) (*RegisterResponse, error) {
	req, _ := http.NewRequest("POST", EndPoint+"/api/v1/scb/register", marshalJSON(param))
	req.Header.Add("x-auth-license", c.license)
	req.Header.Add("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	var response RegisterResponse
	if err := parseResponse(res, &response); err != nil {
		return nil, err
	}
	if response.Error.MsgTH != "สำเร็จ" {
		return nil, errors.New(response.Error.MsgTH)
	}
	return &response, nil
}

type (
	RegisterOTPParam struct {
		OTP string `json:"otp"`
	}
	RegisterOTPResponse struct {
		Error struct {
			Code  int    `json:"code"`
			MsgTH string `json:"msg_th"`
		} `json:"error"`
		Result struct {
			Message  string `json:"msg"`
			DeviceID string `json:"deviceid"`
		} `json:"result"`
	}
)

func (c *Client) RegisterOTP(param RegisterOTPParam) (*RegisterOTPResponse, error) {
	req, _ := http.NewRequest("POST", EndPoint+"/api/v1/scb/register/"+c.accountNumber, marshalJSON(param))
	req.Header.Add("x-auth-license", c.license)
	req.Header.Add("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	var response RegisterOTPResponse
	if err := parseResponse(res, &response); err != nil {
		return nil, err
	}
	if response.Error.MsgTH != "สำเร็จ" {
		return nil, errors.New(response.Error.MsgTH)
	}
	return &response, nil
}

type (
	GetTransactionsParam struct {
		PreviousDay int `json:"previous_day"`
		PageNumber  int `json:"page_number"`
		PageSize    int `json:"page_size"`

		// optional
		DeviceID      string `json:"deviceid,omitempty"`
		AccountNumber string `json:"account_no,omitempty"`
		PIN           string `json:"pin,omitempty"`
	}
	GetTransactionsResponse struct {
		Error struct {
			Code  int    `json:"code"`
			MsgTH string `json:"msg_th"`
		} `json:"error"`
		Result struct {
			AccountNo      string `json:"account_no"`
			EndOfListFlag  string `json:"endOfListFlag"`
			NextPageNumber string `json:"nextPageNumber"`
			PageSize       int    `json:"pageSize"`
			TxnList        []struct {
				Annotation       string  `json:"annotation"`
				SortSequence     int     `json:"sortSequence"`
				TxnAmount        float64 `json:"txnAmount"`
				TxnDateTime      string  `json:"txnDateTime"`
				TxnCurrency      string  `json:"txnCurrency"`
				TxnRemark        string  `json:"txnRemark"`
				TxnDebitCardFlag string  `json:"txnDebitCreditFlag"`
				TxnSequence      int     `json:"txnSequence"`
				TxnChannel       struct {
					Code        string `json:"code"`
					Description string `json:"description"`
				} `json:"txnChannel"`
				TxnCode struct {
					Code        string `json:"code"`
					Description string `json:"description"`
				} `json:"txnCode"`
			} `json:"txnList"`
		} `json:"result"`
	}
)

func (c *Client) GetTransactions(param GetTransactionsParam) (*GetTransactionsResponse, error) {
	param.DeviceID = c.deviceID
	param.AccountNumber = c.accountNumber
	param.PIN = c.pin

	req, _ := http.NewRequest("POST", EndPoint+"/api/v1/scb/transactions", marshalJSON(param))
	req.Header.Add("x-auth-license", c.license)
	req.Header.Add("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	var response GetTransactionsResponse
	if err := parseResponse(res, &response); err != nil {
		return nil, err
	}
	if response.Error.MsgTH != "สำเร็จ" {
		return nil, errors.New(response.Error.MsgTH)
	}
	return &response, nil
}

type (
	TransferParam struct {
		AccountTo string  `json:"account_to"`
		BankCode  string  `json:"bank_code"`
		Amount    float64 `json:"amount"`

		// optional
		DeviceID      string `json:"deviceid,omitempty"`
		AccountNumber string `json:"account_no,omitempty"`
		PIN           string `json:"pin,omitempty"`
	}
	TransferResponse struct {
		Error struct {
			Code  int    `json:"code"`
			MsgTH string `json:"msg_th"`
		} `json:"error"`
		Result struct {
			TransactionID       string  `json:"transactionId"`
			TransactionDateTime string  `json:"transactionDateTime"`
			RemainingBalance    float64 `json:"remainingBalance"`
			AdditionalMetaData  struct {
				PaymentInfo []struct {
					QRString string `json:"QRstring"`
				} `json:"paymentInfo"`
			} `json:"additionalMetaData"`
		} `json:"result"`
	}
)

func (c *Client) Transfer(param TransferParam) (*TransferResponse, error) {
	param.DeviceID = c.deviceID
	param.AccountNumber = c.accountNumber
	param.PIN = c.pin

	req, _ := http.NewRequest("POST", EndPoint+"/api/v1/scb/transfer", marshalJSON(param))
	req.Header.Add("x-auth-license", c.license)
	req.Header.Add("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	var response TransferResponse
	if err := parseResponse(res, &response); err != nil {
		return nil, err
	}
	if response.Error.MsgTH != "สำเร็จ" {
		return nil, errors.New(response.Error.MsgTH)
	}
	return &response, nil
}

type (
	SummaryParam struct {
		// optional
		DeviceID      string `json:"deviceid,omitempty"`
		AccountNumber string `json:"account_no,omitempty"`
		PIN           string `json:"pin,omitempty"`
	}
	SummaryResponse struct {
		Error struct {
			Code  int    `json:"code"`
			MsgTH string `json:"msg_th"`
		} `json:"error"`
		Result struct {
			TotalAvailableBalance float64 `json:"totalAvailableBalance"`
		} `json:"result"`
	}
)

func (c *Client) Summary(param SummaryParam) (*SummaryResponse, error) {
	param.DeviceID = c.deviceID
	param.AccountNumber = c.accountNumber
	param.PIN = c.pin

	req, _ := http.NewRequest("POST", EndPoint+"/api/v1/scb/summary", marshalJSON(param))
	req.Header.Add("x-auth-license", c.license)
	req.Header.Add("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	var response SummaryResponse
	if err := parseResponse(res, &response); err != nil {
		return nil, err
	}
	if response.Error.MsgTH != "สำเร็จ" {
		return nil, errors.New(response.Error.MsgTH)
	}
	return &response, nil
}

func parseResponse(res *http.Response, ret interface{}) error {
	defer res.Body.Close()
	buff, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(buff, &ret)
}

func marshalJSON(data interface{}) io.Reader {
	buff, _ := json.Marshal(data)
	return bytes.NewReader(buff)
}
