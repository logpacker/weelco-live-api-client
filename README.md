#### Weelco Live API Go client

#### Usage

```go
// Live API addr
c := client.New("https://as.weelco.com/liveapi")
// Create sream
c.CreateStream(*client.Stream{
    Name: "test",
})
// Get stream by hash
s, _ := c.GetStream("HASH")
// Update output addr
s.OutputURL = "http://..."
s.OutputNodeAddr = "localhost"
s.Status = 1 // 0 - New, 1 - Running, 2 - Closed
c.UpdateStream(s)
```
