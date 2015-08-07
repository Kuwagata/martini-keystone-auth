# martini-keystone-auth

#### About:
Validates keystone auth tokens and injects them into the martini context.

#### Installation:

`go get github.com/Kuwagata/martini-keystone-auth`

#### Example usage:
```go
package main

import (
	"github.com/go-martini/martini"
	"github.com/Kuwagata/martini-keystone-auth"
)

func main() {
	identity_tokens_endpoint := "https://identity.api.rackspacecloud.com/v2.0/tokens"

	redis_hostname := "localhost"
	redis_port := "6379"
	redis_pw := ""
	redis_db := int64(0)

	m := martini.Classic()
	m.Use(auth.Keystone(
		auth.GetIdentityValidator(identity_tokens_endpoint),
		auth.GetRedisCache(redis_hostname, redis_port, redis_pw, redis_db)))
	m.Get("/", func(token auth.Token) string {
		return string(token)
	})
	m.Run()
}
```
