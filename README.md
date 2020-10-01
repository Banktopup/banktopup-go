# banktopup-go

## Usage

```go
import "github.com/banktopup/banktopup-go"
```

```go
client := banktopup.NewClient("deviceID", "accountNumber", "pin", "license")

client.Register(banktopup.RegisterParam{})
client.RegisterOTP(banktopup.RegisterOTPParam{})

client.GetTransactions(banktopup.GetTransactionsParam{})
client.Transfer(banktopup.TransferParam{})
client.Summary(banktopup.SummaryParam{})
```
