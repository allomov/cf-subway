package broker

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
	"gopkg.in/yaml.v2"
)

// Broker is the core struct for the Broker webapp
type Broker struct {
	Catalog        []brokerapi.Service
	BackendBrokers []*BackendBroker

	Logger lager.Logger
}

// NewBroker is a constructor for a Broker webapp struct
func NewBroker() (subway *Broker) {
	subway = &Broker{}
	subway.Logger = lager.NewLogger("cf-subway")
	subway.Logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	subway.Logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))
	return subway
}

// LoadCatalog loads the homogenous catalog from one of the brokers
func (subway *Broker) LoadCatalog() error {
	if len(subway.BackendBrokers) == 0 {
		return errors.New("No backend broker available for plan")
	}
	backendBroker := subway.BackendBrokers[0]

	client := &http.Client{}
	url := fmt.Sprintf("%s/v2/catalog", backendBroker.URI)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		subway.Logger.Error("backend-catalog-req", err)
		return err
	}
	req.SetBasicAuth(backendBroker.Username, backendBroker.Password)

	resp, err := client.Do(req)
	if err != nil {
		subway.Logger.Error("backend-catalog-resp", err)
		return err
	}
	defer resp.Body.Close()

	jsonData, err := ioutil.ReadAll(resp.Body)

	catalogResponse := brokerapi.CatalogResponse{}
	err = yaml.Unmarshal(jsonData, &catalogResponse)
	if err != nil {
		return err
	}
	subway.Catalog = catalogResponse.Services
	return nil
}

func (subway *Broker) plans() []brokerapi.ServicePlan {
	if len(subway.Catalog) == 0 {
		subway.LoadCatalog()
	}
	plans := []brokerapi.ServicePlan{}
	for _, service := range subway.Catalog {
		for _, plan := range service.Plans {
			plans = append(plans, plan)
		}
	}
	return plans
}

// Run starts the Martini webapp handler
func (subway *Broker) Run() {
	username := os.Getenv("SUBWAY_USERNAME")
	if username == "" {
		username = "username"
	}

	password := os.Getenv("SUBWAY_PASSWORD")
	if password == "" {
		password = "password"
	}

	credentials := brokerapi.BrokerCredentials{
		Username: username,
		Password: password,
	}
	fmt.Println(credentials)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	brokerAPI := brokerapi.New(subway, subway.Logger, credentials)
	http.Handle("/", brokerAPI)
	subway.Logger.Fatal("http-listen", http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil))
}
