package tasty

import (
	"fmt"
	"net/http"

	"github.com/austinbspencer/tasty-go/models"
)

func (c *Client) GetMarginRequirements(accountNumber string) (models.MarginRequirements, *Error) {
	if c.Session.SessionToken == nil {
		return models.MarginRequirements{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/margin/accounts/%s/requirements", c.baseURL, accountNumber)

	type accountResponse struct {
		MarginRequirements models.MarginRequirements `json:"data"`
	}

	accountsRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, nil, accountsRes)
	if err != nil {
		return models.MarginRequirements{}, err
	}

	return accountsRes.MarginRequirements, nil
}

func (c *Client) MarginRequirementsDryRun(accountNumber string) {
	return
}
