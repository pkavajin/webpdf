# webpdf

webpdf is a small utility which generates a pdf from a website and sends it 
either via smtp or via web call.
It can either be used via docker directly or deployed as k8s cronjob.

## Run via docker

Generate pdf for google.de and send report via http to example.org:

```
docker run -it --rm                           \
    -e URL="https://google.de"                \
    -e CALLBACK_HTTP_URL="http://example.org" \
    -e CALLBACK_HTTP_METHOD="POST"            \
    kavatech/webpdf:v0.1.0
```

Generate pdf for google.de and send report via mail:

```
docker run -it --rm                                        \
    -e URL="https://google.de"                             \
    -e CALLBACK_SMTP_FROM="Foo <foo@example.org>"          \
    -e CALLBACK_SMTP_TO="bar@example.org"                  \
    -e CALLBACK_SMTP_CC="qux@example.org,baz@example.org"  \
    -e CALLBACK_SMTP_BCC="ninja@example.org"               \
    -e CALLBACK_SMTP_SUBJECT="Some Report"                 \
    -e CALLBACK_SMTP_USER="foo"                            \
    -e CALLBACK_SMTP_PASSWORD="secret123"                  \
    -e CALLBACK_SMTP_HOST="mail.example.org"               \
    -e CALLBACK_SMTP_PORT="587"                            \
    -e CALLBACK_SMTP_BODY="See attachment"                 \
    -e CALLBACK_SMTP_ATTACHMENTNAME="report.pdf"           \
    kavatech/webpdf:v0.1.0
```

## Deploy using helm

Example cronjob generating report of google.de once a day and send it to example.org:
```
cd ./charts/webpdf && helm upgrade -i webpdf --namespace webpdf ./  \
    --set url='https://google.de'                                   \
    --set callback.http.url='https://example.org'                   \
    --set schedule='0 1 * * *'
```

Specify either callback.http or callback.smtp

### Values:

| Key                          | Type     | Default                  | Description                                            |
|------------------------------|----------|--------------------------|--------------------------------------------------------|
| image                        | string   | "kavatech/webpdf:v0.1.0" | Image which should be used for the cronjob             |
| schedule                     | string   | ""                       | Crontab like schedule for cronjob                      |
| restartPolicy                | string   | "OnFailure"              | Cronjob's restart policy                               |
| url                          | string   | ""                       | URL for which the pdf should be generated              |
| callback.http.url            | string   | ""                       | URL of the webhook which should be called with the pdf |
| callback.http.method         | string   | "POST"                   | HTTP method for the webhook call                       |
| callback.smtp.from           | string   | ""                       | E-Mail address of the sender                           |
| callback.smtp.to             | []string | []                       | E-Mail address of the recipients                       |
| callback.smtp.cc             | []string | []                       | CC e-mail addresses                                    |
| callback.smtp.bcc            | []string | []                       | BCC e-mail addresses                                   |
| callback.smtp.subject        | string   | "webpdf report"          | E-Mail subject                                         |
| callback.smtp.user           | string   | ""                       | SMTP Username                                          |
| callback.smtp.password       | string   | ""                       | SMTP password                                          |
| callback.smtp.host           | string   | ""                       | SMTP hostname                                          |
| callback.smtp.port           | string   | "587"                    | SMTP port                                              |
| callback.smtp.body           | string   | "see attachment"         | E-Mail body                                            |
| callback.smtp.attachmentName | string   | "report.pdf"             | Filename of the attachment                             |

## License

Licensed under [MIT License](./LICENSE.md)
Uses [wkhtmltopdf](https://wkhtmltopdf.org/) (LGPLv3) for pdf generation, source is available [here](https://github.com/wkhtmltopdf/wkhtmltopdf).