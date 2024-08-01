# Gotil
Golang Utility Plugins

## Installation
```bash
go get -u github.com/nothing2512/gotil
```

## Features
- Elastic Search
- Encryption
- File Uploader
- HTTP Fetch
- JSON Parse
- JSON Stringify
- Mailer
- Object To Struct
- Parse Csv File to Struct
- Parse Excel File to Struct
- Parse JSON File to Struct
- PIN Generator
- RabbitMQ
- UUID Generator
- WebSocket

## Example Object to Struct
```go

package main

import (
    "github.com/nothing2512/gotil"
)

type Data struct {
    Name string `customtag:"name"`
}

func main() {
    var data Data
    originalData := map[string]any{
        "name": "Fulan",
    }
    gotil.ParseStruct(&data, originalData, "customtag")
}
```

## Example Send Mail
```go

package main

import (
    "github.com/nothing2512/gotil"
)

type Data struct {
    Name string `json:"name"`
}

func main() {
    m, err := gotil.NewMailer("email", "pass", "host", "port")
	if err != nil {
		panic(err)
	}
	data := Data{"Fulan"}

	m.Subject("subject")

	m.Recipients("mail@gmail.com")
	m.Cc("mail@gmail.com", "mail@gmail.com")
	m.Bcc("mail@gmail.com", "mail@gmail.com")

	m.SetHTMLFile("template.html", data)
	m.AttachFile("certificate.pdf", []byte{})

	err = m.Send()
	if err != nil {
		panic(err)
	}

	m.Close()
}
```

## Example Elastic Search
```go
package main

import (
    "fmt"
    "github.com/nothing2512/gotil"
)

type Test struct {
	ID   int    `json:"id" es:"id"`
	Name string `json:"name" es:"name"`
}

func (*Test) TableName() string {
	return "tests"
}

func main() {
	es, err := gotil.NewElasticSearch("http://0.0.0.0:9200")
    if err != nil {
        panic(err)
    }

	es.Save(&Test{1, "Fulan"})
	es.Save(&Test{2, "Fulan"})
	es.Save(&Test{3, "Fulan"})

	es.Update(&Test{1, "Fulanah"})

	es.Delete(&Test{3, "Fulan"})
	es.DeleteById("tests", 2)

	data := []Test{}
	es.Search(&data, "tests", "fulan", "name", "message")

	for _, x := range data {
		fmt.Println(x.ID, x.Name)
	}
}
```

## Example Rabbit MQ
```go
package main

import (
    "fmt"
    "github.com/nothing2512/gotil"
)

func main() {
	rabbit, err := gotil.NewRabbitMQ("rabbitmq_user", "rabbitmq_password", "0.0.0.0", "5672")
	if err != nil {
		panic(err)
	}
	rabbit.Publish("channel1", "Hello, RabbitMQ World!")
	rabbit.Consume("channel1", func(data string) {
		fmt.Println(data)
	})
}
```

## Example HTTP Fetch
```go

package main

import (
    "github.com/nothing2512/gotil"
)

type Response struct {
	Status	bool `json:"status"`
	Message	string `json:"message"`
	Data 	gotil.JSON `json:"data"`
}

func main() {
	var data Data
	f := gotil.HTTPFetcher{
		Method: "POST",
		Url: "",
		Headers: gotil.JSON{},
		Body: gotil.JSON
	}
	err := f.fetch(&data)
}
```

## Example Web Socket

- server.go
```go
package main

import (
	"fmt"

	"github.com/nothing2512/gotil"
)

func main() {
	ws := gotil.NewWebSocket("0.0.0.0:8080")

	// Server Handle Incoming Command
	ws.OnCommand(func(m gotil.WebSocketMessage) {
		fmt.Println(m.Command, m.Message)

		// Server Send Reply Message
		ws.Reply(m, "Reply Message")

		// Server blast to all connection
		ws.Blast("Blast Message")
	})

	// Start Server
	err := ws.Server("00000000000000000000000000000000", "1111111111111111")
	if err != nil {
		panic(err)
	}
}
```

- client.go
```go
package main

import (
	"fmt"

	"github.com/nothing2512/gotil"
)

func main() {
	ws := gotil.NewWebSocket("0.0.0.0:8080")

	// Start Client
	err := ws.Client()
	if err != nil {
		panic(err)
	}

	// Client Handle Incoming Message
	ws.OnMessage(func(m gotil.WebSocketMessage) {
		fmt.Println(m.Message)
	})

	// Client Send Message To Other Client
	ws.Send("target-uid", "Hello")

	// Client Send Command To server
	ws.Command(gotil.WebSocketMessage{
		Command: "cmd",
		Message: "msg",
	})
}
```

- Send trough HTTP
```http
POST /send HTTP/1.1
Host: 0.0.0.0:8080
Content-Type: application/json
Content-Length: 237

{
    "token": "32097440af2b367064e37c43f08821daddb6ece61de2f4a8bb5a205bb75f3a9fdc27ea70",
    "to": "5fe7d69e-9432-86a0-0585-fd7c11c39e71",
    "command": "command|send",
    "message": "{\"command\":\"cmd\",\"message\": \"message\"}"
}
```