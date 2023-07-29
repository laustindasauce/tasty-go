package tasty //nolint:testpackage // testing private field

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetMyCustomerInfo(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/customers/me", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, getCustomerResp)
	})

	resp, httpResp, err := client.GetMyCustomerInfo()
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, "me", resp.ID)
	require.Equal(t, "Austin", resp.FirstName)
	require.Equal(t, "Spencer", resp.LastName)
	require.Equal(t, "Spencer", resp.FirstSurname)

	addr := resp.Address

	require.Equal(t, "1000 West Fulton St", addr.StreetOne)
	require.Equal(t, "Chicago", addr.City)
	require.Equal(t, "IL", addr.StateRegion)
	require.Equal(t, "60607", addr.PostalCode)
	require.Equal(t, "US", addr.Country)
	require.False(t, addr.IsForeign)
	require.True(t, addr.IsDomestic)

	addr = resp.MailingAddress

	require.Equal(t, "1000 West Fulton St", addr.StreetOne)
	require.Equal(t, "Chicago", addr.City)
	require.Equal(t, "IL", addr.StateRegion)
	require.Equal(t, "60607", addr.PostalCode)
	require.Equal(t, "US", addr.Country)
	require.False(t, addr.IsForeign)
	require.True(t, addr.IsDomestic)

	cs := resp.CustomerSuitability

	require.Equal(t, 115555, cs.ID)
	require.Equal(t, "MARRIED", cs.MaritalStatus)
	require.Equal(t, 0, cs.NumberOfDependents)
	require.Equal(t, "EMPLOYED", cs.EmploymentStatus)
	require.Equal(t, "Technology", cs.Occupation)
	require.Equal(t, "Unknown", cs.EmployerName)
	require.Equal(t, "Senior AI & ML Engineer", cs.JobTitle)
	require.Equal(t, 10, cs.AnnualNetIncome)
	require.Equal(t, 0, cs.NetWorth)
	require.Equal(t, 0, cs.LiquidNetWorth)
	require.Equal(t, "EXTENSIVE", cs.StockTradingExperience)
	require.Equal(t, "EXTENSIVE", cs.CoveredOptionsTradingExperience)
	require.Equal(t, "EXTENSIVE", cs.UncoveredOptionsTradingExperience)
	require.Equal(t, "LIMITED", cs.FuturesTradingExperience)

	require.Equal(t, "Citizen", resp.UsaCitizenshipType)
	require.False(t, resp.IsForeign)
	require.Equal(t, "+15555555555", resp.MobilePhoneNumber)
	require.Equal(t, "", resp.WorkPhoneNumber)
	require.Equal(t, "", resp.HomePhoneNumber)
	require.Equal(t, "me@austinbspencer.com", resp.Email)
	require.Equal(t, "SSN", resp.TaxNumberType)
	require.Equal(t, "*****5555", resp.TaxNumber)
	require.Equal(t, "1900-01-01", resp.BirthDate)
	require.Equal(t, "C0001234567", resp.ExternalID)
	require.Equal(t, "USA", resp.CitizenshipCountry)
	require.False(t, resp.SubjectToTaxWithholding)
	require.True(t, resp.AgreedToMargining)
	require.True(t, resp.AgreedToTerms)
	require.False(t, resp.HasIndustryAffiliation)
	require.False(t, resp.HasPoliticalAffiliation)
	require.False(t, resp.HasListedAffiliation)
	require.False(t, resp.IsInvestmentAdviser)
	require.False(t, resp.HasInstitutionalAssets)
	require.False(t, resp.IsProfessional)
	require.False(t, resp.HasDelayedQuotes)
	require.True(t, resp.HasPendingOrApprovedApplication)
	require.Equal(t, 6, len(resp.PermittedAccountTypes))

	acct := resp.PermittedAccountTypes[0]

	require.Equal(t, "Individual", acct.Name)
	require.Equal(t, "", acct.Description)
	require.False(t, acct.IsTaxAdvantaged)
	require.False(t, acct.HasMultipleOwners)
	require.True(t, acct.IsPubliclyAvailable)
	require.Equal(t, 5, len(acct.MarginTypes))

	mt := acct.MarginTypes[0]

	require.Equal(t, "Cash", mt.Name)
	require.False(t, mt.IsMargin)

	require.Equal(t, "Person", resp.IdentifiableType)

	p := resp.Person

	require.Equal(t, "P0002345678", p.ExternalID)
	require.Equal(t, "Austin", p.FirstName)
	require.Equal(t, "Spencer", p.LastName)
	require.Equal(t, "1900-01-01", p.BirthDate)
	require.Equal(t, "USA", p.CitizenshipCountry)
	require.Equal(t, "Citizen", p.UsaCitizenshipType)
	require.Equal(t, "MARRIED", p.MaritalStatus)
	require.Equal(t, 0, p.NumberOfDependents)
	require.Equal(t, "EMPLOYED", p.EmploymentStatus)
	require.Equal(t, "Technology", p.Occupation)
	require.Equal(t, "Unknown", p.EmployerName)
	require.Equal(t, "Senior AI & ML Engineer", p.JobTitle)
}

