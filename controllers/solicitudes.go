package controllers

import (
	"github.com/astaxie/beego"
	solicitudesHelper "github.com/udistrital/polux_mid/helpers/solicitudes"
	"github.com/udistrital/utils_oas/errorctrl"
)

// SolicitudesController operations for Solicitudes
type SolicitudesController struct {
	beego.Controller
}

// URLMapping ...
func (c *SolicitudesController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Solicitudes
// @Param	body		body 	models.Solicitudes	true		"body for Solicitudes content"
// @Success 201 {object} models.Solicitudes
// @Failure 403 body is empty
// @router / [post]
func (c *SolicitudesController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Solicitudes by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Solicitudes
// @Failure 403 :id is empty
// @router /:id [get]
func (c *SolicitudesController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Solicitudes
// @Param	user	query	string	false	"email del usuario del que se hace la consula"
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Solicitudes
// @Failure 403
// @router / [get]
func (c *SolicitudesController) GetAll() {
	defer errorctrl.ErrorControlController(c.Controller, "SolicitudesController")

	user := c.GetString("user")

	v, err := solicitudesHelper.GetSolicitudesByUser(user)
	if err != nil {
		beego.Error(err)
		c.Abort("404")
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Solicitudes
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Solicitudes	true		"body for Solicitudes content"
// @Success 200 {object} models.Solicitudes
// @Failure 403 :id is not int
// @router /:id [put]
func (c *SolicitudesController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Solicitudes
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *SolicitudesController) Delete() {

}
