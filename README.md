Remailer
===

Remailer is a simple API server for remailing forms from static sites without an SMTP server.  I use it with single-page AngularJS sites that are hosted on Amazon S3.

It receives form data via post and uses Mailgun to forward it to another address.  You must have your own Mailgun account for this to work, but for most sites this is free.

Configuration is via environment variables:

```
MAILGUN_PRIVATE_KEY   = <your private key>
MAILGUN_PUBLIC_KEY    = <your public key>
MAILGUN_DOMAIN        = <your mailgun domain>

REMAILER_FROM_ADDRESS = <email address you want mail sent from e.g. noreply@yoursite.com>
REMAILER_TO_ADDRESS   = <email address you want form data forwarded to>
REMAILER_SUBJ         = <subject line for emails>
PORT                  = <api port>
```

On your site, you send a POST request to `http://<your_site>/send` with the following JSON payload:
```
{
	"name": <name>
	"email": <email>
	"body": <body>
}
```

Server responses:
```
200:  Successfully queued for sending by Mailgun
400:  Either malformed request or error on Mailgun side
```

By default, it will mail a copy to the originator and your forwarding address.  Note that there is no authentication/authorization included, but since the to address is set server-side, it shouldn't be possible for someone to use your account as a generic remailer.

To build it, you must have Go installed.  It depends on the [Gin](http://github.com/gin-gonic/gin) framework.