func TestGetMyCustomerInfoError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/customers/me", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, httpResp, err := client.GetMyCustomerInfo()
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestGetCustomer(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/customers/me", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, getCustomerResp)
	})

	resp, httpResp, err := client.GetCustomer("me")
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, "me", resp.ID)
	require.Equal(t, "Austin", resp.FirstName)
	require.Equal(t, "Spencer", resp.LastName)
	require.Equal(t, "Spencer", resp.FirstSurname)

	addr := resp.Address

	require.Equal(t, "1000 West Fulton St", addr.StreetOne)
	require.Equal(t, "Chicago", addr.City)
	require.Equal(t, "IL", addr.StateRegion)
	require.Equal(t, "60607", addr.PostalCode)
	require.Equal(t, "US", addr.Country)
	require.False(t, addr.IsForeign)
	require.True(t, addr.IsDomestic)

	addr = resp.MailingAddress

	require.Equal(t, "1000 West Fulton St", addr.StreetOne)
	require.Equal(t, "Chicago", addr.City)
	require.Equal(t, "IL", addr.StateRegion)
	require.Equal(t, "60607", addr.PostalCode)
	require.Equal(t, "US", addr.Country)
	require.False(t, addr.IsForeign)
	require.True(t, addr.IsDomestic)

	cs := resp.CustomerSuitability

	require.Equal(t, 115555, cs.ID)
	require.Equal(t, "MARRIED", cs.MaritalStatus)
	require.Equal(t, 0, cs.NumberOfDependents)
	require.Equal(t, "EMPLOYED", cs.EmploymentStatus)
	require.Equal(t, "Technology", cs.Occupation)
	require.Equal(t, "Unknown", cs.EmployerName)
	require.Equal(t, "Senior AI & ML Engineer", cs.JobTitle)
	require.Equal(t, 10, cs.AnnualNetIncome)
	require.Equal(t, 0, cs.NetWorth)
	require.Equal(t, 0, cs.LiquidNetWorth)
	require.Equal(t, "EXTENSIVE", cs.StockTradingExperience)
	require.Equal(t, "EXTENSIVE", cs.CoveredOptionsTradingExperience)
	require.Equal(t, "EXTENSIVE", cs.UncoveredOptionsTradingExperience)
	require.Equal(t, "LIMITED", cs.FuturesTradingExperience)

	require.Equal(t, "Citizen", resp.UsaCitizenshipType)
	require.False(t, resp.IsForeign)
	require.Equal(t, "+15555555555", resp.MobilePhoneNumber)
	require.Equal(t, "", resp.WorkPhoneNumber)
	require.Equal(t, "", resp.HomePhoneNumber)
	require.Equal(t, "me@austinbspencer.com", resp.Email)
	require.Equal(t, "SSN", resp.TaxNumberType)
	require.Equal(t, "*****5555", resp.TaxNumber)
	require.Equal(t, "1900-01-01", resp.BirthDate)
	require.Equal(t, "C0001234567", resp.ExternalID)
	require.Equal(t, "USA", resp.CitizenshipCountry)
	require.False(t, resp.SubjectToTaxWithholding)
	require.True(t, resp.AgreedToMargining)
	require.True(t, resp.AgreedToTerms)
	require.False(t, resp.HasIndustryAffiliation)
	require.False(t, resp.HasPoliticalAffiliation)
	require.False(t, resp.HasListedAffiliation)
	require.False(t, resp.IsInvestmentAdviser)
	require.False(t, resp.HasInstitutionalAssets)
	require.False(t, resp.IsProfessional)
	require.False(t, resp.HasDelayedQuotes)
	require.True(t, resp.HasPendingOrApprovedApplication)
	require.Equal(t, 6, len(resp.PermittedAccountTypes))

	acct := resp.PermittedAccountTypes[0]

	require.Equal(t, "Individual", acct.Name)
	require.Equal(t, "", acct.Description)
	require.False(t, acct.IsTaxAdvantaged)
	require.False(t, acct.HasMultipleOwners)
	require.True(t, acct.IsPubliclyAvailable)
	require.Equal(t, 5, len(acct.MarginTypes))

	mt := acct.MarginTypes[0]

	require.Equal(t, "Cash", mt.Name)
	require.False(t, mt.IsMargin)

	require.Equal(t, "Person", resp.IdentifiableType)

	p := resp.Person

	require.Equal(t, "P0002345678", p.ExternalID)
	require.Equal(t, "Austin", p.FirstName)
	require.Equal(t, "Spencer", p.LastName)
	require.Equal(t, "1900-01-01", p.BirthDate)
	require.Equal(t, "USA", p.CitizenshipCountry)
	require.Equal(t, "Citizen", p.UsaCitizenshipType)
	require.Equal(t, "MARRIED", p.MaritalStatus)
	require.Equal(t, 0, p.NumberOfDependents)
	require.Equal(t, "EMPLOYED", p.EmploymentStatus)
	require.Equal(t, "Technology", p.Occupation)
	require.Equal(t, "Unknown", p.EmployerName)
	require.Equal(t, "Senior AI & ML Engineer", p.JobTitle)
}
func TestGetCustomerError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/customers/me", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, httpResp, err := client.GetCustomer("me")
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestGetCustomerAccounts(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/customers/me/accounts", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, myAccountsResp)
	})

	resp, httpResp, err := client.GetCustomerAccounts("me")
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, 3, len(resp))

	roth := resp[0]

	require.Equal(t, "5YZ55555", roth.AccountNumber)
	require.Equal(t, "A1d557b2a-e5f1-483a-9798-13923403f442", roth.ExternalID)
	require.Equal(t, "2022-10-27T20:49:52.79Z", roth.OpenedAt.Format(time.RFC3339Nano))
	require.Equal(t, "Roth IRA", roth.Nickname)
	require.Equal(t, "Roth IRA", roth.AccountTypeName)
	require.False(t, roth.DayTraderStatus)
	require.False(t, roth.IsClosed)
	require.False(t, roth.IsFirmError)
	require.False(t, roth.IsFirmProprietary)
	require.False(t, roth.IsFuturesApproved)
	require.False(t, roth.IsTestDrive)
	require.Equal(t, "Cash", roth.MarginOrCash)
	require.False(t, roth.IsForeign)
	require.Equal(t, "2022-11-04", roth.FundingDate)
	require.Equal(t, "GROWTH", roth.InvestmentObjective)
	require.Equal(t, "Defined Risk Spreads", roth.SuitableOptionsLevel)
	require.Equal(t, "2022-10-27T20:49:52.793Z", roth.CreatedAt.Format(time.RFC3339Nano))
}

