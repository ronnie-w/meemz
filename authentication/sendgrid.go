package authentication

import (
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendMail(recipient string, subject string, pTc string, hTc string) {
	from := mail.NewEmail("Meemz", "meemz@mail.com")
	to := mail.NewEmail("Meemz User", recipient)

	message := mail.NewSingleEmail(from, subject, to, pTc, hTc)

	client := sendgrid.NewSendClient("SENDGRID_API_KEY")

	_, err := client.Send(message)
	if err != nil {
		log.Println(err)
	}

}
