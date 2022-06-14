package user

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRpcServer(t *testing.T) {
	rpcHandler := &RPCHandler{handler: &MockHandler{users: make(map[uuid.UUID]Model)}}

	request := Request{
		Nickname: "test",
	}
	reply := 0

	err := rpcHandler.Add(&request, &reply)
	assert.NoError(t, err)

	findArgs := RPCFindArgs{}
	resUsers := []Model{}

	err = rpcHandler.FindAllWithPaginationAndFilter(&findArgs, &resUsers)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(resUsers))

	requestUpdate := Request{
		ID:       resUsers[0].ID.String(),
		Nickname: "test2",
	}

	err = rpcHandler.Update(&requestUpdate, &reply)
	assert.NoError(t, err)

	resUsers = []Model{}
	err = rpcHandler.FindAllWithPaginationAndFilter(&findArgs, &resUsers)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(resUsers))
	assert.Equal(t, requestUpdate.Nickname, resUsers[0].Nickname)

	err = rpcHandler.Delete(&RPCRemoveArgs{ID: resUsers[0].ID.String()}, &reply)
	assert.NoError(t, err)

	resUsers = []Model{}
	err = rpcHandler.FindAllWithPaginationAndFilter(&findArgs, &resUsers)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(resUsers))
}
