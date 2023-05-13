package meteologix

import "testing"

func TestNew(t *testing.T) {
	c := New()
	if c == nil {
		t.Errorf("New failed, expected Client, got nil")
		return
	}
	if c.config == nil {
		t.Errorf("New failed, expected Config, got nil")
	}
	if c.httpClient == nil {
		t.Errorf("New failed, expected HTTPClient, got nil")
	}
}

func TestNew_WithAcceptLanguage(t *testing.T) {
	e := "de"
	c := New(WithAcceptLanguage(e))
	if c == nil {
		t.Errorf("NewWithAcceptLanguage failed, expected Client, got nil")
		return
	}
	if c.config.acceptLang != e {
		t.Errorf("NewWithAcceptLanguage failed, expected lang value: %s, got: %s", e,
			c.config.acceptLang)
	}
	c = New(WithAcceptLanguage(""))
	if c == nil {
		t.Errorf("NewWithAcceptLanguage failed, expected Client, got nil")
		return
	}
	if c.config.acceptLang != DefaultAcceptLang {
		t.Errorf("NewWithAcceptLanguage failed, expected lang value: %s, got: %s", DefaultAcceptLang,
			c.config.acceptLang)
	}
}

func TestNew_WithAPIKey(t *testing.T) {
	e := "API-KEY"
	c := New(WithAPIKey(e))
	if c == nil {
		t.Errorf("NewWithAPIKey failed, expected Client, got nil")
		return
	}
	if c.config.apiKey != e {
		t.Errorf("NewWithAPIKey failed, expected apiKey value: %s, got: %s", e,
			c.config.apiKey)
	}
	c = New(WithAPIKey(""))
	if c == nil {
		t.Errorf("NewWithAPIKey failed, expected Client, got nil")
		return
	}
	if c.config.apiKey != "" {
		t.Errorf("NewWithAPIKey failed, expected empty apiKey, got: %s", c.config.apiKey)
	}
}

func TestNew_WithUsername(t *testing.T) {
	e := "username"
	c := New(WithUsername(e))
	if c == nil {
		t.Errorf("NewWithUsername failed, expected Client, got nil")
		return
	}
	if c.config.authUser != e {
		t.Errorf("NewWithUsername failed, expected username value: %s, got: %s", e,
			c.config.authUser)
	}
	c = New(WithUsername(""))
	if c == nil {
		t.Errorf("NewWithUsername failed, expected Client, got nil")
		return
	}
	if c.config.authUser != "" {
		t.Errorf("NewWithUsername failed, expected empty username value, got: %s", c.config.authUser)
	}
}

func TestNew_WithPassword(t *testing.T) {
	e := "password"
	c := New(WithPassword(e))
	if c == nil {
		t.Errorf("NewWithPassword failed, expected Client, got nil")
		return
	}
	if c.config.authPass != e {
		t.Errorf("NewWithPassword failed, expected password value: %s, got: %s", e,
			c.config.authPass)
	}
	c = New(WithPassword(""))
	if c == nil {
		t.Errorf("NewWithPassword failed, expected Client, got nil")
		return
	}
	if c.config.authPass != "" {
		t.Errorf("NewWithPassword failed, expected empty password value, got: %s", c.config.authPass)
	}
}

func TestNew_WithUserAgent(t *testing.T) {
	e := "User-Agent"
	c := New(WithUserAgent(e))
	if c == nil {
		t.Errorf("NewWithUserAgent failed, expected Client, got nil")
		return
	}
	if c.config.userAgent != e {
		t.Errorf("NewWithUserAgent failed, expected userAgent value: %s, got: %s", e,
			c.config.userAgent)
	}
	c = New(WithUserAgent(""))
	if c == nil {
		t.Errorf("NewWithUserAgent failed, expected Client, got nil")
		return
	}
	if c.config.userAgent != DefaultUserAgent {
		t.Errorf("NewWithUserAgent failed, expected userAgent value: %s, got: %s", DefaultUserAgent,
			c.config.userAgent)
	}
}

func TestNew_WithNil(t *testing.T) {
	c := New(nil)
	if c == nil {
		t.Errorf("NewWithUserAgent failed, expected Client, got nil")
		return
	}
	if c.config.acceptLang != DefaultAcceptLang {
		t.Errorf("NewWithNil failed, expected lang value: %s, got: %s", DefaultUserAgent,
			c.config.userAgent)
	}
	if c.config.userAgent != DefaultUserAgent {
		t.Errorf("NewWithNil failed, expected userAgent value: %s, got: %s", DefaultUserAgent,
			c.config.userAgent)
	}
}
