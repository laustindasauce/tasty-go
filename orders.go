package tasty

import (
	"fmt"
	"net/http"
)

// Reconfirm an order
// ** This is currently untested.
// Requires the order to be an Equity Offering
// Unable to submit equity offering orders even in cert environment
// equity_offering_not_supported.
func (c *Client) ReconfirmOrder(accountNumber string, id int) (Order, error) {
	path := fmt.Sprintf("/accounts/%s/orders/%d/reconfirm", accountNumber, id)

	type ordersResponse struct {
		Order Order `json:"data"`
	}

	ordersRes := new(ordersResponse)

	err := c.request(http.MethodPost, path, nil, nil, ordersRes)
	if err != nil {
		return Order{}, err
	}

	return ordersRes.Order, nil
}

// Create an order and then runs the preflights without placing the order.
func (c *Client) SubmitOrderDryRun(accountNumber string, order NewOrder) (OrderResponse, *OrderErrorResponse, error) {
	path := fmt.Sprintf("/accounts/%s/orders/dry-run", accountNumber)

	type ordersResponse struct {
		OrderResponse OrderResponse       `json:"data"`
		OrderError    *OrderErrorResponse `json:"error"`
	}

	ordersRes := new(ordersResponse)

	err := c.request(http.MethodPost, path, nil, order, ordersRes)
	if err != nil {
		return OrderResponse{}, nil, err
	}

	return ordersRes.OrderResponse, ordersRes.OrderError, nil
}

// Create an order for the client.
func (c *Client) SubmitOrder(accountNumber string, order NewOrder) (OrderResponse, *OrderErrorResponse, error) {
	path := fmt.Sprintf("/accounts/%s/orders", accountNumber)

	type ordersResponse struct {
		OrderResponse OrderResponse       `json:"data"`
		OrderError    *OrderErrorResponse `json:"error"`
	}

	ordersRes := new(ordersResponse)

	err := c.request(http.MethodPost, path, nil, order, ordersRes)
	if err != nil {
		return OrderResponse{}, nil, err
	}

	return ordersRes.OrderResponse, ordersRes.OrderError, nil
}

// Returns a list of live orders for the resource.
func (c *Client) GetAccountLiveOrders(accountNumber string) ([]Order, error) {
	path := fmt.Sprintf("/accounts/%s/orders/live", accountNumber)

	type ordersResponse struct {
		Data struct {
			Orders []Order `json:"items"`
		} `json:"data"`
	}

	ordersRes := new(ordersResponse)

	err := c.request(http.MethodGet, path, nil, nil, ordersRes)
	if err != nil {
		return []Order{}, err
	}

	return ordersRes.Data.Orders, nil
}

// Returns a paginated list of the account's orders (as identified by the provided
// authentication token) based on sort param. If no sort is passed in, it defaults
// to descending order.
func (c *Client) GetAccountOrders(accountNumber string, query OrdersQuery) ([]Order, error) {
	path := fmt.Sprintf("/accounts/%s/orders", accountNumber)

	type ordersResponse struct {
		Data struct {
			Orders []Order `json:"items"`
		} `json:"data"`
	}

	ordersRes := new(ordersResponse)

	err := c.request(http.MethodGet, path, query, nil, ordersRes)
	if err != nil {
		return []Order{}, err
	}

	return ordersRes.Data.Orders, nil
}

// Runs through preflights for cancel-replace and edit without routing.
func (c *Client) SubmitOrderECRDryRun(accountNumber string, id int, orderECR NewOrderECR) (OrderResponse, error) {
	path := fmt.Sprintf("/accounts/%s/orders/%d/dry-run", accountNumber, id)

	type ordersResponse struct {
		OrderResponse OrderResponse `json:"data"`
	}

	ordersRes := new(ordersResponse)

	err := c.request(http.MethodPost, path, nil, orderECR, ordersRes)
	if err != nil {
		return OrderResponse{}, err
	}

	return ordersRes.OrderResponse, nil
}

// Returns a single order based on the id.
func (c *Client) GetOrder(accountNumber string, id int) (Order, error) {
	path := fmt.Sprintf("/accounts/%s/orders/%d", accountNumber, id)

	type ordersResponse struct {
		Order Order `json:"data"`
	}

	ordersRes := new(ordersResponse)

	err := c.request(http.MethodGet, path, nil, nil, ordersRes)
	if err != nil {
		return Order{}, err
	}

	return ordersRes.Order, nil
}

// Requests order cancellation.
func (c *Client) CancelOrder(accountNumber string, id int) (Order, error) {
	path := fmt.Sprintf("/accounts/%s/orders/%d", accountNumber, id)

	type ordersResponse struct {
		Order Order `json:"data"`
	}

	ordersRes := new(ordersResponse)

	err := c.request(http.MethodDelete, path, nil, nil, ordersRes)
	if err != nil {
		return Order{}, err
	}

	return ordersRes.Order, nil
}

// Replaces a live order with a new one. Subsequent fills of the original
// order will abort the replacement.
func (c *Client) ReplaceOrder(accountNumber string, id int, orderECR NewOrderECR) (Order, error) {
	path := fmt.Sprintf("/accounts/%s/orders/%d", accountNumber, id)

	type ordersResponse struct {
		Order Order `json:"data"`
	}

	ordersRes := new(ordersResponse)

	err := c.request(http.MethodPut, path, nil, orderECR, ordersRes)
	if err != nil {
		return Order{}, err
	}

	return ordersRes.Order, nil
}

// Edit price and execution properties of a live order by replacement.
// Subsequent fills of the original order will abort the replacement.
func (c *Client) PatchOrder(accountNumber string, id int, orderECR NewOrderECR) (Order, error) {
	path := fmt.Sprintf("/accounts/%s/orders/%d", accountNumber, id)

	type ordersResponse struct {
		Order Order `json:"data"`
	}

	ordersRes := new(ordersResponse)

	err := c.request(http.MethodPatch, path, nil, orderECR, ordersRes)
	if err != nil {
		return Order{}, err
	}

	return ordersRes.Order, nil
}

// Returns a list of live orders for the resource.
// Requires account numbers param to pull orders from.
func (c *Client) GetCustomerLiveOrders(customerID string, query OrdersQuery) ([]Order, error) {
	path := fmt.Sprintf("/customers/%s/orders/live", customerID)

	type ordersResponse struct {
		Data struct {
			Orders []Order `json:"items"`
		} `json:"data"`
	}

	ordersRes := new(ordersResponse)

	err := c.request(http.MethodGet, path, query, nil, ordersRes)
	if err != nil {
		return []Order{}, err
	}

	return ordersRes.Data.Orders, nil
}

// Returns a paginated list of the customer's orders (as identified by the provided
// authentication token) based on sort param. If no sort is passed in, it defaults
// to descending order. Requires account numbers param to pull orders from.
func (c *Client) GetCustomerOrders(customerID string, query OrdersQuery) ([]Order, error) {
	path := fmt.Sprintf("/customers/%s/orders", customerID)

	type ordersResponse struct {
		Data struct {
			Orders []Order `json:"items"`
		} `json:"data"`
	}

	ordersRes := new(ordersResponse)

	err := c.request(http.MethodGet, path, query, nil, ordersRes)
	if err != nil {
		return []Order{}, err
	}

	return ordersRes.Data.Orders, nil
}
