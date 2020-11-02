package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/zgiber/js2svg"
)

type diagramSet struct {
	url           string
	diagrams      []string
	parsedSchemas map[string]map[string]interface{}
}

var (
	sources = map[string]diagramSet{
		"accounts": {
			url: "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/master/dist/openapi/account-info-openapi.yaml",
			diagrams: []string{
				"OBReadConsent1", "OBReadConsentResponse1", "OBReadAccount6", "OBReadBalance1",
				"OBReadTransaction6", "OBReadBeneficiary5", "OBReadDirectDebit2", "OBReadStandingOrder6", "OBReadOffer1",
				"OBParty2", "OBReadParty3", "OBReadScheduledPayment3", "OBReadStatement2",
				"OBReadProduct2", "CreditInterest", "Overdraft", "Product",
			},
		},
		"payments": {
			url: "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/master/dist/openapi/payment-initiation-openapi.yaml",
			diagrams: []string{
				"OBRisk1", "OBCharge2", "OBAuthorisation1", "OBMultiAuthorisation1", "OBDomesticRefundAccount1", "OBInternationalRefundAccount1", "OBWritePaymentDetails1", "OBSCASupportData1",
				"OBDomestic2", "OBWriteDomesticConsent4", "OBWriteDomesticConsentResponse5", "OBWriteFundsConfirmationResponse1", "OBWriteDomestic2", "OBWriteDomesticResponse5", "OBWritePaymentDetailsResponse1",
				"OBDomesticScheduled2", "OBWriteDomesticScheduledConsent4", "OBWriteDomesticScheduledConsentResponse5", "OBWriteDomesticScheduled2", "OBWriteDomesticScheduledResponse5",
				"OBDomesticStandingOrder3", "OBWriteDomesticStandingOrderConsent5", "OBWriteDomesticStandingOrderConsentResponse6",
				"OBWriteDomesticStandingOrder3", "OBWriteDomesticStandingOrderResponse6",
				"OBInternational3", "OBExchangeRate2", "OBWriteInternationalConsent5", "OBWriteInternationalConsentResponse6",
				"OBWriteInternational3", "OBWriteInternationalResponse5",
				"OBInternationalScheduled3", "OBWriteInternationalScheduledConsent5", "OBWriteInternationalScheduledConsentResponse6",
				"OBWriteInternationalScheduled3", "OBWriteInternationalScheduledResponse6",
				"OBInternationalStandingOrder4", "OBWriteInternationalStandingOrderConsent6", "OBWriteInternationalStandingOrderConsentResponse7",
				"OBWriteInternationalStandingOrder4", "OBWriteInternationalStandingOrderResponse7",
				"OBFile2", "OBWriteFileConsent3", "OBWriteFileConsentResponse4",
				"OBWriteFile2", "OBWriteFileResponse3",
			},
		},
		"funds": {
			url: "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/master/dist/openapi/confirmation-funds-openapi.yaml",
			diagrams: []string{
				"OBFundsConfirmationConsent1", "OBFundsConfirmationConsentResponse1",
				"OBFundsConfirmation1", "OBFundsConfirmationResponse1",
			},
		},
	}
)

func main() {
	wg := &sync.WaitGroup{}
	for name, src := range sources {
		wg.Add(1)
		n, s := name, src

		go func() {
			err := renderDiagrams(n, s.url, s.diagrams)
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func renderDiagrams(name, url string, diagrams []string) error {
	schema, err := openRemoteSrc(url)
	if err != nil {
		return err
	}
	defer schema.Close()

	parsedSchema, err := js2svg.ParseToMap(schema, "components.schemas")
	if err != nil {
		return err
	}

	// create directory for the rendered diagrams
	cwd, _ := os.Getwd()
	dir := path.Join(cwd, name)
	if err := os.Mkdir(dir, 0755); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	// write the rendered diagrams to files
	for _, diagramName := range diagrams {
		objectSchema, ok := parsedSchema[diagramName].(map[string]interface{})
		if !ok {
			log.Printf("can't find object in %s: 'components.schemas.%s'", name, diagramName)
			continue
		}

		diagram, err := js2svg.MakeDiagram(objectSchema, diagramName)
		if err != nil {
			return err
		}

		dst, err := os.Create(path.Join(dir, diagramName+".svg"))
		if err != nil {
			return err
		}

		if err != diagram.Render(dst) {
			dst.Close()
			return err
		}
		dst.Close()
	}

	return nil
}

func openRemoteSrc(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("remote returned status %v", resp.StatusCode)
	}

	return resp.Body, nil
}
