package tasty

type Customer struct {
	ID                              string              `json:"id"`
	PrefixName                      string              `json:"prefix-name"`
	FirstName                       string              `json:"first-name"`
	MiddleName                      string              `json:"middle-name"`
	LastName                        string              `json:"last-name"`
	SuffixName                      string              `json:"suffix-name"`
	FirstSurname                    string              `json:"first-surname"`
	SecondSurname                   string              `json:"second-surname"`
	Address                         Address             `json:"address"`
	MailingAddress                  Address             `json:"mailing-address"`
	CustomerSuitability             CustomerSuitability `json:"customer-suitability"`
	UsaCitizenshipType              string              `json:"usa-citizenship-type"`
	IsForeign                       bool                `json:"is-foreign"`
	MobilePhoneNumber               string              `json:"mobile-phone-number"`
	WorkPhoneNumber                 string              `json:"work-phone-number"`
	HomePhoneNumber                 string              `json:"home-phone-number"`
	Email                           string              `json:"email"`
	TaxNumberType                   string              `json:"tax-number-type"`
	TaxNumber                       string              `json:"tax-number"`
	ForeignTaxNumber                string              `json:"foreign-tax-number"`
	BirthDate                       string              `json:"birth-date"`
	ExternalID                      string              `json:"external-id"`
	CitizenshipCountry              string              `json:"citizenship-country"`
	BirthCountry                    string              `json:"birth-country"`
	VisaType                        string              `json:"visa-type"`
	VisaExpirationDate              string              `json:"visa-expiration-date"`
	SubjectToTaxWithholding         bool                `json:"subject-to-tax-withholding"`
	AgreedToMargining               bool                `json:"agreed-to-margining"`
	AgreedToTerms                   bool                `json:"agreed-to-terms"`
	SignatureOfAgreement            bool                `json:"signature-of-agreement"`
	Gender                          string              `json:"gender"`
	HasIndustryAffiliation          bool                `json:"has-industry-affiliation"`
	IndustryAffiliationFirm         string              `json:"industry-affiliation-firm"`
	HasPoliticalAffiliation         bool                `json:"has-political-affiliation"`
	PoliticalOrganization           string              `json:"political-organization"`
	FamilyMemberNames               string              `json:"family-member-names"`
	HasListedAffiliation            bool                `json:"has-listed-affiliation"`
	ListedAffiliationSymbol         string              `json:"listed-affiliation-symbol"`
	IsInvestmentAdviser             bool                `json:"is-investment-adviser"`
	HasInstitutionalAssets          bool                `json:"has-institutional-assets"`
	DeskCustomerID                  string              `json:"desk-customer-id"`
	UserID                          string              `json:"user-id"`
	IsProfessional                  bool                `json:"is-professional"`
	HasDelayedQuotes                bool                `json:"has-delayed-quotes"`
	HasPendingOrApprovedApplication bool                `json:"has-pending-or-approved-application"`
	PermittedAccountTypes           []AccountType       `json:"permitted-account-types"`
	IdentifiableType                string              `json:"identifiable-type"`
	Person                          Person              `json:"person"`
	Entity                          Entity              `json:"entity"`
}

type Address struct {
	StreetOne   string `json:"street-one"`
	StreetTwo   string `json:"street-two"`
	StreetThree string `json:"street-three"`
	City        string `json:"city"`
	StateRegion string `json:"state-region"`
	PostalCode  string `json:"postal-code"`
	Country     string `json:"country"`
	IsForeign   bool   `json:"is-foreign"`
	IsDomestic  bool   `json:"is-domestic"`
}

type CustomerSuitability struct {
	ID                                int    `json:"id"`
	CustomerID                        int    `json:"customer-id"`
	MaritalStatus                     string `json:"marital-status"`
	NumberOfDependents                int    `json:"number-of-dependents"`
	EmploymentStatus                  string `json:"employment-status"`
	Occupation                        string `json:"occupation"`
	EmployerName                      string `json:"employer-name"`
	JobTitle                          string `json:"job-title"`
	AnnualNetIncome                   int    `json:"annual-net-income"`
	NetWorth                          int    `json:"net-worth"`
	LiquidNetWorth                    int    `json:"liquid-net-worth"`
	StockTradingExperience            string `json:"stock-trading-experience"`
	CoveredOptionsTradingExperience   string `json:"covered-options-trading-experience"`
	UncoveredOptionsTradingExperience string `json:"uncovered-options-trading-experience"`
	FuturesTradingExperience          string `json:"futures-trading-experience"`
	TaxBracket                        string `json:"tax-bracket"`
}

