package application_opener_test

import (
	"github.com/cloudfoundry/cli/testhelpers/api/url_opener"

	. "github.com/cloudfoundry/cli/cf/api/application_opener"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("URL Opener", func() {
	It("invokes commands through its command provider", func() {
		fakeProvider := &url_opener.FakeCommandProvider{}
		opener := NewURLOpener(fakeProvider)

		opener.OpenURL("http://example.org")

		Expect(len(fakeProvider.CommandsProvided)).To(Equal(1))
		cmd := fakeProvider.CommandsProvided[0]
		Expect(cmd.Name).To(Equal("open"))
		Expect(cmd.Args).To(Equal([]string{"http://example.org"}))
		Expect(cmd.RunWasCalled).To(BeTrue())
	})
})
