package gomahawk

type tomahawkImpl struct {
	name string
	uuid string
	*controlConnection
}

func (t *tomahawkImpl) Name() string {
	return t.name
}

func (t *tomahawkImpl) UUID() string {
	return t.uuid
}

//func (t *tomahawkImpl) RequestStreamConnection(id string) (msg.StreamConnection, error) {
//	return t.sendStreamRequest(id), nil
//}

func (t *tomahawkImpl) RequestDBConnection(DBConnection) error {
	return t.sendDBSyncOffer()
}

func (t *tomahawkImpl) TriggerDBChanges() {
}

func newTomahawk(c *controlConnection) (*tomahawkImpl, error) {
	t := new(tomahawkImpl)
	t.controlConnection = c

	return t, nil
}
