package application_opener_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestURLOpener(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "URL Opener Suite")
}
