package transactions

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/tdadadavid/block/pkg/toolkit"
)

type Transaction struct {
	id string
	inputs []TxnInput
	outputs []TxnOutput
}


func (t *Transaction) NewCoinbase(data, to string) (txn *Transaction) {
	if data == "" {
		data = fmt.Sprintf("Reward to %s", to)
	}

	txn = &Transaction {
		id: "",
		inputs: []TxnInput{
			{
				txnId: "",
				output: -1,
				scriptSignature: data,
			},
		},
		outputs: []TxnOutput{
			{
				value: 100,
				scriptPubKey: to,
			},
		},
	}

	return txn
}

// GenId generates id for a transaction
//
// Process
// 	- Serializes the transaction
//  - Hash it then convert it to hexadecimal string
//  - Store hexcode in transaction id
func (t *Transaction) GenId() {
	bytez, err := t.Serialize()
	if err != nil {
		fmt.Println("err serializing transactions: " + err.Error())
		return
	}
	hash := sha256.Sum256(bytez)
	hex := hex.EncodeToString(hash[:])
	t.id = hex
}

func (t *Transaction) IsCoinbase() bool {
	return len(t.outputs) == 1 && t.inputs[0].txnId == "" && t.inputs[0].output == -1
}


// Serialize converts a Transaction into a byte slice
func (t *Transaction) Serialize() (val []byte, err error) {
	buf := new(bytes.Buffer)

	// Serialize Transaction ID
	if err = toolkit.SerializeString(buf, t.id); err != nil {
		return val, err
	}

	// Serialize Inputs
	inputCount := uint32(len(t.inputs))
	if err := binary.Write(buf, binary.LittleEndian, inputCount); err != nil {
		return val, err
	}
	for _, input := range t.inputs {
		if err := toolkit.SerializeString(buf, input.txnId); err != nil {
			return val, err
		}
		if err := binary.Write(buf, binary.LittleEndian, input.output); err != nil {
			return val, err
		}
		if err := toolkit.SerializeString(buf, input.scriptSignature); err != nil {
			return val, err
		}
	}

	// Serialize Outputs
	outputCount := uint32(len(t.outputs))
	if err := binary.Write(buf, binary.LittleEndian, outputCount); err != nil {
		return val, err
	}
	for _, output := range t.outputs {
		if err := binary.Write(buf, binary.LittleEndian, output.value); err != nil {
			return val, err
		}
		if err := toolkit.SerializeString(buf, output.scriptPubKey); err != nil {
			return val, err
		}
	}

	val = buf.Bytes()
	return val, nil
}

// Deserialize converts a byte slice back into a Transaction
func (t *Transaction) Deserialize(data []byte) (err error) {
	buf := bytes.NewReader(data)

	// Deserialize Transaction ID
	id, err := toolkit.DeserializeString(buf)
	if err != nil {
		return err
	}
	t.id = id

	// Deserialize Inputs
	var inputCount uint32
	if err = binary.Read(buf, binary.LittleEndian, &inputCount); err != nil {
		return err
	}
	t.inputs = make([]TxnInput, inputCount)
	for i := uint32(0); i < inputCount; i++ {
		txnId, err := toolkit.DeserializeString(buf)
		if err != nil {
			return err
		}

		var output int32
		if err := binary.Read(buf, binary.LittleEndian, &output); err != nil {
			return err
		}

		scriptSig, err := toolkit.DeserializeString(buf)
		if err != nil {
			return err
		}

		t.inputs[i] = TxnInput{txnId, output, scriptSig}
	}

	// Deserialize Outputs
	var outputCount uint32
	if err := binary.Read(buf, binary.LittleEndian, &outputCount); err != nil {
		return err
	}
	t.outputs = make([]TxnOutput, outputCount)
	for i := uint32(0); i < outputCount; i++ {
		var value int
		if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
			return err
		}

		scriptPubKey, err := toolkit.DeserializeString(buf)
		if err != nil {
			return err
		}

		t.outputs[i] = TxnOutput{value, scriptPubKey}
	}

	return err
}