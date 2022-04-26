//
// Notification email sender.
//
// by z0gSh1u @ github.com/z0gSh1u/hushbackup
//
package main

import (
	"crypto/tls"
	"time"

	gomail "gopkg.in/gomail.v2"
)

func SendEMailNotification(from string, to string, filepath string, smtp string, port int, username string, password string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	timeNow := time.Now().Format("20060102150405") // stupid design
	m.SetHeader("Subject", "[hushbackup] Notification for backup @ "+timeNow)
	m.SetBody("text/html", "The backup is done successfully with tarball saved at <i>"+filepath+"</i><br><br>https://github.com/z0gSh1u/hushbackup")

	d := gomail.NewDialer(smtp, port, username, password)
	d.TLSConfig = &tls.Config{ServerName: smtp, InsecureSkipVerify: false}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
