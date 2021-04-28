package template

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"gondola/template/assets"

	"github.com/rainycape/vfs"
)

type templateTest struct {
	tmpl   string
	data   interface{}
	result string
}

type testType struct {
}

func (t *testType) Foo() string {
	return "bar"
}

func (t *testType) Bar(s string) string {
	return "bared-" + s
}

type testInt int

func (t *testInt) Next() int {
	return int(*t) + 1
}

var (
	fortyTwo = 42
	ftests   = []*templateTest{
		{"{{ $one := 1 }}{{ $two := 2 }}{{ $three := 3 }}{{ $one }}+{{ $two }}+{{ $three }}={{ add $one $two $three }}", nil, "1+2+3=6"},
		{"{{ add 2 3 }}", nil, "5"},
		{"{{ to_lower .foo }}", map[string]string{"foo": "BAR"}, "bar"},
		{"{{ to_upper .foo }}", map[string]string{"foo": "bar"}, "BAR"},
		{"{{ join .chars .sep }}", map[string]interface{}{"chars": []string{"a", "b", "c"}, "sep": ","}, "a,b,c"},
		{"{{ to_html .s }}", map[string]string{"s": "<foo\nbar"}, "&lt;foo<br>bar"},
		{"{{ mul 2 1.1 }}", nil, "2.2"},
		{"{{ mulf 2 1.1 }}", nil, "2.2"},
		{"{{ muli 2 1.1 }}", nil, "2"},
		{"{{ int (mul 2 1.1) }}", nil, "2"},
		{"{{ sub 1 2 3 }}", nil, "-4"},
		{"{{ subf 1 0.9 -3 }}", nil, "3.1"},
		{"{{ concat \"foo\" \"bar\" }}", nil, "foobar"},
		{"{{ concat (concat \"foo\" \"bar\") \"baz\" }}", nil, "foobarbaz"},
		{"{{ if divisible 5 2 }}1{{ else }}0{{ end }}", nil, "0"},
		{"{{ if divisible 4 2 }}1{{ else }}0{{ end }}", nil, "1"},
		{"{{ indirect . }}", (*int)(nil), "0"},
		{"{{ indirect . }}", fortyTwo, "42"},
		{"{{ indirect . }}", &fortyTwo, "42"},
	}
	compilerTests = []*templateTest{
		{"{{ \"output\" | printf \"%s\" }}", nil, "output"},
		{"{{ call .foo }}", map[string]interface{}{"foo": func() string { return "bar" }}, "bar"},
		{"{{ .Foo }}", struct{ Foo string }{"bar"}, "bar"},
		{"{{ .Foo }}", &testType{}, "bar"},
		{"{{ .Bar \"this\" }}", &testType{}, "bared-this"},
		{"{{ .t.Bar .foo }}", map[string]interface{}{"t": &testType{}, "foo": "foo"}, "bared-foo"},
		{"{{ .t.Bar (concat .foo \"bar\") }}", map[string]interface{}{"t": &testType{}, "foo": "foo"}, "bared-foobar"},
		{"{{ with .A }}{{ . }}{{ else }}no{{ end }}", map[string]string{"A": "yes"}, "yes"},
		{"{{ with .A }}{{ . }}{{ else }}no{{ end }}", nil, "no"},
		{"{{ with .A }}{{ . }}{{ end }}", nil, ""},
		{"{{ range . }}{{ . }}{{ end }}", []int{1, 2, 3}, "123"},
		{"{{ range . }}{{ . }}{{ end }}{{ . }}", []int{1, 2, 3}, "123[1 2 3]"},
		{"{{ range $idx, $el := . }}{{ $idx }}{{ $el }}{{ end }}", []int{1, 2, 3}, "011223"},
		{"{{ range $el := . }}{{ $el }}{{ end }}", []int{1, 2, 3}, "123"},
		{"{{ range $el := . }}{{ . }}{{ end }}", []int{1, 2, 3}, "123"},
		{"{{ range $idx, $el := . }}{{ . }}{{ end }}", []int{1, 2, 3}, "123"},
		{"{{ range . }}{{ else }}nope{{ end }}", nil, "nope"},
		{"{{ range $k, $v := . }}{{ $k }}={{ $v }}{{ end }}", map[string]int{"b": 2, "c": 3, "a": 1}, "a=1b=2c=3"},
		{"{{ range . }}{{ range . }}{{ if even . }}{{ . }}{{ end }}{{ end }}{{ end }}", [][]int{[]int{1, 2, 3, 4, 5, 6}}, "246"},
		{"{{ define \"a\" }}a{{ end }}{{ range . }}{{ template \"a\" . }}{{ end }}", []int{1, 2, 3}, "aaa"},
		{"{{ define \"a\" }}a{{ . }}{{ . }}{{ end }}{{ range . }}{{ template \"a\" . }}{{ end }}", []int{1, 2, 3}, "a11a22a33"},
		{"{{ define \"a\" }}a{{ . }}{{ . }}{{ end }}{{ if . }}{{ template \"a\" . }}{{ end }}", 0, ""},
		{"{{ define \"a\" }}a{{ . }}{{ . }}{{ end }}{{ if . }}{{ template \"a\" . }}{{ end }}", 1, "a11"},
	}
	varTests = []*templateTest{
		{"{{ $Vars.Foo }}", nil, "foo"},
		{"{{ @Foo }}", nil, "foo"},
		{"{{ $Vars.Foo }}{{ concat \"@bar\" }}@baz", nil, "foo@bar@baz"},
		{"{{ @Foo }}{{ concat \"@bar\" }}@baz", nil, "foo@bar@baz"},
		{"{{ $Vars.Foo }}{{ concat \"@bar\" @Foo }}@baz", nil, "foo@barfoo@baz"},
		{"{{ @Foo }}{{ concat \"@bar\" @Foo }}@baz", nil, "foo@barfoo@baz"},
	}
	compilerErrorTests = []*templateTest{
		//{"{{ #invalid }}", nil, "template.html:1:9: can't range over int"},
		{"{{ range . }}{{ else }}nope{{ end }}", 5, "template.html:1:9: can't range over int"},
		{"{{ . }}\n{{ range . }}{{ else }}nope{{ end }}", 5, "template.html:2:9: can't range over int"},
		{"{{ . }}\n{{ range .foo }}{{ else }}nope{{ end }}\n{{ range .bar }}{{ . }}{{ end }} ", map[string]interface{}{"foo": []int{}, "bar": ""}, "template.html:3:9: can't range over string"},
		{"{{ define \"foo\" }}\n{{ range . }}{{ else }}nope{{ end }}{{ end }}\n{{ template \"foo\" . }}", 5, "template.html:2:9: can't range over int"},
		{"{{ .Foo }}", testType{}, "template.html:1:3: method \"Foo\" requires pointer receiver (*template.testType)"},
		{"{{ .Next }}", testInt(0), "template.html:1:3: method \"Next\" requires pointer receiver (*template.testInt)"},
	}
)

