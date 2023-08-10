package translator

import (
	"context"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
)

type Translation map[string]string
type Translations map[string]Translation

var (
	unitrans *ut.UniversalTranslator
	transKey = struct{ name string }{name: "translator"}
)

func init() {
	english := en.New() // fallback locale
	unitrans = ut.New(english, english, ru.New())
}

// FromContext gets the translator from the context, or fallback if not found
func FromContext(ctx context.Context) ut.Translator {
	trans, ok := ctx.Value(transKey).(ut.Translator)
	if !ok {
		return unitrans.GetFallback()
	}

	return trans
}

// WithContext appends translator to the existing context
func WithContext(ctx context.Context, trans ut.Translator) context.Context {
	return context.WithValue(ctx, transKey, trans)
}

func Get(locale string) (ut.Translator, bool) { return unitrans.GetTranslator(locale) }

func Find(locales ...string) (ut.Translator, bool) { return unitrans.FindTranslator(locales...) }

func Verify() error { return unitrans.VerifyTranslations() }
