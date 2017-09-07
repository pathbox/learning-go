func main() {
  client, err := redis.Connect("127.0.0.1:6379", 10)
  if err != nil {
    log.Println("Error Connectiong to redis")
  }

  if err = client.Ping(); err != nil {
    log.Println(err)
  }

  val, _ := client.Set("hc:req:0", "Adieu")
}