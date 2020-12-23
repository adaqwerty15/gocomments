package controllers

import (
	"github.com/revel/revel"
	"gocomments/app"
	"log"
	"net/http"
	m "gocomments/app/models"
)

type App struct {
	*revel.Controller
}

type Comments struct {
    *revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}


func (c Comments) List() revel.Result {
	status := c.Params.Query.Get("status")
	company := c.Params.Query.Get("companyid")

	var err error

	sql := `SELECT comments.id, pages.id, pages.url, users.id, 
			first_name, last_name, timestamp, text, status, important from users, comments, pages
			where pages.id=comments.page_id and users.id=comments.user_id and pages.company_id=$1
			and comments.status = $2 order by timestamp DESC;`
    rows, err := app.DB.Query(sql, company, status)

    if err != nil {
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
         User:m.User{cm.UserId, cm.UserFirstName, cm.UserFirstName}, Timestamp: cm.Timestamp, 
         Text: cm.Text, Status: cm.Status, Important: cm.Important})
    }	

    data := make(map[string]interface{})
    data["comments"] = cmts

    c.Response.SetStatus(http.StatusOK)
	c.Response.ContentType = "application/json"
	return c.RenderJSON(data)
}
