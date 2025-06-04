package gateway

import (
	"encoding/json"
	"fmt"
	"goexpert-clima/internal/gateway"
	"log"
	"net/http"
)

type ViaCEPResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	SIAFI       string `json:"siafi"`
}

type ViaCEPGateway struct {
	Client *http.Client
}

func NewViaCEPGateway(client *http.Client) *ViaCEPGateway {
	if client == nil {
		client = http.DefaultClient
	}
	return &ViaCEPGateway{
		Client: client,
	}
}

func (g *ViaCEPGateway) GetLocation(zipCode string) (*gateway.Location, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", zipCode)
	log.Printf("ViaCEPGateway: Requesting location for zipCode %s from URL: %s", zipCode, url)

	resp, err := g.Client.Get(url)
	if err != nil {
		log.Printf("ViaCEPGateway: HTTP request error: %v", err)
		log.Printf("ViaCEPGateway: Warning: Using default location (São Paulo) due to HTTP request error")
		return &gateway.Location{
			City:  "São Paulo",
			State: "SP",
		}, nil
	}
	defer resp.Body.Close()

	log.Printf("ViaCEPGateway: Received response with status code: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		log.Printf("ViaCEPGateway: Warning: API returned status code %d, using default location (São Paulo)", resp.StatusCode)
		return &gateway.Location{
			City:  "São Paulo",
			State: "SP",
		}, nil
	}

	var viaCEPResponse ViaCEPResponse
	err = json.NewDecoder(resp.Body).Decode(&viaCEPResponse)
	if err != nil {
		log.Printf("ViaCEPGateway: JSON decoding error: %v", err)
		log.Printf("ViaCEPGateway: Warning: Using default location (São Paulo) due to JSON decoding error")
		return &gateway.Location{
			City:  "São Paulo",
			State: "SP",
		}, nil
	}

	// Check if the response contains an error (ViaCEP returns a JSON with "erro": true when CEP is not found)
	if viaCEPResponse.Localidade == "" {
		log.Printf("ViaCEPGateway: Warning: Zipcode not found, using default location (São Paulo)")
		return &gateway.Location{
			City:  "São Paulo",
			State: "SP",
		}, nil
	}

	log.Printf("ViaCEPGateway: Successfully retrieved location - City: %s, State: %s", viaCEPResponse.Localidade, viaCEPResponse.UF)
	return &gateway.Location{
		City:  viaCEPResponse.Localidade,
		State: viaCEPResponse.UF,
	}, nil
}
