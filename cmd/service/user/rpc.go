package user

import (
	"net"
	"net/http"
	"net/rpc"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type RPCFindArgs struct {
	Offset   int
	Count    int
	Nickname string
}

type RPCRemoveArgs struct {
	ID string
}

type RPCHandler struct {
	handler UniqueHandler
}

func (r *RPCHandler) FindAllWithPaginationAndFilter(args *RPCFindArgs, reply *[]Model) error {
	users, err := r.handler.FindAllWithPaginationAndFilter(args.Offset, args.Count, map[string]string{
		FilterUserNickname: args.Nickname,
	})
	if err != nil {
		log.Error().Err(err).Msg("can't find users")
		return err
	}

	*reply = append(*reply, users...)
	_ = reply

	return nil
}

func (r *RPCHandler) Add(request *Request, reply *int) error {
	err := r.handler.Add(request.FirstName, request.LastName, request.Nickname, request.Password, request.Email, request.Country)
	if err != nil {
		log.Error().Err(err).Msg("can't add a new user")
		return err
	}

	return nil
}

func (r *RPCHandler) Update(request *Request, reply *int) error {
	id, err := uuid.Parse(request.ID)
	if err != nil {
		log.Error().Err(err).Msg("can't parse id as uuid")
		return err
	}

	err = r.handler.Update(id, request.FirstName, request.LastName, request.Nickname, request.Country)
	if err != nil {
		log.Error().Err(err).Msg("can't update existing user")
		return err
	}

	return nil
}

func (r *RPCHandler) Delete(args *RPCRemoveArgs, reply *int) error {
	id, err := uuid.Parse(args.ID)
	if err != nil {
		log.Error().Err(err).Msg("can't parse id as uuid")
		return err
	}

	err = r.handler.Remove(id)
	if err != nil {
		log.Error().Err(err).Msg("can't update existing user")
		return err
	}

	return nil
}

// RpcServer delcares RPC server
func RpcServer(port string, handler UniqueHandler) error {
	err := rpc.Register(&RPCHandler{handler: handler})
	if err != nil {
		return err
	}

	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	log.Debug().Msg("User rpc server is starting")

	err = http.Serve(listener, nil)
	if err != nil {
		return err
	}

	return nil
}
