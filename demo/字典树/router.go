package router

import (
	"errors"
	"strings"
	"sync"
)

// 路由匹配，支持 * 和 **
// *代表匹配一个节点
// **代表匹配全节点
// **必须为path结尾

const (
	OneStar = "*"
	TwoStar = "**"
)

//路由器
type Router struct {
	rawPath     map[string]interface{} //原始路径
	dirTree     *dirNode               //目录树
	routerMutex *sync.RWMutex          //路由锁
}

//目录节点
type dirNode struct {
	name              string              //名称
	value             interface{}         //数据
	isPathEnd         bool                //是否为路径结尾
	subDir            map[string]*dirNode //子目录
	subDirHaveOneStar bool                //子目录有 *
	subDirHaveTwoStar bool                //子目录有**
	parentDir         *dirNode            //父目录
	referenceNumber   int                 //引用计数，经过当前节点的路径个数
}

//新建目录节点
func newDirNode(name string, value interface{}, isPathEnd bool, parentDir *dirNode) *dirNode {
	return &dirNode{
		name:              name,
		value:             value,
		isPathEnd:         isPathEnd,
		subDir:            make(map[string]*dirNode),
		parentDir:         parentDir,
		subDirHaveOneStar: false,
		subDirHaveTwoStar: false,
		referenceNumber:   0,
	}
}

//新建路由
func New() *Router {
	return &Router{
		rawPath:     make(map[string]interface{}),
		dirTree:     newDirNode("/", nil, false, nil),
		routerMutex: new(sync.RWMutex),
	}
}

//切割路径
func (r *Router) splitPath(path string) ([]string, error) {
	dirs := make([]string, 0, 8)

	s, i := 0, 0
	for ; i < len(path); i++ {
		if path[i] == '/' {
			if s == i {
				s++
			} else {
				dirs = append(dirs, path[s:i])
				s = i + 1
			}
		}
	}
	if s != i {
		dirs = append(dirs, path[s:i])
	}

	return dirs, nil
}

//格式化路径
func (r *Router) formatPath(dirs []string) string {
	return "/" + strings.Join(dirs, "/")
}

//添加路径
func (r *Router) AddPath(path string, value interface{}) error {
	//写锁
	r.routerMutex.Lock()
	defer r.routerMutex.Unlock()

	//参数检查
	path = strings.TrimSpace(path)
	path = strings.ToLower(path)
	if path == "" { //不能为空
		return errors.New("len(path) cannot be equal to 0")
	}
	if path[0] != '/' { //必须以 '/' 开头
		return errors.New("the path must begin with '/'")
	}

	//路径切割
	dirs, e := r.splitPath(path)
	if e != nil {
		return e
	}

	//目录检查，TwoStar必须为path的结尾
	for i, dir := range dirs {
		if dir == TwoStar && i != len(dirs)-1 {
			return errors.New("** Must be at the end of the path")
		}
	}

	//路径格式化
	path = r.formatPath(dirs)

	//添加原始路径
	if _, ok := r.rawPath[path]; !ok {
		r.rawPath[path] = value
		if path != "/" {
			r.rawPath[path+"/"] = value
		}
	} else {
		return nil
	}

	//添加目录树
	p := r.dirTree
	for i, dir := range dirs {
		if dir == OneStar || dir == TwoStar {
			if len(p.subDir) != 0 {
				delete(r.rawPath, path)
				if path != "/" {
					delete(r.rawPath, path+"/")
				}
				return errors.New("dir is OneStar or TwoStar, but len(subDir) != 0")
			} else {
				switch dir {
				case OneStar:
					p.subDirHaveOneStar = true
				case TwoStar:
					p.subDirHaveTwoStar = true
				}
			}
		} else if p.subDirHaveOneStar || p.subDirHaveTwoStar {
			delete(r.rawPath, path)
			if path != "/" {
				delete(r.rawPath, path+"/")
			}
			return errors.New("subDir is OneStar or TwoStar")
		}
		if _, ok := p.subDir[dir]; !ok {
			p.subDir[dir] = newDirNode(dir, nil, false, p)
		}
		if i == len(dirs)-1 {
			p.subDir[dir].isPathEnd = true
			p.subDir[dir].value = value
		}
		p = p.subDir[dir]
		p.referenceNumber++
	}

	return nil
}

