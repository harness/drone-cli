package mocks

import "github.com/drone/drone-go/drone"
import "github.com/stretchr/testify/mock"

import "io"

type Client struct {
	mock.Mock
}

func (_m *Client) Self() (*drone.User, error) {
	ret := _m.Called()

	var r0 *drone.User
	if rf, ok := ret.Get(0).(func() *drone.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*drone.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Client) User(_a0 string) (*drone.User, error) {
	ret := _m.Called(_a0)

	var r0 *drone.User
	if rf, ok := ret.Get(0).(func(string) *drone.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*drone.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Client) UserList() ([]*drone.User, error) {
	ret := _m.Called()

	var r0 []*drone.User
	if rf, ok := ret.Get(0).(func() []*drone.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*drone.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Client) UserPost(_a0 *drone.User) (*drone.User, error) {
	ret := _m.Called(_a0)

	var r0 *drone.User
	if rf, ok := ret.Get(0).(func(*drone.User) *drone.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*drone.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*drone.User) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Client) UserPatch(_a0 *drone.User) (*drone.User, error) {
	ret := _m.Called(_a0)

	var r0 *drone.User
	if rf, ok := ret.Get(0).(func(*drone.User) *drone.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*drone.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*drone.User) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Client) UserDel(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *Client) UserFeed() ([]*drone.Activity, error) {
	ret := _m.Called()

	var r0 []*drone.Activity
	if rf, ok := ret.Get(0).(func() []*drone.Activity); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*drone.Activity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Client) Repo(_a0 string, _a1 string) (*drone.Repo, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *drone.Repo
	if rf, ok := ret.Get(0).(func(string, string) *drone.Repo); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*drone.Repo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Client) RepoList() ([]*drone.Repo, error) {
	ret := _m.Called()

	var r0 []*drone.Repo
	if rf, ok := ret.Get(0).(func() []*drone.Repo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*drone.Repo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Client) RepoPost(_a0 string, _a1 string) (*drone.Repo, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *drone.Repo
	if rf, ok := ret.Get(0).(func(string, string) *drone.Repo); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*drone.Repo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Client) RepoPatch(_a0 *drone.Repo) (*drone.Repo, error) {
	ret := _m.Called(_a0)

	var r0 *drone.Repo
	if rf, ok := ret.Get(0).(func(*drone.Repo) *drone.Repo); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*drone.Repo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*drone.Repo) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Client) RepoDel(_a0 string, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *Client) RepoKey(_a0 string, _a1 string) (*drone.Key, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *drone.Key
	if rf, ok := ret.Get(0).(func(string, string) *drone.Key); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*drone.Key)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Client) Build(_a0 string, _a1 string, _a2 int) (*drone.Build, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *drone.Build
	if rf, ok := ret.Get(0).(func(string, string, int) *drone.Build); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*drone.Build)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, int) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Client) BuildList(_a0 string, _a1 string) ([]*drone.Build, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*drone.Build
	if rf, ok := ret.Get(0).(func(string, string) []*drone.Build); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*drone.Build)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Client) BuildStart(_a0 string, _a1 string, _a2 int) (*drone.Build, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *drone.Build
	if rf, ok := ret.Get(0).(func(string, string, int) *drone.Build); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*drone.Build)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, int) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Client) BuildStop(_a0 string, _a1 string, _a2 int, _a3 int) error {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, int, int) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *Client) BuildLogs(_a0 string, _a1 string, _a2 int, _a3 int) (io.ReadCloser, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 io.ReadCloser
	if rf, ok := ret.Get(0).(func(string, string, int, int) io.ReadCloser); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(io.ReadCloser)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, int, int) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Client) Node(_a0 int64) (*drone.Node, error) {
	ret := _m.Called(_a0)

	var r0 *drone.Node
	if rf, ok := ret.Get(0).(func(int64) *drone.Node); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*drone.Node)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Client) NodeList() ([]*drone.Node, error) {
	ret := _m.Called()

	var r0 []*drone.Node
	if rf, ok := ret.Get(0).(func() []*drone.Node); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*drone.Node)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Client) NodePost(_a0 *drone.Node) (*drone.Node, error) {
	ret := _m.Called(_a0)

	var r0 *drone.Node
	if rf, ok := ret.Get(0).(func(*drone.Node) *drone.Node); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*drone.Node)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*drone.Node) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Client) NodeDel(_a0 int64) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
