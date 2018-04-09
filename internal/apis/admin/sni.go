package admin

import (
	"encoding/json"
	"net/url"

	adminv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/admin/v1"
)

type SNIGetter interface {
	SNI() SNIInterface
}

type SNIInterface interface {
	List(params url.Values) (*adminv1.SNIList, error)
	Get(name string) (*adminv1.SNI, *APIResponse)
	Create(sni *adminv1.SNI) (*adminv1.SNI, *APIResponse)
	Delete(name string) error
}

type sniAPI struct {
	client APIInterface
}

func (a *sniAPI) Create(target *adminv1.SNI) (*adminv1.SNI, *APIResponse) {
	out := &adminv1.SNI{}
	err := a.client.Create(target, out)
	return out, err
}

func (a *sniAPI) Get(name string) (*adminv1.SNI, *APIResponse) {
	out := &adminv1.SNI{}
	err := a.client.Get(name, out)
	return out, err
}

func (a *sniAPI) List(params url.Values) (*adminv1.SNIList, error) {
	targets := &adminv1.SNIList{}
	request := a.client.RestClient().Get().Resource("snis")
	for k, vals := range params {
		for _, v := range vals {
			request.Param(k, v)
		}
	}
	data, err := request.DoRaw()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, targets); err != nil {
		return nil, err
	}

	if len(targets.NextPage) > 0 {
		params.Add("offset", targets.Offset)
		result, err := a.List(params)
		if err != nil {
			return nil, err
		}
		targets.Items = append(targets.Items, result.Items...)
	}

	return targets, err
}

func (a *sniAPI) Delete(name string) error {
	return a.client.Delete(name)
}
