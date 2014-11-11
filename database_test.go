package main

import (
	"testing"
)

func Test_NewDatabase(t *testing.T) {
	db := NewDatabase(":6379", "", "")
	defer db.Close()
	err := db.Ping()
	if err != nil {
		t.Error("Unable to ping Redis", err)
	}
}

func Test_RecordMessage(t *testing.T) {
	db := NewDatabase(":6379", "", "")
	defer db.Close()
	cleanup(db)
	defer cleanup(db)

	var err error
	var length int64

	// verify item is not on list
	length, err = db.NumberOfMessagesForCallsign("foo")
	if 0 != length {
		t.Error("List length should be one", length)
	}

	// verify item is not on list
	length, err = db.NumberOfCallsigns()
	if 0 != length {
		t.Error("List length should be one", length)
	}

	// push item onto list
	message := &AprsMessage{}
	err = db.RecordMessage("foo", message)
	if err != nil {
		t.Error("Error while LPUSHing", err)
	}

	// verify item is on list
	length, err = db.NumberOfMessagesForCallsign("foo")
	if 1 != length {
		t.Error("List length should be one", length)
	}

	// verify item is on list
	length, err = db.NumberOfCallsigns()
	if 1 != length {
		t.Error("List length should be one", length)
	}
}

func Benchmark_RetrieveMostRecentEntriesForCallsign(b *testing.B) {
	db := NewDatabase(":6379", "", "")
	defer db.Close()
	cleanup(db)
	defer cleanup(db)

	p := NewParser()
	defer p.Finish()
	msg, _ := p.parseAprsPacket("K7SSW>APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5", false)

	var i int
	for i = 0; i < 10000; i++ {
		db.RecordMessage("foo", msg)
	}

	b.ResetTimer()
	for i = 0; i < b.N; i++ {
		db.GetRecordsForCallsign("foo", 1)
	}
}

func Benchmark_RetrieveMiddleEntriesForCallsign(b *testing.B) {
	db := NewDatabase(":6379", "", "")
	defer db.Close()
	cleanup(db)
	defer cleanup(db)

	p := NewParser()
	defer p.Finish()
	msg, _ := p.parseAprsPacket("K7SSW>APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5", false)

	var i int
	for i = 0; i < 10000; i++ {
		db.RecordMessage("foo", msg)
	}

	b.ResetTimer()
	for i = 0; i < b.N; i++ {
		db.GetRecordsForCallsign("foo", 500)
	}
}

func cleanup(db *Database) {
	db.Delete("callsign.foo")
	db.Delete("callsigns.set")
}
