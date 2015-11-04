## ersatz

*Note: this is a work in progress!*

Ersatz is a tool for mocking RESTful JSON API requests written in Go. It uses simple directory and file structures to emulate the requests and responses that might be associated with a live API server. Response endpoints can be 'varied' (that is, issue different response codes, headers and bodies) using a very simple RESTful command API endpoint.

### Installation

Assuming you have the Go toolchain installed, run `go get github.com/homemade/ersatz`. The command line utility `ersatz` will be added to your `$GOPATH/bin` directory. If you've set up Go in the standard way, that will be included in your `$PATH`, you can now run ersatz by typing `ersatz`.

### Basic use

The server can be started with the command `ersatz server <port> </path/to/definitions/directory>`.

#### File structure
Ersatz requires a drectory of endpoint definitions. The expected file layout is simple and consistent: `{definitions/dir}/<{some/api/endpoint}>/{HTTP_VERB}/{variation}.json`. Often you won't want to use an explicit variation, so the tool defauls to `default`. An example directory tree might look like this:

```    
    |-- /path/to/definitions/directory
        |-- products
            |-- GET
                |-- default.json
                |-- no-products.json
            |-- POST
                |-- default.json
                |-- validation-failed.json
            |-- 123
                |-- GET
                    |-- default.json
                    |-- not-found.json
                    |-- permission-denied.json
                |-- PUT
                    |-- default.json
                    |-- not-found.json
                    |-- permission-denied.json
```

Using that definitions directory, ersatz would serve the following, unvaried endpoints:

```
GET /products
POST /products
GET /products/123
PUT /products/123 
```

#### Endpoint file format

The JSON files follow a fairly simple format:

```
{  
    "response_code": 200, # Any valid HTTP response code
    "headers":{  
        "header-name":"header-value"
    },
    "body":{  
        "any":"legal",
        "json":"values",
        "go":"here"
    }
}
```

If the JSON above were in the file `/path/to/definitions/directory/products/GET/default.json`, then ersatz would serve `GET /products` with one additional header, and the value of `body` from the JSON.

### Variation

Ersatz has a special command endpoint, `/__ersatz`. This endpoint currently only accepts POST requests, and can be used to vary the next request made to an endpoint.

For example, the request `POST /__ersatz`, with the raw request body

```
{
    "command": "vary",
    "endpoint": {
        "url": "/products/123",
        "method": "GET",
        "variant": "permission-denied"
    }
}
```

instructs the ersatz to respond to the next request to `GET /products/123` with the JSON file `/path/to/definitions/directory/products/123/GET/permission-denied.json` instead of the default `/path/to/definitions/directory/products/123/GET/default.json`. This is the core of the API mocking operation: allowing clients to control the response they expect.

### Roadmap

This is a very early release of ersatz. Everything is very likely to change and, if you use it, breaks will no doubt occur as updates are rolled out! 









