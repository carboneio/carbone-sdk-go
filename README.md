# Carbone Render Go SDK
![Version](https://img.shields.io/badge/version-1.0.0-blue.svg?cacheSeconds=2592000)
[![Documentation](https://img.shields.io/badge/documentation-yes-brightgreen.svg)](https://carbone.io/api-reference.html#carbone-sdk-go)

> The golang SDK to use Carbone Render easily.

### 🏠 [Homepage](https://github.com/Ideolys/carbone-sdk-go)
### 🔖 [Documentation](https://carbone.io/api-reference.html#carbone-sdk-go)

## TODO
- [ ] Ajouter la documentation sur le site carbone.io
- [ ] Ajouter le package à la librairie public go

## Install

```sh
go get github.com/Ideolys/carbone-sdk-go
```

## Usage

```go
package main

import (
	carbone "github.com/github.com/Ideolys/carbone-sdk-go"
)

func main() {
	// ...
	// csdk := carbone.
}
```

## Run tests

To run all the tests (-v for verbose output):
```bash
$ go test -v
```

To run only one test:
```bash
$ go test -v -run NameOfTheTest
```

If you need to test the generation of templateId, you can use the nodejs `main.js` to test the sha256 generation.
```bash
$ node ./tests/main.js
```

## 👤 Author

- [**@steevepay**](https://github.com/steevepay)

## 🤝 Contributing

Contributions, issues and feature requests are welcome!

Feel free to check [issues page](https://github.com/Ideolys/carbone-sdk-go/issues).

## Show your support

Give a ⭐️ if this project helped you!