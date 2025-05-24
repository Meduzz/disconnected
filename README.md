# disconnected
something's disconnected, for sure!

```go
    // gin wrapper
    start := disconnected.HttpServer("/service1", func(srv *Server) {
        srv.Quickapi("/api", Entity1, Entity2, Entity3)
        srv.Static("/static", "./public")
        srv.SPA("./public/index.html")
        srv.WithRouter(func(r *gin.Engine) {
            r.POST("/test", func(c *gin.Context) {
                c.JSON(200, gin.H{"message": "Hello World"})
            }
            // ...
        })
    })
```

```go
    // "RPC" server
    disconnected.RpcServer("service1", func(srv *Server) {
        srv.Quickapi("api", Entity1, Entity2, Entity3)
		s.HandleEvent("event", "service1", func(ec *rpc.EventContext) {})
		s.HandleRPC("rpc", "service1", func(rc *rpc.RpcContext) {})
		s.WithConn(func(r *nats.Conn) {
            // ...
		})
    })
```