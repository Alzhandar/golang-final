package iiko

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type IikoWaiterService struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewIikoWaiterService(apiURL, apiKey string) *IikoWaiterService {
	return &IikoWaiterService{
		baseURL:    apiURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: time.Second * 5},
	}
}

func (s *IikoWaiterService) makeRequest(ctx context.Context, endpoint string, payload interface{}) (map[string]interface{}, error) {
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга запроса: %w", err)
	}

	url := fmt.Sprintf("%s/notifications/mobile/%s", s.baseURL, endpoint)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("неверный код ответа: %d, тело: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	if len(body) == 0 {
		return map[string]interface{}{"status": "success"}, nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	return result, nil
}

func (s *IikoWaiterService) CallWaiter(ctx context.Context, departmentID string, tableNumber int, userID string) (map[string]interface{}, error) {
	if departmentID == "" {
		return nil, fmt.Errorf("некорректный department_id: %s", departmentID)
	}

	payload := map[string]interface{}{
		"departmentId": departmentID,
		"table":        tableNumber,
		"userId":       userID,
	}

	log.Printf("Отправка запроса вызова официанта: %v", payload)
	return s.makeRequest(ctx, "waiter-call", payload)
}

func (s *IikoWaiterService) RequestCashPayment(ctx context.Context, userID, departmentID string, tableNumber int) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"departmentId": departmentID,
		"table":        tableNumber,
		"userId":       userID,
	}

	log.Printf("Отправка запроса на оплату наличными: %v", payload)
	return s.makeRequest(ctx, "waiter-call/cash-payment", payload)
}

func (s *IikoWaiterService) RequestCardPayment(ctx context.Context, userID, departmentID string, tableNumber int) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"departmentId": departmentID,
		"table":        tableNumber,
		"userId":       userID,
	}

	log.Printf("Отправка запроса на оплату картой: %v", payload)
	return s.makeRequest(ctx, "waiter-call/card-payment", payload)
}

func (s *IikoWaiterService) NotifyNewOrder(ctx context.Context, userID, departmentID string, tableNumber int) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"departmentId": departmentID,
		"tableNumber":  tableNumber,
		"userId":       userID,
	}

	return s.makeRequest(ctx, "new-order", payload)
}

func (s *IikoWaiterService) BroadcastNewOrder(ctx context.Context, departmentID string, tableNumber int) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"departmentId": departmentID,
		"tableNumber":  tableNumber,
	}

	return s.makeRequest(ctx, "new-order-broadcast", payload)
}

func (s *IikoWaiterService) OrderPaid(ctx context.Context, userID, departmentID string, orderNumber, tableNumber int, orderID string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"userId":       userID,
		"departmentId": departmentID,
		"orderNumber":  orderNumber,
		"orderId":      orderID,
		"tableNumber":  tableNumber,
		"leftToPay":    0,
	}

	log.Printf("Отправка запроса на оплату заказа: %v", payload)
	return s.makeRequest(ctx, "order-paid", payload)
}
