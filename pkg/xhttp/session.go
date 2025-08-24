// @author xiangqian
// @date 2025/07/27 12:37
package xhttp

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"sort"
	"sync"
	"time"
)

// 读写互斥锁
var rwMutex sync.RWMutex

// 会话集
var sessions map[string]*Session

func init() {
	sessions = make(map[string]*Session)
}

// Expired 会话是否已过期
func Expired(r *http.Request) bool {
	// 获取会话
	session, err := GetSession(r)
	// 会话是否已过期
	return err != nil || session == nil || session.ExpiresAt.Before(time.Now())
}

// GetSession 获取会话
func GetSession(r *http.Request) (*Session, error) {
	// 获取读锁，允许多个读操作同时进行
	rwMutex.RLock()
	// 释放读锁
	defer rwMutex.RUnlock()

	// 获取会话 id
	id, err := GetCookie(r, "session_id")
	if err != nil {
		return nil, err
	}

	// 获取会话
	session, ok := sessions[id]
	if !ok {
		return nil, nil
	}
	return session, nil
}

// SetSession 设置会话
func SetSession(w http.ResponseWriter) error {
	// 获取写锁，阻塞其他所有读写操作
	rwMutex.Lock()
	// 释放写锁
	defer rwMutex.Unlock()

	// 生成会话id
	// 生成 16 字节（128 位）的随机数
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return err
	}
	id := base64.RawURLEncoding.EncodeToString(buf)

	// 设置会话
	var maxAge = 12 * 60 * 60 // 设置会话过期时间为 12 个小时
	session := &Session{
		Id:        id,
		ExpiresAt: time.Now().Add(time.Duration(maxAge) * time.Second),
	}
	sessions[id] = session

	// 设置 Cookie
	SetCookie(w, "session_id", id, maxAge)

	// 限制最大登录数
	const maxLen = 100
	var sLen = len(sessions)
	if sLen > maxLen {
		// 会话根据过期时间升序排序
		var sessionArr = make([]*Session, 0, sLen)
		for _, session = range sessions {
			sessionArr = append(sessionArr, session)
		}
		sort.Slice(sessionArr, func(i, j int) bool {
			return sessionArr[i].ExpiresAt.Before(sessionArr[j].ExpiresAt)
		})

		// 移除超过限制最大登录数的会话
		var i = 0
		for l := sLen - maxLen; i < l; i++ {
			session = sessionArr[i]
			delete(sessions, session.Id)
		}
		// 移除过期的会话
		for l := sLen; i < l; i++ {
			session = sessionArr[i]
			if session.ExpiresAt.Before(time.Now()) {
				delete(sessions, session.Id)
			} else {
				break
			}
		}
	}

	return nil
}

// DelSession 删除会话
func DelSession(w http.ResponseWriter, r *http.Request) error {
	// 获取写锁，阻塞其他所有读写操作
	rwMutex.Lock()
	// 释放写锁
	defer rwMutex.Unlock()

	// 获取会话 id
	id, err := GetCookie(r, "session_id")
	if err != nil {
		return err
	}

	// 删除会话
	delete(sessions, id)

	// 设置 Cookie
	SetCookie(w, "session_id", "", -1)

	return nil
}

// GetCookie 获取 Cookie
func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	value := cookie.Value
	if value != "" {
		value, _ = url.QueryUnescape(value)
	}
	return value, nil
}

// SetCookie 设置 Cookie
func SetCookie(w http.ResponseWriter, name, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,                   // Cookie 名称
		Value:    url.QueryEscape(value), // Cookie 值
		Path:     "/",                    // Cookie 有效路径，"/" 表示对整个网站有效
		HttpOnly: true,                   // 设置为 true 防止 JavaScript 通过 document.cookie 访问，增强安全性
		MaxAge:   maxAge,                 // Cookie 有效期（单位：秒），设置为正数表示多少秒后过期，设置为 0 表示立即删除 Cookie，设置为负数表示会话 Cookie（浏览器关闭后删除）
	})
}

// Session 会话
type Session struct {
	Id        string
	ExpiresAt time.Time
}