func TestGetCustomerAccountsError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/customers/me/accounts", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, httpResp, err := client.GetCustomerAccounts("me")
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestGetCustomerAccount(t *testing.T) {
	setup()
	defer teardown()

	customerID := "me"
	accountNumber := "5WV48989"

	mux.HandleFunc(fmt.Sprintf("/customers/%s/accounts/%s", customerID, accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, getCustomerAccountResp)
	})

	resp, httpResp, err := client.GetCustomerAccount(customerID, accountNumber)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, "5WV48989", resp.AccountNumber)
	require.Equal(t, "2023-06-13T23:00:29.903Z", resp.OpenedAt.Format(time.RFC3339Nano))
	require.Equal(t, "Individual", resp.Nickname)
	require.Equal(t, "Individual", resp.AccountTypeName)
	require.False(t, resp.DayTraderStatus)
	require.False(t, resp.IsClosed)
	require.False(t, resp.IsFirmError)
	require.False(t, resp.IsFirmProprietary)
	require.False(t, resp.IsFuturesApproved)
	require.False(t, resp.IsTestDrive)
	require.Equal(t, "Margin", resp.MarginOrCash)
	require.False(t, resp.IsForeign)
	require.Equal(t, "SPECULATION", resp.InvestmentObjective)
	require.Equal(t, "No Restrictions", resp.SuitableOptionsLevel)
	require.Equal(t, "2023-06-13T23:00:29.903Z", resp.CreatedAt.Format(time.RFC3339Nano))
}

