package parser

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type SQSConfig struct {
	QueueURL          string                       `toml:"queue_url"`
	DelaySeconds      int32                        `toml:"delay_seconds"`
	MessageAttributes map[string]SQSMessageAttrVal `toml:"message_attributes"`
}

type SQSMessageAttrVal struct {
	DataType    string `toml:"DataType"`
	StringValue string `toml:"StringValue"`
	BinaryValue []byte `toml:"BinaryValue"`
}

type SQSAction struct {
	Kind string    `toml:"kind"`
	SQS  SQSConfig `toml:"sqs"`
}

func TestUnmarshalComplex(t *testing.T) {
	data := []byte(`
kind = "sqs.send_message"
[sqs]
queue_url = "https://sqs.us-east-1.amazonaws.com/123456789012/MyQueue"
delay_seconds = 10
[sqs.message_attributes.Attribute1]
DataType = "String"
StringValue = "Value1"
`)
	var action SQSAction
	err := Unmarshal(data, &action)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if action.Kind != "sqs.send_message" {
		t.Errorf("Expected Kind 'sqs.send_message', got '%s'", action.Kind)
	}
	if action.SQS.QueueURL != "https://sqs.us-east-1.amazonaws.com/123456789012/MyQueue" {
		t.Errorf("Unexpected QueueURL: %s", action.SQS.QueueURL)
	}
	if action.SQS.DelaySeconds != 10 {
		t.Errorf("Expected DelaySeconds 10, got %d", action.SQS.DelaySeconds)
	}

	attr, ok := action.SQS.MessageAttributes["Attribute1"]
	if !ok {
		t.Fatal("MessageAttribute 'Attribute1' not found")
	}
	if attr.DataType != "String" || attr.StringValue != "Value1" {
		t.Errorf("Unexpected attribute values: %+v", attr)
	}
}

func TestMarshalComplex(t *testing.T) {
	action := SQSAction{
		Kind: "sqs.send_message",
		SQS: SQSConfig{
			QueueURL: "https://example.com",
			MessageAttributes: map[string]SQSMessageAttrVal{
				"Attr1": {
					DataType:    "String",
					StringValue: "Val1",
					BinaryValue: []byte{}, // Initialize to empty slice to match Unmarshal behavior if needed
				},
			},
		},
	}

	data, err := Marshal(action)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var action2 SQSAction
	err = Unmarshal(data, &action2)
	if err != nil {
		t.Fatalf("Unmarshal of marshaled data failed: %v", err)
	}

	if !reflect.DeepEqual(action, action2) {
		t.Errorf("Marshal result mismatch\nExpected: %s\nGot:      %s", fmt.Sprintf("%#v", action), fmt.Sprintf("%#v", action2))
	}
}

func TestJSONCompatibility(t *testing.T) {
	type Mixed struct {
		Name string `toml:"name" json:"name"`
	}

	data := []byte(`name = "test"`)
	var m Mixed
	err := Unmarshal(data, &m)
	if err != nil {
		t.Fatal(err)
	}

	if m.Name != "test" {
		t.Errorf("Expected 'test', got '%s'", m.Name)
	}

	jb, _ := json.Marshal(m)
	if string(jb) != `{"name":"test"}` {
		t.Errorf("JSON output mismatch: %s", string(jb))
	}
}
