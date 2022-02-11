package object

import "github.com/jinzhu/copier"

// CopyProperties 拷贝属性, 支持struct, slice, map等
func CopyProperties(source interface{}, dst interface{}) {
	copier.Copy(dst, source)
}
