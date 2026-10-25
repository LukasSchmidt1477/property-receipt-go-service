package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	client, err := NewEmailClient()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/receipts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var order Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		id, err := SendReceipt(client, order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"message_id": id})
	})
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
