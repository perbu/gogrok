package fragments

import (
	"fmt"
	"github.com/perbu/gogrok/analytics"
	"net/url"
	"strconv"
)

// use templ to generate the Go code:
//go:generate go run github.com/a-h/templ/cmd/templ@latest generate fragments.templ

func s(number any) string {
	switch typedNumber := number.(type) {
	case int:
		return strconv.Itoa(number.(int))
	case int64:
		return strconv.FormatInt(typedNumber, 10)
	default:
		return fmt.Sprintf("%v", number)
	}
}

func slen[T any](slice []T) string {
	return strconv.Itoa(len(slice))
}

func moduleUrl(mod *analytics.Module) string {
	return fmt.Sprintf("/module/%s", mod.Path)
}

func packageUrl(pkg *analytics.Package) string {
	u, err := url.Parse(fmt.Sprintf("/package/%s", pkg.Module.Path))
	if err != nil {
		panic(err)
	}
	q := u.Query()
	q.Set("package", pkg.Name)
	u.RawQuery = q.Encode()
	return u.String()
}

func fileUrl(file *analytics.File) string {
	u, err := url.Parse(fmt.Sprintf("/file/%s", file.Module.Path))
	if err != nil {
		panic(err)
	}
	q := u.Query()
	q.Set("package", file.Package.Name)
	q.Set("file", file.Name)
	u.RawQuery = q.Encode()
	return u.String()
}

// truncateVersion truncates a version string if it's longer than 10 characters
// and adds "..." at the end
func truncateVersion(version string) string {
	if len(version) > 10 {
		return version[:10] + "..."
	}
	return version
}
