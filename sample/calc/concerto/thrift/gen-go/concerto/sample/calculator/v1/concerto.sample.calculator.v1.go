// Code generated by Thrift Compiler (0.16.0). DO NOT EDIT.

package v1

import (
	"bytes"
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
	thrift "github.com/apache/thrift/lib/go/thrift"
	"github.com/gorhythm/concerto/sample/calc/concerto/thrift/gen-go/concerto/rpc/status"

)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = context.Background
var _ = time.Now
var _ = bytes.Equal

var _ = status.GoUnusedProtection__
type Op int64
const (
  Op_ADD Op = 0
  Op_SUBTRACT Op = 1
  Op_MULTIPLY Op = 2
  Op_DIVIDE Op = 3
)

func (p Op) String() string {
  switch p {
  case Op_ADD: return "ADD"
  case Op_SUBTRACT: return "SUBTRACT"
  case Op_MULTIPLY: return "MULTIPLY"
  case Op_DIVIDE: return "DIVIDE"
  }
  return "<UNSET>"
}

func OpFromString(s string) (Op, error) {
  switch s {
  case "ADD": return Op_ADD, nil 
  case "SUBTRACT": return Op_SUBTRACT, nil 
  case "MULTIPLY": return Op_MULTIPLY, nil 
  case "DIVIDE": return Op_DIVIDE, nil 
  }
  return Op(0), fmt.Errorf("not a valid Op string")
}


func OpPtr(v Op) *Op { return &v }

func (p Op) MarshalText() ([]byte, error) {
return []byte(p.String()), nil
}

func (p *Op) UnmarshalText(text []byte) error {
q, err := OpFromString(string(text))
if (err != nil) {
return err
}
*p = q
return nil
}

func (p *Op) Scan(value interface{}) error {
v, ok := value.(int64)
if !ok {
return errors.New("Scan value is not int64")
}
*p = Op(v)
return nil
}

func (p * Op) Value() (driver.Value, error) {
  if p == nil {
    return nil, nil
  }
return int64(*p), nil
}
type CalculatorService interface {
  // Parameters:
  //  - Op
  //  - Num1
  //  - Num2
  Calculate(ctx context.Context, op Op, num1 int64, num2 int64) (_r int64, _err error)
}

type CalculatorServiceClient struct {
  c thrift.TClient
  meta thrift.ResponseMeta
}

func NewCalculatorServiceClientFactory(t thrift.TTransport, fn thrift.TProtocolFactory) *CalculatorServiceClient {
  return &CalculatorServiceClient{
    c: thrift.NewTStandardClient(fn.GetProtocol(t), fn.GetProtocol(t)),
  }
}

func NewCalculatorServiceClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *CalculatorServiceClient {
  return &CalculatorServiceClient{
    c: thrift.NewTStandardClient(iprot, oprot),
  }
}

func NewCalculatorServiceClient(c thrift.TClient) *CalculatorServiceClient {
  return &CalculatorServiceClient{
    c: c,
  }
}

func (p *CalculatorServiceClient) Client_() thrift.TClient {
  return p.c
}

func (p *CalculatorServiceClient) LastResponseMeta_() thrift.ResponseMeta {
  return p.meta
}

func (p *CalculatorServiceClient) SetLastResponseMeta_(meta thrift.ResponseMeta) {
  p.meta = meta
}

// Parameters:
//  - Op
//  - Num1
//  - Num2
func (p *CalculatorServiceClient) Calculate(ctx context.Context, op Op, num1 int64, num2 int64) (_r int64, _err error) {
  var _args0 CalculatorServiceCalculateArgs
  _args0.Op = op
  _args0.Num1 = num1
  _args0.Num2 = num2
  var _result2 CalculatorServiceCalculateResult
  var _meta1 thrift.ResponseMeta
  _meta1, _err = p.Client_().Call(ctx, "calculate", &_args0, &_result2)
  p.SetLastResponseMeta_(_meta1)
  if _err != nil {
    return
  }
  switch {
  case _result2.E1!= nil:
    return _r, _result2.E1
  }

  return _result2.GetSuccess(), nil
}

type CalculatorServiceProcessor struct {
  processorMap map[string]thrift.TProcessorFunction
  handler CalculatorService
}

func (p *CalculatorServiceProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
  p.processorMap[key] = processor
}

