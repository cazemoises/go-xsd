package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	schemas "go-xsd/output/schemas"
)

func main() {
	// Abrir o arquivo XML de entrada
	xmlFile, err := os.Open("schemas/test.xml")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo XML:", err)
		return
	}
	defer xmlFile.Close()

	// Ler o conteúdo do arquivo XML
	xmlData, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo XML:", err)
		return
	}

	// Fazer o Unmarshal do XML para as structs Go geradas
	// Lembre-se de ajustar a struct de acordo com o que foi gerado (substitua ACCCDOC pela raiz gerada)
	var doc schemas.ACCCDOC // Essa struct foi gerada pelo xgen no arquivo `output_structs.go`
	err = xml.Unmarshal(xmlData, &doc)
	if err != nil {
		fmt.Println("Erro ao fazer unmarshal do XML:", err)
		return
	}

	// Converter a struct Go para JSON
	jsonData, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return
	}

	// Exibir o JSON gerado
	fmt.Println("JSON gerado:")
	fmt.Println(string(jsonData))

	// Opcional: salvar o JSON em um arquivo
	err = ioutil.WriteFile("output.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Erro ao salvar o JSON:", err)
		return
	}

	fmt.Println("Arquivo JSON salvo com sucesso!")
}
