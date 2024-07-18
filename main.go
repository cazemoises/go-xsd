package main
 
import (
    "encoding/json"
    "encoding/xml"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "os/exec"
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
 
type ACCCDOC struct {
    XMLName     xml.Name `xml:"ACCCDOC" json:"-"`
    BCARQ       BCARQ    `xml:"BCARQ" json:"BCARQ"`
    SISARQ      SISARQ   `xml:"SISARQ" json:"SISARQ"`
    ESTARQ      string   `xml:"ESTARQ" json:"ESTARQ"`
}
 
type BCARQ struct {
    NomArq           string   `xml:"NomArq" json:"NomArq"`
    NumCtrlEmis      string   `xml:"NumCtrlEmis" json:"NumCtrlEmis"`
    ISPBEmissor      string   `xml:"ISPBEmissor" json:"ISPBEmissor"`
    ISPBDestinatario string   `xml:"ISPBDestinatario" json:"ISPBDestinatario"`
    DtHrArq          string   `xml:"DtHrArq" json:"DtHrArq"`
    SitReq           string   `xml:"SitReq" json:"SitReq"`
    GrupoSeq         GrupoSeq `xml:"Grupo_Seq" json:"Grupo_Seq"`
    DtRef            string   `xml:"DtRef" json:"DtRef"`
}
 
type GrupoSeq struct {
    NumSeq   string `xml:"NumSeq" json:"NumSeq"`
    IndrCont string `xml:"IndrCont" json:"IndrCont"`
}
 
type SISARQ struct {
    ACCC471 ACCC471 `xml:"ACCC471" json:"ACCC471"`
}
 
type ACCC471 struct {
    GrupoGarSCR GrupoGarSCR `xml:"Grupo_ACCC471_GarSCR" json:"Grupo_ACCC471_GarSCR"`
}
 
type GrupoGarSCR struct {
    IdentdPartAdmdo   string `xml:"IdentdPartAdmdo" json:"IdentdPartAdmdo"`
    NumCtrlIFGar      string `xml:"NumCtrlIFGar" json:"NumCtrlIFGar"`
    TpGar             string `xml:"TpGar" json:"TpGar"`
    TpGarSCR          string `xml:"TpGarSCR" json:"TpGarSCR"`
    SubTpGarSCR       string `xml:"SubTpGarSCR" json:"SubTpGarSCR"`
    SitGar            string `xml:"SitGar" json:"SitGar"`
    IndrBemFincd      string `xml:"IndrBemFincd" json:"IndrBemFincd"`
    PercGar           string `xml:"PercGar" json:"PercGar"`
    VlrOrGar          string `xml:"VlrOrGar" json:"VlrOrGar"`
    VlrGarDtReaval    string `xml:"VlrGarDtReaval" json:"VlrGarDtReaval"`
    DtReaval          string `xml:"DtReaval" json:"DtReaval"`
    ClassRscCesta     string `xml:"ClassRscCesta" json:"ClassRscCesta"`
    GrupoVeic         GrupoVeic `xml:"Grupo_ACCC471_Veic" json:"Grupo_ACCC471_Veic"`
    GrupoImovel       GrupoImovel `xml:"Grupo_ACCC471_Imovel" json:"Grupo_ACCC471_Imovel"`
    GrupoFidejussoria GrupoFidejussoria `xml:"Grupo_ACCC471_GarFidejussoria" json:"Grupo_ACCC471_GarFidejussoria"`
}
 
type GrupoVeic struct {
    VlrEntdVeic  string `xml:"VlrEntdVeic" json:"VlrEntdVeic"`
    IdentdChassi string `xml:"IdentdChassi" json:"IdentdChassi"`
    TpVeic       string `xml:"TpVeic" json:"TpVeic"`
    TpTabVeicl   string `xml:"TpTabVeicl" json:"TpTabVeicl"`
    CodTabVeicl  string `xml:"CodTabVeicl" json:"CodTabVeicl"`
    UFCodTabVeicl string `xml:"UFCodTabVeicl" json:"UFCodTabVeicl"`
    IndrVeicUsado string `xml:"IndrVeicUsado" json:"IndrVeicUsado"`
    NumPlaca      string `xml:"NumPlaca" json:"NumPlaca"`
    UFNumPlaca    string `xml:"UFNumPlaca" json:"UFNumPlaca"`
    KM            string `xml:"KM" json:"KM"`
    RENAVAM       string `xml:"RENAVAM" json:"RENAVAM"`
    NumNota       string `xml:"NumNota" json:"NumNota"`
    NumSerNota    string `xml:"NumSerNota" json:"NumSerNota"`
    VlrNota       string `xml:"VlrNota" json:"VlrNota"`
    DtEmiss       string `xml:"DtEmiss" json:"DtEmiss"`
    GrupoIdentcVeic GrupoIdentcVeic `xml:"Grupo_ACCC471_IdentcVeic" json:"Grupo_ACCC471_IdentcVeic"`
}
 
type GrupoIdentcVeic struct {
    VlrMercVeic  string `xml:"VlrMercVeic" json:"VlrMercVeic"`
    CodCatg      string `xml:"CodCatg" json:"CodCatg"`
    CodMarca     string `xml:"CodMarca" json:"CodMarca"`
    CodModl      string `xml:"CodModl" json:"CodModl"`
    AnoModlVeicl string `xml:"AnoModlVeicl" json:"AnoModlVeicl"`
    AnoFabrccVeicl string `xml:"AnoFabrccVeicl" json:"AnoFabrccVeicl"`
}
 
type GrupoImovel struct {
    TpImovl             string `xml:"TpImovl" json:"TpImovl"`
    TpUsoImovl          string `xml:"TpUsoImovl" json:"TpUsoImovl"`
    NumInscrMuncplImovl string `xml:"NumInscrMuncplImovl" json:"NumInscrMuncplImovl"`
    NumMatriclImovl     string `xml:"NumMatriclImovl" json:"NumMatriclImovl"`
    IdCartrio           string `xml:"IdCartrio" json:"IdCartrio"`
    GrupoEndImovel      GrupoEndImovel `xml:"Grupo_ACCC471_EndImovel" json:"Grupo_ACCC471_EndImovel"`
}
 
type GrupoEndImovel struct {
    TpEndImovl string `xml:"TpEndImovl" json:"TpEndImovl"`
    EndImovl   string `xml:"EndImovl" json:"EndImovl"`
    CEPImovl   string `xml:"CEPImovl" json:"CEPImovl"`
    CidImovl   string `xml:"CidImovl" json:"CidImovl"`
    UFImovl    string `xml:"UFImovl" json:"UFImovl"`
    PaisImovl  string `xml:"PaisImovl" json:"PaisImovl"`
}
 
type GrupoFidejussoria struct {
    SeqGarFidjssria       string `xml:"SeqGarFidjssria" json:"SeqGarFidjssria"`
    TpPessoaGarFidjssria  string `xml:"TpPessoaGarFidjssria" json:"TpPessoaGarFidjssria"`
    CNPJ_CPFPessoaGarFidjssria string `xml:"CNPJ_CPFPessoaGarFidjssria" json:"CNPJ_CPFPessoaGarFidjssria"`
    GrupoGarFidjssoriaPF  GrupoGarFidjssoriaPF `xml:"Grupo_ACCC471_GarFidjssriaPF" json:"Grupo_ACCC471_GarFidjssriaPF"`
}
 
type GrupoGarFidjssoriaPF struct {
    NomPessoaGarFidjssria string `xml:"NomPessoaGarFidjssria" json:"NomPessoaGarFidjssria"`
    NomMae                string `xml:"NomMae" json:"NomMae"`
    EstadoCivil           string `xml:"EstadoCivil" json:"EstadoCivil"`
    CPFConjuge            string `xml:"CPFConjuge" json:"CPFConjuge"`
    NomConjuge            string `xml:"NomConjuge" json:"NomConjuge"`
    PortePessoaGarFidjssriaPF string `xml:"PortePessoaGarFidjssriaPF" json:"PortePessoaGarFidjssriaPF"`
    TpIdentc              string `xml:"TpIdentc" json:"TpIdentc"`
    NumIdentc             string `xml:"NumIdentc" json:"NumIdentc"`
}
 
func XMLToJSON(xmlData []byte) ([]byte, error) {
    var accc ACCCDOC
    if err := xml.Unmarshal(xmlData, &accc); err != nil {
        return nil, fmt.Errorf("erro ao deserializar XML: %v", err)
    }
 
    jsonData, err := json.MarshalIndent(accc, "", "  ")
    if err != nil {
        return nil, fmt.Errorf("erro ao serializar para JSON: %v", err)
    }
 
    return jsonData, nil
}
 
func JSONToXML(jsonData []byte) ([]byte, error) {
    var accc ACCCDOC
    if err := json.Unmarshal(jsonData, &accc); err != nil {
        return nil, fmt.Errorf("erro ao deserializar JSON: %v", err)
    }
 
    xmlData, err := xml.MarshalIndent(accc, "", "  ")
    if err != nil {
        return nil, fmt.Errorf("erro ao serializar para XML: %v", err)
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
 
    jsonData, err := XMLToJSON(xmlData)
    if err != nil {
        sendError(w, http.StatusInternalServerError, 1004, fmt.Sprintf("Erro ao converter XML para JSON: %v", err))
        return
    }
 
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)
}
 
func handleJSONToXML(w http.ResponseWriter, r *http.Request) {
    jsonData, err := ioutil.ReadAll(r.Body)
    if err != nil {
        sendError(w, http.StatusBadRequest, 1005, "Falha ao ler o JSON do corpo da requisição.")
        return
    }
 
    xmlData, err := JSONToXML(jsonData)
    if err != nil {
        sendError(w, http.StatusInternalServerError, 1006, fmt.Sprintf("Erro ao converter JSON para XML: %v", err))
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