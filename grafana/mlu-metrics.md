# Cluster MLU usage
```promql
sum((count(cnmon_mlu_temp{container_name=~".+"} or vector(0))-1) / count(cnmon_mlu_temp{uuid=~".+"}) or vector(0))*100
```

# Used MLUs
```promql
count(cnmon_mlu_board_utilization{container_name=~".+"} or vector(0))-1
```

# Total MLUs
```promql
count(cnmon_mlu_temp{uuid=~".+"} or vector(0)) - 1
```

# Cluster MLU utilization (avg)
```promql
sum((sum(cnmon_mlu_board_utilization{pod_name=~".*"}) / count(cnmon_mlu_board_utilization{pod_name=~".*"}) ) or vector(0))
```

# Cluster MLU virtual memory utilization(avg)
```promql
sum((sum(cnmon_virtual_memory_used{pod_name=~".+"}) / (sum(cnmon_virtual_memory_used) + sum(cnmon_virtual_memory_free))) or vector(0))
```

# Cluster MLU physical memory utilization(avg)
```promql
sum((sum(cnmon_physical_memory_free{pod_name=~".+"}) / (sum(cnmon_physical_memory_free) + sum(cnmon_physical_memory_used))) or vector(0))
```