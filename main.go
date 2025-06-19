package main

import (
	"fmt"
	"os"
	"log"

	"gopkg.in/yaml.v3"
)

type Tip struct {
	Amount  *float64 `yaml:"amount"`
	Percent *float64 `yaml:"percent"`
}

type Item struct {
	Name      string   `yaml:"name"`
	Price     float64  `yaml:"price"`
	Quantity  int      `yaml:"quantity"`   // optional, defaults to 1 if 0
	Attendees []string `yaml:"attendees"`
}


type Receipt struct {
	Restaurant string  `yaml:"restaurant"`
	Tax        float64 `yaml:"tax"`
	Surcharge  float64 `yaml:"surcharge"`
	Tip        Tip     `yaml:"tip"`
	Items      []Item  `yaml:"items"`
}

func main() {
	data, err := os.ReadFile("receipt.yaml")
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	var receipt Receipt
	err = yaml.Unmarshal(data, &receipt)
	if err != nil {
		log.Fatalf("error parsing YAML: %v", err)
	}

	subtotals := make(map[string]float64)
	var total float64

	for _, item := range receipt.Items {
		quantity := item.Quantity
		if quantity == 0 {
			quantity = 1
		}
	
		totalItemCost := item.Price * float64(quantity)
		splitPrice := totalItemCost / float64(len(item.Attendees))
	
		for _, person := range item.Attendees {
			subtotals[person] += splitPrice
		}
		total += totalItemCost
	}
	

	// Calculate tip value from amount or percent
	var tipValue float64
	if receipt.Tip.Amount != nil {
		tipValue = *receipt.Tip.Amount
	} else if receipt.Tip.Percent != nil {
		tipValue = (*receipt.Tip.Percent / 100) * total
	} else {
		tipValue = 0.0
	}

	var grandTotal float64
	fmt.Println("Receipt Breakdown:")
	for person, subtotal := range subtotals {
		shareRatio := subtotal / total
		taxShare := shareRatio * receipt.Tax
		tipShare := shareRatio * tipValue
		surchargeShare := shareRatio * receipt.Surcharge
		totalDue := subtotal + taxShare + tipShare + surchargeShare

		grandTotal += totalDue

		fmt.Printf("\n%s owes:\n", person)
		fmt.Printf("  Subtotal:  $%.2f\n", subtotal)
		fmt.Printf("  Tax:       $%.2f\n", taxShare)
		fmt.Printf("  Tip:       $%.2f\n", tipShare)
		fmt.Printf("  Surcharge: $%.2f\n", surchargeShare)
		fmt.Printf("  Total:     $%.2f\n", totalDue)
	}

	fmt.Printf("\n🔢 Grand Total (all people): $%.2f\n", grandTotal)
	fmt.Printf("💳 Expected Total (from receipt): $%.2f\n", total + receipt.Tax + tipValue + receipt.Surcharge)

}
