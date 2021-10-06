### WTF is it?

Oneliner to create go-redis client with env based parameters and repeat policy

### Env Vars:
 - REDIS_HOST | default redis:6379. If specified with no port it uses port from REDIS_PORT
 - REDIS_PORT | defaule 6379. Used if REDIS_HOST has no port specified
 - REDIS_PASSWORD | default blank
 - REDIS_DB | default 0

### How to use this shit?

    import "github.com/ekiyanov/redisclient"
    redisclient.NewRedisClient()

Or alternatively with a context with timeout

    ctx,_:=context.WithTimeout(10*time.Second)
    redisclient.NewRedisClientCtx(ctx)


