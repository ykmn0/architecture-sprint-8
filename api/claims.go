package main

import "github.com/golang-jwt/jwt/v5"

const (
	RealmAccess = "realm_access" // Key for accessing realm information
	Roles       = "roles"        // Key for accessing roles
)

// RolesClaims is a structure that holds JWT claims with a method to extract data as a map.
type RolesClaims struct {
	Claims jwt.MapClaims
}

type mapClaims map[string]interface{}
type sliceClaims []interface{}

// emptyMap is an empty map used for error handling or absence of data.
var emptyMap = map[string]interface{}{}

// ToMap method returns a mapClaims type extracted from Claims using the RealmAccess key.
func (r RolesClaims) ToMap() mapClaims {
	if ra, ok1 := r.Claims[RealmAccess]; !ok1 {
		return emptyMap
	} else if realmAccess, ok2 := ra.(map[string]interface{}); !ok2 {
		return emptyMap
	} else {
		return realmAccess
	}
}

// ToSlice method returns a sliceClaims type extracted from mapClaims using the Roles key.
func (mc mapClaims) ToSlice() sliceClaims {
	if a, ok1 := mc[Roles]; !ok1 {
		return []interface{}{}
	} else if roles, ok2 := a.([]interface{}); !ok2 {
		return []interface{}{}
	} else {
		return roles
	}
}

// ToRoles method converts sliceClaims into a set of roles represented as a map.
func (sc sliceClaims) ToRoles() map[string]struct{} {
	result := make(map[string]struct{})
	for _, value := range sc {
		if role, ok := value.(string); ok {
			result[role] = struct{}{}
		}
	}
	return result
}
