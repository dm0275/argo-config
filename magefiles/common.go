package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"strings"
)

func createNamespace(namespace string) (string, error) {
	output, err := run(fmt.Sprintf("kubectl create namespace %s", namespace))
	if err != nil {
		if strings.Contains(fmt.Sprintf("namespaces \"%s\" already exists", namespace), output) {
			return output, nil
		}
	}
	return output, err
}

func run(command string) (string, error) {
	return sh.Output("bash", "-c", command)
}
