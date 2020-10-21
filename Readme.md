### WTF is it?

Oneliner to create go-redis client with env based parameters and repeat policy

### How to use this shit?

    import "github.com/ekiyanov/redisclient"
    redisclient.NewRedisClient()

Or alternatively with a context with timeout

    ctx,_:=context.WithTimeout(10*time.Second)
    redisclient.NewRedisClientCtx(ctx)


