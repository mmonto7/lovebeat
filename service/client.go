package service

import (
	"github.com/boivie/lovebeat-go/backend"
)

type ServiceIf interface {
	Beat(name string)
	DeleteService(name string)

	SetWarningTimeout(name string, timeout int)
	SetErrorTimeout(name string, timeout int)

	CreateOrUpdateView(name string, regexp string, alertMail string)
	DeleteView(name string)
	GetServices(view string) []backend.StoredService
	GetService(name string) backend.StoredService
	GetViews() []backend.StoredView
	GetView(name string) backend.StoredView
}

const (
	ACTION_SET_WARN = "set-warn"
	ACTION_SET_ERR  = "set-err"
	ACTION_BEAT     = "beat"
	ACTION_DELETE   = "delete"
)

type serviceCmd struct {
	Action  string
	Service string
	Value   int
}

type upsertViewCmd struct {
	View      string
	Regexp    string
	AlertMail string
}

type getServicesCmd struct {
	View  string
	Reply chan []backend.StoredService
}

type getServiceCmd struct {
	Name  string
	Reply chan backend.StoredService
}

type getViewsCmd struct {
	Reply chan []backend.StoredView
}

type getViewCmd struct {
	Name  string
	Reply chan backend.StoredView
}

type client struct {
	svcs *Services
}

func (c *client) Beat(name string) {
	c.svcs.serviceCmdChan <- &serviceCmd{
		Action:  ACTION_BEAT,
		Service: name,
		Value:   1,
	}
}

func (c *client) DeleteService(name string) {
	c.svcs.serviceCmdChan <- &serviceCmd{
		Action:  ACTION_DELETE,
		Service: name,
	}
}

func (c *client) DeleteView(name string) {
	c.svcs.deleteViewCmdChan <- name
}

func (c *client) SetWarningTimeout(name string, timeout int) {
	c.svcs.serviceCmdChan <- &serviceCmd{
		Action:  ACTION_SET_WARN,
		Service: name,
		Value:   timeout,
	}
}
func (c *client) SetErrorTimeout(name string, timeout int) {
	c.svcs.serviceCmdChan <- &serviceCmd{
		Action:  ACTION_SET_ERR,
		Service: name,
		Value:   timeout,
	}
}

func (c *client) CreateOrUpdateView(name string, regexp string, alertMail string) {
	c.svcs.upsertViewCmdChan <- &upsertViewCmd{
		View:      name,
		Regexp:    regexp,
		AlertMail: alertMail,
	}
}

func (c *client) GetServices(view string) []backend.StoredService {
	myc := make(chan []backend.StoredService)
	c.svcs.getServicesChan <- &getServicesCmd{View: view, Reply: myc}
	ret := <-myc
	return ret
}

func (c *client) GetService(name string) backend.StoredService {
	myc := make(chan backend.StoredService)
	c.svcs.getServiceChan <- &getServiceCmd{Name: name, Reply: myc}
	ret := <-myc
	return ret
}

func (c *client) GetViews() []backend.StoredView {
	myc := make(chan []backend.StoredView)
	c.svcs.getViewsChan <- &getViewsCmd{Reply: myc}
	ret := <-myc
	return ret
}

func (c *client) GetView(name string) backend.StoredView {
	myc := make(chan backend.StoredView)
	c.svcs.getViewChan <- &getViewCmd{Name: name, Reply: myc}
	ret := <-myc
	return ret
}

func (svcs *Services) GetClient() ServiceIf {
	return &client{svcs: svcs}
}
