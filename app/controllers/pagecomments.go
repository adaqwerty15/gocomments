package controllers

import (
	"github.com/revel/revel"
	"gocomments/app"
	"log"
	"net/http"
	m "gocomments/app/models"
)

type PageComments struct {
    *revel.Controller
}

func (c PageComments) List() revel.Result {	
	url := c.Params.Query.Get("url")
	company := c.Params.Query.Get("companyid")

	var err error

	sql := `SELECT comments.id, pages.id, pages.url, users.id, 
			first_name, last_name, timestamp, text, status, important from users, comments, pages
			where pages.id=comments.page_id and users.id=comments.user_id and pages.company_id=$1
			and pages.url = $2 and status='published' order by timestamp DESC;`

	rows, err := app.DB.Query(sql, company, url)

	if err != nil {
		log.Println(err)
		c.Response.ContentType = "application/json"
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON("{}")
    }
    defer rows.Close()

    cmts := make([]m.CommentInfo,0)

    for rows.Next() {
        cm := new(m.Comment)
        err := rows.Scan(&cm.Id, &cm.PageId, &cm.PageUrl, &cm.UserId, &cm.UserFirstName, 
        	&cm.UserLastName, &cm.Timestamp, &cm.Text, &cm.Status, &cm.Important)
        if err != nil {
            log.Fatal(err)
        }
        cmts = append(cmts, m.CommentInfo{Id: cm.Id, Page:m.Page{cm.PageId, cm.PageUrl},
         User:m.User{cm.UserId, cm.UserFirstName, cm.UserLastName}, Timestamp: cm.Timestamp, 
         Text: cm.Text, Status: cm.Status, Important: cm.Important})
    }	

    data := make(map[string]interface{})
    data["comments"] = cmts

    c.Response.SetStatus(http.StatusOK)
	c.Response.ContentType = "application/json"
	return c.RenderJSON(data)
}	


func (c PageComments) Add() revel.Result {	
	var jsonData map[string]interface{}
    c.Params.BindJSON(&jsonData)

    company := jsonData["company_id"]
    url := jsonData["url"]

    var err error

    var is_moderated bool
    sql := `SELECT is_moderated from companies where id=$1;`
    err = app.DB.QueryRow(sql, company).Scan(&is_moderated)

    if err != nil {
		log.Fatal(err)
		c.Response.ContentType = "application/json"
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON("{}")
	}

	status := "unmoderated"

	if !is_moderated {
		status = "published"
	}

	var id int

	sql = `SELECT id from pages where url=$1 and company_id=$2;`
	err = app.DB.QueryRow(sql, url, company).Scan(&id)

    if err != nil {
		log.Println(err)
		stmt, err2 := app.DB.Prepare("INSERT INTO pages(company_id, url) VALUES($1, $2) RETURNING id;")
		if err2 != nil {
			log.Println(err2)
			c.Response.ContentType = "application/json"
			c.Response.SetStatus(http.StatusInternalServerError)
			return c.RenderJSON("{}")
		}
		err2 = stmt.QueryRow(company, url).Scan(&id)

		if err2 != nil {
			log.Println(err2)
			c.Response.ContentType = "application/json"
			c.Response.SetStatus(http.StatusInternalServerError)
			return c.RenderJSON("{}")
		}

		
	}

	sql = `INSERT INTO comments(
	page_id, user_id, text, "timestamp", status, important)
	VALUES ($1, $2, $3, now(), $4, $5) RETURNING id;`

	var commentId int

	err = app.DB.QueryRow(sql, id, jsonData["user_id"], jsonData["text"], status, false).Scan(&commentId)


	if err != nil {
		log.Println(err)
		c.Response.ContentType = "application/json"
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON("{}")
	}


	log.Println(commentId)

	sql = `SELECT comments.id, pages.id, pages.url, users.id, 
			first_name, last_name, timestamp, text, status, important from users, comments, pages
			where pages.id=comments.page_id and users.id=comments.user_id and comments.id=$1;`

	cm := new(m.Comment)
   	
	err = app.DB.QueryRow(sql, commentId).Scan(&cm.Id, 
			&cm.PageId, &cm.PageUrl, &cm.UserId, &cm.UserFirstName, 
        	&cm.UserLastName, &cm.Timestamp, &cm.Text, &cm.Status, &cm.Important)

    if err != nil {
		c.Response.ContentType = "application/json"
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON("{}")
	}    		

    cmi := m.CommentInfo{Id: cm.Id, Page:m.Page{cm.PageId, cm.PageUrl},
         User:m.User{cm.UserId, cm.UserFirstName, cm.UserLastName}, Timestamp: cm.Timestamp, 
         Text: cm.Text, Status: cm.Status, Important: cm.Important}    		

    
	data := make(map[string]interface{})
    data["comment"] = cmi

    c.Response.SetStatus(http.StatusOK)
	c.Response.ContentType = "application/json"
	return c.RenderJSON(data)
}