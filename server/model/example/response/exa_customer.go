package response

import "git.echol.cn/loser/menu-server/model/example"

type ExaCustomerResponse struct {
	Customer example.ExaCustomer `json:"customer"`
}
