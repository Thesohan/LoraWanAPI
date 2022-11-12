# LoraWanAPI

## Backend

### Abstract

Each MachineMax sensor has a unique 16-character (hex) identifier called a DevEUI. As part of the manufacturing process, it is written onto the internal storage of the sensor. The DevEUI is also printed on a label on the side of the sensor alongside a 5-character code (the last 5 characters of the DevEUI). For example, a DevEUI of 78111FFFE452555B would have a short code of 2555B

The sensors communicate with the MachineMax cloud though a LoRaWAN provider and the
LoRaWAN provider uses the DevEUI to identify the sensor. This means we first have to register
the DevEUI with the provider before we can use it. We pay for every device registered with the
LoRaWAN provider, so it is important that we only register DevEUIs that we use (this is not strictly
true, but it makes the challenge more interesting).

When a customer registers a new sensor, they will enter the 5-character short-form code instead
of the full DevEUI, so it is essential that each DevEUI in the batch has a unique 5-char code (for
lookups).

## Challenge
Write an application (CLI) that can be run by the technicians on the production line just before
assembling the sensor units to create a batch of 100 unique DevEUIs and registers them with the
LoRaWAN API. They will note the output and feed it into the production system. The technicians
can sometimes be impatient and may kill the process if it takes too long.

## Requirements
1. The application must return every DevEUI that it registers with the LoRaWAN provider
(e.g., if the application is killed it must wait for in-flight requests to finish otherwise, we
would have registered those DevEUIs but would not be using them)
2. It must handle user interrupts gracefully (SIGINT)
3. It must register exactly 100 DevEUIs (no more) with the provider (to avoid paying for
DevEUIs that we do not use)
4. It should make multiple requests concurrently (but there must never be more than 10 requests in-flight, to avoid throttling)

## LoRaWAN API
The registration operation is an HTTP request to the LoRaWAN provider’s API. The details are below

```
host:
    europe-west1-machinemax-dev-d524.cloudfunctions.net
paths:
    /sensor-onboarding-sample:
post:
    consumes:
        - application/json
    parameters:
        in: body
        required: true
        type: string
        name: deveui
    responses:
        200:
            description: The device has bee successfully registered
        422:
            description: The DevEUI has already been used
```
### Helper functions

```
const allowedChars = "ABCDEF0123456789"

func generateHexString(length int) (string, error) {
    max := big.NewInt(int64(len(allowedChars)))
    b := make([]byte, length)
    for i := range b {
        n, err := rand.Int(rand.Reader, max)
        if err != nil {
            return "", err
        }
        b[i] = allowedChars[n.Int64()]
    }
    return string(b), nil
}
```