package typeguessing

import (
	"testing"
)

func TestLearner_Int(t *testing.T) {
	l := NewLearner()
	l.Feed("148")
	l.Feed("0")
	l.Feed("-7")
	l.Feed("")
	l.Feed("29")

	guess := l.BestGuess()
	if guess.(int64) != INT {
		t.Errorf("%v != %d", guess, INT)
	}
}

func TestLearner_Bool(t *testing.T) {
	l := NewLearner()
	l.Feed("false")
	l.Feed("True")
	l.Feed("true")
	l.Feed("False")
	l.Feed("")
	l.Feed("false")

	guess := l.BestGuess()
	if guess.(bool) != BOOL {
		t.Errorf("%v != %d", guess, BOOL)
	}
}

func TestLearner_Float(t *testing.T) {
	l := NewLearner()
	//even though most of these parse as int, the
	//more compatible type (float64) should be chosen
	l.Feed("1")
	l.Feed(".1038738947937594728")
	l.Feed("3.487")
	l.Feed("2")
	l.Feed("")
	l.Feed("50")

	guess := l.BestGuess()
	if guess.(float64) != FLOAT {
		t.Errorf("%v != %d", guess, FLOAT)
	}
}

func TestLearner_InconsistentTyping_Guesses_String(t *testing.T) {
	l := NewLearner()
	//even though none of these parse as string, string
	//type should be chosen because of compatibility
	l.Feed("1")
	l.Feed(".1038738947937594728")
	l.Feed("3.487")
	l.Feed("2")
	l.Feed("")
	l.Feed("50")
	l.Feed("true")
	l.Feed("")
	l.Feed(".38473")

	guess := l.BestGuess()
	if guess.(string) != STRING {
		t.Errorf("%v != %d", guess, STRING)
	}
}

func TestGuessString_Int(t *testing.T) {
	v := GuessString("1")
	if v.(int64) != 1 {
		t.Errorf("%v != 1", v)
	}
}

func TestGuessString_Float(t *testing.T) {
	v := GuessString(".49258")
	if v.(float64) != 0.49258 {
		t.Errorf("%v != 0.49258", v)
	}
}

func TestGuessString_Bool(t *testing.T) {
	v := GuessString("True")
	if v.(bool) != true {
		t.Errorf("%v != true", v)
	}
}

func TestGuessString_Nil(t *testing.T) {
	v := GuessString("")
	if v != nil {
		t.Errorf("%v != nil", v)
	}
}

func TestGuessStrings(t *testing.T) {
	vals := GuessStrings([]string{
		"345",
		"Test",
		"123.456",
		"false",
	})
	if vals[0].(int64) != 345 {
		t.Errorf("%v != 345", vals[0])
	}
	if vals[1].(string) != "Test" {
		t.Errorf("%v != 'Test'", vals[1])
	}
	if vals[2].(float64) != 123.456 {
		t.Errorf("%v != 123.456", vals[2])
	}
	if vals[3].(bool) != false {
		t.Errorf("%v != false", vals[3])
	}
}
