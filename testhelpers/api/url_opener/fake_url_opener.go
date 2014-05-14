package url_opener

type FakeURLOpener struct {
	OpenURLReceived struct {
		URL string
	}

	OpenURLReturns struct {
		Output string
		Error  error
	}
}

func (fake *FakeURLOpener) OpenURL(url string) (string, error) {
	fake.OpenURLReceived.URL = url
	return fake.OpenURLReturns.Output, fake.OpenURLReturns.Error
}