func (p *CalculatorServiceProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
  processor, ok = p.processorMap[key]
  return processor, ok
}

func (p *CalculatorServiceProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
  return p.processorMap
}

func NewCalculatorServiceProcessor(handler CalculatorService) *CalculatorServiceProcessor {

  self3 := &CalculatorServiceProcessor{handler:handler, processorMap:make(map[string]thrift.TProcessorFunction)}
  self3.processorMap["calculate"] = &calculatorServiceProcessorCalculate{handler:handler}
return self3
}

func (p *CalculatorServiceProcessor) Process(ctx context.Context, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  name, _, seqId, err2 := iprot.ReadMessageBegin(ctx)
  if err2 != nil { return false, thrift.WrapTException(err2) }
  if processor, ok := p.GetProcessorFunction(name); ok {
    return processor.Process(ctx, seqId, iprot, oprot)
  }
  iprot.Skip(ctx, thrift.STRUCT)
  iprot.ReadMessageEnd(ctx)
  x4 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function " + name)
  oprot.WriteMessageBegin(ctx, name, thrift.EXCEPTION, seqId)
  x4.Write(ctx, oprot)
  oprot.WriteMessageEnd(ctx)
  oprot.Flush(ctx)
  return false, x4

}

type calculatorServiceProcessorCalculate struct {
  handler CalculatorService
}

func (p *calculatorServiceProcessorCalculate) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := CalculatorServiceCalculateArgs{}
  var err2 error
  if err2 = args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err2.Error())
    oprot.WriteMessageBegin(ctx, "calculate", thrift.EXCEPTION, seqId)
    x.Write(ctx, oprot)
    oprot.WriteMessageEnd(ctx)
    oprot.Flush(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  // Start a goroutine to do server side connectivity check.
  if thrift.ServerConnectivityCheckInterval > 0 {
    var cancel context.CancelFunc
    ctx, cancel = context.WithCancel(ctx)
    defer cancel()
    var tickerCtx context.Context
    tickerCtx, tickerCancel = context.WithCancel(context.Background())
    defer tickerCancel()
    go func(ctx context.Context, cancel context.CancelFunc) {
      ticker := time.NewTicker(thrift.ServerConnectivityCheckInterval)
      defer ticker.Stop()
      for {
        select {
        case <-ctx.Done():
          return
        case <-ticker.C:
          if !iprot.Transport().IsOpen() {
            cancel()
            return
          }
        }
      }
    }(tickerCtx, cancel)
  }

  result := CalculatorServiceCalculateResult{}
  var retval int64
  if retval, err2 = p.handler.Calculate(ctx, args.Op, args.Num1, args.Num2); err2 != nil {
    tickerCancel()
  switch v := err2.(type) {
    case *status.Error:
  result.E1 = v
    default:
    if err2 == thrift.ErrAbandonRequest {
      return false, thrift.WrapTException(err2)
    }
    x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing calculate: " + err2.Error())
    oprot.WriteMessageBegin(ctx, "calculate", thrift.EXCEPTION, seqId)
    x.Write(ctx, oprot)
    oprot.WriteMessageEnd(ctx)
    oprot.Flush(ctx)
    return true, thrift.WrapTException(err2)
  }
  } else {
    result.Success = &retval
  }
  tickerCancel()
  if err2 = oprot.WriteMessageBegin(ctx, "calculate", thrift.REPLY, seqId); err2 != nil {
    err = thrift.WrapTException(err2)
  }
  if err2 = result.Write(ctx, oprot); err == nil && err2 != nil {
    err = thrift.WrapTException(err2)
  }
  if err2 = oprot.WriteMessageEnd(ctx); err == nil && err2 != nil {
    err = thrift.WrapTException(err2)
  }
  if err2 = oprot.Flush(ctx); err == nil && err2 != nil {
    err = thrift.WrapTException(err2)
  }
  if err != nil {
    return
  }
  return true, err
}


// HELPER FUNCTIONS AND STRUCTURES

// Attributes:
//  - Op
//  - Num1
//  - Num2
type CalculatorServiceCalculateArgs struct {
  Op Op `thrift:"op,1" db:"op" json:"op"`
  Num1 int64 `thrift:"num1,2" db:"num1" json:"num1"`
  Num2 int64 `thrift:"num2,3" db:"num2" json:"num2"`
}