func parseNamedText(tb testing.TB, name string, text string, funcs map[string]interface{}, contentType string,
	beforeCompile func(*Template)) *Template {
	fs, err := vfs.Map(map[string]*vfs.File{name: &vfs.File{Data: []byte(text)}})
	if err != nil {
		tb.Fatal(err)
	}
	tmpl := New(fs, nil)
	tmpl.RawFuncs(funcs)
	tmpl.RawFuncs(map[string]interface{}{"t": func(s string) string { return s }})
	tmpl.contentType = contentType
	err = tmpl.Parse(name)
	if err != nil {
		tb.Errorf("error parsing %q: %s", text, err)
		return nil
	}
	if beforeCompile != nil {
		beforeCompile(tmpl)
	}
	if err := tmpl.Compile(); err != nil {
		tb.Errorf("error compiling %q: %s", text, err)
		return nil
	}
	return tmpl
}

func parseText(tb testing.TB, text string) *Template {
	return parseNamedText(tb, "template.html", text, nil, "", nil)
}

func parseTestTemplate(tb testing.TB, name string) *Template {
	fs, err := vfs.FS("_testdata")
	if err != nil {
		tb.Fatal(err)
	}
	tmpl := New(fs, assets.New(fs, ""))
	tmpl.RawFuncs(map[string]interface{}{"t": func(s string) string { return s }})
	if err := tmpl.Parse(name); err != nil {
		tb.Errorf("error parsing %q: %s", name, err)
		return nil
	}
	if err := tmpl.Compile(); err != nil {
		tb.Errorf("error compiling %q: %s", name, err)
		return nil
	}
	return tmpl
}

