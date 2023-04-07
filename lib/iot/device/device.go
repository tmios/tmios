package device

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	validator "github.com/go-playground/validator/v10"
)

var (
	validate    *validator.Validate
	typeContext = reflect.TypeOf((*context.Context)(nil)).Elem()
	typeDevice  = reflect.TypeOf((*Device)(nil)).Elem()
	typeError   = reflect.TypeOf((*error)(nil)).Elem()

	ErrInvalidAction = errors.New("invalid action")
)

func jsonStr(data interface{}) string {
	s, _ := json.Marshal(data)
	return string(s)
}

func init() {
	validate = validator.New()
}

type PropMeta struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Validate string `json:"validate"`
	Desc     string `json:"desc"`
	Extras   string `json:"extras"`

	propType reflect.Type
}

type PropsMeta struct {
	Props []*PropMeta
	Type  reflect.Type `json:"-"`
}

func (p PropsMeta) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Props)
}

func (p *PropMeta) Check(val interface{}) error {
	if p.propType != reflect.TypeOf(val) {
		return fmt.Errorf("value is not %s type", p.Type)
	}

	return nil
}

func (meta PropsMeta) Get(name string) *PropMeta {
	for _, prop := range meta.Props {
		if prop.Name == name {
			return prop
		}
	}

	return nil
}

type PropVal json.RawMessage

func (propVal PropVal) Cast(propMeta *PropMeta) (interface{}, error) {
	val := reflect.New(propMeta.propType)
	if err := json.Unmarshal(propVal, val.Interface()); err != nil {
		return nil, err
	}

	return val.Elem().Interface(), nil
}

type ActionMeta struct {
	Name string    `json:"name"`
	Desc string    `json:"desc"`
	Args PropsMeta `json:"args"`
	Rets PropsMeta `json:"rets"`

	fun      reflect.Value
	argsType reflect.Type
	retsType reflect.Type
}

func (meta ActionMeta) Action(ctx context.Context, dv Device, args []byte) ([]byte, error) {
	var (
		argsVal = reflect.New(meta.argsType)
		retsVal = reflect.New(meta.retsType)
	)

	if dv == nil || ctx == nil {
		panic("context and device cannot be nil")
	}

	if err := json.Unmarshal(args, argsVal.Interface()); err != nil {
		return nil, err
	}
	if err := validate.Struct(argsVal.Interface()); err != nil {
		return nil, err
	}

	retVals := meta.fun.Call([]reflect.Value{
		reflect.ValueOf(ctx), reflect.ValueOf(dv),
		reflect.ValueOf(argsVal.Interface()),
		reflect.ValueOf(retsVal.Interface()),
	})

	err, ok := retVals[0].Interface().(error)
	if ok {
		return nil, err
	}

	return json.Marshal(retsVal.Interface())
}

type ActionsMeta []ActionMeta

type IntervalFunc func(Device) error

type Interval struct {
	Name     string       `json:"name"`
	Interval int64        `json:"interval"`
	Desc     string       `json:"desc"`
	Func     IntervalFunc `json:"-"`
}

type Intervals []Interval
type InitFunc func() error
type ForeignIDFunc func(config []byte) string

type DeviceMeta struct {
	Name     string    `json:"name"`     // 设备名称
	Variable string    `json:"variable"` // 变量名
	Brand    string    `json:"brand"`    // 品牌
	Model    string    `json:"model"`    // 型号
	Type     string    `json:"type"`     // 设备类型
	Desc     string    `json:"desc"`     // 描述
	Config   PropsMeta `json:"config"`

	Properties PropsMeta   `json:"properties"`
	Actions    ActionsMeta `json:"actions"`

	ForeignIDFunc ForeignIDFunc `json:"-"`
	InitFunc      InitFunc      `json:"-"`
	Intervals     Intervals     `json:"intervals"`
}

func (meta *DeviceMeta) GetAction(name string) *ActionMeta {
	for _, actionMeta := range meta.Actions {
		if actionMeta.Name == name {
			return &actionMeta
		}
	}

	return nil
}

func (meta *DeviceMeta) GetProp(name string) *PropMeta {
	for _, propMeta := range meta.Properties.Props {
		if propMeta.Name == name {
			return propMeta
		}
	}

	return nil
}

