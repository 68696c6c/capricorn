package src

import (
	"fmt"
	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/utils"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const logPrefix = "CreateDomains"

func CreateDomains(spec *models.Project, logger *logrus.Logger) error {
	err := utils.CreateDir(spec.Paths.Domains)
	if err != nil {
		return errors.Wrapf(err, "failed to create domains directory '%s'", spec.Paths.Domains)
	}

	for _, d := range spec.Domains {
		logger.Infof("%s | creating domain %s", logPrefix, d.Name)

		domainPath := fmt.Sprintf("%s/%s", spec.Paths.Domains, d.Name)
		err := utils.CreateDir(domainPath)
		if err != nil {
			return errors.Wrapf(err, "failed to create domain directory '%s'", spec.Paths.Domains)
		}

		err = createDomainModel(domainPath, d.Model, logger)
		if err != nil {
			return errors.Wrap(err, "failed to generate domain")
		}

		err = createDomainRepo(domainPath, d.Repo, logger)
		if err != nil {
			return errors.Wrap(err, "failed to generate domain")
		}

		err = createDomainService(domainPath, d.Service, logger)
		if err != nil {
			return errors.Wrap(err, "failed to generate domain")
		}

		err = createDomainController(domainPath, &d.Controller, logger)
		if err != nil {
			return errors.Wrap(err, "failed to generate domain")
		}
		spec.Controllers = append(spec.Controllers, d.Controller)

		logger.Infof("%s | done generating domain %s", logPrefix, d.Name)
	}

	return nil
}

func createDomainModel(basePath string, m models.Model, logger *logrus.Logger) error {
	logger.Infof("%s | creating model %s", logPrefix, m.Filename)

	err := utils.GenerateFile(basePath, m.Filename, modelTemplate, m)
	if err != nil {
		return errors.Wrap(err, "failed to generate model")
	}

	return nil
}

func createDomainRepo(basePath string, r models.Repo, logger *logrus.Logger) error {
	logger.Infof("%s | creating repo %s", logPrefix, r.Filename)

	for _, m := range r.Methods {
		logger.Infof("%s | creating repo method %s", logPrefix, m.Name)

		switch m.Name {
		case "Save":
			mt, err := utils.ParseTemplateToString("repo_save", repoSaveTemplate, m)
			if err != nil {
				return errors.Wrap(err, "failed to generate repo method 'Save'")
			}
			r.MethodTemplates = append(r.MethodTemplates, mt)
			r.InterfaceTemplates = append(r.InterfaceTemplates, m.Signature)

		case "GetByID":
			mt, err := utils.ParseTemplateToString("repo_get", repoGetByIDTemplate, m)
			if err != nil {
				return errors.Wrap(err, "failed to generate repo method 'GetByID'")
			}
			r.MethodTemplates = append(r.MethodTemplates, mt)
			r.InterfaceTemplates = append(r.InterfaceTemplates, m.Signature)

		case "List":
			mt, err := utils.ParseTemplateToString("repo_list", repoListTemplate, m)
			if err != nil {
				return errors.Wrap(err, "failed to generate repo method 'List'")
			}
			r.MethodTemplates = append(r.MethodTemplates, mt)
			r.InterfaceTemplates = append(r.InterfaceTemplates, m.Signature)

		case "SetQueryTotal":
			mt, err := utils.ParseTemplateToString("query_total", repoSetQueryTotalTemplate, m)
			if err != nil {
				return errors.Wrap(err, "failed to generate repo method 'SetQueryTotal'")
			}
			r.MethodTemplates = append(r.MethodTemplates, mt)
			r.InterfaceTemplates = append(r.InterfaceTemplates, m.Signature)

		case "Delete":
			mt, err := utils.ParseTemplateToString("repo_delete", repoDeleteTemplate, m)
			if err != nil {
				return errors.Wrap(err, "failed to generate repo method 'Delete'")
			}
			r.MethodTemplates = append(r.MethodTemplates, mt)
			r.InterfaceTemplates = append(r.InterfaceTemplates, m.Signature)
		}
	}

	err := utils.GenerateFile(basePath, r.Filename, repoTemplate, r)
	if err != nil {
		return errors.Wrap(err, "failed to generate repo")
	}

	return nil
}

func createDomainService(basePath string, s models.Service, logger *logrus.Logger) error {
	if len(s.Methods) == 0 {
		logger.Infof("%s | no service to generate", logPrefix, s.Filename)
		return nil
	}
	logger.Infof("%s | creating service %s", logPrefix, s.Filename)

	for _, m := range s.Methods {
		logger.Infof("%s | creating service method %s", logPrefix, m.Name)

		mt, err := utils.ParseTemplateToString("service_method", serviceMethodTemplate, m)
		if err != nil {
			return errors.Wrapf(err, "failed to generate service method '%s'", m.Name)
		}
		s.MethodTemplates = append(s.MethodTemplates, mt)
		s.InterfaceTemplates = append(s.InterfaceTemplates, m.Signature)
	}

	err := utils.GenerateFile(basePath, s.Filename, serviceTemplate, s)
	if err != nil {
		return errors.Wrap(err, "failed to generate service")
	}

	return nil
}

func createDomainController(basePath string, c *models.Controller, logger *logrus.Logger) error {
	logger.Debugf("%s | controller: %s", logPrefix, c.Filename)

	// Create requests.
	if createRequest, ok := c.Requests["create"]; ok {
		rt, err := utils.ParseTemplateToString("create_request", createRequestTemplate, createRequest)
		if err != nil {
			return errors.Wrap(err, "failed to generate controller request 'create'")
		}
		c.RequestTemplates = append(c.RequestTemplates, rt)
	}

	if updateRequest, ok := c.Requests["update"]; ok {
		rt, err := utils.ParseTemplateToString("update_request", updateRequestTemplate, updateRequest)
		if err != nil {
			return errors.Wrap(err, "failed to generate controller request 'update'")
		}
		c.RequestTemplates = append(c.RequestTemplates, rt)
	}

	// Create responses.
	if viewResponse, ok := c.Responses["view"]; ok {
		rt, err := utils.ParseTemplateToString("view_response", viewResponseTemplate, viewResponse)
		if err != nil {
			return errors.Wrap(err, "failed to generate controller response 'view'")
		}
		c.ResponseTemplates = append(c.ResponseTemplates, rt)
	}

	if listResponse, ok := c.Responses["list"]; ok {
		rt, err := utils.ParseTemplateToString("list_response", listResponseTemplate, listResponse)
		if err != nil {
			return errors.Wrap(err, "failed to generate controller response 'list'")
		}
		c.ResponseTemplates = append(c.ResponseTemplates, rt)
	}

	// Create handlers.
	for _, h := range c.Handlers {
		logger.Infof("%s | creating controller handler %s", logPrefix, h.Name)

		switch h.Name {
		case "Create":
			ht, err := utils.ParseTemplateToString("controller_create", handlerCreateTemplate, h)
			if err != nil {
				return errors.Wrap(err, "failed to generate controller handler 'Create'")
			}
			c.HandlerTemplates = append(c.HandlerTemplates, ht)
			rt, err := utils.ParseTemplateToString("route_create", routeCreateTemplate, c)
			if err != nil {
				return errors.Wrap(err, "failed to generate controller route 'Create'")
			}
			c.RoutesTemplates = append(c.RoutesTemplates, rt)

		case "Update":
			ht, err := utils.ParseTemplateToString("controller_update", handlerUpdateTemplate, h)
			if err != nil {
				return errors.Wrap(err, "failed to generate controller handler 'Update'")
			}
			c.HandlerTemplates = append(c.HandlerTemplates, ht)
			rt, err := utils.ParseTemplateToString("route_update", routeUpdateTemplate, c)
			if err != nil {
				return errors.Wrap(err, "failed to generate controller route 'Update'")
			}
			c.RoutesTemplates = append(c.RoutesTemplates, rt)

		case "List":
			ht, err := utils.ParseTemplateToString("controller_list", handlerListTemplate, h)
			if err != nil {
				return errors.Wrap(err, "failed to generate controller handler 'List'")
			}
			c.HandlerTemplates = append(c.HandlerTemplates, ht)
			rt, err := utils.ParseTemplateToString("route_list", routeListTemplate, c)
			if err != nil {
				return errors.Wrap(err, "failed to generate controller route 'List'")
			}
			c.RoutesTemplates = append(c.RoutesTemplates, rt)

		case "View":
			ht, err := utils.ParseTemplateToString("controller_view", handlerViewTemplate, h)
			if err != nil {
				return errors.Wrap(err, "failed to generate controller handler 'View'")
			}
			c.HandlerTemplates = append(c.HandlerTemplates, ht)
			rt, err := utils.ParseTemplateToString("route_lview", routeViewTemplate, c)
			if err != nil {
				return errors.Wrap(err, "failed to generate controller route 'View'")
			}
			c.RoutesTemplates = append(c.RoutesTemplates, rt)

		case "Delete":
			ht, err := utils.ParseTemplateToString("controller_delete", handlerDeleteTemplate, h)
			if err != nil {
				return errors.Wrap(err, "failed to generate controller handler 'Delete'")
			}
			c.HandlerTemplates = append(c.HandlerTemplates, ht)
			rt, err := utils.ParseTemplateToString("route_delete", routeDeleteTemplate, c)
			if err != nil {
				return errors.Wrap(err, "failed to generate controller route 'Delete'")
			}
			c.RoutesTemplates = append(c.RoutesTemplates, rt)
		}
	}

	err := utils.GenerateFile(basePath, c.Filename, controllerTemplate, c)
	if err != nil {
		return errors.Wrap(err, "failed to generate controller")
	}

	return nil
}
