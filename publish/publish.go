package publish

import (
	"io/ioutil"
	"log"

	"github.com/nats-io/stan.go"
)

func main() {

	source_model := []string{"test_model3.json"}

	sc, err := stan.Connect("test-cluster", "publisher", stan.NatsURL(stan.DefaultNatsURL))

	if err != nil {
		log.Panicln("Error connecting to NATS as publisher", err)
		return
	}

	defer sc.Close()

	for _, value := range source_model {
		values, err := ioutil.ReadFile(value)
		if err != nil {
			log.Panicln("Error reading file", err)
			return
		}
		sc.Publish("orders_model", values)
	}

	//sc.Publish("foo", []byte("Hello World"))

}
