# Cleanupcheck

## Goal

Detecting HTTP calls solely by the http identifier missed aliased imports and client methods

## Updates

- `net/http` the package ise identified using type information.
- `web.Get()` calls made using import aliases are supported.
- `http.Client.Get()` or `http.Client.Do()` supported.
- `http.Header.Get()` methods that do not make HTTP requests are excluded.


## Example

A warning is now expected for the following usage because the response is not closed:


~~~go
func fetch(client *http.Client) error {
    resp, err := client.Get("https://example.com")
    if err != nil {
        return err
    }

    fmt.Println(resp.StatusCode)
    return nil
}
~~~

## Checks

- [ ] `go builds ./...` success. 
- [ ] HTTP calls made using import aliases are detected.
- [ ] HTTP calls made through an HTTP client are detected.
- [ ] No warning is reported for closed response.
- [ ] `http.Header.Get()` for no warning is reported.


## Out of Scope

This change does not improve the existing Close() detection logic. Early returns and conditional close analysis can be addressed separately.


## License 

cleanupcheck is open-source software licensed under the MIT License.