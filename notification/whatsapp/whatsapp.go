package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mtrossbach/waechter/internal/cfg"
	"github.com/mtrossbach/waechter/internal/log"
	"io"
	"net/http"
	"time"
)

type WhatsApp struct {
	client *http.Client
}

func NewWhatsApp() *WhatsApp {
	return &WhatsApp{client: &http.Client{
		Timeout: 60 * time.Second,
	}}
}

func (w *WhatsApp) send(phone string, template string, lang string, parameters []string) error {
	if len(phone) < 5 {
		log.Error().Str("phone", phone).Msg("Could not send WhatsApp message. Invalid phone number")
		return fmt.Errorf("invalid phone number: %v", phone)
	}
	if len(template) < 1 {
		log.Error().Str("template", template).Msg("Could not send WhatsApp message. Invalid template name")
		return fmt.Errorf("invalid template name: %v", template)
	}
	var ps []Parameter
	for _, s := range parameters {
		ps = append(ps, Parameter{
			Type: "text",
			Text: s,
		})
	}
	payload := MessagePayload{
		MessagingProduct: "whatsapp",
		To:               phone,
		Type:             "template",
		Template: Template{
			Name:     template,
			Language: Language{Code: lang},
			Components: []Component{{
				Type:       "body",
				Parameters: ps,
			}},
		},
	}

	var response interface{}

	r, err := w.post(cfg.GetString(cPhoneId), payload, &response)
	if err != nil {
		log.Error().Err(err).Str("phone", phone).Msg("Could not send WhatsApp message")
		return err
	}
	if r.StatusCode >= 300 {
		log.Error().Str("phone", phone).Int("status-code", r.StatusCode).Msg("Could not send WhatsApp message")
		return fmt.Errorf("could not send message to whatsapp, statuscode is %v", r.StatusCode)
	}
	log.Info().Str("phone", phone).Msg("Successfully sent message via WhatsApp")
	return nil
}

func (w *WhatsApp) post(phoneId string, payload MessagePayload, response interface{}) (*http.Response, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v13.0/%v/messages", phoneId)

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", cfg.GetString(cToken)))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return resp, nil
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, err
	}

	_ = resp.Body.Close()
	return resp, nil
}

/*
func (w *WhatsApp) NotifyAlarm(recipient system.Recipient, systemName string, alarmType system.AlarmState, device *device.Device) bool {
	err := w.send(recipient.Phone, cfg.GetString(cAlarmTemplateName), recipient.Lang, []string{
		systemName, i18n.TranslateAlarm(recipient.Lang, alarmType), device.Name,
	})

	return err == nil
}

func (w *WhatsApp) NotifyRecovery(recipient system.Recipient, systemName string, device *device.Device) bool {
	err := w.send(recipient.Phone, cfg.GetString(cRecoverTemplateName), recipient.Lang, []string{
		systemName,
	})

	return err == nil
}

func (w *WhatsApp) NotifyLowBattery(recipient system.Recipient, systemName string, device *device.Device, batteryLevel float32) bool {
	err := w.send(recipient.Phone, cfg.GetString(cNotificationTemplateName), recipient.Lang, []string{
		systemName, device.Name, i18n.Translate(recipient.Lang, i18n.WALowBattery),
	})

	return err == nil
}

func (w *WhatsApp) NotifyLowLinkQuality(recipient system.Recipient, systemName string, device *device.Device, quality float32) bool {
	err := w.send(recipient.Phone, cfg.GetString(cNotificationTemplateName), recipient.Lang, []string{
		systemName, device.Name, i18n.Translate(recipient.Lang, i18n.WALowLinkQuality),
	})

	return err == nil
}

func (w *WhatsApp) NotifyAutoArm(recipient system.Recipient, systemName string) bool {
	err := w.send(recipient.Phone, cfg.GetString(cAutoArmTemplateName), recipient.Lang, []string{
		systemName,
	})

	return err == nil
}

func (w *WhatsApp) NotifyAutoDisarm(recipient system.Recipient, systemName string) bool {
	err := w.send(recipient.Phone, cfg.GetString(cAutoDisarmTemplateName), recipient.Lang, []string{
		systemName,
	})

	return err == nil
}
*/
