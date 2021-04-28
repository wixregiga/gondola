package form

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
	"strconv"

	"gondola/app"
	"gondola/crypto/password"
	"gondola/form/input"
	"gondola/html"
	"gondola/i18n"
	"gondola/net/mail"
	"gondola/util/stringutil"
	"gondola/util/structs"
	"gondola/util/types"
)

//go:generate gondola compile-messages

var (
	formTags = []string{"form", "gondola"}
)

type attrMap map[string]html.Attrs

type Form struct {
	ctx       *app.Context
	id        string
	renderer  Renderer
	values    []reflect.Value
	structs   []*structs.Struct
	fields    []*Field
	attrs     attrMap
	options   *Options
	validated bool
	// Don't include the field name in the error
	NamelessErrors bool
	// DisableCSRF disables CSRF protection when set
	// to true. Note that changing this after the
	// form has been rendered or validated has no effect.
	DisableCSRF bool
	hasCSRF     bool
}

func (f *Form) validate() {
	if err := f.addCSRF(); err != nil {
		panic(err)
	}
	for _, v := range f.fields {
		inp := f.ctx.FormValue(v.HTMLName)
		label := v.Label.TranslatedString(f.ctx)
		if f.NamelessErrors {
			label = ""
		}
		if v.Type.HasChoices() {
			if inp == NotChosen {
				v.err = i18n.Errorfc("form", "You must choose a value").Err(f.ctx)
				continue
			}
			// Verify that the input mathces one of the available choices
			choices := f.fieldChoices(v)
			found := false
			for _, c := range choices {
				if inp == toHTMLValue(c.Value) {
					found = true
					break
				}
			}
			if !found {
				v.err = i18n.Errorfc("form", "%v is not a valid choice", inp).Err(f.ctx)
				continue
			}
		}
		if v.Type == FILE {
			file, header, err := f.ctx.R.FormFile(v.HTMLName)
			if err != nil && !v.Tag().Optional() {
				v.err = input.RequiredInputError(label)
				continue
			}
			if file != nil && header != nil {
				value := &formFile{
					file:   file,
					header: header,
				}
				v.value.Set(reflect.ValueOf(value))
			}
		} else {
			if err := input.InputNamed(label, inp, v.SettableValue(), v.Tag(), true); err != nil {
				v.err = i18n.TranslatedError(err, f.ctx)
				continue
			}
			if v.Type == EMAIL && !v.Tag().Has("novalidate") {
				// Don't validate empty emails. If we reached this point
				// with an empty one, the field is optional.
				if email, ok := v.Value().(string); ok && email != "" {
					if _, err := mail.Validate(email, true); err != nil {
						v.err = i18n.Errorfc("form", "%q is not a valid email address", email).Err(f.ctx)
						continue
					}
				}
			}
		}
		if err := structs.Validate(v.sval.Addr().Interface(), v.Name, f.ctx); err != nil {
			v.err = i18n.TranslatedError(err, f.ctx)
			continue
		}
	}
}

func (f *Form) makeField(name string) (*Field, error) {
	var s *structs.Struct
	idx := -1
	var fieldValue reflect.Value
	var sval reflect.Value
	for ii, v := range f.structs {
		pos, ok := v.QNameMap[name]
		if ok {
			if s != nil {
				return nil, fmt.Errorf("duplicate field %q (found in %v and %v)", name, s.Type, v.Type)
			}
			s = v
			idx = pos
			sval = f.values[ii]
			fieldValue = fieldByIndex(sval, s.Indexes[pos])
			// Check the validation function, so if the function is not valid
			// the error is generated at form instantiation.
			if _, err := structs.ValidationFunction(sval, name); err != nil {
				return nil, err
			}
		}
	}
	if idx < 0 {
		return nil, fmt.Errorf("can't map form field %q", name)
	}
	tag := s.Tags[idx]
	label := tag.Value("label")
	if label == "" {
		label = stringutil.CamelCaseToWords(name, " ")
	}
	var typ Type
	if tag.Has("hidden") {
		typ = HIDDEN
	} else if tag.Has("radio") {
		typ = RADIO
	} else if tag.Has("select") {
		typ = SELECT
	} else {
		switch s.Types[idx].Kind() {
		case reflect.Func:
			return nil, nil
		case reflect.String:
			if s.Types[idx] == reflect.TypeOf(password.Password("")) || tag.Has("password") {
				typ = PASSWORD
			} else if tag.Has("email") {
				typ = EMAIL
			} else {
				if ml, ok := tag.MaxLength(); ok && ml > 0 {
					typ = TEXT
				} else if tag.Has("singleline") || tag.Has("line") {
					typ = TEXT
				} else if tag.Has("password") {
					typ = PASSWORD
				} else {
					typ = TEXTAREA
				}
			}
		case reflect.Bool:
			typ = CHECKBOX
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			typ = TEXT
		default:
			if s.Types[idx] == fileType {
				typ = FILE
				break
			}
			return nil, fmt.Errorf("field %q has invalid type %v", name, s.Types[idx])
		}
	}
	// Check if the struct implements the ChoicesProvider interface
	if typ == RADIO || typ == SELECT {
		container := sval.Addr().Interface()
		if _, ok := container.(ChoicesProvider); !ok {
			return nil, fmt.Errorf("field %q requires choices, but %T does not implement ChoicesProvider", name, container)
		}
	}
	htmlName := f.toHTMLName(name)
	field := &Field{
		Type:        typ,
		Name:        name,
		HTMLName:    htmlName,
		Label:       i18n.String(label),
		Placeholder: i18n.String(tag.Value("placeholder")),
		Help:        i18n.String(tag.Value("help")),
		id:          htmlName,
		value:       fieldValue,
		s:           s,
		sval:        sval,
		pos:         idx,
	}
	return field, nil
}

