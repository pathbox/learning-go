package recaptcha

import (
	"html/template"
	"net/http"

	"github.com/haisum/recaptcha"
)

var (
	recap Info
)

type Info struct {
	Enabled bool
	Secret  string
	SiteKey string
}

func Configure(c Info) {
	recap = c
}

func ReadConfig() Info {
	return recap
}

func Verified(r *http.Request) bool {
	if !recap.Enabled {
		return true
	}

	re := recaptcha.R{
		Secret: recap.Secret,
	}
	return re.Verify(*r)
}

func Plugin() template.FuncMap {
	f := make(template.FuncMap)

	f["RECAPCHA_SITEKEY"] = func() template.HTML {
		if ReadConfig().Enabled {
			return template.HTML(ReadConfig().SiteKey)
		}

		return template.HTML("")
	}

	return f
}
