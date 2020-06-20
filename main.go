package main

import (
    "fmt"
    "net"
    "net/http"
    "net/smtp"
    "os"
    "strings"
    "time"

    "github.com/jordan-wright/email"
    "github.com/sirupsen/logrus"
)

func main() {
    url := os.Getenv("URL")
    callbackHTTPUrl := os.Getenv("CALLBACK_HTTP_URL")
    callbackHTTPMethod := os.Getenv("CALLBACK_HTTP_METHOD")
    callbackSMTPFrom := os.Getenv("CALLBACK_SMTP_FROM")
    callbackSMTPTo := os.Getenv("CALLBACK_SMTP_TO")
    callbackSMTPCC := os.Getenv("CALLBACK_SMTP_CC")
    callbackSMTPBCC := os.Getenv("CALLBACK_SMTP_BCC")
    callbackSMTPSubject := os.Getenv("CALLBACK_SMTP_SUBJECT")
    callbackSMTPUser := os.Getenv("CALLBACK_SMTP_USER")
    callbackSMTPPassword := os.Getenv("CALLBACK_SMTP_PASSWORD")
    callbackSMTPHost := os.Getenv("CALLBACK_SMTP_HOST")
    callbackSMTPPort := os.Getenv("CALLBACK_SMTP_PORT")
    callbackSMTPBody := os.Getenv("CALLBACK_SMTP_BODY")
    callbackSMTPAttachmentName := os.Getenv("CALLBACK_SMTP_ATTACHMENTNAME")

    if url == "" {
        logrus.Fatal("missing URL")
    }

    if callbackHTTPMethod == "" {
        callbackHTTPMethod = "POST"
    }

    if callbackSMTPAttachmentName == "" {
        callbackSMTPAttachmentName = "document.pdf"
    }

    if callbackHTTPUrl != "" && callbackSMTPTo != "" {
        logrus.Fatal("ambigious configuration: both CALLBACK_HTTP_URL and CALLBACK_SMTP_TO set")
    }

    if callbackHTTPUrl != "" {
        httpClient := &http.Client{
            Timeout: time.Second * 30,
            Transport: &http.Transport{
                Dial:                (&net.Dialer{Timeout: 5 * time.Second}).Dial,
                TLSHandshakeTimeout: 5 * time.Second,
            },
        }

        req, err := http.NewRequest(callbackHTTPMethod, url, os.Stdin)
        if err != nil {
            logrus.Fatalf("couldn't create http request, see: %v", err)
        }

        req.Header.Set("Content-Type", "application/pdf")

        resp, err := httpClient.Do(req)
        if err != nil {
            logrus.Fatalf("couldn't send request, see: %v", err)
        }

        if resp.StatusCode < 200 || resp.StatusCode > 299 {
            logrus.Fatal("non-success http status code")
        }

        logrus.Infof("pdf sent to %s", url)
    } else if callbackSMTPTo != "" {
        e := email.NewEmail()
        e.From = callbackSMTPFrom
        e.To = splitTrimSpace(callbackSMTPTo, ",")
        e.Bcc = splitTrimSpace(callbackSMTPBCC, ",")
        e.Cc = splitTrimSpace(callbackSMTPCC, ",")
        e.Subject = callbackSMTPSubject
        e.Text = []byte(callbackSMTPBody)
        _, err := e.Attach(os.Stdin, callbackSMTPAttachmentName, "application/pdf")
        if err != nil {
            logrus.Fatal("couldn't add attachment, see: %v", err)
        }

        err = e.Send(fmt.Sprintf("%s:%s", callbackSMTPHost, callbackSMTPPort), LoginAuth(callbackSMTPUser, callbackSMTPPassword))
        if err != nil {
            logrus.Fatalf("couldn't send mail, see: %v", err)
        }

        logrus.Infof("pdf sent to %s (cc: %s) (bcc: %s)", callbackSMTPTo, callbackSMTPCC, callbackSMTPBCC)
    } else {
        logrus.Fatal("neither CALLBACK_HTTP_URL nor CALLBACK_SMTP_TO given")
    }
}

func splitTrimSpace(str string, delimiter string) []string {
    splitted := strings.Split(str, delimiter)
    ret := make([]string, 0)
    for _, s := range splitted {
        s = strings.TrimSpace(s)

        if s != "" {
            ret = append(ret, s)
        }
    }

    return ret
}

type loginAuth struct {
    username, password string
}

func LoginAuth(username, password string) smtp.Auth {
    return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
    return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
    if more {
        switch string(fromServer) {
        case "Username:":
            return []byte(a.username), nil
        case "Password:":
            return []byte(a.password), nil
        default:
            return nil, fmt.Errorf("Unkown fromServer: %s", fromServer)
        }
    }
    return nil, nil
}