func TestGetCustomerAccountError(t *testing.T) {
	setup()
	defer teardown()

	customerID := "me"
	accountNumber := "5WV48989"

	mux.HandleFunc(fmt.Sprintf("/customers/%s/accounts/%s", customerID, accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, httpResp, err := client.GetCustomerAccount(customerID, accountNumber)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestGetMyAccount(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5WV48989"

	mux.HandleFunc(fmt.Sprintf("/customers/me/accounts/%s", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, getCustomerAccountResp)
	})

	resp, httpResp, err := client.GetMyAccount(accountNumber)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, "5WV48989", resp.AccountNumber)
	require.Equal(t, "2023-06-13T23:00:29.903Z", resp.OpenedAt.Format(time.RFC3339Nano))
	require.Equal(t, "Individual", resp.Nickname)
	require.Equal(t, "Individual", resp.AccountTypeName)
	require.False(t, resp.DayTraderStatus)
	require.False(t, resp.IsClosed)
	require.False(t, resp.IsFirmError)
	require.False(t, resp.IsFirmProprietary)
	require.False(t, resp.IsFuturesApproved)
	require.False(t, resp.IsTestDrive)
	require.Equal(t, "Margin", resp.MarginOrCash)
	require.False(t, resp.IsForeign)
	require.Equal(t, "SPECULATION", resp.InvestmentObjective)
	require.Equal(t, "No Restrictions", resp.SuitableOptionsLevel)
	require.Equal(t, "2023-06-13T23:00:29.903Z", resp.CreatedAt.Format(time.RFC3339Nano))
}

func TestGetMyAccountError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5WV48989"

	mux.HandleFunc(fmt.Sprintf("/customers/me/accounts/%s", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, httpResp, err := client.GetMyAccount(accountNumber)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestGetQuoteStreamerTokens(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/quote-streamer-tokens", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, quoteStreamerTokensResp)
	})

	resp, httpResp, err := client.GetQuoteStreamerTokens()
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, "example-token-here", resp.Token)
	require.Equal(t, "tasty-live.dxfeed.com:7301", resp.StreamerURL)
	require.Equal(t, "https://tasty-live-web.dxfeed.com/live", resp.WebsocketURL)
	require.Equal(t, "live", resp.Level)
}

func TestGetQuoteStreamerTokensError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/quote-streamer-tokens", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, httpResp, err := client.GetQuoteStreamerTokens()
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