func (meta *DeviceMeta) CheckConfig(data []byte) ([]byte, error) {
	configVal := reflect.New(meta.Config.Type)
	if err := json.Unmarshal(data, configVal.Interface()); err != nil {
		return nil, err
	}

	if err := validate.Struct(configVal.Interface()); err != nil {
		return nil, err
	}

	return json.Marshal(configVal.Interface())
}

// propStruct can be reflect.Type or struct instance
func GetPropsMeta(propStruct interface{}) PropsMeta {
	var (
		ok    bool
		typ   reflect.Type
		props []*PropMeta
	)

	if typ, ok = propStruct.(reflect.Type); !ok {
		typ = reflect.TypeOf(propStruct)
	}

	numField := typ.NumField()

	for i := 0; i < numField; i++ {
		var (
			field     = typ.Field(i)
			name      = field.Tag.Get("json")
			filedType = field.Type.Name()
		)

		if name == "" {
			name = field.Name
		}
		props = append(props, &PropMeta{
			Name:     name,
			Type:     filedType,
			Desc:     field.Tag.Get("desc"),
			Validate: field.Tag.Get("validate"),
			Extras:   field.Tag.Get("extras"),
			propType: field.Type,
		})
	}

	return PropsMeta{
		Props: props,
		Type:  typ,
	}
}

func ToActionMeta(name string, actionFunc interface{}, desc string) ActionMeta {
	var (
		error = func() {
			panic(fmt.Sprintf("actionFunc: %s, %#v layout error.", name, actionFunc))
		}
	)

	funcType := reflect.TypeOf(actionFunc)
	if funcType.Kind() != reflect.Func {
		error()
	}

	if funcType.NumIn() != 4 {
		error()
	}

	if !funcType.In(0).Implements(typeContext) {
		error()
	}

	if !funcType.In(1).Implements(typeDevice) {
		error()
	}

	if funcType.NumOut() != 1 {
		error()
	}

	if funcType.Out(0) != typeError {
		error()
	}

	getArgsMeta := func(idx int) PropsMeta {
		if funcType.In(idx).Kind() != reflect.Ptr {
			error()
		}

		typeField := funcType.In(idx).Elem()
		if typeField.Kind() != reflect.Struct {
			error()
		}

		return GetPropsMeta(typeField)
	}

	return ActionMeta{
		Name: name,
		Desc: desc,
		Args: getArgsMeta(2),
		Rets: getArgsMeta(3),

		fun:      reflect.ValueOf(actionFunc),
		argsType: funcType.In(2).Elem(),
		retsType: funcType.In(3).Elem(),
	}
}

func Action(ctx context.Context, dv Device, name string, args []byte) ([]byte, error) {
	actionMeta := dv.Meta().GetAction(name)
	if actionMeta == nil {
		return nil, ErrInvalidAction
	}

	return actionMeta.Action(ctx, dv, args)
}

func NewSetValOptions(opts ...SetValOption) *SetValOptions {
	opt := SetValOptions{
		WriteRedis:  true,
		WriteInflux: true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return &opt
}

func SetVals(d Device, vals map[string]interface{}) error {
	for name, val := range vals {
		if err := d.SetVal(name, val); err != nil {
			return err
		}
	}

	return nil
}

type SetValOption func(o *SetValOptions)

type SetValOptions struct {
	WriteRedis  bool
	WriteInflux bool
}

var ErrNil = errors.New("RedisNil")

type Storage interface {
	// redis
	LRange(ctx context.Context, key string, start, end int) ([]string, error)
	LPush(ctx context.Context, key string, values ...string) error
	RPop(ctx context.Context, key string, c int) error

	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error

	// influxdb
	WritePoint(ctx context.Context, measurement string, tags map[string]string,
		fields map[string]interface{}, ts time.Time) error
}

type CommitAttr struct {
	KeepAlive bool
	UpdateAt  int64
}

type CommitOption func(attr *CommitAttr)

type Device interface {
	Meta() *DeviceMeta
	Action(ctx context.Context, name string, args []byte) ([]byte, error)

	GetConfig(config interface{}) error
	GetStorage() Storage

	Tags() map[string]string // For influxdb storage
	DebugMode() bool
	POM() POM

	// GetVal & SetVal only write data to memory. You should manual commit
	// data to redis & influxdb
	GetVal(name string) (interface{}, error)
	GetPropVals(in interface{}) error // unmarshal vals in
	SetVal(name string, val interface{}, opts ...SetValOption) error
	SetVals(vals map[string]interface{}, opts ...SetValOption) error
	Commit(opts ...CommitOption) error
}
