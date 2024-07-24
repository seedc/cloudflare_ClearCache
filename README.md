## cloudflare_ClearCache

### 清除所有
POST -H "Content-Type: application/json" -d '{"domain":"all"}' http://cf-cache.default.svc.cluster.local:8000/api/v1/domain

### 单域名清除
POST -H "Content-Type: application/json" -d '{"domain":"xxxx.ai"}' http://cf-cache.default.svc.cluster.local:8000/api/v1/domain

### 单域名清除
POST -H "Content-Type: application/json" -d '{"domain":"xxxx.com"}' http://cf-cache.default.svc.cluster.local:8000/api/v1/domain