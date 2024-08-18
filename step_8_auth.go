package main

import (
	"crypto/sha256"
	"encoding/json"

	//	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"

	_ "github.com/mattn/go-sqlite3"
)

// струкрура для реквеста при аутенификации
type SignRequest struct {
	Password string `json:"password"`
}

// структура для ответа при аутетификации
type SignResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// формируем переменную секретного ключа для подписи токена
// строка получена с помощью openssl rand -base64 32
var secretKey = []byte("9wsz2ew8lF2pxS4LEg1pHxq9jVhztkKQD5O/5OfvPdE=")

// signinHandler обрабатывает запрос на аутентификацию пользователя
func signinHandler(w http.ResponseWriter, r *http.Request) {
	//Декодировка  JSON-тела запроса в структуру SignRequest
	var req SignRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	//Сравнение введённого пользователем пароля с паролем в переменной окружения TODO_PASSWORD
	if req.Password != os.Getenv("TODO_PASSWORD") {
		resp := SignResponse{Error: "Неверный пароль"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	//Формируется JWT-токен с хэшем пароля в полезной нагрузке
	passwordHash := sha256.Sum256([]byte(req.Password))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"passwordHash": fmt.Sprintf("%x", passwordHash),
	})

	// Подпись токена секретным ключом
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Возвращаем токен в поле token JSON-объекта
	resp := SignResponse{Token: tokenString}
	json.NewEncoder(w).Encode(resp)
}

// auth функция является middleware, который проверяет аутентификацию пользователя.
// принимает один аргумент: next (http.HandlerFunc), который представляет следующий обработчик в цепочке.
func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверка наличия пароля в переменной окружения TODO_PASSWORD
		//по ТЗ проверка аутентификации происходит, только если определён пароль в TODO_PASSWORD
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			//Если пароль определен, функция получает хэш пароля и преобразует его в строковый формат.
			passwordHash := sha256.Sum256([]byte(pass))
			passwordHashString := fmt.Sprintf("%x", passwordHash)

			// Получение JWT-токен из куки запроса
			var jwtToken string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtToken = cookie.Value
			}

			//Парсит JWT-токен и проверяет его валидность
			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				// Проверяем, что метод подписи токена - HMAC
				//это дополнительная проверка от злоумышленников
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				// Возвращаем секретный ключ для проверки подписи токена
				return secretKey, nil
			})
			if err != nil {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
			// Проверка, что токен валиден
			if !token.Valid {
				// Если токен не валиден, возвращается ошибка авторизации с кодом 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			// Извлечение  полезные данные из токена и проверяем хэш пароля
			// приводим поле Claims к типу jwt.MapClaims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok { // если Сlaims вдруг оказжется другого типа, мы получим панику
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
			//Так как jwt.Claims — словарь вида map[string]inteface{}, используем синтакис получения значения по ключу
			hRaw := claims["passwordHash"]
			h, ok := hRaw.(string)
			if !ok {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			if h != passwordHashString {
				// Если извлеченный хэш пароля не совпадает с хэшем пароля в переменной окружения TODO_PASSWORD,
				// возвращаем ошибку авторизации с кодом 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}

		//Если пароль не определен или токен валиден, функция передает запрос следующему обработчику в цепочке.
		next(w, r)

	})
}
