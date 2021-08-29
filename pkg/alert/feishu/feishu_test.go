package feishu

import (
	"fmt"
	"testing"
)

func Test_Alert(t *testing.T) {
	alert := NewFeishuRot(``)
	err := alert.AlertText("警告", `test`, `test`)
	fmt.Println(err)
}
