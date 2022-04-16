package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type context interface {
	*resty.Client
}

type authenticate func(c *resty.Client) bool
type available func() bool
type prebook func() bool
type book func() bool
type confirm func() bool
type cancel func() bool
type refund func() bool

type Provider struct {
	name    string
	actions ProviderActions
}

func newProvider(name string, actions ProviderActions) *Provider {
	client := resty.New()
	checkActionsInstance(actions)
	provider := &Provider{
		name:    name,
		actions: actions,
	}
	provider.actions.authenticate(client)
	return provider
}

type ProviderActions struct {
	authenticate authenticate
	available    available
	prebook      prebook
	book         book
	confirm      confirm
	cancel       cancel
	refund       refund
}

func checkActionsInstance(a ProviderActions) {
	if a.authenticate == nil {
		panic("authenticate method in passed actions is not defined")
	}
	if a.available == nil {
		panic("available method in passed actions is not defined")
	}
	if a.prebook == nil {
		panic("prebook method in passed actions is not defined")
	}
	if a.book == nil {
		panic("book method in passed actions is not defined")
	}
	if a.confirm == nil {
		panic("confirm method in passed actions is not defined")
	}
	if a.cancel == nil {
		panic("cancel method in passed actions is not defined")
	}
	if a.refund == nil {
		panic("refund method in passed actions is not defined")
	}
}

func newContainer() {

}

func main() {
	actions := ProviderActions{
		authenticate: func(client *resty.Client) bool {
			resp, err := client.R().Get("https://www.example.com")
			fmt.Println(resp, err)
			return true
		},
		available: func() bool {
			return true
		},
		prebook: func() bool {
			return true
		},
		book: func() bool {
			return true
		},
		confirm: func() bool {
			return true
		},
		cancel: func() bool {
			return true
		},
		refund: func() bool {
			return true
		},
	}
	newProvider("ADOTEL", actions)
}
