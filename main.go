package main

import "github.com/gin-gonic/gin"
import "github.com/mailgun/mailgun-go"
import "os"
import "fmt"
import "github.com/BTBurke/remailer/middleware"


type Inquiry struct {
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Body string `json:"body"`
}

func bindmailgunSend() func(c *gin.Context) {
	privkey := os.Getenv("MAILGUN_PRIVATE_KEY")
	pubkey := os.Getenv("MAILGUN_PUBLIC_KEY")
	domain  := os.Getenv("MAILGUN_DOMAIN")
	redirect := os.Getenv("REMAILER_TO_ADDRESS")
	sender := os.Getenv("REMAILER_FROM_ADDRESS")
	subj := os.Getenv("REMAILER_SUBJ")


	if (privkey == "" || pubkey == "" || domain == "" || redirect == "" || sender == "" || subj == "") {
		fmt.Printf("ERROR: Environment variables not properly set.\nMAILGUN_PRIVATE_KEY: %v\nMAILGUN_PUBLIC_KEY: %v\nMAILGUN_DOMAIN: %v\nREMAILER_TO_ADDRESS: %v\nREMAILER_FROM_ADDRESS: %v\nREMAILER_SUBJ: %v\n", privkey, pubkey, domain, redirect, sender, subj)
		os.Exit(1)
	}

	return func(c *gin.Context) {
		from := sender
		to1 := redirect
		sub := subj
		gun := mailgun.NewMailgun(domain, privkey, pubkey)
		var inq Inquiry

		c.Bind(&inq)

		// Handle no email case
		if len(inq.Email) == 0 {
			c.String(400, "failed no email")
			c.Set("email", "User provided no email address")
			c.Set("body", inq.Body)
			return
		} 
		to2 := inq.Name + " <" + inq.Email + ">"
		c.Set("email", to2)
		

		// Handle case where user sends no message body
		if len(inq.Body) == 0 {
			inq.Body = "User sent zero length email body."
		}
		body := "From: " + to2 + "\n\n" + inq.Body
		c.Set("body", inq.Body)

		m := mailgun.NewMessage(from, sub, body, to1, to2)
		response, id, err := gun.Send(m)
		c.Set("response", response)
		c.Set("id", id)
		if err != nil {
			c.String(400, "failed")
			fmt.Printf("Response: %v\nID: %v\nError: %v\n", response, id, err)
		} else {
			c.String(200, "success")
		}
	}	
}



func main() {
	port := os.Getenv("MAILER_API_PORT")
	if len(port) == 0 {
		fmt.Println("ERROR: Port not specified. Set in MAILER_API_PORT.")
		os.Exit(1)
	}
	mailgunSend := bindmailgunSend()

	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(middleware.CORSaccept())
	r.Use(gin.Recovery())


	r.GET("/ping", func(c *gin.Context){
		c.String(200, "pong")
	})
	r.POST("/send", mailgunSend)

	r.Run(":"+port)
	
}