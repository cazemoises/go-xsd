package app

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"go-xsd/internal/utils"
	validator "go-xsd/pkg/xmlvalidator"
)

func ValidateXMLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendError(w, http.StatusMethodNotAllowed, 1001, "Invalid request method")
		return
	}

	xsdFile := r.URL.Query().Get("xsd")
	if xsdFile == "" {
		utils.SendError(w, http.StatusBadRequest, 1002, "Missing XSD file parameter")
		return
	}

	xsdPath := fmt.Sprintf("xsds/%s", xsdFile)
	if _, err := ioutil.ReadFile(xsdPath); err != nil {
		utils.SendError(w, http.StatusNotFound, 1003, "XSD file not found")
		return
	}

	xmlFile, _, err := r.FormFile("file")
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, 1004, "Error reading file")
		return
	}
	defer xmlFile.Close()

	xmlData, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, 1005, "Error reading file data")
		return
	}

	err = validator.ValidateXML(xmlData, xsdPath)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, 1006, "Invalid XML: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("XML is valid"))
}
