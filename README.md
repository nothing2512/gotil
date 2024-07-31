# Gotil
Golang Utility Plugins

## Installation
```bash
go get -u github.com/nothing2512/gotil
```

## Features
- Elastic Search
- Encryption
- HTTP Fetch
- Mailer
- Object To Struct
- PIN Generator
- RabbitMQ
- File Uploader
- UUID Generator

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