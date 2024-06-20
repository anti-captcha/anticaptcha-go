## Official Anti-Captcha.com Go package ##

Official anti-captcha.com Golang package for solving images with text, Recaptcha v2/v3 Enterprise/non-Enterpise, Funcaptcha, GeeTest, HCaptcha Enterprise/non-Enterprise.

[Anti-captcha](https://anti-captcha.com) is an oldest and cheapest web service dedicated to solving captchas by human workers from around the world. By solving captchas with us you help people in poorest regions of the world to earn money, which not only cover their basic needs, but also gives them ability to financially help their families, study and avoid jobs where they're simply not happy.

To use the service you need to [register](https://anti-captcha.com/clients/) and topup your balance. Prices start from $0.0005 per image captcha and $0.002 for Recaptcha. That's $0.5 per 1000 for images and $2 for 1000 Recaptchas.

For more technical information and articles visit our [documentation](https://anti-captcha.com/apidoc) page. 

**Examples how to solve:**

- [Image Captcha](#solve-image-captcha-with-go-example)
- [Recaptcha V2](#solve-recaptcha-v2)
- [Recaptcha V3](#solve-recaptcha-v3)
- [hCaptcha](#solve-hcaptcha)
- [FunCaptcha](#solve-funcaptcha)
- [GeeTest](#solve-geetest)
- [Turnstile](#solve-turnstile)
- [Image to coordinates](#image-to-coordinates)
- [AntiGate (custom tasks)](#solve-antigate-custom-tasks)

### Solve image captcha with Go example:
```go
package main

import (
    "fmt"
    "github.com/anti-captcha/anticaptcha-go"
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
        WebsiteURL:    "https://some-website.com/", // Optional to gather stats in the dashboard by this website
    })
    // OR 
    // solution, err := ac.SolveImage("image-encoded-in-base64", anticaptcha.ImageSettings{})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Captcha Solution:", solution)
}
```
&nbsp;

### Solve Recaptcha V2
```go
package main

import (
    "fmt"
    "github.com/anti-captcha/anticaptcha-go"
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
    
    // Solve Recaptcha V2
    solution, err := ac.SolveRecaptchaV2(anticaptcha.RecaptchaV2{
        WebsiteURL:  "https://huev.com/",
        WebsiteKey:  "6Lcyu8UZAAAAACwSh6Xf58WrNXTu0LLu4F85xf20",
        IsInvisible: false, // Set to 'true' if you are solving an invisible Recaptcha V2
        DataSValue:  "",    // Fill this value if you are solving a ReCaptcha V2 with the 'data-s' parameter, typically found at google.com websites
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Recaptcha g-response token:", solution)
}
```
Also with [proxy](https://anti-captcha.com/apidoc/task-types/RecaptchaV2Task):
```go
// Solve Recaptcha V2 with proxy
solution, err := ac.SolveRecaptchaV2(anticaptcha.RecaptchaV2{
    WebsiteURL:  "https://huev.com/",
    WebsiteKey:  "6Lcyu8UZAAAAACwSh6Xf58WrNXTu0LLu4F85xf20",
    IsInvisible: false, // Set to 'true' if you are solving an invisible Recaptcha V2
    DataSValue:  "",    // Fill this value if you are solving a ReCaptcha V2 with the 'data-s' parameter, typically found at google.com websites
    UserAgent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
    Proxy: &anticaptcha.Proxy{
        Type:      "http",
        IPAddress: "1.2.3.4",
        Port:      1234,
        Login:     "login-optional",
        Password:  "pass-optional",
    },
```

&nbsp;

### Solve Recaptcha V3
```go
package main

import (
    "fmt"
    "github.com/anti-captcha/anticaptcha-go"
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
    
    // Solve Recaptcha V3
    solution, err := ac.SolveRecaptchaV3(anticaptcha.RecaptchaV3{
        WebsiteURL: "https://onlyfans.com/",
        WebsiteKey: "6LcvNcwdAAAAAMWAuNRXH74u3QePsEzTm6GEjx0J",
        PageAction: "somefun",
        MinScore:   0.9,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Recaptcha g-response token:", solution)
}
```
&nbsp;

### Solve Hcaptcha
```go
package main

import (
    "fmt"
    "github.com/anti-captcha/anticaptcha-go"
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
    
    // Solve Hcaptcha without proxy
    solution, err := ac.SolveHcaptcha(anticaptcha.Hcaptcha{
        WebsiteURL: "https://www.website.com/",
        WebsiteKey: "00000000-1111-2222-3333-444444444444",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Hcaptcha Token:", solution)
    // Use this user-agent for the form submission
    fmt.Println("User-Agent:", ac.HcaptchaUserAgent)
    // Optional "respkey" value, you may need it too
    fmt.Println("respkey:", ac.HcaptchaRespKey)
}
```
Also with [proxy](https://anti-captcha.com/apidoc/task-types/HCaptchaTask):
```go
// Solve Hcaptcha with proxy
solution, err := ac.SolveHcaptchaProxyOn(anticaptcha.Hcaptcha{
    WebsiteURL: "https://www.website.com/",
    WebsiteKey: "00000000-1111-2222-3333-444444444444",
    Proxy: &anticaptcha.Proxy{
        Type:      "http",
        IPAddress: "1.2.3.4",
        Port:      1234,
        Login:     "login-optional",
        Password:  "pass-optional",
    },
```
&nbsp;
### Solve FunCaptcha
```go
package main

import (
    "fmt"
    "github.com/anti-captcha/anticaptcha-go"
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
    
    // Solve FunCaptcha
    solution, err := ac.SolveFunCaptcha(anticaptcha.FunCaptcha{
        WebsiteURL:       "https://www.website.com/",
        WebsitePublicKey: "00000000-1111-2222-3333-444444444444",
        //make sure to find and set this correctly, look for URL like https://somewebsite-api.arkoselabs.com/v2/00000000-1111-2222-3333-444444444444/api.js
        ApiSubdomain:     "somewebsite-api.arkoselabs.com", 
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Funcaptcha Token:", solution)
}
```
Also with [proxy](https://anti-captcha.com/apidoc/task-types/FunCaptchaTask):
```go
// Solve FunCaptcha with proxy
solution, err := ac.SolveFunCaptchaProxyOn(anticaptcha.FunCaptcha{
    WebsiteURL:       "https://www.website.com/",
    WebsitePublicKey: "00000000-1111-2222-3333-444444444444",
    ApiSubdomain:     "somewebsite-api.arkoselabs.com", 
    Proxy: &anticaptcha.Proxy{
        Type:      "http",
        IPAddress: "1.2.3.4",
        Port:      1234,
        Login:     "login-optional",
        Password:  "pass-optional",
    },
```

&nbsp;
### Solve Turnstile
```go
package main

import (
    "fmt"
    "github.com/anti-captcha/anticaptcha-go"
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
    
    // Solve Turnstile
    solution, err := ac.SolveTurnstile(anticaptcha.Turnstile{
        WebsiteURL: "https://www.website.com/",
        WebsiteKey: "0x4AAAAAAABD2Inoxs-yJ8bz",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Turnstile Token:", solution)
}
```
Also with [proxy](https://anti-captcha.com/apidoc/task-types/TurnstileTask):
```go
// Solve Turnstile with proxy
solution, err := ac.SolveTurnstileProxyOn(anticaptcha.Turnstile{
    WebsiteURL: "https://www.website.com/",
    WebsiteKey: "0x4AAAAAAABD2Inoxs-yJ8bz",
    Proxy: &anticaptcha.Proxy{
        Type:      "http",
        IPAddress: "1.2.3.4",
        Port:      1234,
        Login:     "login-optional",
        Password:  "pass-optional",
    },
```


### Solve GeeTest
GeeTest has 2 versions, number 3 and 4. Number 3 requires parameter "challenge". Number 4 has optional setting "InitParameters".
```go
package main

import (
    "fmt"
    "github.com/anti-captcha/anticaptcha-go"
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
    
    //Solve Geetest
    solution, err := ac.SolveGeeTest(anticaptcha.GeeTest{
        WebsiteURL: "https://bitget.com/",
        Gt:         "e9ca9c9ca19ad540a8017f5c107b2d0f",
        // Solve GeeTest 4:
        Version:    4,
        InitParameters: map[string]interface{}{
            "riskType": "slide",
        },
        
        // Solve GeeTest 3:
        //Version:  3,
        //Challenge:  "1234567890abcdef1234567890abcdef",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Captcha Solution:", solution)
}
```
&nbsp;

Also with [proxy](https://anti-captcha.com/apidoc/task-types/GeeTestTask):
```go
// Solve Turnstile with proxy
solution, err := ac.SolveGeeTestProxyOn(anticaptcha.Turnstile{
    WebsiteURL: "https://bitget.com/",
    Gt:         "e9ca9c9ca19ad540a8017f5c107b2d0f",
    Version:    4,
    Proxy: &anticaptcha.Proxy{
        Type:      "http",
        IPAddress: "1.2.3.4",
        Port:      1234,
        Login:     "login-optional",
        Password:  "pass-optional",
    },
```




### Image to coordinates:

```go
package main

import (
    "encoding/base64"
    "fmt"
    "github.com/anti-captcha/anticaptcha-go"
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
    
    // Solve image-to-coordinates captcha
    imageData, err := ac.ReadImageFile("coordinates.jpg")
    if err != nil {
        log.Fatal(err)
    }
    solution, err := ac.SolveImageToCoordinates(base64.StdEncoding.EncodeToString(imageData), anticaptcha.ImageToCoordinates{
        Comment: "Select object in the specified order",
        Mode:    "points",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Objects X,Y coordinates:", solution)
}
```
&nbsp;

### Solve AntiGate (custom tasks)
```go
package main

import (
    "fmt"
    "github.com/anti-captcha/anticaptcha-go"
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
    
    //Solve AntiGate
    solution, err := ac.SolveAntiGate(anticaptcha.AntiGate{
        WebsiteURL:   "http://antigate.com/logintest.php",
        TemplateName: "Sign-in and wait for control text",
        Variables: map[string]interface{}{
            "login_input_css":      "#login",
            "login_input_value":    "the login",
            "password_input_css":   "#password",
            "password_input_value": "the password",
            "control_text":         "You have been logged successfully",
        },
        Proxy: &anticaptcha.Proxy{
            Type:      "http",
            IPAddress: "1.2.3.4",
            Port:      1234,
            Login:     "login-optional",
            Password:  "pass-optional",
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Captcha Solution:", solution)
    fmt.Println("Cookies:", solution["cookies"])
    fmt.Println("localStorage:", solution["localStorage"])
}
```
&nbsp;