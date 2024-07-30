package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/clbanning/mxj/v2"
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

	mv, err := mxj.NewMapXml(xmlData)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter XML para Map: %v", err)
	}

	jsonData, err := json.MarshalIndent(mv, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("erro ao converter Map para JSON: %v", err)
	}

	return jsonData, nil
}


func JSONToXML(jsonData []byte) ([]byte, error) {
	var mv map[string]interface{}
	if err := json.Unmarshal(jsonData, &mv); err != nil {
		return nil, fmt.Errorf("erro ao converter JSON para Map: %v", err)
	}

	xmlData, err := mxj.Map(mv).XmlIndent("", "  ")
	if err != nil {
		return nil, fmt.Errorf("erro ao converter Map para XML: %v", err)
	}

	xmlData = append([]byte(xml.Header), xmlData...)

	return xmlData, nil
}

func handleValidateXML(w http.ResponseWriter, r *http.Request) {
	xmlData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendError(w, http.StatusBadRequest, 1001, "Falha ao ler o XML do corpo da requisição.")
		return
	}

	if err := validateXMLWithXSD(xmlData, "./ACCC471.xsd"); err != nil {
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

	if err := validateXMLWithXSD(xmlData, "./ACCC471.xsd"); err != nil {
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

	if err := validateXMLWithXSD(xmlData, "./ACCC471.xsd"); err != nil {
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
