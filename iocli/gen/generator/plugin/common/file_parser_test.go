package common

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinMethodLine(t *testing.T) {
	goFile := `
func (s *ComplexService) RPCBasicType(name string,// some comments
	ageF64Ptr *float64, // some comments
	ageF64Ptr2 *float64,	// some comments
	ageF64Ptr3 *float64,			// some comments
) (string, error){
	return "", nil
}`
	fileLines := strings.Split(goFile, "\n")
	afterJoinedFile := fmt.Sprintf("%s", joinMethodLine(fileLines))
	assert.Equal(t, `[ 
 func (s *ComplexService)RPCBasicType(name string, ageF64Ptr *float64, ageF64Ptr2 *float64, ageF64Ptr3 *float64)(string, error){ 
 return "", nil 
 } 
]`, afterJoinedFile)

	goFile = `
func (c *context) filterAndGetRecord(ctx *aop.InvocationContext) (methodInvocationRecordIOCInterface, bool) {
	// filter invocations that is not monitored
	if c.sdid != "" && ctx.SDID != c.sdid {
		return nil, false
	} else if c.methodName != "" && ctx.MethodName != c.methodName {
		return nil, false
	}

	// monitor the invocation
	invocationMethodKey := common.GetMethodUniqueKey(ctx.SDID, ctx.MethodName)
	c.methodUniqueNameInvocationRecordMapLock.Lock()
	defer c.methodUniqueNameInvocationRecordMapLock.Unlock()
	methodRecord, ok := c.methodUniqueNameInvocationRecordMap[invocationMethodKey]
	if !ok {
		methodRecord, _ = GetmethodInvocationRecordIOCInterface()
		c.methodUniqueNameInvocationRecordMap[invocationMethodKey] = methodRecord
	}
	return methodRecord, true
}`
	fileLines = strings.Split(goFile, "\n")
	afterJoinedFile = fmt.Sprintf("%s", joinMethodLine(fileLines))

	assert.Equal(t, `[ 
 func (c *context) filterAndGetRecord(ctx *aop.InvocationContext) (methodInvocationRecordIOCInterface, bool) { 
  
 if c.sdid != "" && ctx.SDID != c.sdid { 
 return nil, false 
 } else if c.methodName != "" && ctx.MethodName != c.methodName { 
 return nil, false 
 } 
  
  
 invocationMethodKey := common.GetMethodUniqueKey(ctx.SDID, ctx.MethodName) 
 c.methodUniqueNameInvocationRecordMapLock.Lock() 
 defer c.methodUniqueNameInvocationRecordMapLock.Unlock() 
 methodRecord, ok := c.methodUniqueNameInvocationRecordMap[invocationMethodKey] 
 if !ok { 
 methodRecord, _ = GetmethodInvocationRecordIOCInterface() 
 c.methodUniqueNameInvocationRecordMap[invocationMethodKey] = methodRecord 
 } 
 return methodRecord, true 
 } 
]`, afterJoinedFile)
}
