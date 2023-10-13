package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	// Atribuindo a função BuscaCEP a uma URL
	http.HandleFunc("/", BuscaCepHandler)

	// Subindo um servidor HTTP na porta 8080
	http.ListenAndServe(":8080", nil)
}

// Recebendo o Request e Response
func BuscaCepHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// Se a rota for diferente, ele retorna Not Found
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Pegando qualquer parâmetro CEP da URL
	cepParam := r.URL.Query().Get("cep")
	if cepParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Chamando a função BuscaCEP
	cep, err := BuscaCEP(cepParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Setando um Header para o retorno ser no formato de JSON
	w.Header().Set("Content-Type", "application/json")
	//Retornando um Status de OK
	w.WriteHeader(http.StatusOK)

	// Mostrando os dados retornados pela API
	json.NewEncoder(w).Encode(cep)
}

func BuscaCEP(cep string) (*ViaCEP, error) {
	response, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, err
	}

	// Fechando a conexão com o servidor
	defer response.Body.Close()

	// Lendo o Body do Response
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Transformando o JSON em Struct
	var c ViaCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
