package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {

		// Random Keypair
		pair := generateKeypair()
		// Generate Account
		if err := createTestnetAccount(pair); err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"message": err,
			})
		}

		// Send Fund
		destination := c.Query("addr")
		source := pair.Seed()

		// Make sure destination account exists
		if _, err := horizon.DefaultTestNetClient.LoadAccount(destination); err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"message": err,
			})
		}
		tx, err := build.Transaction(
			build.TestNetwork,
			build.SourceAccount{source},
			build.AutoSequence{horizon.DefaultTestNetClient},
			build.Payment(
				build.Destination{destination},
				build.NativeAmount{"9000"},
			),
		)
		txe, err := tx.Sign(source)
		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"message": err,
			})
		}
		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"message": err,
			})
		}
		txeB64, err := txe.Base64()
		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"message": err,
			})
		}

		// And finally, send it off to Stellar!
		resp, err := horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"message": err,
			})
		}
		c.JSON(200, gin.H{
			"success": true,
			"message": resp,
		})
	})
	router.Run(":80") // listen and serve on 0.0.0.0:8080
}

func generateKeypair() *keypair.Full {
	kp, _ := keypair.Random()
	return kp
}

func createTestnetAccount(pair *keypair.Full) error {
	address := pair.Address()
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + address)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}
