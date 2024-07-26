#!/bin/bash
# 清理实例
curl -X POST -H "Content-Type: application/json" -d "{\"domain\":\"${1}\"}" http://cf-cache.default.svc.cluster.local:8000/api/v1/domain