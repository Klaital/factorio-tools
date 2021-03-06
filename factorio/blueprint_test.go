package factorio

import (
	log "github.com/sirupsen/logrus"
	"strings"
	"testing"
)

func TestParseBlueprintString(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	testString := `0eNqdm91u6kYUhV+l8jUcef78w2VfoZfVUUWIG1kCg4ypGkW8e02JzqHpzPhbXEUh4mPtySzPzNrDR/Gyv3SnsR+mYvNR9LvjcC42v38U5/5t2O5vr03vp67YFP3UHYpVMWwPt9+mcTucT8dxWr90+6m4rop+eO3+Ljbmulp883k6Dt36z8s4bHfdw3vt9fuq6Iapn/ruruLfX97/GC6Hl26c4anPXxWn43l+23G4feiMWhu/Kt7nn2HGv/Zjt7v/0d7kfaFaneqvEY4TOBar8wLVYWr4Qe2HczdO82sRXpmpthJ0lVhXLVANpjak2jZTbPtz7h62+/26288fN/a79em47yKsJsO6jQatscElGsEbLadaMnJVrlrBFRXXJbii5lTkipCrVnAFf0AZwRWBU5ErbK5a1RYuA7OCLfjT0wq24E9Pi2xhctVyW/CHp1XWCk4lrshNE8tNIfxjuSeEUoklspNYdER2O8ENwR8ljvuB28ERO+QWCcfdwNcIx93AH5qOmKHOlcrNwBd/x83AV0NHzJDdIjrRDbkdmOduELabntuBb5o8sYPJrQ6e+0HYBXuvH3RcVF7QQTYKqhJnwtgZ4tunS+23EGXVmNUsoRqMqpZQLUb5BVQoMcouoQxGLQ17sBRllkiOktwSyVNSWCIFSqqXSHiqt0skPNHN4pA/TPTTvp8S+9hPyDwVQLDQys+G8r9UHwsbSv2JY6KphUEVu0TFVQxp5dynBFSnJDTxWr0c0ZTL/+AqCGFKXFelRhVEVq2mKgTaPJf7xMtuhSwlSqhLNfUANdZGDWgI1ApJSrxWp2YeRJZX4xkCDUKOEq+1UkMKIqtW8xQCbZ5LfOJlt0KKEiU0pRp4gBobI0YzhGl5hhKv1IlxBxHlxVyGMANPUOKFVmJAQUTVYpJCmM1TSU+85pbnJ1FAW4pRByiwNWImQ5iWpyfxQp0YdBBRXkxkCDMI2Um80kpMJoiqWs1QCLR5KuSJF90KyYmJN9dKNToBNZqSm8FYTuUnB+M4VVgqPKf+dMlleO3Gt/E4/0wutvfHjPkCXv24DzCcLrdbA5HP0bMlu3x+NGUlYw3B1tKouPyoHC9Tclga+ZCZ8EcrnCvJEJhSz/AScYR5aGIvni2RNitngklp7stjbr89nJInwPjYP3StF0+AqLwg55TJ8ip+YEPSajn3TEoT+tRIWivnqClpDz3r7KRwmUlhDT8FkfKslbPdZHkOn1qQMq9GxUlh/JSBhFVq8pwUVqMJ4XPzocGHAVRbq2bhqdoe+tFLm3cizBk1Wk8Ks3yzjZQ5NapPKvN8P4uUBTn6T0rTN0SeKKxlrCNYNXf63GolrqQpeyBStS/5vgXxzFO5c7xaL9zhQ+Ic3yognufrO+KFpzLKxOBVfHVG4mq8pCIcv7CEcE9dX4qPXCjxOkakBYNXH4QT1gzEc8+kIImRE1YNpE0/UqPLspV8JEXYWr7wjrCN3AxF2Fa+sE6wlXxHHFHlO+KIatUGIqI6tdeJqHIDEFGD2qtE1Ept4CFqrfYaEbURe2YI2qrdPXSXvhS7XghqxPYcglqxbYWgTuyvIagXe0UIGsSmFoKq128RVL1+i6CN2qxB1FZsK6Evksg9l/9Rv6/u3xHcPHwfcVXstzNlfu3X7bnf/fLbYUb2w9v8h7+68Xy/OtXU1ph5d+iq6/Ufs1ywUg==`

	blueprint, bpErr := ParseBlueprintString(actualBase64String)
	if bpErr != nil {
		t.Errorf("Failed to parse BP string: %s", bpErr.Error())
		t.Fail()
	}

	if "Basic Smelting" != blueprint.Details.Label {
		t.Errorf("Failed to read item name. Expected %s, got %s", "Basic Smelting", blueprint.Details.Label)
	}
}

//func TestParseBpBookString(t *testing.T) {
//	log.SetLevel(log.DebugLevel)
//	testString, err := ioutil.ReadFile("testdata/bp_book1.txt")
//	if err != nil {
//		t.Errorf("Failed to read BP Book from file: %v", err)
//		t.Fail()
//	}
//
//	actualBase64String := strings.TrimPrefix(string(testString), "0")
//	bpBook, bpErr := ParseBlueprintBookString(actualBase64String)
//	if bpErr != nil {
//		t.Errorf("Failed to parse BP book string: %s", bpErr.Error())
//		t.Fail()
//	}
//	//	log.Debugln("Book:", *bpBook)
//
//	if bpBook == nil {
//		t.Error("No Blueprint Book was generated")
//		t.Fail()
//	}
//
//	if "blueprint-book" != bpBook.Item {
//		t.Errorf("Failed to read item name. Expected %s, got %s", "blueprint-book", bpBook.Item)
//	}
//}
