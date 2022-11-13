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
            description: The device has been successfully registered
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

## Run the application
``` go build```
```./LoraWanAPI```

## Output:

```
API call number : deveui
9 :  99707
7 :  64C28
5 :  AD86C
6 :  5D4C1
4 :  033AA
0 :  C4D08
10 :  E3E75
12 :  346FF
8 :  61185
18 :  4F7B1
17 :  9181F
16 :  D6B87
15 :  C74DD
13 :  3F919
2 :  5BAB1
21 :  CA373
22 :  39908
25 :  014C1
26 :  7F3CE
14 :  5A075
3 :  0D03D
11 :  2E8BC
31 :  C4A2D
19 :  4C869
33 :  E0CB7
23 :  361E1
35 :  5B714
27 :  45414
28 :  8514C
29 :  11225
38 :  BFCEB
20 :  A8D34
40 :  0ADB8
1 :  FBFE8
41 :  72ECB
30 :  EA229
42 :  B9FAD
43 :  152B8
45 :  7D437
39 :  FAEA5
47 :  F0201
32 :  8AFEA
37 :  E5CCC
50 :  44469
49 :  1A669
34 :  E1443
53 :  BEEFA
24 :  2C82A
^C received C-c - shutting down
telling other goroutines to stop
36 :  A5324
58 goroutine ended
59 goroutine ended
60 goroutine ended
61 goroutine ended
62 goroutine ended
63 goroutine ended
64 goroutine ended
65 goroutine ended
66 goroutine ended
67 goroutine ended
68 goroutine ended
69 goroutine ended
70 goroutine ended
71 goroutine ended
72 goroutine ended
73 goroutine ended
74 goroutine ended
75 goroutine ended
76 goroutine ended
77 goroutine ended
78 goroutine ended
79 goroutine ended
80 goroutine ended
81 goroutine ended
82 goroutine ended
83 goroutine ended
84 goroutine ended
85 goroutine ended
86 goroutine ended
87 goroutine ended
88 goroutine ended
89 goroutine ended
90 goroutine ended
91 goroutine ended
92 goroutine ended
93 goroutine ended
94 goroutine ended
95 goroutine ended
96 goroutine ended
97 goroutine ended
98 goroutine ended
99 goroutine ended
46 :  34A7E
48 goroutine ended
51 goroutine ended
52 :  0B372
44 goroutine ended
55 :  31CB1
57 :  2671E
56 :  2C9C3
54 :  E938A
cleanup: all running goroutines are completed.
```