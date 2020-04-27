package cep

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HandlerFunc func(ctx context.Context, cep string) CEP

func searchCorreios(ctx context.Context, cep string) CEP {
	payload := fmt.Sprintf(`<x:Envelope xmlns:x="http://schemas.xmlsoap.org/soap/envelope/" xmlns:cli="http://cliente.bean.master.sigep.bsb.correios.com.br/">
<x:Body>
<cli:consultaCEP>
<cep>%s</cep>
</cli:consultaCEP>
</x:Body>
</x:Envelope>`, cep)
	endpoint := "https://apps.correios.com.br/SigepMasterJPA/AtendeClienteService/AtendeCliente"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader([]byte(payload)))
	if err != nil {
		return CEP{}
	}

	req.Header.Set("Content-type", "text/xml; charset=utf-8")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	response, err := client.Do(req)
	if err != nil {
		return CEP{}
	}
	defer response.Body.Close()

	correio := Correio{}
	err = xml.NewDecoder(response.Body).Decode(&correio)
	if err != nil {
		return CEP{}
	}
	data := correio.Body.ConsultaCEPResponse.Return
	if data.Uf == "" {
		return CEP{}
	}
	return CEP{
		Logradouro: data.End,
		Bairro:     data.Bairro,
		Cidade:     data.Cidade,
		UF:         data.Uf,
		CEP:        cep,
		Base:       "correios",
	}
}

func searchPostmon(ctx context.Context, cep string) CEP {
	var data postmon
	get("https://api.postmon.com.br/v1/cep/%s", cep, &data)
	return CEP{
		Logradouro: data.Logradouro,
		Bairro:     data.Bairro,
		Cidade:     data.Cidade,
		UF:         data.UF,
		CEP:        cep,
		Base:       "postmon",
	}
}

func searchRepublicaVirtual(ctx context.Context, cep string) CEP {
	var data republicaVirtual
	get("https://republicavirtual.com.br/web_cep.php?cep=%s&formato=json", cep, &data)
	return CEP{
		Logradouro: fmt.Sprintf("%s %s", data.TipoLogradouro, data.Logradouro),
		Bairro:     data.Bairro,
		Cidade:     data.Cidade,
		UF:         data.UF,
		CEP:        cep,
		Base:       "republicavirtual",
	}
}

func searchViaCEP(ctx context.Context, cep string) CEP {
	var data viaCEP
	get("https://viacep.com.br/ws/%s/json/", cep, &data)
	return CEP{
		Logradouro: data.Logradouro,
		Bairro:     data.Bairro,
		Cidade:     data.Localidade,
		UF:         data.UF,
		CEP:        cep,
		Base:       "viacep",
	}
}

func get(url, cep string, data interface{}) {
	url = fmt.Sprintf(url, cep)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if body == nil || err != nil {
		return
	}
	json.Unmarshal(body, &data)
}
