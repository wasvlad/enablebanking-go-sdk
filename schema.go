package eb

type BalanceResource struct {
	Name                    string        `json:"name"`
	BalanceAmount           AmountType    `json:"balance_amount"`
	BalanceType             BalanceStatus `json:"balance_type"`
	LastChangeDateTime      *string       `json:"last_change_date_time,omitempty"`
	ReferenceDate           *string       `json:"reference_date,omitempty"`
	LastCommitedTransaction *string       `json:"last_commited_transaction,omitempty"`
}

type AmountType struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type BalanceStatus int

const (
	CLAV BalanceStatus = iota
	CLBD
	FWAV
	INFO
	ITAV
	ITBD
	OPAV
	OPBD
	OTHR
	PRCD
	VALU
	XPCD
)

type Transaction struct {
	EntryReference                          *string                             `json:"entry_reference"`
	MerchantCategoryCode                    *string                             `json:"merchant_category_code"`
	TransactionAmount                       *AmountType                         `json:"transaction_amount"`
	Creditor                                *PartyIdentification                `json:"creditor"`
	CreditorAccount                         *AccountIdentification              `json:"creditor_account"`
	CreditorAgent                           *FinancialInstitutionIdentification `json:"creditor_agent"`
	Debtor                                  *PartyIdentification                `json:"debtor"`
	DebtorAccount                           *AccountIdentification              `json:"debtor_account"`
	DebtorAgent                             *FinancialInstitutionIdentification `json:"debtor_agent"`
	BankTransactionCode                     *BankTransactionCode                `json:"bank_transaction_code"`
	CreditDebitIndicator                    *CreditDebitIndicator               `json:"credit_debit_indicator"`
	Status                                  *TransactionStatus                  `json:"status"`
	BookingDate                             *string                             `json:"booking_date"`
	ValueDate                               *string                             `json:"value_date"`
	TransactionDate                         *string                             `json:"transaction_date"`
	BalanceAfterTransaction                 *AmountType                         `json:"balance_after_transaction"`
	ReferenceNumber                         string                              `json:"reference_number"`
	ReferenceNumberSchema                   *ReferenceNumberScheme              `json:"reference_number_schema"`
	RemmitanceInformation                   []string                            `json:"remmitance_information"`
	DebtorAccountAdditionalIdentification   []GenericIdentification             `json:"debtor_account_additional_identification"`
	CreditorAccountAdditionalIdentification []GenericIdentification             `json:"creditor_account_additional_identification"`
	ExchangeRate                            *ExchangeRate                       `json:"exchange_rate"`
	Note                                    *string                             `json:"note"`
	TransactionId                           *string                             `json:"transaction_id"`
}

type PartyIdentification struct {
	Name           *string                `json:"name"`
	PostalAddress  *PostalAddress         `json:"postal_address"`
	OrganisationId *GenericIdentification `json:"organisation_id"`
	PrivateId      *GenericIdentification `json:"private_id"`
	ContactDetails *ContactDetails        `json:"contact_details"`
}

type AccountIdentification struct {
	Iban  *string                `json:"IBAN"`
	Other *GenericIdentification `json:"other"`
}

type FinancialInstitutionIdentification struct {
	BicFi                  *string                             `json:"bic_fi"`
	ClearingSystemMemberId *ClearingSystemMemberIdentification `json:"clearing_system_member_id"`
	Name                   *string                             `json:"name"`
}

type BankTransactionCode struct {
	Description *string `json:"description"`
	Code        *string `json:"code"`
	SubCode     *string `json:"sub_code"`
}

type CreditDebitIndicator int

const (
	CRDT CreditDebitIndicator = iota
	DBIT
)

type TransactionStatus int

const (
	TxBOOK TransactionStatus = iota
	TxTNCL
	TxHOLD
	TxOTHR
	TxPDNG
	TxRJCT
	TxSCHD
)

type ReferenceNumberScheme int

const (
	BERF ReferenceNumberScheme = iota
	FIRF
	INTL
	NORF
	SDDM
	SEBG
)

type GenericIdentification struct {
	Identification string     `json:"identification"`
	SchemeName     SchemeName `json:"scheme_name"`
	Issuer         *string    `json:"issuer"`
}

type ExchangeRate struct {
	UnitCurrency           *CurrencyCode `json:"unit_currency"`
	ExchangeRate           *string       `json:"exchange_rate"`
	RateType               *RateType     `json:"rate_type"`
	ContractIdentification *string       `json:"contract_identification"`
	InstructedAmount       *AmountType   `json:"instructed_amount"`
}

type ClearingSystemMemberIdentification struct {
	ClearingSystemId *string `json:"clearing_system_id"`
	MemberId         *string `json:"member_id"`
}

type ContactDetails struct {
	EmailDetails *string `json:"email_detaild"`
	PhoneNumber  *string `json:"phone_number"`
}

type PostalAddress struct {
	AddressType        *AddressType `json:"address_type"`
	Department         *string      `json:"department"`
	SubDepartment      *string      `json:"sub_department"`
	StreetName         *string      `json:"street_name"`
	BuildingNumber     *string      `json:"building_number"`
	PostCode           *string      `json:"post_code"`
	TownName           *string      `json:"town_name"`
	CountrySubDivision *string      `json:"country_sub_division"`
	Country            *string      `json:"country"`
	AddressLine        []string     `json:"address_line"`
}

type AddressType int

const (
	Business AddressType = iota
	Correspondence
	DeliveryTo
	MailTo
	POBox
	Postal
	Residential
	Statement
)

type SchemeName int

type CurrencyCode string

type RateType int

const (
	AGRD RateType = iota
	SALE
	SPOT
)
