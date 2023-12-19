package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/nats-io/stan.go"
)

const clusterID = "test-cluster"
const clientID = "producer-client"
const subject = "zxc"

// var verifier *oidc.IDTokenVerifier

// func init() {
// 	ctx := context.Background()
// 	provider, err := oidc.NewProvider(ctx, "http://keycloak:9080/auth/realms/example-realm")
// 	if err != nil {
// 		log.Fatalf("Failed to get provider: %v", err)
// 	}

// 	oidcConfig := &oidc.Config{
// 		ClientID: "my-client", // Замените на ID вашего клиента в Keycloak
// 	}
// 	verifier = provider.Verifier(oidcConfig)
// }

// func authenticate(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		rawToken := r.Header.Get("Authorization")
// 		// Проверка наличия токена и его формата
// 		if rawToken == "" || !strings.HasPrefix(rawToken, "Bearer ") {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}
//zczczczczcz
// 		// Извлечение и проверка токена
// 		token, err := verifier.Verify(r.Context(), strings.TrimPrefix(rawToken, "Bearer "))
// 		if err != nil {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}

// 		// Верификация прошла успешно, продолжаем обработку запроса
// 		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", token)))
// 	})
// }

func main() {
	var tmpl = template.Must(template.New("order").Parse(htmlTemplate))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		tmpl.Execute(w, nil) // Отобразить форму, если метод не POST

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusInternalServerError)
			return
		}
		id := r.FormValue("id")

		order, err := loadOrderAndUpdateUID(id) // Загрузка и обновление заказа
		if err != nil {
			http.Error(w, "Error loading order", http.StatusInternalServerError)
			return
		}
		jsonData, err := json.Marshal(order) // Сериализация обновлённого заказа в JSON
		if err != nil {
			http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
			return
		}

		err = publishToNATS(jsonData) // Отправка данных в NATS
		if err != nil {
			http.Error(w, "Error publishing to NATS", http.StatusInternalServerError)
			return
		}

		log.Println("sended message with id:", order.OrderUID)
	})

	log.Println("Server started at http://localhost:8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Error ListenAndServe: ", err)
	}
}

func loadOrderAndUpdateUID(orderUID string) (Order, error) {
	data, err := ioutil.ReadFile("/ord.json") // Убедитесь, что путь к файлу указан верно
	if err != nil {
		return Order{}, err
	}

	var order Order
	err = json.Unmarshal(data, &order)
	if err != nil {
		return Order{}, err
	}

	order.OrderUID = orderUID      // Обновление OrderUID согласно введённому пользователем значению
	order.DateCreated = time.Now() // Обновление даты создания заказа

	return order, nil
}

func publishToNATS(data []byte) error {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("nats://nats-streaming:4222"))
	if err != nil {
		return err
	}
	defer sc.Close()

	err = sc.Publish(subject, data)
	if err != nil {
		return err
	}
	return nil
}
