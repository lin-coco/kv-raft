package acltool

import (
	"strings"
)

type Permission string

const (
    Read      Permission = "r"
    ReadWrite Permission = "rw"
	Deny      Permission = "deny"
)

func HasPermission(permission Permission) bool {
	if permission != Read &&
		 permission != ReadWrite && permission != Deny {
			return false
	}
	return true
}

type MatchType string

const (
    MatchPrefix  MatchType = "prefix"
    MatchSuffix  MatchType = "suffix"
    MatchContains MatchType = "contains"
    MatchExact   MatchType = "exact"
    // MatchRegex MatchType = "regex" // 可选扩展
)

type ACLRule struct {
	// 要匹配的值（比如 "*account:*"）
    Pattern    string `json:"pattern"`
	// 权限 r / rw / 
    Permission Permission `json:"permission"`
}

// 前缀（account:*）
// 后缀（*:id）
// 包含（*config*）
// 精确匹配（tenant:123）
func KeyPermission(key string, accountACLRule []ACLRule) Permission {
	for _, aclRule := range accountACLRule {
		pat := strings.ReplaceAll(aclRule.Pattern, "*", "")
		if strings.HasPrefix(key, pat) || 
			strings.HasSuffix(key, pat) ||
			strings.Contains(key, pat) ||
			key == pat {
			return aclRule.Permission
		}
	}
	return Deny
}