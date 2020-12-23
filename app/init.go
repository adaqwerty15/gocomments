package app

import (
	"github.com/revel/revel"
	 _"github.com/lib/pq"
    "database/sql"
    "log"
    "fmt"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string

	DB *sql.DB
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.BeforeAfterFilter,       // Call the before and after filter functions
		revel.ActionInvoker,           // Invoke the action.
	}

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	// revel.OnAppStart(ExampleStartupScript)
	revel.OnAppStart(InitDB)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

//func ExampleStartupScript() {
//	// revel.DevMod and revel.RunMode work here
//	// Use this script to check for dev mode and set dev/prod startup scripts here!
//	if revel.DevMode == true {
//		// Dev mode
//	}
//}

func InitDB() {

  var err error           

  // conn := fmt.Sprintf("user=%s password='%s' host=%s port=%d dbname=%s", "adaqwerty15", "1", "localhost", 6543, "gocomments")
  conn := fmt.Sprintf("user=%s password='%s' host=%s port=%d dbname=%s", "postgres", "1", "db", 5432, "gocomments")
  //conn := postgresql://postgres:1@postgres/gocomments
    
  DB, err = sql.Open("postgres", conn)

  if err != nil {
	log.Println("DB Error", err)
  }
  
  log.Println("DB Connected")

  sql := `CREATE TABLE if not exists users (
		    id serial NOT NULL,
		    auth_type character varying(100),
		    first_name character varying(100),
		    last_name character varying(100),
		    photo character varying(255),
		    CONSTRAINT users_pkey PRIMARY KEY (id));
		     
		 INSERT INTO users(id,
		 auth_type, first_name, last_name)
		 VALUES (0, 'in', 'Ivan', 'Ivanov')
		 ON CONFLICT DO NOTHING;

		 CREATE TABLE if not exists companies (
		    id serial NOT NULL,
		    name character varying(100),
		    website character varying(100),
		    is_moderated boolean,
		    is_authed boolean,
		    CONSTRAINT companies_pkey PRIMARY KEY (id));

		 INSERT INTO companies(
			id, name, website, is_moderated, is_authed)
			VALUES (0, 'GGKL', 'ggkl.me', true, false)
			ON CONFLICT DO NOTHING;   

		 CREATE TABLE if not exists pages (
		    id serial NOT NULL,
		    url character varying(500),
		    company_id integer,
		    CONSTRAINT pages_pkey PRIMARY KEY (id));

		 INSERT INTO pages(id,
			url, company_id)
			VALUES (0, '/test', 0)
			ON CONFLICT DO NOTHING;  

		 CREATE TABLE if not exists comments (
		    id serial NOT NULL,
		    page_id integer,
		    user_id integer,
		    text character varying(1000),
		    timestamp timestamp,
		    status character varying(30),
		    important boolean,
		    CONSTRAINT comments_pkey PRIMARY KEY (id));	

		 INSERT INTO public.comments(
			id, page_id, user_id, text, "timestamp", status, important)
			VALUES (-1, 0, 0, 'My first comment!', now(), 'unmodified', false)
			ON CONFLICT DO NOTHING; 

		 INSERT INTO public.comments(
			id, page_id, user_id, text, "timestamp", status, important)
			VALUES (0, 0, 0, 'My second comment!', now(), 'published', false)
			ON CONFLICT DO NOTHING;	  
		 `

   _, err2 := DB.Query(sql)

   if err2 != nil {
	 log.Println("DB Error users", err)
   }	

}
