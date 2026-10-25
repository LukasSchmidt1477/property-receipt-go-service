package main

import (
	"fmt"
	"time"
)

type MaintenanceRequest struct{ ID, Summary string }
type TenantDocument struct{ Name, URL string }
type InspectionReminder struct{ Due time.Time }
type Order struct {
	ID, TenantEmail string
	AmountCents     int
	Maintenance     []MaintenanceRequest
	Documents       []TenantDocument
	Inspection      InspectionReminder
}

func receiptHTML(o Order) string {
	return fmt.Sprintf("<p>Receipt for order %s: $%.2f</p><p>Maintenance items: %d; documents: %d; inspection due: %s</p>", o.ID, float64(o.AmountCents)/100, len(o.Maintenance), len(o.Documents), o.Inspection.Due.Format("2006-01-02"))
}

func SendReceipt(client *EmailClient, o Order) (string, error) {
	if o.TenantEmail == "" || o.ID == "" || o.AmountCents < 0 {
		return "", fmt.Errorf("invalid order")
	}
	return client.Send(o.TenantEmail, "Property receipt "+o.ID, receiptHTML(o))
}
