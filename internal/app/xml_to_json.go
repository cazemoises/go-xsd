package app

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"go-xsd/internal/utils"
	validator "go-xsd/pkg/xmlvalidator"

	xml2json "github.com/basgys/goxml2json"
)

func XMLToJSON(xmlData []byte) ([]byte, error) {
	jsonData, err := xml2json.Convert(bytes.NewReader(xmlData))
	if err != nil {
		return nil, fmt.Errorf("erro ao converter XML para JSON: %v", err)
	}

	return jsonData.Bytes(), nil
}

func HandleXMLToJSON(w http.ResponseWriter, r *http.Request) {
	xmlData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, 1003, "Falha ao ler o XML do corpo da requisição.")
		return
	}

	if len(xmlData) == 0 {
		utils.SendError(w, http.StatusBadRequest, 1004, "O corpo da requisição XML está vazio.")
		return
	}

	xsdPath := "./xsds/ACCC471.xsd"

	if err := validator.ValidateXML(xmlData, xsdPath); err != nil {
		utils.SendError(w, http.StatusBadRequest, 1004, fmt.Sprintf("Falha na validação do XML antes da conversão: %v", err))
		return
	}

	jsonData, err := XMLToJSON(xmlData)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, 1005, fmt.Sprintf("Erro ao converter XML para JSON: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
