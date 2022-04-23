
package wordfilter

import (
	"testing"
)

func TestInitConfigFilter(t *testing.T) {
	InitConfigFilter("../cfgs")
   	rStr :=	 FilterChack("as1111")
   	t.Log(rStr)
	rStr =	 FilterChack("4r5e 5h1t")
	t.Log(rStr)
	rStr =	 FilterChack("assf a sshole")
	t.Log(rStr)
	rStr =	 FilterChack("a 5 5")
	t.Log(rStr)
}