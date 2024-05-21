package pdf

import (
	"context"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	pdftype "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"pdf-service/pdf-tool/internal/svc"
	"pdf-service/pdf-tool/internal/types"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendPdfFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	w      http.ResponseWriter
	r      *http.Request
}

func NewSendPdfFileLogic(ctx context.Context, svcCtx *svc.ServiceContext, w http.ResponseWriter, r *http.Request) *SendPdfFileLogic {
	return &SendPdfFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		w:      w,
		r:      r,
	}
}

func (l *SendPdfFileLogic) SendPdfFile() (resp *types.SendPdfFileReq, err error) {
	file, _, err := l.r.FormFile("file")
	if err != nil {
		logx.Error("the value of file cannot be found")
	}
	defer file.Close()

	path := ""

	ch := make(chan bool)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		path = saveFile(l.r, &wg)
	}()
	wg.Wait()
	addWatermark(path, ch)

	pub()

	return
}

func addWatermark(path string, ch chan<- bool) {
	fmt.Println("get path:", path)
	onTop := false
	update := false
	wm, _ := api.TextWatermark("bellingham", "", onTop, update, pdftype.POINTS)
	if err := api.AddWatermarksFile(path, "", nil, wm, nil); err != nil {
		fmt.Println("add watermark err:", err)
	}
	ch <- true
}

func saveFile(r *http.Request, wg *sync.WaitGroup) string {
	defer wg.Done()
	file, header, err := r.FormFile("file")
	if err != nil {
		logx.Error("the value of file cannot be found")
	}

	// 获取文件名
	fileName := filepath.Base(header.Filename)

	outDir := "./pdfcpu/out"
	pages := []string{"2"}

	// 检查并创建输出目录
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err := os.MkdirAll(outDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating output directory: %s\n", err)
		}
	}

	err = api.ExtractPages(file, outDir, fileName, []string{"2"}, nil)
	if err != nil {
		log.Fatalf("Error extracting pages: %s\n", err)
	}

	baseName := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))

	path := fmt.Sprintf("./pdfcpu/out/%s_page_%s.pdf", baseName, pages[0])
	return path
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func pub() {
	conn, err := amqp.Dial("amqp://guest:guest/url/")
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

	body := "Processing PDF completed"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [pub] Sent %s", body)
}