func testCompiler(t *testing.T, tests []*templateTest, vars VarMap) {
	for _, v := range tests {
		tmpl := parseText(t, v.tmpl)
		if tmpl == nil {
			continue
		}
		var buf bytes.Buffer
		if err := tmpl.ExecuteContext(&buf, v.data, nil, vars); err != nil {
			t.Errorf("error executing %q: %s", v.tmpl, err)
			continue
		}
		if buf.String() != v.result {
			t.Errorf("expecting %q executing %q, got %q", v.result, v.tmpl, buf.String())
		}
		// grab the state from the pool and check its stack
		state := getState()
		if state != nil && len(state.stack) > 0 {
			t.Errorf("template %q left %d elements on stack: %s", v.tmpl, len(state.stack), state.stack)
		}
	}
}

func TestFunctions(t *testing.T) {
	testCompiler(t, ftests, nil)
}

func TestCompiler(t *testing.T) {
	testCompiler(t, compilerTests, nil)
}

func TestVariables(t *testing.T) {
	testCompiler(t, varTests, VarMap{"Foo": "foo"})
}

func TestCompilerErrors(t *testing.T) {
	for _, v := range compilerErrorTests {
		tmpl := parseText(t, v.tmpl)
		if tmpl == nil {
			continue
		}
		var buf bytes.Buffer
		err := tmpl.Execute(&buf, v.data)
		if err == nil {
			t.Errorf("expecting an error when executing %q, got nil", v.tmpl)
			continue
		}
		if err.Error() != v.result {
			t.Logf("template is %q", v.tmpl)
			t.Errorf("expecting error %q, got %q", v.result, err.Error())
		}
	}
}

func TestBigTemplate(t *testing.T) {
	const name = "1.html"
	tmpl := parseTestTemplate(t, name)
	if tmpl != nil {
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, nil); err != nil {
			t.Errorf("error executing template %s: %s", name, err)
		}
	}
}

func benchmarkTests() []*templateTest {
	var tests []*templateTest
	tests = append(tests, ftests...)
	tests = append(tests, compilerTests...)
	return tests
}

func benchmarkTemplate(b *testing.B, tests []*templateTest) {
	b.ReportAllocs()
	templates := make([]*Template, len(tests))
	for ii, v := range tests {
		tmpl := parseText(b, v.tmpl)
		if tmpl == nil {
			b.Fatalf("can't parse %q", v.tmpl)
		}
		templates[ii] = tmpl
	}
	var buf bytes.Buffer
	b.ResetTimer()
	for ii := 0; ii < b.N; ii++ {
		for ii, v := range templates {
			v.Execute(&buf, tests[ii].data)
		}
		buf.Reset()
	}
}

func benchmarkHTMLTemplate(b *testing.B, tests []*templateTest) {
	b.ReportAllocs()
	templates := make([]*template.Template, len(tests))
	for ii, v := range tests {
		tmpl := template.New("template.html")
		tmpl.Funcs(template.FuncMap(templateFuncs.asTemplateFuncMap()))
		_, err := tmpl.Parse(v.tmpl)
		if err != nil {
			b.Fatalf("can't parse %q: %s", v.tmpl, err)
		}
		// Execute once to add the escaping hooks
		tmpl.Execute(ioutil.Discard, nil)
		templates[ii] = tmpl
	}
	var buf bytes.Buffer
	b.ResetTimer()
	for ii := 0; ii < b.N; ii++ {
		for ii, v := range templates {
			v.Execute(&buf, tests[ii].data)
		}
		buf.Reset()
	}
}

