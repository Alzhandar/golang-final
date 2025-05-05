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

	"github.com/go-redis/redis/v8"
)

type IikoService struct {
	baseURL      string
	baseURLMenu  string
	apiLogin     string
	httpClient   *http.Client
	redisClient  *redis.Client
	tokenKey     string
	tokenTimeout time.Duration
}

func NewIikoService(apiLogin string, redisClient *redis.Client) *IikoService {
	return &IikoService{
		baseURL:      "https://api-ru.iiko.services/api/1",
		baseURLMenu:  "https://api-ru.iiko.services/api/2",
		apiLogin:     apiLogin,
		httpClient:   &http.Client{Timeout: time.Second * 30},
		redisClient:  redisClient,
		tokenKey:     "iiko_token",
		tokenTimeout: time.Hour,
	}
}

func (s *IikoService) EnsureTokenInRedis(ctx context.Context) (string, error) {
	token, err := s.redisClient.Get(ctx, s.tokenKey).Result()
	if err == redis.Nil || token == "" {
		log.Println("Токен в Redis не найден. Генерируем новый токен")
		token, err = s.GetNewToken(ctx)
		if err != nil {
			return "", fmt.Errorf("ошибка получения нового токена: %w", err)
		}
		err = s.redisClient.Set(ctx, s.tokenKey, token, s.tokenTimeout).Err()
		if err != nil {
			log.Printf("Ошибка сохранения токена в Redis: %v", err)
		}
		log.Printf("Сгенерирован и сохранён новый IIKO Token: %s", token)
	} else {
		log.Printf("Токен уже существует в Redis: %s", token)
	}
	return token, nil
}

func (s *IikoService) GetNewToken(ctx context.Context) (string, error) {
	type tokenRequest struct {
		ApiLogin string `json:"apiLogin"`
	}

	type tokenResponse struct {
		Token string `json:"token"`
	}

	reqData := tokenRequest{ApiLogin: s.apiLogin}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return "", fmt.Errorf("ошибка маршалинга запроса токена: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/access_token", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("неверный код ответа: %d, тело: %s", resp.StatusCode, string(body))
	}

	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	if tokenResp.Token == "" {
		return "", fmt.Errorf("пустой токен в ответе")
	}

	log.Printf("Новый IIKO Token: %s", tokenResp.Token)
	return tokenResp.Token, nil
}

func (s *IikoService) GetTokenAndOrderByTable(ctx context.Context, organizationID string, tableUUID string, statuses []string) (map[string]interface{}, error) {
	if organizationID == "" || tableUUID == "" {
		return nil, fmt.Errorf("organization_id и table_uuid обязательны для получения заказа")
	}

	var token string
	var err error

	token, err = s.redisClient.Get(ctx, s.tokenKey).Result()
	if err == redis.Nil || token == "" {
		log.Println("Токен в Redis отсутствует, генерируем новый перед запросом заказа...")
		token, err = s.GetNewToken(ctx)
		if err != nil {
			return nil, err
		}
		if err = s.redisClient.Set(ctx, s.tokenKey, token, s.tokenTimeout).Err(); err != nil {
			log.Printf("Не удалось сохранить токен в Redis: %v", err)
		}
	} else {
		log.Printf("Токен найден в Redis: %s", token)
	}

	payload := map[string]interface{}{
		"organizationIds": []string{organizationID},
		"tableIds":        []string{tableUUID},
	}

	if len(statuses) > 0 {
		payload["statuses"] = statuses
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга запроса: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/order/by_table", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		log.Println("Получен 401. Пробуем обновить токен и повторить запрос.")
		token, err = s.GetNewToken(ctx)
		if err != nil {
			return nil, fmt.Errorf("ошибка обновления токена: %w", err)
		}

		if err = s.redisClient.Set(ctx, s.tokenKey, token, s.tokenTimeout).Err(); err != nil {
			log.Printf("Не удалось сохранить обновленный токен в Redis: %v", err)
		}

		req, err = http.NewRequestWithContext(ctx, "POST", s.baseURL+"/order/by_table", bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, fmt.Errorf("ошибка создания повторного запроса: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err = s.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("ошибка выполнения повторного запроса: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			log.Println("Получен 401 после обновления токена. Ошибка авторизации.")
			return nil, fmt.Errorf("ошибка авторизации")
		}
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("неверный код ответа: %d, тело: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	log.Println("Запрос текущего счёта (by_table) выполнен успешно.")
	return result, nil
}
