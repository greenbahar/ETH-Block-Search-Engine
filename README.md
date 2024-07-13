Ethereum Blcok discovery
====
This task is going to implement a service which uses the Ethereum JSON-RPC API to store the following information in a local datastore for the most recent 50 blocks and provide a basic API:


Instructions
-----

1. Get and store the block and all transaction hashes in the block
2. Get and store all events related to each transaction in each block
3. Expose an endpoint that allows a user to query all events related to a particular address Given


Important notes
-----

You can use the following endpoints for Ethereum RPC access. For getting the API KEY, go to Infura's website, then sign up for a free account. After verifing your email, login. Then in the dashboard go for "CREATE NEW API KEY".


1. https://mainnet.infura.io/v3/<YOUR API KEY>
2. wss://mainnet.infura.io/ws/v3/<YOUR API KEY>
3. [go-ethereum](github.com/ethereum/go-ethereum) client for the rpc calls. 



Thoughts
--------

1. Based on the size of data, 50 blocks, map data structure could be good as data storage
2. Providing easy-to-query API, data `indexed` in some maps
3. external package `go-ethereum ` used as ethereum client for rpc calls
4. lock management to prevent deadlock. However, moving forward single thread db communication we can be free on lock management but it costs performance.
5. Race condition prevented, also data modification of db methods' consumers restricted
4. To improve the performance, workerool speeded of blockprocessing nearly by `5 times` (based on benchmark on my pc)
5. gurilla/mux as web framework. Basically the routing is `dicuppled` an can be replaced
6. Modular structur, components talk to each other through interfaces. Hence, well decupled in case of decision to change a module.
7. For example, the inmemory database can easily be replaced by Cassandra (Cassandra can be a good candidate as per simple query patterns and scale freindly) for the entire chain history.



How to run by docker 
----
There need to be no input for the Load functionin the main.go; like `godotenv.Load()` . 

```json5
cd <project path>/cmd/bash
chmod +x build.sh
chmod +x run.sh
./build.sh
./run.sh
```

How to run manually 
----
Please set `godotenv.Load("../../.env")` in the main.go. But if you are running the project by docker you need to leave no input like `godotenv.Load()`. I can provide a flag for this two ways of running manually (debug mode) or by docker.

```json5
go mod download
cd <project path>/cmd/ethereum-tracker-app
go run main.go
```

Sample output
---------------

Generated JSON response of the API call test via postman:

Method: GET

URL: 127.0.0.1:8000/v1/events/0x388C818CA8B9251b393131C08a736A67ccB19297

```json5
[
    {
        "address": "0x388c818ca8b9251b393131c08a736a67ccb19297",
        "topics": [
            "0x27f12abfe35860a9a927b465bb3d4a9c23c8428174b83f278fe45ed7b4da2662"
        ],
        "data": "0x00000000000000000000000000000000000000000000000007778f259c4d3b9f",
        "blockNumber": "0x132eee1",
        "transactionHash": "0xc127cef074e536054c81640cff38756a4072d235af78d1b7ef84d1406e1ba0a6",
        "transactionIndex": "0x71",
        "blockHash": "0x75c05edbe6b4b843ad6b1a438508718c0b2960d73c0b545f99f63fe2c05a1f1f",
        "logIndex": "0x211",
        "removed": false
    },
    {
        "address": "0x388c818ca8b9251b393131c08a736a67ccb19297",
        "topics": [
            "0x27f12abfe35860a9a927b465bb3d4a9c23c8428174b83f278fe45ed7b4da2662"
        ],
        "data": "0x00000000000000000000000000000000000000000000000000d8358d228b2cc8",
        "blockNumber": "0x132eee5",
        "transactionHash": "0x8cbb5e9c5cccfdffc88743d6abe2dc51c16d1df5571f1b2afad036e3c1692c88",
        "transactionIndex": "0x7d",
        "blockHash": "0x8e7fa066a3b6b2e6576c0a20267dcc0bbd1e4e2665d2c003b22252cdb54a386a",
        "logIndex": "0x139",
        "removed": false
    },
    ...
]    
```

adds-on
---
__already added adds-on__:
1. Adding a worker pool for processing the transactions in each block in parallel
2. Standard Go project layout
3. concurrency safe consideration
4. some performance improvements
5. dockerized + bash commands
6. tests
7. Update the datastore by subscribing new headers
8. swager for API doc
9. more systematic error handling
    - using "github.com/pkg/errors" to wrap errors with proper stack tracing, also for better monitoring.
    - From go 1.13 standard "errors" package provides wraping functionality in fmt.Errorf(".... %w", err), but not tracing.
    - Removing unnecesarry errors from goroutines funcs. Not all functions need to return errors, especially when they are designed to run as goroutines or in other contexts where immediate error handling by the caller is impractical or unnecessary. Instead, logging errors can suffice in some situations.
    - Errors within goroutines are logged, ensuring they are not ignored.
    - Granular error codes provided. Hence, having a midleware on routes, after handler function execution, can provide some metrics to provide observability (Nevertheless other observability approaches can be taken).
10. Enhance configuration and remove hardcodes
11. Logging  
12. Graceful Shutdown 

__nice to have adds-on__:
1. Security related middlewares
2. github actions 
3. gather metrics by prometheus, for later monitoring purpose on datadog, grafana
4. variadic function to pass ...options to constructors. passing option functions instead of one-by-one entities
5. better data retreival in rpc calls. for example in case of BlockByNumber faced error, repeat for some times to get the data. This can be done for all rpc calls.
```golang
if number == nil {
    return nil, errors.New("block number cannot be nil")
}

var block *types.Block
var err error

// retry strategy
retryCount := 3
for i := 0; i < retryCount; i++ {
    block, err = ec.httpClient.BlockByNumber(ctx, number)
    if err == nil {
        return block, nil
    }
    
    // Log the error for visibility of which blocks are missed
    log.Printf("Error retrieving block by number %s: %v (retrying %d/%d)", number.String(), err, i+1, retryCount)
    
    // Exponential backoff, sleep before retrying
    time.Sleep(time.Duration(i+1) * time.Second)
    
    select {
    case <-ctx.Done():
        return nil, errors.Wrap(ctx.Err(), "context canceled or timed out")
    default:
        Continue
    }
}
// If we reach here, it means all retries have failed
return nil, errors.Wrapf(err, "failed to retrieve block by number %s after %d retries", number.String(), retryCount)

```
