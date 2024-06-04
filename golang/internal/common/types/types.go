package types

type Any interface{}
type Param interface{}
type StringMap map[string]string
type JSON map[string]interface{} // This does not include Arrays
type Func func(...Param)
