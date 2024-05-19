package communication

import (
	"CurrencyChecking/config"
	"CurrencyChecking/database"
	"CurrencyChecking/service"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"net/smtp"
	"time"
)

// SMTP server configuration for Gmail
const (
	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
)

type UserEmailSender struct {
	userService service.UserServiceInterface
	DbUser      database.DbInterface
}

func NewUserEmailSender(service service.UserServiceInterface, DbUser database.DbInterface) *UserEmailSender {
	return &UserEmailSender{userService: service, DbUser: DbUser}
}

func (ues *UserEmailSender) ScheduleEmailSender() {
	c := cron.New()

	todayDate := time.Now()
	// Format the date to "YYYY-MM-DD"
	date := todayDate.Format("2006-01-02")
	time.Now().Date()
	subject := "Currency rate. Update " + date
	body, err, _ := ues.userService.GetRate()
	// Schedule the task to run at a specific time every day
	emails, err := ues.DbUser.SelectUsersEmail()
	if err != nil {

	}
	_, err = c.AddFunc("0 18 * * *", func() {
		// You can adjust the recipient, subject, and body as needed
		for _, email := range emails {
			sendEmail(email, subject, body)
		}
	})
	if err != nil {
		log.Fatalf("Error scheduling email: %v", err)
	}

	// Start the cron scheduler
	c.Start()

	log.Println("Email scheduler started successfully")
}

func sendEmail(to, subject, body string) error {
	cfg := config.LoadENV(".env")

	message := fmt.Sprintf("From: %s\r\n", cfg.Email) +
		fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"\r\n" +
		body

	// Set up authentication information
	auth := smtp.PlainAuth("", cfg.Email, cfg.EmailPassword, smtpHost)

	// Connect to the SMTP server
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, cfg.Email, []string{to}, []byte(message))
	if err != nil {
		return err
	}

	return nil
}
