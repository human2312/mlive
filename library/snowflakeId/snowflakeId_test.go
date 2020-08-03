package snowflakeId

import (
	"log"
	"testing"
)

func TestGetSnowflakeId(t *testing.T)  {

	id := GetSnowflakeId()

	log.Println("id:",id)

}