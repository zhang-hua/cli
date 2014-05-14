package application_opener // TODO: rename this package

type URLOpener interface {
	OpenURL(url string) (output string, err error)
}

type urlOpener struct {
	commandProvider CommandProvider
}

func NewURLOpener(provider CommandProvider) URLOpener {
	return urlOpener{commandProvider: provider}
}

func (opener urlOpener) OpenURL(url string) (string, error) {
	cmd := opener.commandProvider.NewCommand("open", url)
	cmd.Run()
	return "", nil
}
