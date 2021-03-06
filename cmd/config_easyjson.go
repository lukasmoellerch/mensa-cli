// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package cmd

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson6615c02eDecodeGithubComLukasmoellerchMensaCliCmd(in *jlexer.Lexer, out *group) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Name = string(in.String())
		case "refs":
			if in.IsNull() {
				in.Skip()
				out.Refs = nil
			} else {
				in.Delim('[')
				if out.Refs == nil {
					if !in.IsDelim(']') {
						out.Refs = make([]canteenRef, 0, 2)
					} else {
						out.Refs = []canteenRef{}
					}
				} else {
					out.Refs = (out.Refs)[:0]
				}
				for !in.IsDelim(']') {
					var v1 canteenRef
					(v1).UnmarshalEasyJSON(in)
					out.Refs = append(out.Refs, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6615c02eEncodeGithubComLukasmoellerchMensaCliCmd(out *jwriter.Writer, in group) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"refs\":"
		out.RawString(prefix)
		if in.Refs == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Refs {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v group) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6615c02eEncodeGithubComLukasmoellerchMensaCliCmd(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v group) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6615c02eEncodeGithubComLukasmoellerchMensaCliCmd(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *group) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6615c02eDecodeGithubComLukasmoellerchMensaCliCmd(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *group) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6615c02eDecodeGithubComLukasmoellerchMensaCliCmd(l, v)
}
func easyjson6615c02eDecodeGithubComLukasmoellerchMensaCliCmd1(in *jlexer.Lexer, out *config) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "enabled_providers":
			if in.IsNull() {
				in.Skip()
				out.EnabledProviders = nil
			} else {
				in.Delim('[')
				if out.EnabledProviders == nil {
					if !in.IsDelim(']') {
						out.EnabledProviders = make([]string, 0, 4)
					} else {
						out.EnabledProviders = []string{}
					}
				} else {
					out.EnabledProviders = (out.EnabledProviders)[:0]
				}
				for !in.IsDelim(']') {
					var v4 string
					v4 = string(in.String())
					out.EnabledProviders = append(out.EnabledProviders, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "groups":
			if in.IsNull() {
				in.Skip()
				out.Groups = nil
			} else {
				in.Delim('[')
				if out.Groups == nil {
					if !in.IsDelim(']') {
						out.Groups = make([]group, 0, 1)
					} else {
						out.Groups = []group{}
					}
				} else {
					out.Groups = (out.Groups)[:0]
				}
				for !in.IsDelim(']') {
					var v5 group
					(v5).UnmarshalEasyJSON(in)
					out.Groups = append(out.Groups, v5)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6615c02eEncodeGithubComLukasmoellerchMensaCliCmd1(out *jwriter.Writer, in config) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"enabled_providers\":"
		out.RawString(prefix[1:])
		if in.EnabledProviders == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v6, v7 := range in.EnabledProviders {
				if v6 > 0 {
					out.RawByte(',')
				}
				out.String(string(v7))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"groups\":"
		out.RawString(prefix)
		if in.Groups == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Groups {
				if v8 > 0 {
					out.RawByte(',')
				}
				(v9).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v config) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6615c02eEncodeGithubComLukasmoellerchMensaCliCmd1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v config) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6615c02eEncodeGithubComLukasmoellerchMensaCliCmd1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *config) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6615c02eDecodeGithubComLukasmoellerchMensaCliCmd1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *config) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6615c02eDecodeGithubComLukasmoellerchMensaCliCmd1(l, v)
}
func easyjson6615c02eDecodeGithubComLukasmoellerchMensaCliCmd2(in *jlexer.Lexer, out *canteenRef) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "provider":
			out.Provider = string(in.String())
		case "id":
			out.Id = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6615c02eEncodeGithubComLukasmoellerchMensaCliCmd2(out *jwriter.Writer, in canteenRef) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"provider\":"
		out.RawString(prefix[1:])
		out.String(string(in.Provider))
	}
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.String(string(in.Id))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v canteenRef) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6615c02eEncodeGithubComLukasmoellerchMensaCliCmd2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v canteenRef) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6615c02eEncodeGithubComLukasmoellerchMensaCliCmd2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *canteenRef) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6615c02eDecodeGithubComLukasmoellerchMensaCliCmd2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *canteenRef) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6615c02eDecodeGithubComLukasmoellerchMensaCliCmd2(l, v)
}
