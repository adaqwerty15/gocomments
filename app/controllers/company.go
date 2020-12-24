package controllers

import (
	"github.com/revel/revel"
	"gocomments/app"
	"log"
	"net/http"
	m "gocomments/app/models"
)

type Company struct {
    *revel.Controller
}

func (c Company) Info() revel.Result {
	id := c.Params.Query.Get("id")

	var err error

	sql := `SELECT companies.id, user_main_id, first_name, last_name, name, website, is_moderated, is_authed
	FROM companies, users where users.id=companies.user_main_id and companies.id=$1;`

	cm := new(m.Company)
   	
	err = app.DB.QueryRow(sql, id).Scan(&cm.Id, 
			&cm.UserId,  &cm.FirstName, 
        	&cm.LastName, &cm.Name, &cm.Website, &cm.IsModerated, &cm.IsAuthed)

    if err != nil {
		c.Response.ContentType = "application/json"
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON("{}")
	}    		

    cmi := m.CompanyInfo{Id: cm.Id,
         User:m.User{cm.UserId, cm.FirstName, cm.FirstName}, Name: cm.Name, 
         Website: cm.Website, IsModerated: cm.IsModerated, IsAuthed: cm.IsAuthed}    		

    
	data := make(map[string]interface{})
    data["company"] = cmi

    c.Response.SetStatus(http.StatusOK)
	c.Response.ContentType = "application/json"
	return c.RenderJSON(data)
}


func (c Company) Stats() revel.Result {
	period := c.Params.Query.Get("period")
	company := c.Params.Query.Get("id")

	var err error
	var sql string

	if period=="year" {
		sql = `select count(c.id), c.status from comments as c, pages as p 
		where p.id=c.page_id and p.company_id=$1
		and date_part('year', current_timestamp)=date_part('year', c.timestamp)
		GROUP BY status`
	} else if period=="month" {
		sql = `select count(c.id), c.status from comments as c, pages as p 
		where p.id=c.page_id and p.company_id=$1
		and date_part('year', current_timestamp)=date_part('year', c.timestamp)
		and date_part('month', current_timestamp)=date_part('month', c.timestamp)
		GROUP BY status`
	} else if period=="day" {
		sql = `select count(c.id), c.status from comments as c, pages as p 
		where p.id=c.page_id and p.company_id=$1
		and date_part('year', current_timestamp)=date_part('year', c.timestamp)
		and date_part('month', current_timestamp)=date_part('month', c.timestamp)
		and date_part('day', current_timestamp)=date_part('day', c.timestamp)
		GROUP BY status`
	}

    rows, err := app.DB.Query(sql, company)

    if err != nil {
		c.Response.ContentType = "application/json"
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON("{}")
    }
    defer rows.Close()

    cmts := make([]*m.Stat,0)

    for rows.Next() {
        cm := new(m.Stat)
        err := rows.Scan(&cm.Value, &cm.Name)
        if err != nil {
            log.Fatal(err)
        }
        cmts = append(cmts, cm)
    }	

    data := make(map[string]interface{})
    data["stats"] = cmts

    c.Response.SetStatus(http.StatusOK)
	c.Response.ContentType = "application/json"
	return c.RenderJSON(data)
}


func (c Company) ChangeModeration() revel.Result {

	is_mod := c.Params.Query.Get("is_mod")
	id := c.Params.Query.Get("id")

	var is_mod_b bool

	var err error

	if (is_mod=="true") {
		is_mod_b = true
	} else {
		is_mod_b = false
	}

	stmt, err := app.DB.Prepare("UPDATE companies SET is_moderated=$2 WHERE id=$1;")
	if err != nil {
		log.Println(err)
		c.Response.ContentType = "application/json"
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON("{}")
	}
	_, err = stmt.Exec(id, is_mod_b)

	if err != nil {
		log.Println(err)
		c.Response.ContentType = "application/json"
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON("{}")
	}

	sql := `SELECT companies.id, user_main_id, first_name, last_name, name, website, is_moderated, is_authed
	FROM companies, users where users.id=companies.user_main_id and companies.id=$1;`

	cm := new(m.Company)
   	
	err = app.DB.QueryRow(sql, id).Scan(&cm.Id, 
			&cm.UserId,  &cm.FirstName, 
        	&cm.LastName, &cm.Name, &cm.Website, &cm.IsModerated, &cm.IsAuthed)

    if err != nil {
		c.Response.ContentType = "application/json"
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON("{}")
	}    		

    cmi := m.CompanyInfo{Id: cm.Id,
         User:m.User{cm.UserId, cm.FirstName, cm.FirstName}, Name: cm.Name, 
         Website: cm.Website, IsModerated: cm.IsModerated, IsAuthed: cm.IsAuthed}    		

    
	data := make(map[string]interface{})
    data["company"] = cmi

    c.Response.SetStatus(http.StatusOK)
	c.Response.ContentType = "application/json"
	return c.RenderJSON(data)
}



func (c Company) ChangeAuth() revel.Result {

	is_mod := c.Params.Query.Get("is_auth")
	id := c.Params.Query.Get("id")

	var is_mod_b bool

	var err error

	if (is_mod=="true") {
		is_mod_b = true
	} else {
		is_mod_b = false
	}

	stmt, err := app.DB.Prepare("UPDATE companies SET is_authed=$2 WHERE id=$1;")
	if err != nil {
		log.Println(err)
		c.Response.ContentType = "application/json"
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON("{}")
	}
	_, err = stmt.Exec(id, is_mod_b)

	if err != nil {
		log.Println(err)
		c.Response.ContentType = "application/json"
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON("{}")
	}

	sql := `SELECT companies.id, user_main_id, first_name, last_name, name, website, is_moderated, is_authed
	FROM companies, users where users.id=companies.user_main_id and companies.id=$1;`

	cm := new(m.Company)
   	
	err = app.DB.QueryRow(sql, id).Scan(&cm.Id, 
			&cm.UserId,  &cm.FirstName, 
        	&cm.LastName, &cm.Name, &cm.Website, &cm.IsModerated, &cm.IsAuthed)

    if err != nil {
		c.Response.ContentType = "application/json"
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON("{}")
	}    		

    cmi := m.CompanyInfo{Id: cm.Id,
         User:m.User{cm.UserId, cm.FirstName, cm.FirstName}, Name: cm.Name, 
         Website: cm.Website, IsModerated: cm.IsModerated, IsAuthed: cm.IsAuthed}    		

    
	data := make(map[string]interface{})
    data["company"] = cmi

    c.Response.SetStatus(http.StatusOK)
	c.Response.ContentType = "application/json"
	return c.RenderJSON(data)
}
