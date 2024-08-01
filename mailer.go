package gotil

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"strings"
)

// mailer base model
type Mailer struct {
	email    string
	password string
	host     string
	port     string
	client   *smtp.Client

	buffer *bytes.Buffer
	writer *multipart.Writer

	recipients []string
	subject    string

	cc  []string
	bcc []string

	from string
}

// crate new mailer
func NewMailer(email, password, host, port string) (*Mailer, error) {
	auth := smtp.PlainAuth("", email, password, host)
	client, err := smtp.Dial(host + ":" + port)
	if err != nil {
		return nil, err
	}
	if err := client.StartTLS(&tls.Config{
		InsecureSkipVerify: true,
	}); err != nil {
		return nil, err
	}
	if err := client.Auth(auth); err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	return &Mailer{
		email:    email,
		password: password,
		host:     host,
		port:     port,
		client:   client,
		writer:   writer,
		buffer:   &buffer,
	}, nil
}

// set sender email
func (m *Mailer) From(from string) {
	m.from = from
}

// set cc emails
func (m *Mailer) Cc(emails ...string) {
	m.cc = emails
}

// set bcc emails
func (m *Mailer) Bcc(emails ...string) {
	m.bcc = emails
}

// generate smtp headers
func (m *Mailer) getHeaders() string {
	headers := make(JSON)
	headers["From"] = m.email
	if m.from != "" {
		headers["From"] = fmt.Sprintf("%v <%v>", m.from, m.email)
	}

	headers["To"] = strings.Join(m.recipients, ", ")
	if len(m.cc) > 0 {
		headers["Cc"] = strings.Join(m.cc, ", ")
	}
	headers["Subject"] = m.subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "multipart/mixed; boundary=" + m.writer.Boundary()

	h := ""

	for k, v := range headers {
		h += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	return h + ";\r\n\n"
}

// set mail subject
func (m *Mailer) Subject(subject string) {
	m.subject = subject
}

// set mail recipients
func (m *Mailer) Recipients(emails ...string) {
	m.recipients = emails
}

// set mail text content
func (m *Mailer) SetText(text string) error {
	part, err := m.writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": {"text/plain; charset=\"utf-8\""},
	})
	if err != nil {
		return err
	}
	part.Write([]byte(text))
	return nil
}

// set mail html content
func (m *Mailer) SetHTML(content string) error {
	textPart, err := m.writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": {"text/html; charset=\"utf-8\""},
	})
	if err != nil {
		return err
	}
	textPart.Write([]byte(content))
	return nil
}

// set mail html file with data
func (m *Mailer) SetHTMLFile(file string, data any) error {

	var body bytes.Buffer
	t, err := template.ParseFiles(file)
	if err != nil {
		panic(err)
	}
	_ = t.Execute(&body, data)

	textPart, err := m.writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": {"text/html; charset=\"utf-8\""},
	})
	if err != nil {
		return err
	}
	textPart.Write(body.Bytes())
	return nil
}

// attach file to mail
func (m *Mailer) AttachFile(filename string, data []byte) error {
	part, err := m.writer.CreatePart(textproto.MIMEHeader{
		"Content-Disposition":       {"attachment; filename=\"" + filename + "\""},
		"Content-Type":              {"application/octet-stream"},
		"Content-Transfer-Encoding": {"base64"},
	})
	if err != nil {
		return err
	}

	b64 := base64.NewEncoder(base64.StdEncoding, part)
	if _, err = b64.Write(data); err != nil {
		panic(err)
	}
	b64.Close()
	return nil
}

// send email
func (m *Mailer) Send() error {
	m.writer.Close()
	headers := m.getHeaders()
	content := m.buffer.String()

	m.client.Mail(m.email)

	for _, x := range append(m.recipients, append(m.cc, m.bcc...)...) {
		m.client.Rcpt(x)
	}

	w, err := m.client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(headers + content))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	m.buffer = &buffer
	m.writer = writer
	return nil
}

func (m *Mailer) Close() {
	m.client.Close()
}
