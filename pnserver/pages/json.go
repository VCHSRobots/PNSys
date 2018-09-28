// --------------------------------------------------------------------
// jsondata.go -- Handles Json data requests
//
// Created 2018-09-22 DLB
// --------------------------------------------------------------------

package pages

import (
	"github.com/gin-gonic/gin"
)

type TProject struct {
	Project    string
	SubSystems []string
}

func init() {
	RegisterPage("JsonData", Invoke_GET, handle_json_data)
}

func handle_json_data(c *gin.Context) {
	data := make([]*TProject, 0, 10)
	p1 := &TProject{"F18 - Fall 2018", []string{"CH -- Chasssis", "FS -- Fresby Shooter",
		"MS -- Miscellaous"}}
	p2 := &TProject{"P18 -- Practice 2018", []string{"CH", "AB", "TW", "DR", "EL"}}
	p3 := &TProject{"C18 -- Competition 2018", []string{"CH", "AM", "BC", "DE"}}
	data = append(data, p1)
	data = append(data, p2)
	data = append(data, p3)

	c.JSON(200, data)
}
