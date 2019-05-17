package ccp

func SendErrMsg(protocol *CcpProtocol, errMessage string) error {
	if err := protocol.Send(Message{
		RESPONSE_ERROR,
		[]byte{},
		"",
	}); err != nil {
		return err
	} else {
		return nil
	}
}