func BenchmarkTemplate(b *testing.B) {
	benchmarkTemplate(b, benchmarkTests())
}

func BenchmarkTemplateGo(b *testing.B) {
	benchmarkHTMLTemplate(b, benchmarkTests())
}

func BenchmarkBig(b *testing.B) {
	b.ReportAllocs()
	const name = "1.html"
	tmpl := parseTestTemplate(b, name)
	if tmpl == nil {
		return
	}
	var buf bytes.Buffer
	tmpl.Execute(&buf, nil)
	buf.Reset()
	b.ResetTimer()
	for ii := 0; ii < b.N; ii++ {
		tmpl.Execute(&buf, nil)
		buf.Reset()
	}
}

func BenchmarkBigGo(b *testing.B) {
	b.ReportAllocs()
	tmpl := template.New("")
	tmpl.Funcs(template.FuncMap{"t": func(s string) string { return s }})
	tmpl.Funcs(template.FuncMap(templateFuncs.asTemplateFuncMap()))
	readFile := func(name string) string {
		data, err := ioutil.ReadFile(filepath.Join("_testdata", name))
		if err != nil {
			b.Fatal(err)
		}
		return "{{ $Vars := .Vars }}\n" + string(data)
	}
	if _, err := tmpl.Parse(readFile("1.html")); err != nil {
		b.Fatal(err)
	}
	t2 := tmpl.New("2.html")
	if _, err := t2.Parse(readFile("2.html")); err != nil {
		b.Fatal(err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		b.Fatal(err)
	}
	buf.Reset()
	b.ResetTimer()
	for ii := 0; ii < b.N; ii++ {
		tmpl.Execute(&buf, nil)
		buf.Reset()
	}
}

func rangeTests() []*templateTest {
	var tests []*templateTest
	for _, v := range benchmarkTests() {
		if strings.Contains(v.tmpl, "range") {
			tests = append(tests, v)
		}
	}
	return tests
}

func TestContextFunc(t *testing.T) {
	tmpl := parseNamedText(t, "context", "{{ f_interface }} - {{ f_int }}", map[string]interface{}{
		"f_interface": &Func{Fn: func(ctx interface{}) int { return ctx.(int) }, Traits: FuncTraitContext},
		"f_int":       &Func{Fn: func(ctx int) int { return ctx }, Traits: FuncTraitContext},
	}, "text/plain", nil)

	ctx := 42
	expected := fmt.Sprintf("%d - %d", ctx, ctx)
	var buf bytes.Buffer
	if err := tmpl.ExecuteContext(&buf, nil, ctx, nil); err != nil {
		t.Fatal(err)
	}
	if buf.String() != expected {
		t.Errorf("expecting %q, got %q instead", expected, buf.String())
	}
}

func TestBadContextFunc(t *testing.T) {
	tmpl := parseNamedText(t, "context", "{{ f_float64 }}", map[string]interface{}{
		"f_float64": &Func{Fn: func(ctx float64) float64 { return ctx }, Traits: FuncTraitContext},
	}, "text/plain", nil)
	err := tmpl.ExecuteContext(ioutil.Discard, nil, int(42), nil)
	if err == nil {
		t.Error("expecting an error when executing bad context function")
	} else {
		expected := `context:1:3: context function "f_float64" requires a context of type float64, not int`
		if err.Error() != expected {
			t.Errorf("expected error %q, got %q instead", expected, err.Error())
		}
	}
}

func TestTemplateNoPipeMissedInstructions(t *testing.T) {
	tmpl := parseNamedText(t, "template-nopipe", "{{ define \"foo\" }}2{{ end }}1 {{ template \"foo\" }}", nil, "text/plain", nil)
	expected := "1 2"
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		t.Fatal(err)
	}
	if buf.String() != expected {
		t.Errorf("expecting %q, got %q instead", expected, buf.String())
	}
}

func BenchmarkRange(b *testing.B) {
	benchmarkTemplate(b, rangeTests())
}

func BenchmarkRangeGo(b *testing.B) {
	benchmarkHTMLTemplate(b, rangeTests())
}
