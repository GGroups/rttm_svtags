package tagmd

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	COMM "github.com/GGroups/rttm_login/comm"
	UR "github.com/GGroups/rttm_login/user"
	"github.com/go-kit/kit/endpoint"
)

const (
	ACL_CODE = "8001"
)

type ITagm interface {
	GetOneTag(usr UR.Usr, id int, one *TagM) error
	GetTagList(usr UR.Usr, all *[]TagM) error
	SetOneTag(usr UR.Usr, one *TagM) error
	AddOneTag(usr UR.Usr, one *TagM) error
}

func (s TagM) GetOneTag(usr UR.Usr, id int, one *TagM) error {
	return nil
}

func (s TagM) GetTagList(usr UR.Usr, all *[]TagM) error {
	return GetAllTags(all)
}

func (s TagM) AddOneTag(usr UR.Usr, one *TagM) error {
	doit := hasAccessRole(usr)
	if doit {
		return CreatTag(one)
	}
	return errors.New("无执行权限")
}

func (s TagM) SetOneTag(usr UR.Usr, one *TagM) error {
	doit := hasAccessRole(usr)
	if doit {
		return SetTag(one)
	}
	return errors.New("无执行权限")
}

func hasAccessRole(usr UR.Usr) bool {
	roles := strings.Split(usr.Roles, ",")
	doit := false
	for _, r := range roles {
		if strings.TrimSpace(r) == ACL_CODE {
			doit = true
		}
	}
	return doit
}

/*===============================================================================
EndPoint
====================*/

func MakeEndPointGetAllTag(sv ITagm) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r, ok := request.(COMM.RequestWarp)
		if !ok {
			return COMM.ErrBody(), nil
		}

		var tgs []TagM
		err = sv.GetTagList(r.Usr, &tgs)
		if err != nil {
			return TagM{}, COMM.RepErr(http.StatusBadRequest, err)
		} else {
			return tgs, nil
		}
	}
}

type RespGetOneTag struct {
	Id int `json:"tid"`
}

func MakeEndPointGetOneTag(sv ITagm) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r, ok := request.(RespGetOneTag)
		if !ok {
			return RespGetOneTag{}, nil
		}

		var tg TagM
		var u UR.Usr
		err = sv.GetOneTag(u, r.Id, &tg)
		if err != nil {
			return TagM{}, COMM.RepErr(http.StatusBadRequest, err)
		} else {
			return tg, nil
		}
	}
}

type RespAddOneTag struct {
	Tag TagM `json:"tag"`
}

func MakeEndPointAddOneTag(sv ITagm) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r, ok := request.(COMM.RequestWarp)
		if !ok {
			return COMM.ErrBody(), nil
		}
		resp, ok := r.Resp.(RespAddOneTag)
		if !ok {
			return COMM.ErrBody(), nil
		}

		err = sv.AddOneTag(r.Usr, &resp.Tag)
		if err != nil {
			return COMM.ErrBody(), COMM.RepErr(http.StatusBadRequest, err)
		} else {
			return COMM.OkBody(), nil
		}
	}
}

func MakeEndPointSetOneTag(sv ITagm) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r, ok := request.(COMM.RequestWarp)
		if !ok {
			return COMM.ErrBody(), nil
		}
		resp, ok := r.Resp.(RespAddOneTag)
		if !ok {
			return COMM.ErrBody(), nil
		}

		err = sv.SetOneTag(r.Usr, &resp.Tag)
		if err != nil {
			return COMM.ErrBody(), COMM.RepErr(http.StatusBadRequest, err)
		} else {
			return COMM.OkBody(), nil
		}
	}
}

/*===============================================================================
Transport
====================*/

func DecodeRequestAddOneTag(c context.Context, request *http.Request) (interface{}, error) {
	if request.Method != "POST" {
		return nil, errors.New("#must POST")
	}
	u := UR.Usr{}
	code, err := COMM.GetUserFromToken(request, &u)
	if err != nil {
		return nil, COMM.RepErr(code, err)
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, COMM.RepErr(http.StatusBadRequest, err)
	}
	var obj RespAddOneTag
	err = json.Unmarshal(body, &obj)
	if err != nil {
		return nil, COMM.RepErr(http.StatusBadRequest, err)
	}

	w := COMM.RequestWarp{Usr: u, Resp: obj}
	return w, nil
}

func DecodeRequestEmptyReq(c context.Context, request *http.Request) (interface{}, error) {
	if request.Method != "POST" {
		return nil, errors.New("#must POST")
	}
	u := UR.Usr{}
	code, err := COMM.GetUserFromToken(request, &u)
	if err != nil {
		return nil, COMM.RepErr(code, err)
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, COMM.RepErr(http.StatusBadRequest, err)
	}
	var obj COMM.EmptyReqRep
	err = json.Unmarshal(body, &obj)
	if err != nil {
		return nil, COMM.RepErr(http.StatusBadRequest, err)
	}

	w := COMM.RequestWarp{Usr: u, Resp: obj}
	return w, nil
}
