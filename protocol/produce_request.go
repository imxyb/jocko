package protocol

type Data struct {
	Partition int32
	RecordSet []byte
}

type TopicData struct {
	Topic string
	Data  []*Data
}

type ProduceRequest struct {
	// unsupported: Acks int16
	Timeout   int32
	TopicData []*TopicData
}

func (r *ProduceRequest) Encode(e PacketEncoder) error {
	e.PutInt32(r.Timeout)
	e.PutArrayLength(len(r.TopicData))
	for _, td := range r.TopicData {
		e.PutString(td.Topic)
		e.PutArrayLength(len(td.Data))
		for _, d := range td.Data {
			e.PutInt32(d.Partition)
			e.PutBytes(d.RecordSet)
		}
	}
	return nil
}

func (r *ProduceRequest) Decode(d PacketDecoder) error {
	var err error
	r.Timeout, err = d.Int32()
	if err != nil {
		return err
	}
	tdlen, err := d.ArrayLength()
	r.TopicData = make([]*TopicData, tdlen)
	for i := range r.TopicData {
		td := new(TopicData)
		r.TopicData[i] = td
		td.Topic, err = d.String()
		if err != nil {
			return err
		}
		dlen, err := d.ArrayLength()
		if err != nil {
			return err
		}
		td.Data = make([]*Data, dlen)
		for j := range td.Data {
			data := new(Data)
			td.Data[j] = data
			data.Partition, err = d.Int32()
			if err != nil {
				return err
			}
			data.RecordSet, err = d.Bytes()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *ProduceRequest) Key() int16 {
	return 0
}

func (r *ProduceRequest) Version() int16 {
	return 2
}
