# powerschool-api

Go package for PowerSchool API access.

---

## Install

* Requires Go 1.18 or above.
* PowerSchool [OAuth.](https://support.powerschool.com/developer/#/page/oauth)
* PowerSchool [Access Request.](https://support.powerschool.com/developer/#/page/access-request---field-access)

Install with `go get github.com/ledfree/powerschool-api`

---

## Index

[type ApiConfig](#type-ApiConfig)

### type ApiConfig

```go
type ApiConfig struct {
  Access_Token string `json:"access_token"`
  Token_Type   string `json:"token_type"`
  Expires_In   string `json:"expires_in"`
}
```