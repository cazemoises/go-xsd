package xmlvalidator

import (
	"bytes"
	"fmt"
	"os/exec"
)

func ValidateXML(xmlData []byte, xsdPath string) error {

	cmd := exec.Command("xmllint", "--noout", "--schema", xsdPath, "-")
	cmd.Stdin = bytes.NewReader(xmlData)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("efro de validação: %s: %v", string(output), err)
	}

	return nil
}
