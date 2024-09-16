package main

import (
	"fmt"

	"github.com/xuri/xgen"
)

func main() {
	// Inicializar as opções com os mapas necessários
	options := &xgen.Options{
		FilePath:            "schemas/ACCCTIPOS.xsd", // Arquivo XSD principal
		OutputDir:           "./output",              // Diretório de saída
		InputDir:            "./schemas",             // Diretório contendo todos os XSDs
		Lang:                "go",                    // Linguagem Go para a geração
		Package:             "schemas",               // Nome do pacote
		Extract:             false,
		IncludeMap:          make(map[string]bool),
		LocalNameNSMap:      make(map[string]string),
		NSSchemaLocationMap: make(map[string]string),
		ParseFileList:       make(map[string]bool),
		ParseFileMap:        make(map[string][]interface{}),
		RemoteSchema:        make(map[string][]byte),
	}

	// Fazer o parsing e gerar as structs
	err := options.Parse()
	if err != nil {
		fmt.Println("Erro ao fazer parsing:", err)
		return
	}

	fmt.Println("Parsing e geração de structs concluídos!")
}
