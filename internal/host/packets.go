package host

import (
	"encoding/json"
	"math"
	"reflect"
	"sync"

	"github.com/bedrock-gophers/plugins/internal/native"
	dfunsafe "github.com/bedrock-gophers/unsafe"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Packets owns callback-scoped handles to packets already decoded by
// gophertunnel. Handles become invalid as soon as the intercept callback ends.
type Packets struct {
	mu      sync.RWMutex
	next    native.PacketHandle
	packets map[native.PacketHandle]borrowedPacket
	players *Players
}

type borrowedPacket struct {
	value   packet.Packet
	mutable bool
}

func NewPackets(players *Players) *Packets {
	return &Packets{packets: map[native.PacketHandle]borrowedPacket{}, players: players}
}

func (p *Packets) WritePlayerPacket(invocation native.InvocationID, id native.PlayerID, handle native.PacketHandle) bool {
	value, ok := p.packet(handle)
	if !ok || p.players == nil {
		return false
	}
	return p.players.mutatePlayer(invocation, id, func(connected *player.Player) {
		dfunsafe.WritePacket(connected, value)
	})
}

func (p *Packets) packet(handle native.PacketHandle) (packet.Packet, bool) {
	p.mu.RLock()
	borrowed, ok := p.packets[handle]
	p.mu.RUnlock()
	return borrowed.value, ok && borrowed.value != nil
}

func (p *Packets) Borrow(value packet.Packet, mutable bool) (native.PacketHandle, func(), bool) {
	if p == nil || value == nil {
		return 0, func() {}, false
	}
	p.mu.Lock()
	if p.next == native.PacketHandle(math.MaxUint64) {
		p.mu.Unlock()
		return 0, func() {}, false
	}
	p.next++
	handle := p.next
	p.packets[handle] = borrowedPacket{value: value, mutable: mutable}
	p.mu.Unlock()
	return handle, func() {
		p.mu.Lock()
		delete(p.packets, handle)
		p.mu.Unlock()
	}, true
}

func (p *Packets) PacketField(handle native.PacketHandle, index uint32) (native.PacketFieldValue, bool) {
	field, _, ok := p.field(handle, index)
	if !ok {
		return native.PacketFieldValue{}, false
	}
	return readPacketField(field)
}

func (p *Packets) SetPacketField(handle native.PacketHandle, index uint32, value native.PacketFieldValue) bool {
	field, mutable, ok := p.field(handle, index)
	return ok && mutable && field.CanSet() && writePacketField(field, value)
}

func (p *Packets) field(handle native.PacketHandle, index uint32) (reflect.Value, bool, bool) {
	p.mu.RLock()
	borrowed, ok := p.packets[handle]
	p.mu.RUnlock()
	if !ok {
		return reflect.Value{}, false, false
	}
	root := reflect.ValueOf(borrowed.value)
	if root.Kind() != reflect.Pointer || root.IsNil() || root.Elem().Kind() != reflect.Struct || uint64(index) >= uint64(root.Elem().NumField()) {
		return reflect.Value{}, false, false
	}
	field := root.Elem().Field(int(index))
	return field, borrowed.mutable, field.CanInterface()
}

func readPacketField(field reflect.Value) (native.PacketFieldValue, bool) {
	result := native.PacketFieldValue{}
	switch field.Kind() {
	case reflect.Bool:
		result.Kind = native.PacketFieldBool
		if field.Bool() {
			result.Unsigned = 1
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result.Kind, result.Signed = native.PacketFieldSigned, field.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result.Kind, result.Unsigned = native.PacketFieldUnsigned, field.Uint()
	case reflect.Float32, reflect.Float64:
		result.Kind, result.Number = native.PacketFieldFloat, field.Float()
	case reflect.String:
		result.Kind, result.Data = native.PacketFieldString, []byte(field.String())
	case reflect.Slice:
		if field.Type().Elem().Kind() != reflect.Uint8 {
			return jsonPacketField(field)
		}
		result.Kind, result.Data = native.PacketFieldBytes, append([]byte(nil), field.Bytes()...)
	case reflect.Array:
		switch {
		case field.Type().PkgPath() == "github.com/go-gl/mathgl/mgl32" && field.Type().Name() == "Vec2":
			result.Kind = native.PacketFieldVec2
			result.X, result.Y = field.Index(0).Float(), field.Index(1).Float()
		case field.Type().PkgPath() == "github.com/go-gl/mathgl/mgl32" && field.Type().Name() == "Vec3":
			result.Kind = native.PacketFieldVec3
			result.X, result.Y, result.Z = field.Index(0).Float(), field.Index(1).Float(), field.Index(2).Float()
		case field.Len() == 16 && field.Type().Elem().Kind() == reflect.Uint8:
			result.Kind = native.PacketFieldUUID
			reflect.Copy(reflect.ValueOf(result.UUID[:]), field)
		default:
			return jsonPacketField(field)
		}
	default:
		return jsonPacketField(field)
	}
	return result, true
}

func jsonPacketField(field reflect.Value) (native.PacketFieldValue, bool) {
	data, err := json.Marshal(field.Interface())
	return native.PacketFieldValue{Kind: native.PacketFieldJSON, Data: data}, err == nil
}

func writePacketField(field reflect.Value, value native.PacketFieldValue) bool {
	switch value.Kind {
	case native.PacketFieldBool:
		if field.Kind() != reflect.Bool || value.Unsigned > 1 {
			return false
		}
		field.SetBool(value.Unsigned != 0)
	case native.PacketFieldSigned:
		if field.Kind() < reflect.Int || field.Kind() > reflect.Int64 || field.OverflowInt(value.Signed) {
			return false
		}
		field.SetInt(value.Signed)
	case native.PacketFieldUnsigned:
		if field.Kind() < reflect.Uint || field.Kind() > reflect.Uint64 || field.OverflowUint(value.Unsigned) {
			return false
		}
		field.SetUint(value.Unsigned)
	case native.PacketFieldFloat:
		if (field.Kind() != reflect.Float32 && field.Kind() != reflect.Float64) || field.OverflowFloat(value.Number) {
			return false
		}
		field.SetFloat(value.Number)
	case native.PacketFieldString:
		if field.Kind() != reflect.String {
			return false
		}
		field.SetString(string(value.Data))
	case native.PacketFieldBytes:
		if field.Kind() != reflect.Slice || field.Type().Elem().Kind() != reflect.Uint8 {
			return false
		}
		field.SetBytes(append([]byte(nil), value.Data...))
	case native.PacketFieldVec2, native.PacketFieldVec3:
		want := 2
		if value.Kind == native.PacketFieldVec3 {
			want = 3
		}
		if field.Kind() != reflect.Array || field.Len() != want || field.Type().Elem().Kind() != reflect.Float32 {
			return false
		}
		values := [...]float64{value.X, value.Y, value.Z}
		for index := range want {
			field.Index(index).SetFloat(values[index])
		}
	case native.PacketFieldUUID:
		if field.Kind() != reflect.Array || field.Len() != 16 || field.Type().Elem().Kind() != reflect.Uint8 {
			return false
		}
		reflect.Copy(field, reflect.ValueOf(value.UUID[:]))
	default:
		return false
	}
	return true
}
