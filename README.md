## Official Anti-Captcha.com npm module ##

Official anti-captcha.com Golang library for solving images with text, Recaptcha v2/v3 Enterprise/non-Enterpise, Funcaptcha, GeeTest, HCaptcha Enterprise/non-Enterprise.

[Anti-captcha](https://anti-captcha.com) is an oldest and cheapest web service dedicated to solving captchas by human workers from around the world. By solving captchas with us you help people in poorest regions of the world to earn money, which not only cover their basic needs, but also gives them ability to financially help their families, study and avoid jobs where they're simply not happy.

To use the service you need to [register](https://anti-captcha.com/clients/) and topup your balance. Prices start from $0.0005 per image captcha and $0.002 for Recaptcha. That's $0.5 per 1000 for images and $2 for 1000 Recaptchas.

For more technical information and articles visit our [documentation](https://anti-captcha.com/apidoc) page. 

### Solve image captcha with Go example:
```go
package main

import (
	"fmt"
	// Sorry for this many "anticaptcha" in one import :-D,
	// otherwise the package would be named "anticaptcha_go" and that's ugly
	"github.com/anti-captcha/anticaptcha-go/anticaptcha"
	"log"
)

func main() {
    // Create API client and set the API Key
	ac := anticaptcha.NewClient("API_KEY_HERE")

	// set to 'false' to turn off debug output
	ac.IsVerbose = true

	// Specify softId to earn 10% commission with your app.
	// Get your softId here: https://anti-captcha.com/clients/tools/devcenter
	//ac.SoftId = 1187

    // Make sure the API key funds balance is positive
    balance, err := ac.GetBalance()
    if err != nil {
        log.Fatal(err)
        // Exit program to make sure you don't DDoS API with requests, while having empty balance
        return
    }
    fmt.Println("Balance:", balance)
    
    // Solve image captcha
    solution, err := ac.SolveImageFile("captcha.jpg", anticaptcha.ImageSettings{
	    CaseSensitive: true,
	    MaxLength:     5,
    })
    // OR 
    // solution, err := ac.SolveImage("image-encoded-in-base64", anticaptcha.ImageSettings{})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Captcha Solution:", solution)
}

```