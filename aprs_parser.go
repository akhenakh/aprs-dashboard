package main

// #cgo pkg-config: --libs libfap
// #include <fap.h>
// #include <stdlib.h>
import "C"
import "unsafe"
import "errors"

type AprsMessage struct {
	Timestamp           int64   `json:"timestamp"`
	SourceCallsign      string  `json:"src_callsign"`
	DestinationCallsign string  `json:"dst_callsign"`
	Status              string  `json:"status"`
	Symbol              string  `json:"symbol"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
	Altitude            float64 `json:"altitude"`
	Speed               float64 `json:"speed"`
	Temperature         float64 `json:"temperature"`
	Humidity            float64 `json:"humidity"`
	RawMessage          string  `json:"raw_message"`
}

func parseAprsPacket(message string, isAX25 bool) (*AprsMessage, error) {
	message_cstring := C.CString(message)
	message_length := C.uint(len(message))

	C.fap_init()

	var ax25Flag int
	if isAX25 {
		ax25Flag = 1
	} else {
		ax25Flag = 0
	}
	packet := C.fap_parseaprs(message_cstring, message_length, C.short(ax25Flag))
	if packet.error_code != nil {
		return &AprsMessage{}, errors.New("Unable to parse APRS message")
	}

	parsedMsg := AprsMessage{
		SourceCallsign:      C.GoString(packet.src_callsign),
		DestinationCallsign: C.GoString(packet.dst_callsign),
	}
	parsedMsg.Latitude = float64(C.double(*packet.latitude))
	parsedMsg.Longitude = float64(C.double(*packet.longitude))
	parsedMsg.RawMessage = C.GoStringN(packet.body, C.int(packet.body_len))

	C.fap_free(packet)
	C.fap_cleanup()
	C.free(unsafe.Pointer(message_cstring))

	return &parsedMsg, nil
}