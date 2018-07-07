package person

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type Person struct {
	Name string `redis:"name"`
	Age  int    `redis:"age"`
}

func RetrievePerson(conn redis.Conn, id string) (Person, error) {
	var person Person

	values, err := redis.Values(conn.Do("HGETALL", fmt.Sprintf("person:%s", id)))
	if err != nil {
		return person, err
	}

	//ScanStruct scans alternating names and values from src to a struct. The HGETALL and CONFIG GET commands return replies in this format.
	err = redis.ScanStruct(values, &person)
	return person, err

}
