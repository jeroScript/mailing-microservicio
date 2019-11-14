package main

import (
	"fmt"
	"log"

	"google.golang.org/api/iterator"
)

type template struct {
	TemplateName string `json:"template_name" firestore:"-"`
	TemplateBody string `binding:"required" json:"template_body" firestore:"template_body,omitempty"`
}

type templateParsed struct {
	TemplateBody string
}

func (t template) postTemplateFirebase() error {

	client, ctx, err := getFireStore()
	if err != nil {
		panic(err)
	}
	result, err := client.Collection("template").Doc(t.TemplateName).Set(ctx, t)
	if err != nil {
		return err
	}
	fmt.Print(result)
	return nil
}

func getTemplateFirebase(tempName string) ([]template, template, error) {

	client, ctx, err := getFireStore()
	if err != nil {
		panic(err)
	}

	//obtener los template por nombre
	if tempName == "" {
		sliceTemp := make([]template, 0)
		array := client.Collection("template").Documents(ctx)
		i := 0
		for {
			var tempElem template
			doc, err := array.Next()

			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, template{"", ""}, err
			}
			tempElem.TemplateName = doc.Ref.ID
			doc.DataTo(&tempElem)
			sliceTemp = append(sliceTemp, tempElem)
			i++
		}
		return sliceTemp, template{"", ""}, nil
	} else {
		var tempzero template
		doc, err := client.Collection("template").Doc(tempName).Get(ctx)
		if err != nil {
			return nil, template{"", ""}, err
		}
		tempzero.TemplateName = doc.Ref.ID
		err = doc.DataTo(&tempzero)
		if err != nil {
			log.Print("error datato ", err)
		}
		return nil, tempzero, nil
	}
}
