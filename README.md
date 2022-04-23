# GRPC Redis Client LoadBalancer

## Redis connection
There's no restriction for redis package.

You can create a struct which wrap the redis package you're using, and implement the `resolver.Redis` interface.
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
The option `registry.WithStaticProvider("127.0.0.1")` is used to register the service with a known IP.

### AWS ECS
The option `registry.WithECSProvider()` is used to register the service which is running in aws ECS. It will detect container IP through ECS metadata automatically.

