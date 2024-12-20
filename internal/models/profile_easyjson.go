// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

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

func easyjson521a5691DecodeGithubCom20242BetterCallFirewallInternalModels(in *jlexer.Lexer, out *ShortProfile) {
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
		case "id":
			out.ID = uint32(in.Uint32())
		case "first_name":
			out.FirstName = string(in.String())
		case "last_name":
			out.LastName = string(in.String())
		case "is_author":
			out.IsAuthor = bool(in.Bool())
		case "is_friend":
			out.IsFriend = bool(in.Bool())
		case "is_subscriber":
			out.IsSubscriber = bool(in.Bool())
		case "is_subscription":
			out.IsSubscription = bool(in.Bool())
		case "avatar":
			out.Avatar = Picture(in.String())
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
func easyjson521a5691EncodeGithubCom20242BetterCallFirewallInternalModels(out *jwriter.Writer, in ShortProfile) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint32(uint32(in.ID))
	}
	{
		const prefix string = ",\"first_name\":"
		out.RawString(prefix)
		out.String(string(in.FirstName))
	}
	{
		const prefix string = ",\"last_name\":"
		out.RawString(prefix)
		out.String(string(in.LastName))
	}
	{
		const prefix string = ",\"is_author\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsAuthor))
	}
	{
		const prefix string = ",\"is_friend\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsFriend))
	}
	{
		const prefix string = ",\"is_subscriber\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsSubscriber))
	}
	{
		const prefix string = ",\"is_subscription\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsSubscription))
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.Avatar))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ShortProfile) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson521a5691EncodeGithubCom20242BetterCallFirewallInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ShortProfile) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson521a5691EncodeGithubCom20242BetterCallFirewallInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ShortProfile) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson521a5691DecodeGithubCom20242BetterCallFirewallInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ShortProfile) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson521a5691DecodeGithubCom20242BetterCallFirewallInternalModels(l, v)
}
func easyjson521a5691DecodeGithubCom20242BetterCallFirewallInternalModels1(in *jlexer.Lexer, out *FullProfile) {
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
		case "id":
			out.ID = uint32(in.Uint32())
		case "first_name":
			out.FirstName = string(in.String())
		case "last_name":
			out.LastName = string(in.String())
		case "bio":
			out.Bio = string(in.String())
		case "is_author":
			out.IsAuthor = bool(in.Bool())
		case "is_friend":
			out.IsFriend = bool(in.Bool())
		case "is_subscriber":
			out.IsSubscriber = bool(in.Bool())
		case "is_subscription":
			out.IsSubscription = bool(in.Bool())
		case "avatar":
			out.Avatar = Picture(in.String())
		case "pics":
			if in.IsNull() {
				in.Skip()
				out.Pics = nil
			} else {
				in.Delim('[')
				if out.Pics == nil {
					if !in.IsDelim(']') {
						out.Pics = make([]Picture, 0, 4)
					} else {
						out.Pics = []Picture{}
					}
				} else {
					out.Pics = (out.Pics)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Picture
					v1 = Picture(in.String())
					out.Pics = append(out.Pics, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "posts":
			if in.IsNull() {
				in.Skip()
				out.Posts = nil
			} else {
				in.Delim('[')
				if out.Posts == nil {
					if !in.IsDelim(']') {
						out.Posts = make([]*Post, 0, 8)
					} else {
						out.Posts = []*Post{}
					}
				} else {
					out.Posts = (out.Posts)[:0]
				}
				for !in.IsDelim(']') {
					var v2 *Post
					if in.IsNull() {
						in.Skip()
						v2 = nil
					} else {
						if v2 == nil {
							v2 = new(Post)
						}
						(*v2).UnmarshalEasyJSON(in)
					}
					out.Posts = append(out.Posts, v2)
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
func easyjson521a5691EncodeGithubCom20242BetterCallFirewallInternalModels1(out *jwriter.Writer, in FullProfile) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint32(uint32(in.ID))
	}
	{
		const prefix string = ",\"first_name\":"
		out.RawString(prefix)
		out.String(string(in.FirstName))
	}
	{
		const prefix string = ",\"last_name\":"
		out.RawString(prefix)
		out.String(string(in.LastName))
	}
	{
		const prefix string = ",\"bio\":"
		out.RawString(prefix)
		out.String(string(in.Bio))
	}
	{
		const prefix string = ",\"is_author\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsAuthor))
	}
	{
		const prefix string = ",\"is_friend\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsFriend))
	}
	{
		const prefix string = ",\"is_subscriber\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsSubscriber))
	}
	{
		const prefix string = ",\"is_subscription\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsSubscription))
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.Avatar))
	}
	{
		const prefix string = ",\"pics\":"
		out.RawString(prefix)
		if in.Pics == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v3, v4 := range in.Pics {
				if v3 > 0 {
					out.RawByte(',')
				}
				out.String(string(v4))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"posts\":"
		out.RawString(prefix)
		if in.Posts == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Posts {
				if v5 > 0 {
					out.RawByte(',')
				}
				if v6 == nil {
					out.RawString("null")
				} else {
					(*v6).MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v FullProfile) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson521a5691EncodeGithubCom20242BetterCallFirewallInternalModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v FullProfile) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson521a5691EncodeGithubCom20242BetterCallFirewallInternalModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *FullProfile) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson521a5691DecodeGithubCom20242BetterCallFirewallInternalModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *FullProfile) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson521a5691DecodeGithubCom20242BetterCallFirewallInternalModels1(l, v)
}
func easyjson521a5691DecodeGithubCom20242BetterCallFirewallInternalModels2(in *jlexer.Lexer, out *ChangePasswordReq) {
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
		case "old_password":
			out.OldPassword = string(in.String())
		case "new_password":
			out.NewPassword = string(in.String())
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
func easyjson521a5691EncodeGithubCom20242BetterCallFirewallInternalModels2(out *jwriter.Writer, in ChangePasswordReq) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"old_password\":"
		out.RawString(prefix[1:])
		out.String(string(in.OldPassword))
	}
	{
		const prefix string = ",\"new_password\":"
		out.RawString(prefix)
		out.String(string(in.NewPassword))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ChangePasswordReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson521a5691EncodeGithubCom20242BetterCallFirewallInternalModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ChangePasswordReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson521a5691EncodeGithubCom20242BetterCallFirewallInternalModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ChangePasswordReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson521a5691DecodeGithubCom20242BetterCallFirewallInternalModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ChangePasswordReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson521a5691DecodeGithubCom20242BetterCallFirewallInternalModels2(l, v)
}
