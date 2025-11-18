package telegram

import (
	"bytes"
	"context"
	"embed"
	"text/template"
	"time"

	"github.com/bahmN/rocket-factory/notification/internal/client/http"
	"github.com/bahmN/rocket-factory/notification/internal/config"
	"github.com/bahmN/rocket-factory/notification/internal/model"
)

//go:embed templates/paid_notification.tmpl templates/assembled_notification.tmpl
var FS embed.FS

type orderPaidTemplateData struct {
	OrderUUID       string
	PaymentMethod   string
	TransactionUUID string
	PaymentDate     string
}

type orderAssembledTemplateData struct {
	OrderUUID string
}

var (
	orderPaidTemplate      = template.Must(template.ParseFS(FS, "templates/paid_notification.tmpl"))
	orderAssembledTemplate = template.Must(template.ParseFS(FS, "templates/assembled_notification.tmpl"))
)

type service struct {
	telegramClient http.TelegramClient
}

func NewService(tgClient http.TelegramClient) *service {
	return &service{
		telegramClient: tgClient,
	}
}

func (s *service) SendOrderPaidNotify(ctx context.Context, msg model.OrderPaidEvent) error {
	chatId := config.AppConfig().TgBot.ChatID()

	message, err := s.buildPaidMsg(msg)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatId, message)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) SendOrderAssemblyNotify(ctx context.Context, msg model.OrderAssembledEvent) error {
	chatId := config.AppConfig().TgBot.ChatID()

	message, err := s.buildAssemblyMsg(msg)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatId, message)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) buildPaidMsg(msg model.OrderPaidEvent) (string, error) {
	data := orderPaidTemplateData{
		OrderUUID:       msg.OrderUUID,
		PaymentMethod:   msg.PaymentMethod,
		TransactionUUID: msg.TransactionUUID,
		PaymentDate:     time.Now().Format("2006-01-02 15:04:05"),
	}

	var buf bytes.Buffer
	err := orderPaidTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *service) buildAssemblyMsg(msg model.OrderAssembledEvent) (string, error) {
	data := orderAssembledTemplateData{
		OrderUUID: msg.OrderUUID,
	}

	var buf bytes.Buffer
	err := orderAssembledTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
