package types

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/ehwjh2010/cobra/config"
	"github.com/ehwjh2010/cobra/util/intutils"
	"github.com/ehwjh2010/cobra/util/jsonutils"
	"github.com/ehwjh2010/cobra/util/timeutils"
	"time"
)

//********************int64*****************************

// NullInt64 is an alias for sql.NullInt64 data type
type NullInt64 struct {
	sql.NullInt64
}

//IsNil 是否是Nil
func (ni *NullInt64) IsNil() bool {
	return !ni.NullInt64.Valid
}

//Equal 比较是否相等
func (ni *NullInt64) Equal(v NullInt64) bool {
	return ni.Valid == v.Valid && (!ni.Valid || ni.Int64 == v.Int64)
}

//GetValue 获取值
func (ni *NullInt64) GetValue() int64 {
	return ni.Int64
}

func NewInt64(v int64) NullInt64 {
	return NullInt64{NullInt64: sql.NullInt64{
		Int64: v,
		Valid: true,
	}}
}

func NewInt64Null() NullInt64 {
	return NullInt64{}
}

// MarshalJSON for NullInt64
func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return jsonutils.Marshal(ni.Int64)
}

//UnmarshalJSON for NullInt64
func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, config.NullBytes) {
		ni.Valid = false
		return nil
	}

	err := jsonutils.Unmarshal(b, &ni.Int64)

	if err != nil {
		ni.Valid = false
	} else {
		ni.Valid = true
	}

	return err
}

//********************int64*****************************

// NullInt is an alias for sql.NullInt64 data type
type NullInt struct {
	sql.NullInt64
}

//IsNil 是否是Nil
func (ni *NullInt) IsNil() bool {
	return !ni.Valid
}

//GetValue 获取值
func (ni *NullInt) GetValue() int {
	return intutils.Int64ToInt(ni.Int64)
}

func NewInt(v int) NullInt {
	return NullInt{NullInt64: sql.NullInt64{
		Int64: int64(v),
		Valid: true,
	}}
}

func NewIntNull() NullInt {
	return NullInt{}
}

// MarshalJSON for NullInt
func (ni NullInt) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return jsonutils.Marshal(intutils.Int64ToInt(ni.Int64))
}

//UnmarshalJSON for NullInt
func (ni *NullInt) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, config.NullBytes) {
		ni.Valid = false
		return nil
	}

	err := jsonutils.Unmarshal(b, &ni.Int64)
	if err != nil {
		ni.Valid = false
	} else {
		ni.Valid = true
	}

	return err
}

//Equal 比较是否相等
func (ni *NullInt) Equal(v NullInt) bool {
	return ni.Valid == v.Valid && (!ni.Valid || ni.Int64 == v.Int64)
}

//********************bool*****************************

// NullBool is an alias for sql.NullBool data type
type NullBool struct {
	sql.NullBool
}

//IsNil 是否是Nil
func (nb *NullBool) IsNil() bool {
	return !nb.NullBool.Valid
}

//GetValue 获取值
func (nb *NullBool) GetValue() bool {
	return nb.NullBool.Bool
}

// MarshalJSON for NullBool
func (nb NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	return jsonutils.Marshal(nb.Bool)
}

func NewBool(v bool) NullBool {
	return NullBool{NullBool: sql.NullBool{
		Bool:  v,
		Valid: true,
	}}
}

func NewBoolNull() NullBool {
	return NullBool{}
}

//UnmarshalJSON for NullBool
func (nb *NullBool) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, config.NullBytes) {
		nb.Valid = false
		return nil
	}

	err := jsonutils.Unmarshal(b, &nb.Bool)
	if err != nil {
		nb.Valid = false
	} else {
		nb.Valid = true
	}

	return err
}

//Equal 比较是否相等
func (nb *NullBool) Equal(v NullBool) bool {
	return nb.Valid == v.Valid && (!nb.Valid || nb.Bool == v.Bool)
}

//********************float64*****************************

// NullFloat64 is an alias for sql.NullFloat64 data type
type NullFloat64 struct {
	sql.NullFloat64
}

//IsNil 是否是Nil
func (nf *NullFloat64) IsNil() bool {
	return !nf.NullFloat64.Valid
}

//GetValue 获取值
func (nf *NullFloat64) GetValue() float64 {
	return nf.NullFloat64.Float64
}

// MarshalJSON for NullFloat64
func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return jsonutils.Marshal(nf.Float64)
}

func NewFloat64(v float64) NullFloat64 {
	return NullFloat64{NullFloat64: sql.NullFloat64{
		Float64: v,
		Valid:   true,
	}}
}

func NewFloat64Null() NullFloat64 {
	return NullFloat64{}
}

// UnmarshalJSON for NullFloat64
func (nf *NullFloat64) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, config.NullBytes) {
		nf.Valid = false
		return nil
	}

	err := jsonutils.Unmarshal(b, &nf.Float64)
	if err != nil {
		nf.Valid = false
	} else {
		nf.Valid = true
	}

	return err
}

//Equal 比较是否相等
func (nf *NullFloat64) Equal(v NullFloat64) bool {
	return nf.Valid == v.Valid && (!nf.Valid || nf.Float64 == v.Float64)
}

//********************string*****************************

// NullString is an alias for sql.NullString data type
type NullString struct {
	sql.NullString
}

//IsNil 是否是Nil
func (ns *NullString) IsNil() bool {
	return !ns.NullString.Valid
}

//GetValue 获取值
func (ns *NullString) GetValue() string {
	return ns.NullString.String
}

func NewStr(str string) NullString {
	return NullString{NullString: sql.NullString{
		String: str,
		Valid:  true,
	}}
}

func NewEmptyStr() NullString {
	return NullString{NullString: sql.NullString{
		String: "",
		Valid:  true,
	}}
}

func NewStrNull() NullString {
	return NullString{}
}

// MarshalJSON for NullString
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return jsonutils.Marshal(ns.String)
}

// UnmarshalJSON for NullString
func (ns *NullString) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, config.NullBytes) {
		ns.Valid = false
		return nil
	}

	err := jsonutils.Unmarshal(b, &ns.String)
	if err != nil {
		ns.Valid = false
	} else {
		ns.Valid = true
	}

	return err
}

//Equal 比较是否相等
func (ns *NullString) Equal(v NullString) bool {
	return ns.Valid == v.Valid && (!ns.Valid || ns.String == v.String)
}

//********************time*****************************

// NullTime is an alias for mysql.NullTime data type
type NullTime struct {
	sql.NullTime
}

//IsNil 是否是Nil
func (nt *NullTime) IsNil() bool {
	return !nt.Valid
}

//GetValue 获取值
func (nt *NullTime) GetValue() time.Time {
	return nt.Time
}

// MarshalJSON for NullTime
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", nt.Time.In(timeutils.GetBJLocation()).Format(config.DefaultTimePattern))
	return []byte(val), nil
}

func NewTime(t time.Time) NullTime {
	return NullTime{NullTime: sql.NullTime{
		Time:  t,
		Valid: true,
	}}
}

func NewTimeNull() NullTime {
	return NullTime{}
}

// UnmarshalJSON for NullTime
func (nt *NullTime) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, config.NullBytes) {
		nt.Valid = false
		return nil
	}

	err := jsonutils.Unmarshal(b, &nt.Time)
	if err != nil {
		nt.Valid = false
	} else {
		nt.Valid = true
	}

	return err
}

//Equal 比较是否相等
func (nt *NullTime) Equal(v NullTime) bool {
	return nt.Valid == v.Valid && (!nt.Valid || nt.Time == v.Time)
}
