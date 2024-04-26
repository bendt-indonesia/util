package util

func NilBool(value bool) *bool {
	return &value
}

func NilString(value string, allowEmptyValue bool) *string {
	if !allowEmptyValue && value == "" {
		return nil
	}
	return &value
}

func NilInt(value int, allowEmptyValue bool) *int {
	if !allowEmptyValue && value == 0 {
		return nil
	}
	return &value
}

func NilInt8(value int8, allowEmptyValue bool) *int8 {
	if !allowEmptyValue && value == 0 {
		return nil
	}
	return &value
}

func NilInt16(value int16, allowEmptyValue bool) *int16 {
	if !allowEmptyValue && value == 0 {
		return nil
	}
	return &value
}

func NilInt32(value int32, allowEmptyValue bool) *int32 {
	if !allowEmptyValue && value == 0 {
		return nil
	}
	return &value
}

func NilInt64(value int64, allowEmptyValue bool) *int64 {
	if !allowEmptyValue && value == 0 {
		return nil
	}
	return &value
}

func NilUint(value uint, allowEmptyValue bool) *uint {
	if !allowEmptyValue && value == 0 {
		return nil
	}
	return &value
}

func NilUint8(value uint8, allowEmptyValue bool) *uint8 {
	if !allowEmptyValue && value == 0 {
		return nil
	}
	return &value
}

func NilUint16(value uint16, allowEmptyValue bool) *uint16 {
	if !allowEmptyValue && value == 0 {
		return nil
	}
	return &value
}

func NilUint32(value uint32, allowEmptyValue bool) *uint32 {
	if !allowEmptyValue && value == 0 {
		return nil
	}
	return &value
}

func NilUint64(value uint64, allowEmptyValue bool) *uint64 {
	if !allowEmptyValue && value == 0 {
		return nil
	}
	return &value
}

func NilFloat32(value float32, allowEmptyValue bool) *float32 {
	if !allowEmptyValue && value == 0 {
		return nil
	}
	return &value
}

func NilFloat64(value float64, allowEmptyValue bool) *float64 {
	if !allowEmptyValue && value == 0 {
		return nil
	}
	return &value
}
