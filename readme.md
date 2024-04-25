# Sliding Window Rate Limting with Go and Lua in Redis

Rate limiting is a technique used to control the rate of requests sent or received by a system. It is used to prevent abuse and to ensure that the system is not overwhelmed by too many requests. 
There are many ways to implement rate limiting; fixed window, sliding window, token bucket, etc.,
In this article, we will delve into the implementation of **sliding window** rate limiting and evaluate its effectiveness across three different setups:
- Using Go and Redis alone
- Integrating Go, Redis, and Lua scripting with the `EVAL` command
- Combining Go, Redis, and Lua scripting using the `EVALSHA` command

&nbsp;
## What is a Sliding Window?
A **sliding window** tracks requests in the current time window (now - window duration).
Say we were to limit users to 10 requests per minute (60 second window): 
```js
14:01:00 -> req #1: processed  
14:01:05 -> req #2: processed  
...   
14:01:40 -> req #10: processed  
14:01:45 -> req #11: rejected // current window (now - 60s = 14:00:45 to 14:01:45) already processed 10 reqs 
// All reqs before 14:02:00 (when req #1 becomes outside of the current window) will be rejected 
// Continuing...
14:01:55 -> req #12: rejected  
14:02:00 -> req #13: processed  
14:02:03 -> req #14: rejected (next available slot in window is at 14:02:05)  

And so on
```


&nbsp;
## Redis, Sorted Sets, and Lua - Interlude
*Skip this section if you are already familiar with Redis and Sorted Sets*
- **Redis** is an in-memory data store that can be used as a database, cache, and message broker, it is very fast and can be used to store and retrieve data quickly, perfect for rate limiting. 
- It is **single-threaded**, which means it can only execute one command at a time.
- It has many data structures, one of them is **Sorted Sets**. A sorted set is a collection of unique elements ordered by a score.
- Taking a class example, some common operations include:
  - Adding an element (e.g. a student) with a score: `ZADD key score member`
  ```py
  ZADD class 90 "John"
  ZADD class 85 "Jane"
  ZADD class 80 "Joe"
  ```
  
  - Getting the rank of an element: `ZRANK key member` (ascending order)
  ```py
  ZRANK class "John" # 2
  ZRANK class "Jane" # 1
  ZRANK class "Joe"  # 0
  ```
  
  - Getting the score of an element: `ZSCORE key member`
  ```py
  ZSCORE class "John" # 90
  ZSCORE class "Jane" # 85
  ZSCORE class "Joe"  # 80
  ```
  
  - Getting the elements within a rank range: `ZRANGE key start stop`
  ```py
  ZRANGE class 0 1 # Joe, Jane
  ```

  - Getting the elements within a score range: `ZRANGEBYSCORE key min max`
  ```py
  ZRANGEBYSCORE class 85 100 # Jane, John
  ```

- **Lua** is a lightweight, high-level programming language designed primarily for embedded use in applications. It is used in **Redis** to write scripts that can be executed atomically. 
- **Atomically**? This means is that the script will run without interruption from other commands, ensuring that the script is executed in its entirety. So if we're [1] checking if limit exceeded and [2] adding timestamp of req to sorted set, we can be sure that no other command will run in between, and that any subsequent requests will read the updated data.
- Why is this important? In the context of highly concurrent systems, where a user can make multiple requests at the same time, there's a chance that the rate limiting logic can be bypassed. Suppose a user still has one req left to be made in the current window, and they make two requests at the same time: 
```js
req #1: check limit -> one req left -> passed
req #2: check limit -> one req left -> passed
req #1: add timestamp 
req #2: add timestamp
```

In this case, both requests will be processed, even though the user has exceeded the limit. This is where Lua scripting comes in handy, we'll see how in the upcoming sections. Just keep in mind that we can integrate Lua in Redis using either the `EVAL` or `EVALSHA` command.


&nbsp;
## The Plan
- Create a HTTP server and 3 middlewares to handle the rate limiting logic:
  1. **RateLimit()**: Go and Redis alone (logs success in `nolua.log`)
  2. **RateLimitLua()**: Go, Redis, and Lua using `EVAL` (logs in `lua.log`)
  3. **RateLimitLuaSha**: Go, Redis, and Lua using `EVALSHA` (logs in `luasha.log`)
- Each middleware logs success of incoming requests (true if req processed, false if rate-limited)

- Create `CallConcurrentReqs()` to mimick user behaviour and make calls concurrently to the server 
  - Why concurrent? To simulate real-world scenarios where multiple users are making requests at the same time, maybe even multiple requests from same user. And mainly to showcase how `Lua` prevents race conditions
- Run `CallConcurrentReqs()` for each middleware, for a certain duration
- Analyze the logs and compare the results

