package auth

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	_ "github.com/mattn/go-sqlite3"

	"github.com/O1V1/go_final_project/pkg/controller/config"
)

var (
	secretKey    = config.SecretKey
	todoPassword = config.TodoPassword
)

// для авторизации
type AuthMiddleware struct{}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

// auth функция является middleware, который проверяет авторизацию пользователя.
// принимает один аргумент: next (http.HandlerFunc), который представляет следующий обработчик в цепочке.
func (m *AuthMiddleware) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//по ТЗ проверка аутентификации происходит, только если определён пароль в TODO_PASSWORD
		if len(todoPassword) > 0 {
			//Если пароль определен, функция получает хэш пароля и преобразует его в строковый формат.
			passwordHash := sha256.Sum256([]byte(todoPassword))
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

// для аутентификации
type (
	// струкрура для реквеста при аутенификации
	SignRequest struct {
		Password string `json:"password"`
	}

	// структура для ответа при аутентификации
	SignResponse struct {
		Token string `json:"token,omitempty"`
		Error string `json:"error,omitempty"`
	}

	//структура обработчика аутентификации
	AuthHandlerImpl struct {
		todoPassword string
		secretKey    []byte
	}
)

// конструктор нового экземпляра структуры AuthHandlerImpl
func NewAuthHandler(todoPassword string, secretKey []byte) *AuthHandlerImpl {
	return &AuthHandlerImpl{todoPassword: todoPassword, secretKey: secretKey}
}

// signinHandler обрабатывает запрос на аутентификацию пользователя
func (h *AuthHandlerImpl) SigninHandler(w http.ResponseWriter, r *http.Request) {
	//Декодировка  JSON-тела запроса в структуру SignRequest
	var req SignRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	//Сравнение введённого пользователем пароля с паролем в переменной окружения TODO_PASSWORD
	if req.Password != h.todoPassword {
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
	tokenString, err := token.SignedString(h.secretKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Возвращаем токен в поле token JSON-объекта
	resp := SignResponse{Token: tokenString}
	json.NewEncoder(w).Encode(resp)
}
