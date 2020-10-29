package main

import (
	"os"

	svg "github.com/zgiber/js2svg"
)

var (

	// all of these would be parsed from a jsonschema with proper description:
	DeliveryAddress = svg.Object{
		Name: "DeliveryAddress",
		Properties: []svg.Property{
			{Name: "AddressLine", Relationship: "0..2", Description: "Description could show validation rules for example."},
			{Name: "StreetName", Relationship: "0..1", Description: "Description could show validation rules for example."},
			{Name: "BuildingNumber", Relationship: "0..1", Description: "Description could show validation rules for example."},
			{Name: "PostCode", Relationship: "0..1", Description: "Description could show validation rules for example."},
			{Name: "TownName", Relationship: "0..1", Description: "Description could show validation rules for example."},
			{Name: "CountrySubdivision", Relationship: "0..1", Description: "Description could show validation rules for example."},
			{Name: "Country", Relationship: "0..1", Description: "Description could show validation rules for example."},
		},
	}

	Risk = svg.Object{
		Name: "RemittanceInformation",
		Properties: []svg.Property{
			{Name: "PaymentContextCode", Relationship: "0..1", Description: "Description could show validation rules for example."},
			{Name: "MerchantCategoryCode", Relationship: "0..1", Description: "Description could show validation rules for example."},
			{Name: "MerchantCustomerIdentification", Relationship: "0..1", Description: "Description could show validation rules for example."},
		},
		ComposedOf: []svg.Composition{
			{Relationship: "1..1", Object: &DeliveryAddress},
		},
	}

	RemittanceInformation = svg.Object{
		Name: "RemittanceInformation",
		Properties: []svg.Property{
			{Name: "Unstructured", Relationship: "0..1", Description: "Description could show validation rules for example."},
			{Name: "Reference", Relationship: "0..1", Description: "Description could show validation rules for example."},
		},
	}

	SupplementaryData = svg.Object{
		Name: "SupplementaryData",
	}

	PostalAddressCreditor = svg.Object{
		Name: "PostalAddress",
		Properties: []svg.Property{
			{Name: "AddressType", Relationship: "0..1", Description: "Description could show validation rules for example."},
			{Name: "Department", Relationship: "0..1", Description: "Description could show validation rules for example."},
			{Name: "SubDepartment", Relationship: "0..1", Description: "Description could show validation rules for example."},
			{Name: "StreetName", Relationship: "0..1", Description: "Description could show validation rules for example."},
			{Name: "BuildingNumber", Relationship: "0..1", Description: "Description could show validation rules for example."},
			{Name: "PostCode", Relationship: "0..1", Description: "Description could show validation rules for example."},
			{Name: "TownName", Relationship: "0..1", Description: "Description could show validation rules for example."},
			{Name: "CountrySubdivision", Relationship: "0..1", Description: "Description could show validation rules for example."},
			{Name: "Country", Relationship: "0..1", Description: "Description could show validation rules for example."},
			{Name: "AddressLine", Relationship: "0..7", Description: "Description could show validation rules for example."},
		},
	}

	PostalAddressCreditorAgent = PostalAddressCreditor

	CreditorAgent = svg.Object{
		Name: "CreditorAgent",
		Properties: []svg.Property{
			{Name: "SchemeName", Relationship: "0..1", Description: "Description could show validation rules for example."},
			{Name: "Identification", Relationship: "0..1", Description: "Description could show validation rules for example."},
			{Name: "Name", Relationship: "0..1", Description: "Description could show validation rules for example."},
		},
		ComposedOf: []svg.Composition{
			{Relationship: "0..1", Object: &PostalAddressCreditorAgent},
		},
	}

	Creditor = svg.Object{
		Name: "Creditor",
		Properties: []svg.Property{
			{Name: "Name", Relationship: "1..1", Description: "Description could show validation rules for example."},
		},
		ComposedOf: []svg.Composition{
			{Relationship: "0..1", Object: &PostalAddressCreditor},
		},
	}

	DebtorAccount = svg.Object{
		Name: "DebtorAccount",
		Properties: []svg.Property{
			{Name: "SchemeName", Relationship: "1..1", Description: "Description could show validation rules for example."},
			{Name: "Identification", Relationship: "1..1", Description: "Description could show validation rules for example."},
			{Name: "Name", Relationship: "1..1", Description: "Description could show validation rules for example."},
			{Name: "SecondaryIdentification", Relationship: "1..1", Description: "Description could show validation rules for example."},
		},
	}

	CreditorAccount = svg.Object{
		Name: "DebtorAccount",
		Properties: []svg.Property{
			{Name: "SchemeName", Relationship: "1..1", Description: "Description could show validation rules for example."},
			{Name: "Identification", Relationship: "1..1", Description: "Description could show validation rules for example."},
			{Name: "Name", Relationship: "1..1", Description: "Description could show validation rules for example."},
			{Name: "SecondaryIdentification", Relationship: "1..1", Description: "Description could show validation rules for example."},
		},
	}

	ExchangeRateInformation = svg.Object{
		Name: "ExchangeRateInformation",
		Properties: []svg.Property{
			{Name: "UnitCurrency", Relationship: "1..1", Description: "Description could show validation rules for example."},
			{Name: "ExchangeRate", Relationship: "1..1", Description: "Description could show validation rules for example."},
			{Name: "RateType", Relationship: "1..1", Description: "Description could show validation rules for example."},
			{Name: "ContractIdentification", Relationship: "1..1", Description: "Description could show validation rules for example."},
		},
	}

	InstructedAmount = svg.Object{
		Name: "InstructedAmount",
		Properties: []svg.Property{
			{Name: "Amount", Relationship: "1..1", Description: "Description could show validation rules for example."},
			{Name: "Currency", Relationship: "1..1", Description: "Description could show validation rules for example."},
		},
	}

	Initiation = svg.Object{
		Name: "Initiation",
		Properties: []svg.Property{
			{Name: "InstructionIdentification", Relationship: "1..1", Description: "Description could show type and validation rules for example."},
			{Name: "EndToEndIdentification", Relationship: "1..1", Description: ""},
			{Name: "LocalInstrument", Relationship: "1..0", Description: ""},
			{Name: "InstructionPriority", Relationship: "1..0", Description: ""},
			{Name: "Purpose", Relationship: "1..0", Description: ""},
			{Name: "ExtendedPurpose", Relationship: "1..0", Description: ""},
			{Name: "ChargeBearer", Relationship: "1..0", Description: ""},
			{Name: "CurrencyOfTransfer", Relationship: "1..0", Description: ""},
			{Name: "DestinationCountryCode", Relationship: "1..0", Description: ""},
		},
		ComposedOf: []svg.Composition{
			{Relationship: "1..1", Object: &InstructedAmount},
			{Relationship: "1..0", Object: &ExchangeRateInformation},
			{Relationship: "1..1", Object: &DebtorAccount},
			{Relationship: "1..1", Object: &Creditor},
			{Relationship: "1..1", Object: &CreditorAgent},
			{Relationship: "1..1", Object: &CreditorAccount},
			{Relationship: "1..1", Object: &RemittanceInformation},
			{Relationship: "1..1", Object: &SupplementaryData},
		},
	}
)

func main() {

	Data := svg.Object{
		Name: "Data",
		Properties: []svg.Property{
			{Name: "ConsentId", Relationship: "1..1", Description: "The ID of the PSU consent for this operation."},
		},
		ComposedOf: []svg.Composition{
			{Relationship: "1..1", Object: &Initiation},
		},
	}

	OBWriteInternational3 := svg.Object{
		Name:        "OBWriteInternational3",
		Description: "The response object for xxx call... whatever",
		ComposedOf: []svg.Composition{
			{Relationship: "1..1", Object: &Data},
			{Relationship: "1..1", Object: &Risk},
		},
	}

	d := svg.Diagram{Root: &OBWriteInternational3}
	d.Render(os.Stdout)
}