//查找路径
func (r *Router) SearchPath(path string) (interface{}, bool) {
	//读上锁
	r.routerMutex.RLock()

	//转换为小写
	path = strings.ToLower(path)

	//快速匹配
	if v, ok := r.rawPath[path]; ok {
		//读解锁
		r.routerMutex.RUnlock()
		return v, true
	}

	//路径切割
	dirs, e := r.splitPath(path)
	if e != nil {
		//读解锁
		r.routerMutex.RUnlock()
		return nil, false
	}

	//对 "/" 特殊处理
	if len(dirs) == 0 {
		if v, ok := r.rawPath["/"]; ok {
			//读解锁
			r.routerMutex.RUnlock()
			return v, true
		} else {
			return nil, false
		}
	}

	//目录树匹配
	p := r.dirTree
	for _, dir := range dirs {
		if _, ok := p.subDir[dir]; ok {
			p = p.subDir[dir]
		} else if p.subDirHaveTwoStar {
			//读解锁
			r.routerMutex.RUnlock()
			return p.subDir[TwoStar].value, true
		} else if p.subDirHaveOneStar {
			p = p.subDir[OneStar]
		} else {
			//读解锁
			r.routerMutex.RUnlock()
			return nil, false
		}
	}

	//返回
	if p.isPathEnd {
		//读解锁
		r.routerMutex.RUnlock()
		return p.value, true
	} else {
		//读解锁
		r.routerMutex.RUnlock()
		return nil, false
	}
}

//修改路径
func (r *Router) ModifyPath(path string, value interface{}) error {
	//写锁
	r.routerMutex.Lock()
	defer r.routerMutex.Unlock()

	//转换为全小写
	path = strings.ToLower(path)

	//参数过滤
	if path == "" {
		return errors.New("len(path) cannot be equal to 0")
	}

	//路径切割
	dirs, e := r.splitPath(path)
	if e != nil {
		return e
	}

	//格式化路径
	path = r.formatPath(dirs)

	//修改原始路径
	if _, ok := r.rawPath[path]; !ok {
		return errors.New("path is not exists")
	} else {
		r.rawPath[path] = value
		if path != "/" {
			r.rawPath[path+"/"] = value
		}
	}

	//修改目录树
	p := r.dirTree
	for i, dir := range dirs {
		if i == len(dirs)-1 {
			p.subDir[dir].value = value
		}
		p = p.subDir[dir]
	}

	return nil
}

//删除路径
func (r *Router) DelPath(path string) error {
	//写锁
	r.routerMutex.Lock()
	defer r.routerMutex.Unlock()

	//转换为全小写
	path = strings.ToLower(path)

	//参数过滤
	if path == "" {
		return errors.New("len(path) cannot be equal to 0")
	}

	//路径切割
	dirs, e := r.splitPath(path)
	if e != nil {
		return e
	}

	//格式化路径
	path = r.formatPath(dirs)

	//删除原始路径
	if _, ok := r.rawPath[path]; !ok {
		return errors.New("path is not exists")
	} else {
		delete(r.rawPath, path)
		if path != "/" {
			delete(r.rawPath, path+"/")
		}
	}

	//删除目录树
	p := r.dirTree
	for _, dir := range dirs {
		if p.subDir[dir].referenceNumber == 1 {
			delete(p.subDir, dir)
			switch dir {
			case OneStar:
				p.subDirHaveOneStar = false
			case TwoStar:
				p.subDirHaveTwoStar = false
			}
			break
		} else {
			p = p.subDir[dir]
			p.referenceNumber--
		}
	}

	return nil
}

//获取Path路径字典
func (r *Router) GetPathMap() map[string]struct{} {
	r.routerMutex.RLock()
	defer r.routerMutex.RUnlock()

	rtn := make(map[string]struct{})
	for k := range r.rawPath {
		rtn[k] = struct{}{}
	}

	return rtn
}
