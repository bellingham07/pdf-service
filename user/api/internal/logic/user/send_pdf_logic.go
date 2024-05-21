package user

import (
	"bytes"
	"context"
	"github.com/HiBugEnterprise/gotools/httpc"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"pdf-service/user/api/internal/svc"
)

type SendPdfLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewSendPdfLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *SendPdfLogic {
	return &SendPdfLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *SendPdfLogic) SendPdf() error {
	file, header, err := l.r.FormFile("file")
	if err != nil {
		return err
	}

	err = UploadFile(filepath.Base(header.Filename), file)
	if err != nil {
		return err
	}

	return nil
}

type FileUploadResp struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func UploadFile(fileName string, file multipart.File) (err error) {
	URL := "http://127.0.0.1:8888/pdf/send"

	var reqBody bytes.Buffer
	resp := new(FileUploadResp)
	writer := multipart.NewWriter(&reqBody)

	filePart, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		err = errors.Wrap(err, "Error creating form file")
		return
	}

	// copy the file content directly to the form field
	_, err = io.Copy(filePart, file)
	if err != nil {
		err = errors.Wrap(err, "Error copying file")
		return
	}

	writer.Close()

	// send upload request
	if _, err = httpc.Post(URL,
		httpc.SetHeader("Content-Type", writer.FormDataContentType()),
		httpc.SetBody(&reqBody), httpc.SetResult(&resp)); err != nil {
		err = errors.Wrap(err, "upload file tools error")
		return
	}

	if resp == nil {
		err = errors.New("hibug-file-svc interface exception")
		return
	}

	if resp.Code != 200 {
		err = errors.Wrap(errors.New(resp.Msg), "file upload failed")
		return
	}

	sub()

	return
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func sub() {
	conn, err := amqp.Dial("amqp://guest:guest@url/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
