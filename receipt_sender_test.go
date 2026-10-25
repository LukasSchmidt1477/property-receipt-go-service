package main

import (
	"strings"
	"testing"
	"time"
)

func TestReceiptHTMLKeepsComplianceDetails(t *testing.T) {
	cases := []struct {
		name  string
		order Order
		want  []string
	}{
		{"complete", Order{ID: "ord-7", AmountCents: 1250, Maintenance: []MaintenanceRequest{{ID: "m1"}}, Documents: []TenantDocument{{Name: "lease"}}, Inspection: InspectionReminder{Due: time.Date(2026, 2, 3, 0, 0, 0, 0, time.UTC)}}, []string{"ord-7", "$12.50", "Maintenance items: 1", "documents: 1", "2026-02-03"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := receiptHTML(tc.order)
			for _, want := range tc.want {
				if !strings.Contains(got, want) {
					t.Errorf("receipt missing %q: %s", want, got)
				}
			}
		})
	}
}
