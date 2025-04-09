package transactions

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/tdadadavid/block/pkg/toolkit"
)

var COINBASE_DATA = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

// Transaction represents the shape of transactions that occurs on the chain
type Transaction struct {
	Id      string      `json:"id"`
	Inputs  []TxnInput  `json:"vin"`
	Outputs []TxnOutput `json:"vout"`
}

// GetId returns transaction id
func (t *Transaction) GetId() string {
	return t.Id
}

// GetInputs returns transaction inputs
func (t *Transaction) GetInputs() []TxnInput {
	return t.Inputs
}

// GetOutputs returns transaction outputs
func (t *Transaction) GetOutputs() []TxnOutput {
	return t.Outputs
}

// NewCoinbase creates the first coin on the chain
//
// Parameters
//   - `data string`: The input to the transaction
//   - `to string`: The address where this should be //FIXME
//
// # Process
//
// Returns
//   - `txn *Transaction`: The new coinbase transactions.
func NewCoinbase(to, data string) (txn *Transaction) {
	if data == "" {
		data = fmt.Sprintf("Reward to %s", to)
	}

	txn = &Transaction{
		Id: "",
		Inputs: []TxnInput{
			{
				TxnId:           "",
				Output:          -1,
				ScriptSignature: data,
			},
		},
		Outputs: []TxnOutput{
			{
				Value:        100,
				ScriptPubKey: to,
			},
		},
	}

	return txn
}

// GenId generates id for a transaction
//
// Process
//   - Serializes the transaction
//   - Hash it then convert it to hexadecimal string
//   - Store hexcode in transaction id
func (t *Transaction) GenId() {
	bytez, err := t.Serialize()
	if err != nil {
		fmt.Println("err serializing transactions: " + err.Error())
		return
	}
	hash := sha256.Sum256(bytez)
	hex := hex.EncodeToString(hash[:])
	t.Id = hex
}

// IsCoinbase checks if the transaction is the first transaction
//
// Process
//   - Check if outputs is equal to 1,the first input transaction ID is an empty string
//     the first input output is -1
//
// Returns
//   - bool: True or false informing the caller whether it is coinbase transaction
func (t *Transaction) IsCoinbase() bool {
	return len(t.Outputs) == 1 && t.Inputs[0].TxnId == "" && t.Inputs[0].Output == -1
}

// Serialize converts a Transaction into a byte slice
//
// Process
//   - First serializes the transaction ID, then iterates through the `Outputs` & `Input` and serializes them
//
// Returns
//   - `val []byte`: The byte value of the current transaction
//   - `err error`: Any error that occurs during serialization
func (t *Transaction) Serialize() (val []byte, err error) {
	buf := new(bytes.Buffer)

	// Write Transaction ID
	if err = toolkit.SerializeString(buf, t.Id); err != nil {
		return val, err
	}

	// Serialize Inputs
	// Write  the number of inputs in for this transaction into the buffer
	inputCount := uint32(len(t.Inputs))
	if err := binary.Write(buf, binary.LittleEndian, inputCount); err != nil {
		return val, err
	}

	for _, input := range t.Inputs {
		// Write transaction ID
		if err := toolkit.SerializeString(buf, input.TxnId); err != nil {
			return val, err
		}

		// Write transaction input ouput
		if err := binary.Write(buf, binary.LittleEndian, input.Output); err != nil {
			return val, err
		}

		// Write ScriptSignature
		if err := toolkit.SerializeString(buf, input.ScriptSignature); err != nil {
			return val, err
		}
	}

	// Serialize Outputs
	outputCount := uint32(len(t.Outputs))
	// Write  the number of outputs in for this transaction into the buffer
	if err := binary.Write(buf, binary.LittleEndian, outputCount); err != nil {
		return val, err
	}
	for _, output := range t.Outputs {
		// Write ScriptPubKey
		if err := toolkit.SerializeString(buf, output.ScriptPubKey); err != nil {
			return val, err
		}
		// Write Value
		if err := binary.Write(buf, binary.LittleEndian, output.Value); err != nil {
			return val, err
		}
	}

	val = buf.Bytes()
	return val, err
}

// Deserialize converts a byte slice back into a Transaction
//
// Process
//   - First reads the transaction ID, then iterates through the `Outputs` & `Input` and deserilizes them
//     into the given transaction struct
//
// NOTE
//   - Something I learnt about binary serialization and deserialization is that you deserialize in the same order you serialized.
//
// Returns
//   - `err error`: Any error that occurs during deserialization
func (t *Transaction) Deserialize(data []byte) (err error) {
	buf := bytes.NewReader(data)

	// Deserialize Transaction ID
	id, err := toolkit.DeserializeString(buf)
	if err != nil {
		return err
	}
	t.Id = id

	// Deserialize Inputs
	var inputCount uint32
	if err = binary.Read(buf, binary.LittleEndian, &inputCount); err != nil {
		return err
	}
	t.Inputs = make([]TxnInput, inputCount)
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

		t.Inputs[i] = TxnInput{txnId, output, scriptSig}
	}

	// Deserialize Outputs
	var outputCount int32
	if err := binary.Read(buf, binary.LittleEndian, &outputCount); err != nil {
		return err
	}
	t.Outputs = make([]TxnOutput, outputCount)
	for i := int32(0); i < outputCount; i++ {
		scriptPubKey, err := toolkit.DeserializeString(buf)
		if err != nil {
			return err
		}

		var value int32
		if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
			return err
		}

		t.Outputs[i] = TxnOutput{int64(value), scriptPubKey}
	}

	return err
}
