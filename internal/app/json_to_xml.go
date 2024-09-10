package app

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"go-xsd/internal/utils"
	"go-xsd/pkg/xmlvalidator"
)

func jsonToXMLHelper(v interface{}) ([]byte, error) {
	switch vv := v.(type) {
	case map[string]interface{}:
		var buf bytes.Buffer
		buf.WriteString("<root>")
		for key, value := range vv {
			child, err := jsonToXMLHelper(value)
			if err != nil {
				return nil, err
			}
			buf.WriteString(fmt.Sprintf("<%s>%s</%s>", key, child, key))
		}
		buf.WriteString("</root>")
		return buf.Bytes(), nil
	case []interface{}:
		var buf bytes.Buffer
		for _, value := range vv {
			child, err := jsonToXMLHelper(value)
			if err != nil {
				return nil, err
			}
			buf.Write(child)
		}
		return buf.Bytes(), nil
	default:
		return []byte(fmt.Sprintf("%v", vv)), nil
	}
}

func JSONToXML(jsonData []byte) ([]byte, error) {
	var mv map[string]interface{}
	if err := json.Unmarshal(jsonData, &mv); err != nil {
		return nil, fmt.Errorf("erro ao converter JSON para Map: %v", err)
	}

	if _, ok := mv["ACCCDOC"]; !ok {
		mv["ACCCDOC"] = map[string]interface{}{}
	}

	acccdoc := mv["ACCCDOC"].(map[string]interface{})
	acccdoc["-xmlns"] = "http://www.cip-bancos.org.br/ARQ/ACCC471.xsd"
	acccdoc["-xmlns:xs"] = "http://www.w3.org/2001/XMLSchema"
	acccdoc["-xmlns:cat"] = "http://www.cip-bancos.org.br/catalogomsg"
	acccdoc["-xmlns:xsi"] = "http://www.w3.org/2001/XMLSchema-instance"
	acccdoc["-xsi:schemaLocation"] = "http://www.cip-bancos.org.br/ARQ/ACCC471.xsd ACCC471.xsd"

	xmlData, err := jsonToXMLHelper(mv)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter Map para XML: %v", err)
	}

	xmlData = append([]byte(xml.Header), xmlData...)

	return xmlData, nil
}

func HandleJSONToXML(w http.ResponseWriter, r *http.Request) {
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, 1006, "Falha ao ler o JSON do corpo da requisição.")
		return
	}

	xmlData, err := JSONToXML(jsonData)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, 1007, fmt.Sprintf("Erro ao converter JSON para XML: %v", err))
		return
	}

	if err := xmlvalidator.ValidateXML(xmlData, "./xsds/ACCC471.xsd"); err != nil {
		utils.SendError(w, http.StatusBadRequest, 1008, fmt.Sprintf("Falha na validação do XML após a conversão: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(xmlData)
}
