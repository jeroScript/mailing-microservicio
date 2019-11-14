package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	req := gin.Default()
	v1 := req.Group("/mailing/v1")

	templ := v1.Group("/template")
	templ.POST("/", func(c *gin.Context) { postTemplate(c) })
	templ.POST("/:template_name", func(c *gin.Context) { postTemplate(c) })
	templ.GET("/:template_name", func(c *gin.Context) { getTemplate(c) })
	templ.GET("/", func(c *gin.Context) { getTemplate(c) })

	send := v1.Group("/send")

	send.POST("/specific/", func(c *gin.Context) { specific(c) })
	send.POST("/broadcast/", func(c *gin.Context) { broadcast(c) })
	go startRabbit()
	req.Run(":22000")

}

func postTemplate(c *gin.Context) {
	tempname := c.Param("template_name")
	var bindeo template
	err := c.BindJSON(&bindeo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error bindeo": &err})
	} else {
		if tempname == "" {
			err := bindeo.postTemplateFirebase()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error firebase": &err})
			} else {
				c.JSON(http.StatusAccepted, gin.H{"Template creado!: ": &bindeo.TemplateName})
			}
		} else {
			err := bindeo.postTemplateFirebase()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error firebase": &err})
			} else {
				c.JSON(http.StatusAccepted, gin.H{"Template modificado! :": &tempname})
			}

		}
	}
	return
}

func getTemplate(c *gin.Context) {
	nombre := c.Param("template_name")
	var arrayTemp []template
	arrayTemp, _, err := getTemplateFirebase(nombre)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error getting template from firebase": err})
		//c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		if err != nil {
			log.Print(err)
		}
		c.JSON(http.StatusAccepted, gin.H{"templatesEncontrados": arrayTemp})
	}
}

func specific(c *gin.Context) {
	var body callBody
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		c.JSON(http.StatusBadRequest, gin.H{"error bindeo": err})
	} else {
		err := body.sendSpecific()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error sending": &err})
		} else {
			c.JSON(http.StatusAccepted, gin.H{"success": "Se envio correctamente el mail"})
		}
	}

}

func broadcast(c *gin.Context) {
	var body callBody
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		c.JSON(http.StatusBadRequest, gin.H{"error bindeo": err})
	} else {
		c.JSON(http.StatusAccepted, gin.H{"mensaje broadcast enviado del template ": body.TemplateName})
		// c.JSON(http.StatusAccepted, gin.H{"mensaje": err.template_name"))})
	}
}
