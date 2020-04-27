package cep

import "encoding/xml"

type CEP struct {
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Cidade     string `json:"cidade"`
	UF         string `json:"estado"`
	CEP        string `json:"cep"`
	Base       string `json:"base"`
}

func (c CEP) FromCorreios() bool {
	return c.Base == "correios"
}

func (c CEP) Valid() bool {
	return c.UF != ""
}

type Correio struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soap    string   `xml:"soap,attr"`
	Body    struct {
		Text                string `xml:",chardata"`
		ConsultaCEPResponse struct {
			Text   string `xml:",chardata"`
			Ns2    string `xml:"ns2,attr"`
			Return struct {
				Text         string `xml:",chardata"`
				Bairro       string `xml:"bairro"`
				Cep          string `xml:"cep"`
				Cidade       string `xml:"cidade"`
				Complemento2 string `xml:"complemento2"`
				End          string `xml:"end"`
				Uf           string `xml:"uf"`
			} `xml:"return"`
		} `xml:"consultaCEPResponse"`
	} `xml:"Body"`
}

type cepLa struct {
	Bairro     string `json:"bairro"`
	Cidade     string `json:"cidade"`
	Logradouro string `json:"logradouro"`
	UF         string `json:"estado"`
}

type postmon struct {
	Bairro     string `json:"bairro"`
	Cidade     string `json:"cidade"`
	Logradouro string `json:"logradouro"`
	UF         string `json:"estado"`
}

type republicaVirtual struct {
	Bairro         string `json:"bairro"`
	Cidade         string `json:"cidade"`
	Logradouro     string `json:"logradouro"`
	TipoLogradouro string `json:"tipo_logradouro"`
	UF             string `json:"uf"`
}

type viaCEP struct {
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Logradouro string `json:"logradouro"`
	UF         string `json:"uf"`
}
