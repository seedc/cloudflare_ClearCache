#!/bin/bash
# 清理实例
# curl -X POST -H "Content-Type: application/json" -d "{\"domain\":\"${1}\"}" http://cf-cache.default.svc.cluster.local:8000/api/v1/domain
curl -X POST -H "Content-Type: application/json" -d "{\"domain\":\"${1}\"}" http://127.0.0.1:8000/api/v1/domain
