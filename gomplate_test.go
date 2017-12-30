package main

import (
	"bytes"
	"net/http/httptest"
	"os"
	"testing"

	"text/template"

	"github.com/hairyhenderson/gomplate/aws"
	"github.com/hairyhenderson/gomplate/conv"
	"github.com/hairyhenderson/gomplate/data"
	"github.com/hairyhenderson/gomplate/env"
	"github.com/stretchr/testify/assert"
)

func testTemplate(g *Gomplate, template string) string {
	var out bytes.Buffer
	g.RunTemplate(&input{name: "testtemplate", contents: template, target: &out})
	return out.String()
}

func TestGetenvTemplates(t *testing.T) {
	g := &Gomplate{
		funcMap: template.FuncMap{
			"getenv": env.Getenv,
			"bool":   conv.Bool,
		},
	}
	assert.Empty(t, testTemplate(g, `{{getenv "BLAHBLAHBLAH"}}`))
	assert.Equal(t, os.Getenv("USER"), testTemplate(g, `{{getenv "USER"}}`))
	assert.Equal(t, "default value", testTemplate(g, `{{getenv "BLAHBLAHBLAH" "default value"}}`))
}

func TestBoolTemplates(t *testing.T) {
	g := &Gomplate{
		funcMap: template.FuncMap{
			"bool": conv.Bool,
		},
	}
	assert.Equal(t, "true", testTemplate(g, `{{bool "true"}}`))
	assert.Equal(t, "false", testTemplate(g, `{{bool "false"}}`))
	assert.Equal(t, "false", testTemplate(g, `{{bool "foo"}}`))
	assert.Equal(t, "false", testTemplate(g, `{{bool ""}}`))
}

func TestEc2MetaTemplates(t *testing.T) {
	createGomplate := func(status int, body string) (*Gomplate, *httptest.Server) {
		server, ec2meta := aws.MockServer(status, body)
		return &Gomplate{funcMap: template.FuncMap{"ec2meta": ec2meta.Meta}}, server
	}

	g, s := createGomplate(404, "")
	defer s.Close()
	assert.Equal(t, "", testTemplate(g, `{{ec2meta "foo"}}`))
	assert.Equal(t, "default", testTemplate(g, `{{ec2meta "foo" "default"}}`))

	s.Close()
	g, s = createGomplate(200, "i-1234")
	defer s.Close()
	assert.Equal(t, "i-1234", testTemplate(g, `{{ec2meta "instance-id"}}`))
	assert.Equal(t, "i-1234", testTemplate(g, `{{ec2meta "instance-id" "default"}}`))
}

func TestEc2MetaTemplates_WithJSON(t *testing.T) {
	server, ec2meta := aws.MockServer(200, `{"foo":"bar"}`)
	defer server.Close()
	g := &Gomplate{
		funcMap: template.FuncMap{
			"ec2meta":    ec2meta.Meta,
			"ec2dynamic": ec2meta.Dynamic,
			"json":       data.JSON,
		},
	}

	assert.Equal(t, "bar", testTemplate(g, `{{ (ec2meta "obj" | json).foo }}`))
	assert.Equal(t, "bar", testTemplate(g, `{{ (ec2dynamic "obj" | json).foo }}`))
}

func TestJSONArrayTemplates(t *testing.T) {
	g := &Gomplate{
		funcMap: template.FuncMap{
			"jsonArray": data.JSONArray,
		},
	}

	assert.Equal(t, "[foo bar]", testTemplate(g, `{{jsonArray "[\"foo\",\"bar\"]"}}`))
	assert.Equal(t, "bar", testTemplate(g, `{{ index (jsonArray "[\"foo\",\"bar\"]") 1 }}`))
}

func TestYAMLTemplates(t *testing.T) {
	g := &Gomplate{
		funcMap: template.FuncMap{
			"yaml":      data.YAML,
			"yamlArray": data.YAMLArray,
		},
	}

	assert.Equal(t, "bar", testTemplate(g, `{{(yaml "foo: bar").foo}}`))
	assert.Equal(t, "[foo bar]", testTemplate(g, `{{yamlArray "- foo\n- bar\n"}}`))
	assert.Equal(t, "bar", testTemplate(g, `{{ index (yamlArray "[\"foo\",\"bar\"]") 1 }}`))
}

func TestSliceTemplates(t *testing.T) {
	g := &Gomplate{
		funcMap: template.FuncMap{
			"slice": conv.Slice,
		},
	}
	assert.Equal(t, "foo", testTemplate(g, `{{index (slice "foo") 0}}`))
	assert.Equal(t, `[foo bar 42]`, testTemplate(g, `{{slice "foo" "bar" 42}}`))
	assert.Equal(t, `helloworld`, testTemplate(g, `{{range slice "hello" "world"}}{{.}}{{end}}`))
}

func TestHasTemplate(t *testing.T) {
	g := &Gomplate{
		funcMap: template.FuncMap{
			"yaml": data.YAML,
			"has":  conv.Has,
		},
	}
	assert.Equal(t, "true", testTemplate(g, `{{has ("foo:\n  bar: true" | yaml) "foo"}}`))
	assert.Equal(t, "true", testTemplate(g, `{{has ("foo:\n  bar: true" | yaml).foo "bar"}}`))
	assert.Equal(t, "false", testTemplate(g, `{{has ("foo: true" | yaml) "bah"}}`))
	tmpl := `{{- $data := yaml "foo: bar\nbaz: qux\n" }}
{{- if (has $data "baz") }}
{{- $data.baz }}
{{- end }}`
	assert.Equal(t, "qux", testTemplate(g, tmpl))
	tmpl = `{{- $data := yaml "foo: bar\nbaz: qux\n" }}
{{- if (has $data "quux") }}
{{- $data.quux }}
{{- else }}
{{- $data.foo }}
{{- end }}`
	assert.Equal(t, "bar", testTemplate(g, tmpl))
}

func TestCustomDelim(t *testing.T) {
	g := &Gomplate{
		leftDelim:  "[",
		rightDelim: "]",
		funcMap:    template.FuncMap{},
	}
	assert.Equal(t, "hi", testTemplate(g, `[print "hi"]`))
}
