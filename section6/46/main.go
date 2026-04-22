package main

import "fmt"

type UserNotFound struct {
	Username string
}

func (e *UserNotFound) Error() string {
	return fmt.Sprintf("User not found: %s", e.Username)
}

func myFunc() error {
	ok := false
	if ok {
		return nil
	} else {
		return &UserNotFound{Username: "mike"}
}
}

func main() {
	if err := myFunc(); err != nil {
		fmt.Println(err)
	}
}