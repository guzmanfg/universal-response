package home

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
)

func Home(c *gin.Context){
	r := c.Request

	message := fmt.Sprintf("[%s]\n", r.Method)

	parameters := r.URL.Query()
	if len(parameters) > 0 {
		message += "QUERY PARAMETERS\n"
		message = formatValues(parameters, message)
	}

	// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
	if err := r.ParseForm(); err != nil {
		message += fmt.Sprintf("Error: %v", err)
	} else if len(r.PostForm) > 0 {
		message += "BODY PARAMETERS\n"
		form := r.PostForm
		message = formatValues(form, message)
	}

	c.String(200, message)
}

func formatValues(body url.Values, message string) string {
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
