/*
 * @Author: your name
 * @Date: 2022-01-13 21:06:08
 * @LastEditTime: 2022-01-13 21:11:07
 * @LastEditors: Please set LastEditors
 * @Description: 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 * @FilePath: /IAM/pkg/storage/storage.go
 */

package storage

import (
	"crypto/sha256"
	"fmt"
	"hash"

	"github.com/spaolacci/murmur3"
)

// Defines algorithm constant.
var (
	HashSha256    = "sha256"
	HashMurmur32  = "murmur32"
	HashMurmur64  = "murmur64"
	HashMurmur128 = "murmur128"
)

func hashFunction(algorithm string) (hash.Hash, error) {
	switch algorithm {
	case HashSha256:
		return sha256.New(), nil
	case HashMurmur32:
		return murmur3.New32(), nil
	case HashMurmur64:
		return murmur3.New64(), nil
	case HashMurmur128:
		return murmur3.New128(), nil
	default:
		return murmur3.New32(), fmt.Errorf("unknown key hash function: %s. Falling back to murmur32", algorithm)
	}
}

// HashStr return hash the give string and return.
func HashStr(in string) string {
	h, _ := hashFunction()
}