type Person struct {
	ExternalID         string `json:"external-id"`
	PrefixName         string `json:"prefix-name"`
	FirstName          string `json:"first-name"`
	MiddleName         string `json:"middle-name"`
	LastName           string `json:"last-name"`
	SuffixName         string `json:"suffix-name"`
	BirthDate          string `json:"birth-date"`
	BirthCountry       string `json:"birth-country"`
	CitizenshipCountry string `json:"citizenship-country"`
	UsaCitizenshipType string `json:"usa-citizenship-type"`
	VisaType           string `json:"visa-type"`
	VisaExpirationDate string `json:"visa-expiration-date"`
	MaritalStatus      string `json:"marital-status"`
	NumberOfDependents int    `json:"number-of-dependents"`
	EmploymentStatus   string `json:"employment-status"`
	Occupation         string `json:"occupation"`
	EmployerName       string `json:"employer-name"`
	JobTitle           string `json:"job-title"`
}

type Entity struct {
	ID                               string            `json:"id"`
	LegalName                        string            `json:"legal-name"`
	TaxNumber                        string            `json:"tax-number"`
	IsDomestic                       bool              `json:"is-domestic"`
	EntityType                       string            `json:"entity-type"`
	Email                            string            `json:"email"`
	PhoneNumber                      string            `json:"phone-number"`
	BusinessNature                   string            `json:"business-nature"`
	HasForeignInstitutionAffiliation bool              `json:"has-foreign-institution-affiliation"`
	ForeignInstitution               string            `json:"foreign-institution"`
	HasForeignBankAffiliation        bool              `json:"has-foreign-bank-affiliation"`
	GrantorFirstName                 string            `json:"grantor-first-name"`
	GrantorMiddleName                string            `json:"grantor-middle-name"`
	GrantorLastName                  string            `json:"grantor-last-name"`
	GrantorEmail                     string            `json:"grantor-email"`
	GrantorBirthDate                 string            `json:"grantor-birth-date"`
	GrantorTaxNumber                 string            `json:"grantor-tax-number"`
	Address                          Address           `json:"address"`
	EntitySuitability                EntitySuitability `json:"entity-suitability"`
	EntityOfficers                   []EntityOfficer   `json:"entity-officers"`
}

type EntitySuitability struct {
	ID                                string `json:"id"`
	EntityID                          int    `json:"entity-id"`
	AnnualNetIncome                   int    `json:"annual-net-income"`
	NetWorth                          int    `json:"net-worth"`
	LiquidNetWorth                    int    `json:"liquid-net-worth"`
	StockTradingExperience            string `json:"stock-trading-experience"`
	TaxBracket                        string `json:"tax-bracket"`
	CoveredOptionsTradingExperience   string `json:"covered-options-trading-experience"`
	UncoveredOptionsTradingExperience string `json:"uncovered-options-trading-experience"`
	FuturesTradingExperience          string `json:"futures-trading-experience"`
}

type EntityOfficer struct {
	ID                   string  `json:"id"`
	ExternalID           string  `json:"external-id"`
	PrefixName           string  `json:"prefix-name"`
	FirstName            string  `json:"first-name"`
	MiddleName           string  `json:"middle-name"`
	LastName             string  `json:"last-name"`
	SuffixName           string  `json:"suffix-name"`
	BirthDate            string  `json:"birth-date"`
	BirthCountry         string  `json:"birth-country"`
	CitizenshipCountry   string  `json:"citizenship-country"`
	UsaCitizenshipType   string  `json:"usa-citizenship-type"`
	IsForeign            bool    `json:"is-foreign"`
	VisaType             string  `json:"visa-type"`
	VisaExpirationDate   string  `json:"visa-expiration-date"`
	TaxNumberType        string  `json:"tax-number-type"`
	TaxNumber            string  `json:"tax-number"`
	MobilePhoneNumber    string  `json:"mobile-phone-number"`
	HomePhoneNumber      string  `json:"home-phone-number"`
	WorkPhoneNumber      string  `json:"work-phone-number"`
	Email                string  `json:"email"`
	MaritalStatus        string  `json:"marital-status"`
	NumberOfDependents   int     `json:"number-of-dependents"`
	EmploymentStatus     string  `json:"employment-status"`
	EmployerName         string  `json:"employer-name"`
	Occupation           string  `json:"occupation"`
	JobTitle             string  `json:"job-title"`
	RelationshipToEntity string  `json:"relationship-to-entity"`
	OwnerOfRecord        bool    `json:"owner-of-record"`
	Address              Address `json:"address"`
}