func (f *Form) lookupField(name string) (*Field, error) {
	for _, v := range f.fields {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, fmt.Errorf("form has no field named %q", name)
}

func (f *Form) makeFields(names []string) error {
	for _, v := range names {
		field, err := f.makeField(v)
		if err != nil {
			return err
		}
		if field != nil {
			f.fields = append(f.fields, field)
		}
	}
	return nil
}

type formStructConfigurator struct{}

func (formStructConfigurator) DecomposeField(s *structs.Struct, typ reflect.Type, tag *structs.Tag) bool {
	if typ == fileType {
		return false
	}
	return true
}

func (f *Form) appendVal(val interface{}) error {
	v, err := types.SettableValue(val)
	if err != nil {
		return err
	}
	s, err := structs.New(val, formTags, formStructConfigurator{})
	if err != nil {
		return err
	}
	f.values = append(f.values, v)
	f.structs = append(f.structs, s)
	return nil
}

func (f *Form) valid() bool {
	for _, f := range f.fields {
		if f.err != nil {
			return false
		}
	}
	return true
}

func (f *Form) HasErrors() bool {
	return f.Submitted() && !f.IsValid()
}

func (f *Form) Submitted() bool {
	return f.ctx.R.Method == "POST" || f.ctx.FormValue("submitted") != ""
}

func (f *Form) IsValid() bool {
	if !f.validated {
		f.validate()
		f.validated = true
	}
	return f.valid()
}

func (f *Form) Fields() []*Field {
	return f.fields
}

func (f *Form) FieldNames() []string {
	names := make([]string, len(f.fields))
	for ii, v := range f.fields {
		names[ii] = v.Name
	}
	return names
}

func (f *Form) FieldByName(name string) (*Field, error) {
	return f.lookupField(name)
}

func (f *Form) Renderer() Renderer {
	return f.renderer
}

func (f *Form) writeTag(buf *bytes.Buffer, tag string, attrs html.Attrs, closed bool) {
	buf.WriteByte('<')
	if closed {
		buf.WriteByte('/')
		buf.WriteString(tag)
	} else {
		buf.WriteString(tag)
		if attrs != nil {
			attrs.WriteTo(buf)
		}
	}
	buf.WriteByte('>')
}

func (f *Form) openTag(buf *bytes.Buffer, tag string, attrs html.Attrs) {
	f.writeTag(buf, tag, attrs, false)
}

func (f *Form) closeTag(buf *bytes.Buffer, tag string) {
	f.writeTag(buf, tag, nil, true)
}

func (f *Form) prepareFieldAttributes(field *Field, attrs html.Attrs, pos int) error {
	if f.renderer != nil {
		fattrs, err := f.renderer.FieldAttributes(field, pos)
		if err != nil {
			return err
		}
		for k, v := range fattrs {
			attrs[k] = v
		}
	}
	return nil
}

func (f *Form) fieldChoices(field *Field) []*Choice {
	// The type was asserted on form creation
	provider := field.sval.Addr().Interface().(ChoicesProvider)
	return provider.FieldChoices(f.ctx, field)
}

func (f *Form) beginInput(buf *bytes.Buffer, field *Field, pos int) error {
	if r := f.renderer; r != nil {
		placeholder := field.Placeholder.TranslatedString(f.ctx)
		if err := r.BeginInput(buf, field, placeholder, pos); err != nil {
			return err
		}
		for _, a := range field.addons {
			if a.Position == AddOnPositionBefore {
				err := r.WriteAddOn(buf, field, a)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (f *Form) endInput(buf *bytes.Buffer, field *Field, pos int) error {
	if r := f.renderer; r != nil {
		for _, a := range field.addons {
			if a.Position == AddOnPositionAfter {
				if err := r.WriteAddOn(buf, field, a); err != nil {
					return err
				}
			}
		}
		if err := r.EndInput(buf, field, pos); err != nil {
			return err
		}
	}
	return nil
}

func (f *Form) writeField(buf *bytes.Buffer, field *Field) error {
	var closed bool
	if field.Type != HIDDEN {
		closed = field.Type != CHECKBOX
		label := field.Label.TranslatedString(f.ctx)
		if err := f.writeLabel(buf, field, field.Id(), label, closed, -1); err != nil {
			return err
		}
	}
	var err error
	switch field.Type {
	case TEXT:
		err = f.writeInput(buf, "text", field)
	case PASSWORD:
		err = f.writeInput(buf, "password", field)
	case EMAIL:
		err = f.writeInput(buf, "email", field)
	case HIDDEN:
		err = f.writeInput(buf, "hidden", field)
	case FILE:
		err = f.writeInput(buf, "file", field)
	case TEXTAREA:
		attrs := html.Attrs{
			"id":   field.Id(),
			"name": field.HTMLName,
		}
		if _, ok := field.Tag().IntValue("rows"); ok {
			attrs["rows"] = field.Tag().Value("rows")
		}
		if err := f.prepareFieldAttributes(field, attrs, -1); err != nil {
			return err
		}
		f.openTag(buf, "textarea", attrs)
		buf.WriteString(toHTMLValue(field.Value()))
		f.closeTag(buf, "textarea")
	case CHECKBOX:
		err = f.writeInput(buf, "checkbox", field)
	case RADIO:
		for ii, v := range f.fieldChoices(field) {
			var value interface{}
			id := fmt.Sprintf("%s_%d", field.Id(), ii)
			if err := f.writeLabel(buf, field, id, v.TranslatedName(f.ctx), false, ii); err != nil {
				return err
			}
			if err := f.beginInput(buf, field, ii); err != nil {
				return err
			}
			attrs := html.Attrs{
				"id":   id,
				"name": field.HTMLName,
				"type": "radio",
			}
			if v.Value != nil {
				attrs["value"] = toHTMLValue(v.Value)
				value = v.Value
			} else {
				value = v.Name
			}
			if reflect.DeepEqual(value, field.Value()) {
				attrs["checked"] = "checked"
			}
			if err := f.prepareFieldAttributes(field, attrs, ii); err != nil {
				return err
			}
			f.openTag(buf, "input", attrs)
			if err := f.endLabel(buf, field, v.TranslatedName(f.ctx), ii); err != nil {
				return err
			}
			if err := f.endInput(buf, field, ii); err != nil {
				return err
			}
		}
	case SELECT:
		attrs := html.Attrs{
			"id":   field.Id(),
			"name": field.HTMLName,
		}
		if field.Tag().Has("multiple") {
			attrs["multiple"] = "multiple"
		}
		if err := f.prepareFieldAttributes(field, attrs, -1); err != nil {
			return err
		}
		f.openTag(buf, "select", attrs)
		for ii, v := range f.fieldChoices(field) {
			var value interface{}
			oattrs := html.Attrs{}
			if v.Value != nil {
				oattrs["value"] = toHTMLValue(v.Value)
				value = v.Value
			} else {
				value = v.Name
			}
			if reflect.DeepEqual(value, field.Value()) {
				oattrs["selected"] = "selected"
			}
			if err := f.prepareFieldAttributes(field, attrs, ii); err != nil {
				return err
			}
			f.openTag(buf, "option", oattrs)
			buf.WriteString(html.Escape(v.TranslatedName(f.ctx)))
			f.closeTag(buf, "option")
		}
		f.closeTag(buf, "select")
	}
	return err
}

func (f *Form) writeLabel(buf *bytes.Buffer, field *Field, id, label string, closed bool, pos int) error {
	attrs := html.Attrs{}
	if r := f.renderer; r != nil {
		label := field.Label.TranslatedString(f.ctx)
		if err := r.BeginLabel(buf, field, label, pos); err != nil {
			return err
		}
		lattrs, err := r.LabelAttributes(field, pos)
		if err != nil {
			return err
		}
		for k, v := range lattrs {
			attrs[k] = v
		}
	}
	if id != "" {
		attrs["for"] = id
	}
	f.openTag(buf, "label", attrs)
	if closed {
		return f.endLabel(buf, field, label, pos)
	}
	return nil
}

func (f *Form) endLabel(buf *bytes.Buffer, field *Field, label string, pos int) error {
	buf.WriteString(html.Escape(label))
	f.closeTag(buf, "label")
	if f.renderer != nil {
		if err := f.renderer.EndLabel(buf, field, pos); err != nil {
			return err
		}
	}
	return nil
}

func (f *Form) writeInput(buf *bytes.Buffer, itype string, field *Field) error {
	if err := f.beginInput(buf, field, -1); err != nil {
		return err
	}
	attrs := html.Attrs{
		"id":   field.Id(),
		"type": itype,
		"name": field.HTMLName,
	}
	if err := f.prepareFieldAttributes(field, attrs, -1); err != nil {
		return err
	}
	switch field.Type {
	case CHECKBOX:
		if t, ok := types.IsTrue(field.value.Interface()); t && ok {
			attrs["checked"] = "checked"
		}
	case TEXT, PASSWORD, EMAIL, HIDDEN:
		attrs["value"] = html.Escape(types.ToString(field.Value()))
		if field.Placeholder != "" {
			attrs["placeholder"] = html.Escape(field.Placeholder.TranslatedString(f.ctx))
		}
		if ml, ok := field.Tag().MaxLength(); ok {
			attrs["maxlength"] = strconv.Itoa(ml)
		}
	case FILE:
	default:
		panic("unreachable")
	}
	f.openTag(buf, "input", attrs)
	if field.Type == CHECKBOX {
		// Close the label before calling EndInput
		if err := f.endLabel(buf, field, field.Label.TranslatedString(f.ctx), -1); err != nil {
			return err
		}
	}
	if err := f.endInput(buf, field, -1); err != nil {
		return err
	}
	return nil
}

func (f *Form) renderField(buf *bytes.Buffer, field *Field) (err error) {
	if provider, ok := field.sval.Addr().Interface().(AddOnProvider); ok {
		field.addons = provider.FieldAddOns(f.ctx, field)
	}
	r := f.renderer
	if r != nil {
		err = r.BeginField(buf, field)
		if err != nil {
			return
		}
	}
	err = f.writeField(buf, field)
	if err != nil {
		return
	}
	if r != nil {
		if ferr := field.Err(); ferr != nil {
			ferr = i18n.TranslatedError(ferr, f.ctx)
			err = r.WriteError(buf, field, ferr)
			if err != nil {
				return
			}
		}
		if field.Help != "" {
			err = r.WriteHelp(buf, field, field.Help.TranslatedString(f.ctx))
			if err != nil {
				return
			}
		}
		err = r.EndField(buf, field)
		if err != nil {
			return
		}
	}
	return
}

func (f *Form) toHTMLName(name string) string {
	return stringutil.CamelCaseToLower(name, "_")
}

func (f *Form) render(fields []*Field) (template.HTML, error) {
	if err := f.addCSRF(); err != nil {
		return template.HTML(""), err
	}
	var buf bytes.Buffer
	var err error
	for _, v := range fields {
		if err = f.renderField(&buf, v); err != nil {
			break
		}
	}
	return template.HTML(buf.String()), err
}

// EncType returns the recommended enctype for the form, which
// is derived from its field types. Typically, this will be
// called from an HTML template like:
//
//	<form action="..." method="..." enctype="{{ .MyForm.EncType }}">
func (f *Form) EncType() string {
	if f != nil {
		for _, field := range f.fields {
			if field.Type == FILE {
				return "multipart/form-data"
			}
		}
	}
	return "application/x-www-form-urlencoded"
}

// Render renders all the fields in the form, in the order
// specified during construction.
func (f *Form) Render() (template.HTML, error) {
	if err := f.addCSRF(); err != nil {
		return template.HTML(""), err
	}
	return f.render(f.fields)
}

// RenderOnly renders the given fields, identified by their names
// in the struct. If a field does not exist, an error is returned.
// Fields are rendered according to the order of the parameters
// passed to this function.
func (f *Form) RenderOnly(names ...string) (template.HTML, error) {
	if err := f.addCSRF(); err != nil {
		return template.HTML(""), err
	}
	var fields []*Field
	for _, v := range names {
		field, err := f.lookupField(v)
		if err != nil {
			return template.HTML(""), err
		}
		fields = append(fields, field)
	}
	return f.render(fields)
}

// RenderExcept renders all the form's fields except the ones specified
// in the names parameter.
func (f *Form) RenderExcept(names ...string) (template.HTML, error) {
	if err := f.addCSRF(); err != nil {
		return template.HTML(""), err
	}
	n := make(map[string]bool, len(names))
	for _, v := range names {
		n[v] = true
	}
	var fields []*Field
	for _, v := range f.fields {
		if !n[v.Name] {
			fields = append(fields, v)
		}
	}
	return f.render(fields)
}

func (f *Form) addCSRF() error {
	if !f.hasCSRF && !f.DisableCSRF {
		csrf, err := newCSRF(f)
		if err != nil {
			return err
		}
		f.hasCSRF = true
		if err := f.appendVal(csrf); err != nil {
			return err
		}
		return f.makeFields(f.structs[len(f.structs)-1].QNames)
	}
	return nil
}

func (f *Form) makeId() {
	// Use the form pointer to generate the id,
	// to ensure uniqueness
	p, _ := strconv.ParseInt(fmt.Sprintf("%p", f), 0, 64)
	f.SetId(strconv.FormatInt(p%(1024*1024), 36))
}

// Id returns the prefix added to each field id in this form. Keep in mind
// that this function will never return an empty string because the form
// automatically generates a sufficiently unique id on creation.
func (f *Form) Id() string {
	return f.id
}

// SetId sets the prefix to be added to each field
// id attribute when rendering the form.
func (f *Form) SetId(id string) {
	f.id = id
	p := id + "_"
	for _, v := range f.fields {
		v.prefix = p
	}
}

// New is a shorthand for NewOpts(ctx, nil, values...). See
// NewOpts for more information.
func New(ctx *app.Context, values ...interface{}) *Form {
	return NewOpts(ctx, nil, values...)
}

// NewOpts returns a new form using the given context, renderer
// and options.
//
// If no Renderer is specified (either opts is nil or its Renderer field is
// nil), DefaultRenderer will be used to instantiate a renderer. Some
// packages from gondola/frontend, override DefaultRenderer when imported.
//
// The values argument must contains pointers to structs.
//
// Since any error generated during form creation will be a programming error,
// this function panics on errors.
//
// Consult the package documentation for the the tags parsed by the form library.
func NewOpts(ctx *app.Context, opts *Options, values ...interface{}) *Form {
	var r Renderer = nil
	if opts != nil {
		r = opts.Renderer
	}
	if r == nil {
		r = DefaultRenderer()
	}
	form := &Form{
		ctx:      ctx,
		renderer: r,
		options:  opts,
	}
	for _, v := range values {
		err := form.appendVal(v)
		if err != nil {
			panic(err)
		}
	}
	var fieldNames []string
	if opts != nil && len(opts.Fields) > 0 {
		fieldNames = opts.Fields
	} else {
		for _, v := range form.structs {
			fieldNames = append(fieldNames, v.QNames...)
		}
	}
	err := form.makeFields(fieldNames)
	if err != nil {
		panic(err)
	}
	form.makeId()
	return form
}

func fieldByIndex(v reflect.Value, indexes []int) reflect.Value {
	for _, idx := range indexes {
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.IsValid() {
			v = v.FieldByIndex([]int{idx})
		}
	}
	return v
}

func toHTMLValue(val interface{}) string {
	v := reflect.ValueOf(val)
	if v.IsValid() {
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.IsValid() {
			// Avoid enum types with a String() method to be represented
			// as a string. Use their numeric representation.
			k := types.Kind(v.Kind())
			if k == types.Int {
				return strconv.FormatInt(v.Int(), 10)
			}
			if k == types.Uint {
				return strconv.FormatUint(v.Uint(), 10)
			}
		}
	}
	return html.Escape(types.ToString(val))
}
