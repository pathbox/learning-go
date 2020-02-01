package main 

import (
    "crypto/rand"
    "crypto/sha256"
    //"database/sql"
    "fmt"
    "io"
)

const saltSize = 12

func main() {
    email := []byte("john.deo@gmail.com")
    password := []byte("1234567890")

    salt := make([]byte, saltSize) // 创建随机salt用于加密
    _, err := io.ReadFull(rand.Reader, salt)
    if err != nil {
        panic(err) 
    }

    hash := sha256.New()
    hash.Write(password)
    hash.Write(salt)

    fmt.Println("email: ", email)
    fmt.Println("password: ", string(password))
    fmt.Println("salt: ", salt)

    fmt.Println("hash: ", hash.Sum(nil)) // 最后的密码存储的其实是hash的值，当需要验证密码时，需要将salt用于验证，所以salt一起存在了数据库记录中
                                        //  这样即使撞库得到了hash的字符串真正的显示，也很难知道真正的输入密码是什么，由于sha256加密是不可逆的过程
                                        // 所以，难以返回得到真正的输入密码。比单独使用md5加密password的方法有效

    /* stmt, err := db.Prepare("INSERT INTO accounts SET hash = ?, salt = ?, email = ?")
     if err != nil {
         panic(err) 
     }
     result, _ := stmt.Exec(hash, salt, email)
    */
}