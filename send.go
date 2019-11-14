package main

import (
	"bytes"
	html "html/template"
	"log"

	gomail "gopkg.in/gomail.v2"
)

type varsTemp struct {
	selector string
	valor    string
}

type callBody struct {
	TemplateName      string            `binding:"required" json:"template_name"`
	Direccion         []string          `binding:"required" json:"to_address"`
	IDUsuario         []int             `json:"id_usuario"`
	Subject           string            `binding:"required" json:"subject"`
	VariablesTemplate map[string]string `json:"variables_template"`
}

func parseTemplate(templateBody string, templateData interface{}) (string, error) {
	parseo, err := html.New("temp").Parse(templateBody)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = parseo.Execute(buf, templateData); err != nil {
		return "", err
	}
	result := buf.String()
	return result, nil
}

func (b callBody) sendSpecific() error {
	tempname := b.TemplateName
	_, temp, err := getTemplateFirebase(tempname)

	m := gomail.NewMessage()
	m.SetHeader("From", "microservicios.mailing@gmail.com")
	m.SetHeader("To", b.Direccion...)
	m.SetHeader("Subject", b.Subject)
	// recorro todas las variablesTemplate que recibo de la request para setearlas en el template
	res, err := parseTemplate(temp.TemplateBody, b.VariablesTemplate)
	if err != nil {
		log.Print(err)
		return err
	}
	log.Print(res)
	m.SetBody("text/html", res)

	d := gomail.NewDialer("smtp.gmail.com", 587, "microservicios.mailing@gmail.com", "lkohecuypouihukj")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	return nil
}