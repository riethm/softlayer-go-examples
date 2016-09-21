package main

import (
	"fmt"
	"log"
	"time"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/helpers/order"
	"github.com/softlayer/softlayer-go/helpers/product"
	"github.com/softlayer/softlayer-go/helpers/virtual"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

const guestID = 24459505

func main() {
	sess := session.New()
	sess.Debug = true

	// Create a minimal Virtual_Guest object to pass to the upgrade helper
	guestToUpgrade := datatypes.Virtual_Guest{
		Id: sl.Int(guestID),
	}

	// Upgrade to 4 Core, 8 GB
	upgradeOptions := map[string]float64{
		product.CPUCategoryCode:    float64(4),
		product.MemoryCategoryCode: float64(8),
	}

	receipt, err := virtual.UpgradeVirtualGuest(sess, &guestToUpgrade, upgradeOptions)
	if err != nil {
		log.Fatal("Couldn't upgrade virtual guest:", err)
	}

	fmt.Println("Virtual Guest upgrade order submitted")

	// Wait up to 5 minutes for the transaction to complete
	fmt.Println("Waiting up to 5 minutes for order to complete...")
	complete := false
	for i := 0; i < 10; i++ {
		fmt.Println("Sleeping 30 seconds...")
		time.Sleep(30 * time.Second)
		complete, _, err = order.CheckBillingOrderComplete(sess, &receipt)

		if err != nil {
			log.Fatal("Error while checking upgrade order status:", err)
		}

		if complete {
			break
		}
		fmt.Printf("Not complete.  ")
	}
	if !complete {
		fmt.Println("Upgrade transaction did not complete within the allotted time")
	}
}
