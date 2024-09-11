package app

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"go-xsd/internal/utils"
	"go-xsd/pkg/xmlvalidator"
)

// Estrutura completa para ACCCDOC e seus subelementos, incluindo atributos
type ACCCDOC struct {
	XMLName        xml.Name `xml:"ACCCDOC"`
	Xmlns          string   `xml:"xmlns,attr"`
	XmlnsXs        string   `xml:"xmlns:xs,attr"`
	XmlnsCat       string   `xml:"xmlns:cat,attr"`
	XmlnsXsi       string   `xml:"xmlns:xsi,attr"`
	SchemaLocation string   `xml:"xsi:schemaLocation,attr"`
	BCARQ          BCARQ    `xml:"BCARQ"`
	SISARQ         SISARQ   `xml:"SISARQ"`
	ESTARQ         string   `xml:"ESTARQ"`
}

// Defina BCARQ e subestruturas
type BCARQ struct {
	NomArq           string   `xml:"NomArq"`
	NumCtrlEmis      string   `xml:"NumCtrlEmis"`
	ISPBEmissor      string   `xml:"ISPBEmissor"`
	ISPBDestinatario string   `xml:"ISPBDestinatario"`
	DtHrArq          string   `xml:"DtHrArq"`
	SitReq           string   `xml:"SitReq"`
	GrupoSeq         GrupoSeq `xml:"Grupo_Seq"`
	DtRef            string   `xml:"DtRef"`
}

// Defina GrupoSeq
type GrupoSeq struct {
	NumSeq   string `xml:"NumSeq"`
	IndrCont string `xml:"IndrCont"`
}

// Defina SISARQ e subestruturas
type SISARQ struct {
	ACCC471 ACCC471 `xml:"ACCC471"`
}

type ACCC471 struct {
	GrupoGarSCR GrupoGarSCR `xml:"Grupo_ACCC471_GarSCR"`
}

type GrupoGarSCR struct {
	IdentdPartAdmdo string          `xml:"IdentdPartAdmdo"`
	NumCtrlIFGar    string          `xml:"NumCtrlIFGar"`
	TpGar           string          `xml:"TpGar"`
	TpGarSCR        string          `xml:"TpGarSCR"`
	SubTpGarSCR     string          `xml:"SubTpGarSCR"`
	SitGar          string          `xml:"SitGar"`
	IndrBemFincd    string          `xml:"IndrBemFincd"`
	PercGar         string          `xml:"PercGar"`
	VlrOrGar        string          `xml:"VlrOrGar"`
	VlrGarDtReaval  string          `xml:"VlrGarDtReaval"`
	DtReaval        string          `xml:"DtReaval"`
	ClassRscCesta   string          `xml:"ClassRscCesta"`
	GrupoVeic       GrupoVeic       `xml:"Grupo_ACCC471_Veic"`
	GrupoImovel     GrupoImovel     `xml:"Grupo_ACCC471_Imovel"`
	GarFidejussoria GarFidejussoria `xml:"Grupo_ACCC471_GarFidejussoria"`
}

type GrupoVeic struct {
	// Veículo - defina seus atributos conforme necessário
	// Este é um exemplo com os campos do seu XML
	VlrEntdVeic  string `xml:"VlrEntdVeic"`
	IdentdChassi string `xml:"IdentdChassi"`
	TpVeic       string `xml:"TpVeic"`
	NumPlaca     string `xml:"NumPlaca"`
}

type GrupoImovel struct {
	// Imóvel - defina seus atributos conforme necessário
	NumMatriclImovl string         `xml:"NumMatriclImovl"`
	IdCartrio       string         `xml:"IdCartrio"`
	GrupoEndImovel  GrupoEndImovel `xml:"Grupo_ACCC471_EndImovel"`
}

type GrupoEndImovel struct {
	// Endereço do Imóvel
	EndImovl  string `xml:"EndImovl"`
	CEPImovl  string `xml:"CEPImovl"`
	CidImovl  string `xml:"CidImovl"`
	UFImovl   string `xml:"UFImovl"`
	PaisImovl string `xml:"PaisImovl"`
}

type GarFidejussoria struct {
	// Fidejussória - defina seus atributos conforme necessário
	SeqGarFidjssria            string `xml:"SeqGarFidjssria"`
	TpPessoaGarFidjssria       string `xml:"TpPessoaGarFidjssria"`
	CNPJ_CPFPessoaGarFidjssria string `xml:"CNPJ_CPFPessoaGarFidjssria"`
}

// Função para converter JSON para XML
func JSONToXML(jsonData []byte) ([]byte, error) {
	var mv ACCCDOC
	if err := json.Unmarshal(jsonData, &mv); err != nil {
		return nil, fmt.Errorf("erro ao converter JSON para estrutura: %v", err)
	}

	// Definir os namespaces e schemaLocation se estiverem faltando
	if mv.Xmlns == "" {
		mv.Xmlns = "http://www.cip-bancos.org.br/ARQ/ACCC471.xsd"
	}
	if mv.XmlnsXs == "" {
		mv.XmlnsXs = "http://www.w3.org/2001/XMLSchema"
	}
	if mv.XmlnsCat == "" {
		mv.XmlnsCat = "http://www.cip-bancos.org.br/catalogomsg"
	}
	if mv.XmlnsXsi == "" {
		mv.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	}
	if mv.SchemaLocation == "" {
		mv.SchemaLocation = "http://www.cip-bancos.org.br/ARQ/ACCC471.xsd ACCC471.xsd"
	}

	// Serializa a estrutura para XML
	xmlData, err := xml.MarshalIndent(mv, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar XML: %v", err)
	}

	// Adicionar o cabeçalho XML
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
		fmt.Println(err)
		utils.SendError(w, http.StatusBadRequest, 1008, fmt.Sprintf("Falha na validação do XML após a conversão: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(xmlData)
}
