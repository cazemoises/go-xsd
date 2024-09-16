package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	schemas "go-xsd/output/schemas"
	"io/ioutil"
	"os"

	"github.com/xuri/xgen"
)

func main() {
	// Flag para definir o que fazer: gerar structs, converter XML para JSON ou JSON para XML
	mode := flag.String("mode", "generate", "Modo de execução: 'generate' para gerar structs, 'convert' para XML -> JSON ou 'convert-to-xml' para JSON -> XML.")
	flag.Parse()

	switch *mode {
	case "generate":
		// Chama a função para gerar as structs a partir do XSD
		generateStructs()
	case "convert":
		// Chama a função para converter XML para JSON
		convertXMLToJSON()
	case "convert-to-xml":
		// Chama a função para converter JSON para XML
		convertJSONToXML()
	default:
		fmt.Println("Modo inválido. Use 'generate', 'convert', ou 'convert-to-xml'.")
	}
}

func generateStructs() {
	// Inicializar as opções com os mapas necessários
	options := &xgen.Options{
		FilePath:            "schemas/ACCC471.xsd",
		OutputDir:           "./output",
		InputDir:            "./schemas",
		Lang:                "go",
		Package:             "main",
		Extract:             false,
		IncludeMap:          make(map[string]bool),          // Inicializar o mapa
		LocalNameNSMap:      make(map[string]string),        // Inicializar o mapa
		NSSchemaLocationMap: make(map[string]string),        // Inicializar o mapa
		ParseFileList:       make(map[string]bool),          // Inicializar o mapa
		ParseFileMap:        make(map[string][]interface{}), // Inicializar o mapa
		RemoteSchema:        make(map[string][]byte),        // Inicializar o mapa
	}

	// Fazer o parsing e gerar as structs
	err := options.Parse()
	if err != nil {
		fmt.Println("Erro ao fazer parsing:", err)
		return
	}

	fmt.Println("Parsing e geração de structs concluídos!")
}

func convertXMLToJSON() {
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

	// Salvar o JSON em um arquivo
	err = ioutil.WriteFile("output.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Erro ao salvar o JSON:", err)
		return
	}

	fmt.Println("Arquivo JSON salvo com sucesso!")
}

func convertJSONToXML() {
	// Abrir o arquivo JSON de entrada
	jsonFile, err := os.Open("output.json")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo JSON:", err)
		return
	}
	defer jsonFile.Close()

	// Ler o conteúdo do arquivo JSON
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo JSON:", err)
		return
	}

	// Fazer o Unmarshal do JSON para as structs Go geradas
	var doc schemas.ACCCDOC // Substitua ACCCDOC pela struct raiz gerada
	err = json.Unmarshal(jsonData, &doc)
	if err != nil {
		fmt.Println("Erro ao fazer unmarshal do JSON:", err)
		return
	}

	// Converter a struct Go para XML
	xmlData, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		fmt.Println("Erro ao converter para XML:", err)
		return
	}

	// Exibir o XML gerado
	fmt.Println("XML gerado:")
	fmt.Println(string(xmlData))

	// Salvar o XML em um arquivo
	err = ioutil.WriteFile("output.xml", xmlData, 0644)
	if err != nil {
		fmt.Println("Erro ao salvar o XML:", err)
		return
	}

	fmt.Println("Arquivo XML salvo com sucesso!")
}
