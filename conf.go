package missmolly

import "fmt"

// Struct Conf holds the run-time application configuration; the main function
// builds and populates one of these.
//
type Conf struct {
	OhNo string
}

func (c *Conf) Write(s string) {
	fmt.Printf("%+v\n", s)
}

func (c *Conf) RandomCatastrophe() *Other {
	return &Other{1, 2}
}

type Other struct {
	X, Y int
}
