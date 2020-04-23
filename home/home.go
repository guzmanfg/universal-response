package home

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func Home(c *gin.Context) {

	_ = c.Request.ParseForm()

	if len(c.Request.URL.Query()) + len(c.Request.PostForm) > 0 {
		ProcessForm(c)
	} else {
		Index(c)
	}

}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl.html", nil)
}

func ProcessForm(c *gin.Context) {
	r := c.Request
	response := fmt.Sprintf("[%s]\n", r.Method)

	parameters := r.URL.Query()
	if len(parameters) > 0 {
		response += "QUERY PARAMETERS\n"
		response = FormatValues(parameters, response)
	}

	if len(r.PostForm) > 0 {
		response += "BODY PARAMETERS\n"
		form := r.PostForm
		response = FormatValues(form, response)
	}

	c.String(200, response)
}

func FormatValues(body url.Values, message string) string {
	for key, values := range body {
		if len(values) > 1 {
			for index, value := range values {
				message += fmt.Sprintf("%s[%d] = %s\n", key, index, value)
			}
		} else {
			for _, value := range values {
				message += fmt.Sprintf("%s = %s\n", key, value)
			}
		}
	}
	return message
}
