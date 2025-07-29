package asset

import "testing"

func TestAssetMapperGet(t *testing.T) {
	a := NewAssetMapper()
	a.CSSAssets["test.css"] = &Asset{
		Path: "test.css",
		Hash: "123",
	}

	result := a.Get("test.css")
	expected := "/test.css?v=123"

	if expected != result {
		t.Errorf("String should be equal. Expected: %s\nGot:%s\n", expected, result)
	}

	expected = "raw.doc"
	result = a.Get(expected)
	if expected != result {
		t.Errorf("String should be equal. Expected: %s\nGot:%s\n", expected, result)
	}
}

func TestAttributeToString(t *testing.T) {
	s := attributeMapToString(map[string]string{
		"data-test": "value",
	})

	expected := "data-test=\"value\""
	if s != expected {
		t.Errorf("String should be equal. Expected: \"%s\"\nGot: \"%s\"\n", expected, s)
	}

	s = attributeMapToString(map[string]string{
		"shouldEscape<>": ">",
	})

	expected = "shouldEscape&lt;&gt;=\"&gt;\""
	if s != expected {
		t.Errorf("String should be equal. Expected: \"%s\"\nGot: \"%s\"\n", expected, s)
	}
}
