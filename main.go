package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	xml2json "github.com/basgys/goxml2json"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func sendError(w http.ResponseWriter, code int, errorCode int, message string) {
	w.WriteHeader(code)
	errorResponse := ErrorResponse{
		Code:    errorCode,
		Message: message,
	}
	jsonResponse, _ := json.Marshal(errorResponse)
	w.Write(jsonResponse)
}

func validateXMLWithXSD(xmlData []byte, xsdPath string) error {
	tmpFile, err := ioutil.TempFile("", "*.xml")
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo temporário: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(xmlData); err != nil {
		return fmt.Errorf("erro ao escrever no arquivo temporário: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("erro ao fechar arquivo temporário: %v", err)
	}

	cmd := exec.Command("xmllint", "--noout", "--schema", xsdPath, tmpFile.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("falha na validação: %v, output: %s", err, string(output))
	}

	return nil
}

func XMLToJSON(xmlData []byte) ([]byte, error) {
	jsonData, err := xml2json.Convert(bytes.NewReader(xmlData))
	if err != nil {
		return nil, fmt.Errorf("erro ao converter XML para JSON: %v", err)
	}

	return jsonData.Bytes(), nil
}

func JSONToXML(jsonData []byte) ([]byte, error) {
	var mv map[string]interface{}
	if err := json.Unmarshal(jsonData, &mv); err != nil {
		return nil, fmt.Errorf("erro ao converter JSON para Map: %v", err)
	}

	// Adiciona manualmente os namespaces e a raiz do XML se não estiverem presentes
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

func handleValidateXML(w http.ResponseWriter, r *http.Request) {
	xmlData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendError(w, http.StatusBadRequest, 1001, "Falha ao ler o XML do corpo da requisição.")
		return
	}

	if err := validateXMLWithXSD(xmlData, "./xsds/ACCC471.xsd"); err != nil {
		sendError(w, http.StatusBadRequest, 1002, fmt.Sprintf("Falha na validação do XML: %v", err))
		return
	}

	fmt.Fprintln(w, "XML is valid")
}

func handleXMLToJSON(w http.ResponseWriter, r *http.Request) {
	xmlData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendError(w, http.StatusBadRequest, 1003, "Falha ao ler o XML do corpo da requisição.")
		return
	}

	if err := validateXMLWithXSD(xmlData, "./xsds/ACCC471.xsd"); err != nil {
		sendError(w, http.StatusBadRequest, 1004, fmt.Sprintf("Falha na validação do XML antes da conversão: %v", err))
		return
	}

	jsonData, err := XMLToJSON(xmlData)
	if err != nil {
		sendError(w, http.StatusInternalServerError, 1005, fmt.Sprintf("Erro ao converter XML para JSON: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func handleJSONToXML(w http.ResponseWriter, r *http.Request) {
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendError(w, http.StatusBadRequest, 1006, "Falha ao ler o JSON do corpo da requisição.")
		return
	}

	xmlData, err := JSONToXML(jsonData)
	if err != nil {
		sendError(w, http.StatusInternalServerError, 1007, fmt.Sprintf("Erro ao converter JSON para XML: %v", err))
		return
	}

	fmt.Println(string(xmlData))

	if err := validateXMLWithXSD(xmlData, "./xsds/ACCC471.xsd"); err != nil {
		sendError(w, http.StatusBadRequest, 1008, fmt.Sprintf("Falha na validação do XML após a conversão: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(xmlData)
}

func main() {
	path, err := exec.LookPath("xmllint")
	if err != nil {
		fmt.Println("xmllint não encontrado no path")
		os.Exit(1)
	}
	fmt.Println("xmllint encontrado em", path)

	http.HandleFunc("/validate-xml", handleValidateXML)
	http.HandleFunc("/xml-to-json", handleXMLToJSON)
	http.HandleFunc("/json-to-xml", handleJSONToXML)
	fmt.Println("Server started on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed to start:", err)
	}
}
