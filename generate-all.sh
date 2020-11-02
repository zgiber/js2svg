#!/usr/bin/env bash
declare -A sources
declare -A diagrams

sources=( 
    ["accounts"]="https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/master/dist/openapi/account-info-openapi.yaml"
    ["payments"]="https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/master/dist/openapi/payment-initiation-openapi.yaml"
    ["funds"]="https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/master/dist/openapi/confirmation-funds-openapi.yaml"
    ["events"]="https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/master/dist/swagger/events-swagger.yaml"
    ["notifications"]="https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/master/dist/swagger/event-notifications-swagger.yaml"
)

diagrams["accounts"]=( "OBReadConsent1" "OBReadConsentResponse1" "OBReadAccount6" "OBReadBalance1" )
diagrams["accounts"]+=( "OBReadTransaction6" "OBReadBeneficiary5" "OBReadDirectDebit2" "OBReadStandingOrder6" "OBReadOffer1" )
diagrams["accounts"]+=( "OBParty2" "OBReadParty3" "OBReadScheduledPayment3" "OBReadStatement2" )
diagrams["accounts"]+=( "OBReadProduct2" "CreditInterest" "Overdraft" "Product" )

for d in "${!diagrams[@]}" do
    echo "$d"
    echo "${diagrams[$d]}"
done

function fetchSources() {
    for i in "${!sources[@]}"
    do
        if curl -sL --fail "${sources[$i]}" -o "/tmp/$i-openapi.yaml"; then
            echo "fetched $i"
        else
            echo "failed to fetch $i"
            exit 1
        fi
    done
}

function cleanup() {
    for i in "${!sources[@]}"
    do
        rm "/tmp/$i-openapi.yaml"
    done
}
 
 
# fetchSources
exit 0


paymentDiagrams=( "OBRisk1" "OBCharge2" "OBAuthorisation1" "OBMultiAuthorisation1" "OBDomesticRefundAccount1" "OBInternationalRefundAccount1" "OBWritePaymentDetails1" "OBSCASupportData1" ) 
paymentDiagrams+=( "OBDomestic2" "OBWriteDomesticConsent4" "OBWriteDomesticConsentResponse5" "OBWriteFundsConfirmationResponse1" "OBWriteDomestic2" "OBWriteDomesticResponse5" "OBWritePaymentDetailsResponse1" )
paymentDiagrams+=( "OBDomesticScheduled2" "OBWriteDomesticScheduledConsent4" "OBWriteDomesticScheduledConsentResponse5" "OBWriteDomesticScheduled2" "OBWriteDomesticScheduledResponse5")
paymentDiagrams+=( "OBDomesticStandingOrder3" "OBWriteDomesticStandingOrderConsent5" "OBWriteDomesticStandingOrderConsentResponse6" )
paymentDiagrams+=( "OBWriteDomesticStandingOrder3" "OBWriteDomesticStandingOrderResponse6" )
paymentDiagrams+=( "OBInternational3" "OBExchangeRate2" "OBWriteInternationalConsent5" "OBWriteInternationalConsentResponse6" )
paymentDiagrams+=( "OBWriteInternational3" "OBWriteInternationalResponse5" )
paymentDiagrams+=( "OBInternationalScheduled3" "OBWriteInternationalScheduledConsent5" "OBWriteInternationalScheduledConsentResponse6" )
paymentDiagrams+=( "OBWriteInternationalScheduled3" "OBWriteInternationalScheduledResponse6" )
paymentDiagrams+=( "OBInternationalStandingOrder4" "OBWriteInternationalStandingOrderConsent6" "OBWriteInternationalStandingOrderConsentResponse7" )
paymentDiagrams+=( "OBWriteInternationalStandingOrder4" "OBWriteInternationalStandingOrderResponse7" )
paymentDiagrams+=( "OBFile2" "OBWriteFileConsent3" "OBWriteFileConsentResponse4" )
paymentDiagrams+=( "OBWriteFile2" "OBWriteFileResponse3" )

fundsConfirmationDiagrams=( "OBFundsConfirmationConsent1" "OBFundsConfirmationConsentResponse1" )
fundsConfirmationDiagrams+=( "OBFundsConfirmation1" "OBFundsConfirmationResponse1" )

# generate class diagrams for accounts
for name in ${accountDiagrams[@]}; do
  runtime=$(./generate -src=file:///tmp/account-info-openapi.yaml --threads $name)
  allRuntimes+=( $runtime )
done

# generate class diagrams for payments

# generate class diagrams for funds confirmation

rm ./tmp/account-info-openapi.yaml
rm ./tmp/payment-initiation-openapi.yaml
rm ./tmp/confirmation-funds-openapi.yaml
rm ./tmp/events-openapi.yaml
rm ./tmp/event-notifications-openapi.yaml