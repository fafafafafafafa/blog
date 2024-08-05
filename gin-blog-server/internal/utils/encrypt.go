package utils

import "golang.org/x/crypto/bcrypt"

// 加密
// 使用 bcrypt 对比 明文字符串 和 哈希值
func BcryptCheck(plain, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}
