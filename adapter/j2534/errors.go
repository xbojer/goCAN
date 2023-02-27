package j2534

import (
	"errors"
	"fmt"
)

var (
	ErrNotSupported        = errors.New("device cannot support requested functionality mandated in J2534. Device is not fully SAE J2534 compliant")
	ErrInvalidChannelID    = errors.New("invalid ChannelID value")
	ErrInvalidProtocolID   = errors.New("invalid or unsupported ProtocolID, or there is a resource conflict (i.e. trying to connect to multiple mutually exclusive protocols such as J1850PWM and J1850VPW, or CAN and SCI, etc.)")
	ErrNullParameter       = errors.New("NULL pointer supplied where a valid pointer is required")
	ErrInvalidIoctlValue   = errors.New("invalid value for Ioctl parameter")
	ErrInvalidFlags        = errors.New("invalid flag values")
	ErrFailed              = errors.New("undefined error, use PassThruGetLastError() for text description")
	ErrDeviceNotConnected  = errors.New("unable to communicate with device")
	ErrTimeout             = errors.New("read or write timeout")
	ErrInvalidMsg          = errors.New("invalid message structure pointed to by pMsg")
	ErrInvalidTimeInterval = errors.New("invalid TimeInterval value")
	ErrExceededLimit       = errors.New("exceeded maximum number of message IDs or allocated space")
	ErrInvalidMsgID        = errors.New("invalid MsgID value")
	ErrDeviceInUse         = errors.New("device is currently open")
	ErrInvalidIoctlID      = errors.New("invalid IoctlID value")
	ErrBufferEmpty         = errors.New("protocol message buffer empty, no messages available to read")
	ErrBufferFull          = errors.New("protocol message buffer full. All the messages specified may not have been transmitted")
	ErrBufferOverflow      = errors.New("indicates a buffer overflow occurred and messages were lost")
	ErrPinInvalid          = errors.New("invalid pin number, pin number already in use, or voltage already applied to a different pin")
	ErrChannelInUse        = errors.New("channel number is currently connected")
	ErrMsgProtocolID       = errors.New("protocol type in the message does not match the protocol associated with the Channel ID")
	ErrInvalidFilterID     = errors.New("invalid Filter ID value")
	ErrNoFlowControl       = errors.New("no flow control filter set or matched (for ProtocolID ISO15765 only)")
	ErrNotUnique           = errors.New("a CAN ID in pPatternMsg or pFlowControlMsg matches either ID in an existing FLOW_CONTROL_FILTER")
	ErrInvalidBaudrate     = errors.New("the desired baud rate cannot be achieved within the tolerance specified in SAE J2534-1 Section 6.5")
	ErrInvalidDeviceID     = errors.New("device ID invalid")
	ErrUnknown             = errors.New("unknown error")
)

func CheckError(ret uintptr) error {
	switch int32(ret) {
	case STATUS_NOERROR:
		//return errors.New("Function call successful")
		return nil
	case ERR_NOT_SUPPORTED:
		return ErrNotSupported
	case ERR_INVALID_CHANNEL_ID:
		return ErrInvalidChannelID
	case ERR_INVALID_PROTOCOL_ID:
		return ErrInvalidProtocolID
	case ERR_NULL_PARAMETER:
		return ErrNullParameter
	case ERR_INVALID_IOCTL_VALUE:
		return ErrInvalidIoctlValue
	case ERR_INVALID_FLAGS:
		return ErrInvalidFlags
	case ERR_FAILED:
		return ErrFailed
	case ERR_DEVICE_NOT_CONNECTED:
		return ErrDeviceNotConnected
	case ERR_TIMEOUT:
		return ErrTimeout
	case ERR_INVALID_MSG:
		return ErrInvalidMsg
	case ERR_INVALID_TIME_INTERVAL:
		return ErrInvalidTimeInterval
	case ERR_EXCEEDED_LIMIT:
		return ErrExceededLimit
	case ERR_INVALID_MSG_ID:
		return ErrInvalidMsgID
	case ERR_DEVICE_IN_USE:
		return ErrDeviceInUse
	case ERR_INVALID_IOCTL_ID:
		return ErrInvalidIoctlID
	case ERR_BUFFER_EMPTY:
		return ErrBufferEmpty
	case ERR_BUFFER_FULL:
		return ErrBufferFull
	case ERR_BUFFER_OVERFLOW:
		return ErrBufferOverflow
	case ERR_PIN_INVALID:
		return ErrPinInvalid
	case ERR_CHANNEL_IN_USE:
		return ErrChannelInUse
	case ERR_MSG_PROTOCOL_ID:
		return ErrMsgProtocolID
	case ERR_INVALID_FILTER_ID:
		return ErrInvalidFilterID
	case ERR_NO_FLOW_CONTROL:
		return ErrNoFlowControl
	case ERR_NOT_UNIQUE:
		return ErrNotUnique
	case ERR_INVALID_BAUDRATE:
		return ErrInvalidBaudrate
	case ERR_INVALID_DEVICE_ID:
		return ErrInvalidDeviceID
	default:
		return fmt.Errorf("unknown error: %d", ret)
	}
}

/*
func errorHander(ret int32) {
	switch ret {
	case 0:
		fmt.Println("STATUS_NOERROR")
	case 0x01:
		fmt.Println("ERR_NOT_SUPPORTED")
	case 0x02:
		fmt.Println("ERR_INVALID_CHANNEL_ID")
	case 0x03:
		fmt.Println("ERR_INVALID_PROTOCOL_ID")
	case 0x04:
		fmt.Println("ERR_NULL_PARAMETER")
	case 0x05:
		fmt.Println("ERR_INVALID_IOCTL_VALUE")
	case 0x06:
		fmt.Println("ERR_INVALID_FLAGS")
	case 0x07:
		fmt.Println("ERR_FAILED")
	case 0x08:
		fmt.Println("ERR_DEVICE_NOT_CONNECTED")
	case 0x09:
		fmt.Println("ERR_TIMEOUT")
	case 0x0A:
		fmt.Println("ERR_INVALID_MSG")
	case 0x0B:
		fmt.Println("ERR_INVALID_TIME_INTERVAL")
	case 0x0C:
		fmt.Println("ERR_EXCEEDED_LIMIT")
	case 0x0D:
		fmt.Println("ERR_INVALID_MSG_ID")
	case 0x0E:
		fmt.Println("ERR_DEVICE_IN_USE")
	case 0x0F:
		fmt.Println("ERR_INVALID_IOCTL_ID")
	case 0x10:
		fmt.Println("ERR_BUFFER_EMPTY")
	case 0x11:
		fmt.Println("ERR_BUFFER_FULL")
	case 0x12:
		fmt.Println("ERR_BUFFER_OVERFLOW")
	case 0x13:
		fmt.Println("ERR_PIN_INVALID")
	case 0x14:
		fmt.Println("ERR_CHANNEL_IN_USE")
	case 0x15:
		fmt.Println("ERR_MSG_PROTOCOL_ID")
	case 0x16:
		fmt.Println("ERR_INVALID_FILTER_ID")
	case 0x17:
		fmt.Println("ERR_NO_FLOW_CONTROL")
	case 0x18:
		fmt.Println("ERR_NOT_UNIQUE")
	case 0x19:
		fmt.Println("ERR_INVALID_BAUDRATE")
	case 0x1A:
		fmt.Println("ERR_INVALID_DEVICE_ID")
	}
}
*/
