package controller

import (
	"errors"
	"github.com/sirupsen/logrus"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

const (
	sep = "/"
)

var (
	ErrLocked = errors.New("acquire Lock failed")
	ErrBadKey = errors.New("invalid key")
)

type Lock struct {
	c   *zk.Conn
	acl []zk.ACL
}

func NewLock(addr string) (*Lock, error) {
	l := &Lock{}

	c, _, err := zk.Connect([]string{addr}, time.Second)
	if err != nil {
		return nil, err
	}

	l.c = c
	l.acl = zk.WorldACL(zk.PermAll)

	return l, nil
}

func (l *Lock) getRealPath(path string) string {
	path = "/jcs" + path
	return path
}

func (l *Lock) Lock(path string) error {
	// 检查前缀格式
	if len(path) == 0 {
		return ErrBadKey
	} else if path[0] != '/' {
		path = "/" + path
	}
	path = "/jcs" + path
	// 检查path格式
	paths := strings.Split(strings.Trim(path, sep), sep)
	if len(paths) == 0 {
		return ErrBadKey
	}

	// 生成子node数组
	nodes := []string{sep + paths[0]}
	for i := 1; i < len(paths); i++ {
		nodes = append(nodes, nodes[i-1]+sep+paths[i])
	}

	// 查询未创建的node
	var idx int
	for idx = 0; idx < len(nodes); idx++ {
		exist, stat, err := l.c.Exists(nodes[idx])
		if err != nil {
			return err
		}
		if exist && stat.NumChildren == 0 {
			return ErrLocked
		}
		if !exist {
			break
		}
	}
	nodes = nodes[idx:]

	// 全部存在意味着已被锁定
	if len(nodes) == 0 {
		return ErrLocked
	}

	// 创建请求
	var ops []interface{}
	for _, p := range nodes {
		ops = append(ops, &zk.CreateRequest{p, []byte{}, l.acl, 0})
	}

	// 发送多个子请求
	_, err := l.c.Multi(ops...)
	if err != nil {
		if err == zk.ErrNodeExists {
			return ErrLocked
		}
		return err
	}

	return nil
}

func (l *Lock) UnLockAll(path string) {
	realPath := l.getRealPath(path)
	_, _, err := l.c.Get(realPath)
	if err != nil {
		return
	}

	children, _, err := l.c.Children(realPath)

	if len(children) != 0 {
		for _, child := range children {
			l.UnLockAll(path + sep + child)
		}
	}
	err = l.UnLock(path)
	if err != nil {
		logrus.Errorf("unlock err: %v", err)
	}
	return
}

func (l *Lock) UnLock(path string) error {
	// 检查前缀格式
	if len(path) == 0 {
		return ErrBadKey
	} else if path[0] != '/' {
		path = "/" + path
	}
	path = "/jcs" + path
	// 检查path格式
	paths := strings.Split(strings.Trim(path, sep), sep)
	if len(paths) == 0 {
		logrus.Infof("unlock fail")
		return ErrBadKey
	}

	// 生成子node数组
	nodes := []string{sep + paths[0]}
	for i := 1; i < len(paths); i++ {
		nodes = append(nodes, nodes[i-1]+sep+paths[i])
	}
	for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	}

	// 查询未创建的node
	var idx int
	for idx = 0; idx < len(nodes); idx++ {
		_, stat, err := l.c.Exists(nodes[idx])
		if err != nil {
			logrus.Infof("unlock fail: %v", err)
			return err
		}
		if stat.NumChildren > 1 {
			break
		}
	}
	nodes = nodes[:idx]

	// 创建请求
	var ops []interface{}
	for _, n := range nodes {
		ops = append(ops, &zk.DeleteRequest{n, -1})
	}

	// 发送多个子请求
	_, err := l.c.Multi(ops...)
	if err != nil {
		logrus.Warnf("unlock %v fail: %v", path, err)
	} else {
		logrus.Infof("unlock success")
	}
	return nil
}