func NewCalculatorServiceCalculateArgs() *CalculatorServiceCalculateArgs {
  return &CalculatorServiceCalculateArgs{}
}


func (p *CalculatorServiceCalculateArgs) GetOp() Op {
  return p.Op
}

func (p *CalculatorServiceCalculateArgs) GetNum1() int64 {
  return p.Num1
}

func (p *CalculatorServiceCalculateArgs) GetNum2() int64 {
  return p.Num2
}
func (p *CalculatorServiceCalculateArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I32 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *CalculatorServiceCalculateArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI32(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  temp := Op(v)
  p.Op = temp
}
  return nil
}

func (p *CalculatorServiceCalculateArgs)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.Num1 = v
}
  return nil
}

func (p *CalculatorServiceCalculateArgs)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  p.Num2 = v
}
  return nil
}

func (p *CalculatorServiceCalculateArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "calculate_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *CalculatorServiceCalculateArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "op", thrift.I32, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:op: ", p), err) }
  if err := oprot.WriteI32(ctx, int32(p.Op)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.op (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:op: ", p), err) }
  return err
}

func (p *CalculatorServiceCalculateArgs) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "num1", thrift.I64, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:num1: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.Num1)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.num1 (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:num1: ", p), err) }
  return err
}

func (p *CalculatorServiceCalculateArgs) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "num2", thrift.I64, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:num2: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.Num2)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.num2 (3) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:num2: ", p), err) }
  return err
}

func (p *CalculatorServiceCalculateArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("CalculatorServiceCalculateArgs(%+v)", *p)
}

// Attributes:
//  - Success
//  - E1
type CalculatorServiceCalculateResult struct {
  Success *int64 `thrift:"success,0" db:"success" json:"success,omitempty"`
  E1 *status.Error `thrift:"e1,1" db:"e1" json:"e1,omitempty"`
}

func NewCalculatorServiceCalculateResult() *CalculatorServiceCalculateResult {
  return &CalculatorServiceCalculateResult{}
}

var CalculatorServiceCalculateResult_Success_DEFAULT int64
func (p *CalculatorServiceCalculateResult) GetSuccess() int64 {
  if !p.IsSetSuccess() {
    return CalculatorServiceCalculateResult_Success_DEFAULT
  }
return *p.Success
}
var CalculatorServiceCalculateResult_E1_DEFAULT *status.Error
func (p *CalculatorServiceCalculateResult) GetE1() *status.Error {
  if !p.IsSetE1() {
    return CalculatorServiceCalculateResult_E1_DEFAULT
  }
return p.E1
}
func (p *CalculatorServiceCalculateResult) IsSetSuccess() bool {
  return p.Success != nil
}

func (p *CalculatorServiceCalculateResult) IsSetE1() bool {
  return p.E1 != nil
}

func (p *CalculatorServiceCalculateResult) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 0:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField0(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 1:
      if fieldTypeId == thrift.STRUCT {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *CalculatorServiceCalculateResult)  ReadField0(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 0: ", err)
} else {
  p.Success = &v
}
  return nil
}

func (p *CalculatorServiceCalculateResult)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  p.E1 = &status.Error{}
  if err := p.E1.Read(ctx, iprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.E1), err)
  }
  return nil
}

func (p *CalculatorServiceCalculateResult) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "calculate_result"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField0(ctx, oprot); err != nil { return err }
    if err := p.writeField1(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *CalculatorServiceCalculateResult) writeField0(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetSuccess() {
    if err := oprot.WriteFieldBegin(ctx, "success", thrift.I64, 0); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 0:success: ", p), err) }
    if err := oprot.WriteI64(ctx, int64(*p.Success)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.success (0) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 0:success: ", p), err) }
  }
  return err
}

func (p *CalculatorServiceCalculateResult) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetE1() {
    if err := oprot.WriteFieldBegin(ctx, "e1", thrift.STRUCT, 1); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:e1: ", p), err) }
    if err := p.E1.Write(ctx, oprot); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.E1), err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 1:e1: ", p), err) }
  }
  return err
}

func (p *CalculatorServiceCalculateResult) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("CalculatorServiceCalculateResult(%+v)", *p)
}


