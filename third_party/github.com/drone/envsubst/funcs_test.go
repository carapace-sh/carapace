package envsubst

import "testing"

func Test_len(t *testing.T) {
	got, want := toLen("Hello World"), "11"
	if got != want {
		t.Errorf("Expect len function to return %s, got %s", want, got)
	}
}

func Test_lower(t *testing.T) {
	got, want := toLower("Hello World"), "hello world"
	if got != want {
		t.Errorf("Expect lower function to return %s, got %s", want, got)
	}
}

func Test_lowerFirst(t *testing.T) {
	got, want := toLowerFirst("HELLO WORLD"), "hELLO WORLD"
	if got != want {
		t.Errorf("Expect lowerFirst function to return %s, got %s", want, got)
	}
	defer func() {
		if recover() != nil {
			t.Errorf("Expect empty string does not panic lowerFirst")
		}
	}()
	toLowerFirst("")
}

func Test_upper(t *testing.T) {
	got, want := toUpper("Hello World"), "HELLO WORLD"
	if got != want {
		t.Errorf("Expect upper function to return %s, got %s", want, got)
	}
}

func Test_upperFirst(t *testing.T) {
	got, want := toUpperFirst("hello world"), "Hello world"
	if got != want {
		t.Errorf("Expect upperFirst function to return %s, got %s", want, got)
	}
	defer func() {
		if recover() != nil {
			t.Errorf("Expect empty string does not panic upperFirst")
		}
	}()
	toUpperFirst("")
}

func Test_default(t *testing.T) {
	got, want := toDefault("Hello World", "Hola Mundo"), "Hello World"
	if got != want {
		t.Errorf("Expect default function uses variable value")
	}

	got, want = toDefault("", "Hola Mundo"), "Hola Mundo"
	if got != want {
		t.Errorf("Expect default function uses default value, when variable empty. Got %s, Want %s", got, want)
	}

	got, want = toDefault("", "Hola Mundo", "-Bonjour le monde", "-Halló heimur"), "Hola Mundo-Bonjour le monde-Halló heimur"
	if got != want {
		t.Errorf("Expect default function to use concatenated args when variable empty. Got %s, Want %s", got, want)
	}
}

func Test_substr(t *testing.T) {
	got, want := toSubstr("123456789123456789", "0", "8"), "12345678"
	if got != want {
		t.Errorf("Expect substr function to cut from beginning to length")
	}

	got, want = toSubstr("123456789123456789", "1", "8"), "23456789"
	if got != want {
		t.Errorf("Expect substr function to cut from offset to length")
	}

	got, want = toSubstr("123456789123456789", "9"), "123456789"
	if got != want {
		t.Errorf("Expect substr function to cut beginnging with offset")
	}

	got, want = toSubstr("123456789123456789", "9", "50"), "123456789"
	if got != want {
		t.Errorf("Expect substr function to ignore length if out of bound")
	}

	got, want = toSubstr("123456789123456789", "-3", "2"), "78"
	if got != want {
		t.Errorf("Expect substr function to count negative offsets from the end")
	}

	got, want = toSubstr("123456789123456789", "-300", "3"), "123"
	if got != want {
		t.Errorf("Expect substr function to cut from the beginning to length for negative offsets exceeding string length")
	}

	got, want = toSubstr("12345678", "9", "1"), ""
	if got != want {
		t.Errorf("Expect substr function to cut entire string if pos is itself out of bound")
	}
}
