# GRPC Redis Client LoadBalancer

## Redis connection
There's no restriction for redis package.

You can create an struct which wrap the package, and implement the `resolver.Redis` interface.
## Discovery

``` golang
// create an redis connection which implements the resolver.Redis interface
var redisConn resolver.Redis

// registe schema
discovery.RegisterSchema(redisConn)

// discovery and dial the server
conn, err := grpc.Dial("redis://service-foo", grpc.WithInsecure())
```


## Registry

``` golang
// create an redis connection which implements the resolver.Redis interface
var redisConn resolver.Redis

ctx := context.Background()
_, err := registry.RegisterService(ctx, redisConn, "service-foo", 8080, registry.WithStaticProvider("127.0.0.1"))
```

## Listening addr providers
### Static IP
the option `registry.WithStaticProvider("127.0.0.1")` is to registry a static ip address into redis.

### AWS ECS
the option `registry.WithECSProvider()` is used when registry service which running in aws ECS, it will detect container IP through ECS metadata automatically.