const getCustomerResp = `{
  "data": {
    "id": "me",
    "first-name": "Austin",
    "last-name": "Spencer",
    "first-surname": "Spencer",
    "address": {
      "street-one": "1000 West Fulton St",
      "city": "Chicago",
      "state-region": "IL",
      "postal-code": "60607",
      "country": "US",
      "is-foreign": false,
      "is-domestic": true
    },
    "mailing-address": {
      "street-one": "1000 West Fulton St",
      "city": "Chicago",
      "state-region": "IL",
      "postal-code": "60607",
      "country": "US",
      "is-foreign": false,
      "is-domestic": true
    },
    "customer-suitability": {
      "id": 115555,
      "marital-status": "MARRIED",
      "number-of-dependents": 0,
      "employment-status": "EMPLOYED",
      "occupation": "Technology",
      "employer-name": "Unknown",
      "job-title": "Senior AI & ML Engineer",
      "annual-net-income": 10,
      "net-worth": 0,
      "liquid-net-worth": 0,
      "stock-trading-experience": "EXTENSIVE",
      "covered-options-trading-experience": "EXTENSIVE",
      "uncovered-options-trading-experience": "EXTENSIVE",
      "futures-trading-experience": "LIMITED"
    },
    "usa-citizenship-type": "Citizen",
    "is-foreign": false,
    "mobile-phone-number": "+15555555555",
    "work-phone-number": "",
    "home-phone-number": "",
    "email": "me@austinbspencer.com",
    "tax-number-type": "SSN",
    "tax-number": "*****5555",
    "birth-date": "1900-01-01",
    "external-id": "C0001234567",
    "citizenship-country": "USA",
    "subject-to-tax-withholding": false,
    "agreed-to-margining": true,
    "agreed-to-terms": true,
    "has-industry-affiliation": false,
    "has-political-affiliation": false,
    "has-listed-affiliation": false,
    "is-investment-adviser": false,
    "has-institutional-assets": false,
    "is-professional": false,
    "has-delayed-quotes": false,
    "has-pending-or-approved-application": true,
    "permitted-account-types": [
      {
        "name": "Individual",
        "description": "",
        "is-tax-advantaged": false,
        "has-multiple-owners": false,
        "is-publicly-available": true,
        "margin-types": [
          { "name": "Cash", "is-margin": false },
          { "name": "Cash Secured Margin", "is-margin": true },
          { "name": "IRA Margin", "is-margin": true },
          { "name": "Reg T", "is-margin": true },
          { "name": "Portfolio Margin", "is-margin": true }
        ]
      },
      {
        "name": "Investment Club",
        "description": "",
        "is-tax-advantaged": false,
        "has-multiple-owners": false,
        "is-publicly-available": true,
        "margin-types": [
          { "name": "Cash", "is-margin": false },
          { "name": "Cash Secured Margin", "is-margin": true },
          { "name": "IRA Margin", "is-margin": true },
          { "name": "Reg T", "is-margin": true },
          { "name": "Portfolio Margin", "is-margin": true }
        ]
      },
      {
        "name": "Traditional IRA",
        "description": "",
        "is-tax-advantaged": true,
        "has-multiple-owners": false,
        "is-publicly-available": true,
        "margin-types": [{ "name": "IRA Margin", "is-margin": true }]
      },
      {
        "name": "Roth IRA",
        "description": "",
        "is-tax-advantaged": true,
        "has-multiple-owners": false,
        "is-publicly-available": true,
        "margin-types": [{ "name": "IRA Margin", "is-margin": true }]
      },
      {
        "name": "Beneficiary Roth IRA",
        "description": "",
        "is-tax-advantaged": true,
        "has-multiple-owners": false,
        "is-publicly-available": true,
        "margin-types": [{ "name": "IRA Margin", "is-margin": true }]
      },
      {
        "name": "Beneficiary Traditional IRA",
        "description": "",
        "is-tax-advantaged": true,
        "has-multiple-owners": false,
        "is-publicly-available": true,
        "margin-types": [{ "name": "IRA Margin", "is-margin": true }]
      }
    ],
    "identifiable-type": "Person",
    "person": {
      "external-id": "P0002345678",
      "first-name": "Austin",
      "last-name": "Spencer",
      "birth-date": "1900-01-01",
      "citizenship-country": "USA",
      "usa-citizenship-type": "Citizen",
      "marital-status": "MARRIED",
      "number-of-dependents": 0,
      "employment-status": "EMPLOYED",
      "occupation": "Technology",
      "employer-name": "Unknown",
      "job-title": "Senior AI & ML Engineer"
    }
  },
  "context": "/customers/me"
}`

const quoteStreamerTokensResp = `{
  "data": {
    "token": "example-token-here",
    "streamer-url": "tasty-live.dxfeed.com:7301",
    "websocket-url": "https://tasty-live-web.dxfeed.com/live",
    "level": "live"
  },
  "context": "/quote-streamer-tokens"
}`

const getCustomerAccountResp = `{
  "data": {
    "account-number": "5WV48989",
    "opened-at": "2023-06-13T23:00:29.903+00:00",
    "nickname": "Individual",
    "account-type-name": "Individual",
    "day-trader-status": false,
    "is-closed": false,
    "is-firm-error": false,
    "is-firm-proprietary": false,
    "is-futures-approved": false,
    "is-test-drive": false,
    "margin-or-cash": "Margin",
    "is-foreign": false,
    "investment-objective": "SPECULATION",
    "suitable-options-level": "No Restrictions",
    "created-at": "2023-06-13T23:00:29.903+00:00"
  },
  "context": "/customers/me/accounts/5WV48989"
}`
